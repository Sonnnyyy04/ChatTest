package chat

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

var (
	clients   = make(map[net.Conn]string)
	mu        sync.Mutex
	broadcast = make(chan string)
	users     = map[string]string{"user1": "qwe1", "user2": "qwe2", "user3": "qwe3"}
)

func handleConn(c net.Conn) {
	defer c.Close()

	reader := bufio.NewReader(c)
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	if validPassword(username, password) {
		c.Write([]byte("Login successful!!!\n"))
	} else {
		c.Write([]byte("Invalid username or password\n"))
		return
	}

	name := username
	mu.Lock()
	clients[c] = name
	mu.Unlock()

	broadcast <- fmt.Sprintf("%s has joined the chat", name)

	for {
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)
		if message == "EXIT" {
			break
		}
		p2pMessage(c, message, name)
		if !strings.HasPrefix(message, "P2P") {
			broadcast <- fmt.Sprintf("%s: %s", name, message)
		}
	}
	mu.Lock()
	delete(clients, c)
	mu.Unlock()
	broadcast <- fmt.Sprintf("%s has left the chat", name)
}

func p2pMessage(c net.Conn, message string, name string) {
	if strings.HasPrefix(message, "P2P") {
		parts := strings.SplitN(message, " ", 3)
		if len(parts) < 3 {
			c.Write([]byte("Invalid format or not enough arguments. Use P2P <username> <message>\n"))
		}
		targetUserName := parts[1]
		privateMessage := parts[2]
		var targetConn net.Conn
		mu.Lock()
		for conn, user := range clients {
			if user == targetUserName {
				targetConn = conn
				break
			}
		}
		mu.Unlock()
		if targetConn != nil {
			_, err := fmt.Fprintf(targetConn, "Private message from %s: %s\n", name, privateMessage)
			if err != nil {
				fmt.Println("Error sending private message:", err)
			}
		} else {
			c.Write([]byte("User not found\n"))
		}
	}
}

func validPassword(username string, password string) bool {
	storedPass, exist := users[username]
	if exist && storedPass == password {
		return true
	}
	return false
}

func broadcastMessages() {
	for {
		msg := <-broadcast
		mu.Lock()
		for conn := range clients {
			_, err := fmt.Fprintf(conn, "%s\n", msg)
			if err != nil {
				fmt.Println("Error broadcasting message:", err)
			}
		}
		mu.Unlock()
	}
}

func StartServer(listener net.Listener) {
	go broadcastMessages()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConn(conn)
	}
}
