package book

import (
	"fmt"
	"html"
	"strconv"
	"strings"
	"time"

	h "github.com/barsanuphe/endive/helpers"
)

const apiRoot = "https://www.goodreads.com/"

// response is the top xml element in goodreads response.
type response struct {
	Book   Info          `xml:"book"`
	Search searchResults `xml:"search"`
}

// searchResults is the main xml element in goodreads search.
type searchResults struct {
	ResultsNumber string `xml:"total-results"`
	Works         []work `xml:"results>work"`
}

// works holds the work information in the xml reponse.
type work struct {
	ID     string `xml:"best_book>id"`
	Author string `xml:"best_book>author>name"`
	Title  string `xml:"best_book>title"`
}

// GetBook returns a GoodreadsBook from its Goodreads ID
func GetBook(id, key string) Info {
	defer h.TimeTrack(time.Now(), "Getting Book info")
	uri := apiRoot + "book/show/" + id + ".xml?key=" + key
	r := response{}
	h.GetXMLData(uri, &r)
	r.Book.Tags.Clean()
	return r.Book
}

func makeSearchQuery(parts ...string) (query string) {
	query = strings.Join(parts, "+")
	r := strings.NewReplacer(" ", "+")
	return html.EscapeString(r.Replace(query))
}

// GetBookIDByQuery gets a Goodreads ID from a query
func GetBookIDByQuery(author, title, key string) (id string) {
	defer h.TimeTrack(time.Now(), "Getting Book ID")

	uri := apiRoot + "search/index.xml?key=" + key + "&q=" + makeSearchQuery(author, title)
	r := response{}
	h.GetXMLData(uri, &r)
	// parsing results
	hits, err := strconv.Atoi(r.Search.ResultsNumber)
	if err != nil {
		fmt.Println("error")
	}
	if hits != 0 {
		for _, work := range r.Search.Works {
			if work.Author == author && work.Title == title {
				return work.ID
			}
		}
		fmt.Println("Could not find exact match, returning first hit.")
		return r.Search.Works[0].ID
	}
	return
}

// GetBookIDByISBN gets a Goodreads ID from an ISBN
func GetBookIDByISBN(isbn, key string) (id string) {
	defer h.TimeTrack(time.Now(), "Getting Book ID")
	// TODO
	return
}
