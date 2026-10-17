package tags

import (
	"bytes"
	"testing"
)

func Test_extract_mp3v1(t *testing.T) {
	//TODO
}

func Test_extract_mp3id3_v2_2(t *testing.T) {
	//TODO
}

func Test_extract_mp3id3_v2_3(t *testing.T) {
	//TODO
}

func Test_getFrameFlags_v3(t *testing.T) {
	//TODO
}

func Test_extract_mp3id3_v2_4(t *testing.T) {
	//TODO
}

func Test_getFrameFlags_v4(t *testing.T) {
	//TODO
}

func Test_getHeaderDetails(t *testing.T) {
	tests := []struct {
		name  string
		bytes []byte
		want  mp3_headerDetails
	}{
		{
			name:  "Correct name extraction",
			bytes: []byte{73, 68, 51, 4, 0, 128, 0, 0, 17, 45},
			want:  mp3_headerDetails{identifier: "ID3", version: 4, revision: 0, flags: mp3_headerFlags{unsynchronisation: true, extended_header: false, experimental_indicator: false, footer_present: false}, size: 2221},
		},
		{
			name:  "Correct name extraction",
			bytes: []byte{73, 68, 51, 3, 0, 0, 0, 3, 22, 100},
			want:  mp3_headerDetails{identifier: "ID3", version: 3, revision: 0, flags: mp3_headerFlags{unsynchronisation: false, extended_header: false, experimental_indicator: false, footer_present: false}, size: 52068},
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header_info := mp3_getHeaderDetails(tt.bytes)
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
			size := mp3_getCorrectTagSize(tt.bytes)
			if size != tt.want {
				t.Errorf("Test %d --- %s\nFUCN: getCorrectTagSize() --- expected = '%v', got = '%v'", i+1, tt.name, tt.want, size)
			}
		})
	}
}

func Test_unsynchronize(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  []byte
	}{
		{
			name:  "Empty slice returns empty slice",
			input: []byte{},
			want:  []byte{},
		},
		{
			name:  "No changes, slice is returned as is",
			input: []byte{0x1A, 0xB7, 0x3C, 0xE2, 0x59, 0xD4, 0x08, 0xAF, 0x71, 0xCE},
			want:  []byte{0x1A, 0xB7, 0x3C, 0xE2, 0x59, 0xD4, 0x08, 0xAF, 0x71, 0xCE},
		},
		{
			name:  "No changes, slice is returned as is",
			input: []byte{0x1A, 0xB7, 0x3C, 0xE2, 0x59, 0xD4, 0x08, 0xAF, 0x71, 0xCE},
			want:  []byte{0x1A, 0xB7, 0x3C, 0xE2, 0x59, 0xD4, 0x08, 0xAF, 0x71, 0xCE},
		},
		{
			name:  "Unsynchronization is applied",
			input: []byte{0x1A, 0xB7, 0x3C, 0xFF, 0x00, 0xD4, 0x08, 0xAF, 0x71, 0xCE},
			want:  []byte{0x1A, 0xB7, 0x3C, 0xFF, 0xD4, 0x08, 0xAF, 0x71, 0xCE},
		},
		{
			name:  "Last item is the one that is detected for unsynchronization but nothing to sync and no changes done",
			input: []byte{0x1A, 0xB7, 0x3C, 0xE2, 0x59, 0xD4, 0x08, 0xAF, 0x71, 0xFF},
			want:  []byte{0x1A, 0xB7, 0x3C, 0xE2, 0x59, 0xD4, 0x08, 0xAF, 0x71, 0xFF},
		},
		{
			name:  "Multiple items for unsynchronize",
			input: []byte{0xFF, 0xB7, 0xFF, 0x00, 0xFF, 0xFF, 0x00, 0xAF, 0x71, 0xFF},
			want:  []byte{0xFF, 0xB7, 0xFF, 0xFF, 0xFF, 0xAF, 0x71, 0xFF},
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mp3_unsynchronize(&tt.input)
			if !bytes.Equal(tt.input, tt.want) {
				t.Errorf("Test %d --- %s\nFUCN: unsynchronize() --- expected = '%02X', got = '%02X'", i+1, tt.name, tt.want, tt.input)
			}
		})
	}
}
