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
	"aletheiaware.com/spacego"
	"encoding/base64"
	"fyne.io/fyne/v2"
)

const (
	SPACE_SCHEME        = "space"
	SPACE_SCHEME_PREFIX = "space:"
)

type SpaceURI interface {
	fyne.URI
	FileHash() []byte
}

type DeltaURI interface {
	SpaceURI
	DeltaHash() []byte
}

type FileURI interface {
	SpaceURI
	Meta() *spacego.Meta
}

type MetaURI interface {
	SpaceURI
}

type PreviewURI interface {
	SpaceURI
	PreviewHash() []byte
}

type TagURI interface {
	SpaceURI
	TagHash() []byte
}

func NewDeltaURI(fileHash, deltaHash []byte) DeltaURI {
	return &deltaURI{
		spaceURI: spaceURI{
			fileHash: fileHash,
		},
		deltaHash: deltaHash,
	}
}

func NewFileURI(fileHash []byte, meta *spacego.Meta) FileURI {
	return &fileURI{
		spaceURI: spaceURI{
			fileHash: fileHash,
		},
		meta: meta,
	}
}

func NewMetaURI(fileHash []byte) MetaURI {
	return &metaURI{
		spaceURI: spaceURI{
			fileHash: fileHash,
		},
	}
}

func NewPreviewURI(fileHash, previewHash []byte) PreviewURI {
	return &previewURI{
		spaceURI: spaceURI{
			fileHash: fileHash,
		},
		previewHash: previewHash,
	}
}

func NewTagURI(fileHash, tagHash []byte) TagURI {
	return &tagURI{
		spaceURI: spaceURI{
			fileHash: fileHash,
		},
		tagHash: tagHash,
	}
}

type spaceURI struct {
	fileHash []byte
}

func (u *spaceURI) Authority() string {
	return ""
}

func (u *spaceURI) FileHash() []byte {
	return u.fileHash
}

func (u *spaceURI) Extension() string {
	return ""
}

func (u *spaceURI) Fragment() string {
	return ""
}

func (u *spaceURI) MimeType() string {
	return ""
}

func (u *spaceURI) Name() string {
	return ""
}

func (u *spaceURI) Path() string {
	return base64.RawURLEncoding.EncodeToString(u.fileHash)
}

func (u *spaceURI) Query() string {
	return ""
}

func (u *spaceURI) Scheme() string {
	return SPACE_SCHEME
}

func (u *spaceURI) String() string {
	return SPACE_SCHEME_PREFIX + u.Path()
}

type deltaURI struct {
	spaceURI
	deltaHash []byte
}

func (u *deltaURI) DeltaHash() []byte {
	return u.deltaHash
}

func (u *deltaURI) MimeType() string {
	return "" // TODO protobuf type
}

func (u *deltaURI) Name() string {
	return base64.RawURLEncoding.EncodeToString(u.deltaHash)
}

func (u *deltaURI) Path() string {
	return u.spaceURI.Path() + "/delta/" + base64.RawURLEncoding.EncodeToString(u.deltaHash)
}

func (u *deltaURI) String() string {
	return SPACE_SCHEME_PREFIX + u.Path()
}

type fileURI struct {
	spaceURI
	meta *spacego.Meta
}

func (u *fileURI) Meta() *spacego.Meta {
	return u.meta
}

func (u *fileURI) MimeType() string {
	return u.meta.Type
}

func (u *fileURI) Name() string {
	return u.meta.Name
}

type metaURI struct {
	spaceURI
}

func (u *metaURI) MimeType() string {
	return "" // TODO protobuf type
}

func (u *metaURI) Name() string {
	return base64.RawURLEncoding.EncodeToString(u.fileHash)
}

func (u *metaURI) Path() string {
	return u.spaceURI.Path() + "/meta"
}

func (u *metaURI) String() string {
	return SPACE_SCHEME_PREFIX + u.Path()
}

type previewURI struct {
	spaceURI
	previewHash []byte
}

func (u *previewURI) MimeType() string {
	return "" // TODO protobuf type
}

func (u *previewURI) Name() string {
	return base64.RawURLEncoding.EncodeToString(u.previewHash)
}

func (u *previewURI) Path() string {
	return u.spaceURI.Path() + "/preview/" + u.Name()
}

func (u *previewURI) PreviewHash() []byte {
	return u.previewHash
}

func (u *previewURI) String() string {
	return SPACE_SCHEME_PREFIX + u.Path()
}

type tagURI struct {
	spaceURI
	tagHash []byte
}

func (u *tagURI) MimeType() string {
	return "" // TODO protobuf type
}

func (u *tagURI) Name() string {
	return base64.RawURLEncoding.EncodeToString(u.tagHash)
}

func (u *tagURI) Path() string {
	return u.spaceURI.Path() + "/tag/" + u.Name()
}

func (u *tagURI) String() string {
	return SPACE_SCHEME_PREFIX + u.Path()
}

func (u *tagURI) TagHash() []byte {
	return u.tagHash
}
