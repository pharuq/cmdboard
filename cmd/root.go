package cmd

import (
	"cmdboard/cmd/utils"
	"cmdboard/cmd/viewer"
	"fmt"
	"log"
	"os"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

type OptionsForRoot struct {
	need_copy bool
}

var optionsForRoot = &OptionsForRoot{}
var rootCmd = &cobra.Command{
	Use:   "cmdboard",
	Short: "Displays a list of saved commands.",
	Long:  `Displays a list of saved commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := utils.LoadCommands()
		if err != nil {
			log.Fatal(err)
			panic(err)
		}

		viewer.View(c)

		fmt.Fprintln(os.Stdout, viewer.SelectedText())
		if optionsForRoot.need_copy {
			clipboard.WriteAll(viewer.SelectedText())
		}
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.Flags().BoolVarP(&optionsForRoot.need_copy, "copy", "c", false, "Copy the result to the clipboard.")
}
