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
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"io"
	"io/ioutil"
	"log"
)

type TextPlainViewer struct {
	widget.BaseWidget
	label    *widget.Label
	scroller *container.Scroll
}

func NewTextPlainViewer() *TextPlainViewer {
	v := &TextPlainViewer{}
	// Create label to hold text
	v.label = &widget.Label{
		Wrapping: fyne.TextWrapWord,
	}
	v.scroller = container.NewVScroll(v.label)
	v.ExtendBaseWidget(v)
	return v
}

func (v *TextPlainViewer) CreateRenderer() fyne.WidgetRenderer {
	return v.scroller.CreateRenderer()
}

func (v *TextPlainViewer) SetSource(source io.Reader) {
	bytes, err := ioutil.ReadAll(source)
	if err != nil {
		log.Println("Error:", err)
		return
	}
	v.label.SetText(string(bytes))
	v.scroller.Refresh()
}
