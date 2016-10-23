package book

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetadata(t *testing.T) {
	fmt.Println("+ Testing MetaData...")
	assert := assert.New(t)
	for i, testEpub := range epubs {
		e := NewBook(ui, i, testEpub.filename, standardTestConfig, true)

		// testing HasAny
		hasMetadata := e.Metadata.HasAny()
		assert.False(hasMetadata, "Error, should not have metadata yet.")

		// reading info
		info, err := e.MainEpub().ReadMetadata()
		assert.NotNil(err, "Error should be found (no ISBN in test epubs) for "+e.FullPath())
		if err != nil {
			assert.Equal("ISBN not found in epub", err.Error(), "Error should only mention missing ISBN")
		}
		e.Metadata = info

		// testing Get, GetFirstValue
		assert.Equal(e.Metadata.Title(), testEpub.expectedTitle, "Error getting title")
		assert.Equal(e.Metadata.Author(), testEpub.expectedAuthor, "Error getting author")
		assert.Equal(e.Metadata.OriginalYear, testEpub.expectedPublicationYear, "Error getting year")
		assert.Equal(e.Metadata.Language, testEpub.expectedLanguage, "Error getting language")

		// testing HasAny
		hasMetadata = e.Metadata.HasAny()
		assert.True(hasMetadata, "Error, should have metadata")

		// testing IsSimilar
		o := Metadata{}
		assert.False(e.Metadata.IsSimilar(o), "Error: metadata should not be similar.")

		// copying manually
		o.Authors = []string{}
		o.Authors = append(o.Authors, e.Metadata.Authors...)
		o.BookTitle = e.Metadata.BookTitle
		// checking again
		assert.True(e.Metadata.IsSimilar(o), "Error: metadata should be similar.")

		// Set
		err = e.Metadata.Set("ISBN", "hihi")
		assert.NotNil(err, "Invalid ISBN Format")
		err = e.Metadata.Set("isbn", "9780340839935")
		assert.Nil(err, "Setting ISBN should work.")
		assert.Equal("9780340839935", e.Metadata.ISBN)
		err = e.Metadata.Set("ISBN", "9-78-0-4--410-13-593")
		assert.Nil(err, "Setting ISBN should work.")
		assert.Equal("9780441013593", e.Metadata.ISBN)

		err = e.Metadata.Set("tags", "hihi, hoho")
		assert.Nil(err)
		assert.Equal("hihi, hoho", e.Metadata.Tags.String())

		err = e.Metadata.Set("year", "hihi")
		assert.NotNil(err, "Invalid year Format")
		err = e.Metadata.Set("edition_year", "2013")
		assert.Nil(err, "Valid year Format")
		assert.Equal("2013", e.Metadata.EditionYear)
	}
}
