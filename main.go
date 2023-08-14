package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"time"

	//UI Elements
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	//Audio / mp3
	"github.com/dhowden/tag"
	//	"github.com/hajimehoshi/go-mp3"
	
)

const (
	padding = 2
	maxWidth = 80
)

var (
	//Lipgloss styles for UI elements
	filesStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Italic(true).Render
	winStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Italic(true).Render
	playlistStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Italic(true).Render
	playerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Italic(true).Render
	infoStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Italic(true).Render
	buttonStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Italic(true).Render
	titleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Italic(true).Render
	mainStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Italic(true).Render
	playingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Italic(true).Render

	//Timing vars
	totalLen = 0
	currentLen = 0

	//Player elements
	paused = false
	playText = buttonStyle("||")

	//Playlist
	rows = []row{
		{table.Row{"Dawson's Christian", "Leslie Fish", "Carmen Miranda's Ghost", ""},"",true,120,0},
	}
	
	columns = []table.Column{
		{Title: "Title", Width: 20},
		{Title: "Artist", Width: 20},
		{Title: "Album", Width: 20},
		{Title: "Time", Width: 20},
	}
)

type model struct {
	playBar  progress.Model
	playlist table.Model
	fileAdd  textinput.Model
}

type row struct {
	visibleRow table.Row
	filename string
	playing bool
	time int
	tablePos int
}

type tickMsg time.Time

func (m model) UpdatePlaylist(newIdx int) {
	if newIdx == 0 { return }
	//Remove everything before the file we're about to play
	rows[newIdx].playing = true
	rows = rows[newIdx:]

	m.playlist.SetRows(MakeVisibleRows())

	//TODO Extract to separate function
	//Play the new file
	log.Default().Printf("Playing %s", rows[newIdx].filename)
	//stop playing whatever is currently playing
	//close the file
	//open the new file
	//play it
}

func MakeVisibleRows() []table.Row {
	var visibleRows []table.Row

	for _,r := range rows {
		r.visibleRow[3] = fmt.Sprintf("%02d:%02d", r.time/60, r.time%60)
		
		if r.playing && !strings.Contains(r.visibleRow[0], ">") {
			r.visibleRow[0] = playingStyle(">" + r.visibleRow[0])
		}
		visibleRows = append(visibleRows, r.visibleRow)

	}
	return visibleRows
}

func AddFileToPlaylist(m model) row {

	fName := m.fileAdd.Value()
	
	var title, album, artist string

	tNewMp3, err := os.Open(fName)
	if err != nil {
		log.Default().Fatalf("Error opening file %s: %s", fName, err)
	}

	fTags, err := tag.ReadFrom(tNewMp3)
	if err != nil {
		if err.Error() == "no tags found" {
			log.Default().Printf("No tags found for %s. Using defauts", fName)
			title = fName
			artist = "No Artist"
			album = "No Album"
		} else if err.Error() == "invalid argument" {
			log.Default().Fatalf("File not found: %s", fName)
		} else {
			log.Default().Fatalf("%v", err)
		}
	} else {
		title = fTags.Title()
		album = fTags.Album()
		artist = fTags.Artist()
		if title == "" { title = fName }
		if artist == "" { artist = "No Artist" }
		if album == "" { album = "No Album" }

		tableIdx := len(m.playlist.Rows())
		rowToAdd := row{table.Row{title, artist, album, ""}, fName, false, 0, tableIdx}
		return rowToAdd
	
	}
	return row{table.Row{}, "", false, 0, -1}
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
		})
}

func (m model) Init() tea.Cmd{
	//Get the 1sec timer started
	return tickCmd()
}
	
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "p", " ":
			paused = !paused
			if paused {
				playText = ">"
			} else {
				playText = "||"
			}
		case "l", ".":
			currentLen = int(math.Min(float64(totalLen), float64(currentLen+5)))
		case "j", ",":
			currentLen = int(math.Max(float64(0), float64(currentLen-5)))
		case "tab", "shift+tab":
			if m.playlist.Focused() {
				m.playlist.Blur()
				m.fileAdd.Focus()
			} else {
				m.fileAdd.Blur()
				m.playlist.Focus()
			}
		case "enter":
			if m.fileAdd.Focused() {
				rows = append(rows, AddFileToPlaylist(m))
				m.playlist.SetRows(MakeVisibleRows())
				m.fileAdd.Reset()
				m.fileAdd.Blur()
				m.playlist.Focus()
			} else if m.playlist.Focused() { m.UpdatePlaylist(m.playlist.Cursor())}
	}
	case tea.WindowSizeMsg:
		m.playBar.Width = int(math.Min(float64(msg.Width - padding*2 - 4), float64(maxWidth)))
		return m, nil

	case tickMsg:
		if m.playBar.Percent() == 1.0 && len(rows) > 1 { m.UpdatePlaylist(1) }

		if !paused { currentLen++ }

		cmd := m.playBar.SetPercent(float64(currentLen)/float64(totalLen))
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

func (m model) View() string {
	pad := strings.Repeat(" ", padding)
	//halfPad := strings.Repeat(" ", padding/2)

	timeProgress := infoStyle(fmt.Sprintf("\t%02d:%02d/%02d:%02d\t", currentLen/60, currentLen%60, totalLen/60, totalLen%60))

	nowPlaying := filesStyle(fmt.Sprintf("%s | %s | %s", m.playlist.Rows()[0][0][1:], m.playlist.Rows()[0][1], m.playlist.Rows()[0][2]))

	playlist := playlistStyle(lipgloss.JoinVertical(lipgloss.Left,
		titleStyle("Playlist"),
		m.playlist.View(),
		m.fileAdd.View(),
	))

	player := playerStyle(lipgloss.JoinVertical(lipgloss.Left,
		nowPlaying,
		lipgloss.JoinHorizontal(lipgloss.Right, playText, pad, m.playBar.View(), timeProgress),
	))

	return mainStyle(lipgloss.JoinHorizontal(lipgloss.Center, player, playlist))	
}

func main() {
	fErr, err := tea.LogToFile("debug.log", "DEBUG: ")
	if err != nil {
		fmt.Println("Fatal: Could not start logger, ", err)
		os.Exit(1)
	}
	defer fErr.Close()

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(MakeVisibleRows()),
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