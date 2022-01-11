package viewer

import (
	"bufio"
	"bytes"
	"cmdboard/cmd/utils"
	"cmdboard/typefile"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const EDIT_FILE_NAME = ".cmdboard_edit"
const COMMAND_HEADER_STRING = "# |<----  Command  ---->|(Don't erase this line)"
const COMMENT_HEADER_STRING = "# |<----  Comment  ---->|(Don't erase this line)"

func Edit(c typefile.Command) (string, string, error) {
	fPath := utils.GetFilePath(EDIT_FILE_NAME)
	message := COMMAND_HEADER_STRING + "\n" + c.Name + "\n\n" + COMMENT_HEADER_STRING + "\n" + c.Comment

	err := makeFile(fPath, message)
	if err != nil {
		fmt.Fprintf(os.Stdout, "failed make edit file. %s\n", err.Error())
		return "", "", err
	}
	defer deleteFile(fPath)

	err = openEditor("vim", fPath)
	if err != nil {
		fmt.Fprintf(os.Stdout, "failed open text editor. %s\n", err.Error())
		return "", "", err
	}

	content, err := ioutil.ReadFile(fPath)
	if err != nil {
		fmt.Fprintf(os.Stdout, "failed read content. %s\n", err.Error())
		return "", "", err
	}

	// Parce read content
	reader := bytes.NewReader(content)
	output_name, output_comment, err := editExtraction(reader)
	if err != nil {
		fmt.Fprintf(os.Stdout, "failed parce content. %s\n", err.Error())
		return "", "", err
	}

	return output_name, output_comment, nil
}

func makeFile(fPath, message string) (err error) {
	return ioutil.WriteFile(fPath, []byte(message), 0644)
}

func deleteFile(fPath string) error {
	return os.Remove(fPath)
}

func openEditor(program string, args ...string) error {
	c := exec.Command(program, args...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

func editExtraction(reader io.Reader) (name, comment string, err error) {
	var nameParts, commentParts []string

	r := regexp.MustCompile(`\S`)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if (line == COMMAND_HEADER_STRING) || (line == COMMENT_HEADER_STRING) {
			continue
		}

		if len(nameParts) == 0 && r.MatchString(line) {
			nameParts = append(nameParts, line)
		} else {
			commentParts = append(commentParts, line)
		}
	}

	if err = scanner.Err(); err != nil {
		return
	}

	name = strings.Join(nameParts, " ")
	name = strings.TrimSpace(name)

	comment = strings.Join(commentParts, "\n")
	comment = strings.TrimSpace(comment)

	return
}
