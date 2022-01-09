package converter

import (
	"image"
	"os"
	"strings"

	"github.com/gosuri/uiprogress"
	"github.com/nfnt/resize"

	"image/color"
	_ "image/jpeg"
)

type AsccierImageSize struct {
	width  int
	height int
}

type AsciierImage struct {
	img       *image.Image
	asciiText *string
	size      *AsccierImageSize
}

type colorWithPosition struct {
	c color.Color
	p image.Point
}

type byteWithPosition struct {
	b    byte
	p    image.Point
	rank int
}

type asciiImageReconstructor struct {
	d      [][]byte
	width  int
	height int
}

const maxWorkers = 100000

func NewAsciierImageFromImageFile(imageFile string) (*AsciierImage, error) {
	reader, err := os.Open(imageFile)
	if err != nil {
		return nil, err
	}
	m, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}
	return &AsciierImage{
		img:       &m,
		asciiText: nil,
		size: &AsccierImageSize{
			width:  m.Bounds().Size().X,
			height: m.Bounds().Size().Y,
		},
	}, nil
}

func (asciierImage *AsciierImage) GetSize() *AsccierImageSize {
	return asciierImage.size
}

func (asciierImage *AsciierImage) Resize(width int, height int) {
	asciierImage.size.width = width
	asciierImage.size.height = height
}

func colorToAsciiWorker(jobs <-chan *colorWithPosition, results chan<- *byteWithPosition, maxX int) {
	for j := range jobs {
		r, g, b, _ := (*j).c.RGBA()
		gray := (r + g + b) / 3
		grayReal := float32(gray) / (0xffff)
		ret := pixelToAscii(grayReal)
		results <- &byteWithPosition{
			b:    ret,
			p:    j.p,
			rank: (j.p.X * maxX) + j.p.Y,
		}
	}
}

func (asciierImage *AsciierImage) ConvertToAscii() string {
	resizedImage := resize.Resize(uint(asciierImage.size.width), uint(asciierImage.size.height), *asciierImage.img, resize.Lanczos3)

	maxX := resizedImage.Bounds().Size().X
	maxY := resizedImage.Bounds().Size().Y

	pixelCount := maxX * maxY

	jobs := make(chan *colorWithPosition, pixelCount)
	results := make(chan *byteWithPosition, pixelCount)

	var workerCount = maxWorkers
	if maxWorkers > pixelCount {
		workerCount = pixelCount
	}

	uiprogress.Start()
	jobAddBar := uiprogress.AddBar(pixelCount).AppendCompleted().PrependElapsed()
	resultGetBar := uiprogress.AddBar(pixelCount).AppendCompleted().PrependElapsed()

	reconstAwait := make(chan asciiImageReconstructor, 1)
	go func() {
		r := asciiImageReconstructor{}
		r.width = maxX
		r.height = maxY
		r.d = make([][]byte, r.height)
		for i := 0; i < r.height; i++ {
			r.d[i] = []byte(strings.Repeat(" ", r.width))
		}
		reconstAwait <- r
	}()

	for i := 0; i < workerCount; i++ {
		go colorToAsciiWorker(jobs, results, maxX)
	}

	for x := 0; x < maxX; x++ {
		for y := 0; y < maxY; y++ {
			jobAddBar.Incr()
			jobs <- &colorWithPosition{
				c: resizedImage.At(x, y),
				p: image.Point{X: x, Y: y},
			}
		}
	}

	reconst := <-reconstAwait

	for i := 0; i < pixelCount; i++ {
		resultGetBar.Incr()
		ret := <-results
		line := reconst.d[ret.p.Y]
		line[ret.p.X] = ret.b
	}

	close(jobs)

	makeStringBar := uiprogress.AddBar(len(reconst.d)).AppendCompleted().PrependElapsed()

	ret := make([]string, reconst.height)
	for i := range reconst.d {
		makeStringBar.Incr()
		ret[i] = string(reconst.d[i])
	}

	return strings.Join(ret, "\n")
}
