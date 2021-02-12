/*
 * Copyright 2021 Aletheia Ware LLC
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
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
)

type ImageViewer struct {
	widget.BaseWidget
	img      *canvas.Image
	scroller *container.Scroll
}

func NewImageViewer() *ImageViewer {
	v := &ImageViewer{}
	// Create image to hold image
	v.img = &canvas.Image{
		FillMode: canvas.ImageFillOriginal,
	}
	v.scroller = container.NewScroll(v.img)
	v.ExtendBaseWidget(v)
	return v
}

func (v *ImageViewer) CreateRenderer() fyne.WidgetRenderer {
	return v.scroller.CreateRenderer()
}

func (v *ImageViewer) SetSource(source io.Reader) {
	i, _, err := image.Decode(source)
	if err != nil {
		log.Println("Error:", err)
		return
	}
	v.img.Image = i
	v.img.Refresh()
	v.scroller.Refresh()
}
