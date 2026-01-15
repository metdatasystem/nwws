package main

import (
	"github.com/joho/godotenv"
	nwws "github.com/metdatasystem/nwws/internal"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	envFile     string
	logLevelInt int
	logLevel    zerolog.Level = 1
	// The root command of our program
	rootCmd = &cobra.Command{
		Use:   "nwws",
		Short: "MDS ingest from the National Weather Wire Service Open Interface.",
		Run: func(cmd *cobra.Command, args []string) {
			nwws.NWWS(logLevel)
		},
	}
)

// Go, go, go
func main() {
	rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Bind our args to the command
	rootCmd.PersistentFlags().StringVar(&envFile, "env", "", "The env file to read.")
	rootCmd.PersistentFlags().IntVar(&logLevelInt, "log", 1, "The logging level to use.")
}

func initConfig() {
	setLogLevel()

	if envFile != "" {
		err := godotenv.Load(envFile)
		if err != nil {
			log.Info().Err(err).Msg("failed to load env file")
		}
	}
}

func setLogLevel() {
	logLevel = zerolog.Level(logLevelInt)
}
