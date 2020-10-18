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
	"flag"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	bcuidata "github.com/AletheiaWareLLC/bcfynego/ui/data"
	"github.com/AletheiaWareLLC/bcgo"
	"github.com/AletheiaWareLLC/spaceclientgo"
	"github.com/AletheiaWareLLC/spacefynego"
	"github.com/AletheiaWareLLC/spacefynego/ui"
	"github.com/AletheiaWareLLC/spacefynego/ui/data"
	"github.com/AletheiaWareLLC/spacego"
	"log"
)

var peer = flag.String("peer", "", "Space peer")

func main() {
	// Parse command line flags
	flag.Parse()

	// Create Application
	a := app.New()

	// Create Window
	w := a.NewWindow("S P A C E")
	w.SetMaster()

	// Create Space Client
	c := spaceclientgo.NewSpaceClient(bcgo.SplitRemoveEmpty(*peer, ",")...)

	// Create Space Fyne
	f := spacefynego.NewSpaceFyne(a, w, c)

	// Create a scrollable list of metas
	l := ui.NewMetaList(func(id string, timestamp uint64, meta *spacego.Meta) {
		f.ShowFile(c, id, timestamp, meta)
	})

	refreshList := func() {
		n, err := f.GetNode(&c.BCClient)
		if err != nil {
			f.ShowError(err)
			return
		}
		l.Update(c, n)
	}

	// Populate list in goroutine
	go refreshList()

	f.OnKeysImported = func(alias string) {
		go refreshList()
	}
	f.OnSignedIn = func(node *bcgo.Node) {
		go l.Update(c, node)
		// TODO FIXME Remove
		go f.ShowWelcome(c, node)
	}

	// Create a toolbar of common operations
	t := widget.NewToolbar(
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {
			log.Println("Add File")
			go f.AddFile(c)
		}),
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
			log.Println("Refresh List")
			go refreshList()
		}),
		widget.NewToolbarAction(theme.SearchIcon(), func() {
			log.Println("Search File")
			go f.SearchFile(c)
		}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.NewThemedResource(data.StorageIcon, nil), func() {
			log.Println("Storage Info")
			go f.ShowStorage(c)
		}),
		widget.NewToolbarAction(bcuidata.NewPrimaryThemedResource(bcuidata.AccountIcon), func() {
			log.Println("Account Info")
			go f.ShowAccount(&c.BCClient)
		}),
		widget.NewToolbarAction(theme.HelpIcon(), func() {
			log.Println("Display Help")
			go f.ShowHelp(c)
		}),
	)

	// Set window content, resize window, center window, show window, and run application
	w.SetContent(container.NewBorder(t, nil, nil, nil, f.GetIcon(), t, l))
	w.Resize(fyne.NewSize(600, 480))
	w.CenterOnScreen()
	w.ShowAndRun()
}
