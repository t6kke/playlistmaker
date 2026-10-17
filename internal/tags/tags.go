package tags

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
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
		_, err := io.ReadFull(f, header)
		if err != nil {
			return Metadata{}, err
		}

		fmt.Println("Header Data")
		fmt.Println(header[0:10])

		header_deatils := mp3_getHeaderDetails(header[0:10])

		fmt.Printf("file header type: %s --- and version: %d.%d --- flags: %+v and tag size of: %d\n", header_deatils.identifier, header_deatils.version, header_deatils.revision, header_deatils.flags, header_deatils.size)

		//ID3 v2.x metadata
		switch header_deatils.version {
		case 2:
			//TODO
			fmt.Println("   ", "mp3 id3 v2.2")
			//function call here
		case 3:
			fmt.Println("   ", "mp3 id3 v2.3")
			err := data.extract_mp3id3v23(&header_deatils, f)
			if err != nil {
				return Metadata{}, err
			}
		case 4:
			fmt.Println("   ", "mp3 id3 v2.4")
			err := data.extract_mp3id3v24(&header_deatils, f)
			if err != nil {
				return Metadata{}, err
			}
		default:
			//TODO don't just give up and retur, try for TAGv1 and then from path
			return Metadata{}, fmt.Errorf("not a valid mp3 ID3 version value: '%d.%d'", header_deatils.version, header_deatils.revision)
		}

		//TAG v1 metadata
		//TODO use this as a backup if data is missing or regular tags are not found
		full_data_size, _ := f.Seek(0, io.SeekEnd) //TODO handle error

		tag_v1_bytes := make([]byte, 128)
		_, _ = f.ReadAt(tag_v1_bytes, full_data_size-128) //TODO handle error
		v1_tags, _ := data.extract_mp3v1(tag_v1_bytes)    //TODO handle error
		fmt.Printf("    v1 tags for backup: %+v\n", v1_tags)

	case "flac":
		//TODO
		//code here

	default:
		return Metadata{}, fmt.Errorf("not a valid file for parsing: '%s'", file_name)
	}

	//TODO
	// this is for backup when metadata is missing or partial from the file.
	// 1. try some web API to search for the track to get metadata from web.
	// 2. use the path details for missing metadata, enable the existing function that is prepared.
	data.fillFromPath(strconv.Itoa(n_th_file_in_dir))

	//fmt.Printf("%+v\n", data)
	return data, nil
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
