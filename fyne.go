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
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/widget"
	"github.com/AletheiaWareLLC/bcfynego"
	bcui "github.com/AletheiaWareLLC/bcfynego/ui"
	"github.com/AletheiaWareLLC/spaceclientgo"
	"github.com/AletheiaWareLLC/spacefynego/ui"
	"github.com/AletheiaWareLLC/spacego"
)

type SpaceFyne struct {
	bcfynego.BCFyne
}

// Adds file
// - SpaceClient.Add(node *bcgo.Node, listener bcgo.MiningListener, name, mime string, reader io.Reader) (*bcgo.Reference, error)
// Adds file using Remote Mining Service
// - SpaceClient.AddRemote(node *bcgo.Node, domain, name, mime string, reader io.Reader) (*bcgo.Reference, error)

// NewFile displays a dialog (name & mime), and adds the resulting file.
func (f *SpaceFyne) NewFile(client *spaceclientgo.SpaceClient) {
	f.ShowError(fmt.Errorf("Not yet implemented: %s", "SpaceFyne.NewFile"))
	/* TODO
	f.ShowNewFileDialog(func(name, mime string) {
		node, err := f.GetNode(&client.BCClient)
		if err != nil {
			f.ShowError(err)
			return
		}
		editor := f.Editor(mime)
		f.ShowConfirmDialog(editor, func(result bool) {
			if result {
				client.Add(node, client.Listener, name, mime, reader)
			}
		})
	})
	*/
}

// UploadFile displays a file picker, and adds the resulting file.
func (f *SpaceFyne) UploadFile(client *spaceclientgo.SpaceClient) {
	node, err := f.GetNode(&client.BCClient)
	if err != nil {
		f.ShowError(err)
		return
	}

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

// Show file owned by key with given hash
// - SpaceClient.Show(node *bcgo.Node, recordHash []byte, callback MetaCallback) error
// Show file shared to key with given hash
// - SpaceClient.ShowShared(node *bcgo.Node, recordHash []byte, callback MetaCallback) error
// Show all files owned by key with given mime-type
// - SpaceClient.ShowAll(node *bcgo.Node, mime string, callback MetaCallback) error
// Show all files shared to key with given mime-type
// - SpaceClient.ShowAllShared(node *bcgo.Node, mime string, callback MetaCallback) error
// Get file by given hash
// - SpaceClient.Get(node *bcgo.Node, recordHash []byte, writer io.Writer) (uint64, error)
// Get file shared to key with given hash
// - SpaceClient.GetShared(node *bcgo.Node, recordHash []byte, writer io.Writer) (uint64, error)
func (f *SpaceFyne) ShowFile(client *spaceclientgo.SpaceClient, id string, meta *spacego.Meta) {
	f.ShowError(fmt.Errorf("Not yet implemented: %s", "SpaceFyne.ShowFile"))
	/* TODO
	node, err := f.GetNode(&client.BCClient)
	if err != nil {
		f.ShowError(err)
		return
	}
	*/
}

// List files owned by key
// - SpaceClient.List(node *bcgo.Node, callback MetaCallback) error
// List files shared with key
// - SpaceClient.ListShared(node *bcgo.Node, callback MetaCallback) error
func (f *SpaceFyne) NewList(client *spaceclientgo.SpaceClient) *ui.MetaList {
	return ui.NewMetaList(func(id string, meta *spacego.Meta) {
		f.ShowFile(client, id, meta)
	})
}

func (f *SpaceFyne) ShowNewFileDialog(callback func(string, string)) {
	name := widget.NewEntry() // TODO Change to SelectEntry
	mime := widget.NewEntry()

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Name", Widget: name},
			{Text: "Mime", Widget: mime},
		},
	}
	f.ShowConfirmDialog(form, func(result bool) {
		if result {
			n := name.Text
			if n == "" {
				// TODO either select lastest canvas or show error
			}
			m := mime.Text
			if m == "" {
				m = "text/plain"
			}
			callback(n, m)
		}
	})
}

func (f *SpaceFyne) ShowConfirmDialog(content fyne.CanvasObject, callback func(bool)) {
	if f.Dialog != nil {
		f.Dialog.Hide()
	}
	f.Dialog = dialog.NewCustomConfirm("New File", "Create", "Cancel", content, callback, f.Window)

	f.Dialog.Show()
}

/*
func (c *SpaceClient) Share(node *bcgo.Node, listener bcgo.MiningListener, recordHash []byte, recipients []string) error {

// Search files owned by key
func (c *SpaceClient) Search(node *bcgo.Node, terms []string, callback MetaCallback) error {

// Search files shared with key
func (c *SpaceClient) SearchShared(node *bcgo.Node, terms []string, callback MetaCallback) error {

// Tag file owned by key
func (c *SpaceClient) Tag(node *bcgo.Node, listener bcgo.MiningListener, recordHash []byte, tag []string) ([]*bcgo.Reference, error) {

// Tag file shared with key
func (c *SpaceClient) TagShared(node *bcgo.Node, listener bcgo.MiningListener, recordHash []byte, tag []string) ([]*bcgo.Reference, error) {


func (c *SpaceClient) ShowTag(node *bcgo.Node, recordHash []byte, callback func(entry *bcgo.BlockEntry, tag *spacego.Tag)) error {


func (c *SpaceClient) Registration(merchant string, callback func(*financego.Registration) error) error {


func (c *SpaceClient) Subscription(merchant string, callback func(*financego.Subscription) error) error {
*/
