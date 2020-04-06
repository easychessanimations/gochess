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

func (eng *UciEngine) GetOptionByName(name string) (int, *utils.UciOption) {
	for index, uo := range board.UCI_OPTIONS {
		if uo.Name == name {
			return index, &uo
		}
	}

	return -1, nil
}

func (eng *UciEngine) SetOption(name string, value string) {
	index, uo := eng.GetOptionByName(name)

	if uo == nil {
		fmt.Println("unknown uci option")
	} else {
		uo.SetFromString(value)
		board.UCI_OPTIONS[index] = *uo

		// handle UCI_Variant as a special option
		if name == "UCI_Variant" {
			eng.Board.ResetVariantFromUciOption()
		}
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
	_, uo := eng.GetOptionByName(name)

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

func (eng *UciEngine) Position(args []string) {
	var i int

	if len(args) <= 0 {
		fmt.Println("missing argument for position command")
	} else if args[0] == "startpos" {
		eng.Board.ResetVariantFromUciOption()
		args = args[1:]
	} else if args[0] == "fen" {
		if len(args) > 1 {
			fenparts := []string{}
			for i = 1; i < len(args); i++ {
				if args[i] == "moves" {
					break
				} else {
					fenparts = append(fenparts, args[i])
				}
			}
			eng.Board.SetFromVariantUciOptionAndFen(strings.Join(fenparts, " "))
			if i >= (len(args) - 1) {
				return
			} else {
				args = args[i+1:]
			}
		} else {
			fmt.Println("no fen specified in fen argument")
		}
	} else {
		fmt.Println("unknown position command")
	}

	if len(args) > 0 {
		if args[0] == "moves" {
			if len(args) > 1 {
				for i := 1; i < len(args); i++ {
					eng.Board.MakeAlgebMove(args[i], board.ADD_SAN)
				}
			} else {
				fmt.Println("no move list specified in moves argument")
			}
		} else {
			fmt.Printf("unknown position argument %s\n", args[0])
		}
	}
}

func (eng *UciEngine) ListUci() {
	for _, uo := range board.UCI_OPTIONS {
		fmt.Println(uo.ToString())
	}
}

func (eng *UciEngine) Go(args []string) {
	depth := board.DEFAULT_SEARCH_DEPTH

	for len(args) > 0 {
		if args[0] == "depth" {
			if len(args) > 1 {
				parseDepth, err := strconv.ParseInt(args[1], 10, 32)
				if err != nil {
					fmt.Printf("invalid depth argument to go command, assuming depth %d\n", depth)
				} else {
					depth = int(parseDepth)
				}
				args = args[2:]
			} else {
				fmt.Printf("empty depth argument to go command, assuming depth %d\n", depth)
				args = args[1:]
			}
		} else {
			args = args[1:]
		}
	}

	go eng.Board.Go(depth)
}

func (eng *UciEngine) ExecuteUciCommand(command string) {
	if (command == "stop") || (command == "s") {
		eng.Stop()
	} else if command == "uci" {
		eng.Uci()
	} else if command == "l" {
		eng.ListUci()
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
		} else if command == "position" {
			eng.Position(args)
		} else if command == "go" {
			eng.Go(args)
		} else {
			go eng.Board.ExecCommand(command)
		}
	}
}

func (eng *UciEngine) ExecuteUciCommands(commandsStr string) {
	commands := strings.Split(commandsStr, "\n")

	for _, command := range commands {
		eng.ExecuteUciCommand(command)
	}
}

func (eng *UciEngine) UciLoop() {
	reader := bufio.NewReader(os.Stdin)

	for {
		text, _ := reader.ReadString('\n')

		command := strings.Trim(text, "\r\n")

		if (command == "quit") || (command == "x") {
			break
		} else {
			alias, ok := board.UCI_COMMAND_ALIASES[command]

			if ok {
				eng.ExecuteUciCommands(alias)
			} else {
				eng.ExecuteUciCommand(command)
			}
		}
	}
}

/////////////////////////////////////////////////////////////////////
