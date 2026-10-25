package tags

import (
	"testing"
)

func Test_extract_flac_meta(t *testing.T) {
	//TODO
}

func Test_flac_getBlockHeader(t *testing.T) {
	tests := []struct {
		name  string
		bytes []byte
		want  flacBlockHeader
	}{
		{
			name:  "first block: Streaminfo",
			bytes: []byte{0, 0, 0, 34},
			want:  flacBlockHeader{flacBlockHeaderFlags{false, 0}, 34},
		},
		{
			name:  "block 4, vorbis comment with accurate size",
			bytes: []byte{4, 0, 0, 143},
			want:  flacBlockHeader{flacBlockHeaderFlags{false, 4}, 143},
		},
		{
			name:  "block 4, vorbis comment with accurate size as last block",
			bytes: []byte{132, 0, 0, 143},
			want:  flacBlockHeader{flacBlockHeaderFlags{true, 4}, 143},
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			block_header := flac_getBlockHeader(tt.bytes)
			if block_header != tt.want {
				t.Errorf("Test %d --- %s\nFUCN: flac_getBlockHeader() --- expected = '%+v', got = '%+v'", i+1, tt.name, tt.want, block_header)
			}
		})
	}
}
