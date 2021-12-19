package metrics

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ZeroNull7/risProducer/pkg/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	srv *http.Server
}

var (
	RisNotificationCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ris_producer_notification_events_total",
		Help: "The total number of ris notification messages received",
	})

	RisOpenCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ris_producer_open_events_total",
		Help: "The total number of ris events messages received",
	})

	RisPeerStateCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ris_producer_peer_state_events_total",
		Help: "The total number of ris peer state messages received",
	})

	RisUnknownsCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ris_producer_unknown_events_total",
		Help: "The total number of unknown messages received",
	})

	RisUpdateAnnouncementsCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ris_producer_update_announcements_events_total",
		Help: "The total number of update-announcement messages received",
	})

	RisUpdateWithdrawalsCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ris_producer_update_withdrawal_events_total",
		Help: "The total number of update-withdrawal messages received",
	})

	RisUpdateUnknownCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ris_producer_update_unknown_events_total",
		Help: "The total number of unknown update messages received",
	})

	RisMessageCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ris_producer_events_total",
		Help: "The total number of messages received",
	})
)

func New(conf config.Metrics) *Server {
	mux := http.DefaultServeMux

	mux.Handle(conf.Path, promhttp.Handler())

	srv := &Server{
		srv: &http.Server{
			Addr:    fmt.Sprintf(":%v", conf.Port),
			Handler: mux,
		},
	}

	return srv
}

func (s *Server) ListenAndServe() {
	s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
