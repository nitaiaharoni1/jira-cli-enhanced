package view

import (
	"bytes"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/ankitpokhrel/jira-cli/pkg/jira"
	"github.com/ankitpokhrel/jira-cli/pkg/tui"
)

// SavedFilterOption is a functional option to wrap saved filter properties.
type SavedFilterOption func(*SavedFilter)

// SavedFilter is a saved filter view.
type SavedFilter struct {
	data   []*jira.SavedFilter
	writer io.Writer
	buf    *bytes.Buffer
}

// NewSavedFilter initializes a saved filter view.
func NewSavedFilter(data []*jira.SavedFilter, opts ...SavedFilterOption) *SavedFilter {
	sf := SavedFilter{
		data: data,
		buf:  new(bytes.Buffer),
	}
	sf.writer = tabwriter.NewWriter(sf.buf, 0, tabWidth, 1, '\t', 0)

	for _, opt := range opts {
		opt(&sf)
	}
	return &sf
}

// WithSavedFilterWriter sets a writer for the saved filter.
func WithSavedFilterWriter(w io.Writer) SavedFilterOption {
	return func(sf *SavedFilter) {
		sf.writer = w
	}
}

// Render renders the saved filter view.
func (sf SavedFilter) Render() error {
	sf.printHeader()

	for _, d := range sf.data {
		favorite := ""
		if d.Favourite {
			favorite = "★"
		}
		_, _ = fmt.Fprintf(sf.writer, "%s\t%s\t%s\t%s\n", d.ID, prepareTitle(d.Name), prepareTitle(d.JQL), favorite)
	}
	if _, ok := sf.writer.(*tabwriter.Writer); ok {
		err := sf.writer.(*tabwriter.Writer).Flush()
		if err != nil {
			return err
		}
	}

	return tui.PagerOut(sf.buf.String())
}

func (sf SavedFilter) header() []string {
	return []string{
		"ID",
		"NAME",
		"JQL",
		"FAVORITE",
	}
}

func (sf SavedFilter) printHeader() {
	headers := sf.header()
	end := len(headers) - 1
	for i, h := range headers {
		_, _ = fmt.Fprintf(sf.writer, "%s", h)
		if i != end {
			_, _ = fmt.Fprintf(sf.writer, "\t")
		}
	}
	_, _ = fmt.Fprintln(sf.writer)
}

