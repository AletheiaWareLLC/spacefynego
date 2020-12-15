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
	"fyne.io/fyne/dialog"
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
	"os"
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

	// Set environment variable
	if a.Settings().BuildType() == fyne.BuildRelease {
		os.Setenv("LIVE", "true")
	}

	// Create Space Client
	c := spaceclientgo.NewSpaceClient(bcgo.SplitRemoveEmpty(*peer, ",")...)

	// Create Space Fyne
	f := spacefynego.NewSpaceFyne(a, w, c)

	// Create a scrollable list of metas
	l := ui.NewMetaList(func(id string, timestamp uint64, meta *spacego.Meta) {
		go f.ShowFile(c, id, timestamp, meta)
	})

	refreshList := func() {
		n, err := f.GetNode(&c.BCClient)
		if err != nil {
			f.ShowError(err)
			return
		}

		// Show progress dialog
		progress := dialog.NewProgressInfinite("Refreshing", "Refreshing File List", f.Window)
		progress.Show()
		defer progress.Hide()

		l.Update(c, n)
	}

	// Populate list in goroutine
	go refreshList()

	f.OnKeysImported = func(alias string) {
		go refreshList()
	}
	onSignedIn := f.OnSignedIn
	f.OnSignedIn = func(node *bcgo.Node) {
		if onSignedIn != nil {
			onSignedIn(node)
		}
		go l.Update(c, node)
	}
	f.OnSignedOut = func() {
		go refreshList()
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
		widget.NewToolbarAction(theme.NewThemedResource(data.StorageIcon, nil), func() {
			go f.ShowStorage(c)
		}),
		widget.NewToolbarAction(bcuidata.NewPrimaryThemedResource(bcuidata.AccountIcon), func() {
			go f.ShowAccount(&c.BCClient)
		}),
		widget.NewToolbarAction(theme.HelpIcon(), func() {
			go f.ShowHelp(c)
		}),
	)

	// Set window content, resize window, center window, show window, and run application
	w.SetContent(container.NewBorder(t, nil, nil, nil, f.GetIcon(), l))
	w.Resize(fyne.NewSize(600, 480))
	w.CenterOnScreen()
	w.ShowAndRun()
}
