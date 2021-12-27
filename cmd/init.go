package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Executes the initialization process.",
	Long:  "Executes the initialization process.",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := os.Stat("cmdboard.json")
		if os.IsNotExist(err) {
			fp, err := os.Create("cmdboard.json")
			if err != nil {
				fmt.Println(err)
				return
			}
			defer fp.Close()

			fp.WriteString("{}")
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
