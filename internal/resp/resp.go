// Reading and parsing RESP commands from a bufio.Reader.

package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func ReadCommand(reader *bufio.Reader) ([]string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	line = strings.TrimSpace(line)
	if len(line) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	if line[0] != '*' {
		return nil, fmt.Errorf("expected array prefix '*', got: %c", line[0])
	}

	arraySize, err := strconv.Atoi(line[1:])
	if err != nil {
		return nil, fmt.Errorf("invalid array size: %v", err)
	}

	var args []string

	for i := 0; i < arraySize; i++ {
		lenLine, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		lenLine = strings.TrimSpace(lenLine)
		if lenLine[0] != '$' {
			return nil, fmt.Errorf("expected bulk string prefix '$', got: %c", lenLine[0])
		}

		strLen, err := strconv.Atoi(lenLine[1:])
		if err != nil {
			return nil, fmt.Errorf("invalid bulk string length: %v", err)
		}

		buf := make([]byte, strLen+2)
		_, err = io.ReadFull(reader, buf)
		if err != nil {
			return nil, err
		}

		args = append(args, string(buf[:strLen]))
	}
	return args, nil
}
