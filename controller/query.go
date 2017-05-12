package controller

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"

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

const (
	errorPolicyReturn = iota
	errorPolicyPanic
)

type queryStatus struct {
	Error  error
	Query  string
	Count  int
	Result [][]interface{}
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
	s.Result = make([][]interface{}, 0, len(pages))
	var count = 0
	for _, pnotes := range pagenotes {
		for _, pnote := range pnotes {
			for _, hlt := range pnote.Highlights {
				count++
				record := make([]interface{}, 0, 4)
				pageTitle := pages[hlt.Id][0]
				pageURI := pages[hlt.Id][1]
				var urlLabel string
				if bytes, ok := pageTitle.([]byte); ok {
					urlLabel = strings.TrimSpace(string(bytes))
				}
				if bytes, ok := pageURI.([]byte); ok {
					urlStr := string(bytes)
					if urlLabel == "" {
						urlLabel = extractURLDomain(urlStr)
					}
					if len(urlLabel) > 30 {
						urlLabel = urlLabel[:30]
					}
					pageURI = template.HTML(fmt.Sprintf("<a href='%s'>%s</a>", urlStr, urlLabel))
				}
				record = append(record, count, hlt.Text, pageURI)
				s.Result = append(s.Result, record)
			}
		}
	}
	s.Count = count
	log.Debug(count, " records found for ", q)
	return s
}

func loadTemplate(name, path string, errorPolicy int32) (*template.Template, error) {
	log.Debug("loading template [", name, "] from ", path)
	templateStr, err := ioutil.ReadFile(path)
	if err != nil {
		log.Error("error loading tempalte", err)
		if errorPolicy == errorPolicyPanic {
			panic(err)
		}
		return nil, err
	}
	t, err := template.New(name).Parse(string(templateStr))
	if err != nil {
		log.Error("error loading tempalte", err)
		if errorPolicy == errorPolicyPanic {
			panic(err)
		}
	}
	return t, err
}

func loadQTemplate() (err error) {
	if conf.IsDebug() || qTemplate == nil {
		qTemplatePath := util.GetAbsRunDirPath() + "/view/q.html.template"
		qTemplate, err = loadTemplate("q", qTemplatePath, errorPolicyReturn)
	}
	return err
}

func init() {

}
