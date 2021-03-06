# Endive

[![GoDoc](https://godoc.org/github.com/barsanuphe/endive?status.svg)](https://godoc.org/github.com/barsanuphe/endive)
[![Build Status](https://travis-ci.org/barsanuphe/endive.svg?branch=master)](https://travis-ci.org/barsanuphe/endive)
[![codecov](https://codecov.io/gh/barsanuphe/endive/branch/master/graph/badge.svg)](https://codecov.io/gh/barsanuphe/endive)
[![GPLv3](https://img.shields.io/badge/license-GPLv3-blue.svg)](http://www.gnu.org/licenses/gpl-3.0.en.html)
[![Go Report Card](https://goreportcard.com/badge/github.com/barsanuphe/endive)](https://goreportcard.com/report/github.com/barsanuphe/endive)

## What it is

*endive* is a CLI epub collection manager, allowing you to: 

- **Centralize metadata**: 
*endive* retrieves metadata from the epub file itself, and also from Goodreads, 
to make sure all the relevant and correct information about your books is 
centralized in its database. 
This includes user metadata such as reading status, rating, or review.
- **Carefully import books**: 
importing makes sure that relevant metadata is always collected for any new
book in the collection. The user can define special folders as retail or 
non-retail sources for epubs, or point directly to epub files.
- **Easily manage the collection**:
with a simple but powerful CLI, to search/modify the database with minimal pain. 
- **Organize from metadata**: 
epub files are organized and named automatically from the metadata, according 
to a user-defined pattern. 
- **Take good care of the books**: 
the library can be checked for changes, allowing you to make sure retail epubs 
have not been modified since import. 
Retail epubs always trump non-retail versions. 
- **Easily export to e-readers**:
once your e-reader is locally mounted, *endive* allows exporting a subset of 
the collection with the same queries used for searching. It also keeps track of 
which books have already been sent to your device.

## Current state

Not stable yet. The database format may change, and you should only feed 
*endive* files you've already backed up somewhere.

## Table of Contents

- [Prerequisites](#Prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [Testing](#testing)
- [Third party libraries](#third-party-libraries)

## Prerequisites

Publishers do not seem to lose much sleep over retail epubs having incomplete
metadata.
This means: missing information, wildly different ways to include vital
information such as ISBN number (if it is included at all, or correctly
included).

To get as faithful as possible information about epub files, *endive* relies
on getting information from Goodreads.
Since this means using the Goodreads API, you must have an account and request
[an API Key](https://www.goodreads.com/api/keys)
([terms and conditions](https://www.goodreads.com/api/terms)).

See the configuration instructions to find out what to do with that key.

## Installation

If you have a [working Go installation](https://golang.org/doc/install), just 
run:

    $ go get -u github.com/barsanuphe/endive
    $ go install ...endive

*endive* uses `libgit2` to automatically version your library database, make 
sure it is installed on your system.

Note that *endive* uses `less` do display information, so it will only work
where *less* is available.

For editing long fields (such as book description), *endive* fires up 
`$EDITOR` and falls back to `nano` if nothing is defined, because `nano` is 
everywhere.

*endive* is a *Works On My Machine* project, but it should work on all Linux
distributions, and maybe elsewhere, provided `less` and `$EDITOR` are 
recognized.

## Usage

Import from retail sources:

    $ endive import retail

Import specific non-retail epub:

    $ endive import nonretail book.epub

List epubs in english written by (Charles) Stross:

    $ endive search language:en +author:stross

Note that search is powered by [bleve](https://github.com/blevesearch/bleve),
and therefore uses its
[syntax](http://www.blevesearch.com/docs/Query-String-Query/).

Available fields are: `author`, `title`, `year`, `language`, `tags`, `series`,
`publisher`, `category`, `type`, `genre`, `description`, `exported`, `progress`, 
`review`, and probably a few more.

Same search, ordered by year:

    $ endive search language:en +author:stross --sort year

Results can be sorted by: `id`, `author`, `title`, `year`.

Show info about a book with a specific *ID*:

    $ endive info *ID*

Edit title for book with specific *ID*:

    $ endive set field title "New Title" *ID*
    
Set type for several books: 

    $ endive set field type novel *ID1* *ID2* *ID3*
    
Mark a book as shortlisted for reading:

    $ endive set shortlisted *ID*

Refresh library after configuration changes:

    $ endive collection refresh

List all books:

    $ endive list

List all books that do not have a retail version:

    $ endive list --nonretail

Sorting all books by author:

    $ endive list --sort author

Sorting all books by year, limit to first 10 results:

    $ endive list --sort year --first 10
    
Exporting shortlisted books to a mounted e-reader:

    $ endive export progress:shortlisted

For other commands, see:

    $ endive --help


## Configuration

*endive* uses a YAML configuration file, in the configuration XDG directory
(which should be `/home/user/.config/endive/`):

    # library location and database filename
    library_root: /home/user/endive
    database_filename: endive.json

    # $a author, $y year, $t title, $i isbn, $l language, $r retail status,
	# $c fiction/nonfiction, $g main genre, $p reading progress, $s series
    epub_filename_format: $a/$a ($y) $t

    # list of directories that will be scraped for ebooks,
    # and automatically flagged as retail or non-retail
    nonretail_source:
        - /home/user/nonretail
        - /home/user/nonretail_source2
    retail_source:
        - /home/user/retail

    # see prerequisites
    goodreads_api_key: XXXXXXXXXXXXXX

    # associate main alias to alternative aliases
    # only the main alias will be used by endive
    author_aliases:
        Alexandre Dumas:
            - Alexandre Dumas Père
        China Miéville:
            - China Mieville
            - China Miévile
        Richard K. Morgan:
            - Richard Morgan
        Robert Silverberg:
            - Robert K. Silverberg
        Jared Diamond:
            - Diamond, Jared
    tag_aliases:
        science-fiction:
            - sci-fi
            - sf
            - sciencefiction
    publisher_aliases:
        Tor:
            - Tom Doherty Associates

## Testing

Testing requires the `GR_API_KEY` environment variable to be set with your very
own Goodreads API key.

    $ export GR_API_KEY=XXXXXXXXX
    $ go test ./...


## Third party libraries

|                 | Library       |
| --------------- |:-------------:|
| Epub parser     | [github.com/barsanuphe/epubgo](https://github.com/barsanuphe/epubgo), forked from [github.com/meskio/epubgo](https://github.com/meskio/epubgo)             |
| Search          | [github.com/blevesearch/bleve](https://github.com/blevesearch/bleve) |
| CLI             | [github.com/docopt/docopt-go](https://github.com/docopt/docopt-go)     |
| Color output    | [github.com/ttacon/chalk](https://github.com/ttacon/chalk)           |
| Tables output   | [github.com/barsanuphe/gotabulate](https://github.com/barsanuphe/gotabulate), forked from [github.com/bndr/gotabulate](https://github.com/bndr/gotabulate) |
| XDG directories | [launchpad.net/go-xdg](https://launchpad.net/go-xdg)                 |
| YAML Parser     | [gopkg.in/yaml.v2](https://github.com/go-yaml/yaml)             |
| ISBN validator  | [github.com/moraes/isbn](https://github.com/moraes/isbn)             |
| Spinner         | [github.com/tj/go-spin](https://github.com/tj/go-spin)               |
| Diff            | [github.com/kylelemons/godebug/pretty](https://github.com/kylelemons/godebug/pretty)               |
| Versioning      | [github.com/libgit2/git2go](https://github.com/libgit2/git2go)               |
