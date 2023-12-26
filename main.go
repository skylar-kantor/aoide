package main

import (
	"fmt"
	//"log"
	//"math"
	"os"
	//"strings"
	"time"

	//UI Elements
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	//Audio / mp3
	//"github.com/dhowden/tag"
	//"github.com/hajimehoshi/go-mp3"
)

var columns = []table.Column{
	{Title: "Title", Width: 20},
	{Title: "Artist", Width: 20},
	{Title: "Album", Width: 20},
	{Title: "Time", Width: 20},
	{Title: "", Width: 0},
	{Title: "", Width: 0},
	{Title: "", Width: 0},
}

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func main() {
	logFile, err := tea.LogToFile("debug.log", "DEBUG: ")
	if err != nil {
		fmt.Println("Fatal: Could not start logger, ", err)
		os.Exit(1)
	}
	defer logFile.Close()

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(3),
	)

	f := textinput.New()
	f.Placeholder = "Add files to playlist..."
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("#861270")).
		Bold(true)
	t.SetStyles(s)

	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithoutPercentage(),
	)

	m := model{
		playBar: p, playlist: t, fileAdd: f,
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Fatal: ", err)
		os.Exit(1)
	}
}
