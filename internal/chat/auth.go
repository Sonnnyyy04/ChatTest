package chat

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func GetCredentials() (string, string) {
	fmt.Print("Enter your username: ")
	reader := bufio.NewReader(os.Stdin)
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("Enter your password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	return username, password
}

func SendCredentials(conn net.Conn, username, password string) {
	SendMessage(conn, username)
	SendMessage(conn, password)
}
func Login(conn net.Conn) bool {
	serverResponse, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Print(serverResponse)
	if strings.Contains(serverResponse, "Invalid") {
		return false
	}
	return true
}
