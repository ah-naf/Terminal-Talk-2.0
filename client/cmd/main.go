package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

var username string

// ANSI escape codes for text formatting
const (
	Reset       = "\033[0m"
	Bold        = "\033[1m"
	FgGreen     = "\033[32m"
	FgCyan      = "\033[36m"
	ClearScreen = "\033[H\033[2J"
	FgRed       = "\033[1;31m" // Bold Red
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Usage: go run client.go [global|group]")
	}
	mode := os.Args[1]

	if mode != "global" && mode != "group" {
		log.Fatalln("Invalid mode. Use 'global' or 'group'.")
	}

	// Display cool welcome message
	fmt.Println(ClearScreen) // Clear the screen
	fmt.Println(Bold + FgCyan + "Welcome to the Cool Chat Server!" + Reset)

	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatalln("Error connecting:", err.Error())
	}
	defer conn.Close()

	scanner := bufio.NewScanner(os.Stdin)

	// Send chat mode to server
	conn.Write([]byte(mode + "\n"))

	// Ask for username until a unique one is provided
	for {
		fmt.Print(FgCyan + "Enter your username: " + Reset)
		scanner.Scan()
		username = scanner.Text()

		// Send username to server
		conn.Write([]byte(username + "\n"))

		// Read server's response
		response, _ := bufio.NewReader(conn).ReadString('\n')
		response = strings.TrimSpace(response) // Trim the response
		if strings.Contains(response, "Username already taken.") {
			fmt.Println(FgRed + "Username already exist. Please enter a new username." + Reset)
		} else {
			fmt.Println(ClearScreen) // Clear the screen after successful username entry
			fmt.Println(FgGreen + "Your username is " + Bold + username + Reset)
			break
		}
	}

	// Start listening for messages from the server
	go listenForMessages(conn)

	// Start reading input from the user and send it to the server
	for {
		fmt.Printf("%s%s > %s", Bold, username, Reset)
		if scanner.Scan() {
			message := scanner.Text()

			// Clear the previous log and print user's own message
			fmt.Print("\r\033[K") // This clears the line

			_, err := conn.Write([]byte(message + "\n"))
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

		// Clear the current input prompt
		fmt.Print("\r\033[K") // This clears the line

		if err != nil {
			os.Exit(0)
		}

		// Print the new message
		fmt.Print(message)

		// Print the message prompt again
		fmt.Printf("%s%s%s > %s",Reset, Bold, username, Reset)
	}
}
