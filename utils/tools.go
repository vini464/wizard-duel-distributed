package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func Input(title string) string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(title)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func ReadFile(filapath string) []byte {
	file_bytes, err := os.ReadFile(filapath)
	if err != nil {
		return nil
	}
	return file_bytes
}

func DecodeFromFile[T any](filepath string, receiver *T, basic T) error {
	data := ReadFile(filepath)
	var err error
	if data != nil {
		err = json.Unmarshal(data, receiver)
	}
	if err != nil {
		*receiver = basic
	}
	return err
}
