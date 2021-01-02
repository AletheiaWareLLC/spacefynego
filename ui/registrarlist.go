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
	"aletheiaware.com/financego"
	"aletheiaware.com/spaceclientgo"
	"aletheiaware.com/spacego"
	"encoding/base64"
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
	"github.com/golang/protobuf/proto"
	"log"
	"sort"
)

type RegistrarList struct {
	widget.List
	ids        []string
	registrars map[string]*spacego.Registrar
	timestamps map[string]uint64
}

func NewRegistrarList(callback func(id string, timestamp uint64, registrar *spacego.Registrar)) *RegistrarList {
	l := &RegistrarList{
		registrars: make(map[string]*spacego.Registrar),
		timestamps: make(map[string]uint64),
		List: widget.List{
			CreateItem: func() fyne.CanvasObject {
				return widget.NewVBox(
					&widget.Label{
						Alignment: fyne.TextAlignLeading,
						TextStyle: fyne.TextStyle{
							Bold: true,
						},
					},
					container.NewGridWithColumns(2,
						&widget.Label{
							Alignment: fyne.TextAlignLeading,
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
					))
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
		var merchant *financego.Merchant
		var service *financego.Service
		r, ok := l.registrars[i]
		if ok {
			merchant = r.Merchant
			service = r.Service
			box := item.(*widget.Box)
			box.Children[0].(*widget.Label).SetText(merchant.Alias)
			items := box.Children[1].(*fyne.Container).Objects
			items[0].(*widget.Label).SetText(service.Country)
			items[1].(*widget.Label).SetText(fmt.Sprintf("%s / %s / %s",
				bcgo.MoneyToString(service.Currency, service.GroupPrice),
				bcgo.DecimalSizeToString(uint64(service.GroupSize)),
				financego.IntervalToString(service.Interval)))
		}
	}
	l.OnSelected = func(id widget.ListItemID) {
		if id < 0 || id >= len(l.ids) {
			return
		}
		i := l.ids[id]
		if r, ok := l.registrars[i]; ok && callback != nil {
			callback(i, l.timestamps[i], r)
		}
	}
	l.ExtendBaseWidget(l)
	return l
}

func (l *RegistrarList) Add(entry *bcgo.BlockEntry, registrar *spacego.Registrar) error {
	id := base64.RawURLEncoding.EncodeToString(entry.RecordHash)
	if _, ok := l.registrars[id]; !ok {
		l.registrars[id] = registrar
		l.timestamps[id] = entry.Record.Timestamp
		l.ids = append(l.ids, id)
		sort.Slice(l.ids, func(i, j int) bool {
			return l.timestamps[l.ids[i]] < l.timestamps[l.ids[j]]
		})
	}
	return nil
}

func (l *RegistrarList) Update(client *spaceclientgo.SpaceClient, node *bcgo.Node) error {
	registrars := node.GetOrOpenChannel(spacego.SPACE_REGISTRAR, func() *bcgo.Channel {
		return spacego.OpenRegistrarChannel()
	})
	if err := registrars.Refresh(node.Cache, node.Network); err != nil {
		log.Println(err)
	}
	if err := bcgo.Read(registrars.Name, registrars.Head, nil, node.Cache, node.Network, "", nil, nil, func(entry *bcgo.BlockEntry, key, data []byte) error {
		// Unmarshal as Registrar
		r := &spacego.Registrar{}
		err := proto.Unmarshal(data, r)
		if err != nil {
			return err
		}
		l.Add(entry, r)
		return nil
	}); err != nil {
		return err
	}
	l.Refresh()
	return nil
}
