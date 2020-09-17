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

package ui

import (
	"encoding/base64"
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"github.com/AletheiaWareLLC/bcgo"
	"github.com/AletheiaWareLLC/spacego"
	"sort"
)

type MetaList struct {
	widget.List
	ids        []string
	metas      map[string]*spacego.Meta
	timestamps map[string]uint64
}

func NewMetaList(callback func(id string, meta *spacego.Meta)) *MetaList {
	l := &MetaList{
		metas:      make(map[string]*spacego.Meta),
		timestamps: make(map[string]uint64),
		List: widget.List{
			CreateItem: func() fyne.CanvasObject {
				return &widget.Label{}
			},
		},
	}
	l.Length = func() int {
		return len(l.ids)
	}
	l.UpdateItem = func(index int, item fyne.CanvasObject) {
		var name string
		m, ok := l.metas[l.ids[index]]
		if ok {
			name = m.Name
		}
		if name == "" {
			name = "(untitled)"
		}
		item.(*widget.Label).SetText(name)
	}
	l.OnItemSelected = func(index int) {
		callback(l.ids[index], l.metas[l.ids[index]])
	}
	l.ExtendBaseWidget(l)
	return l
}

func (l *MetaList) Update(entry *bcgo.BlockEntry, meta *spacego.Meta) error {
	id := base64.RawURLEncoding.EncodeToString(entry.RecordHash)
	if _, ok := l.metas[id]; !ok {
		l.metas[id] = meta
		l.timestamps[id] = entry.Record.Timestamp
		l.ids = append(l.ids, id)
		sort.Slice(l.ids, func(i, j int) bool {
			return l.timestamps[l.ids[i]] < l.timestamps[l.ids[j]]
		})
	}
	return nil
}
