package handler

import (
	"context"
	"encoding/base64"
	"fmt"
	"html/template"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"simple-upload/pkg/mime"
	"strconv"

	log "github.com/sirupsen/logrus"

	"simple-upload/internal/model"
	mongoPkg "simple-upload/pkg/mongo"
)

type Upload interface {
	UploadFile(w http.ResponseWriter, r *http.Request)
}

func NewUpload() Upload {
	return &UploadImpl{}
}

type UploadImpl struct {
	Upload
}

func (h *UploadImpl) UploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check Auth Code
	if auth := r.FormValue("auth"); auth != os.Getenv("AUTH") {
		http.Error(w, "Unauthenticated", http.StatusForbidden)
		return
	}

	file, header, err := r.FormFile("data")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	err = h.checkFileHeader(header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	// create temporary file
	ext := filepath.Ext(header.Filename)
	nameFormat := fmt.Sprintf("img-*%s", ext)
	tempFile, err := ioutil.TempFile("temp-images", nameFormat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// write to temporary file
	tempFile.Write(fileBytes)

	fInfo, _ := tempFile.Stat()

	mdl := model.NewUploadData()
	mdl.Name = fInfo.Name()
	mdl.OriginalName = header.Filename
	mdl.Ext = ext
	mdl.FileSize = header.Size
	mdl.Mime = mime.GetMimeByExtension(mdl.Ext)

	h.insertUploadData(r.Context(), mdl)

	h.uploadOkWriter(w, fileBytes, mdl.Mime)
}

func (h *UploadImpl) checkFileHeader(header *multipart.FileHeader) error {
	maxImageSize, err := strconv.Atoi(os.Getenv("MAX_IMAGE_SIZE_IN_MB"))
	if err != nil {
		maxImageSize = 8
	}

	if header.Size > int64(maxImageSize<<20) {
		return fmt.Errorf("file size exceeded")
	}

	fileType := mime.GetMimeByExtension(filepath.Ext(header.Filename))
	errMustImage := fmt.Errorf("file must be images")
	if fileType == "" {
		return errMustImage
	}

	ct := header.Header["Content-Type"]
	if len(ct) > 0 && ct[0] != fileType {
		return errMustImage
	}

	return nil
}

func (h *UploadImpl) uploadOkWriter(w http.ResponseWriter, bytes []byte, contentType string) {

	str := base64.StdEncoding.EncodeToString(bytes)

	data := struct {
		Image string
		Mime  string
	}{
		Image: str,
		Mime:  contentType,
	}

	tmpl := template.Must(template.ParseFiles("./static/upload_ok.html"))
	err := tmpl.Execute(w, data)

	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *UploadImpl) insertUploadData(ctx context.Context, upl *model.UploadData) error {
	db, err := mongoPkg.NewMongoDB()
	if db != nil {
		defer db.Close(ctx)
	}
	if err != nil {
		log.Error(err)
		return err
	}

	err = db.Insert(ctx, upl.Table(), upl)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
