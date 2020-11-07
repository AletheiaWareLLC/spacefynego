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
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"github.com/AletheiaWareLLC/bcgo"
	"github.com/AletheiaWareLLC/financego"
	"github.com/AletheiaWareLLC/spaceclientgo"
	"github.com/AletheiaWareLLC/spacego"
)

type ProviderList struct {
	widget.List
	ids           []string
	registrars    map[string]*spacego.Registrar
	registrations map[string]*financego.Registration
	subscriptions map[string]*financego.Subscription
}

func NewProviderList(callback func(id string, registrar *spacego.Registrar, registration *financego.Registration, subscription *financego.Subscription)) *ProviderList {
	l := &ProviderList{
		registrars:    make(map[string]*spacego.Registrar),
		registrations: make(map[string]*financego.Registration),
		subscriptions: make(map[string]*financego.Subscription),
		List: widget.List{
			CreateItem: func() fyne.CanvasObject {
				return widget.NewLabel("Template Object")
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
		var merchant string
		if r, ok := l.registrars[i]; ok {
			merchant = r.Merchant.Alias
		}
		item.(*widget.Label).SetText(merchant)
	}
	l.OnSelected = func(id widget.ListItemID) {
		if id < 0 || id >= len(l.ids) {
			return
		}
		i := l.ids[id]
		if r, ok := l.registrars[i]; ok && callback != nil {
			callback(i, r, l.registrations[i], l.subscriptions[i])
		}
	}
	l.ExtendBaseWidget(l)
	return l
}

func (l *ProviderList) Add(registrar *spacego.Registrar, registration *financego.Registration, subscription *financego.Subscription) error {
	id := registrar.Merchant.Alias
	if _, ok := l.registrars[id]; !ok {
		l.registrars[id] = registrar
		l.registrations[id] = registration
		l.subscriptions[id] = subscription
		l.ids = append(l.ids, id)
	}
	return nil
}

func (l *ProviderList) Update(client *spaceclientgo.SpaceClient, node *bcgo.Node) error {
	if err := client.GetRegistrarsForNode(node, l.Add); err != nil {
		return err
	}
	l.Refresh()
	return nil
}
