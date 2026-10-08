package tags

import "testing"

func Test_getMetadataFromPath(t *testing.T) {
	tests := []struct {
		name  string
		input Metadata
		want  path_metadata
	}{
		{
			name:  "Standard extraction from path",
			input: Metadata{Path: "/home/user/Music/all_songs/artist_name/album_name/song_name.mp3"},
			want:  path_metadata{song_name: "song_name", album_name: "album_name", artist_name: "artist_name", track_number: ""},
		},
		{
			name:  "Standard with added track number",
			input: Metadata{Path: "/home/user/Music/all_songs/artist_name/album_name/03 song_name.mp3"},
			want:  path_metadata{song_name: "song_name", album_name: "album_name", artist_name: "artist_name", track_number: "03"},
		},
		{
			name:  "Standard with added track number without space in the file name",
			input: Metadata{Path: "/home/user/Music/all_songs/artist_name/album_name/03song_name.mp3"},
			want:  path_metadata{song_name: "song_name", album_name: "album_name", artist_name: "artist_name", track_number: "03"},
		},
		{
			name:  "Only artist directory",
			input: Metadata{Path: "/home/user/Music/all_songs/album_name/song_name.mp3"},
			want:  path_metadata{song_name: "song_name", album_name: "album_name", artist_name: "", track_number: ""},
		},
		{
			name:  "Album name from the file name",
			input: Metadata{Path: "/home/user/Music/all_songs/artist_name/album_name/album - song_name.mp3"},
			want:  path_metadata{song_name: "song_name", album_name: "album", artist_name: "artist_name", track_number: ""},
		},
		{
			name:  "Album and artist name from the file name",
			input: Metadata{Path: "/home/user/Music/all_songs/artist_name/album_name/artist - album - song_name.mp3"},
			want:  path_metadata{song_name: "song_name", album_name: "album", artist_name: "artist", track_number: ""},
		},
		{
			name:  "Album, artist and track number from the file name",
			input: Metadata{Path: "/home/user/Music/all_songs/artist_name/album_name/artist - album - 03 song_name.mp3"},
			want:  path_metadata{song_name: "song_name", album_name: "album", artist_name: "artist", track_number: "03"},
		},
		{
			name:  "Album, artist and track number from the file name, artist in the album dir name",
			input: Metadata{Path: "/home/user/Music/all_songs/artist_name/artist_nm - album_nm/artist - album - 03 song_name.mp3"},
			want:  path_metadata{song_name: "song_name", album_name: "album", artist_name: "artist", track_number: "03"},
		},
		{
			name:  "Album, artist and track number from drectory",
			input: Metadata{Path: "/home/user/Music/all_songs/artist_name/artist_nm - album_nm/song_name.mp3"},
			want:  path_metadata{song_name: "song_name", album_name: "album_nm", artist_name: "artist_nm", track_number: ""},
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metadata := tt.input.getMetadataFromPath()
			if metadata != tt.want {
				t.Errorf("Test %d --- %s\nFUCN: getMetadataFromPath() --- expected = '%+v', got = '%+v'", i+1, tt.name, tt.want, metadata)
			}
		})
	}
}

func Test_isNumeric(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    bool
		wantErr bool
	}{
		{
			name:    "Is valid number based on int",
			input:   01,
			want:    true,
			wantErr: false,
		},
		{
			name:    "Is valid number based on string",
			input:   "01",
			want:    true,
			wantErr: false,
		},
		{
			name:    "Is valid number based on string",
			input:   "Aa",
			want:    false,
			wantErr: false,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isNumeric(tt.input)
			if (result != tt.want) != tt.wantErr {
				t.Errorf("Test %d --- %s\nFUCN: isNumeric() --- expected = '%v', got = '%v'", i+1, tt.name, tt.want, result)
			}
		})
	}
}
