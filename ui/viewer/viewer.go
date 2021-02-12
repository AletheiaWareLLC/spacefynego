/*
 * Copyright 2020-2021 Aletheia Ware LLC
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
	"io"
)

type Viewer interface {
	fyne.CanvasObject
	SetSource(io.Reader)
}

func ForMime(mime string) Viewer {
	switch mime {
	case spacego.MIME_TYPE_TEXT_PLAIN:
		return NewTextPlainViewer()
	case spacego.MIME_TYPE_IMAGE_GIF:
		fallthrough
	case spacego.MIME_TYPE_IMAGE_JPEG:
		fallthrough
		// TODO	case spacego.MIME_TYPE_IMAGE_SVG:
		// TODO		fallthrough
	case spacego.MIME_TYPE_IMAGE_PNG:
		return NewImageViewer()
	}
	return nil
}
