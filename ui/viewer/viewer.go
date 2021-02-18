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
	"fmt"
	"fyne.io/fyne/v2"
	"io"
	"strings"
)

// generatorTable stores the mapping of mime types to generators of Viewers.
var generatorTable map[string]func() (Viewer, error) = map[string]func() (Viewer, error){}

// Viewer represents a fyne.CanvasObject that can view a file.
type Viewer interface {
	fyne.CanvasObject
	SetSource(io.Reader) error
}

// Register registers a function that can generate a generator.
func Register(mime string, generator func() (Viewer, error)) {
	generatorTable[strings.ToLower(mime)] = generator
}

// ForMime returns the Viewer instance which is registered to handle URIs
// of the given mime.
func ForMime(mime string) (Viewer, error) {
	generator, ok := generatorTable[strings.ToLower(mime)]

	if !ok {
		return nil, fmt.Errorf("no generator registered for mime '%s'", mime)
	}

	return generator()
}
