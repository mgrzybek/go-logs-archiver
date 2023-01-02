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

// Package cmd manages the CLI
package cmd

import (
	"log"

	"go-logs-archiver/internal/core"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// flushCmd represents the flush command
var flushCmd = &cobra.Command{
	Use:   "flush",
	Short: "Flush the configured buffer into the persistent storage",
	Long:  `Create the buffer / producer processes and triggers the flushing procedure.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runFlush(cmd); err != nil {
			log.Panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(flushCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// validateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// validateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runFlush(cmd *cobra.Command) error {
	/*
	 * Ops options: logs
	 */
	logger := configureLogger(cmd)

	defer logger.Sync()
	logger.Info("Starting flush…")

	/*
	 * Cannot work with memory buffer
	 */
	if viper.Get("buffer.type") == "memory" {
		logger.Panic("Cannot flush using a memory buffer")
	}

	/*
	 * Core
	 */
	engine, err := core.NewEngine(
		logger,
		configureProducer(logger),
		configureBuffer(logger),
		configureLock(logger),
	)
	if err != nil {
		logger.Sugar().Panic(err)
	}

	engine.FlushBuffer()
	return nil
}
