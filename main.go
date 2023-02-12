package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Color used to draw face outlines
var drawColor color.Color = color.RGBA{255, 0, 0, 255}

type FaceRectangle struct {
	Width  int `json:"width"`
	Height int `json:"height"`
	Left   int `json:"left"`
	Top    int `json:"top"`
}

// Prepare an HTTP request with headers and body
func createHttpRequest(method string, url string, file *os.File) (*http.Request, error) {
	body, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	req, _ := http.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Ocp-Apim-Subscription-Key", os.Getenv("API_KEY"))

	return req, err
}

// Load the image into a draw.Image
func loadImage(file *os.File) (draw.Image, error) {
	file.Seek(0, io.SeekStart)
	src, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	bounds := src.Bounds()
	img := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(img, img.Bounds(), src, bounds.Min, draw.Src)

	return img, err
}

// Save a copy of the modified image
func saveImage(img draw.Image, name string) {
	file, _ := os.Create(fmt.Sprintf("%s.png", name))
	defer file.Close()

	png.Encode(file, img)
}

// Outline the detected face with a rectangle
func drawRectangle(img draw.Image, rectangle FaceRectangle) {
	for i := rectangle.Left; i <= rectangle.Left+rectangle.Width; i++ {
		img.Set(i, rectangle.Top, drawColor)
		img.Set(i, rectangle.Top+rectangle.Height, drawColor)
	}

	for i := rectangle.Top; i <= rectangle.Top+rectangle.Height; i++ {
		img.Set(rectangle.Left, i, drawColor)
		img.Set(rectangle.Left+rectangle.Width, i, drawColor)
	}
}

// Draw rectangles for all detected faces
func drawRectangles(img draw.Image, faces []map[string]FaceRectangle) {
	for _, face := range faces {
		drawRectangle(img, face["faceRectangle"])
	}
}

func main() {
	godotenv.Load()

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	req, _ := createHttpRequest("POST", os.Getenv("ENDPOINT"), file)
	defer req.Body.Close()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, _ := io.ReadAll(res.Body)

	var jsonRes []map[string]FaceRectangle
	json.Unmarshal(body, &jsonRes)

	img, _ := loadImage(file)
	drawRectangles(img, jsonRes)
	saveImage(img, fmt.Sprintf("%s_output", strings.TrimSuffix(file.Name(), ".png")))
}
