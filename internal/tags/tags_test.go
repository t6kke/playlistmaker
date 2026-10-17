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
