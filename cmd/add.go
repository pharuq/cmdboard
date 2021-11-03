package cmd

import (
	"cmdboard/cmd/utils"
	"cmdboard/typefile"
	"errors"
	"log"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

type OptionsForAdd struct {
	dir     string
	comment string
}

var (
	optionsForAdd = &OptionsForAdd{}
	commands      = map[int]typefile.Command{}
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "register command to tree",
	Long:  "register command to tree",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("specify only one argument")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := utils.LoadCommands()
		if err != nil {
			log.Fatal(err)
			return err
		}
		commands = c

		parentNode := typefile.Command{Id: 0}
		if optionsForAdd.dir != "" {
			dirs := strings.Split(removedExstraCharForDirs(optionsForAdd.dir), "/")

			for _, d := range dirs {
				parentNode = findOrCreateNode(d, parentNode.Id)
			}
		}

		id := getNewId()
		node := typefile.Command{
			Id:       id,
			Name:     args[0],
			Comment:  optionsForAdd.comment,
			ParentId: parentNode.Id,
			IsDir:    false,
		}
		commands[id] = node

		if err := utils.WriteCommand(commands); err != nil {
			log.Fatal(err)
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(&optionsForAdd.dir, "dir", "d", "", "Specify a directory for storing command")
	addCmd.Flags().StringVarP(&optionsForAdd.comment, "comment", "c", "", "Specify a comment for command")
}

func findOrCreateNode(name string, parentId int) typefile.Command {
	node, ok := findNode(name, parentId)
	if !ok {
		id := getNewId()
		node = typefile.Command{
			Id:       id,
			Name:     name,
			Comment:  "",
			ParentId: parentId,
			IsDir:    true,
		}
		commands[id] = node
	}
	return node
}

func findNode(name string, parentId int) (typefile.Command, bool) {
	for _, c := range commands {
		if c.Name == name && c.ParentId == parentId {
			return c, true
		}
	}
	return typefile.Command{}, false
}

func getNewId() int {
	a := []int{}
	for k := range commands {
		a = append(a, k)
	}
	return maxInt(a) + 1
}

func maxInt(a []int) int {
	if len(a) == 0 {
		return 1
	}

	max := a[0]
	for i := 1; i < len(a); i++ {
		if max < a[i] {
			max = a[i]
		}
	}
	return max
}

func removedExstraCharForDirs(s string) string {
	regexp := regexp.MustCompile("(^/|/$)")
	return regexp.ReplaceAllString(s, "")
}
