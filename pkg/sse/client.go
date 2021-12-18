package sse // import "astuart.co/go-sse"

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/simplefxn/ripeIngest/gibson/logger"
	"github.com/simplefxn/ripeIngest/pkg/pb"
)

//SSE name constants
const (
	eName = "event"
	dName = "data"
)

var (
	//ErrNilChan will be returned by Notify if it is passed a nil channel
	ErrNilChan = fmt.Errorf("nil channel given")
)

//Client is the default client used for requests.
var Client = &http.Client{}

func liveReq(verb, uri string, body io.Reader) (*http.Request, error) {
	req, err := GetReq(verb, uri, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "text/event-stream")

	return req, nil
}

//Event is a go representation of an http server-sent event
type Event struct {
	Type string
	Data io.Reader
}

//GetReq is a function to return a single request. It will be used by notify to
//get a request and can be replaces if additional configuration is desired on
//the request. The "Accept" header will necessarily be overwritten.
var GetReq = func(verb, uri string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(verb, uri, body)
}

func clientConnect(uri string) (*http.Response, error) {

	req, err := liveReq("GET", uri, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting sse request: %v", err)
	}

	res, err := Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error performing request for %s: %v", uri, err)
	}
	return res, nil
}

func getEvent(br *bufio.Reader) (*pb.RIS_Message, error) {
	delim := []byte{':', ' '}
	currEvent := &pb.RIS_Message{}

	for {
		bs, err := br.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return nil, err
		}

		if len(bs) < 2 {
			continue
		}

		spl := bytes.Split(bs, delim)

		if len(spl) < 2 {
			continue
		}

		switch string(spl[0]) {
		case eName:
			currEvent.Type = string(bytes.TrimSpace(spl[1]))
		case dName:
			currEvent.Data = string(bytes.TrimSpace(spl[1]))
			return currEvent, nil
		}
		if err == io.EOF {
			return nil, err
		}
	}
}

func getEvents(br *bufio.Reader, evCh chan<- *pb.RIS_Message) error {

	for {
		currEvent, err := getEvent(br)
		if err != nil {
			logger.Log.Errorf("Error getting event", err.Error())
			return err
		}
		evCh <- currEvent
	}
}

func Start(ctx context.Context, uri string, evCh chan<- *pb.RIS_Message) {
	// Make a receive channel for getting messages from the http response
	recvChan := make(chan *pb.RIS_Message)
	ctxDone := false

	if evCh == nil {
		return
	}
	// Main goroutine, connect, fecth event , repeat
	go func() {
		for {
			if ctxDone {
				return
			}
			res, err := clientConnect(uri)
			if err != nil {
				logger.Log.Info("Client connect skip until next cycle.")
				continue
			}
			// GoRoutine that will listen for the context and close the response if the context
			// is closed
			go func(ctx context.Context, res *http.Response) {
				<-ctx.Done()
				ctxDone = true
				logger.Log.Info("Received context close, closing service side response")
				res.Body.Close()
			}(ctx, res)
			// Create bufio reader
			br := bufio.NewReader(res.Body)
			// Loop for all events and send them to the recv Channel
			err = getEvents(br, recvChan)
			// If the goRoutine context is dome
			if err != nil {
				logger.Log.Info("Error from getting events from connection, skip until next cycle")
				res.Body.Close()
				continue
			}
		}
	}()

outside:
	for {
		select {
		// If we receive a message, forward to outside
		case ris := <-recvChan:
			evCh <- ris
			// If context is done , exit
		case <-ctx.Done():
			logger.Log.Info("SSE client receive signal to stop, closing receive channel")
			close(recvChan)
			break outside
		}
	}
}
