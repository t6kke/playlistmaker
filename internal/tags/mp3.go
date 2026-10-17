package tags

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
)

type mp3_headerDetails struct {
	identifier string
	version    int
	revision   int
	flags      mp3_headerFlags
	size       int
	crc        bool
}

type mp3_headerFlags struct {
	unsynchronisation      bool
	extended_header        bool
	experimental_indicator bool
	footer_present         bool
}

type mp3_2_2_headerDetails struct {
	identifier string
	version    int
	revision   int
	flags      mp3_2_2_headerFlags
	size       int
}

type mp3_2_2_headerFlags struct {
	unsynchronisation bool
	compression       bool
}

/*
 * MP3 TAGv1.1
 *
 */

type v1_tags struct {
	title     string
	creator   string
	album     string
	track_nbr int //TODO to string
}

func (m *Metadata) extract_mp3v1(tag_bytes []byte) (v1_tags, error) {

	if string(tag_bytes[0:3]) != "TAG" {
		return v1_tags{}, fmt.Errorf("not valid v1 tag dataset, does not start with 'TAG'")
	}

	v1_tags := v1_tags{
		title:     string(tag_bytes[3:33]),
		creator:   string(tag_bytes[33:63]),
		album:     string(tag_bytes[63:93]),
		track_nbr: int(tag_bytes[126]), //TODO handle existence sice TAG v1.1 uses thsi byte as track number and v1 does not
	}

	return v1_tags, nil
}

/*
 * MP3 ID3 v2.3 section starts
 *
 */

func (m *Metadata) extract_mp3id3_v2_2(initial_header_details *mp3_headerDetails, header_flags_byte byte, f *os.File) error {
	header_details := mp3_2_2_headerDetails{
		identifier: initial_header_details.identifier,
		version:    initial_header_details.version,
		revision:   initial_header_details.revision,
		flags:      mp3_2_2_getHeaderFlags(header_flags_byte),
		size:       initial_header_details.size,
	}

	fmt.Printf("%+v\n", header_details)

	fmt.Println("   ", "Frames Data")
	tag_data := make([]byte, header_details.size)
	_, err := io.ReadFull(f, tag_data)
	if err != nil {
		return err
	}

	frame_identifier_name_map := getFrameIdentifierMap_v2_2()
	frame_identifier_list_to_parse := getFrameIdentifierSliceToParse_v2_2()

	if header_details.flags.unsynchronisation {
		mp3_unsynchronize(&tag_data)
	}

	frame_header_size := 10
	for pos := 0; pos+frame_header_size <= len(tag_data); {
		fmt.Println("   ", tag_data[pos:pos+6])
		//fmt.Println("   ", hex.Dump(tag_data[pos:pos+10]))

		id := string(tag_data[pos : pos+3])

		//stop after final tag item
		if id == "\x00\x00\x00\x00" {
			fmt.Println("   ", " --- no more data")
			break
		}

		size := int(tag_data[pos+3])<<16 | int(tag_data[pos+4])<<8 | int(tag_data[pos+5])
		if size <= 0 || pos+6+size > len(tag_data) {
			break
		}

		fmt.Printf("    frame header id: %s --- and size of: %d, bytes: %d\n", id, size, tag_data[pos+3:pos+6])

		if slices.Contains(frame_identifier_list_to_parse, id) {
			frame := tag_data[pos+6 : pos+6+size]
			fmt.Println("        ", frame)
			//fmt.Println("        ", hex.Dump(frame))
			fmt.Printf("         %s: %v\n", frame_identifier_name_map[id], decodeText(frame))

			switch id {
			case "TT2":
				m.Title = decodeText(frame)
			case "TP1":
				m.Creator = decodeText(frame)
			case "TAL":
				m.Album = decodeText(frame)
			case "TLE":
				m.Duration = decodeText(frame)
			case "TRK":
				m.Track_nbr = decodeText(frame)
				if strings.Contains(m.Track_nbr, "/") {
					m.Track_nbr = strings.Split(m.Track_nbr, "/")[0]
				}
			}
		}
		pos += 6 + size
	}

	return nil
}

