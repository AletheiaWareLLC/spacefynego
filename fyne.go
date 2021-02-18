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
	bcstorage "aletheiaware.com/bcfynego/storage"
	bcui "aletheiaware.com/bcfynego/ui"
	"aletheiaware.com/bcgo"
	"aletheiaware.com/financego"
	"aletheiaware.com/spaceclientgo"
	"aletheiaware.com/spacefynego/storage"
	"aletheiaware.com/spacefynego/ui"
	"aletheiaware.com/spacefynego/ui/data"
	"aletheiaware.com/spacefynego/ui/viewer"
	"aletheiaware.com/spacego"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	fynestorage "fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"strings"
)

type SpaceFyne struct {
	bcfynego.BCFyne
}

func NewSpaceFyne(a fyne.App, w fyne.Window, c *spaceclientgo.SpaceClient) *SpaceFyne {
	f := &SpaceFyne{
		BCFyne: *bcfynego.NewBCFyne(a, w),
	}
	f.OnSignedIn = func(node *bcgo.Node) {
		// Create BC Repository
		bcstorage.NewBCRepository(&c.BCClient).Register()
		// Create Space Repository
		storage.NewSpaceRepository(c, node).Register()
		count := 0
		if err := spacego.GetAllSubscriptionsForNode(node, func(*bcgo.BlockEntry, *financego.Subscription) error {
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
	f.Dialog.Resize(bcui.DialogSize)
}

func (f *SpaceFyne) ShowRegistrarSelectionDialog(client *spaceclientgo.SpaceClient, node *bcgo.Node) {
	// Show progress dialog
	progress := dialog.NewProgressInfinite("Updating", "Getting Registrars", f.Window)
	progress.Show()

	list := ui.NewRegistrarList(f.ShowRegistrarDialog(client, node))

	// Update list
	err := list.Update(client, node)

	// Hide progress dialog
	progress.Hide()

	if err != nil {
		f.ShowError(fmt.Errorf("Error updating registrar list: %s", err))
		return
	}

	if d := f.Dialog; d != nil {
		d.Hide()
	}
	f.Dialog = dialog.NewCustom("Registrars", "Done",
		container.NewBorder(
			&widget.Label{
				Text:     fmt.Sprintf("Your encrypted data will be stored by your choice of storage providers.\nWe recommend choosing at least %d registrars from the list below - the more you choose, the more resilient your data will be against the unexpected.", spacego.GetMinimumRegistrars()),
				Wrapping: fyne.TextWrapWord,
			},
			nil,
			nil,
			nil,
			list,
		),
		f.Window)
	f.Dialog.Show()
	f.Dialog.Resize(bcui.DialogSize)
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

	// Show progress dialog
	progress := dialog.NewProgressInfinite("Updating", "Getting Registrars", f.Window)
	progress.Show()

	domains, err := f.getRegistrarDomainsForNode(client, node)
	if err != nil {
		log.Println(err)
	}

	// Hide progress dialog
	progress.Hide()

	if len(domains) == 0 {
		f.ShowRegistrarSelectionDialog(client, node)
		return
	}

	composeText := widget.NewButtonWithIcon("Text", theme.DocumentCreateIcon(), func() {
		if d := f.Dialog; d != nil {
			d.Hide()
		}
		f.ShowComposeTextDialog(client, node)
	})
	captureImage := widget.NewButtonWithIcon("Image", theme.NewThemedResource(data.CameraPhotoIcon), func() {
		if d := f.Dialog; d != nil {
			d.Hide()
		}
		// TODO
		f.ShowError(fmt.Errorf("Not yet implemented: %s %s", "SpaceFyne.Add", "Image"))
	})
	captureVideo := widget.NewButtonWithIcon("Video", theme.NewThemedResource(data.CameraVideoIcon), func() {
		if d := f.Dialog; d != nil {
			d.Hide()
		}
		// TODO
		f.ShowError(fmt.Errorf("Not yet implemented: %s %s", "SpaceFyne.Add", "Video"))
	})
	captureAudio := widget.NewButtonWithIcon("Audio", theme.NewThemedResource(data.MicrophoneIcon), func() {
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
	// Hide progress dialog
	defer progress.Hide()

	hash, err := base64.RawURLEncoding.DecodeString(id)
	if err != nil {
		f.ShowError(err)
		return
	}

	view, err := viewer.ForMime(meta.GetType())
	if err != nil {
		f.ShowError(err)
		return
	}
	if view == nil {
		f.ShowError(fmt.Errorf("Not yet implemented: %s %s", "SpaceFyne.ShowFile", meta.Type))
		return
	}

	// Create goroutine to load file contents and update viewer
	go func() {
		reader, err := client.ReadFile(node, hash)
		if err != nil {
			f.ShowError(err)
			return
		}
		if err := view.SetSource(reader); err != nil {
			f.ShowError(err)
			return
		}
	}()

	name := meta.Name
	if name == "" {
		name = "(untitled)"
	}
	window := f.App.NewWindow(fmt.Sprintf("%s - %s - %s", bcgo.TimestampToString(timestamp), name, id[:8]))
	window.SetContent(view)
	window.Resize(bcui.WindowSize)
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

	list := ui.NewRegistrarList(f.ShowRegistrarDialog(client, node))

	// Update list
	err = list.Update(client, node)

	// Hide progress dialog
	progress.Hide()

	if err != nil {
		f.ShowError(fmt.Errorf("Error updating registrar list: %s", err))
		return
	}

	if d := f.Dialog; d != nil {
		d.Hide()
	}

	// Show registrar list
	f.Dialog = dialog.NewCustom("Registrars", "OK", list, f.Window)
	f.Dialog.Show()
	f.Dialog.Resize(bcui.DialogSize)
}

func (f *SpaceFyne) ShowRegistrarDialog(client *spaceclientgo.SpaceClient, node *bcgo.Node) func(id string, timestamp uint64, registrar *spacego.Registrar, registration *financego.Registration, subscription *financego.Subscription) {
	return func(id string, timestamp uint64, registrar *spacego.Registrar, registration *financego.Registration, subscription *financego.Subscription) {
		form := widget.NewForm(
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
		)
		if registration != nil {
			form.Append("Customer", &widget.Label{
				Text: registration.CustomerId,
			})
		} else {
			form.Append("", widget.NewButton("Register", func() {
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
			}))
		}
		if subscription != nil {
			form.Append("Subscription", &widget.Label{
				Text: subscription.SubscriptionId,
			})
			form.Append("Subscription Item", &widget.Label{
				Text: subscription.SubscriptionItemId,
			})

			/* TODO
			- Payment Methods
			- Usage Records
			- Invoices & Reciepts
			*/
		} else if registration != nil {
			form.Append("", widget.NewButton("Subscribe", func() {
				u, err := url.Parse(fmt.Sprintf("https://%s/%s", registrar.Merchant.Domain, registrar.Service.SubscribeUrl))
				if err != nil {
					f.ShowError(err)
					return
				}
				params := url.Values{}
				params.Add("alias", node.Alias)
				params.Add("customerId", registration.CustomerId)
				u.RawQuery = params.Encode()
				if err := f.App.OpenURL(u); err != nil {
					f.ShowError(err)
					return
				}
			}))
		}

		// Show detailed information
		info := dialog.NewCustom(registrar.Merchant.Alias, "OK", form, f.Window)
		info.Show()
		info.Resize(bcui.DialogSize)
	}
}

func (f *SpaceFyne) ShowHelp(client *spaceclientgo.SpaceClient) {
	f.ShowError(fmt.Errorf("Not yet implemented: %s", "SpaceFyne.ShowHelp"))
}

// ShowComposeTextDialog displays a dialog for creating a note, and adds the resulting file.
func (f *SpaceFyne) ShowComposeTextDialog(client *spaceclientgo.SpaceClient, node *bcgo.Node) {
	title := widget.NewEntry()
	title.Validator = func(s string) error {
		if s == "" {
			return errors.New("Title cannot be empty")
		}
		return nil
	}
	content := widget.NewMultiLineEntry()
	items := []*widget.FormItem{
		widget.NewFormItem("Title", title),
		widget.NewFormItem("Content", content),
	}
	f.Dialog = dialog.NewForm("Compose", "Save", "Cancel", items, func(b bool) {
		if !b {
			return
		}
		name := title.Text

		// Show progress dialog
		progress := dialog.NewProgress("Uploading", "Uploading "+name, f.Window)
		progress.Show()
		listener := &bcui.ProgressMiningListener{Func: progress.SetValue}

		reference, err := client.Add(node, listener, name, spacego.MIME_TYPE_TEXT_PLAIN, strings.NewReader(content.Text))

		// Hide progress dialog
		progress.Hide()

		if err != nil {
			f.ShowError(err)
			return
		}
		log.Println("Uploaded:", reference)
	}, f.Window)
	f.Dialog.Show()
	f.Dialog.Resize(bcui.DialogSize)
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

		// Show confirmation dialog so user can see preview and change name, mime, etc.
		uri := reader.URI()
		name := widget.NewEntry()
		name.SetText(uri.Name())
		mime := widget.NewSelect(spacego.GetMimeTypes(), nil)
		mime.Selected = uri.MimeType()
		size := widget.NewLabel("0bytes")
		prop := canvas.NewRectangle(color.Transparent)
		prop.SetMinSize(fyne.NewSize(200, 200))
		noPreview := widget.NewLabel("No Preview")
		preview := container.NewMax(prop, noPreview)
		form := widget.NewForm(
			widget.NewFormItem("Name", name),
			widget.NewFormItem("Type", mime),
			widget.NewFormItem("Size", size),
			widget.NewFormItem("Preview", preview),
		)

		var buffer []byte
		loadPreview := func(mime string) {
			if view, err := viewer.ForMime(mime); err != nil || view == nil {
				preview.Objects[1] = noPreview
			} else {
				preview.Objects[1] = view
				// Load file contents and update viewer
				if err := view.SetSource(bytes.NewReader(buffer)); err != nil {
					log.Println(err)
				}
			}
			form.Refresh()
		}
		mime.OnChanged = func(mime string) {
			go loadPreview(mime)
		}

		go func() {
			buffer, err = ioutil.ReadAll(reader)
			if err != nil {
				f.ShowError(err)
				return
			}
			size.SetText(bcgo.BinarySizeToString(uint64(len(buffer))))
			loadPreview(mime.Selected)
		}()

		f.Dialog = dialog.NewCustomConfirm("Upload File", "Upload", "Cancel", form, func(result bool) {
			if result {
				f.UploadFile(client, node, name.Text, mime.Selected, bytes.NewReader(buffer))
			}
		}, f.Window)
		f.Dialog.Show()
		f.Dialog.Resize(bcui.DialogSize)
	}, f.Window)
	f.Dialog.Show()
	f.Dialog.Resize(bcui.DialogSize)
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
	f.Dialog.Resize(bcui.DialogSize)
}

func (f *SpaceFyne) UploadFile(client *spaceclientgo.SpaceClient, node *bcgo.Node, name, mime string, reader io.Reader) {
	// Show progress dialog
	progress := dialog.NewProgress("Uploading", "Uploading "+name, f.Window)
	progress.Show()
	listener := &bcui.ProgressMiningListener{Func: progress.SetValue}

	reference, err := client.Add(node, listener, name, mime, reader)

	// Hide progress dialog
	progress.Hide()

	if err != nil {
		f.ShowError(err)
		return
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

	// TODO show confirmation dialog which lists the files to be uploaded so the user can select/deselect

	// Show progress dialog
	progress := dialog.NewProgress("Uploading", "Uploading "+folder.Name(), f.Window)
	progress.Show()
	// Hide progress dialog
	defer progress.Hide()

	for i, uri := range uris {
		progress.SetValue(float64(i) / float64(count))

		// Check if URI points to Folder
		if lister, ok := uri.(fyne.ListableURI); ok {
			f.UploadFolder(client, node, lister)
			continue
		}
		// Check if URI points to Folder
		if lister, err := fynestorage.ListerForURI(uri); err == nil {
			f.UploadFolder(client, node, lister)
			continue
		}

		// URI points to File
		file, err := fynestorage.OpenFileFromURI(uri)
		if err != nil {
			log.Println(err)
			continue
		}
		f.UploadFile(client, node, uri.Name(), uri.MimeType(), file)
	}
}

func (f *SpaceFyne) getRegistrarDomainsForNode(client *spaceclientgo.SpaceClient, node *bcgo.Node) (domains []string, err error) {
	var net *bcgo.TCPNetwork
	if node.Network != nil {
		if n, ok := node.Network.(*bcgo.TCPNetwork); ok {
			net = n
		}
	}
	err = spacego.GetAllRegistrarsForNode(node, func(registrar *spacego.Registrar, registration *financego.Registration, subscription *financego.Subscription) error {
		if registrar != nil && registration != nil && subscription != nil {
			domain := registrar.Merchant.Domain
			if domain == "" {
				domain = registrar.Merchant.Alias
			}
			if net != nil {
				net.AddPeer(domain)
			}
			domains = append(domains, domain)
		}
		return nil
	})
	return
}
