package utils

import (
	"cmdboard/typefile"
	"encoding/json"
	"io/ioutil"
	"log"
)

func LoadCommands() (commands map[int]typefile.Command, err error) {
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

func WriteCommand(commands map[int]typefile.Command) error {
	commandsJson, err := json.Marshal(commands)
	if err != nil {
		return err
	}
	ioutil.WriteFile("cmdboard.json", commandsJson, 0664)
	return nil
}
