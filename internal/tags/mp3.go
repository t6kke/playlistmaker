package tags

import (
	"encoding/binary"
	"os"
	//"encoding/hex"
	"fmt"
	"io"
	"slices"
	"strings"
)

/*
 * MP3 v1 tags
 *
 */

type v1_tags struct {
	title     string
	creator   string
	album     string
	track_nbr int
}

func (m *Metadata) extract_mp3v1(tag_bytes []byte) (v1_tags, error) {

	if string(tag_bytes[0:3]) != "TAG" {
		return v1_tags{}, fmt.Errorf("not valid v1 tag dataset, does not start with 'TAG'")
	}

	v1_tags := v1_tags{
		title:     string(tag_bytes[3:33]),
		creator:   string(tag_bytes[33:63]),
		album:     string(tag_bytes[63:93]),
		track_nbr: int(tag_bytes[126]),
	}

	return v1_tags, nil
}

/*
 * MP3 ID3 v2.3 section starts
 *
 */

type frameFlags_v3 struct {
	//status flags
	tag_alter_preservation  bool //if true discard frame
	file_alter_preservation bool //if true discard frameFlags
	read_only               bool //frame content is read only, parser should not change anything, not relevant for this application

	//format flags
	in_group      bool //shows if frame belongs to a group
	is_compressed bool
	is_encrypted  bool
}

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
	if header_details.flags.unsynchronisation {
		unsynchronize(&tag_data)
	}

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
		frame_flags := getFrameFlags_v3(tag_data[pos+8 : pos+10])
		fmt.Printf("    frame flags: %+v\n", frame_flags)
		//TODO deal with frame header flags
		//handled: -

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
				if strings.Contains(m.Track_nbr, "/") {
					m.Track_nbr = strings.Split(m.Track_nbr, "/")[0]
				}
			}
		}
		pos += 10 + size
	}

	return nil
}

func getFrameFlags_v3(b []byte) frameFlags_v3 {
	status := b[0]
	format := b[1]

	return frameFlags_v3{
		tag_alter_preservation:  status&0b10000000 != 0,
		file_alter_preservation: status&0b01000000 != 0,
		read_only:               status&0b00100000 != 0,

		is_compressed: format&0b10000000 != 0,
		is_encrypted:  format&0b01000000 != 0,
		in_group:      format&0b00100000 != 0,
	}
}

/*
 * MP3 ID3 v2.4 section starts
 *
 */

type frameFlags_v4 struct {
	//status flags
	tag_alter_preservation  bool //if true discard frame
	file_alter_preservation bool //if true discard frameFlags
	read_only               bool //frame content is read only, parser should not change anything, not relevant for this application

	//format flags
	in_group                  bool //shows if frame belongs to a group
	is_compressed             bool
	is_encrypted              bool
	is_unsynchronized         bool
	has_data_length_indicator bool
}

func (m *Metadata) extract_mp3id3v24(header_details *headerDetails, f *os.File) error {
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
		size := getCorrectTagSize(tag_data[10:15])
		//TODO parse extended header details

		extended_header_size = 4 + size

	}

	//fmt.Println(len(tag_data))
	//fmt.Println(hex.Dump(tag_data))
	is_unsync_done := false
	if header_details.flags.unsynchronisation {
		unsynchronize(&tag_data)
		is_unsync_done = true
	}

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

		//size := int(binary.BigEndian.Uint32(tag_data[pos+4 : pos+8]))
		size := getCorrectTagSize(tag_data[pos+4 : pos+8])
		if size <= 0 || pos+10+size > len(tag_data) {
			break
		}

		fmt.Printf("    frame header id: %s --- and size of: %d, bytes: %d --- flags data: %08b\n", id, size, tag_data[pos+4:pos+8], tag_data[pos+8:pos+10])
		frame_flags := getFrameFlags_v4(tag_data[pos+8 : pos+10])
		//TODO deal with frame header flags
		//handled: data_lenght, unsync

		if slices.Contains(frame_identifier_list_to_parse, id) {
			frame := tag_data[pos+10 : pos+10+size]
			if frame_flags.has_data_length_indicator {
				d_len := getCorrectTagSize(frame[:4])
				fmt.Println(d_len)
				frame = frame[4:]
			}
			if !is_unsync_done && frame_flags.is_unsynchronized {
				fmt.Println("test")
				unsynchronize(&frame)
			}
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
				if strings.Contains(m.Track_nbr, "/") {
					m.Track_nbr = strings.Split(m.Track_nbr, "/")[0]
				}
			}
		}
		pos += 10 + size
	}

	return nil
}

func getFrameFlags_v4(b []byte) frameFlags_v4 {
	status := b[0]
	format := b[1]

	return frameFlags_v4{
		tag_alter_preservation:  status&0b01000000 != 0,
		file_alter_preservation: status&0b00100000 != 0,
		read_only:               status&0b00010000 != 0,

		in_group:                  format&0b01000000 != 0,
		is_compressed:             format&0b00001000 != 0,
		is_encrypted:              format&0b00000100 != 0,
		is_unsynchronized:         format&0b00000010 != 0,
		has_data_length_indicator: format&0b00000001 != 0,
	}
}

/*
 * MP3 ID3 general helper functions for all version
 *
 */

func unsynchronize(b *[]byte) {
	data := *b
	result := make([]byte, 0)
	for i := 0; i < len(data); i++ {
		//fmt.Printf("%02X ", b[i])
		result = append(result, data[i])
		if data[i] == 0xFF && i+1 < len(data) && data[i+1] == 0x00 {
			i++
		}
	}
	*b = result
}
