package main

import (
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"math/rand"
	"os"
)

const characterWidth = 7
const characterHeight = 13
const characterCount = 8

const imageSizeMultiplier = 4

const imageWidth = characterCount*characterWidth + 1
const imageHeight = characterHeight + 44

func generateCaptcha() (string, *image.Gray) {

	drawer, img := createImage()
	generatedString := generateRandomString(characterCount)
	drawer.DrawString(generatedString)

	img = resizeImage(img)
	img = distortImage(img)
	img = addNoise(img)

	return generatedString, img
}

func addNoise(img *image.Gray) *image.Gray {
	imageArray := imageToArray(img)

	for y, h := range imageArray {
		for x := range h {
			img.Set(x, y, color.Gray{Y: imageArray[y][x] + uint8(rand.Intn(255-0)+0)/2})
		}
	}
	return img
}

func distortImage(src *image.Gray) *image.Gray {

	sineModifier := 0.1 + rand.Float64()*(.9-0.1)
	var sineCollection []float64

	// generateCaptcha height modifier for each column
	for i := 0.1; i < imageHeight*imageSizeMultiplier*20; i += .05 {
		sineCollection = append(sineCollection, math.Sin(i*sineModifier))
	}

	imageArray := imageToArray(src)
	modifierRange := 20.0

	for y, h := range imageArray {
		for x := range h {
			if imageArray[y][x] != 255 {
				src.Set(x, y, color.Gray{Y: 255})
				arrayMod := int(math.Round(modifierRange*sineCollection[x])) + int(modifierRange)
				src.Set(x, (((arrayMod+imageHeight*imageSizeMultiplier)+y)/2)-(rand.Intn((2*imageSizeMultiplier/2)-0)+0)-imageHeight/2, color.Gray{Y: 100})
			}

		}
	}

	return src
}

func resizeImage(src image.Image) *image.Gray {
	dst := image.NewGray(image.Rect(0, 0, imageWidth*imageSizeMultiplier, imageHeight*imageSizeMultiplier))
	draw.NearestNeighbor.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)
	return dst
}

func createImage() (*font.Drawer, *image.Gray) {

	img := image.NewGray(image.Rect(0, 0, imageWidth, imageHeight))
	draw.Draw(img, img.Bounds(), image.White, image.Point{}, draw.Src)

	face := basicfont.Face7x13
	drawer := &font.Drawer{
		Dst:  img,
		Src:  image.Black,
		Face: face,
		Dot:  fixed.Point26_6{X: characterWidth * 9, Y: characterHeight * 160},
	}

	return drawer, img
}

func generateRandomString(captchaLength int) string {

	characterSet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var charArray = make([]uint8, captchaLength)

	for i := 0; i < captchaLength; i++ {
		charArray[i] = characterSet[rand.Intn(len(characterSet))]
	}
	captchaCharacters := string(charArray[:])

	return captchaCharacters
}

func saveImage(img image.Gray) {
	f, err := os.Create("hello.png")
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	err = png.Encode(f, &img)
	if err != nil {
		log.Fatalf("Failed to encode image: %v", err)
	}
}

func imageToArray(img *image.Gray) [][]uint8 {

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var pixels [][]uint8
	for y := 0; y < height; y++ {
		var row []uint8
		for x := 0; x < width; x++ {
			row = append(row, grayToUint8(img.At(x, y).RGBA()))
		}
		pixels = append(pixels, row)
	}

	return pixels
}

func grayToUint8(r uint32, g uint32, b uint32, a uint32) uint8 {
	return uint8((r + g + b) / 3)
}
