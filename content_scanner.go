package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/t6kke/playlistmaker/internal/tags"
)

func getAudioFileExtensions() []string {
	return []string{} //TODO remove, this is for testing purpose, revert back to original when actually scanning
	//return []string{"mp3", "flac"}
}

func (mlc *MusicLibraryConfig) scanner() []tags.Metadata {
	var temp_result = make([]tags.Metadata, 0)

	dir_channel := make(chan string)
	go func() {
		defer close(dir_channel)
		getDirForScan(mlc.Songs_dir, dir_channel)
	}()

	for {
		item, ok := <- dir_channel
		if ok == false {
			break
		}
		fmt.Println(item)
	}

	// for item := range dir_channel {
	// 	fmt.Println(item)
	// }

	return temp_result

	//original setup before intial dir scanning to channel
	//all_files_metadata := recursive_scanner(mlc.Songs_dir, "")
	//fmt.Printf("Number of items handled: '%d'\n", len(all_files_metadata))
	//return all_files_metadata
}

func getDirForScan(dir string, dir_channel chan string) {
	//var result = make([]string,0)

	//fmt.Println(dir)
	files, _ := os.ReadDir(dir)

	for _, file := range files {
		if file.IsDir() && strings.Contains(file.Name(), "_soundtracks") {
			//fmt.Printf("Scan This: '%s'\n", file.Name())
			getDirForScan(dir+"/"+file.Name(), dir_channel)
		} else {
			//fmt.Printf("Add This: '%s'\n", file.Name())
			dir_channel <- dir+"/"+file.Name()
		}
	}
}

func recursive_scanner(dir, spacer string) []tags.Metadata {
	files, _ := os.ReadDir(dir)
	file_count := 0
	music_file_count := 0
	var result_data []tags.Metadata
	for _, file := range files {
		file_count += 1
		file_extension := strings.Split(file.Name(), ".")[len(strings.Split(file.Name(), "."))-1]
		if !file.IsDir() && slices.Contains(getAudioFileExtensions(), file_extension) {

			music_file_count += 1
			metadata, err := tags.ExtractMetadata(dir+"/"+file.Name(), music_file_count, file_extension)
			if err != nil {
				return []tags.Metadata{}
			}
			//fmt.Printf("%+v\n", metadata)
			result_data = append(result_data, metadata)
			//fmt.Println(music_file_count, result_data)

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
