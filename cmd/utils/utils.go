package utils

import (
	"cmdboard/typefile"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

func LoadCommands() (commands map[int]typefile.Command, err error) {
	bytes, err := ioutil.ReadFile(StoredFilePath())
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

func WriteCommand(commands map[int]typefile.Command) error {
	commandsJson, err := json.Marshal(commands)
	if err != nil {
		return err
	}
	ioutil.WriteFile(StoredFilePath(), commandsJson, 0664)
	return nil
}

func StoredFilePath() (filePath string) {
	filePath = os.Getenv("CMDBOARD_STORED_FILE_PATH")
	if filePath == "" {
		filePath = os.Getenv("HOME") + "/.cmdboard.json"
	}
	return
}
