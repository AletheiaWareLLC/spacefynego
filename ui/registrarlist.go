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
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"sort"
)

type RegistrarList struct {
	widget.List
	ids           []string
	registrars    map[string]*spacego.Registrar
	timestamps    map[string]uint64
	registrations map[string]*financego.Registration
	subscriptions map[string]*financego.Subscription
}

func NewRegistrarList(callback func(id string, timestamp uint64, registrar *spacego.Registrar, registration *financego.Registration, subscription *financego.Subscription)) *RegistrarList {
	l := &RegistrarList{
		registrars:    make(map[string]*spacego.Registrar),
		timestamps:    make(map[string]uint64),
		registrations: make(map[string]*financego.Registration),
		subscriptions: make(map[string]*financego.Subscription),
		List: widget.List{
			CreateItem: func() fyne.CanvasObject {
				return container.NewVBox(
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
		r, ok := l.registrars[i]
		if ok {
			box := item.(*fyne.Container)
			box.Objects[0].(*widget.Label).SetText(r.Merchant.Alias)
			items := box.Objects[1].(*fyne.Container).Objects
			if _, ok := l.registrations[i]; ok {
				items[0].(*widget.Label).SetText("Registered")
			} else {
				items[0].(*widget.Label).SetText(r.Service.Country)
			}
			if _, ok := l.subscriptions[i]; ok {
				items[1].(*widget.Label).SetText("Subscribed")
			} else {
				items[1].(*widget.Label).SetText(fmt.Sprintf("%s / %s / %s",
					bcgo.MoneyToString(r.Service.Currency, r.Service.GroupPrice),
					bcgo.DecimalSizeToString(uint64(r.Service.GroupSize)),
					financego.IntervalToString(r.Service.Interval)))
			}
		}
	}
	l.OnSelected = func(id widget.ListItemID) {
		if id < 0 || id >= len(l.ids) {
			return
		}
		i := l.ids[id]
		if r, ok := l.registrars[i]; ok && callback != nil {
			callback(i, l.timestamps[i], r, l.registrations[i], l.subscriptions[i])
		}
		l.Unselect(id) // TODO FIXME Hack
	}
	l.ExtendBaseWidget(l)
	return l
}

func (l *RegistrarList) AddRegistrar(entry *bcgo.BlockEntry, registrar *spacego.Registrar) error {
	id := registrar.Merchant.Alias
	if _, ok := l.registrars[id]; !ok {
		l.registrars[id] = registrar
		l.timestamps[id] = entry.Record.Timestamp
		l.ids = append(l.ids, id)
	}
	return nil
}

func (l *RegistrarList) AddRegistration(entry *bcgo.BlockEntry, registration *financego.Registration) error {
	l.registrations[registration.MerchantAlias] = registration
	return nil
}

func (l *RegistrarList) AddSubscription(entry *bcgo.BlockEntry, subscription *financego.Subscription) error {
	l.subscriptions[subscription.MerchantAlias] = subscription
	return nil
}

func (l *RegistrarList) Update(client spaceclientgo.SpaceClient, node bcgo.Node) error {
	if err := spacego.AllRegistrars(node, l.AddRegistrar); err != nil {
		return err
	}
	if err := spacego.AllRegistrationsForNode(node, l.AddRegistration); err != nil {
		return err
	}
	if err := spacego.AllSubscriptionsForNode(node, l.AddSubscription); err != nil {
		return err
	}
	sort.Slice(l.ids, func(i, j int) bool {
		// TODO sort registrations and subscriptions first
		return l.timestamps[l.ids[i]] < l.timestamps[l.ids[j]]
	})
	l.Refresh()
	return nil
}
