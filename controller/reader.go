package controller

import (
	"html/template"
	"net/http"

	"github.com/go-playground/log"
	"github.com/lpimem/hlcsrv/hlccookie"
	"github.com/lpimem/hlcsrv/storage"
)

type reader struct{}

type readingIndex struct {
	Error    error
	Readings []*storage.ReadingProgress
	Books    []*storage.Book
}

// Reader controller
var (
	Reader reader
	tmpl   *template.Template
)

func loadReaderTemplate() (err error) {
	tmpl, err = loadTemplate("reader", "view/reader.html.tmpl", errorPolicyReturn)
	return err
}

func queryReadings(uid storage.UserID) readingIndex {
	var idx readingIndex
	if idx.Books, idx.Error = storage.Reading.QueryBooks("%"); idx.Error != nil {
		return idx
	}
	idx.Readings, idx.Error = storage.Reading.QueryAllProgress(uid)
	return idx
}

func (r *reader) Index(w http.ResponseWriter, req *http.Request) {
	if !requireAuth(w, req) {
		return
	}

	var (
		uid storage.UserID
		err error
	)

	if err = loadReaderTemplate(); err != nil {
		log.Error("Cannot load template: ", err)
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}

	if uid, err = hlccookie.GetRequestUID(req); err != nil {
		log.Error("Cannot get request id ", err)
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}
	indexStatus := queryReadings(uid)
	tmpl.Execute(w, indexStatus)
}
