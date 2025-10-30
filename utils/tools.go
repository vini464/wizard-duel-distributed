package utils

import (
	"bufio"
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
