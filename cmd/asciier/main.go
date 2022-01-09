package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/THasthika/asciier-go/converter"
)

func main() {

	width := flag.Int("w", 0, "output width")
	height := flag.Int("h", 0, "output height")
	outFile := flag.String("out", "./out.txt", "output file path")

	flag.Parse()

	var inFile string
	fmt.Scanln(&inFile)

	img, err := converter.NewAsciierImageFromImageFile(string(inFile))
	if err != nil {
		panic(err)
	}

	if *width != 0 || *height != 0 {
		img.Resize(*width, *height)
	}

	r := img.ConvertToAscii()

	err = os.WriteFile(*outFile, []byte(r), 0644)
	if err != nil {
		panic(err)
	}
}
