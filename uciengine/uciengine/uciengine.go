package uciengine

/////////////////////////////////////////////////////////////////////
// imports

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// types

func (eng *UciEngine) WelcomeMessage() string {
	return fmt.Sprintf("%s %s by %s", eng.Name, eng.Description, eng.Author)
}

func (eng *UciEngine) PrintWelcomeMessage() {
	fmt.Println(eng.WelcomeMessage())
	fmt.Println()
}

func (eng *UciEngine) Uci() {
	fmt.Printf("id name %s\n", eng.Name)
	fmt.Printf("id author %s\n", eng.Author)
	fmt.Println()
	fmt.Println("uciok")
}

func (eng *UciEngine) UciLoop() {
	reader := bufio.NewReader(os.Stdin)

	for {
		text, _ := reader.ReadString('\n')

		command := strings.Trim(text, "\r\n")

		if (command == "quit") || (command == "x") {
			break
		}

		if command == "uci" {
			eng.Uci()
		}

		/*tokens := strings.Split(command, " ")

		command = tokens[0]

		args := []string{}

		if len(tokens) > 1 {
			args = tokens[1:]
		}*/
	}
}

/////////////////////////////////////////////////////////////////////
