package tags

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

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
			err := data.extract_mp3id3_v2_2(&header_deatils, header[5], f)
			if err != nil {
				return Metadata{}, err
			}
		case 3:
			fmt.Println("   ", "mp3 id3 v2.3")
			err := data.extract_mp3id3_v2_3(&header_deatils, f)
			if err != nil {
				return Metadata{}, err
			}
		case 4:
			fmt.Println("   ", "mp3 id3 v2.4")
			err := data.extract_mp3id3_v2_4(&header_deatils, f)
			if err != nil {
				return Metadata{}, err
			}
		default:
			//TODO don't just give up and retur, try for TAGv1 and then from path
			return Metadata{}, fmt.Errorf("not a valid mp3 ID3 version value: '%d.%d'", header_deatils.version, header_deatils.revision)
		}

		// TAGv1.1 metadata
		// this is for backup when metadata is missing or partial from the file.
		if data.hasEmptyFields() {
			tag_v1_bytes := make([]byte, 128)
			full_data_size, _ := f.Seek(0, io.SeekEnd)        //TODO handle error, but don't exit function
			_, _ = f.ReadAt(tag_v1_bytes, full_data_size-128) //TODO handle error, but don't exit function
			v1_tags, _ := data.extract_mp3v1(tag_v1_bytes)    //TODO handle error, but don't exit function
			data.fillFromMP3v1(&v1_tags)

		}

	case "flac":
		//TODO
		//code here

	default:
		return Metadata{}, fmt.Errorf("not a valid file for parsing: '%s'", file_name)
	}

	//TODO
	// this is for backup when metadata is missing or partial from the file.
	// 1. try some web API to search for the track to get metadata from web.
	// 2. use the path details for missing metadata, enable all other parsing is done
	/*if data.hasEmptyFields() {
		path_metadata := data.getMetadataFromPath()
		data.fillFromPath(&path_metadata, strconv.Itoa(n_th_file_in_dir))
	}*/

	//fmt.Printf("%+v\n", data)
	return data, nil
}

func getFileName(path string) string {
	i := strings.LastIndex(path, "/")
	return path[i+1:]
}

func (m *Metadata) hasEmptyFields() bool {
	v := reflect.ValueOf(m).Elem()

	for i := 0; i < v.NumField(); i++ {
		if v.Type().Field(i).Name == "Duration" {
			continue
		}
		if v.Field(i).IsZero() {
			return true
		}
	}

	return false
}

func (m *Metadata) fillFromMP3v1(tags *v1_tags) {
	fmt.Printf("    v1 tags for backup: %+v\n", tags)

	if m.Title == "" {
		m.Title = tags.title
	}
	if m.Creator == "" {
		m.Creator = tags.creator
	}
	if m.Album == "" {
		m.Album = tags.album
	}
	if m.Duration == "" {
		//TODO skipped, not required to have this field
	}
	if m.Track_nbr == "" {
		m.Track_nbr = strconv.Itoa(tags.track_nbr) //TODO handle data type in TAGv1 extraction
	}
}

func (m *Metadata) fillFromPath(pm *path_metadata, n_th_file_in_dir string) {
	fmt.Printf("    path tags for backup: %+v\n", pm)
	if m.Title == "" {
		m.Title = pm.song_name
	}
	if m.Creator == "" {
		m.Creator = pm.artist_name
	}
	if m.Album == "" {
		m.Album = pm.album_name
	}
	if m.Duration == "" {
		//TODO skipped, not required to have this field
	}
	if m.Track_nbr == "" {
		//fmt.Printf("       no track number in frame, using value from file name: %s\n", pm.track_number)
		m.Track_nbr = pm.track_number
	}
	if m.Track_nbr == "" {
		//fmt.Printf("       no track number in frame, nor in file name, using file count value: %d\n", n_th_file_in_dir)
		m.Track_nbr = n_th_file_in_dir
	}
}
