// by mo2zie

package goasciiart

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"log"
	"math"
	"os"
	"reflect"
	"strings"

	"github.com/golang/freetype/truetype"
	"github.com/nfnt/resize"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

const ascii string = "MND8OZ$7I?+=~:,.."

func ScaleImage(img image.Image, w int) (image.Image, int, int) {
	sz := img.Bounds()
	h := (sz.Max.Y * w * 10) / (sz.Max.X * 23)
	img = resize.Resize(uint(w), uint(h), img, resize.Lanczos3)
	return img, w, h
}

//From original goasciiart package
func Convert2Ascii(img image.Image, w, h int) []byte {
	table := []byte(ascii)
	buf := new(bytes.Buffer)
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			g := color.GrayModel.Convert(img.At(j, i))
			y := reflect.ValueOf(g).FieldByName("Y").Uint()
			pos := int(y * 16 / 255)
			_ = buf.WriteByte(table[pos])
		}
		_ = buf.WriteByte('\n')
	}
	return buf.Bytes()
}

var (
	//dpi of image
	dpi float64 = 150
	//Whether we want hinting or not
	hinting = "none"
	//Font size being rasterized onto image
	size float64 = 5
	//Line spacing of text
	spacing = 1.5
)

//Grayscale an image. Might be used later.
func rgbaToGray(img image.Image) *image.Gray {
	bounds := img.Bounds()
	gray := image.NewGray(bounds)
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			var rgba = img.At(x, y)
			gray.Set(x, y, rgba)
		}
	}
	return gray
}

//Convert string to slice by just splitting at new line.
func stringToSlice(input string) []string {
	output := strings.Split(input, "\n")
	output = output[1:]
	return output
}

//Simple algo for image height
func imageHeight(sliceLen int) int {
	height := 0
	for i := 0; i < sliceLen; i++ {
		height += 16
	}
	return height
}

//TextToImage - Converts Ascii string to an image. IT IS NOT ENCODED; It returns an image.Image.
//The caller is responsible for encoding how they please or just using the built in convert and save
//functions.
func TextToImage(ascii string) (image.Image, error) {
	// Read the font data.
	fontBytes, err := Asset("luximr.ttf")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	asciiSlice := stringToSlice(ascii)
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Draw the background and the guidelines.
	fg, bg := image.White, image.Black
	imgW := 770
	imgH := imageHeight(len(asciiSlice))
	rgba := image.NewRGBA(image.Rect(0, 0, imgW, imgH))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)

	// Draw the text.
	h := font.HintingNone
	switch hinting {
	case "full":
		h = font.HintingFull
	}

	d := &font.Drawer{
		Dst: rgba,
		Src: fg,
		Face: truetype.NewFace(f, &truetype.Options{
			Size:    size,
			DPI:     dpi,
			Hinting: h,
		}),
	}

	y := 10 + int(math.Ceil(size*dpi/72))
	dy := int(math.Ceil(size * spacing * dpi / 72))
	for _, s := range asciiSlice {
		d.Dot = fixed.P(10, y)
		d.DrawString(s)
		y += dy
	}
	return rgba, nil
}

//Encodes image.Image into png format.
func ConvertToImage(rgba image.Image) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := png.Encode(buf, rgba)
	if err != nil {
		return nil, err
	}
	log.Println("File successfully converted to Png format.")
	return buf.Bytes(), nil
}

//Saves png image data to file.
func SaveImage(fileName string, imageData []byte) error {
	outfile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer outfile.Close()
	w := bufio.NewWriter(outfile)
	_, err = w.Write(imageData)
	if err != nil {
		return err
	}
	w.Flush()
	fmt.Printf("File '%s' successfully written to.\n", fileName)
	return nil
}
