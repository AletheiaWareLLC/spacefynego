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

package spacefynego

import (
	"encoding/base64"
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"github.com/AletheiaWareLLC/bcfynego"
	"github.com/AletheiaWareLLC/bcgo"
	"github.com/AletheiaWareLLC/spaceclientgo"
	"github.com/AletheiaWareLLC/spacego"
	"log"
)

type SpaceFyne struct {
	bcfynego.BCFyne
}

func (f *SpaceFyne) ShowFile(client *spaceclientgo.SpaceClient, id string, meta *spacego.Meta) {
	// TODO
}

func (f *SpaceFyne) FileList(client *spaceclientgo.SpaceClient) (files []fyne.CanvasObject) {
	node, err := client.GetNode()
	if err != nil {
		f.ShowError(err)
		return
	}
	name := spacego.GetMetaChannelName(node.Alias)
	channel := node.GetOrOpenChannel(name, func() *bcgo.Channel {
		return spacego.OpenMetaChannel(node.Alias)
	})
	if err := channel.Refresh(node.Cache, node.Network); err != nil {
		log.Println(err)
	}
	// Iterate Metas and populate list
	if err := bcgo.Iterate(channel.Name, channel.Head, nil, node.Cache, node.Network, func(hash []byte, block *bcgo.Block) error {
		for _, entry := range block.Entry {
			id := base64.RawURLEncoding.EncodeToString(hash)
			meta, err := spacego.UnmarshalMeta(entry.Record.Payload)
			if err != nil {
				return err
			}
			text := meta.Name
			if text == "" {
				text = id[:8]
			}
			files = append(files, widget.NewHBox(
				//widget.NewLabel(id),
				widget.NewButton(text, func() {
					go f.ShowFile(client, id, meta)
				}),
			))
		}
		return nil
	}); err != nil {
		f.ShowError(err)
		return
	}
	return
}
