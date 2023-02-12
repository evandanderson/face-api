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

type FaceRectangle struct {
	Width  int `json:"width"`
	Height int `json:"height"`
	Left   int `json:"left"`
	Top    int `json:"top"`
}

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

func loadImage(file *os.File) (draw.Image, error) {
	file.Seek(0, io.SeekStart)
	original, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	bounds := original.Bounds()
	img := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(img, img.Bounds(), original, bounds.Min, draw.Src)

	return img, err
}

func saveImage(img draw.Image, name string) {
	file, _ := os.Create(fmt.Sprintf("%s.png", name))
	defer file.Close()

	png.Encode(file, img)
}

func drawRectangle(img draw.Image, faces []map[string]FaceRectangle) {
	color := color.RGBA{255, 0, 0, 255}

	for _, face := range faces {
		rectangle := face["faceRectangle"]

		for i := rectangle.Left; i <= rectangle.Left+rectangle.Width; i++ {
			img.Set(i, rectangle.Top, color)
			img.Set(i, rectangle.Top+rectangle.Height, color)
		}

		for i := rectangle.Top; i <= rectangle.Top+rectangle.Height; i++ {
			img.Set(rectangle.Left, i, color)
			img.Set(rectangle.Left+rectangle.Width, i, color)
		}
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

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, _ := io.ReadAll(resp.Body)

	var res []map[string]FaceRectangle
	json.Unmarshal(body, &res)

	img, _ := loadImage(file)
	drawRectangle(img, res)
	saveImage(img, fmt.Sprintf("%s_output", strings.TrimSuffix(file.Name(), ".png")))
}
