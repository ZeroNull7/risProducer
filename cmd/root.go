/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/ZeroNull7/risProducer/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var opts config.Service

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "risProducer",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.risProducer.yaml)")

	rootCmd.PersistentFlags().StringVar(&opts.Ris.URL, "ris.url", "https://ris-live.ripe.net/v1/stream", "ris url to connect for sse")

	rootCmd.PersistentFlags().StringVar(&opts.Ris.ClientString, "ris.clientstring", "ripe-client", "ris url to connect for sse")

	rootCmd.PersistentFlags().BoolVar(&opts.Ris.LogUnknowns, "ris.logunknowns", true, "log unknown ris messages and types")

	rootCmd.PersistentFlags().StringVar(&opts.Kafka.Host, "kafka.broker", "cluster-kafka.kafka.svc", "host to connect or bind the socket")

	// Kafka.port Port use to connect to the kafka broker
	rootCmd.PersistentFlags().IntVar(&opts.Kafka.Port, "kafka.port", 9092, "Port used to connect to the broker service")

	// Kafka.Cert is the cert to use if TLS is enabled
	rootCmd.PersistentFlags().StringVar(&opts.Kafka.Cert, "kafka.cert", "", "server certificate to use for kafka connections, requires grpc_key, enables TLS")

	// Kafka.key is the key to use if TLS is enabled
	rootCmd.PersistentFlags().StringVar(&opts.Kafka.Key, "kafka.key", "", "server private key to use for kafka connections, requires grpc_cert, enables TLS")

	// Kafka.ca	 is the CA to use if TLS is enabled
	rootCmd.PersistentFlags().StringVar(&opts.Kafka.CA, "kafka.ca", "", "server CA to use for kafka connections, requires TLS, and enforces client certificate check")

	// Kafka.verifyssl Optional verify ssl certificates chain
	rootCmd.PersistentFlags().BoolVar(&opts.Kafka.VerifySSL, "kafka.verifyssl", true, "Optional verify ssl certificates chain")

	// Enable prometheus stats
	rootCmd.PersistentFlags().BoolVar(&opts.Metrics.Enable, "metrics.enable", true, "enable prometheus metrics")

	// Prometheus metrics port to listen to
	rootCmd.PersistentFlags().IntVar(&opts.Metrics.Port, "metrics.port", 8080, "Port to liste for prometheus metrics scraping")

	rootCmd.PersistentFlags().StringVar(&opts.Metrics.Path, "metrics.path", "/metrics", "Path for prometheus metrics scraping ")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".ripe" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".risProducer")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
