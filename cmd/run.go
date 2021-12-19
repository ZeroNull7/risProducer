/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/ZeroNull7/risProducer/pkg/logger"
	"github.com/ZeroNull7/risProducer/pkg/metrics"
	"github.com/ZeroNull7/risProducer/pkg/producer"
	"github.com/ZeroNull7/risProducer/pkg/signals"
	"github.com/ZeroNull7/risProducer/pkg/sse"
	"github.com/spf13/cobra"
	"github.com/valyala/fastjson"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var p fastjson.Parser
		var server *metrics.Server

		logger.RegisterLog()

		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)

		// Stop channel , catch SIGINT and SIGTERM
		sigStop := signals.SetupSignalHandler()

		// Prometheus metrics
		if opts.Metrics.Enable {
			server := metrics.New(opts.Metrics)
			go server.ListenAndServe()
		}

		go func() {
			<-sigStop
			if server != nil {
				shutCtx := context.Background()
				server.Shutdown(shutCtx)
			}
			cancel()
		}()

		uri := fmt.Sprintf("%s/?format=sse&client%s", opts.Ris.URL, opts.Ris.ClientString)
		risProducer := producer.New(opts.Kafka)

		sse.Start(ctx, uri, func(event *sse.SSE_RIS) {
			v, err := p.Parse(event.Data)
			if err != nil {
				logger.Log.Errorf(err.Error())
			} else {
				metrics.RisMessageCounter.Inc()
				switch kind := string(v.GetStringBytes("type")); {
				case kind == "UPDATE":
					if !v.Exists("announcements") && !v.Exists("withdrawals") {
						if opts.Ris.LogUnknowns {
							logger.Log.Infof("Unknown update %v", event.Data)
						}
						metrics.RisUpdateUnknownCounter.Inc()
					} else if v.Exists("announcements") {
						metrics.RisUpdateAnnouncementsCounter.Inc()
						risProducer.Input() <- &sarama.ProducerMessage{
							Topic: "ris-update-announcement",
							Value: sarama.ByteEncoder(event.Data),
						}
					} else if v.Exists("withdrawals") {
						metrics.RisUpdateWithdrawalsCounter.Inc()
						risProducer.Input() <- &sarama.ProducerMessage{
							Topic: "ris-update-withdrawal",
							Value: sarama.ByteEncoder(event.Data),
						}
					}
				case kind == "RIS_PEER_STATE":
					metrics.RisPeerStateCounter.Inc()
					risProducer.Input() <- &sarama.ProducerMessage{
						Topic: "ris-peer-state",
						Value: sarama.ByteEncoder(event.Data),
					}
				case kind == "OPEN":
					metrics.RisOpenCounter.Inc()
					risProducer.Input() <- &sarama.ProducerMessage{
						Topic: "ris-open",
						Value: sarama.ByteEncoder(event.Data),
					}
				case kind == "NOTIFICATION":
					metrics.RisNotificationCounter.Inc()
					risProducer.Input() <- &sarama.ProducerMessage{
						Topic: "ris-notification",
						Value: sarama.ByteEncoder(event.Data),
					}
				case kind == "":
					metrics.RisUnknownsCounter.Inc()
					risProducer.Input() <- &sarama.ProducerMessage{
						Topic: "ris-unknowns",
						Value: sarama.ByteEncoder(event.Data),
					}
				default:
					metrics.RisUnknownsCounter.Inc()
					if opts.Ris.LogUnknowns {
						logger.Log.Infof("Unknown event \"%s\"\n %v", kind, event.Data)
					}

				}
			}
		})
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
