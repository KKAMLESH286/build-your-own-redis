package main

import (
	"fmt"
	"os"

	"github.com/kkamlesh286/build-your-own-redis/internal/server"
)

func main() {
	err := server.Start("0.0.0.0:6380")
	if err != nil {
		fmt.Println("Something went wrong")
		os.Exit(1)
	}
}

// import (
// 	"bufio"
// 	"fmt"
// 	"strings"

// 	"github.com/kkamlesh286/build-your-own-redis/internal/resp"
// )

// func main() {
// 	rawInput := "*1\r\n$4\r\nPING\r\n"

// 	reader := bufio.NewReader(strings.NewReader(rawInput))

// 	args, err := resp.ReadCommand(reader)
// 	if err != nil {
// 		fmt.Println("Error reading command:", err)
// 		return
// 	}

// 	fmt.Println("Parsed command arguments:", args)
// }
