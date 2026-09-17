package tags

import "testing"

func Test_ExtractMetadata_Errors(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		result  metadata
		wantErr bool
	}{
		{
			name:    "Invalid Path -- Output error",
			path:    "dummy",
			result:  metadata{},
			wantErr: true,
		},
		{
			name:    "Invalid Path -- Output error",
			path:    "/home/user/Music/music.mp3",
			result:  metadata{},
			wantErr: true,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metadata, err := ExtractMetadata(tt.path, 1)
			if (err != nil) != tt.wantErr {
				t.Errorf("Test %d --- ExtractMetadata() error = '%v', wantErr '%v'", i+1, err, tt.wantErr)
			}
			if metadata != tt.result {
				t.Errorf("Test %d --- ExtractMetadata() metadat = %v, wantErr %v", i+1, metadata, tt.result)
			}
		})
	}
}

func Test_TagSizeFromBits(t *testing.T) {
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
				t.Errorf("Test %d --- getCorrectTagSize() expected = '%v', got '%v'", i+1, tt.want, size)
			}
		})
	}
}

func Test_FileNameFromPath(t *testing.T) {
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
				t.Errorf("Test %d --- getFileName() expected = '%v', got '%v'", i+1, tt.want, file_name)
			}
		})
	}
}

func Test_IsNumber(t *testing.T) {
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
				t.Errorf("Test %d --- isNumeric() expected = '%v', got '%v'", i+1, tt.want, result)
			}
		})
	}
}
