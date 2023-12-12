package main

import (
	"fmt"
	"log"
	"os"
	
	"github.com/charmbracelet/bubbles/table"
	"github.com/dhowden/tag"
)


var rows = []table.Row{}

func (m model) ChangePosition(newIdx int) model{
	log.Default().Printf("Moving to position %d", newIdx)
	if newIdx == 0 {
		return m
	}
	log.Default().Printf("Playlist is: %v. length is %d", m.playlist.Rows(), len(m.playlist.Rows()))
	//Remove everything before the file we're about to play
	rows = m.playlist.Rows()[newIdx:]
	log.Default().Printf("Playlist is: %v. length is %d", m.playlist.Rows(), len(m.playlist.Rows()))
	m.playlist.SetRows(rows)
	log.Default().Printf("Playlist is: %v. length is %d", m.playlist.Rows(), len(m.playlist.Rows()))
	//TODO Extract to separate function?
	//Play the new file
	log.Default().Printf("Playing %s", rows[0][4])
	//stop playing whatever is currently playing
	//close the file
	//Can't defer close, bec that runs on return, so
	//if isPlaying {
	// Close()
	//}
	//open the new file
	//play it
	return m
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
	time := "0:00"
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