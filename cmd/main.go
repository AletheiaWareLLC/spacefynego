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

package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/AletheiaWareLLC/bcclientgo"
	"github.com/AletheiaWareLLC/bcfynego"
	bcuidata "github.com/AletheiaWareLLC/bcfynego/ui/data"
	"github.com/AletheiaWareLLC/bcgo"
	"github.com/AletheiaWareLLC/spaceclientgo"
	"github.com/AletheiaWareLLC/spacefynego"
	"github.com/AletheiaWareLLC/spacefynego/ui/data"
	"github.com/AletheiaWareLLC/spacego"
	"log"
)

func main() {
	// Create Application
	a := app.New()

	// Create Window
	w := a.NewWindow("S P A C E")
	w.SetMaster()

	peers := append(
		spacego.GetSpaceHosts(), // Add SPACE host as peer
		bcgo.GetBCHost(),        // Add BC host as peer
	)

	// Create Space Client
	c := &spaceclientgo.SpaceClient{
		BCClient: bcclientgo.BCClient{
			Peers: peers,
		},
	}

	// Create Space Fyne
	f := &spacefynego.SpaceFyne{
		BCFyne: bcfynego.BCFyne{
			App:    a,
			Window: w,
		},
	}

	space := canvas.NewImageFromResource(data.SpaceIcon)
	space.FillMode = canvas.ImageFillContain

	// Create a scrollable list of files
	fileList := f.NewList(c)

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
			log.Println("TODO Create File")
			go f.NewFile(c)
		}),
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {
			log.Println("Upload File")
			go f.UploadFile(c)
		}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.SearchIcon(), func() {
			log.Println("TODO Search File")
			// TODO go f.SearchFile(c)
		}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
			log.Println("Refresh List")
			go func() {
				node, err := f.GetNode(&c.BCClient)
				if err != nil {
					f.ShowError(err)
					return
				}
				if err := c.List(node, fileList.Update); err != nil {
					f.ShowError(err)
					return
				}
				fileList.Refresh()
			}()
		}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(bcuidata.NewPrimaryThemedResource(bcuidata.AccountIcon), func() {
			log.Println("Account Info")
			go f.ShowNode(&c.BCClient)
		}),
		widget.NewToolbarAction(theme.HelpIcon(), func() {
			log.Println("TODO Display Help")
			//go c.ShowHelp()
		}),
	)

	// Set window content, resize window, center window, show window, and run application
	w.SetContent(fyne.NewContainerWithLayout(layout.NewBorderLayout(toolbar, nil, nil, nil), space, toolbar, fileList))
	w.Resize(fyne.NewSize(800, 600))
	w.CenterOnScreen()
	w.ShowAndRun()
}
