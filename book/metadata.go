package book

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	e "github.com/barsanuphe/endive/endive"
	h "github.com/barsanuphe/helpers"
	i "github.com/barsanuphe/helpers/ui"
)

const (
	titleField         = "title"
	descriptionField   = "description"
	isbnField          = "isbn"
	yearField          = "year"
	editionYearField   = "edition_year"
	authorField        = "author"
	publisherField     = "publisher"
	tagsField          = "tags"
	seriesField        = "series"
	languageField      = "language"
	categoryField      = "category"
	typeField          = "type"
	genreField         = "genre"
	numPagesField      = "numpages"
	averageRatingField = "averagerating"

	authorUsage      = "Authors can be edited as a comma-separated list of strings."
	categoryUsage    = "A book can be either fiction or nonfiction."
	typeUsage        = "The nature of this book."
	tagsUsage        = "Tags can be edited as a comma-separated list of strings."
	seriesUsage      = "Series can be edited as a comma-separated list of 'series name:index' strings. Index can be empty, or a range."
	yearUsage        = "The year in which the book was written."
	editionYearUsage = "The year in which this edition was published."
	publisherUsage   = "Publisher of this edition."
	languageUsage    = "Language of this edition."
	genreUsage       = "Main genre of this book."
	isbnUsage        = "ISBN13 for this edition."
	titleUsage       = "Title, without series information."
	descriptionUsage = "Description for this edition."

	unknownYear = "XXXX"
	unknown     = "Unknown"

	localSource  = "Epub"
	onlineSource = "Online"

	invalidField   = "Invalid field: %s"
	cannotSetField = "Cannot set field %s"
)

// MetadataFieldNames is a list of valid field names
var MetadataFieldNames = []string{authorField, titleField, yearField, editionYearField, publisherField, descriptionField, languageField, categoryField, typeField, genreField, tagsField, seriesField, isbnField}
var metadataFieldMap = map[string]string{
	authorField:      "Authors",
	titleField:       "BookTitle",
	yearField:        "OriginalYear",
	editionYearField: "EditionYear",
	publisherField:   "Publisher",
	descriptionField: "Description",
	languageField:    "Language",
	categoryField:    "Category",
	typeField:        "Type",
	genreField:       "Genre",
	tagsField:        "Tags",
	seriesField:      "Series",
	isbnField:        "ISBN",
}

// Metadata contains all of the known book metadata.
type Metadata struct {
	BookTitle     string   `json:"title" xml:"title"`
	ImageURL      string   `json:"image_url" xml:"image_url"`
	NumPages      string   `json:"num_pages" xml:"num_pages"`
	Authors       []string `json:"authors" xml:"authors>author>name"`
	ISBN          string   `json:"isbn" xml:"isbn13"`
	OriginalYear  string   `json:"year" xml:"work>original_publication_year"`
	EditionYear   string   `json:"edition_year" xml:"publication_year"`
	Description   string   `json:"description" xml:"description"`
	Series        Series   `json:"series" xml:"series_works>series_work"`
	AverageRating string   `json:"average_rating" xml:"average_rating"`
	Tags          Tags     `json:"tags" xml:"popular_shelves>shelf"`
	Category      string   `json:"category"`
	Type          string   `json:"type"`
	Genre         string   `json:"genre"`
	Language      string   `json:"language" xml:"language_code"`
	Publisher     string   `json:"publisher" xml:"publisher"`
}

// String returns a representation of Metadata
func (i *Metadata) String() string {
	if len(i.Series) != 0 {
		return fmt.Sprintf("%s (%s) %s [%s]", i.Author(), i.OriginalYear, i.Title(), i.MainSeries().String())
	}
	return fmt.Sprintf("%s (%s) %s", i.Author(), i.OriginalYear, i.Title())
}

// HasAny checks if metadata was parsed.
func (i *Metadata) HasAny() bool {
	// if Metadata does not have a title and author, chances are it's empty.
	if i.Title() != "" && i.Author() != "" {
		return true
	}
	return false
}

// IsComplete checks if metadata looks complete
func (i *Metadata) IsComplete() bool {
	hasAuthor := i.Author() != ""
	hasTitle := i.Title() != ""
	hasYear := i.OriginalYear != "" && i.OriginalYear != unknownYear
	hasLanguage := i.Language != ""
	hasDescription := i.Description != ""
	hasCategory := i.Category != "" && i.Category != unknown
	hasType := i.Type != "" && i.Type != unknown
	hasGenre := i.Genre != "" && i.Genre != unknown
	hasISBN := i.ISBN != ""
	hasPublisher := i.Publisher != ""
	hasTags := i.Tags.String() != ""
	return hasAuthor && hasTitle && hasYear && hasLanguage && hasDescription && hasCategory && hasType && hasGenre && hasISBN && hasPublisher && hasTags
}

