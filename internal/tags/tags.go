package tags

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"unicode/utf16"
)

func getFrameIdentifierMap() map[string]string {
	//following details from: https://id3.org/id3v2.3.0#ID3v2_header
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

func getFrameIdentifierSliceToParse() []string {
	return []string{"TIT2", "TPE1", "TALB", "TLEN", "TRCK"}
}

type metadata struct {
	path      string
	title     string
	creator   string
	album     string
	duration  string
	track_nbr string
}

func ExtractMetadata(path string) (metadata, error) {
	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		return metadata{}, err
	}
	defer f.Close()

	header := make([]byte, 10)
	if _, err := io.ReadFull(f, header); err != nil {
		return metadata{}, err
	}

	fmt.Println(path)
	fmt.Println("Header Data")
	fmt.Println(header[0:10])

	tag_size := getCorrectTagSize(header[6:10])

	fmt.Printf("file header type: %s --- and version: %d.%d --- flags: %08b and tag size of: %d, bytes: %d\n", string(header[0:3]), header[3], header[4], header[5], tag_size, header[6:10])

	frame_identifier_name_map := getFrameIdentifierMap()
	frame_identifier_list_to_parse := getFrameIdentifierSliceToParse()

	data := metadata{}
	data.path = path

	if int(header[3]) == 3 {
		fmt.Println("   ", "Frames Data")
		tag_data := make([]byte, tag_size)
		_, err := io.ReadFull(f, tag_data) //TODO handle error
		if err != nil {
			return metadata{}, err
		}

		frameHeaderSize := 10
		for pos := 0; pos+frameHeaderSize <= len(tag_data); {
			fmt.Println("   ", tag_data[pos:pos+10])

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
			if slices.Contains(frame_identifier_list_to_parse, id) {
				frame := tag_data[pos+10 : pos+10+size]
				fmt.Println("        ", frame)
				fmt.Printf("         %s: %v\n", frame_identifier_name_map[id], decodeText(frame))

				//TODO data validation, if something is missing alternative method has to be used to get the info
				switch id {
				case "TIT2":
					//TODO validation, likely from the file name
					data.title = decodeText(frame)
				case "TPE1":
					//TODO validation, likely from the file name
					data.creator = decodeText(frame)
				case "TALB":
					//TODO validation, how? tbd
					data.album = decodeText(frame)
				case "TLEN":
					//TODO validation, calculated other information from the file
					data.duration = decodeText(frame)
				case "TRCK":
					//TODO validation, maybe from the file name, or just based on the order of the files in the directory
					data.track_nbr = decodeText(frame)
				}
			}

			pos += 10 + size
		}
	}
	fmt.Println(data)
	fmt.Println("")

	return data, nil
}

func getCorrectTagSize(b []byte) int {
	//fmt.Printf("%08b  %08b %08b\n",b[2], b[2]&0x7f, int(b[2]&0x7f)<<7)
	full_data := int(b[0]&0x7f)<<21 | int(b[1]&0x7f)<<14 | int(b[2]&0x7f)<<7 | int(b[3]&0x7f)
	//fmt.Printf("%d, %08b", full_data, full_data)

	return full_data
}

func decodeText(data []byte) string {
	encoding := data[0]
	//fmt.Println(encoding)

	data = data[1:]

	switch encoding {
	case 0x00:
		// ISO-8859-1.
		runes := make([]rune, len(data))
		for i, b := range data {
			runes[i] = rune(b)
		}
		return strings.TrimRight(string(runes), "\x00")
	case 0x01:
		// UTF-16 with BOM.
		if len(data) < 2 {
			return ""
		}

		littleEndian := data[0] == 0xff && data[1] == 0xfe
		data = data[2:]

		return decodeUTF16(data, littleEndian)
	case 0x02:
		// UTF-16BE.
		return decodeUTF16(data, false)
	case 0x03:
		// UTF-8.
		return strings.TrimRight(string(data), "\x00")
	default:
		return ""
	}
}

func decodeUTF16(data []byte, littleEndian bool) string {
	if len(data)%2 != 0 {
		data = data[:len(data)-1]
	}

	u16 := make([]uint16, len(data)/2)

	for i := range u16 {
		if littleEndian {
			u16[i] = binary.LittleEndian.Uint16(data[i*2:])
		} else {
			u16[i] = binary.BigEndian.Uint16(data[i*2:])
		}
	}

	return strings.TrimRight(string(utf16.Decode(u16)), "\x00")
}
