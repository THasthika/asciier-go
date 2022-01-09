package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/THasthika/asciier-go/converter"
)

func main() {

	width := flag.Int("w", 0, "output width")
	height := flag.Int("h", 0, "output height")
	outFile := flag.String("o", "", "output file path")

	flag.Parse()

	var inFile string

	if flag.NArg() == 0 {
		reader := bufio.NewReader(os.Stdin)
		var err error
		inFile, err = reader.ReadString('\n')
		if err != nil {
			log.Fatalln("failed to read input")
		}
		inFile = strings.TrimSpace(inFile)
	} else {
		inFile = flag.Arg(0)
	}

	img, err := converter.NewAsciierImageFromImageFile(string(inFile))
	if err != nil {
		log.Fatalln(err)
	}
	if *width != 0 || *height != 0 {
		img.Resize(*width, *height)
	}

	r := img.ConvertToAscii()

	if *outFile == "" {
		fmt.Printf("%s", r)
	} else {
		err = os.WriteFile(*outFile, []byte(r), 0644)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
