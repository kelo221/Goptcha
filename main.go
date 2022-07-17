package main

import (
	"fmt"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"math"
	"math/rand"
	"os"
	"time"
)

const imageHeight = 20
const imageWidth = 48
const capchaLenght = 6
const scaleFactor = 20
const noiseLevel = 0
const pngOutput = true

type Pixel struct {
	R int
	G int
	B int
	A int
}

func addLabel(img *image.RGBA, x, y int, label string) {
	col := color.RGBA{
		R: uint8(50 + rand.Intn(120-50)),
		G: uint8(50 + rand.Intn(120-50)),
		B: uint8(50 + rand.Intn(120-50)),
		A: 255,
	}
	point := fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}

func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
}

func pixelToRGBA(pixel Pixel) color.RGBA {
	return color.RGBA{
		R: uint8(pixel.R * 257),
		G: uint8(pixel.G * 257),
		B: uint8(pixel.B * 257),
		A: uint8(pixel.A * 257),
	}
}

func imageToArray(img *image.RGBA) [][]Pixel {

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var pixels [][]Pixel
	for y := 0; y < height; y++ {
		var row []Pixel
		for x := 0; x < width; x++ {
			row = append(row, rgbaToPixel(img.At(x, y).RGBA()))
		}
		pixels = append(pixels, row)
	}

	return pixels
}

func main() {

	rand.Seed(time.Now().UnixNano())

	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var c [capchaLenght]uint8

	for i := 0; i < capchaLenght; i++ {
		c[i] = charset[rand.Intn(len(charset))]
	}

	capchaCharacters := string(c[:])

	fmt.Println(capchaCharacters)

	img := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))
	addLabel(img, capchaLenght-2, imageHeight/2+10, capchaCharacters)

	var fileExtension string

	if pngOutput {
		fileExtension = "png"
	} else {
		fileExtension = "jpg"
	}

	filename := "capcha." + fileExtension

	output, _ := os.Create(filename)
	defer output.Close()

	dst := image.NewRGBA(image.Rect(0, 0, imageWidth*scaleFactor, imageHeight*scaleFactor))
	draw.NearestNeighbor.Scale(dst, dst.Rect, img, img.Bounds(), draw.Over, nil)

	v := 0.1 + rand.Float64()*(.9-0.1)
	fmt.Println("sine variable", v)

	var sineCollection []float64

	// generate sine values into an array
	for i := 0.1; i < 480; i += .05 {
		sineCollection = append(sineCollection, math.Sin(i*v))
	}

	array := imageToArray(dst)

	// 24 total pixel to move
	// 1 to -1, range of 2
	// each unit = 8 pixels
	// if the pixel is not empty get the sine value for the current column
	// move the current pixels up or down depending on sine value
	// clear the old pixels
	// +8 to offset to 0

	modifierRange := math.Round(240) / 10

	// You cannot modify constants
	arrayCountHorizontal := imageWidth
	arrayCountVertical := imageHeight

	fmt.Println(arrayCountHorizontal * scaleFactor)

	for x := 0; x < (arrayCountHorizontal * scaleFactor); x++ {
		for y := 0; y < (arrayCountVertical * scaleFactor); y++ {
			if array[y][x] != (Pixel{0, 0, 0, 0}) {
				sineModifier := int(math.Round(modifierRange*sineCollection[x])) + int(modifierRange)
				array[((sineModifier+arrayCountVertical*4)+y)/2][x] = array[y][x]
				array[y][x] = Pixel{0, 0, 0, 0}
			}
		}
	}

	if noiseLevel > 0 {
		for x := 0; x < (arrayCountHorizontal * scaleFactor); x++ {
			for y := 0; y < (arrayCountVertical * scaleFactor); y++ {
				array[y][x].R = (array[y][x].R + noiseLevel*(0+rand.Intn(255-0))) / (1 + noiseLevel)
				array[y][x].G = (array[y][x].G + noiseLevel*(0+rand.Intn(255-0))) / (1 + noiseLevel)
				array[y][x].B = (array[y][x].B + noiseLevel*(0+rand.Intn(255-0))) / (1 + noiseLevel)
				array[y][x].A = 255
			}
		}
	}

	newImage := image.NewRGBA(image.Rect(0, 0, img.Bounds().Max.X*scaleFactor, img.Bounds().Max.Y*scaleFactor))

	// convert back to RGBA
	for x := 0; x < imageHeight*scaleFactor; x++ {
		for y := 0; y < imageWidth*scaleFactor; y++ {
			newImage.SetRGBA(y, x, pixelToRGBA(array[x][y]))
		}
	}

	// Set the expected size that you want:
	dst2 := image.NewRGBA(image.Rect(0, 0, newImage.Bounds().Max.X/3, newImage.Bounds().Max.Y/3))

	// Resize:
	draw.NearestNeighbor.Scale(dst2, dst2.Rect, newImage, newImage.Bounds(), draw.Over, nil)

	if !pngOutput {
		jpeg.Encode(output, dst2, &jpeg.Options{Quality: 90})
	} else {
		png.Encode(output, dst2)
	}

}
