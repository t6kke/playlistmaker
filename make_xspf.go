package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"

	"github.com/t6kke/playlistmaker/internal/tags"
)

type XSPF_Playlist struct {
	XMLName   xml.Name `xml:"playlist"`
	Version   string   `xml:"version,attr"`
	Xmlns     string   `xml:"xmlns,attr"`
	TrackList struct {
		Track []XSPF_Playlist_track `xml:"track"`
	} `xml:"trackList"`
}

type XSPF_Playlist_track struct {
	Location string `xml:"location"`
	Title    string `xml:"title"`
	Creator  string `xml:"creator"`
	Album    string `xml:"album"`
	TrackNum string `xml:"trackNum"`
	Image    string `xml:"image"`
}

const playlist_file_name = "all_songs.xspf"

func (mlc *MusicLibraryConfig) createAllSongsPL(all_songs []tags.Metadata) error {
	var tracks []XSPF_Playlist_track

	for _, item := range all_songs {
		data := XSPF_Playlist_track{
			Location: item.Path,
			Title:    item.Title,
			Creator:  item.Creator,
			Album:    item.Album,
			TrackNum: item.Track_nbr,
			Image:    "(embedded)",
		}
		tracks = append(tracks, data)
	}

	playlist := XSPF_Playlist{
		Version: "1",
		Xmlns:   "http://xspf.org/ns/0/",
	}
	playlist.TrackList.Track = tracks

	playlist_file_full_path := mlc.Songs_playlist_dir + "/" + playlist_file_name
	file, err := os.Create(filepath.Clean(playlist_file_full_path))
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := xml.NewEncoder(file)
	encoder.Indent("", "  ")

	fmt.Print(xml.Header)
	_, err = file.WriteString(xml.Header)
	if err != nil {
		return err
	}

	err = encoder.Encode(playlist)
	if err != nil {
		return err
	}

	return nil
}
