package main

import (
	"fmt"
	"log"
	"math"
	"reflect"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	playBar  progress.Model
	playlist table.Model
	fileAdd  textinput.Model
}

const (
	padding  = 2
	maxWidth = 80
	
)

var (
	StyleFilePicker = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Italic(true).Render
	StyleWindow     = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Italic(true).Render
	StylePlaylist   = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Italic(true).Render
	StylePlayer     = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Italic(true).Render
	StyleTime       = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Italic(true).Render
	StylePlayText   = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Italic(true).Render
	//StyleTitle      = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Italic(true).Render
	StyleApp        = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Italic(false).Render

	paused   = true
	playText = StylePlayText("⏵")

	totalLen   = 0
	currentLen = 0
	pad = strings.Repeat(" ", padding)
)

func (m model) Init() tea.Cmd {
	//Get the 1sec timer started
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		var c tea.Cmd
		m, c = m.HandleKeyPress(msg.String())
		if reflect.ValueOf(c).Pointer() == reflect.ValueOf(tea.Quit).Pointer() {
			return m, c
		}
	case tea.WindowSizeMsg:
		m.playBar.Width = int(math.Min(float64(msg.Width-padding*2-4), float64(maxWidth)))
		return m, nil

	case tickMsg:
		if m.playBar.Percent() == 1.0 && len(m.playlist.Rows()) > 1 {
			log.Default().Print("Song ended, moving on...")
			m = m.ChangePosition(1)
		}

		if !paused {
			currentLen++
		}

		cmd := m.playBar.SetPercent(float64(currentLen) / float64(totalLen))
		return m, tea.Batch(tickCmd(), cmd)
	case progress.FrameMsg:
		progressModel, cmd := m.playBar.Update(msg)
		m.playBar = progressModel.(progress.Model)
		return m, cmd
	default:
		return m, nil

	}
	var tblCmd, tiCmd tea.Cmd
	m.playlist, tblCmd = m.playlist.Update(msg)
	m.fileAdd, tiCmd = m.fileAdd.Update(msg)
	return m, tea.Batch(tiCmd, tblCmd)
}

func (m model) HandleKeyPress(key string) (model, tea.Cmd){
	switch key {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "p", " ":
		if len(m.playlist.Rows()) > 0 && m.playlist.Rows()[0][3] != "0:00" {
		paused = !paused
		} else {
			return m, nil
		}
		if paused {
			playText = "⏵"
		} else {
			playText = "⏸"
		}
	case "l", ".":
		currentLen = int(math.Min(float64(totalLen), float64(currentLen+5)))
	case "j", ",":
		currentLen = int(math.Max(float64(0), float64(currentLen-5)))
	case "tab", "shift+tab":
		m = m.switchFocus()
	case "enter":
		if m.fileAdd.Focused() {
			m = m.AddFile()
			log.Default().Printf("%v", m.playlist.Rows())
			m.fileAdd.Reset()
			m = m.switchFocus()
		} else if m.playlist.Focused() {
			m = m.ChangePosition(m.playlist.Cursor())
		}
	}
	return m, nil
}

func (m model) switchFocus() model {
	if m.playlist.Focused() {
		m.playlist.Blur()
		m.fileAdd.Focus()
	} else {
		m.fileAdd.Blur()
		m.playlist.Focus()
	}
	return m
}

func (m model) View() string {
	
	timeProgress := StyleTime(fmt.Sprintf(" %02d:%02d/%02d:%02d ", currentLen/60, currentLen%60, totalLen/60, totalLen%60))
	var nowPlaying string
	if len(m.playlist.Rows()) > 0 {
		nowPlaying = StyleFilePicker(fmt.Sprintf(" Now Playing:\n %s | %s | %s", m.playlist.Rows()[0][0], m.playlist.Rows()[0][1], m.playlist.Rows()[0][2]))
	} else {
		nowPlaying = StyleFilePicker(" Now Playing:\n \t | \t | \t")
	}

	playlist := StylePlaylist(lipgloss.JoinVertical(lipgloss.Left,
		
		m.playlist.View(),
		m.fileAdd.View(),
	))

	player := StylePlayer(lipgloss.JoinVertical(lipgloss.Left,
		nowPlaying,
		lipgloss.JoinHorizontal(lipgloss.Right, playText, pad, m.playBar.View(), timeProgress),
	))

	return StyleApp(lipgloss.JoinHorizontal(lipgloss.Top, player, playlist))
}