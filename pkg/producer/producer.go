package producer

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	"github.com/Shopify/sarama"
	"github.com/ZeroNull7/risProducer/pkg/config"
	"github.com/ZeroNull7/risProducer/pkg/logger"
)

type RIS struct {
	RisEventProducer sarama.AsyncProducer
}

func New(conf config.Kafka) *RIS {
	risProducer := RIS{
		RisEventProducer: newProducer(conf),
	}

	return &risProducer
}

func (r RIS) Input() chan<- *sarama.ProducerMessage {
	return r.RisEventProducer.Input()
}

func (r RIS) Close() error {
	return r.RisEventProducer.Close()
}

func newProducer(conf config.Kafka) sarama.AsyncProducer {
	sConfig := sarama.NewConfig()
	tlsConfig := createTlsConfiguration(conf)
	if tlsConfig != nil {
		sConfig.Net.TLS.Enable = true
		sConfig.Net.TLS.Config = tlsConfig
	}
	sConfig.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
	sConfig.Producer.Compression = sarama.CompressionSnappy   // Compress messages
	sConfig.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms

	brokers := []string{fmt.Sprintf("%s:%d", conf.Host, conf.Port)}
	producer, err := sarama.NewAsyncProducer(brokers, sConfig)

	if err != nil {
		logger.Log.Fatalf("Failed to start Sarama producer:", err)
	}
	return producer
}

func createTlsConfiguration(conf config.Kafka) (t *tls.Config) {
	if conf.Cert != "" && conf.Key != "" && conf.CA != "" {
		cert, err := tls.LoadX509KeyPair(conf.Cert, conf.Key)
		if err != nil {
			logger.Log.Fatal(err)
		}

		caCert, err := os.ReadFile(conf.CA)
		if err != nil {
			logger.Log.Fatal(err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		t = &tls.Config{
			Certificates:       []tls.Certificate{cert},
			RootCAs:            caCertPool,
			InsecureSkipVerify: conf.VerifySSL,
		}
	}
	// will be nil by default if nothing is provided
	return t
}
