package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

const command = "/usr/bin/convert" // ImageMagick

type Response struct {
	Images []string `json:"base64Images"`
}

type Meta struct {
	Page        int          `json:"page"`
	Resolutions []Resolution `json:"resolutions"`
}

type Resolution struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

func (m *Meta) prepareCommand(pdfFileName string) ([]string, []*os.File, error) {
	var args []string
	var files []*os.File
	args = append(args, fmt.Sprintf("%s[%d]", pdfFileName, m.Page-1))
	args = append(args, "-flatten")
	for i, res := range m.Resolutions {
		args = append(args, "-thumbnail")
		args = append(args, fmt.Sprintf("%dx%d!", res.Width, res.Height))
		if i != len(m.Resolutions)-1 {
			args = append(args, "-write") // intermediate step
		}
		file, err := ioutil.TempFile("", "*.jpg")
		if err != nil {
			return nil, nil, fmt.Errorf("create temp file: %v", err)
		}
		args = append(args, file.Name())
		files = append(files, file)
	}
	return args, files, nil
}

func main() {
	http.HandleFunc("/v1/canary", canary)
	http.HandleFunc("/v1/generatemultiple", generatemultiple)
	port := os.Getenv("DUMBNAILER_PORT")
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}

func canary(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func generatemultiple(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fail(w, http.StatusMethodNotAllowed, "called with method %s", r.Method)
		return
	}
	fileReader, _, err := r.FormFile("file")
	if err != nil {
		fail(w, http.StatusBadRequest, "'file' (PDF) missing: %v", err)
		return
	}
	var meta Meta
	if err := json.Unmarshal([]byte(r.FormValue("meta")), &meta); err != nil {
		fail(w, http.StatusBadRequest, "error decoding meta JSON: %v", err)
		return
	}
	pdfFileName, err := store(fileReader)
	if err != nil {
		fail(w, http.StatusInternalServerError, "storing PDF: %v", err)
		return
	}
	defer os.Remove(pdfFileName)
	args, files, err := meta.prepareCommand(pdfFileName)
	if err != nil {
		fail(w, http.StatusInternalServerError, "preparing dumbnailer command: %v", err)
		return
	}
	cmd := exec.Command(command, args...)
	started := time.Now()
	err = cmd.Run()
	finished := time.Now()
	if err != nil {
		fail(w, http.StatusInternalServerError, "executing '%s %v': %v", command, strings.Join(args, " "), err)
		return
	}
	duration := finished.Sub(started)
	log.Printf("%s %s [%v]", command, strings.Join(args, " "), duration)
	var base64Thumbnails []string
	for _, f := range files {
		content, err := ioutil.ReadAll(f)
		if err != nil {
			fail(w, http.StatusInternalServerError, "reading from file %s: %v", f.Name(), err)
		}
		encoded := base64.StdEncoding.EncodeToString(content)
		base64Thumbnails = append(base64Thumbnails, encoded)
	}
	var response Response
	response.Images = base64Thumbnails
	jsonThubnails, err := json.Marshal(response)
	if err != nil {
		fail(w, http.StatusInternalServerError, "encoding thumbnails to base64 JSON: %v", err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(jsonThubnails)
}

func store(pdfReader io.Reader) (string, error) {
	tmpFile, err := ioutil.TempFile("", "*.pdf")
	if err != nil {
		return "", fmt.Errorf("create temp file: %v", err)
	}
	defer tmpFile.Close()
	if _, err := io.Copy(tmpFile, pdfReader); err != nil {
		return "", fmt.Errorf("store PDF as %s: %v", tmpFile.Name(), err)
	}
	return tmpFile.Name(), nil
}

func fail(w http.ResponseWriter, httpStatus int, format string, params ...interface{}) {
	message := fmt.Sprintf(format, params...)
	w.WriteHeader(httpStatus)
	w.Write([]byte(message))
	log.Println(message)
}
