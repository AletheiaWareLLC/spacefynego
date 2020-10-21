/*
 * Copyright 2020 Aletheia Ware LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package viewer

import (
	"bytes"
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"
	"github.com/AletheiaWareLLC/spacego"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log"
)

func GetViewer(meta *spacego.Meta, callback func(io.Writer) uint64) fyne.CanvasObject {
	switch meta.GetType() {
	case spacego.MIME_TYPE_TEXT_PLAIN:
		return NewTextPlain(callback)
	case spacego.MIME_TYPE_IMAGE_GIF:
		return NewImageGif(callback)
	case spacego.MIME_TYPE_IMAGE_JPEG:
		return NewImageJpeg(callback)
	case spacego.MIME_TYPE_IMAGE_PNG:
		return NewImagePng(callback)
	}
	return nil
}

func NewTextPlain(callback func(io.Writer) uint64) fyne.CanvasObject {
	// Create label to hold text
	label := &widget.Label{
		Wrapping: fyne.TextWrapWord,
	}
	scroller := widget.NewVScrollContainer(label)

	// Create goroutine to load file contents and update label
	go func() {
		var buffer bytes.Buffer
		count := callback(&buffer)
		log.Println("Count:", count)
		if count > 0 {
			label.SetText(buffer.String())
			scroller.Refresh()
		}
	}()

	return scroller
}

func NewImageGif(callback func(io.Writer) uint64) fyne.CanvasObject {
	// Create image to hold image
	img := &canvas.Image{
		FillMode: canvas.ImageFillContain,
	}
	scroller := widget.NewScrollContainer(img)

	// Create goroutine to load file contents and update image
	go func() {
		var buffer bytes.Buffer
		count := callback(&buffer)
		log.Println("Count:", count)
		if count > 0 {
			i, err := gif.Decode(&buffer)
			if err != nil {
				log.Println("Error:", err)
				return
			}
			img.Image = i
			img.Refresh()
			scroller.Refresh()
		}
	}()

	return scroller
}

func NewImageJpeg(callback func(io.Writer) uint64) fyne.CanvasObject {
	// Create image to hold image
	img := &canvas.Image{
		FillMode: canvas.ImageFillContain,
	}
	scroller := widget.NewScrollContainer(img)

	// Create goroutine to load file contents and update image
	go func() {
		var buffer bytes.Buffer
		count := callback(&buffer)
		log.Println("Count:", count)
		if count > 0 {
			i, err := jpeg.Decode(&buffer)
			if err != nil {
				log.Println("Error:", err)
				return
			}
			img.Image = i
			img.Refresh()
			scroller.Refresh()
		}
	}()

	return scroller
}

func NewImagePng(callback func(io.Writer) uint64) fyne.CanvasObject {
	// Create image to hold image
	img := &canvas.Image{
		FillMode: canvas.ImageFillContain,
	}
	scroller := widget.NewScrollContainer(img)

	// Create goroutine to load file contents and update image
	go func() {
		var buffer bytes.Buffer
		count := callback(&buffer)
		log.Println("Count:", count)
		if count > 0 {
			i, err := png.Decode(&buffer)
			if err != nil {
				log.Println("Error:", err)
				return
			}
			img.Image = i
			img.Refresh()
			scroller.Refresh()
		}
	}()

	return scroller
}