// Title returns Metadata's main title.
func (i *Metadata) Title() string {
	return i.BookTitle
}

// Clean cleans up the Metadata
func (i *Metadata) Clean(cfg e.Config) {
	// default year
	if i.OriginalYear == "" {
		if i.EditionYear != "" {
			i.OriginalYear = i.EditionYear
		} else {
			i.OriginalYear = unknownYear
		}
	}
	if i.EditionYear == "" {
		if i.OriginalYear != "" {
			i.EditionYear = i.OriginalYear
		} else {
			i.EditionYear = unknownYear
		}
	}
	// clean description
	i.Description = cleanHTML(i.Description)
	// clean language
	i.Language = cleanLanguage(i.Language)
	// use config aliases
	i.useAliases(cfg)
	// clean tags
	i.Tags.Clean()
	// autofill category
	if i.Category == "" {
		for _, possibleCategory := range validCategories {
			if isIn, _ := i.Tags.Has(Tag{Name: possibleCategory}); isIn {
				i.Category = possibleCategory
				i.Tags.RemoveFromNames(possibleCategory)
				break
			}
		}
	}
	// if nothing valid found...
	if i.Category == "" {
		i.Category = unknown
	}
	if cat, err := cleanCategory(i.Category); err == nil {
		i.Category = cat
		i.Tags.RemoveFromNames(cat)
	}

	// autofill type
	if i.Type == "" {
		for _, possibleType := range validTypes {
			if isIn, _ := i.Tags.Has(Tag{Name: possibleType}); isIn {
				i.Type = possibleType
				i.Tags.RemoveFromNames(possibleType)
				break
			}
		}
	}
	// if nothing found, unknown.
	if i.Type == "" {
		i.Type = unknown
	}
	if tp, err := cleanType(i.Type); err == nil {
		i.Type = tp
		i.Tags.RemoveFromNames(tp)
	}

	// MainGenre
	if i.Genre == "" && len(i.Tags) != 0 {
		if cleanName, err := cleanTagName(i.Tags[0].Name); err == nil {
			i.Genre = cleanName
			i.Tags.RemoveFromNames(i.Genre)
		}
	}
	// if nothing valid found...
	if i.Genre == "" {
		i.Genre = unknown
	}

	// clean series
	for j := range i.Series {
		i.Series[j].Name = strings.TrimSpace(i.Series[j].Name)
	}
	// clean publisher
	i.Publisher = strings.TrimSpace(i.Publisher)
	// use config aliases, again, to clean up new values for maingenre, category, etc
	i.useAliases(cfg)
}

// useAliases updates Metadata fields, using the configuration file.
func (i *Metadata) useAliases(cfg e.Config) {
	// author aliases
	for j, author := range i.Authors {
		for mainAlias, aliases := range cfg.AuthorAliases {
			_, isIn := h.StringInSlice(author, aliases)
			if isIn {
				i.Authors[j] = mainAlias
				break
			}
		}
	}
	// tag aliases
	cleanTags := Tags{}
	for _, tag := range i.Tags {
		added := false
		for mainAlias, aliases := range cfg.TagAliases {
			_, isIn := h.StringInSlice(tag.Name, aliases)
			if isIn {
				cleanTags.AddFromNames(mainAlias)
				added = true
				break
			}
		}
		// if no alias found, add directly
		if !added {
			cleanTags.AddFromNames(tag.Name)
		}
	}
	i.Tags = cleanTags
	// genre aliases (same as tags)
	for mainAlias, aliases := range cfg.TagAliases {
		_, isIn := h.StringInSlice(i.Genre, aliases)
		if isIn {
			i.Genre = mainAlias
			break
		}
	}
	// type aliases (same as tags)
	for mainAlias, aliases := range cfg.TagAliases {
		_, isIn := h.StringInSlice(i.Type, aliases)
		if isIn {
			i.Type = mainAlias
			break
		}
	}
	// publisher aliases
	for mainAlias, aliases := range cfg.PublisherAliases {
		_, isIn := h.StringInSlice(i.Publisher, aliases)
		if isIn {
			i.Publisher = mainAlias
			break
		}
	}
}

// Author returns Metadata's main author.
func (i *Metadata) Author() string {
	if len(i.Authors) != 0 {
		return strings.Join(i.Authors, ", ")
	}
	return unknown
}

// MainSeries return the main Series of Metadata.
func (i *Metadata) MainSeries() SingleSeries {
	if len(i.Series) != 0 {
		return i.Series[0]
	}
	return SingleSeries{}
}

