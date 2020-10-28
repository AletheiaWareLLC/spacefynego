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
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/AletheiaWareLLC/bcfynego"
	bcui "github.com/AletheiaWareLLC/bcfynego/ui"
	"github.com/AletheiaWareLLC/bcgo"
	"github.com/AletheiaWareLLC/spaceclientgo"
	"github.com/AletheiaWareLLC/spacefynego/ui"
	"github.com/AletheiaWareLLC/spacefynego/ui/data"
	"github.com/AletheiaWareLLC/spacefynego/ui/viewer"
	"github.com/AletheiaWareLLC/spacego"
	"io"
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
	f.OnSignedUp = func(node *bcgo.Node) {
		f.ShowWelcome(c, node)
	}
	return f
}

// ShowWelcome displays a wizard to welcome a new user and walk them through the setup process.
func (f *SpaceFyne) ShowWelcome(client *spaceclientgo.SpaceClient, node *bcgo.Node) {
	if f.Dialog != nil {
		f.Dialog.Hide()
	}
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
	go l.Update(client, node)
	f.Dialog = dialog.NewCustom("Welcome", "OK",
		container.NewGridWithRows(3,
			widget.NewLabel(fmt.Sprintf("Hello %s, welcome to S P A C E", node.Alias)),
			widget.NewLabel("Your encrypted data will be stored by your choice of storage providers.\nChoose at least two providers from the list below;"),
			l,
		),
		f.Window)
	f.Dialog.Resize(fyne.NewSize(0, 300))
	f.Dialog.Show()
}

func (f *SpaceFyne) GetIcon() fyne.CanvasObject {
	icon := canvas.NewImageFromResource(data.SpaceIcon)
	icon.FillMode = canvas.ImageFillContain
	return icon
}

// AddFile displays a dialog (write text, take a picture, upload an existing file), and adds the result.
func (f *SpaceFyne) AddFile(client *spaceclientgo.SpaceClient) {
	composeText := widget.NewButtonWithIcon("Text", theme.DocumentCreateIcon(), func() {
		if f.Dialog != nil {
			f.Dialog.Hide()
		}
		// TODO
		f.ShowError(fmt.Errorf("Not yet implemented: %s %s", "SpaceFyne.AddFile", "Text"))
	})
	captureImage := widget.NewButtonWithIcon("Image", theme.NewThemedResource(data.CameraPhotoIcon, nil), func() {
		if f.Dialog != nil {
			f.Dialog.Hide()
		}
		// TODO
		f.ShowError(fmt.Errorf("Not yet implemented: %s %s", "SpaceFyne.AddFile", "Image"))
	})
	captureVideo := widget.NewButtonWithIcon("Video", theme.NewThemedResource(data.CameraVideoIcon, nil), func() {
		if f.Dialog != nil {
			f.Dialog.Hide()
		}
		// TODO
		f.ShowError(fmt.Errorf("Not yet implemented: %s %s", "SpaceFyne.AddFile", "Video"))
	})
	captureAudio := widget.NewButtonWithIcon("Audio", theme.NewThemedResource(data.MicrophoneIcon, nil), func() {
		if f.Dialog != nil {
			f.Dialog.Hide()
		}
		// TODO
		f.ShowError(fmt.Errorf("Not yet implemented: %s %s", "SpaceFyne.AddFile", "Audio"))
	})
	uploadFile := widget.NewButtonWithIcon("Document", theme.FileIcon(), func() {
		if f.Dialog != nil {
			f.Dialog.Hide()
		}
		go f.UploadFile(client)
	})
	content := container.NewAdaptiveGrid(3,
		composeText,
		captureImage,
		captureAudio,
		captureVideo,
		uploadFile,
	)
	if f.Dialog != nil {
		f.Dialog.Hide()
	}
	f.Dialog = dialog.NewCustom("Add File", "Cancel", content, f.Window)
	f.Dialog.Show()

	/*
		edit := viewer.GetEditor(meta, func(writer io.Writer) uint64 {
			hash, err := base64.RawURLEncoding.DecodeString(id)
			if err != nil {
				f.ShowError(err)
				return 0
			}
			count, err := client.Get(node, hash, writer)
			if err != nil {
				f.ShowError(err)
				return 0
			}
			return count
		})

		if edit == nil {
			f.ShowError(fmt.Errorf("Not yet implemented: %s %s", "SpaceFyne.NewFile", mime))
			return
		}

		if name == "" {
			name = "(untitled)"
		}
		window := f.App.NewWindow(name)
		window.SetContent(edit)
		window.Resize(fyne.NewSize(800, 600))
		window.CenterOnScreen()
		window.Show()
	*/
}

