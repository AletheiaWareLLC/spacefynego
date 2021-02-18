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
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"io"
	"io/ioutil"
)

func init() {
	Register(spacego.MIME_TYPE_TEXT_PLAIN, func() (Viewer, error) {
		return NewTextPlainViewer(), nil
	})
}

type TextPlainViewer struct {
	widget.BaseWidget
	text string
}

func NewTextPlainViewer() *TextPlainViewer {
	v := &TextPlainViewer{}
	v.ExtendBaseWidget(v)
	return v
}

func (v *TextPlainViewer) CreateRenderer() fyne.WidgetRenderer {
	v.ExtendBaseWidget(v)
	r := &textPlainViewerRenderer{
		viewer: v,
		label: &widget.Label{
			Wrapping: fyne.TextWrapWord,
		},
	}
	r.scroller = container.NewVScroll(r.label)
	r.objects = []fyne.CanvasObject{r.scroller}
	return r
}

func (v *TextPlainViewer) MinSize() fyne.Size {
	v.ExtendBaseWidget(v)
	return v.BaseWidget.MinSize()
}

func (v *TextPlainViewer) SetSource(source io.Reader) error {
	bytes, err := ioutil.ReadAll(source)
	if err != nil {
		return err
	}
	v.text = string(bytes)
	v.Refresh()
	return nil
}

type textPlainViewerRenderer struct {
	viewer   *TextPlainViewer
	label    *widget.Label
	scroller *container.Scroll
	objects  []fyne.CanvasObject
}

func (r *textPlainViewerRenderer) Destroy() {}

func (r *textPlainViewerRenderer) Layout(size fyne.Size) {
	r.scroller.Resize(size)
}

func (r *textPlainViewerRenderer) MinSize() fyne.Size {
	return r.scroller.MinSize()
}

func (r *textPlainViewerRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *textPlainViewerRenderer) Refresh() {
	r.label.Text = r.viewer.text
	r.label.Refresh()
	r.scroller.Refresh()
}
