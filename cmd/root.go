package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/micmonay/keybd_event"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

type Command struct {
	Id       int
	Name     string
	Comment  string
	ParentId int
	IsDir    bool
}

var commands = map[int]Command{}

// var nodeByCommandId = map[int]*tview.TreeNode{}
var parentNodeById = map[int]*tview.TreeNode{}

var rootCmd = &cobra.Command{
	Use:   "cmdboard",
	Short: "A brief description of your application",
	Long:  `TDB`,
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := LoadCommands()
		if err != nil {
			log.Fatal(err)
			return err
		}
		commands = c

		// Create tree view
		rootDir := "."
		root := tview.NewTreeNode(rootDir).
			SetColor(tcell.ColorRed)
		tree := tview.NewTreeView().
			SetRoot(root).
			SetCurrentNode(root)

		outputText := ""
		for _, c := range commands {
			if c.ParentId == 0 {
				addNode(c, root)
			}
		}

		page := tview.NewPages()
		modal := tview.NewModal().
			SetText("Do you want to delete the command?").
			AddButtons([]string{"Delete", "Cancel"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				if buttonLabel == "Delete" {
					currentNode := tree.GetCurrentNode()
					currentNodeId := currentNode.GetReference().(int)
					parentNode := parentNodeById[currentNodeId]
					parentNode.RemoveChild(currentNode)
					delete(parentNodeById, currentNodeId)
					delete(commands, currentNodeId)
					writeCommand()
				}
				page.SwitchToPage("tree")
			})
		pages := page.
			AddPage("tree", tree, true, true).
			AddPage("modal", modal, true, false)

		tree.SetSelectedFunc(func(node *tview.TreeNode) {
			reference := node.GetReference()
			if reference == nil {
				return // Selecting the root node does nothing.
			}
			nodeColor := node.GetColor()
			if nodeColor == tcell.ColorGreen {
				if len(node.GetChildren()) == 0 {
					for _, c := range commands {
						if reference == c.ParentId {
							addNode(c, node)
						}
					}
				} else {
					// Collapse if visible, expand if collapsed.
					node.SetExpanded(!node.IsExpanded())
				}
			} else {
				outputText = node.GetText()

				err := exitCmd()
				if err != nil {
					panic(err)
				}
			}
		})

		tree.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyRune:
				switch event.Rune() {
				case 'e':
					return nil
				case 'd':
					page.SwitchToPage("modal")
					return nil
				}
			}
			return event
		})

		if err := tview.NewApplication().SetRoot(pages, true).Run(); err != nil {
			panic(err)
		}
		fmt.Fprintln(os.Stdout, outputText)

		return nil
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
}

func LoadCommands() (commands map[int]Command, err error) {
	bytes, err := ioutil.ReadFile("cmdboard.json")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if err := json.Unmarshal(bytes, &commands); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return commands, nil
}

func exitCmd() error {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		panic(err)
	}

	// For linux, it is very important to wait 2 seconds
	if runtime.GOOS == "linux" {
		time.Sleep(2 * time.Second)
	}

	// exit tree view
	kb.HasCTRL(true)
	kb.SetKeys(keybd_event.VK_C)
	err = kb.Launching()
	if err != nil {
		panic(err)
	}
	return nil
}

func addNode(c Command, parentNode *tview.TreeNode) {
	newNode := tview.NewTreeNode(c.Name).
		SetReference(c.Id)
	if c.IsDir {
		newNode.SetColor(tcell.ColorGreen)
	}
	parentNode.AddChild(newNode)
	// nodeByCommandId[c.Id] = newNode
	parentNodeById[c.Id] = parentNode
}

func writeCommand() error {
	commandsJson, err := json.Marshal(commands)
	if err != nil {
		return err
	}
	ioutil.WriteFile("cmdboard.json", commandsJson, 0664)
	return nil
}
