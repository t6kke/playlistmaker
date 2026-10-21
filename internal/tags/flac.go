package tags

import (
	"encoding/binary"
	//"encoding/hex"
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
	comment_handled := false
	pos := int64(4)
	for !final_block && !comment_handled {
		block_header := make([]byte, 4)
		_, _ = f.ReadAt(block_header, pos) //TODO handle error
		fmt.Printf("    block header: %v --- %08b\n", block_header, block_header)

		//fmt.Println(hex.Dump(block_header))

		header := getBlockHeader(block_header)

		if header.flags.blck_type == vorbis_comment {
			fmt.Printf("        Block type 4. Here is metadata info. Content size: '%d'\n", header.size)
			v_comment_payload := make([]byte, header.size)
			_, _ = f.ReadAt(v_comment_payload, pos+int64(4)) //TODO handle error
			//fmt.Println(hex.Dump(v_comment_payload))
			vendor_size := binary.LittleEndian.Uint32(v_comment_payload[0:4])
			fmt.Println("       ", vendor_size)
			vendor_string := v_comment_payload[4:vendor_size]
			fmt.Println("       ", vendor_string, string(vendor_string))
			comment_count := binary.LittleEndian.Uint32(v_comment_payload[4+vendor_size:4+vendor_size+4])
			fmt.Println("       ", v_comment_payload[4+vendor_size:4+vendor_size+4],comment_count)
			comment_pos := 4+vendor_size+4
			for i := 1; i <= int(comment_count); i++ {
				c_size := binary.LittleEndian.Uint32(v_comment_payload[comment_pos:comment_pos+4])
				fmt.Printf("        comment_size '%d', '%d'\n", v_comment_payload[comment_pos:comment_pos+4], c_size)

				comment_val := v_comment_payload[comment_pos+4:comment_pos+4+c_size]
				fmt.Printf("        comment_val '%d', '%s'\n", v_comment_payload[comment_pos+4:comment_pos+4+c_size], comment_val)

				comment_pos += 4 + c_size
			}

			comment_handled = true //TODO technically on one vorbis comment should be used and we don't need to handle any others if they exists but need to validate that we actually got metadata we needed from the first one
		}

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
	data := blockType(val)
	return data
}
