package cmd

import (
	"errors"
	"log"
	"strings"

	"github.com/spf13/cobra"
)

type Options struct {
	dir     string
	comment string
}

var o = &Options{}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "register command",
	Long:  "register command",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("specify only one argument")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := LoadCommands()
		if err != nil {
			log.Fatal(err)
			return err
		}
		commands = c

		parentNode := Command{Id: 0}
		// [TODO]o.dirの最初と最後が"/"の場合は取り除く
		if o.dir != "" {
			dirs := strings.Split(o.dir, "/")

			for _, d := range dirs {
				parentNode = findOrCreateNode(d, parentNode.Id)
			}
		}

		id := getNewId()
		node := Command{
			Id:       id,
			Name:     args[0],
			Comment:  o.comment,
			ParentId: parentNode.Id,
			IsDir:    false,
		}
		commands[id] = node

		if err := writeCommand(); err != nil {
			log.Fatal(err)
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(&o.dir, "dir", "d", "", "Directory for storing commands")
	addCmd.Flags().StringVarP(&o.comment, "comment", "c", "", "Command's comment")
}

func findOrCreateNode(name string, parentId int) Command {
	node, ok := findNode(name, parentId)
	if !ok {
		id := getNewId()
		node = Command{
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

func findNode(name string, parentId int) (Command, bool) {
	for _, c := range commands {
		if c.Name == name && c.ParentId == parentId {
			return c, true
		}
	}
	return Command{}, false
}

func getNewId() int {
	a := []int{}
	for k, _ := range commands {
		a = append(a, k)
	}
	return maxInt(a) + 1
}

func maxInt(a []int) int {
	max := a[0]
	for i := 1; i < len(a); i++ {
		if max < a[i] {
			max = a[i]
		}
	}
	return max
}
