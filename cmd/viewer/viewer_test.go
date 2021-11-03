package viewer

import (
	"cmdboard/typefile"
	"encoding/json"
	"testing"
	"time"
)

func TestInitPages(t *testing.T) {
	var commands map[int]typefile.Command
	t.Run("with valid args", func(t *testing.T) {
		message := json.RawMessage(`
			{
	  			"1": {
	      			"Id": 1,
	      			"Name": "hoge",
	      			"Comment": "comment",
	      			"ParentId": 0,
	      			"IsDir": false
	  			}
			}
		`)
		if err := json.Unmarshal(message, &commands); err != nil {
			t.Error(err)
		}

		viewer := &Viewer{
			commands: commands,
		}

		viewer.initPages()

		time.Sleep(time.Microsecond)

		if len(viewer.tree.GetRoot().GetChildren()) != 1 {
			t.Errorf("The number of TreeNodes is not the intended value. nodes count: %d", len(viewer.tree.GetRoot().GetChildren()))
		}
	})
}
