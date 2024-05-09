package controller

import (
	"fmt"
	"io"
	configpackage "miniwebserver/config"
	loggerpackage "miniwebserver/logger"
	"net/http"
	"os"
	"strings"
)

type Controller interface {
	Root(w http.ResponseWriter, r *http.Request)
	Options(w http.ResponseWriter, r *http.Request)
	About(w http.ResponseWriter, r *http.Request)
	ServeFile(w http.ResponseWriter, r *http.Request)
	UploadPage(w http.ResponseWriter, r *http.Request)
	UploadFile(w http.ResponseWriter, r *http.Request)
}

type DefaultController struct {
	config   configpackage.Configuration
	logger   loggerpackage.Logger
	response HTTPResponse
}

func NewHTTPController(config configpackage.Configuration, logger loggerpackage.Logger, res HTTPResponse) (defaultController *DefaultController) {

	defaultController = &DefaultController{
		config:   config,
		logger:   logger,
		response: res,
	}

	return defaultController
}

func (u *DefaultController) Root(w http.ResponseWriter, r *http.Request) {
	defer u.response.Panic(w, r)

	u.response.DefaultText(w, http.StatusOK, true, *u.config.GetConfiguration().Version)
	return
}

func (u *DefaultController) Options(w http.ResponseWriter, r *http.Request) {
	defer u.response.Panic(w, r)

	return
}

func (u *DefaultController) About(w http.ResponseWriter, r *http.Request) {
	defer u.response.Panic(w, r)
	u.response.Default(w, http.StatusOK, true, *u.config.GetConfiguration().Version)
	return
}

func (u *DefaultController) ServeFile(w http.ResponseWriter, r *http.Request) {
	defer u.response.Panic(w, r)

	// fmt.Println(r.URL.Query())

	var path = strings.Split(r.URL.Path, "/")[2]
	var file = strings.Split(r.URL.Path, "/")[3]

	fmt.Printf("%#v\n", r.URL.Path)

	// var identifier = r.URL.Query().Get("identifier")
	_filesep := *u.config.GetConfiguration().FileSep
	var servePath string

	fmt.Printf("path: %v, file: %v\n", path, file)

	switch path {
	case "book":
		servePath = *u.config.GetConfiguration().PathFile + _filesep + file
	case "media":
		servePath = *u.config.GetConfiguration().PathMedia + _filesep + file
	default:
		u.response.Default(w, http.StatusNotFound, false, "file tidak ditemukan")
		return
	}

	fmt.Println(servePath)

	download := r.URL.Query().Get("download")
	if download != "" {
		w.Header().Add("Content-Disposition", "attachment")
	}

	// filepath := servePath + _filesep + file

	http.ServeFile(w, r, servePath)
}

func (u *DefaultController) UploadPage(w http.ResponseWriter, r *http.Request) {
	defer u.response.Panic(w, r)

	http.ServeFile(w, r, *u.config.GetConfiguration().PathTemp+*u.config.GetConfiguration().FileSep+`html`+*u.config.GetConfiguration().FileSep+`upload.html`)
}

func (u *DefaultController) UploadFileSingle(w http.ResponseWriter, r *http.Request) {
	defer u.response.Panic(w, r)

	// Parse the form data, including the file uploaded
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve file from form data
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error Retrieving the File", http.StatusBadRequest)
		return
	}
	defer file.Close()

	filename := *u.config.GetConfiguration().PathTemp + *u.config.GetConfiguration().FileSep + `html` + *u.config.GetConfiguration().FileSep + handler.Filename

	// Create a new file in the server's upload directory
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// Copy the file to the destination path
	_, err = io.Copy(f, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File uploaded successfully: %s", handler.Filename)
}

func (u *DefaultController) UploadFile(w http.ResponseWriter, r *http.Request) {
	defer u.response.Panic(w, r)

	fmt.Printf("multipart : %#v\n", r.MultipartForm)

	err := r.ParseMultipartForm(1 << 30) // 10 MB limit
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve files from form data
	files := r.MultipartForm.File["files[]"]

	fmt.Println(files)

	// Iterate over the uploaded files
	for _, fileHeader := range files {
		// Open the uploaded file
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Error Retrieving the File", http.StatusBadRequest)
			return
		}
		defer file.Close()

		filename := *u.config.GetConfiguration().PathTemp + *u.config.GetConfiguration().FileSep + `html` + *u.config.GetConfiguration().FileSep + fileHeader.Filename

		// Create a new file in the server's upload directory
		f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		// Copy the file to the destination path
		_, err = io.Copy(f, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
