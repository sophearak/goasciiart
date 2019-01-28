package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	"log"
	"os"
	"strings"

	ascii "github.com/cantasaurus/goasciiart"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Init() (image.Image, int) {
	widthPtr := flag.Int("w", 120, "Use -w <num>")
	fpath := flag.String("p", "test.jpg", "Use -p <filesource>")
	flag.Parse()
	f, err := os.Open(*fpath)
	check(err)
	img, _, err := image.Decode(f)
	check(err)
	f.Close()
	return img, *widthPtr
}

func main() {
	asciiBytes := ascii.Convert2Ascii(ascii.ScaleImage(Init()))
	asciiText := string(asciiBytes)
	fmt.Print(asciiText)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Would you like to save this as an image? (Y or N): ")
	text, _ := reader.ReadString('\n')

	if strings.ToLower(strings.TrimSpace(text)) == "yes" || strings.ToLower(strings.TrimSpace(text)) == "y" {
		fmt.Print("What would you like to name it?: ")
		text, _ := reader.ReadString('\n')
		rgba, err := ascii.TextToImage(asciiText)
		if err != nil {
			log.Println(err)
			return
		}
		text = strings.TrimSuffix(text, "\n")
		imageData, err := ascii.ConvertToImage(rgba)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = ascii.SaveImage(text, imageData)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
