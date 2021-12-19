/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"

	"github.com/ZeroNull7/risProducer/pkg/logger"
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
		logger.RegisterLog()

		ctx := context.Background()

		uri := fmt.Sprintf("%s/?format=sse&client%s", opts.Ris.URL, opts.Ris.ClientString)

		var p fastjson.Parser

		sse.Start(ctx, uri, func(ris *sse.SSE_RIS) {
			//logger.Log.Infof("Got this back %v", ris.Type)
			v, err := p.Parse(ris.Data)
			if err != nil {
				logger.Log.Errorf(err.Error())
			} else {

				switch kind := string(v.GetStringBytes("type")); {
				case kind == "UPDATE":
					if !v.Exists("announcements") && !v.Exists("withdrawals") {
						if opts.Ris.LogUnknowns {
							logger.Log.Infof("Unknown update %v", ris.Data)
						}
					}
				case kind == "RIS_PEER_STATE":
					break
				case kind == "OPEN":
					break
				case kind == "NOTIFICATION":
					break
				case kind == "":
					break
				default:
					if opts.Ris.LogUnknowns {
						logger.Log.Infof("Unknown event \"%s\"\n %v", kind, ris.Data)
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
