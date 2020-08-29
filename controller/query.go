package controller

import (
	"errors"
	"html/template"
	"net/http"
	"net/url"
	"fmt"

	"strings"

	"github.com/go-playground/log"
	"github.com/lpimem/hlcsrv/conf"
	"github.com/lpimem/hlcsrv/hlccookie"
	"github.com/lpimem/hlcsrv/storage"
	"github.com/lpimem/hlcsrv/util"
)

var (
	qTemplate *template.Template
)

type QueryRecord struct {
	Count int; 
	Text string;
	UrlLabel string;
	Url template.URL;
}

type queryStatus struct {
	Error  error
	Query  string
	Count  int
	Result []QueryRecord
}

// HandleQuery process query requrests
func HandleQuery(w http.ResponseWriter, req *http.Request) {
	if !requireAuth(w, req) {
		log.Warn("Not authorized...")
		return
	}
	if req.Method == "GET" {
		getQueryPage(w, req)
	} else if req.Method == "POST" {
		postQuery(w, req)
	} else {
		http.Error(w, "invalid request", http.StatusBadRequest)
	}
}

func getQueryPage(w http.ResponseWriter, req *http.Request) {
	err := loadQTemplate()
	if err != nil {
		log.Error("cannot load q page template", err)
	}
	qTemplate.Execute(w, queryStatus{err, "", 0, nil})
}

func postQuery(w http.ResponseWriter, req *http.Request) {
	err := loadQTemplate()
	if err != nil {
		qTemplate.Execute(w, queryStatus{err, "", 0, nil})
		return
	}
	prefix := req.FormValue("query")
	uid, err := hlccookie.GetRequestUID(req)
	if err != nil {
		log.Error("cannot extract request uid from cookie: ", err)
		qTemplate.Execute(w, queryStatus{errors.New("cannot perform query now"), prefix, 0, nil})
		return
	}
	pagenotes, pages, err := storage.QueryUserPagenoteByURI(uid, prefix)
	if err != nil {
		qTemplate.Execute(w, queryStatus{err, prefix, 0, nil})
		return
	}
	status := buildQueryStatus(prefix, pagenotes, pages)
	qTemplate.Execute(w, status)
}

func extractURLDomain(uri string) string {
	u, err := url.Parse(uri)
	if err != nil {
		log.Warn("Error parsing url", uri, err)
		return uri
	}
	return u.Hostname()
}

func buildQueryStatus(q string, pagenotes storage.PagenoteDict, pages storage.PagenoteAddon) *queryStatus {
	var s = &queryStatus{}
	s.Query = q
	s.Result = []QueryRecord{}
	var count = 0
	for _, pnotes := range pagenotes {
		for _, pnote := range pnotes {
			for _, hlt := range pnote.Highlights {
				count++
				pageTitle := pages[hlt.Id][0]
				pageURI := pages[hlt.Id][1]
				var urlLabel string
				var urlStr string
				if bytes, ok := pageTitle.([]byte); ok {
					urlLabel = strings.TrimSpace(string(bytes))
				} else if astring, ok := pageTitle.(string); ok {
					urlLabel = astring
				}
				if bytes, ok := pageURI.([]byte); ok {
					urlStr = string(bytes)
				} else {
					if astring, ok := pageURI.(string); !ok{
						defer log.Error(fmt.Sprintf("Cannot parse URL: %s", pageURI))
						continue
				} else {
						urlStr = astring
					}
				}
				if urlLabel == "" {
					urlLabel = extractURLDomain(urlStr)
				}
				if len(urlLabel) > 30 {
					urlLabel = urlLabel[:30]
				}
				url := template.URL(urlStr)
				s.Result = append(s.Result, 
					QueryRecord{count, hlt.Text, urlLabel, url})
			}
		}
	}
	s.Count = count
	log.Debug(count, " records found for '", q, "'")
	return s
}

func loadQTemplate() (err error) {
	if conf.IsDebug() || qTemplate == nil {
		qTemplatePath := util.GetHLCRoot() + "/view/q.html.template"
		qTemplate, err = loadTemplate("q", qTemplatePath, errorPolicyReturn)
	}
	return err
}

func init() {

}
