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
	"aletheiaware.com/spacego"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
)

func init() {
	generator := func() (Viewer, error) {
		return NewImageViewer(), nil
	}
	Register(spacego.MIME_TYPE_IMAGE_GIF, generator)
	Register(spacego.MIME_TYPE_IMAGE_JPEG, generator)
	Register(spacego.MIME_TYPE_IMAGE_JPG, generator)
	// TODO	Register(spacego.MIME_TYPE_IMAGE_SVG, generator)
	Register(spacego.MIME_TYPE_IMAGE_PNG, generator)
}

type ImageViewer struct {
	widget.BaseWidget
	img image.Image
}

func NewImageViewer() *ImageViewer {
	v := &ImageViewer{}
	v.ExtendBaseWidget(v)
	return v
}

func (v *ImageViewer) CreateRenderer() fyne.WidgetRenderer {
	v.ExtendBaseWidget(v)
	r := &imageViewerRenderer{
		viewer: v,
		img: &canvas.Image{
			FillMode: canvas.ImageFillOriginal,
		},
	}
	r.scroller = container.NewScroll(r.img)
	r.objects = []fyne.CanvasObject{r.scroller}
	return r
}

func (v *ImageViewer) MinSize() fyne.Size {
	v.ExtendBaseWidget(v)
	return v.BaseWidget.MinSize()
}

func (v *ImageViewer) SetSource(source io.Reader) error {
	i, _, err := image.Decode(source)
	if err != nil {
		return err
	}
	v.img = i
	v.Refresh()
	return nil
}

type imageViewerRenderer struct {
	viewer   *ImageViewer
	img      *canvas.Image
	scroller *container.Scroll
	objects  []fyne.CanvasObject
}

func (r *imageViewerRenderer) Destroy() {}

func (r *imageViewerRenderer) Layout(size fyne.Size) {
	r.scroller.Resize(size)
}

func (r *imageViewerRenderer) MinSize() fyne.Size {
	return r.scroller.MinSize()
}

func (r *imageViewerRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *imageViewerRenderer) Refresh() {
	r.img.Image = r.viewer.img
	r.img.Refresh()
	r.scroller.Refresh()
}
