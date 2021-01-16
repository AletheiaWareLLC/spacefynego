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
	"aletheiaware.com/bcfynego"
	bcui "aletheiaware.com/bcfynego/ui"
	"aletheiaware.com/bcgo"
	"aletheiaware.com/financego"
	"aletheiaware.com/spaceclientgo"
	"aletheiaware.com/spacefynego/ui"
	"aletheiaware.com/spacefynego/ui/data"
	"aletheiaware.com/spacefynego/ui/viewer"
	"aletheiaware.com/spacego"
	"encoding/base64"
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/storage"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"log"
	"net/url"
)

type SpaceFyne struct {
	bcfynego.BCFyne
}

func NewSpaceFyne(a fyne.App, w fyne.Window, c *spaceclientgo.SpaceClient) *SpaceFyne {
	f := &SpaceFyne{
		BCFyne: *bcfynego.NewBCFyne(a, w),
	}
	f.OnKeysImported = func(alias string) {
		// TODO show success dialog, and tell the user to sign in with their newly-imported alias, and password
	}
	f.OnSignedIn = func(node *bcgo.Node) {
		count := 0
		if err := c.GetRegistrarsForNode(node, func(*spacego.Registrar, *financego.Registration, *financego.Subscription) error {
			count++
			return nil
		}); err != nil {
			log.Println(err)
		}
		if count < spacego.GetMinimumRegistrars() {
			f.ShowRegistrarSelectionDialog(c, node)
		}
	}
	f.OnSignedUp = func(node *bcgo.Node) {
		f.ShowWelcome(c, node)
	}
	return f
}

// ShowWelcome displays a wizard to welcome a new user and walk them through the setup process.
func (f *SpaceFyne) ShowWelcome(client *spaceclientgo.SpaceClient, node *bcgo.Node) {
	if d := f.Dialog; d != nil {
		d.Hide()
	}
	f.Dialog = dialog.NewCustomConfirm("Welcome", "Next", "Cancel",
		widget.NewLabel(fmt.Sprintf("Hello %s, Welcome to S P A C E!", node.Alias)),
		func(result bool) {
			if result {
				f.ShowRegistrarSelectionDialog(client, node)
			}
		},
		f.Window)
	f.Dialog.Show()
}

func (f *SpaceFyne) ShowRegistrarSelectionDialog(client *spaceclientgo.SpaceClient, node *bcgo.Node) {
	if d := f.Dialog; d != nil {
		d.Hide()
	}

	// Show progress dialog
	f.Dialog = dialog.NewProgressInfinite("Updating", "Getting Registrars", f.Window)
	f.Dialog.Show()

	l := ui.NewRegistrarList(func(id string, timestamp uint64, registrar *spacego.Registrar) {
		u, err := url.Parse(fmt.Sprintf("https://%s/%s", registrar.Merchant.Domain, registrar.Merchant.RegisterUrl))
		if err != nil {
			f.ShowError(err)
			return
		}
		params := url.Values{}
		params.Add("alias", node.Alias)
		params.Add("next", registrar.Service.SubscribeUrl)
		u.RawQuery = params.Encode()
		if err := f.App.OpenURL(u); err != nil {
			f.ShowError(err)
			return
		}
	})

	// Update list
	l.Update(client, node)

	if d := f.Dialog; d != nil {
		d.Hide()
	}
	f.Dialog = dialog.NewCustom("Registrars", "Done",
		container.NewBorder(
			&widget.Label{
				Text:     "Your encrypted data will be stored by your choice of storage providers.",
				Wrapping: fyne.TextWrapWord,
			},
			&widget.Label{
				Text:     fmt.Sprintf("Choose at least %d providers from the list above.", spacego.GetMinimumRegistrars()),
				Wrapping: fyne.TextWrapWord,
			},
			nil,
			nil,
			l,
		),
		f.Window)
	f.Dialog.Resize(fyne.NewSize(400, 300))
	f.Dialog.Show()
}

func (f *SpaceFyne) GetIcon() fyne.CanvasObject {
	icon := canvas.NewImageFromResource(data.SpaceIcon)
	icon.FillMode = canvas.ImageFillContain
	return icon
}