// UploadFile displays a file picker, and adds the resulting file.
func (f *SpaceFyne) UploadFile(client *spaceclientgo.SpaceClient) {
	node, err := f.GetNode(&client.BCClient)
	if err != nil {
		f.ShowError(err)
		return
	}

	// TODO ensure alias has chosen providers, and those providers are added to node.Network.Peers

	if f.Dialog != nil {
		f.Dialog.Hide()
	}
	f.Dialog = dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err == nil && reader == nil {
			return
		}
		if err != nil {
			f.ShowError(err)
			return
		}

		// Show progress dialog
		progress := dialog.NewProgress("Uploading", "Uploading "+reader.Name(), f.Window)
		progress.Show()
		defer progress.Hide()
		listener := &bcui.ProgressMiningListener{Func: progress.SetValue}

		reference, err := client.Add(node, listener, reader.Name(), reader.URI().MimeType(), reader)
		if err != nil {
			f.ShowError(err)
		}
		fmt.Println("Uploaded:", reference)
	}, f.Window)
	//f.Dialog.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".txt"}))
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

	view := viewer.GetViewer(meta, func(writer io.Writer) uint64 {
		hash, err := base64.RawURLEncoding.DecodeString(id)
		if err != nil {
			f.ShowError(err)
			return 0
		}
		// TODO display and update progress bar
		count, err := client.Read(node, hash, writer)
		if err != nil {
			f.ShowError(err)
			return 0
		}
		return count
	})
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

// Adds file
// - SpaceClient.Add(node *bcgo.Node, listener bcgo.MiningListener, name, mime string, reader io.Reader) (*bcgo.Reference, error)
// Adds file using Remote Mining Service
// - SpaceClient.AddRemote(node *bcgo.Node, domain, name, mime string, reader io.Reader) (*bcgo.Reference, error)
// Get file owned by key with given hash
// - SpaceClient.Get(node *bcgo.Node, recordHash []byte, callback MetaCallback) error
// Get file shared to key with given hash
// - SpaceClient.GetShared(node *bcgo.Node, recordHash []byte, callback MetaCallback) error
// Get all files owned by key with given mime-type
// - SpaceClient.GetAll(node *bcgo.Node, mime string, callback MetaCallback) error
// Get all files shared to key with given mime-type
// - SpaceClient.GetAllShared(node *bcgo.Node, mime string, callback MetaCallback) error
// List files owned by key
// - SpaceClient.List(node *bcgo.Node, callback MetaCallback) error
// List files shared with key
// - SpaceClient.ListShared(node *bcgo.Node, callback MetaCallback) error
// Search files owned by key
// - SpaceClient.Search(node *bcgo.Node, terms []string, callback MetaCallback) error
// Search files shared with key
// - SpaceClient.SearchShared(node *bcgo.Node, terms []string, callback MetaCallback) error
// Read file by given hash
// - SpaceClient.Read(node *bcgo.Node, recordHash []byte, writer io.Writer) (uint64, error)
// Read file shared to key with given hash
// - SpaceClient.ReadShared(node *bcgo.Node, recordHash []byte, writer io.Writer) (uint64, error)
// Share file with recipients
// - SpaceClient.Share(node *bcgo.Node, listener bcgo.MiningListener, recordHash []byte, recipients []string) error
// Tag file owned by key
// - SpaceClient.Tag(node *bcgo.Node, listener bcgo.MiningListener, recordHash []byte, tag []string) ([]*bcgo.Reference, error)
// Tag file shared with key
// - SpaceClient.TagShared(node *bcgo.Node, listener bcgo.MiningListener, recordHash []byte, tag []string) ([]*bcgo.Reference, error)
// Get all tags for the file with the given hash
// - SpaceClient.GetTag(node *bcgo.Node, recordHash []byte, callback func(entry *bcgo.BlockEntry, tag *spacego.Tag)) error
// - SpaceClient.GetRegistration(merchant string, callback func(*financego.Registration) error) error
// - SpaceClient.GetSubscription(merchant string, callback func(*financego.Subscription) error) error

func (f *SpaceFyne) ShowStorage(client *spaceclientgo.SpaceClient) {
	f.ShowError(fmt.Errorf("Not yet implemented: %s", "SpaceFyne.ShowStorage"))
}

func (f *SpaceFyne) ShowHelp(client *spaceclientgo.SpaceClient) {
	f.ShowError(fmt.Errorf("Not yet implemented: %s", "SpaceFyne.ShowHelp"))
}
