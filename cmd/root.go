/*
Copyright © 2022 Mathieu GRZYBEK

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"go.uber.org/zap"

	"go-logs-archiver/internal/buffer"
	"go-logs-archiver/internal/lock"
	"go-logs-archiver/internal/consumer"
	"go-logs-archiver/internal/core"
	"go-logs-archiver/internal/producer"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-logs-archiver",
	Short: "Tool used to send JSON messages to a persistent backend",
	Long: `Reads the incoming messages from the configured consumer driver 
and send them to the backend.

For example: get messages from a kafka topic and send them to a S3 storage.`,
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

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.go-logs-archiver.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Logs management
	rootCmd.PersistentFlags().StringP("log-level", "l", viper.GetString("LOG_LEVEL"), "Level of verbosity (development, production).")
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

		// Search config in home directory with name ".go-logs-archiver" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".go-logs-archiver")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func configureLogger(cmd *cobra.Command) *zap.Logger {
	var logger *zap.Logger
	var err error
	var zapPreset string

	logLevel, _ := cmd.Flags().GetString("log-level")

	switch strings.ToLower(logLevel) {
	case "production":
		zapPreset = "production"
		logger, err = zap.NewProduction()
	default:
		zapPreset = "development"
		logger, err = zap.NewDevelopment()
	}

	if err == nil {
		logger.Sugar().Infof("Logger configured as %s", zapPreset)
		return logger
	}

	logger.Sugar().Panic(err)
	return nil
}

func configureConsumer(logger *zap.Logger, engine *core.Engine) core.MessagesConsumer {
	if viper.Get("consumer.type") == "console" {
		logger.Info("Consumer created")
		result, err := consumer.NewConsole(logger, engine)

		if err != nil {
			logger.Sugar().Panic(err)
		}

		return result
	}

	logger.Sugar().Panic("the given consumer type is not found.")
	return nil
}

func configureProducer(logger *zap.Logger) core.MessagesProducer {
	if viper.Get("producer.type") == "console" {
		logger.Info("Producer created")
		result, err := producer.NewConsole()

		if err != nil {
			logger.Sugar().Panic(err)
		}

		return result
	}

	logger.Sugar().Panic("the given producer type is not found.")
	return nil
}

func configureBuffer(logger *zap.Logger) core.MessagesBuffer {
	if viper.Get("buffer.type") == "memory" {
		result, err := buffer.NewMemoryBuffer(logger, viper.GetInt64("buffer.step"))

		if err != nil {
			logger.Sugar().Panic(err)
		}

		logger.Info("Buffer created")
		return result
	}

	logger.Sugar().Panic("the given buffer type is not found.")
	return nil
}

func configureLock(logger *zap.Logger) core.LockingSystem {
	if viper.Get("locking.type") == "local" {
		result, err := lock.NewLockingSystem(logger, nil)

		if err != nil {
			logger.Sugar().Panic(err)
		}

		logger.Info("Lock created")
		return result
	}

	logger.Sugar().Panic("the given locking type is not found.")
	return nil
}