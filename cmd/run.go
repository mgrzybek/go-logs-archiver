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
	"go-logs-archiver/internal/core"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start processing",
	Long:  `Start the service, consuming and producing the messages.`,
	Run:   func(cmd *cobra.Command, args []string) { runRun(cmd) },
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

func runRun(cmd *cobra.Command) {
	/*
	 * Ops options: logs
	 */
	logger := configureLogger(cmd)

	defer logger.Sync()
	logger.Info("Starting validate…")

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
	go engine.TriggerFlush()

	/*
	 * Consumer
	 */
	consumer := configureConsumer(logger, &engine)
	consumer.Run()
}
