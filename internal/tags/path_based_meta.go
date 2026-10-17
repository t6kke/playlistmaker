package tags

import (
	"strconv"
	"strings"
)

type path_metadata struct {
	song_name    string
	album_name   string
	artist_name  string
	track_number string
}

// TODO a lot of various edge cases and special situation still to handle
func (m *Metadata) getMetadataFromPath() path_metadata {
	var result_data path_metadata

	parts := strings.Split(m.Path, "/")
	parts_len := len(parts)

	for i := parts_len - 1; i >= 0; i-- {
		switch i {
		case parts_len - 1:
			var album_name_from_file string
			var artist_name_from_file string

			file_name_parts := strings.Split(parts[i], "-")

			switch len(file_name_parts) {
			case 1:
				if isNumeric(file_name_parts[0][:2]) {
					result_data.track_number = file_name_parts[0][:2]
					result_data.song_name = strings.TrimSpace(strings.Split(file_name_parts[0][2:], ".")[0])
				} else {
					result_data.song_name = strings.Split(file_name_parts[0], ".")[0]
				}
			case 2:
				album_name_from_file = strings.TrimSpace(file_name_parts[0])
				//last_index := strings.LastIndex(file_name_parts[1], ".")
				//potential_song_name := strings.TrimSpace(file_name_parts[1][last_index:])
				potential_song_name := strings.TrimSpace(strings.Split(file_name_parts[1], ".")[0])
				if isNumeric(potential_song_name[:2]) {
					result_data.track_number = potential_song_name[:2]
					result_data.song_name = potential_song_name[3:]
				} else {
					result_data.song_name = potential_song_name
				}
			case 3:
				artist_name_from_file = strings.TrimSpace(file_name_parts[0])
				album_name_from_file = strings.TrimSpace(file_name_parts[1])
				potential_song_name := strings.TrimSpace(strings.Split(file_name_parts[2], ".")[0])
				if isNumeric(potential_song_name[:2]) {
					result_data.track_number = potential_song_name[:2]
					result_data.song_name = potential_song_name[3:]
				} else {
					result_data.song_name = potential_song_name
				}
			}
			result_data.album_name = album_name_from_file
			result_data.artist_name = artist_name_from_file

		case parts_len - 2:
			if result_data.album_name == "" {
				parts_from_dir := strings.Split(parts[i], "-")

				switch len(parts_from_dir) {
				case 1:
					result_data.album_name = parts[i]
				case 2:
					result_data.artist_name = strings.TrimSpace(parts_from_dir[0])
					result_data.album_name = strings.TrimSpace(parts_from_dir[1])
				}
			}

		case parts_len - 3:
			if parts[i] == "all_songs" {
				break
			}
			if result_data.artist_name == "" {
				result_data.artist_name = parts[i]
			}
		}
	}
	return result_data
}

func isNumeric(input any) bool {
	if _, ok := input.(int); ok {
		return true
	}
	if data, ok := input.(string); ok {
		_, err := strconv.Atoi(data)
		if err == nil {
			return true
		}
	}
	return false
}
