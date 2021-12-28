package cmd

import (
	"cmdboard/cmd/utils"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Executes the initialization process.",
	Long:  "Executes the initialization process.",
	Run: func(cmd *cobra.Command, args []string) {
		filePath := utils.StoredFilePath()

		_, err := os.Stat(filePath)
		if os.IsNotExist(err) {
			fp, err := os.Create(filePath)
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
