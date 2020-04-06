package uciengine

/////////////////////////////////////////////////////////////////////
// imports

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/easychessanimations/gochess/board"
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

func (eng *UciEngine) PrintUciOptions() {
	for _, uo := range UCI_OPTIONS {
		if uo.Kind == "button" {
			fmt.Printf("option name %s type button", uo.Name)
		} else if uo.Kind == "spin" {
			fmt.Printf("option name %s type spin default %d min %d max %d", uo.Name, uo.DefaultInt, uo.MinInt, uo.MaxInt)
		} else if uo.ValueKind == "int" {
			fmt.Printf("option name %s type %s default %d", uo.Name, uo.Kind, uo.DefaultInt)
		} else if uo.ValueKind == "bool" {
			fmt.Printf("option name %s type %s default %v", uo.Name, uo.Kind, uo.DefaultBool)
		} else if uo.ValueKind == "stringarray" {
			fmt.Printf("option name %s type %s default %s", uo.Name, uo.Kind, strings.Join(uo.DefaultStringArray, " "))
		} else {
			fmt.Printf("option name %s type %s default %s", uo.Name, uo.Kind, uo.Default)
		}

		fmt.Println()
	}
}

func (eng *UciEngine) Uci() {
	fmt.Printf("id name %s\n", eng.Name)
	fmt.Printf("id author %s\n", eng.Author)
	fmt.Println()

	eng.PrintUciOptions()

	fmt.Println("uciok")
}

func (eng *UciEngine) GetOptionByName(name string) *UciOption {
	for _, uo := range UCI_OPTIONS {
		if uo.Name == name {
			return &uo
		}
	}

	return nil
}

func (eng *UciEngine) SetOptionFromString(uo *UciOption, value string) {
	if uo.ValueKind == "string" {
		uo.Value = value
	} else if uo.ValueKind == "int" {
		i, err := strconv.ParseInt(value, 10, 32)

		if err == nil {
			uo.ValueInt = int(i)
		} else {
			fmt.Println("option value should be an integer")
		}
	} else if uo.ValueKind == "bool" {
		uo.ValueBool = false

		if value == "true" {
			uo.ValueBool = true
		}
	}
}

func (eng *UciEngine) SetOption(name string, value string) {
	uo := eng.GetOptionByName(name)

	if uo == nil {
		fmt.Println("unknown uci option")
	} else {
		eng.SetOptionFromString(uo, value)
	}
}

func (eng *UciEngine) SetOptionFromTokens(tokens []string) {
	parseName := true
	nameBuff := ""
	valueBuff := ""

	if tokens[0] != "name" {
		fmt.Println("missing option name")
		return
	}

	for _, token := range tokens[1:] {
		if parseName {
			if token == "value" {
				parseName = false
			} else {
				if nameBuff == "" {
					nameBuff = token
				} else {
					nameBuff += " " + token
				}
			}
		} else {
			if valueBuff == "" {
				valueBuff = token
			} else {
				valueBuff += " " + token
			}
		}
	}

	eng.SetOption(nameBuff, valueBuff)
}

func (eng *UciEngine) LogPrefixedContent(prefix string, content string) {
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		fmt.Println(prefix, line)
	}
}

func (eng *UciEngine) LogInfoContent(content string) {
	eng.LogPrefixedContent("info string", content)
}

func (eng *UciEngine) LogAnalysisInfo(content string) {
	eng.LogPrefixedContent("info", content)
}

func (eng *UciEngine) Init() {
	eng.Name = ENGINE_NAME
	eng.Description = ENGINE_DESCRIPTION
	eng.Author = ENGINE_AUTHOR

	eng.Board.Init(board.VARIANT_STANDARD)

	eng.Board.Reset()

	eng.Board.LogFunc = eng.LogInfoContent
	eng.Board.LogAnalysisInfoFunc = eng.LogAnalysisInfo
}

func (eng *UciEngine) Stop() {
	eng.Board.Stop()
}

func (eng *UciEngine) UciLoop() {
	reader := bufio.NewReader(os.Stdin)

	for {
		text, _ := reader.ReadString('\n')

		command := strings.Trim(text, "\r\n")

		if (command == "quit") || (command == "x") {
			break
		} else if command == "stop" {
			eng.Stop()
		} else if command == "uci" {
			eng.Uci()
		} else if command == "i" {
			eng.Interactive = true
			eng.Board.Print()
		} else {
			tokens := strings.Split(command, " ")

			command = tokens[0]

			args := []string{}

			if len(tokens) > 1 {
				args = tokens[1:]
			}

			if command == "setoption" {
				eng.SetOptionFromTokens(args)
			} else if eng.Interactive {
				go eng.Board.ExecCommand(command)
			}
		}
	}
}

/////////////////////////////////////////////////////////////////////
