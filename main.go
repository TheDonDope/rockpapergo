package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	score   int
	chain   []string // Chain of moves, with the most recent move at the front.
	input   string   // Player's current input.
	message string   // Status or prompt message.
}

func initialModel() model {
	return model{
		score:   0,
		chain:   []string{"rock"}, // Game always starts with "rock".
		input:   "",
		message: "What beats rock?",
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			playerMove := strings.TrimSpace(m.input)
			if playerMove == "" {
				return m, nil
			}
			// Validate player's move using custom logic.
			if checkMove(m.chain[0], playerMove) {
				// Valid move: prepend the new move to the chain.
				m.chain = append([]string{playerMove}, m.chain...)
				m.score++
				// m.chain[1] holds the previous move.
				m.message = fmt.Sprintf("Good! %q beats %q. What beats %q?", playerMove, m.chain[1], playerMove)
			} else {
				m.message = fmt.Sprintf("Invalid move! %q does not beat %q. Final score: %d. Press q to quit.", playerMove, m.chain[0], m.score)
			}
			m.input = ""
			return m, nil
		case "backspace":
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
		default:
			// Append any other key press to the input.
			m.input += msg.String()
		}
	}
	return m, nil
}

func (m model) View() string {
	var chainDisplay string
	if len(m.chain) > 1 {
		chainDisplay = fmt.Sprintf("Guessed so far: %s\n", strings.Join(m.chain, " 🤜 "))
	}
	return fmt.Sprintf(
		"%s\n\nScore: %d\n%s\nYour answer: %s\n",
		m.message, m.score, chainDisplay, m.input,
	)
}

// checkMove is a placeholder for your custom logic.
// Replace this with your rules to determine if one move beats another.
func checkMove(current, move string) bool {
	return move != ""
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
