package tags

import (
	"encoding/binary"
	"os"
	//"encoding/hex"
	"fmt"
	"io"
	"slices"
)

func (m *Metadata) extract_mp3id3v23(header_details *headerDetails, f *os.File) error {
	fmt.Println("   ", "Frames Data")
	tag_data := make([]byte, header_details.size)
	_, err := io.ReadFull(f, tag_data)
	if err != nil {
		return err
	}

	frame_identifier_name_map := getFrameIdentifierMap()
	frame_identifier_list_to_parse := getFrameIdentifierSliceToParse()

	extended_header_size := 0
	if header_details.flags.extended_header {
		size := int(binary.BigEndian.Uint32(tag_data[10:15]))
		extended_flags := binary.BigEndian.Uint16(tag_data[15:17])
		header_details.crc = extended_flags&0x8000 != 0 //CRC is used for checksum of flags data

		extended_header_size = 4 + size
	}

	//fmt.Println(len(tag_data))
	//fmt.Println(hex.Dump(tag_data))
	unsynchronize(tag_data)

	frame_header_size := 10 + extended_header_size
	for pos := 0; pos+frame_header_size <= len(tag_data); {
		fmt.Println("   ", tag_data[pos:pos+10])
		//fmt.Println("   ", hex.Dump(tag_data[pos:pos+10]))

		id := string(tag_data[pos : pos+4])

		//stop after final tag item
		if id == "\x00\x00\x00\x00" {
			fmt.Println("   ", " --- no more data")
			break
		}

		size := int(binary.BigEndian.Uint32(tag_data[pos+4 : pos+8]))
		if size <= 0 || pos+10+size > len(tag_data) {
			break
		}

		fmt.Printf("    frame header id: %s --- and size of: %d, bytes: %d --- flags data: %08b\n", id, size, tag_data[pos+4:pos+8], tag_data[pos+8:pos+10])
		//TODO deal with frame header flags

		if slices.Contains(frame_identifier_list_to_parse, id) {
			frame := tag_data[pos+10 : pos+10+size]
			fmt.Println("        ", frame)
			//fmt.Println("        ", hex.Dump(frame))
			fmt.Printf("         %s: %v\n", frame_identifier_name_map[id], decodeText(frame))

			switch id {
			case "TIT2":
				m.Title = decodeText(frame)
			case "TPE1":
				m.Creator = decodeText(frame)
			case "TALB":
				m.Album = decodeText(frame)
			case "TLEN":
				m.Duration = decodeText(frame)
			case "TRCK":
				m.Track_nbr = decodeText(frame)
			}
		}
		pos += 10 + size
	}

	return nil
}

//TODO find a better location for this function
func unsynchronize(b []byte) []byte {
	result := make([]byte, 0)
	for i := 0; i < len(b); i++ {
		//fmt.Printf("%02X ", b[i])
		result = append(result, b[i])
		if b[i] == 0xFF && i+1 < len(b) && b[i+1] == 0x00 {
			i++
		}
	}
	return result
}
