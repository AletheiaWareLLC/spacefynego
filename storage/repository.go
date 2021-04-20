/*
 * Copyright 2021 Aletheia Ware LLC
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

package storage

import (
	"aletheiaware.com/bcfynego/storage"
	"aletheiaware.com/bcgo"
	"aletheiaware.com/spaceclientgo"
	"aletheiaware.com/spacego"
	"encoding/base64"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage/repository"
	"strings"
)

var RootURI = &spaceURI{}

type SpaceRepository interface {
	repository.Repository
	repository.CustomURIRepository
	repository.CopyableRepository
	repository.HierarchicalRepository
	repository.ListableRepository
	repository.MovableRepository
	repository.WritableRepository
	Register()
}

type spaceRepository struct {
	client spaceclientgo.SpaceClient
	node   bcgo.Node
}

func NewSpaceRepository(client spaceclientgo.SpaceClient, node bcgo.Node) SpaceRepository {
	return &spaceRepository{
		client: client,
		node:   node,
	}
}

func (r *spaceRepository) CanList(u fyne.URI) (bool, error) {
	if u == RootURI {
		return true, nil
	}
	// TODO
	return false, fmt.Errorf("%s: Not Yet Implemented", "SpaceRepository.CanList")
}

func (r *spaceRepository) CanRead(u fyne.URI) (bool, error) {
	if u == RootURI {
		return false, nil
	}
	// TODO
	return false, fmt.Errorf("%s: Not Yet Implemented", "SpaceRepository.CanRead")
}

func (r *spaceRepository) CanWrite(u fyne.URI) (bool, error) {
	if u == RootURI {
		return false, nil
	}
	// TODO
	return false, fmt.Errorf("%s: Not Yet Implemented", "SpaceRepository.CanWrite")
}

func (r *spaceRepository) Child(u fyne.URI, c string) (fyne.URI, error) {
	if u == RootURI {
		return r.ParseURI(SPACE_SCHEME_PREFIX + c)
	}
	// TODO
	return nil, fmt.Errorf("%s: Not Yet Implemented", "SpaceRepository.Child")
}

func (r *spaceRepository) Copy(src, dest fyne.URI) error {
	// TODO
	return fmt.Errorf("%s: Not Yet Implemented", "SpaceRepository.Copy")
}

func (r *spaceRepository) CreateListable(u fyne.URI) error {
	// TODO
	return fmt.Errorf("%s: Not Yet Implemented", "SpaceRepository.CreateListable")
}

func (r *spaceRepository) Delete(u fyne.URI) error {
	// BC is indelible
	return repository.ErrOperationNotSupported
}

func (r *spaceRepository) Destroy(string) {
	// Do nothing
}

func (r *spaceRepository) Exists(u fyne.URI) (bool, error) {
	// TODO
	return false, fmt.Errorf("%s: Not Yet Implemented", "SpaceRepository.Exists")
}

func (r *spaceRepository) List(u fyne.URI) ([]fyne.URI, error) {
	// TODO
	return nil, fmt.Errorf("%s: Not Yet Implemented", "SpaceRepository.List")
}

func (r *spaceRepository) Move(fyne.URI, fyne.URI) error {
	// BC is immutable
	return repository.ErrOperationNotSupported
}

func (r *spaceRepository) Parent(fyne.URI) (fyne.URI, error) {
	// TODO
	return nil, fmt.Errorf("%s: Not Yet Implemented", "SpaceRepository.Parent")
}

func (r *spaceRepository) ParseURI(s string) (fyne.URI, error) {
	if !strings.HasPrefix(s, SPACE_SCHEME_PREFIX) {
		return nil, storage.ErrInvalidURI
	}
	s = strings.TrimPrefix(s, SPACE_SCHEME_PREFIX)
	s = strings.TrimSuffix(s, "/")

	if s == "" {
		return RootURI, nil
	}

	parts := strings.Split(s, "/")

	var fileHash []byte
	h, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, err
	}
	fileHash = h

	if len(parts) == 1 {
		var meta *spacego.Meta
		if err := r.client.MetaForHash(r.node, fileHash, func(e *bcgo.BlockEntry, m *spacego.Meta) error {
			meta = m
			return nil
		}); err != nil {
			return nil, err
		}
		if meta == nil {
			return nil, fmt.Errorf("Could not load metadata for %s", s)
		}
		return NewFileURI(fileHash, meta), nil
	}

	var recordHash []byte
	if len(parts) > 2 {
		h, err := base64.RawURLEncoding.DecodeString(parts[2])
		if err != nil {
			return nil, err
		}
		recordHash = h
	}

	switch parts[1] {
	case "delta":
		return NewDeltaURI(fileHash, recordHash), nil
	case "meta":
		return NewMetaURI(fileHash), nil
	case "preview":
		return NewPreviewURI(fileHash, recordHash), nil
	case "tag":
		return NewTagURI(fileHash, recordHash), nil
	}
	return nil, storage.ErrInvalidURI
}

func (r *spaceRepository) Reader(u fyne.URI) (fyne.URIReadCloser, error) {
	// TODO
	return nil, fmt.Errorf("%s: Not Yet Implemented", "SpaceRepository.Reader")
}

func (r *spaceRepository) Register() {
	repository.Register(SPACE_SCHEME, r)
}

func (r *spaceRepository) Writer(u fyne.URI) (fyne.URIWriteCloser, error) {
	// TODO
	return nil, fmt.Errorf("%s: Not Yet Implemented", "SpaceRepository.Writer")
}
