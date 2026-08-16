package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

func (mlc *MusicLibraryConfig) scanner() {
	recursive_scanner(mlc.Songs_dir, "")
}

func recursive_scanner(dir, spacer string) {
	audio_file_extensions := []string{"mp3", "flac"} //TODO need to map out more items
	files, _ := os.ReadDir(dir)
	for _, file := range files {
		if !file.IsDir() && slices.Contains(audio_file_extensions, strings.Split(file.Name(), ".")[len(strings.Split(file.Name(), "."))-1]) {
			fmt.Println(spacer, dir+"/"+file.Name()) //TODO this is an actual audio file need to do metadata extraction and then use that to build playlist
		} else if file.IsDir() {
			fmt.Println(spacer, file.Name())
			if string(file.Name()[0]) == "_" && strings.Contains(file.Name(), "_soundtracks") {
				//TODO handle separately
			} else {
				recursive_scanner(dir+"/"+file.Name(), spacer+"   ")
			}

		}
	}
}
