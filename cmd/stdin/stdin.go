package stdin

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// Read reads all input from stdin
func Read(reader io.Reader) (string, error) {
	bufReader := bufio.NewReader(reader)
	var builder strings.Builder

	for {
		line, err := bufReader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				builder.WriteString(line) // Write the last line if it doesn't end with newline
				break
			}
			return "", err
		}
		builder.WriteString(line)
	}

	return builder.String(), nil
}

// HasData is true if stdin is not empty
func HasData() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}
