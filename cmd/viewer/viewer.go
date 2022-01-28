package viewer

import (
	"cmdboard/cmd/utils"
	"cmdboard/typefile"
	"sort"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	parentNodeById = map[int]*tview.TreeNode{}
	selectedText   string
)

type Viewer struct {
	commands    map[int]typefile.Command
	app         *tview.Application
	pages       *tview.Pages
	flex        *tview.Flex
	tree        *tview.TreeView
	commentView *tview.TextView
	helpView    *tview.TextView
	modal       *tview.Modal
}

func View(commands map[int]typefile.Command) {
	app := tview.NewApplication()
	viewer := &Viewer{
		app:      app,
		commands: commands,
	}

	viewer.initPages()

	if err := app.SetRoot(viewer.pages, true).Run(); err != nil {
		panic(err)
	}
}

func SelectedText() string {
	return selectedText
}

func (v *Viewer) initPages() {
	v.pages = tview.NewPages()

	v.initTree()
	v.initModal()
	v.initHelpView()

	v.pages.
		AddPage("tree", v.flex, true, true).
		AddPage("modal", v.modal, true, false).
		AddPage("help", v.helpView, true, false)
}

func (v *Viewer) initTree() {
	v.initTreeView()

	v.flex = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(v.tree, 0, 3, true).
		AddItem(v.commentView, 0, 1, false)
}

func (v *Viewer) initTreeView() {
	rootDir := "."
	root := tview.NewTreeNode(rootDir).
		SetColor(tcell.ColorRed)

	v.tree = tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	go sortedFor(v.commands, func(c typefile.Command) {
		if c.ParentId == 0 {
			addNode(c, root)
		}
	})

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
				go sortedFor(v.commands, func(c typefile.Command) {
					if reference == c.ParentId {
						addNode(c, node)
					}
				})
			} else {
				// Collapse if visible, expand if collapsed.
				node.SetExpanded(!node.IsExpanded())
			}
		} else {
			selectedText = node.GetText()
			v.app.Stop()
		}
	})

	v.tree.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'e':
				v.app.Suspend(v.editMode)
				return nil
			case 'd':
				v.pages.SwitchToPage("modal")
				return nil
			case 'h':
				v.pages.SwitchToPage("help")
				return nil
			case 'q':
				v.app.Stop()
			}
		}
		return event
	})
}

func (v *Viewer) editMode() {
	node := v.tree.GetCurrentNode()
	c := v.getCommandfromNode(node)
	newName, newComment, err := Edit(c)
	if err != nil {
		panic(err)
	}
	c.Name = newName
	c.Comment = newComment
	v.commands[c.Id] = c
	if err := utils.WriteCommand(v.commands); err != nil {
		panic(err)
	}
	node.SetText(newName)
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

func (v *Viewer) initHelpView() {
	v.helpView = tview.NewTextView().
		SetText(`
  /////////////////////////////
 /     This is help view     /
/////////////////////////////

Enter: Select / Expand directory
k: Up
j: Down
e: Edit command
d: Delete command
h: Open help view
q: Quit`)
	v.helpView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'q':
				v.pages.SwitchToPage("tree")
				return nil
			}
		}
		return event
	})
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

func sortedFor(m map[int]typefile.Command, f func(typefile.Command)) {
	var keys []int
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		f(m[k])
	}
}

func (v *Viewer) getCommandfromNode(node *tview.TreeNode) typefile.Command {
	reference := node.GetReference().(int)
	return v.commands[reference]
}
