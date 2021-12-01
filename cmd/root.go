package cmd

import (
	"cmdboard/cmd/utils"
	"cmdboard/cmd/viewer"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cmdboard",
	Short: "A brief description of your application",
	Long:  `TDB`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := utils.LoadCommands()
		if err != nil {
			log.Fatal(err)
			panic(err)
		}

		viewer.View(c)

		fmt.Fprintln(os.Stdout, viewer.SelectedText())
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
}
