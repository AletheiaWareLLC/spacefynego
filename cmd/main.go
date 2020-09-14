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
	"github.com/AletheiaWareLLC/spacefynego"
	"github.com/AletheiaWareLLC/spacefynego/ui/data"
	"log"
)

func main() {
	// Create application
	a := app.New()

	// Create window
	w := a.NewWindow("S P A C E")
	w.SetMaster()

	// Create Client
	c := &spacefynego.SpaceClient{
		BCFyneClient: bcfynego.BCFyneClient{
			BCClient: bcclientgo.BCClient{
				Peers: []string{
					// TODO bcgo.GetBCHost(),         // Add BC host as peer
					// TODO spacego.GetSpaceHost(), // Add SPACE host as peer
				},
			},
			App:    a,
			Window: w,
		},
	}

	space := canvas.NewImageFromResource(data.SpaceIcon)
	space.FillMode = canvas.ImageFillContain

	// Create a scrollable list of files
	fileBox := widget.NewVBox()
	refreshList := func() {
		fileBox.Children = c.FileList()
		// Trigger list redraw
		fileBox.Refresh()
	}
	go refreshList()
	fileList := widget.NewScrollContainer(fileBox)

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {
			log.Println("Create File")
			go c.NewFile()
		}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.SearchIcon(), func() {
			log.Println("Search File")
			go c.SearchFile()
		}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
			log.Println("Refresh List")
			go refreshList()
		}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(data.NewPrimaryThemedResource(data.AccountIcon), func() {
			log.Println("Account Info")
			go c.ShowNode()
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
