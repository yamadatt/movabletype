package movabletype_test

import (
	"bytes"
	"io"
	"reflect"
	"strings"
	"testing"
	"time"

	. "github.com/yamadatt/movabletype"
)

func TestParse(t *testing.T) {
	buf := bytes.NewBufferString(`AUTHOR: catatsuy
TITLE: ポエム
BASENAME: poem
STATUS: Publish
ALLOW COMMENTS: 1
ALLOW PINGS: 1
CONVERT BREAKS: 0
DATE: 04/22/2017 20:41:58
PRIMARY CATEGORY: ブログ
CATEGORY: ポエム
CATEGORY: 技術系
-----
BODY:
<p>body</p>
<p>bodybody</p>
<p>bodybodybody</p>
-----
EXTENDED BODY:
<p>extended body</p>
<p>extended body body</p>
<p>extended body body body</p>
-----
EXCERPT:
ここに概要が表示されます。
-----
--------
AUTHOR: catatsuy
TITLE: 風邪で声を失った話
BASENAME: 2017/04/09/194939
STATUS: Publish
ALLOW COMMENTS: 1
CONVERT BREAKS: 0
DATE: 04/09/2017 07:49:39 PM
CATEGORY: 日常
-----
BODY:
<p>bodybodybody</p>
-----
EXTENDED BODY:
<p>extended body body body</p>
-----
EXCERPT:
ここに概要が表示されます。
-----
KEYWORDS:
keywords
-----
COMMENT:
AUTHOR: 紗菜
EMAIL: 
-----

--------
`)

	expected := []*Entry{
		&Entry{
			Author:          "catatsuy",
			Title:           "ポエム",
			Basename:        "poem",
			Status:          "Publish",
			AllowComments:   1,
			AllowPings:      1,
			ConvertBreaks:   "0",
			Date:            time.Date(2017, time.April, 22, 20, 41, 58, 0, time.UTC),
			PrimaryCategory: "ブログ",
			Category:        []string{"ポエム", "技術系"},
			Body:            "<p>body</p>\n<p>bodybody</p>\n<p>bodybodybody</p>\n",
			ExtendedBody:    "<p>extended body</p>\n<p>extended body body</p>\n<p>extended body body body</p>\n",
			Excerpt:         "ここに概要が表示されます。\n",
		},
		&Entry{
			Author:        "catatsuy",
			Title:         "風邪で声を失った話",
			Basename:      "2017/04/09/194939",
			Status:        "Publish",
			AllowComments: 1,
			AllowPings:    -1,
			ConvertBreaks: "0",
			Date:          time.Date(2017, time.April, 9, 19, 49, 39, 0, time.UTC),
			Category:      []string{"日常"},
			Body:          "<p>bodybodybody</p>\n",
			ExtendedBody:  "<p>extended body body body</p>\n",
			Excerpt:       "ここに概要が表示されます。\n",
			Keywords:      "keywords\n",
			Comment:       "AUTHOR: 紗菜\nEMAIL: \n",
		},
	}

	mts, err := Parse(buf)

	if err != nil {
		t.Fatalf("got error %q", err)
	}

	if !reflect.DeepEqual(mts, expected) {
		t.Errorf("Error parsing, expected %#v; got %#v", expected, mts)
	}
}

func TestParseStatusNotAllowed(t *testing.T) {
	buf := bytes.NewBufferString(`STATUS: Published`)

	_, err := Parse(buf)

	if err == nil || err.Error() != "STATUS column is allowed only Draft or Publish or Future. Got Published" {
		t.Errorf("Error parsing, got %q", err)
	}
}

func TestParseDate(t *testing.T) {
	var featuretests = []struct {
		buf io.Reader
		t   time.Time
	}{
		{
			bytes.NewBufferString("DATE: 04/22/2017 08:41:58 PM\n--------\n"),
			time.Date(2017, time.April, 22, 20, 41, 58, 0, time.UTC),
		},
		{
			bytes.NewBufferString("DATE: 04/22/2017 08:41:58 AM\n--------\n"),
			time.Date(2017, time.April, 22, 8, 41, 58, 0, time.UTC),
		},
		{
			bytes.NewBufferString("DATE: 04/22/2017 20:41:58\n--------\n"),
			time.Date(2017, time.April, 22, 20, 41, 58, 0, time.UTC),
		},
	}

	for _, ft := range featuretests {
		mts, err := Parse(ft.buf)

		if err != nil {
			t.Fatalf("got error %q", err)
		}

		if mts[0].Date != ft.t {
			t.Errorf("m.Date got %v; want %v", mts[0].Date, ft.t)
		}
	}
}

func TestNewMT(t *testing.T) {
	m := NewEntry()

	if m.AllowComments != DefaultAllowComments {
		t.Errorf("By default, AllowComments is %d, got %d", DefaultAllowComments, m.AllowComments)
	}

	if m.AllowPings != DefaultAllowPings {
		t.Errorf("By default, AllowComments is %d, got %d", DefaultAllowPings, m.AllowPings)
	}
}

// TestParseEmptyInput tests how the parser handles an empty input
func TestParseEmptyInput(t *testing.T) {
	buf := bytes.NewBufferString("")

	mts, err := Parse(buf)

	if err != nil {
		t.Fatalf("Parse should not return error on empty input, got: %v", err)
	}

	if len(mts) != 0 {
		t.Errorf("Expected empty slice for empty input, got %d entries", len(mts))
	}
}

