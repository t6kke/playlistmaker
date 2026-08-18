package tags

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf16"
)

func Test(path string) {
	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		return
	}
	defer f.Close()

	header := make([]byte, 10)
	if _, err := io.ReadFull(f, header); err != nil {
		return
	}

	fmt.Println(path)
	fmt.Println(header[0:10])
	for i := range 10 {
		fmt.Printf("byte: %d values: ", i)
		fmt.Print(header[i])
		fmt.Printf(" %08b\n", header[i])
		fmt.Printf("decimal: %d --- hex: %02x --- binary: %08b --- char: %c\n", header[i], header[i], header[i], header[i])
	}

	tag_size := fix_tagsize_staring_bits(header[6:10])
	fmt.Println(tag_size)

	fmt.Printf("file header type: %s, and header version: %d.%d, flags: %08b and tag size of: %d\n", string(header[0:3]), header[3], header[4], header[5], header[6:10])

	if int(header[3]) == 3 {
		tag_data := make([]byte, tag_size)
		_, err := io.ReadFull(f, tag_data) //TODO handle error
		if err != nil {
			return
		}

		frameHeaderSize := 10
		for pos := 0; pos+frameHeaderSize <= len(tag_data); {
			id := string(tag_data[pos : pos+4])

			//stop after final tag item
			if id == "\x00\x00\x00\x00" {
				break
			}

			size := int(binary.BigEndian.Uint32(tag_data[pos+4 : pos+8]))
			if size <= 0 || pos+10+size > len(tag_data) {
				break
			}

			switch id {
			case "TIT2":
				fmt.Println(id)
				fmt.Println("Title")
				frame := tag_data[pos+10 : pos+10+size]
				fmt.Println(frame)
				decoded_frame_info := decodeText(frame)
				fmt.Println(decoded_frame_info)
			case "TPE1":
				fmt.Println(id)
				fmt.Println("Artist")
				frame := tag_data[pos+10 : pos+10+size]
				fmt.Println(frame)
				decoded_frame_info := decodeText(frame)
				fmt.Println(decoded_frame_info)
			case "TALB":
				fmt.Println(id)
				fmt.Println("Album")
				frame := tag_data[pos+10 : pos+10+size]
				fmt.Println(frame)
				decoded_frame_info := decodeText(frame)
				fmt.Println(decoded_frame_info)
			case "TLEN":
				fmt.Println(id)
				fmt.Println("Duration")
				frame := tag_data[pos+10 : pos+10+size]
				fmt.Println(frame)
				decoded_frame_info := decodeText(frame)
				fmt.Println(decoded_frame_info)
			case "TRCK":
				fmt.Println(id)
				fmt.Println("Track")
				frame := tag_data[pos+10 : pos+10+size]
				fmt.Println(frame)
				decoded_frame_info := decodeText(frame)
				fmt.Println(decoded_frame_info)
			}

			pos += 10 + size
		}
	}
}

func fix_tagsize_staring_bits(b []byte) int {
	//fmt.Printf("%08b  %08b %08b\n",b[2], b[2]&0x7f, int(b[2]&0x7f)<<7)
	full_data := int(b[0]&0x7f)<<21 | int(b[1]&0x7f)<<14 | int(b[2]&0x7f)<<7 | int(b[3]&0x7f)
	//fmt.Printf("%d, %08b", full_data, full_data)

	return full_data
}

func decodeText(data []byte) string {
	encoding := data[0]
	fmt.Println(encoding)

	data = data[1:]

	switch encoding {
	case 0x00:
		// ISO-8859-1.
		/*runes := make([]rune, len(data))
		for i, b := range data {
			runes[i] = rune(b)
		}
		return strings.TrimRight(string(runes), "\x00")*/
		//TODO
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
		//return decodeUTF16(data, false)
		//TODO

	case 0x03:
		// UTF-8.
		//return strings.TrimRight(string(data), "\x00")
		//TODO
	default:
		return ""
	}

	return ""
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
