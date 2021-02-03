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

package main

import (
	bcui "aletheiaware.com/bcfynego/ui"
	bcuidata "aletheiaware.com/bcfynego/ui/data"
	"aletheiaware.com/bcgo"
	"aletheiaware.com/spaceclientgo"
	"aletheiaware.com/spacefynego"
	"aletheiaware.com/spacefynego/ui"
	"aletheiaware.com/spacefynego/ui/data"
	"aletheiaware.com/spacego"
	"flag"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"log"
)

var peer = flag.String("peer", "", "Space peer")

func main() {
	// Parse command line flags
	flag.Parse()

	// Set log flags
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Create Application
	a := app.NewWithID("com.aletheiaware.space")

	// Create Window
	w := a.NewWindow("S P A C E")

	// Create Space Client
	c := spaceclientgo.NewSpaceClient(bcgo.SplitRemoveEmpty(*peer, ",")...)

	// Create Space Fyne
	f := spacefynego.NewSpaceFyne(a, w, c)

	// Create a scrollable list of metas
	l := ui.NewMetaList(func(id string, timestamp uint64, meta *spacego.Meta) {
		go f.ShowFile(c, id, timestamp, meta)
	})

	refreshListWithNode := func(n *bcgo.Node) {
		// Show progress dialog
		progress := dialog.NewProgressInfinite("Refreshing", "Refreshing File List", f.Window)
		progress.Show()
		defer progress.Hide()

		l.Update(c, n)
	}

	refreshList := func() {
		n, err := f.GetNode(&c.BCClient)
		if err != nil {
			f.ShowError(err)
			return
		}
		refreshListWithNode(n)
	}

	// Populate list in goroutine
	go refreshList()

	onSignedIn := f.OnSignedIn
	f.OnSignedIn = func(node *bcgo.Node) {
		if onSignedIn != nil {
			onSignedIn(node)
		}
		go refreshListWithNode(node)
	}
	f.OnSignedOut = func() {
		go l.Clear()
	}

	// Create a toolbar of common operations
	t := widget.NewToolbar(
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {
			go f.Add(c)
		}),
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
			go refreshList()
		}),
		widget.NewToolbarAction(theme.SearchIcon(), func() {
			go f.SearchFile(c)
		}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.NewThemedResource(data.StorageIcon), func() {
			go f.ShowStorage(c)
		}),
		widget.NewToolbarAction(theme.NewThemedResource(bcuidata.AccountIcon), func() {
			go f.ShowAccount(&c.BCClient)
		}),
		widget.NewToolbarAction(theme.HelpIcon(), func() {
			go f.ShowHelp(c)
		}),
	)

	// Set window content, resize window, center window, show window, and run application
	w.SetContent(container.NewBorder(t, nil, nil, nil, l))
	w.Resize(bcui.WindowSize)
	w.CenterOnScreen()
	w.ShowAndRun()
}
