package uciengine

/////////////////////////////////////////////////////////////////////
// imports

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/easychessanimations/gochess/board"
	"github.com/easychessanimations/gochess/utils"
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
	for _, uo := range board.UCI_OPTIONS {
		uo.PrintUci()
	}
}

func (eng *UciEngine) Uci() {
	fmt.Printf("id name %s\n", eng.Name)
	fmt.Printf("id author %s\n", eng.Author)
	fmt.Println()

	eng.PrintUciOptions()

	fmt.Println("uciok")
}

func (eng *UciEngine) GetOptionByName(name string) *utils.UciOption {
	for _, uo := range board.UCI_OPTIONS {
		if uo.Name == name {
			return &uo
		}
	}

	return nil
}

func (eng *UciEngine) SetOption(name string, value string) {
	uo := eng.GetOptionByName(name)

	if uo == nil {
		fmt.Println("unknown uci option")
	} else {
		uo.SetFromString(value)
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

func (eng *UciEngine) GetUciOptionByNameWithDefault(name string, defaultUciOption utils.UciOption) utils.UciOption {
	uo := eng.GetOptionByName(name)

	if uo != nil {
		return *uo
	}

	return defaultUciOption
}

func (eng *UciEngine) Init() {
	eng.Name = ENGINE_NAME
	eng.Description = ENGINE_DESCRIPTION
	eng.Author = ENGINE_AUTHOR

	eng.Board.Init(utils.VARIANT_STANDARD)

	eng.Board.Reset()

	eng.Board.LogFunc = eng.LogInfoContent
	eng.Board.LogAnalysisInfoFunc = eng.LogAnalysisInfo
	eng.Board.GetUciOptionByNameWithDefaultFunc = eng.GetUciOptionByNameWithDefault
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
