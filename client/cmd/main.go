package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

var username string
func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Usage: go run client.go [global|group]")
	}
	mode := os.Args[1]

	if mode != "global" && mode != "group" {
		log.Fatalln("Invalid mode. Use 'global' or 'group'.")
	}

	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatalln("Error connecting:", err.Error())
	}
	defer conn.Close()
	log.Println("Connected to the chat server.")

	scanner := bufio.NewScanner(os.Stdin)

	// Ask for username until a unique one is provided
	for {
		fmt.Print("Enter your username: ")
		scanner.Scan()
		username = scanner.Text()

		// Send username to server
		conn.Write([]byte(username + "\n"))

		// Read server's response
		response, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print(response)

		if response != "Username already taken. Please enter a different username:\n" {
			break
		}
	}

	// Send chat mode to server
	conn.Write([]byte(mode + "\n"))

	// Start listening for messages from the server
	go listenForMessages(conn)

	// Start reading input from the user and send it to the server
	for {
		fmt.Printf("%s > ", username)
		if scanner.Scan() {
			message := scanner.Text() + "\n"
			_, err := conn.Write([]byte(message))
			// Clear the previouse log
			fmt.Print("\r\033[K") // This clears the line
			
            if err != nil {
				log.Println("Error sending message:", err.Error())
				return
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading from input:", err.Error())
		}
	}
}

func listenForMessages(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalln("Server disconnected.")
		}
		// Clear the current input prompt
		fmt.Print("\r\033[K") // This clears the line

		// Print the new message
		fmt.Print(message)

		// Print the message prompt again
		fmt.Printf("%s > ", username)
	}
}
