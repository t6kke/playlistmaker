package tags

import (
	"fmt"
	"os"
	//"slices"
	//"strings"
)

type blockType int

const (
	streaminfo blockType = iota
	padding
	application
	seek_table
	vorbis_comment
	cuesheet
	picture
)

func (b blockType) String() string {
	return [...]string{"streaminfo", "padding", "application", "seek_table", "vorbis_comment", "cuesheet", "picture"}[b]
}

type flacBlockHeader struct {
	flags flacBlockHeaderFlags
	size  int
}

type flacBlockHeaderFlags struct {
	is_last   bool
	blck_type blockType
}

func (m *Metadata) extract_flac_meta(f *os.File) error {

	final_block := false
	pos := int64(4)
	for !final_block {
		block_header := make([]byte, 4)
		_, _ = f.ReadAt(block_header, pos)
		fmt.Printf("    block header: %v --- %08b\n", block_header, block_header)

		header := getBlockHeader(block_header)

		//block_header_flags := getBlockHeaderFlags(block_header[0])
		if header.flags.is_last {
			final_block = true
		}

		pos += int64(4) + int64(header.size)
	}

	return nil
}

func getBlockHeader(b []byte) flacBlockHeader {
	return flacBlockHeader{
		flags: getBlockHeaderFlags(b[0]),
		size:  int(b[1])<<16 | int(b[2])<<8 | int(b[3]),
	}
}

func getBlockHeaderFlags(b byte) flacBlockHeaderFlags {
	return flacBlockHeaderFlags{
		is_last:   b&0b10000000 != 0,
		blck_type: getBlockHeaderTypeInfo(b),
	}
}

func getBlockHeaderTypeInfo(b byte) blockType {
	val := int(b & 0x7f)
	fmt.Println(val)

	data := blockType(val)
	fmt.Println(data)

	return data
}
