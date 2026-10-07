package tags

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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

type Metadata struct {
	Path      string
	Title     string
	Creator   string
	Album     string
	Duration  string
	Track_nbr string
}

type headerDetails struct {
	identifier string
	version    int
	revision   int
	flags      headerFlags
	size       int
}

type headerFlags struct {
	unsynchronisation      bool
	extended_header        bool
	experimental_indicator bool
	footer_present         bool
}

func ExtractMetadata(path string, n_th_file_in_dir int, file_type string) (Metadata, error) {
	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		return Metadata{}, err
	}
	defer f.Close()

	fmt.Println(path)
	file_name := getFileName(path)
	fmt.Println(file_name)

	data := Metadata{}
	data.Path = path

	switch file_type {
	case "mp3":
		header := make([]byte, 10)
		if _, err := io.ReadFull(f, header); err != nil {
			return Metadata{}, err
		}

		fmt.Println("Header Data")
		fmt.Println(header[0:10])

		header_deatils := getHeaderDetails(header[0:10])

		fmt.Printf("file header type: %s --- and version: %d.%d --- flags: %+v and tag size of: %d\n", header_deatils.identifier, header_deatils.version, header_deatils.revision, header_deatils.flags, header_deatils.size)

		switch header_deatils.version {
		case 2:
			//TODO
			//function call here
		case 3:
			err := data.extract_mp3id3v23(&header_deatils, f)
			if err != nil {
				return Metadata{}, err
			}
		case 4:
			//TODO
			//function call here
		default:
			return Metadata{}, fmt.Errorf("not a valid mp3 ID3 version value: '%d.%d'", header_deatils.version, header_deatils.revision)
		}

	case "flac":
		//TODO
		//code here

	default:
		return Metadata{}, fmt.Errorf("not a valid file for parsing: '%s'", file_name)
	}

	//data.fillFromPath(path, strconv.Itoa(n_th_file_in_dir)) //TODO enable following once we have all types of metadata collections implemented, this is only backup when we have no more data to find

	//fmt.Printf("%+v\n", data)
	return data, nil
}

func getHeaderDetails(b []byte) headerDetails {
	result_data := headerDetails{}

	result_data.identifier = string(b[0:3])
	result_data.version = int(b[3])
	result_data.revision = int(b[4])
	result_data.flags = getHeaderFlags(b[5])
	result_data.size = getCorrectTagSize(b[6:10])

	//fmt.Printf("header struct: %+v\n", result_data)

	return result_data
}

func getCorrectTagSize(b []byte) int {
	//TODO validate lenght of bytes slice and return error if it's not 4

	//fmt.Printf("%08b  %08b %08b\n",b[2], b[2]&0x7f, int(b[2]&0x7f)<<7)
	full_data := int(b[0]&0x7f)<<21 | int(b[1]&0x7f)<<14 | int(b[2]&0x7f)<<7 | int(b[3]&0x7f)
	//fmt.Printf("%d, %08b", full_data, full_data)

	return full_data
}

func getHeaderFlags(b byte) headerFlags {
	h_flags := headerFlags{
		unsynchronisation:      false,
		extended_header:        false,
		experimental_indicator: false,
		footer_present:         false,
	}
	for i := 7; i >= 4; i-- {
		value := (b >> i) & 1
		switch i {
		case 7:
			if value == 1 {
				h_flags.unsynchronisation = true
			}
		case 6:
			if value == 1 {
				h_flags.extended_header = true
			}
		case 5:
			if value == 1 {
				h_flags.experimental_indicator = true
			}
		case 4:
			if value == 1 {
				h_flags.footer_present = true
			}
		}
	}
	return h_flags
}

func getFileName(path string) string {
	i := strings.LastIndex(path, "/")
	return path[i+1:]
}

func (m *Metadata) fillFromPath(n_th_file_in_dir string) {
	if m.Title == "" {
		m.Title = m.getMetadataFromPath().song_name
	}
	if m.Creator == "" {
		m.Creator = m.getMetadataFromPath().artist_name
	}
	if m.Album == "" {
		m.Album = m.getMetadataFromPath().album_name
	}
	if m.Duration == "" {
		//TODO skipped, not required to have this field
	}
	if m.Track_nbr == "" {
		metadata_from_file := m.getMetadataFromPath()
		//fmt.Printf("       no track number in frame, using value from file name: %s\n", metadata_from_file.track_number)
		m.Track_nbr = metadata_from_file.track_number
	}
	if m.Track_nbr == "" {
		//fmt.Printf("       no track number in frame, nor in file name, using file count value: %d\n", n_th_file_in_dir)
		m.Track_nbr = n_th_file_in_dir
	}
}
