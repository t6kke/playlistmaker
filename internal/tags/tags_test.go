package tags

import "testing"

func Test_ExtractMetadata_Errors(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		result  Metadata
		wantErr bool
	}{
		{
			name:    "Invalid Path -- Output error",
			path:    "dummy",
			result:  Metadata{},
			wantErr: true,
		},
		{
			name:    "Invalid Path -- Output error",
			path:    "/home/user/Music/music.mp3",
			result:  Metadata{},
			wantErr: true,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metadata, err := ExtractMetadata(tt.path, 1)
			if (err != nil) != tt.wantErr {
				t.Errorf("Test %d --- %s\nFUCN: ExtractMetadata() --- error = '%v', wantErr = '%v'", i+1, tt.name, err, tt.wantErr)
			}
			if metadata != tt.result {
				t.Errorf("Test %d --- %s\nFUCN: ExtractMetadata() --- metadat = %v, wantErr = %v", i+1, tt.name, metadata, tt.result)
			}
		})
	}
}

func Test_getCorrectTagSize(t *testing.T) {
	tests := []struct {
		name    string
		bytes   []byte
		want    int
		wantErr bool
	}{
		{
			name:    "Accurate bytes calculation",
			bytes:   []byte{0, 3, 22, 100},
			want:    52068,
			wantErr: false,
		},
		{
			name:    "Falze bytes calculation",
			bytes:   []byte{0, 3, 22, 100},
			want:    1,
			wantErr: true,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			size := getCorrectTagSize(tt.bytes)
			if (size != tt.want) != tt.wantErr {
				t.Errorf("Test %d --- %s\nFUCN: getCorrectTagSize() --- expected = '%v', got = '%v'", i+1, tt.name, tt.want, size)
			}
		})
	}
}

func Test_getFileName(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		want    string
		wantErr bool
	}{
		{
			name:    "Correct name extraction",
			path:    "/home/user/Music/music.mp3",
			want:    "music.mp3",
			wantErr: false,
		},
		{
			name:    "Correct name extraction",
			path:    "/home/user/Music/artist - album - song.mp3",
			want:    "artist - album - song.mp3",
			wantErr: false,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file_name := getFileName(tt.path)
			if (file_name != tt.want) != tt.wantErr {
				t.Errorf("Test %d --- %s\nFUCN: getFileName() --- expected = '%v', got = '%v'", i+1, tt.name, tt.want, file_name)
			}
		})
	}
}

func Test_getMetadataFromPath(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  path_metadata
	}{
		{
			name:  "Standard extraction from path",
			input: "/home/user/Music/all_songs/artist_name/album_name/song_name.mp3",
			want:  path_metadata{song_name: "song_name", album_name: "album_name", artist_name: "artist_name", track_number: ""},
		},
		{
			name:  "Standard with added track number",
			input: "/home/user/Music/all_songs/artist_name/album_name/03 song_name.mp3",
			want:  path_metadata{song_name: "song_name", album_name: "album_name", artist_name: "artist_name", track_number: "03"},
		},
		{
			name:  "Standard with added track number without space in the file name",
			input: "/home/user/Music/all_songs/artist_name/album_name/03song_name.mp3",
			want:  path_metadata{song_name: "song_name", album_name: "album_name", artist_name: "artist_name", track_number: "03"},
		},
		{
			name:  "Only artist directory",
			input: "/home/user/Music/all_songs/album_name/song_name.mp3",
			want:  path_metadata{song_name: "song_name", album_name: "album_name", artist_name: "", track_number: ""},
		},
		{
			name:  "Album name from the file name",
			input: "/home/user/Music/all_songs/artist_name/album_name/album - song_name.mp3",
			want:  path_metadata{song_name: "song_name", album_name: "album", artist_name: "artist_name", track_number: ""},
		},
		{
			name:  "Album and artist name from the file name",
			input: "/home/user/Music/all_songs/artist_name/album_name/artist - album - song_name.mp3",
			want:  path_metadata{song_name: "song_name", album_name: "album", artist_name: "artist", track_number: ""},
		},
		{
			name:  "Album, artist and track number from the file name",
			input: "/home/user/Music/all_songs/artist_name/album_name/artist - album - 03 song_name.mp3",
			want:  path_metadata{song_name: "song_name", album_name: "album", artist_name: "artist", track_number: "03"},
		},
		{
			name:  "Album, artist and track number from the file name, artist in the album dir name",
			input: "/home/user/Music/all_songs/artist_name/artist_nm - album_nm/artist - album - 03 song_name.mp3",
			want:  path_metadata{song_name: "song_name", album_name: "album", artist_name: "artist", track_number: "03"},
		},
		{
			name:  "Album, artist and track number from drectory",
			input: "/home/user/Music/all_songs/artist_name/artist_nm - album_nm/song_name.mp3",
			want:  path_metadata{song_name: "song_name", album_name: "album_nm", artist_name: "artist_nm", track_number: ""},
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metadata := getMetadataFromPath(tt.input)
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
