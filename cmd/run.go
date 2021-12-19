/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"sync"

	"github.com/ZeroNull7/risProducer/pkg/producer"
	"github.com/ZeroNull7/risProducer/pkg/signals"
	"github.com/ZeroNull7/risProducer/pkg/sse"
	"github.com/spf13/cobra"
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
		var wg sync.WaitGroup

		ris := producer.New(opts)

		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		// Stop channel , catch SIGINT and SIGTERM
		stopCh := signals.SetupSignalHandler()

		uri := "https://ris-live.ripe.net/v1/stream/?format=sse&client=ripe-client"

		// Open Server side Events for RIS message production
		wg.Add(1)
		go func() {
			sse.Start(ctx, uri, ris)
			wg.Done()
		}()

		select {
		case <-stopCh:
			cancel()
			break
		case <-ctx.Done():
			break
		}
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