// IsSimilar checks if metadata is similar to known Metadata.
func (i *Metadata) IsSimilar(o Metadata) bool {
	// check isbn
	if i.ISBN != "" && o.ISBN != "" && i.ISBN == o.ISBN {
		return true
	}
	// similar == same author/title, for now
	if i.Author() == o.Author() && i.Title() == o.Title() {
		return true
	}
	return false
}

// OutputDiffTable returns differences between Metadatas.
func (i *Metadata) OutputDiffTable(o *Metadata, diffOnly bool) [][]string {
	var rows [][]string
	for _, field := range MetadataFieldNames {
		iValue, err := i.Get(field)
		if err != nil {
			iValue = couldNotRetrieveValue
		}
		oValue, err := o.Get(field)
		if err != nil {
			oValue = couldNotRetrieveValue
		}
		if !diffOnly || iValue != oValue {
			rows = append(rows, []string{iValue, oValue})
		}
	}
	return rows
}

// Diff returns differences between Metadatas.
func (i *Metadata) Diff(o *Metadata, firstHeader, secondHeader string) string {
	var rows [][]string
	rows = append(rows, []string{i.String(), o.String()})
	rows = append(rows, i.OutputDiffTable(o, false)...)
	return e.TabulateRows(rows, firstHeader, secondHeader)
}

// Merge with another Metadata.
func (i *Metadata) Merge(o *Metadata, cfg e.Config, ui i.UserInterface, diffOnly bool) (err error) {
	for _, field := range MetadataFieldNames {
		err = i.MergeField(o, field, cfg, ui, diffOnly)
		if err != nil {
			return
		}
	}
	// automatically fill fields usually not found in epubs.
	i.ImageURL = getLargeGRUrl(o.ImageURL)
	i.NumPages = o.NumPages
	i.AverageRating = o.AverageRating
	i.Clean(cfg)
	return
}

// Get Metadata field value
func (i *Metadata) Get(field string) (string, error) {
	publicFieldName, structField, _, err := getField(i, metadataFieldMap, field)
	if err != nil {
		return "", err
	}
	switch publicFieldName {
	case tagsField:
		return i.Tags.String(), nil
	case seriesField:
		return i.Series.rawString(), nil
	case authorField:
		return i.Author(), nil
	default:
		return structField.String(), nil
	}
}

// Set Metadata field with a string value
func (i *Metadata) Set(field, value string) error {
	publicFieldName, structField, canBeSet, err := getField(i, metadataFieldMap, field)
	if err != nil {
		return err
	}
	if !canBeSet {
		return fmt.Errorf(cannotSetField, field)
	}
	// set value
	switch publicFieldName {
	case tagsField:
		value = strings.ToLower(value)
		i.Tags = Tags{}
		i.Tags.AddFromNames(strings.Split(value, ",")...)
	case seriesField:
		i.Series = Series{}
		if value != "" {
			for _, s := range strings.Split(value, ",") {
				if _, err := i.Series.AddFromString(s); err != nil {
					return err
				}
			}
		}
	case authorField:
		i.Authors = strings.Split(value, ",")
		for j := range i.Authors {
			i.Authors[j] = strings.TrimSpace(i.Authors[j])
		}
	case yearField, editionYearField:
		// check it's a correct year
		_, err := strconv.Atoi(value)
		if err != nil {
			return errors.New("Invalid year value: " + value)
		}
		structField.SetString(value)
	case isbnField:
		// check it's a correct isbn
		isbn, err := e.CleanISBN(value)
		if err != nil {
			return err
		}
		structField.SetString(isbn)
	case categoryField:
		cleanCategory, err := cleanCategory(value)
		if err != nil {
			return err
		}
		structField.SetString(cleanCategory)
	case typeField:
		cleanType, err := cleanType(value)
		if err != nil {
			return err
		}
		structField.SetString(cleanType)
	case descriptionField:
		structField.SetString(cleanHTML(value))
	case languageField:
		structField.SetString(cleanLanguage(value))
	default:
		structField.SetString(value)
	}
	return nil
}

func listLocalAndRemoteOnly(ui i.UserInterface, local, online string, options *[]string, thingsToClean ...string) {
	*options = append(*options, local, online)
	CleanSliceAndTagEntries(ui, local, online, options, thingsToClean...)
}

