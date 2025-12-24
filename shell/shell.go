package shell

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Handler — функция-обработчик команды.
type Handler func(args []string)

// Command описывает одну команду в shell.
type Command struct {
	Name        string  // имя команды
	Description string  // краткое описание
	Usage       string  // пример использования
	Handler     Handler // вызываемая функция
}

// Shell хранит зарегистрированные команды и пейджер.
type Shell struct {
	commands map[string]*Command
	prompt   string
	pager    *Pager
}

// New создаёт Shell с базовыми командами.
func New() *Shell {
	sh := &Shell{
		commands: make(map[string]*Command),
		prompt:   "> ",
		pager:    NewPager(),
	}

	// регистрируем команды
	sh.Register(&Command{
		Name:        "help",
		Description: "список команд или подробности по конкретной",
		Usage:       "help [command]",
		Handler:     sh.helpHandler,
	})
	sh.Register(&Command{
		Name:        "stats",
		Description: "показать статистику",
		Usage:       "stats",
		Handler:     sh.statsHandler,
	})
	sh.Register(&Command{
		Name:        "start",
		Description: "начать загрузку",
		Usage:       "start",
		Handler:     sh.startHandler,
	})
	sh.Register(&Command{
		Name:        "set-uspeed",
		Description: "задать скорость отдачи в Mbit/s",
		Usage:       "set-uspeed <mbit>",
		Handler:     sh.setUspeedHandler,
	})
	sh.Register(&Command{
		Name:        "clear",
		Description: "очистить экран",
		Usage:       "clear",
		Handler:     sh.clearHandler,
	})
	sh.Register(&Command{
		Name:        "exit",
		Description: "выход из shell",
		Usage:       "exit",
		Handler:     sh.exitHandler,
	})

	return sh
}

// Register добавляет команду в shell.
func (sh *Shell) Register(cmd *Command) {
	sh.commands[cmd.Name] = cmd
}

// Run запускает цикл ввода команд.
func (sh *Shell) Run() {
	fmt.Println("Welcome to torrent-cli shell! Type 'help' to list commands.")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(sh.prompt)
		if !scanner.Scan() {
			fmt.Println("\nExiting shell.")
			return
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		name, args := parts[0], parts[1:]
		cmd, ok := sh.commands[name]
		if !ok {
			fmt.Printf("Unknown command: %s. Type 'help' for list.\n", name)
			continue
		}
		cmd.Handler(args)
	}
}

// helpHandler — список команд или подробный help через пейджер.
func (sh *Shell) helpHandler(args []string) {
	if len(args) == 0 {
		fmt.Println("Available commands:")
		for _, cmd := range sh.commands {
			fmt.Printf("  %-12s %s\n", cmd.Name, cmd.Description)
		}
		return
	}
	if len(args) == 1 {
		name := args[0]
		cmd, ok := sh.commands[name]
		if !ok {
			fmt.Printf("No such command: %s\n", name)
			return
		}
		text := fmt.Sprintf("%s\n\nUsage:\n  %s\n", cmd.Description, cmd.Usage)
		if err := sh.pager.Page(text); err != nil {
			// fallback на случай, если пейджер не работает
			fmt.Print(text)
		}
		return
	}
	fmt.Println("Usage:")
	fmt.Println("  help             — список команд")
	fmt.Println("  help <command>   — подробная справка по команде")
}

// Заглушки остальных команд

func (sh *Shell) statsHandler(args []string) {
	fmt.Println("stats: not implemented yet")
}

func (sh *Shell) startHandler(args []string) {
	fmt.Println("start: not implemented yet")
}

func (sh *Shell) setUspeedHandler(args []string) {
	if len(args) != 1 {
		fmt.Println("Usage:", sh.commands["set-uspeed"].Usage)
		return
	}
	fmt.Printf("set-uspeed %s: not implemented yet\n", args[0])
}

func (sh *Shell) clearHandler(args []string) {
	// ANSI escape: move cursor home + clear screen
	fmt.Print("\033[H\033[2J")
}

func (sh *Shell) exitHandler(args []string) {
	fmt.Println("exit: not implemented yet")
}
