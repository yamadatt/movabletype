// Package movabletype provides parsing functionality for "The Movable Type Import / Export Format".
// It can parse Movable Type's text-based export format into Go structures.
package movabletype

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

// Default
const (
	// DefaultAllowComments is the default value for the AllowComments field (-1).
	DefaultAllowComments = -1

	// DefaultAllowPings is the default value for the AllowPings field (-1).
	DefaultAllowPings = -1
)

// Entry represents a single blog post entry in the Movable Type Import Format.
// Each entry contains metadata and content fields as defined in the MT format specification.
type Entry struct {
	Author   string // Author name
	Title    string // Entry title
	Basename string // URL basename
	Status   string // Publication status: "Draft", "Publish", or "Future"

	// AllowComments indicates whether comments are allowed for this entry (0 or 1).
	// If not initialized, it defaults to DefaultAllowComments.
	AllowComments int

	// AllowPings indicates whether trackbacks/pingbacks are allowed for this entry (0 or 1).
	// If not initialized, it defaults to DefaultAllowPings.
	AllowPings int

	ConvertBreaks   string    // Convert line breaks setting
	Date            time.Time // Publication date and time
	PrimaryCategory string    // Primary category name
	Category        []string  // List of categories
	Image           string    // Featured image path
	Body            string    // Main content
	ExtendedBody    string    // Extended/additional content
	Excerpt         string    // Entry excerpt/summary
	Keywords        string    // SEO keywords
	Comment         string    // Comments on the entry
}

// NewEntry creates a new Entry with default values.
func NewEntry() *Entry {
	return &Entry{
		AllowComments: DefaultAllowComments,
		AllowPings:    DefaultAllowPings,
	}
}

// Parse reads Movable Type formatted data from an io.Reader and returns a slice of Entry structures.
// It returns an error if the input is malformed or if required fields have invalid values.
//
// Example usage:
//
//	entries, err := movabletype.Parse(os.Stdin)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, entry := range entries {
//	    fmt.Println(entry.Title)
//	}
func Parse(r io.Reader) ([]*Entry, error) {
	mts := []*Entry{}

	scanner := bufio.NewScanner(r)

	var err error

	m := NewEntry()

	for scanner.Scan() {
		ss := strings.Split(scanner.Text(), ": ")

		if len(ss) <= 1 {
			value := ss[0]

			if value == "--------" {
				mts = append(mts, m)
				m = NewEntry()
				continue
			}

			if value == "-----" {
				continue
			}

			switch value {
			case "BODY:":
				for scanner.Scan() {
					line := scanner.Text()

					if line == "-----" {
						break
					}

					m.Body += line + "\n"
				}
				break
			case "EXTENDED BODY:":
				for scanner.Scan() {
					line := scanner.Text()

					if line == "-----" {
						break
					}

					m.ExtendedBody += line + "\n"
				}
				break
			case "EXCERPT:":
				for scanner.Scan() {
					line := scanner.Text()

					if line == "-----" {
						break
					}

					m.Excerpt += line + "\n"
				}
				break
			case "KEYWORDS:":
				for scanner.Scan() {
					line := scanner.Text()

					if line == "-----" {
						break
					}

					m.Keywords += line + "\n"
				}
				break
			case "COMMENT:":
				for scanner.Scan() {
					line := scanner.Text()
					// fmt.Println(line)
					if line == "-----" {
						break
					}

					m.Comment += line + "\n"
				}
				break
			}

			continue
		}

		key, value := ss[0], ss[1]
		if value != "COMMENT:" {
			switch key {
			case "AUTHOR":
				m.Author = value
				break
			case "TITLE":
				m.Title = value
				break
			case "BASENAME":
				m.Basename = value
				break
			case "STATUS":
				if value == "Draft" || value == "Publish" || value == "Future" {
					m.Status = value
				} else {
					return nil, fmt.Errorf("STATUS column is allowed only Draft or Publish or Future. Got %s", value)
				}
				break
			case "ALLOW COMMENTS":
				m.AllowComments, err = strconv.Atoi(value)
				if err != nil {
					return nil, fmt.Errorf("ALLOW COMMENTS column is allowed only 0 or 1: %w", err)
				}
				if m.AllowComments != 0 && m.AllowComments != 1 {
					return nil, fmt.Errorf("ALLOW COMMENTS column is allowed only 0 or 1. Got %d", m.AllowComments)
				}
				break
			case "ALLOW PINGS":
				m.AllowPings, err = strconv.Atoi(value)
				if err != nil {
					return nil, fmt.Errorf("ALLOW PINGS column is allowed only 0 or 1: %w", err)
				}
				if m.AllowComments != 0 && m.AllowComments != 1 {
					return nil, fmt.Errorf("ALLOW PINGS column is allowed only 0 or 1. Got %d", m.AllowPings)
				}
				break
			case "CONVERT BREAKS":
				m.ConvertBreaks = value
				break
			case "DATE":
				if strings.HasSuffix(value, "AM") || strings.HasSuffix(value, "PM") {
					m.Date, err = time.Parse("01/02/2006 03:04:05 PM", value)
				} else {
					m.Date, err = time.Parse("01/02/2006 15:04:05", value)
				}
				if err != nil {
					return nil, fmt.Errorf("Parsing error on DATE column: %w", err)
				}
				break
			case "PRIMARY CATEGORY":
				m.PrimaryCategory = value
				break
			case "CATEGORY":
				m.Category = append(m.Category, value)
				break
			case "IMAGE":
				m.Image = value
				break
			}
		}
	}

	return mts, nil
}
