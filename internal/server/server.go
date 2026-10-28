// TCP connection and command handling for the Redis-like server.

package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/kkamlesh286/build-your-own-redis/internal/resp"
	"github.com/kkamlesh286/build-your-own-redis/internal/store"
)

var valStore = store.NewStore()

func Start(port string) error {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}

	defer listener.Close()
	fmt.Printf("Server is listening on %s\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("failed to accept connection: %v\n", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		args, err := resp.ReadCommand(reader)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Client disconnected")
				return
			}
			fmt.Printf("Error reading command: %v\n", err)
			return
		}

		fmt.Printf("Received command: %v\n", args)
		// Here you would typically process the command and send a response back to the client.
		command := strings.ToUpper(args[0])

		switch command {
		case "PING":
			conn.Write([]byte("+PONG\r\n"))
		case "ECHO":
			if len(args) < 2 {
				conn.Write([]byte("-ERR wrong number of arguments for 'ECHO' command\r\n"))
			} else {
				fmt.Fprintf(conn, "+%s\r\n", args[1])
			}
		case "SET":
			if len(args) < 3 {
				conn.Write([]byte("-ERR wrong number of arguments for 'SET' command\r\n"))
			}
			key, value := args[1], args[2]
			var duration *time.Duration

			if len(args) >= 5 && strings.ToUpper(args[3]) == "EX" {
				sec, err := time.ParseDuration(args[4])
				if err != nil {
					conn.Write([]byte("-ERR invalid expiration time\r\n"))
					continue
				}
				d := time.Duration(sec) * time.Second
				duration = &d
			}
			valStore.Set(key, value, duration)
			conn.Write([]byte("+OK\r\n"))

		case "GET":
			if len(args) < 2 {
				conn.Write([]byte("-ERR wrong number of arguments for 'GET' command\r\n"))
				continue
			}
			key := args[1]
			value, exists := valStore.Get(key)
			if !exists {
				conn.Write([]byte("$-1\r\n"))
			} else {
				fmt.Fprintf(conn, "$%d\r\n%s\r\n", len(value), value)
			}
		case "DEL":
			if len(args) < 2 {
				conn.Write([]byte("-ERR wrong number of arguments for 'DEL' command\r\n"))
				continue
			}
			keys := args[1:]
			deletedCount := valStore.Delete(keys)
			fmt.Fprintf(conn, ":%d\r\n", deletedCount)

		default:
			conn.Write([]byte("-ERR unknown command\r\n"))
		}
	}
}
