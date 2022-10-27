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
	"log"

	"github.com/spf13/cobra"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate the given configuration",
	Long: `Create the consumer / producer processes and test the connections
against these resources.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runValidate(cmd); err != nil {
			log.Panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// validateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// validateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runValidate(cmd *cobra.Command) error {
	/*
	 * Ops options: logs
	 */
	logger := configureLogger(cmd)

	defer logger.Sync()
	logger.Info("Starting validate…")

	logger.Sugar().Debugf("producer: %v", configureProducer(logger))
	logger.Sugar().Debugf("buffer: %v", configureBuffer(logger))
	logger.Sugar().Debugf("lock: %v", configureLock(logger))
	logger.Sugar().Debugf("consumer: %v", configureConsumer(logger, nil))

	logger.Info("Validate finished")
	return nil
}
