package tags

import "testing"

func Test_ExtractMetadata_Errors(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		file_type string
		result    Metadata
		wantErr   bool
	}{
		{
			name:      "Invalid Path -- Output error",
			path:      "dummy",
			file_type: "mp3",
			result:    Metadata{},
			wantErr:   true,
		},
		{
			name:      "Invalid Path -- Output error",
			path:      "/home/user/Music/music.mp3",
			file_type: "mp3",
			result:    Metadata{},
			wantErr:   true,
		},
		{
			name:      "Invalid file type -- Output error",
			path:      "/home/user/Music/music.mp3",
			file_type: "png",
			result:    Metadata{},
			wantErr:   true,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metadata, err := ExtractMetadata(tt.path, 1, tt.file_type)
			if (err != nil) != tt.wantErr {
				t.Errorf("Test %d --- %s\nFUCN: ExtractMetadata() --- error = '%v', wantErr = '%v'", i+1, tt.name, err, tt.wantErr)
			}
			if metadata != tt.result {
				t.Errorf("Test %d --- %s\nFUCN: ExtractMetadata() --- metadat = %v, wantErr = %v", i+1, tt.name, metadata, tt.result)
			}
		})
	}
}

func Test_getHeaderDetails(t *testing.T) {
	tests := []struct {
		name  string
		bytes []byte
		want  headerDetails
	}{
		{
			name:  "Correct name extraction",
			bytes: []byte{73, 68, 51, 4, 0, 128, 0, 0, 17, 45},
			want:  headerDetails{identifier: "ID3", version: 4, revision: 0, flags: headerFlags{unsynchronisation: true, extended_header: false, experimental_indicator: false, footer_present: false}, size: 2221},
		},
		{
			name:  "Correct name extraction",
			bytes: []byte{73, 68, 51, 3, 0, 0, 0, 3, 22, 100},
			want:  headerDetails{identifier: "ID3", version: 3, revision: 0, flags: headerFlags{unsynchronisation: false, extended_header: false, experimental_indicator: false, footer_present: false}, size: 52068},
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header_info := getHeaderDetails(tt.bytes)
			if header_info != tt.want {
				t.Errorf("Test %d --- %s\nFUCN: getHeaderDetails() --- expected = '%+v', got = '%+v'", i+1, tt.name, tt.want, header_info)
			}
		})
	}
}

func Test_getCorrectTagSize(t *testing.T) {
	tests := []struct {
		name  string
		bytes []byte
		want  int
	}{
		{
			name:  "Accurate bytes calculation",
			bytes: []byte{0, 3, 22, 100},
			want:  52068,
		},
		{
			name:  "Accurate bytes calculation",
			bytes: []byte{0, 0, 17, 45},
			want:  2221,
		},
		{
			name:  "Accurate bytes calculation",
			bytes: []byte{0, 7, 88, 00},
			want:  125952,
		},
		{
			name:  "Accurate bytes calculation",
			bytes: []byte{0, 28, 10, 58},
			want:  460090,
		},
		{
			name:  "Accurate bytes calculation",
			bytes: []byte{0, 0, 1, 14},
			want:  142,
		},
		{
			name:  "Accurate bytes calculation",
			bytes: []byte{4, 0, 11, 20},
			want:  8390036,
		},
		{
			name:  "Accurate bytes calculation, size zero",
			bytes: []byte{0, 0, 0, 0},
			want:  0,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			size := getCorrectTagSize(tt.bytes)
			if size != tt.want {
				t.Errorf("Test %d --- %s\nFUCN: getCorrectTagSize() --- expected = '%v', got = '%v'", i+1, tt.name, tt.want, size)
			}
		})
	}
}

func Test_getFileName(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "Correct name extraction",
			path: "/home/user/Music/music.mp3",
			want: "music.mp3",
		},
		{
			name: "Correct name extraction",
			path: "/home/user/Music/artist - album - song.mp3",
			want: "artist - album - song.mp3",
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file_name := getFileName(tt.path)
			if file_name != tt.want {
				t.Errorf("Test %d --- %s\nFUCN: getFileName() --- expected = '%v', got = '%v'", i+1, tt.name, tt.want, file_name)
			}
		})
	}
}
