/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/ZeroNull7/risProducer/pkg/config"
)
var cfgFile string
var opts config.Kafka

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
	
	rootCmd.PersistentFlags().StringVar(&opts.Host, "grpc.host", "localhost", "host to connect or bind the socket")

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ripe.yaml)")

	// Kafka.port is the port to listen on for kafka. If not set or zero, don't listen.
	rootCmd.PersistentFlags().IntVar(&opts.Port, "kafka.port", 9091, "Port to listen on for kafka calls")

	// Kafka.Cert is the cert to use if TLS is enabled
	rootCmd.PersistentFlags().StringVar(&opts.Cert, "kafka.cert", "", "server certificate to use for kafka connections, requires grpc_key, enables TLS")

	// Kafka.key is the key to use if TLS is enabled
	rootCmd.PersistentFlags().StringVar(&opts.Key, "kafka.key", "", "server private key to use for kafka connections, requires grpc_cert, enables TLS")

	// Kafka.ca	 is the CA to use if TLS is enabled
	rootCmd.PersistentFlags().StringVar(&opts.CA, "kafka.ca", "", "server CA to use for kafka connections, requires TLS, and enforces client certificate check")

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
		viper.SetConfigName(".ripe")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

