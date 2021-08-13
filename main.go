package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"

	imggg "github.com/disintegration/imaging"
)

var ASCIISTR = `$@B%8&WM#*oahkbdpqwmZO0QLCJUYXzcvunxrjft7/\|()1{}[]?-_+~<>i!lI;:,"^'.  `

func GetImage() (image.Image, int) {
	type urlStr struct {
		URL string `json:"url"`
	}
	urlstruct := urlStr{}
	res, err := http.Get("https://api.waifu.pics/sfw/neko")
	if err != nil {
		log.Println("GET image error#1:", err)
		return nil, 0
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("GET image parse error:", err)
		return nil, 0
	}
	if err = json.Unmarshal(body, &urlstruct); err != nil {
		log.Println("GET unmarshal error:", err)
		return nil, 0
	}
	res.Body.Close()
	res, err = http.Get(urlstruct.URL)
	if err != nil {
		log.Println("GET image error#2:", err)
		return nil, 0
	}
	defer res.Body.Close()
	img, _, err := image.Decode(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return img, 180

}

func ScaleImage(img image.Image, w int) (*image.NRGBA, int, int) {
	sz := img.Bounds()
	h := (sz.Max.Y * w * 10) / (sz.Max.X * 16)
	img1 := imggg.Grayscale(img)
	img1 = imggg.Resize(img1, w, h, imggg.Lanczos)
	return img1, w, h
}

func Convert2Ascii(img *image.NRGBA, w, h int) []byte {
	table := []byte(ASCIISTR)
	bufr := new(bytes.Buffer)
	var max int = 0
	var min int = 70
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			r, _, _, a := img.At(j, i).RGBA()
			avg := (float64(r))
			pos := int(avg * 18100 / float64(a) / 255)
			if pos < min {
				min = pos
			}
			if pos > max {
				max = pos
			}
			if err := bufr.WriteByte(table[pos]); err != nil {
				log.Println("Convert2Ascii write error:", err)
				return nil
			}
		}
		if err := bufr.WriteByte('\n'); err != nil {
			log.Println("Convert2Ascii write error:", err)
			return nil
		}

	}
	log.Println(min, max)
	return bufr.Bytes()
}

func main() {
	// Convert2Ascii(ScaleImage(GetImage()))

	r := Convert2Ascii(ScaleImage(GetImage()))
	fmt.Println(string(r))

}
