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
	"aletheiaware.com/bcgo"
	"aletheiaware.com/spaceclientgo"
	"aletheiaware.com/spacego"
	"encoding/base64"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"sort"
)

type MetaList struct {
	widget.List
	ids        []string
	metas      map[string]*spacego.Meta
	timestamps map[string]uint64
}

func NewMetaList(callback func(id string, timestamp uint64, meta *spacego.Meta)) *MetaList {
	l := &MetaList{
		metas:      make(map[string]*spacego.Meta),
		timestamps: make(map[string]uint64),
		List: widget.List{
			CreateItem: func() fyne.CanvasObject {
				return container.NewGridWithColumns(3,
					&widget.Label{
						TextStyle: fyne.TextStyle{
							Bold: true,
						},
						Wrapping: fyne.TextTruncate,
					},
					&widget.Label{
						Alignment: fyne.TextAlignTrailing,
						TextStyle: fyne.TextStyle{
							Monospace: true,
						},
						Wrapping: fyne.TextTruncate,
					},
					&widget.Label{
						Alignment: fyne.TextAlignTrailing,
						TextStyle: fyne.TextStyle{
							Monospace: true,
						},
						Wrapping: fyne.TextTruncate,
					},
				)
			},
		},
	}
	l.Length = func() int {
		return len(l.ids)
	}
	l.UpdateItem = func(id widget.ListItemID, item fyne.CanvasObject) {
		if id < 0 || id >= len(l.ids) {
			return
		}
		i := l.ids[id]
		var name string
		m, ok := l.metas[i]
		if ok {
			name = m.Name
		}
		if name == "" {
			name = "(untitled)"
		}
		items := item.(*fyne.Container).Objects
		items[0].(*widget.Label).SetText(name)
		items[1].(*widget.Label).SetText(m.Type)
		items[2].(*widget.Label).SetText(bcgo.TimestampToString(l.timestamps[i]))
	}
	l.OnSelected = func(id widget.ListItemID) {
		if id < 0 || id >= len(l.ids) {
			return
		}
		i := l.ids[id]
		if m, ok := l.metas[i]; ok && callback != nil {
			callback(i, l.timestamps[i], m)
		}
		l.Unselect(id) // TODO FIXME Hack
	}
	l.ExtendBaseWidget(l)
	return l
}

func (l *MetaList) Add(entry *bcgo.BlockEntry, meta *spacego.Meta) error {
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

func (l *MetaList) Clear() {
	for k := range l.metas {
		delete(l.metas, k)
	}
	l.ids = nil
	l.Refresh()
}

func (l *MetaList) Update(client spaceclientgo.SpaceClient, node bcgo.Node) error {
	if err := client.AllMetas(node, l.Add); err != nil {
		return err
	}
	l.Refresh()
	return nil
}