func mp3_2_2_getHeaderFlags(b byte) mp3_2_2_headerFlags {
	return mp3_2_2_headerFlags{
		unsynchronisation: b&0b10000000 != 0,
		compression:       b&0b01000000 != 0,
	}
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

func (m *Metadata) extract_mp3id3_v2_3(header_details *mp3_headerDetails, f *os.File) error {
	fmt.Println("   ", "Frames Data")
	tag_data := make([]byte, header_details.size)
	_, err := io.ReadFull(f, tag_data)
	if err != nil {
		return err
	}

	frame_identifier_name_map := getFrameIdentifierMap_v2_3()
	frame_identifier_list_to_parse := getFrameIdentifierSliceToParse_v2_3()

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
		mp3_unsynchronize(&tag_data)
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

func (m *Metadata) extract_mp3id3_v2_4(header_details *mp3_headerDetails, f *os.File) error {
	fmt.Println("   ", "Frames Data")
	tag_data := make([]byte, header_details.size)
	_, err := io.ReadFull(f, tag_data)
	if err != nil {
		return err
	}

	frame_identifier_name_map := getFrameIdentifierMap_v2_3()
	frame_identifier_list_to_parse := getFrameIdentifierSliceToParse_v2_3()

	extended_header_size := 0
	if header_details.flags.extended_header {
		size := mp3_getCorrectTagSize(tag_data[10:15])
		//TODO parse extended header details

		extended_header_size = 4 + size

	}

	//fmt.Println(len(tag_data))
	//fmt.Println(hex.Dump(tag_data))
	is_unsync_done := false
	if header_details.flags.unsynchronisation {
		mp3_unsynchronize(&tag_data)
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
		size := mp3_getCorrectTagSize(tag_data[pos+4 : pos+8])
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
				d_len := mp3_getCorrectTagSize(frame[:4])
				fmt.Println(d_len)
				frame = frame[4:]
			}
			if !is_unsync_done && frame_flags.is_unsynchronized {
				fmt.Println("test")
				mp3_unsynchronize(&frame)
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

func mp3_getHeaderDetails(b []byte) mp3_headerDetails {
	result_data := mp3_headerDetails{}

	result_data.identifier = string(b[0:3])
	result_data.version = int(b[3])
	result_data.revision = int(b[4])
	result_data.flags = mp3_getHeaderFlags(b[5])
	result_data.size = mp3_getCorrectTagSize(b[6:10])

	//fmt.Printf("header struct: %+v\n", result_data)

	return result_data
}

func mp3_getCorrectTagSize(b []byte) int {
	//TODO validate lenght of bytes slice and return error if it's not 4

	//fmt.Printf("%08b  %08b %08b\n",b[2], b[2]&0x7f, int(b[2]&0x7f)<<7)
	full_data := int(b[0]&0x7f)<<21 | int(b[1]&0x7f)<<14 | int(b[2]&0x7f)<<7 | int(b[3]&0x7f)
	//fmt.Printf("%d, %08b", full_data, full_data)

	return full_data
}

func mp3_getHeaderFlags(b byte) mp3_headerFlags {
	return mp3_headerFlags{
		unsynchronisation:      b&0b10000000 != 0,
		extended_header:        b&0b01000000 != 0,
		experimental_indicator: b&0b00100000 != 0,
		footer_present:         b&0b00010000 != 0,
	}
}

func mp3_unsynchronize(b *[]byte) {
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

/*
 * MP3 frame identifier information
 * both ID3 v2.2 and v2.3; v2.4
 *
 */

func getFrameIdentifierMap_v2_2() map[string]string {
	//following details from: https://id3.org/id3v2-00
	//initally only set up information that are relevant for me, maybe will expand but since v2.2 is legacy setup it's likely not needed'
	return map[string]string{
		"TAL": "album name",
		"TLE": "lenght",
		"TP1": "leade artist",
		"TRK": "track number",
		"TT2": "title of the track",
	}
}

func getFrameIdentifierSliceToParse_v2_2() []string {
	return []string{"TT2", "TP1", "TAL", "TLE", "TRK"}
}

func getFrameIdentifierMap_v2_3() map[string]string {
	//following details from: https://id3.org/id3v2.3.0#ID3v2_header
	//TODO maybe not all the details here are actually needed
	return map[string]string{
		"TALB": "album name",
		"TBPM": "beats per minute",
		"TCOM": "composers name",
		"TCON": "content type", //TODO more details here for genres
		"TCOP": "copyright infromation",
		"TDAT": "date information in DDMM format",
		"TDLY": "playlist details, silence between songs in a playlist",
		"TENC": "person or organization that encoded the file",
		"TEXT": "writers of the content",
		"TFLT": "file type, indicates type of audio this tag defines", //TODO more details
		"TIME": "time of the recording in HHMM format",
		"TIT1": "content group description",
		"TIT2": "name of the song/piece",
		"TIT3": "subtitle or descrition of the song/piece",
		"TKEY": "musical key in which the sound starts",
		"TLAN": "language(s) of the lyrics of the song",
		"TLEN": "length of the audiofile in milliseconds represented as numeric string",
		"TMED": "description of the media the sound originated from", //TODO more details
		"TOAL": "original album of the song",
		"TOFN": "orignial/preferred filename",
		"TOLY": "original writers of the lyrics",
		"TOPE": "original performer",
		"TORY": "original release year",
		"TOWN": "license owner name",
		"TPE1": "main artist name",
		"TPE2": "additional information on the performers",
		"TPE3": "conductor name",
		"TPE4": "information on the authors of the remix",
		"TPOS": "information what part of the set the audio came from(CD1, CD2 etc.), in numeric string format",
		"TPUB": "lable or publisher name",
		"TRCK": "track number",
		"TRDA": "reckording date, used as acompliment to TYER, TDAT and TIME",
		"TRSN": "internet radio station name",
		"TRSO": "internet radio station owner",
		"TSIZ": "size of the audio file in bytes excluding the tag",
		"TSRC": "International Standard Recording Code(ISRC)",
		"TSSE": "audio encoder used and it's setting when it was encoded",
		"TYER": "year of the recording, in numeric string format and always four characters long",
	}
}

func getFrameIdentifierSliceToParse_v2_3() []string {
	return []string{"TIT2", "TPE1", "TALB", "TLEN", "TRCK"}
}