// TestParseSingleEntry tests parsing a single entry without optional fields
func TestParseSingleEntry(t *testing.T) {
	buf := bytes.NewBufferString(`AUTHOR: testauthor
TITLE: Test Title
STATUS: Draft
DATE: 01/02/2023 15:30:45
-----
BODY:
Simple body content
-----
--------
`)

	expected := []*Entry{
		{
			Author:        "testauthor",
			Title:         "Test Title",
			Status:        "Draft",
			AllowComments: DefaultAllowComments,
			AllowPings:    DefaultAllowPings,
			Date:          time.Date(2023, time.January, 2, 15, 30, 45, 0, time.UTC),
			Body:          "Simple body content\n",
		},
	}

	mts, err := Parse(buf)

	if err != nil {
		t.Fatalf("got error %q", err)
	}

	if !reflect.DeepEqual(mts, expected) {
		t.Errorf("Error parsing simple entry, expected %#v; got %#v", expected, mts)
	}
}

// TestParseInvalidAllowComments tests error handling for invalid ALLOW COMMENTS values
func TestParseInvalidAllowComments(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "non-numeric value",
			input:    "ALLOW COMMENTS: yes",
			expected: "ALLOW COMMENTS column is allowed only 0 or 1:",
		},
		{
			name:     "out of range value",
			input:    "ALLOW COMMENTS: 2",
			expected: "ALLOW COMMENTS column is allowed only 0 or 1. Got 2",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buf := bytes.NewBufferString(test.input)
			_, err := Parse(buf)

			if err == nil {
				t.Fatalf("Expected error for invalid ALLOW COMMENTS, got nil")
			}

			if !strings.Contains(err.Error(), test.expected) {
				t.Errorf("Error mismatch, expected: %q, got: %q", test.expected, err.Error())
			}
		})
	}
}

// TestParseInvalidAllowPings tests error handling for invalid ALLOW PINGS values
func TestParseInvalidAllowPings(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "non-numeric value",
			input:    "ALLOW PINGS: true",
			expected: "ALLOW PINGS column is allowed only 0 or 1:",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buf := bytes.NewBufferString(test.input)
			_, err := Parse(buf)

			if err == nil {
				t.Fatalf("Expected error for invalid ALLOW PINGS, got nil")
			}

			if !strings.Contains(err.Error(), test.expected) {
				t.Errorf("Error mismatch, expected: %q, got: %q", test.expected, err.Error())
			}
		})
	}
}

// TestParseInvalidDate tests error handling for invalid DATE formats
func TestParseInvalidDate(t *testing.T) {
	buf := bytes.NewBufferString("DATE: 2023-01-02 15:30:45")

	_, err := Parse(buf)

	if err == nil {
		t.Fatal("Expected error for invalid date format, got nil")
	}

	if !strings.Contains(err.Error(), "Parsing error on DATE column") {
		t.Errorf("Error mismatch, expected message about DATE parsing, got: %s", err.Error())
	}
}

// TestMultipleEntries tests parsing multiple entries in a single input
func TestMultipleEntries(t *testing.T) {
	buf := bytes.NewBufferString(`AUTHOR: author1
TITLE: Title 1
STATUS: Publish
DATE: 01/01/2023 12:00:00
-----
BODY:
Body 1
-----
--------
AUTHOR: author2
TITLE: Title 2
STATUS: Draft
DATE: 01/02/2023 12:00:00
-----
BODY:
Body 2
-----
--------
AUTHOR: author3
TITLE: Title 3
STATUS: Future
DATE: 01/03/2023 12:00:00
-----
BODY:
Body 3
-----
--------
`)

	mts, err := Parse(buf)

	if err != nil {
		t.Fatalf("got error %q", err)
	}

	if len(mts) != 3 {
		t.Fatalf("Expected 3 entries, got %d", len(mts))
	}

	expected := []struct {
		author string
		title  string
		status string
		date   time.Time
	}{
		{
			author: "author1",
			title:  "Title 1",
			status: "Publish",
			date:   time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			author: "author2",
			title:  "Title 2",
			status: "Draft",
			date:   time.Date(2023, time.January, 2, 12, 0, 0, 0, time.UTC),
		},
		{
			author: "author3",
			title:  "Title 3",
			status: "Future",
			date:   time.Date(2023, time.January, 3, 12, 0, 0, 0, time.UTC),
		},
	}

	for i, entry := range expected {
		if mts[i].Author != entry.author {
			t.Errorf("Entry %d, author mismatch, expected %q, got %q", i, entry.author, mts[i].Author)
		}
		if mts[i].Title != entry.title {
			t.Errorf("Entry %d, title mismatch, expected %q, got %q", i, entry.title, mts[i].Title)
		}
		if mts[i].Status != entry.status {
			t.Errorf("Entry %d, status mismatch, expected %q, got %q", i, entry.status, mts[i].Status)
		}
		if !mts[i].Date.Equal(entry.date) {
			t.Errorf("Entry %d, date mismatch, expected %v, got %v", i, entry.date, mts[i].Date)
		}
	}
}
