package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/t6kke/playlistmaker/internal/tags"
)

func getAudioFileExtensions() []string {
	//return []string{} //TODO remove, this is for testing purpose, revert back to original when actually scanning
	//return []string{"flac"} //TODO remove, this is for testing purpose, revert back to original when actually scanning
	return []string{"mp3", "flac"}
}

type playlist_obj struct {
	playlist_name  string
	metadata_items []tags.Metadata
}

func (mlc *MusicLibraryConfig) scanner() {
	dir_for_scan := getDirForScan(mlc.Songs_dir)
	//fmt.Println(len(dir_for_scan), "items:", dir_for_scan)

	playlists := map[string]playlist_obj{
		"all_tracks":    {playlist_name: "all_tracks", metadata_items: make([]tags.Metadata, 0)},
		"regular_songs": {playlist_name: "regular_songs", metadata_items: make([]tags.Metadata, 0)},
	}
	//TODO additional custom playlits from MusicLibraryConfig once it's added to the configuration

	scanForMetadata(dir_for_scan, &playlists)

	fmt.Println(len(playlists), "playlist object items:")
	for _, pl_item := range playlists {
		fmt.Println("Playlist Name:", pl_item.playlist_name)
		//fmt.Printf("%+v\n", pl_item)
		fmt.Println(pl_item.playlist_name, "len: ", len(pl_item.metadata_items))
	}

	err := mlc.writeXSPFplaylists(&playlists)
	if err != nil {
		panic(err)
	}
}

func getDirForScan(dir string) []string {
	var result = make([]string, 0)

	files, _ := os.ReadDir(dir)

	for _, file := range files {
		if file.IsDir() && strings.Contains(file.Name(), "_soundtracks") {
			//fmt.Printf("Scan This: '%s'\n", file.Name())
			data := getDirForScan(dir + "/" + file.Name())
			result = append(result, data...)
		} else {
			//fmt.Printf("Add This: '%s'\n", file.Name())
			result = append(result, dir+"/"+file.Name())
		}
	}
	return result
}

func scanForMetadata(root_dir_for_scan []string, playlists *map[string]playlist_obj) {
	//TODO metadata scanning should be only done once per file not once each time it shows up it a given playlist
	for _, dir := range root_dir_for_scan {
		metadata := recursive_scanner(dir, "")
		//fmt.Print("ALL: ", dir)
		//fmt.Println(" --- nbr", len(metadata))
		item := (*playlists)["all_tracks"]
		item.metadata_items = append(item.metadata_items, metadata...)
		(*playlists)["all_tracks"] = item
		//fmt.Println("count ALL: ", len((*playlists)["all_tracks"].metadata_items))
		item = playlist_obj{}

		if strings.Contains(dir, "/_") && strings.Contains(dir, "_soundtracks") {
			parts := strings.Split(dir, "/_")
			pl_name := strings.Split(parts[1], "_soundtracks/")[0] + "_soundtracks"
			pl_ident := "_" + pl_name

			_, exists := (*playlists)[pl_ident]

			if !exists {
				//fmt.Println(" !!! NEW PLAYLIST !!!")
				//fmt.Print(strings.ToUpper(pl_ident)+": ", dir)
				new_pl_obj := playlist_obj{
					playlist_name:  pl_name,
					metadata_items: make([]tags.Metadata, 0),
				}
				//fmt.Println(len(metadata))
				new_pl_obj.metadata_items = append(item.metadata_items, metadata...)
				(*playlists)[pl_ident] = new_pl_obj
				//fmt.Println("count",strings.ToUpper(pl_ident)+": ", len((*playlists)[pl_ident].metadata_items))
			} else {
				//fmt.Print(strings.ToUpper(pl_ident)+": ", dir)
				//fmt.Println(" --- nbr", len(metadata))
				item := (*playlists)[pl_ident]
				item.metadata_items = append(item.metadata_items, metadata...)
				(*playlists)[pl_ident] = item
				//fmt.Println("count",strings.ToUpper(pl_ident)+": ", len((*playlists)[pl_ident].metadata_items))
			}
		} else {
			//fmt.Print("SONGS: ", dir)
			//fmt.Println(" --- nbr", len(metadata))
			item := (*playlists)["regular_songs"]
			item.metadata_items = append(item.metadata_items, metadata...)
			(*playlists)["regular_songs"] = item
			//fmt.Println("count SONGS: ", len((*playlists)["regular_songs"].metadata_items))
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
			data := recursive_scanner(dir+"/"+file.Name(), spacer+"   ")
			result_data = append(result_data, data...)
		}
	}
	return result_data
}