// Add displays a dialog (write text, take a picture, upload an existing file or folder), and adds the result.
func (f *SpaceFyne) Add(client *spaceclientgo.SpaceClient) {
	node, err := f.GetNode(&client.BCClient)
	if err != nil {
		f.ShowError(err)
		return
	}

	if d := f.Dialog; d != nil {
		d.Hide()
	}
	// Show progress dialog
	f.Dialog = dialog.NewProgressInfinite("Updating", "Getting Registrars", f.Window)
	f.Dialog.Show()

	domains, err := f.getRegistrarDomainsForNode(client, node)
	if err != nil {
		log.Println(err)
	}

	if len(domains) == 0 {
		f.ShowRegistrarSelectionDialog(client, node)
		return
	}

	if d := f.Dialog; d != nil {
		d.Hide()
	}
	composeText := widget.NewButtonWithIcon("Text", theme.DocumentCreateIcon(), func() {
		if d := f.Dialog; d != nil {
			d.Hide()
		}
		f.ShowAddTextDialog(client, node)
	})
	captureImage := widget.NewButtonWithIcon("Image", theme.NewThemedResource(data.CameraPhotoIcon, nil), func() {
		if d := f.Dialog; d != nil {
			d.Hide()
		}
		// TODO
		f.ShowError(fmt.Errorf("Not yet implemented: %s %s", "SpaceFyne.Add", "Image"))
	})
	captureVideo := widget.NewButtonWithIcon("Video", theme.NewThemedResource(data.CameraVideoIcon, nil), func() {
		if d := f.Dialog; d != nil {
			d.Hide()
		}
		// TODO
		f.ShowError(fmt.Errorf("Not yet implemented: %s %s", "SpaceFyne.Add", "Video"))
	})
	captureAudio := widget.NewButtonWithIcon("Audio", theme.NewThemedResource(data.MicrophoneIcon, nil), func() {
		if d := f.Dialog; d != nil {
			d.Hide()
		}
		// TODO
		f.ShowError(fmt.Errorf("Not yet implemented: %s %s", "SpaceFyne.Add", "Audio"))
	})
	uploadFile := widget.NewButtonWithIcon("Document", theme.FileIcon(), func() {
		if d := f.Dialog; d != nil {
			d.Hide()
		}
		go f.ShowUploadFileDialog(client, node)
	})
	uploadFolder := widget.NewButtonWithIcon("Directory", theme.FolderIcon(), func() {
		if d := f.Dialog; d != nil {
			d.Hide()
		}
		go f.ShowUploadFolderDialog(client, node)
	})
	content := container.NewAdaptiveGrid(3,
		composeText,
		captureImage,
		captureAudio,
		captureVideo,
		uploadFile,
		uploadFolder,
	)
	if d := f.Dialog; d != nil {
		d.Hide()
	}
	f.Dialog = dialog.NewCustom("Add", "Cancel", content, f.Window)
	f.Dialog.Show()
}

func (f *SpaceFyne) SearchFile(client *spaceclientgo.SpaceClient) {
	f.ShowError(fmt.Errorf("Not yet implemented: %s", "SpaceFyne.SearchFile"))
}

func (f *SpaceFyne) ShowFile(client *spaceclientgo.SpaceClient, id string, timestamp uint64, meta *spacego.Meta) {
	node, err := f.GetNode(&client.BCClient)
	if err != nil {
		f.ShowError(err)
		return
	}

	// Show progress dialog
	progress := dialog.NewProgressInfinite("Loading", "Reading "+meta.Name, f.Window)
	progress.Show()
	defer progress.Hide()

	hash, err := base64.RawURLEncoding.DecodeString(id)
	if err != nil {
		f.ShowError(err)
		return
	}

	reader, err := client.ReadFile(node, hash)
	if err != nil {
		f.ShowError(err)
		return
	}

	view := viewer.GetViewer(meta, reader)
	if view == nil {
		f.ShowError(fmt.Errorf("Not yet implemented: %s %s", "SpaceFyne.ShowFile", meta.Type))
		return
	}

	name := meta.Name
	if name == "" {
		name = "(untitled)"
	}
	window := f.App.NewWindow(fmt.Sprintf("%s - %s - %s", bcgo.TimestampToString(timestamp), name, id[:8]))
	window.SetContent(view)
	window.Resize(fyne.NewSize(800, 600))
	window.CenterOnScreen()
	window.Show()
}

func (f *SpaceFyne) ShowStorage(client *spaceclientgo.SpaceClient) {
	node, err := f.GetNode(&client.BCClient)
	if err != nil {
		f.ShowError(err)
		return
	}

	// Show progress dialog
	progress := dialog.NewProgressInfinite("Updating", "Getting Registrars", f.Window)
	progress.Show()

	list := ui.NewProviderList(f.ShowRegistrarDialog)

	// Update list
	list.Update(client, node)

	// Hide progress dialog
	progress.Hide()

	if d := f.Dialog; d != nil {
		d.Hide()
	}
	// Show registrar list
	f.Dialog = dialog.NewCustom("Registrars", "OK", list, f.Window)
	f.Dialog.Resize(fyne.NewSize(320, 320))
	f.Dialog.Show()
}