// MergeField with another Metadata.
func (i *Metadata) MergeField(o *Metadata, field string, cfg e.Config, ui i.UserInterface, diffOnly bool) (err error) {
	var userInput string
	options := []string{}
	usage, ok := usageMap[field]
	if !ok {
		usage = ""
	}
	currentValue, err := i.Get(field)
	if err != nil {
		return err
	}
	otherValue, err := o.Get(field)
	if err != nil {
		return err
	}
	// if only merging the differences and values are the same, just return
	if diffOnly && currentValue == otherValue {
		return nil
	}

	switch field {
	case yearField, editionYearField:
		listLocalAndRemoteOnly(ui, currentValue, otherValue, &options, unknownYear)
		title := "Publication year"
		if field == yearField {
			title = "Original " + title
		}
		userInput, err = ui.SelectOption(title, usage, options, false)
	case languageField:
		listLocalAndRemoteOnly(ui, cleanLanguage(i.Language), cleanLanguage(o.Language), &options, unknown)
		userInput, err = ui.SelectOption(strings.Title(field), usage, options, false)
	case categoryField:
		options = append(options, validCategories...)
		CleanSliceAndTagEntries(ui, currentValue, otherValue, &options, unknown)
		userInput, err = ui.SelectOption(strings.Title(field), usage, options, false)
	case typeField:
		options = append(options, validTypes...)
		CleanSliceAndTagEntries(ui, currentValue, otherValue, &options, unknown)
		userInput, err = ui.SelectOption(strings.Title(field), usage, options, false)
	case descriptionField:
		listLocalAndRemoteOnly(ui, cleanHTML(i.Description), cleanHTML(o.Description), &options, unknown)
		userInput, err = ui.SelectOption(strings.Title(field), usage, options, true)
	default:
		listLocalAndRemoteOnly(ui, currentValue, otherValue, &options, unknown)
		userInput, err = ui.SelectOption(strings.Title(field), usage, options, false)
	}

	// checking SelectOption err
	if err != nil {
		return
	}
	// set the field
	err = i.Set(field, userInput)
	if err != nil {
		return
	}
	i.Clean(cfg)
	return
}

// getOnlineMetadata retrieves the online info for this book.
func (i *Metadata) getOnlineMetadata(ui i.UserInterface, cfg e.Config) (*Metadata, error) {
	if cfg.GoodReadsAPIKey == "" {
		return nil, e.WarningGoodReadsAPIKeyMissing
	}
	var err error
	var g RemoteLibraryAPI
	g = GoodReads{}
	id := ""

	// If not ISBN is found, ask for input
	if i.ISBN == "" {
		ui.Warning("Could not find ISBN.")
		isbn, err := e.AskForISBN(ui)
		if err == nil {
			i.ISBN = isbn
		}
	}
	// search by ISBN preferably
	if i.ISBN != "" {
		id, err = g.GetBookIDByISBN(i.ISBN, cfg.GoodReadsAPIKey)
		if err != nil {
			return nil, err
		}
	}
	// if no ISBN or nothing was found
	if id == "" {
		// TODO: if unsure, show hits
		id, err = g.GetBookIDByQuery(i.Author(), i.Title(), cfg.GoodReadsAPIKey)
		if err != nil {
			return nil, err
		}
	}
	// if still nothing was found...
	if id == "" {
		return nil, errors.New("Could not find online data for " + i.String())
	}
	// get book info
	onlineInfo, err := g.GetBook(id, cfg.GoodReadsAPIKey)
	if err == nil {
		onlineInfo.Clean(cfg)
	}
	return &onlineInfo, nil
}

// SearchOnline tries to find metadata from online sources.
func (i *Metadata) SearchOnline(ui i.UserInterface, cfg e.Config, fields ...string) (err error) {
	onlineInfo, err := i.getOnlineMetadata(ui, cfg)
	if err != nil {
		ui.Debug(err.Error())
		ui.Warning("Could not retrieve information from GoodReads. Manual review.")
		err = i.Merge(&Metadata{}, cfg, ui, false)
		if err != nil {
			ui.Error(err.Error())
		}
		return err
	}

	// show diff between epub and GR versions, then ask what to do.
	fmt.Println(i.Diff(onlineInfo, localSource, onlineSource))
	ui.Choice("[E]dit or [A]bort : ")
	validChoice := false
	errs := 0
	for !validChoice {
		choice, scanErr := ui.GetInput()
		if scanErr != nil {
			return scanErr
		}
		switch strings.ToLower(choice) {
		case "a":
			err = errors.New("Abort")
			validChoice = true
		case "e":
			if len(fields) == 0 {
				if err := i.Merge(onlineInfo, cfg, ui, false); err != nil {
					return err
				}
			} else {
				for _, f := range fields {
					if err := i.MergeField(onlineInfo, f, cfg, ui, false); err != nil {
						return err
					}
				}
				// automatically fill fields usually not found in epubs.
				i.ImageURL = getLargeGRUrl(onlineInfo.ImageURL)
				i.NumPages = onlineInfo.NumPages
				i.AverageRating = onlineInfo.AverageRating
				i.Clean(cfg)
			}
			validChoice = true
		default:
			fmt.Println("Invalid choice.")
			errs++
			if errs > 10 {
				return errors.New("Too many invalid choices.")
			}
		}
	}
	return
}
