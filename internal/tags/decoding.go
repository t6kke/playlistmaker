package tags

import (
	"encoding/binary"
	"strings"
	"unicode/utf16"
)

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
