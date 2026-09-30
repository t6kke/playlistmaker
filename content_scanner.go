package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/t6kke/playlistmaker/internal/tags"
)

func getAudioFileExtensions() []string {
	return []string{"mp3"} //TODO need to map out more filetypes when adding metadata extraction to them.
	//return []string{"mp3", "flac"}
}

func (mlc *MusicLibraryConfig) scanner() []tags.Metadata {
	all_files_metadata := recursive_scanner(mlc.Songs_dir, "")
	fmt.Println(len(all_files_metadata))
	return all_files_metadata
}

func recursive_scanner(dir, spacer string) []tags.Metadata {
	audio_file_extensions := getAudioFileExtensions()
	files, _ := os.ReadDir(dir)
	file_count := 0
	music_file_count := 0
	var result_data []tags.Metadata
	for _, file := range files {
		file_count += 1
		if !file.IsDir() && slices.Contains(audio_file_extensions, strings.Split(file.Name(), ".")[len(strings.Split(file.Name(), "."))-1]) {
			if file_count == 1 {
				music_file_count += 1
				metadata, err := tags.ExtractMetadata(dir+"/"+file.Name(), music_file_count)
				if err != nil {
					return []tags.Metadata{}
				}
				//fmt.Printf("%+v\n", metadata)
				result_data = append(result_data, metadata)
			}
		} else if file.IsDir() {
			//fmt.Println(spacer, file.Name())
			if string(file.Name()[0]) == "_" && strings.Contains(file.Name(), "_soundtracks") {
				//TODO handle separately
			} else {
				data := recursive_scanner(dir+"/"+file.Name(), spacer+"   ")
				result_data = append(result_data, data...)
			}

		}
	}
	return result_data
}
