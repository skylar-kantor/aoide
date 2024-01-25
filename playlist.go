package main

import (
	"fmt"
	"log"
	"os"
	
	"github.com/charmbracelet/bubbles/table"
	"github.com/dhowden/tag"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
)


var rows = []table.Row{}
var player oto.Player
func (m model) ChangePosition(newIdx int) (model, error) {
	log.Default().Printf("Moving to position %d", newIdx)
	if newIdx == 0 {
		return m, nil
	}
	log.Default().Printf("Playlist is: %v. length is %d", m.playlist.Rows(), len(m.playlist.Rows()))
	//Remove everything before the file we're about to play
	m.playlist.SetRows(m.playlist.Rows()[newIdx:])
	log.Default().Printf("Now, Playlist is: %v. length is %d", m.playlist.Rows(), len(m.playlist.Rows()))
	//TODO Extract to separate function?
	//Play the new file
	fileToPlay, err := os.Open(m.playlist.Rows()[0][5])
	if err != nil {
		return m, err
	}
	decodedMP3, err := mp3.NewDecoder(fileToPlay)
	if err != nil {
		return m, err
	}
	c, ready, err := oto.NewContext(decodedMP3.SampleRate(), 2, 2)
	if err != nil {
		return m, err
	}
	<-ready
	if player.IsPlaying() {
		player.Close()
	}
	player = c.NewPlayer(decodedMP3)
	player.Play()
	return m, nil
}

func (m model) AddFile()  model {
	var title, album, artist, time string
	fName := m.fileAdd.Value()
	fTags, err := FileTags(fName)
	if err != nil {
		log.Default().Printf("%v", err)
		return m
	} else {
		title, artist, album, time = SongInfo(fTags)
		if title == "No Title" {
			title = fName
		}
	    tableIdx := len(m.playlist.Rows())
		rowToAdd := table.Row{title, artist, album, time, "false", fName, fmt.Sprint(tableIdx)}
		m.playlist.SetRows(append(m.playlist.Rows(), rowToAdd))
		log.Default().Printf("%v", m.playlist.Rows())
		return m
	}

}

func FileTags(fName string) (tag.Metadata, error) {
	tNewMp3, err := os.Open(fName)
	
	if err != nil {
		//log.Default().Printf("Error opening file %s: %s", fName, err)
		return nil, err
	}
	defer tNewMp3.Close() 
	fTags, err := tag.ReadFrom(tNewMp3)
	if err != nil {
		return nil, err
	} else {
	return fTags, nil
	}
}

func SongInfo(fTags tag.Metadata) (string,string,string, string) {
	title := fTags.Title()
	album := fTags.Album()
	artist := fTags.Artist()
	time := "2:00"
	if title == "" {
		log.Default().Print("No title found, using default")
		title = "No Title"
	}
	if artist == "" {
		log.Default().Print("No artist found, using default")
		artist = "No Artist"
	}
	if album == "" {
		log.Default().Print("No album found, using default")
		album = "No Album"
	}
	//log.Default().Printf("%s, %s, %s", title, artist, album)
	return title, artist, album, time
}