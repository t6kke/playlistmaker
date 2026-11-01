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
	return []string{"mp3", "MP3", "flac"}
}

func getNotRelevantFileExtensions() []string {
	return []string{"jpg", "jpeg", "png", "txt", "nfo", "ini", "pdf", "url", "db", "rtf", "m3u"}
}

type playlist_obj struct {
	playlist_name  string
	ident_category identCategory
	ident_keyword  string
	ident_wild     bool
	metadata_items []tags.Metadata
}

func (mlc *MusicLibraryConfig) scanner() {
	//default playlists
	playlists := map[string]playlist_obj{
		"all_tracks":    {playlist_name: "all_tracks", metadata_items: make([]tags.Metadata, 0)},
		"regular_songs": {playlist_name: "regular_songs", metadata_items: make([]tags.Metadata, 0)},
	}
	//user custom playlists
	for _, user_pl_config := range mlc.Custom_playlists {
		if user_pl_config.Playlist_name == "example_name" {
			continue
		}
		playlists[user_pl_config.Identifier_keyword] = playlist_obj{
			playlist_name:  user_pl_config.Playlist_name,
			ident_category: user_pl_config.Identifier_category,
			ident_keyword:  user_pl_config.Identifier_keyword,
			ident_wild:     user_pl_config.Is_lookup_wildcard,
			metadata_items: make([]tags.Metadata, 0),
		}
	}

	dir_for_scan := getDirForScanAndDirPlaylists(mlc.Songs_dir, &playlists)
	//fmt.Println(len(dir_for_scan), "items:", dir_for_scan)

	var missing_handling_for_files []string

	scanForMetadata(dir_for_scan, &playlists, &missing_handling_for_files)
	fmt.Printf("Missing parsing of file types: %v\n", missing_handling_for_files)

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

func getDirForScanAndDirPlaylists(dir string, playlists *map[string]playlist_obj) []string {
	var result = make([]string, 0)

	files, _ := os.ReadDir(dir)

	for _, file := range files {
		if file.IsDir() && strings.Contains(file.Name(), "_soundtracks") {
			//fmt.Printf("Scan This: '%s'\n", file.Name())

			//make new playlist based on the dir
			full_dir := dir + "/" + file.Name()
			parts := strings.Split(full_dir, "/_")
			pl_name := strings.Split(parts[1], "_soundtracks")[0] + "_soundtracks"
			pl_ident := "_" + pl_name

			_, exists := (*playlists)[pl_ident]
			if !exists {
				//fmt.Println(" !!! NEW PLAYLIST !!!")
				//fmt.Print(strings.ToUpper(pl_ident)+": ", full_dir)
				new_pl_obj := playlist_obj{
					playlist_name:  pl_name,
					metadata_items: make([]tags.Metadata, 0),
				}
				//fmt.Println(len(metadata))
				(*playlists)[pl_ident] = new_pl_obj
				//fmt.Println("count",strings.ToUpper(pl_ident)+": ", len((*playlists)[pl_ident].metadata_items))
			}

			data := getDirForScanAndDirPlaylists(dir+"/"+file.Name(), playlists)
			result = append(result, data...)
		} else {
			//fmt.Printf("Add This: '%s'\n", file.Name())
			result = append(result, dir+"/"+file.Name())
		}
	}
	return result
}

func scanForMetadata(root_dir_for_scan []string, playlists *map[string]playlist_obj, missing_handling_for_files *[]string) {
	for _, dir := range root_dir_for_scan {
		metadata := recursive_scanner(dir, "", missing_handling_for_files)
		//fmt.Print("ALL: ", dir)
		//fmt.Println(" --- nbr", len(metadata))
		if len(metadata) == 0 {
			continue
		}

		item := (*playlists)["all_tracks"]
		item.metadata_items = append(item.metadata_items, metadata...)
		(*playlists)["all_tracks"] = item
		//fmt.Println("count ALL: ", len((*playlists)["all_tracks"].metadata_items))

		for pl_id, playlist := range *playlists {
			if isUserCustomPL(pl_id) {
				switch playlist.ident_category {
				case 1:
					if strings.EqualFold(metadata[0].Creator, playlist.ident_keyword) {
						item := (*playlists)[pl_id]
						item.metadata_items = append(item.metadata_items, metadata...)
						(*playlists)[pl_id] = item
					}
				case 2:
					if playlist.ident_wild {
						if strings.Contains(strings.ToLower(metadata[0].Album), strings.ToLower(playlist.ident_keyword)) {
							item := (*playlists)[pl_id]
							item.metadata_items = append(item.metadata_items, metadata...)
							(*playlists)[pl_id] = item
						}
					} else {
						if strings.EqualFold(metadata[0].Album, playlist.ident_keyword) {
							item := (*playlists)[pl_id]
							item.metadata_items = append(item.metadata_items, metadata...)
							(*playlists)[pl_id] = item
						}
					}
				}
			}
		}

		if strings.Contains(dir, "/_") && strings.Contains(dir, "_soundtracks") {
			parts := strings.Split(dir, "/_")
			pl_name := strings.Split(parts[1], "_soundtracks/")[0] + "_soundtracks"
			pl_ident := "_" + pl_name

			_, exists := (*playlists)[pl_ident]
			if !exists {
				//TODO handle situation where playlist based on directory does not exists. but this situation should never happen
			}

			//fmt.Print(strings.ToUpper(pl_ident)+": ", dir)
			//fmt.Println(" --- nbr", len(metadata))
			item := (*playlists)[pl_ident]
			item.metadata_items = append(item.metadata_items, metadata...)
			(*playlists)[pl_ident] = item
			//fmt.Println("count",strings.ToUpper(pl_ident)+": ", len((*playlists)[pl_ident].metadata_items))

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

func isUserCustomPL(pl_identifier string) bool {
	if strings.HasPrefix(pl_identifier, "_") && strings.Contains(pl_identifier, "_soundtracks") {
		return false
	}
	switch pl_identifier {
	case "all_tracks":
		return false
	case "regular_songs":
		return false
	default:
		return true
	}
}

func recursive_scanner(dir, spacer string, missing_handling_for_files *[]string) []tags.Metadata {
	files, _ := os.ReadDir(dir)
	file_count := 0
	music_file_count := 0
	var result_data []tags.Metadata
	for _, file := range files {
		file_count += 1
		file_extension := strings.ToLower(strings.Split(file.Name(), ".")[len(strings.Split(file.Name(), "."))-1])
		if !file.IsDir() && !slices.Contains(getAudioFileExtensions(), file_extension) && !slices.Contains(getNotRelevantFileExtensions(), file_extension) {
			//fmt.Printf("Missing file extention parsing for: '%s'\n",file_extension)
			if !slices.Contains(*missing_handling_for_files, file_extension) {
				*missing_handling_for_files = append(*missing_handling_for_files, file_extension)
			}
		}
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
			data := recursive_scanner(dir+"/"+file.Name(), spacer+"   ", missing_handling_for_files)
			result_data = append(result_data, data...)
		}
	}
	return result_data
}
