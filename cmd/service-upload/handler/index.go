package handler

import (
	"html/template"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

type Index interface {
	GetIndex(w http.ResponseWriter, r *http.Request)
}

func NewIndex() Index {
	return &IndexImpl{}
}

type IndexImpl struct {
	Index
}

// Serve Front Page
func (m *IndexImpl) GetIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	data := struct {
		Auth string
	}{
		Auth: os.Getenv("AUTH"),
	}

	tmpl := template.Must(template.ParseFiles("./static/index.html"))
	err := tmpl.Execute(w, data)

	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
