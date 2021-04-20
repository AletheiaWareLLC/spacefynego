package storage_test

import (
	"aletheiaware.com/spacefynego/storage"
	"aletheiaware.com/spacego"
	"encoding/base64"
	"testing"
	"github.com/stretchr/testify/assert"
)

const (
	FileName    = "Test.txt"
	FileType    = spacego.MIME_TYPE_TEXT_PLAIN
	FileHash    = "abcd1234"
	DeltaHash   = "efgh5678"
	PreviewHash = "ijkl9012"
	TagHash     = "mnop3456"

	DeltaURI   = "space:abcd1234/delta/efgh5678"
	FileURI    = "space:abcd1234"
	MetaURI    = "space:abcd1234/meta"
	PreviewURI = "space:abcd1234/preview/ijkl9012"
	TagURI     = "space:abcd1234/tag/mnop3456"
)

func TestURI_DeltaURI_String(t *testing.T) {
	assert.Equal(t, DeltaURI, newDeltaURI(t).String())
}

func TestURI_FileURI_String(t *testing.T) {
	assert.Equal(t, FileURI, newFileURI(t).String())
}

func TestURI_MetaURI_String(t *testing.T) {
	assert.Equal(t, MetaURI, newMetaURI(t).String())
}

func TestURI_PreviewURI_String(t *testing.T) {
	assert.Equal(t, PreviewURI, newPreviewURI(t).String())
}

func TestURI_TagURI_String(t *testing.T) {
	assert.Equal(t, TagURI, newTagURI(t).String())
}

func fileHash(t *testing.T) []byte {
	hash, err := base64.RawURLEncoding.DecodeString(FileHash)
	assert.Nil(t, err)
	return hash
}

func newDeltaURI(t *testing.T) storage.DeltaURI {
	hash, err := base64.RawURLEncoding.DecodeString(DeltaHash)
	assert.Nil(t, err)
	return storage.NewDeltaURI(fileHash(t), hash)
}

func newFileURI(t *testing.T) storage.FileURI {
	return storage.NewFileURI(fileHash(t), &spacego.Meta{
		Name: FileName,
		Type: FileType,
	})
}

func newMetaURI(t *testing.T) storage.MetaURI {
	return storage.NewMetaURI(fileHash(t))
}

func newPreviewURI(t *testing.T) storage.PreviewURI {
	hash, err := base64.RawURLEncoding.DecodeString(PreviewHash)
	assert.Nil(t, err)
	return storage.NewPreviewURI(fileHash(t), hash)
}

func newTagURI(t *testing.T) storage.TagURI {
	hash, err := base64.RawURLEncoding.DecodeString(TagHash)
	assert.Nil(t, err)
	return storage.NewTagURI(fileHash(t), hash)
}
