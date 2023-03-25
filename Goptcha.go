package Goptcha

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

// characterWidth is the width of default font characters in pixels.
const characterWidth = 7

// characterHeight is the height of default font characters in pixels.
const characterHeight = 13

type Config struct {
	CharacterCount      int
	ImageSizeMultiplier int
	imageWidth          int
	imageHeight         int
	CharSet             string
	Opacity             uint8
	NoiseModifier       uint8
}

var cDefs = Config{
	CharacterCount:      8,
	ImageSizeMultiplier: 4,
	imageWidth:          8*characterWidth + 1,
	imageHeight:         characterHeight + 22,
	CharSet:             "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
	Opacity:             100,
	NoiseModifier:       10,
}

// Configure modifies default captcha generation settings. (OPTIONAL)
func Configure(config *Config) {

	if config.ImageSizeMultiplier == 0 {
		log.Fatal("Image size multiplier cannot be 0")
	}

	if config.CharSet == "" {
		log.Fatal("You must specify at least one character.")
	} else {
		cDefs.CharacterCount = config.CharacterCount
		cDefs.CharSet = config.CharSet
		cDefs.imageWidth = config.CharacterCount*characterWidth + 1
	}

	cDefs.Opacity = config.Opacity

}

func GenerateCaptcha() (string, *image.Gray) {

	drawer, img := createImage()
	generatedString := generateRandomString(cDefs.CharacterCount)
	drawer.DrawString(generatedString)

	img = resizeImage(img)
	img = distortImage(img)

	img = addNoise(img)

	return generatedString, img
}

func addNoise(img *image.Gray) *image.Gray {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			currentColor := img.GrayAt(x, y).Y
			noiseValue := uint8(rand.Intn(255)) / 2
			img.SetGray(x, y, color.Gray{Y: currentColor + noiseValue})
		}
	}

	return img
}

func distortImage(src *image.Gray) *image.Gray {
	sineModifier := 0.1 + rand.Float64()*(.32-0.1)
	var sineCollection []float64

	// generate height modifier for each column
	for i := 0.1; i < float64(cDefs.imageHeight*cDefs.ImageSizeMultiplier)*20; i += .05 {
		sineCollection = append(sineCollection, math.Sin(i*sineModifier))
	}

	imageArray := imageToArray(src)
	modifierRange := 35.0

	for y, h := range imageArray {
		for x := range h {
			if imageArray[y][x] != 255 {
				src.Set(x, y, color.Gray{Y: 255})
				arrayMod := int(math.Round(modifierRange*sineCollection[x])) + int(modifierRange)
				newY := y + arrayMod - cDefs.ImageSizeMultiplier*cDefs.imageHeight/2
				src.Set(x, newY, color.Gray{Y: cDefs.Opacity})
			}
		}
	}

	return src
}

func resizeImage(src image.Image) *image.Gray {
	dst := image.NewGray(image.Rect(0, 0, cDefs.imageWidth*cDefs.ImageSizeMultiplier, cDefs.imageHeight*cDefs.ImageSizeMultiplier))
	draw.NearestNeighbor.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)
	return dst
}

func createImage() (*font.Drawer, *image.Gray) {

	img := image.NewGray(image.Rect(0, 0, cDefs.imageWidth, cDefs.imageHeight))
	draw.Draw(img, img.Bounds(), image.White, image.Point{}, draw.Src)

	face := basicfont.Face7x13
	drawer := &font.Drawer{
		Dst:  img,
		Src:  image.Black,
		Face: face,
		Dot:  fixed.Point26_6{X: characterWidth * 9, Y: fixed.Int26_6(characterHeight * cDefs.imageHeight * 3)},
	}

	return drawer, img
}

func generateRandomString(captchaLength int) string {

	var charArray = make([]uint8, captchaLength)

	for i := 0; i < captchaLength; i++ {
		charArray[i] = cDefs.CharSet[rand.Intn(len(cDefs.CharSet))]
	}
	captchaCharacters := string(charArray[:])

	return captchaCharacters
}

func SaveImage(img image.Gray) {
	f, err := os.Create("hello.png")
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	err = png.Encode(f, &img)
	if err != nil {
		log.Fatalf("Failed to encode image: %v", err)
	}
}

func grayToUint8(r uint32, g uint32, b uint32, a uint32) uint8 {
	return uint8((r + g + b) / 3)
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
