package viewer

import (
	"cmdboard/cmd/utils"
	"cmdboard/typefile"
	"runtime"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/micmonay/keybd_event"
	"github.com/rivo/tview"
)

var parentNodeById = map[int]*tview.TreeNode{}
var selectedText string

type Viewer struct {
	commands    map[int]typefile.Command
	pages       *tview.Pages
	flex        *tview.Flex
	tree        *tview.TreeView
	commentView *tview.TextView
	modal       *tview.Modal
}

func View(commands map[int]typefile.Command) {
	viewer := &Viewer{
		commands: commands,
	}

	viewer.initPages()

	if err := tview.NewApplication().SetRoot(viewer.pages, true).Run(); err != nil {
		panic(err)
	}
}

func SelectedText() string {
	return selectedText
}

func (v *Viewer) initPages() {
	v.pages = tview.NewPages()

	v.initFlex()
	v.initModal()

	v.pages.
		AddPage("tree", v.flex, true, true).
		AddPage("modal", v.modal, true, false)
}

func (v *Viewer) initFlex() {
	v.initTree()

	v.flex = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(v.tree, 0, 3, true).
		AddItem(v.commentView, 0, 1, false)
}

func (v *Viewer) initTree() {
	rootDir := "."
	root := tview.NewTreeNode(rootDir).
		SetColor(tcell.ColorRed)

	v.tree = tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	for _, c := range v.commands {
		if c.ParentId == 0 {
			addNode(c, root)
		}
	}

	v.commentView = tview.NewTextView()
	v.commentView.SetBorder(true)
	v.tree.SetChangedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference == nil {
			return // Selecting the root node does nothing.
		}
		v.commentView.SetText(v.commands[reference.(int)].Comment)
	})

	v.tree.SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference == nil {
			return // Selecting the root node does nothing.
		}
		nodeColor := node.GetColor()
		if nodeColor == tcell.ColorGreen {
			if len(node.GetChildren()) == 0 {
				for _, c := range v.commands {
					if reference == c.ParentId {
						addNode(c, node)
					}
				}
			} else {
				// Collapse if visible, expand if collapsed.
				node.SetExpanded(!node.IsExpanded())
			}
		} else {
			selectedText = node.GetText()

			err := exitCmd()
			if err != nil {
				panic(err)
			}
		}
	})

	v.tree.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'e':
				return nil
			case 'd':
				v.pages.SwitchToPage("modal")
				return nil
			}
		}
		return event
	})
}

func (v *Viewer) initModal() {
	v.modal = tview.NewModal().
		SetText("Do you want to delete the command?").
		AddButtons([]string{"Delete", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Delete" {
				currentNode := v.tree.GetCurrentNode()
				currentNodeId := currentNode.GetReference().(int)
				parentNode := parentNodeById[currentNodeId]
				parentNode.RemoveChild(currentNode)
				delete(parentNodeById, currentNodeId)
				delete(v.commands, currentNodeId)
				utils.WriteCommand(v.commands)
			}
			v.pages.SwitchToPage("tree")
		})
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

func addNode(c typefile.Command, parentNode *tview.TreeNode) {
	newNode := tview.NewTreeNode(c.Name).
		SetReference(c.Id)
	if c.IsDir {
		newNode.SetColor(tcell.ColorGreen)
	}
	parentNode.AddChild(newNode)
	parentNodeById[c.Id] = parentNode
}
