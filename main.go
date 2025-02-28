package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// model holds the state of our game.
type model struct {
	mode       string   // "difficulty" or "game"
	difficulty int      // the chosen difficulty level (1-10)
	threshold  int      // actual allowed maximum Levenshtein distance (inverse mapping)
	chain      []string // chain of moves (most recent move is at the front)
	input      string   // current user input
	score      int      // number of valid moves made
	message    string   // current prompt or feedback message
}

func initialModel() model {
	return model{
		mode:    "difficulty",
		chain:   []string{"rock"}, // game always starts with "rock"
		input:   "",
		score:   0,
		message: "🪨 📜 🚀 Welcome to Rock Paper Go!\n\nSelect a difficulty level (1-10), with 1 being easiest and 10 hardest:",
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
			trimmed := strings.TrimSpace(m.input)
			if m.mode == "difficulty" {
				// Attempt to parse the difficulty level.
				level, err := strconv.Atoi(trimmed)
				if err != nil || level < 1 || level > 10 {
					m.message = "Please enter a valid number between 1 and 10:"
					m.input = ""
					return m, nil
				}
				// Save the chosen difficulty.
				m.difficulty = level
				// Reverse map: input 1 becomes threshold 10, 2 becomes 9, etc.
				m.threshold = 11 - level
				// Switch to game mode.
				m.mode = "game"
				m.message = "What beats rock 🪨?"
				m.input = ""
				return m, nil
			} else {
				// We're in game mode.
				if trimmed == "" {
					return m, nil
				}
				previous := m.chain[0]
				if checkMove(previous, trimmed, m.threshold) {
					// Valid move: prepend new move.
					m.chain = append([]string{trimmed}, m.chain...)
					m.score++
					m.message = fmt.Sprintf("Good! %q beats %q. What beats %q?", trimmed, previous, trimmed)
				} else {
					m.message = fmt.Sprintf("Invalid answer! %q doesn't meet the closeness requirement to %q. Final score: %d. Press q to quit.", trimmed, previous, m.score)
				}
				m.input = ""
				return m, nil
			}

		case "backspace":
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
		default:
			m.input += msg.String()
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.mode == "difficulty" {
		// During difficulty mode, only show the prompt for difficulty.
		return fmt.Sprintf(
			"%s\n\nYour difficulty level: %s\n",
			m.message, m.input,
		)
	}

	// In game mode, display the chosen difficulty above the score.
	var chainDisplay string
	if len(m.chain) > 1 {
		chainDisplay = fmt.Sprintf("Guessed so far: %s\n", strings.Join(m.chain, " 🤜 "))
	}
	return fmt.Sprintf(
		"%s\n\nDifficulty: %d\nScore: %d\n%sYour answer: %s\n",
		m.message, m.difficulty, m.score, chainDisplay, m.input,
	)
}

// checkMove returns true if the Levenshtein distance between the previous move
// and the new move is greater than 0 and less than or equal to the threshold.
func checkMove(previous, move string, threshold int) bool {
	dist := levenshteinDistance(strings.ToLower(previous), strings.ToLower(move))
	return dist > 0 && dist <= threshold
}

// levenshteinDistance calculates the Levenshtein distance between two strings s and t.
// It is the minimum number of single-character edits (insertions, deletions, or substitutions)
// required to change s into t.
func levenshteinDistance(s, t string) int {
	m, n := len(s), len(t)
	dp := make([][]int, m+1)
	for i := 0; i <= m; i++ {
		dp[i] = make([]int, n+1)
		dp[i][0] = i
	}
	for j := 0; j <= n; j++ {
		dp[0][j] = j
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if s[i-1] == t[j-1] {
				dp[i][j] = dp[i-1][j-1]
			} else {
				dp[i][j] = min(
					dp[i-1][j]+1,   // deletion
					dp[i][j-1]+1,   // insertion
					dp[i-1][j-1]+1, // substitution
				)
			}
		}
	}
	return dp[m][n]
}

func min(a, b, c int) int {
	if a < b && a < c {
		return a
	} else if b < c {
		return b
	}
	return c
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
