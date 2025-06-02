package main

import (
	"image"
	"image/jpeg"
	"os"
	"strings"
	"time"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/paint"
)

func ticker(w *app.Window, imgs *[]image.Image, imgIndex *int, imageDelay *int) {
	// Start ticker in a separate goroutine to update imgIndex and invalidate window
	ticker := time.NewTicker(time.Duration(*imageDelay) * time.Second)
	defer ticker.Stop()
	for {
		<-ticker.C
		*imgIndex = (*imgIndex + 1) % len(*imgs)
		w.Invalidate() // Request redraw
	}
}

func getEntries(imageDir *string) []os.DirEntry {
	entries, err := os.ReadDir(*imageDir)
	if err != nil {
		panic(err)
	}
	return entries
}

func getImages(imageDir *string, entries *[]os.DirEntry) []image.Image {
	var imgs []image.Image

	for _, entry := range *entries {
		if entry.IsDir() {
			continue
		}
		fileName := entry.Name()
		fileExt := fileName[len(fileName)-4:]
		if len(fileName) > 4 && strings.ToLower(fileExt) == ".jpg" {
			f, err := os.Open(*imageDir + "/" + entry.Name())
			if err != nil {
				panic(err)
			}
			img, err := jpeg.Decode(f)
			f.Close()
			if err != nil {
				panic(err)
			}
			imgs = append(imgs, img)
		}
	}
	if len(imgs) == 0 {
		panic("No images found in the directory")
	}
	return imgs
}

func loop(w *app.Window, imgs *[]image.Image, imgIndex *int) {
	var ops op.Ops
	for {
		evt := w.Event()
		switch typ := evt.(type) {
		case app.FrameEvent:
			{
				gtx := app.NewContext(&ops, typ)
				img := (*imgs)[*imgIndex]
				imgOp := paint.NewImageOp(img)
				imgBounds := img.Bounds()
				imgWidth := float32(imgBounds.Dx())
				imgHeight := float32(imgBounds.Dy())
				winWidth := float32(gtx.Constraints.Max.X)
				winHeight := float32(gtx.Constraints.Max.Y)
				// Calculate scale factors
				scaleX := winWidth / imgWidth
				scaleY := winHeight / imgHeight
				// Apply scaling transform
				m := op.Affine(f32.Affine2D{}.Scale(f32.Pt(0, 0), f32.Pt(scaleX, scaleY)))
				m.Add(gtx.Ops)
				imgOp.Add(gtx.Ops)
				paint.PaintOp{}.Add(gtx.Ops)
				typ.Frame(gtx.Ops)
			}
		case app.DestroyEvent:
			{
				os.Exit(0)
			}
		}
	}
}

func main() {
	imageDir := "./images" // Images found in the directory
	imgIndex := 0          // Start with the first
	imageDelay := 3600        // seconds between images

	go func() {
		// create new window
		w := new(app.Window)
		w.Option(app.Fullscreen.Option()) // Set fullscreen first

		entries := getEntries(&imageDir)
		imgs := getImages(&imageDir, &entries)
		go ticker(w, &imgs, &imgIndex, &imageDelay)

		// listen for events in the window
		loop(w, &imgs, &imgIndex)

	}()
	app.Main()
}