func (f *SpaceFyne) ShowRegistrarDialog(id string, registrar *spacego.Registrar, registration *financego.Registration, subscription *financego.Subscription) {
	// Show detailed information
	info := dialog.NewCustom(registrar.Merchant.Alias, "OK", widget.NewForm(
		widget.NewFormItem("Domain", &widget.Label{
			Text:     registrar.Merchant.Domain,
			Wrapping: fyne.TextWrapBreak,
		}),
		widget.NewFormItem("Country", &widget.Label{
			Text: registrar.Service.Country,
		}),
		widget.NewFormItem("Cost", &widget.Label{
			Text: fmt.Sprintf("%s / %s / %s",
				bcgo.MoneyToString(registrar.Service.Currency, registrar.Service.GroupPrice),
				bcgo.DecimalSizeToString(uint64(registrar.Service.GroupSize)),
				financego.IntervalToString(registrar.Service.Interval)),
		}),
		widget.NewFormItem("Customer", &widget.Label{
			Text: registration.CustomerId,
		}),
		widget.NewFormItem("Subscription", &widget.Label{
			Text: subscription.SubscriptionId,
		}),
		widget.NewFormItem("Subscription Item", &widget.Label{
			Text: subscription.SubscriptionItemId,
		}),
		/* TODO
		- Payment Methods
		- Usage Records
		- Invoices & Reciepts
		*/
	), f.Window)
	info.Resize(fyne.NewSize(320, 320))
	info.Show()
}

func (f *SpaceFyne) ShowHelp(client *spaceclientgo.SpaceClient) {
	f.ShowError(fmt.Errorf("Not yet implemented: %s", "SpaceFyne.ShowHelp"))
}

// ShowAddNoteDialog displays a dialog for creating a note, and adds the resulting file.
func (f *SpaceFyne) ShowAddTextDialog(client *spaceclientgo.SpaceClient, node *bcgo.Node) {
	// TODO
}

// ShowUploadFileDialog displays a file picker, and adds the resulting file.
func (f *SpaceFyne) ShowUploadFileDialog(client *spaceclientgo.SpaceClient, node *bcgo.Node) {
	f.Dialog = dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			f.ShowError(err)
			return
		}
		if reader == nil {
			return
		}
		f.UploadFile(client, node, reader)
	}, f.Window)
	f.Dialog.Show()
}

// ShowUploadFolderDialog displays a folder picker, and adds the resulting folder.
func (f *SpaceFyne) ShowUploadFolderDialog(client *spaceclientgo.SpaceClient, node *bcgo.Node) {
	f.Dialog = dialog.NewFolderOpen(func(lister fyne.ListableURI, err error) {
		if err != nil {
			f.ShowError(err)
			return
		}
		if lister == nil {
			return
		}
		f.UploadFolder(client, node, lister)
	}, f.Window)
	f.Dialog.Show()
}

func (f *SpaceFyne) UploadFile(client *spaceclientgo.SpaceClient, node *bcgo.Node, file fyne.URIReadCloser) {
	// Show progress dialog
	progress := dialog.NewProgress("Uploading", "Uploading "+file.Name(), f.Window)
	progress.Show()
	defer progress.Hide()
	listener := &bcui.ProgressMiningListener{Func: progress.SetValue}

	reference, err := client.Add(node, listener, file.Name(), file.URI().MimeType(), file)
	if err != nil {
		f.ShowError(err)
	}
	log.Println("Uploaded:", reference)
}

func (f *SpaceFyne) UploadFolder(client *spaceclientgo.SpaceClient, node *bcgo.Node, folder fyne.ListableURI) {
	uris, err := folder.List()
	if err != nil {
		f.ShowError(err)
		return
	}
	count := len(uris)

	// Show progress dialog
	progress := dialog.NewProgress("Uploading", "Uploading "+folder.Name(), f.Window)
	progress.Show()
	defer progress.Hide()

	for i, uri := range uris {
		progress.SetValue(float64(i) / float64(count))

		// Check if URI points to Folder
		if lister, ok := uri.(fyne.ListableURI); ok {
			f.UploadFolder(client, node, lister)
			continue
		}
		// Check if URI points to Folder
		if lister, err := storage.ListerForURI(uri); err == nil {
			f.UploadFolder(client, node, lister)
			continue
		}

		// URI points to File
		file, err := storage.OpenFileFromURI(uri)
		if err != nil {
			log.Println(err)
			continue
		}
		f.UploadFile(client, node, file)
	}
}

func (f *SpaceFyne) getRegistrarDomainsForNode(client *spaceclientgo.SpaceClient, node *bcgo.Node) (domains []string, err error) {
	var net *bcgo.TCPNetwork
	if node.Network != nil {
		if n, ok := node.Network.(*bcgo.TCPNetwork); ok {
			net = n
		}
	}
	err = client.GetRegistrarsForNode(node, func(registrar *spacego.Registrar, registration *financego.Registration, subscription *financego.Subscription) error {
		if registrar != nil && registration != nil && subscription != nil {
			domain := registrar.Merchant.Domain
			if domain == "" {
				domain = registrar.Merchant.Alias
			}
			// Add any missing domains to network
			if net != nil {
				if _, ok := net.Peers[domain]; !ok {
					net.Peers[domain] = 0
				}
			}
			domains = append(domains, domain)
		}
		return nil
	})
	return
}
