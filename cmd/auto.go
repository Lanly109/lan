/*
Copyright © 2022 Lanly
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var autoCmd = &cobra.Command{
	Use:   "auto",
	Short: "one command to run all",
	Long:  `clean -> check -> valid -> md5`,
	Run: func(cmd *cobra.Command, args []string) {
		commands := []*cobra.Command{cleanCmd, checkCmd, validCmd, md5Cmd}

		reader := bufio.NewReader(os.Stdin)

		for _, command := range commands {
			fmt.Printf("\nPress Enter to run '%s', or type anything else to stop...\n", command.Name())
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading input:", err)
				return
			}

			// Trim any extra whitespace, and check if it's a plain Enter
			if strings.TrimSpace(input) != "" {
				fmt.Println("Stopping execution.")
				return
			}

			// Execute the command if Enter is pressed
			new_args := []string{command.Name()}
			if len(args) > 1 {
				new_args = append(new_args, args[1:]...)
			}
			rootCmd.SetArgs(new_args)
			if err := rootCmd.Execute(); err != nil {
				fmt.Printf("Error executing '%s': %v\n", command.Name(), err)
				return
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(autoCmd)
}
