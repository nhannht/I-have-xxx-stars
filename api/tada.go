package api

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/fogleman/gg"
	"golang.org/x/image/font/opentype"
	"image"
	"image/color"
	"net/http"
	"os"
)

const (
	W        = 1000
	H        = 1000
	textSize = 48
	textDPI  = 72
)

//go:embed larvar-shit-without-background.png
var larvarShitPNG []byte

//go:embed JetBrainsMonoNL-Bold.ttf
var jetBrainsMonoTTF []byte

func TadaHandler(w http.ResponseWriter, r *http.Request) {
	dc := gg.NewContext(W, H)

	// Load the embedded image
	img, _, err := image.Decode(bytes.NewReader(larvarShitPNG))
	if err != nil {
		http.Error(w, "Unable to load embedded image.", http.StatusInternalServerError)
		return
	}

	// Calculate the position to center the image
	imgWidth := img.Bounds().Dx()
	imgHeight := img.Bounds().Dy()
	x := (dc.Width() - imgWidth) / 2
	y := (dc.Height() - imgHeight) / 2

	// Draw the embedded image onto the context
	dc.DrawImage(img, x, y)

	// Load the embedded font
	font, err := opentype.Parse(jetBrainsMonoTTF)
	if err != nil {
		http.Error(w, "Unable to parse font.", http.StatusInternalServerError)
		return
	}

	face, err := opentype.NewFace(font, &opentype.FaceOptions{
		Size: textSize,
		DPI:  textDPI,
	})
	if err != nil {
		http.Error(w, "Unable to create font face.", http.StatusInternalServerError)
		return
	}
	dc.SetFontFace(face)

	// Calculate the position to draw the string

	textX := float64(dc.Width()) / 2
	textY := float64(y + imgHeight + 50) // Adjust the 50 to add some padding below the image

	// get the GITHUB_TOKEN
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		http.Error(w, "GitHub token is not set", http.StatusUnauthorized)
		return
	}

	// Fetch the stargazer count from the GitHub API
	repo := "nhannht/i-have-xxx-stars"
	url := fmt.Sprintf("https://api.github.com/repos/%s", repo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, "Failed to create request to GitHub API", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Authorization", "token "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to send request to GitHub API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to fetch repository details from GitHub API", http.StatusInternalServerError)
		return
	}

	var repoDetails struct {
		StargazersCount int `json:"stargazers_count"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&repoDetails); err != nil {
		http.Error(w, "Failed to parse response from GitHub API", http.StatusInternalServerError)
		return
	}

	// Draw the string  under the image
	dc.SetRGB(255, 255, 255)
	dc.DrawStringAnchored(fmt.Sprintf("Stargazers: %d", repoDetails.StargazersCount), textX, textY, 0.5, 0.5)

	// add gradient to text
	// get the context as an alpha mask
	mask := dc.AsMask()

	// set a gradient
	g := gg.NewLinearGradient(textX/2, textY, 1000, textY+48)
	g.AddColorStop(0, color.RGBA{R: 255, A: 255})
	g.AddColorStop(1, color.RGBA{B: 255, A: 255})
	dc.SetFillStyle(g)

	// using the mask, fill the context with the gradient
	if err := dc.SetMask(mask); err != nil {
		http.Error(w, "Unable to set mask.", http.StatusInternalServerError)
		return
	}
	dc.DrawRectangle(textX/2, textY-10, W, textY+textSize) // adjust the rectangle with minus 10 in order to fill all text
	dc.Fill()

	// Set the content type to image/png
	w.Header().Set("Content-Type", "image/png")

	// Encode the image to PNG and write it to the response writer
	if err := dc.EncodePNG(w); err != nil {
		http.Error(w, "Unable to encode image.", http.StatusInternalServerError)
	}
}
