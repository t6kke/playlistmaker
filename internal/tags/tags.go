package tags

import (
	"fmt"
	"io"
	"os"
)

func Test(path string) {
	f, err := os.Open(path)
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
	fmt.Printf("file header type: %s, and header version: %d.%d, flags: %08b and tag size of: %d\n", string(header[0:3]), header[3], header[4], header[5], header[6:10])
}
