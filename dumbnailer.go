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
	"sort"
	"strings"
	"time"

	"dumbnailer/params"
)

var command = "/usr/bin/convert" // ImageMagick

type Response struct {
	Images []string `json:"base64Images"`
}

func main() {
	http.HandleFunc("/v1/canary", canary)
	http.HandleFunc("/v1/generatemultiple", generatemultiple)
	imgmgck := os.Getenv("IMAGE_MAGICK")
	if imgmgck != "" {
		command = imgmgck
	}
	port := os.Getenv("PORT")
	log.Printf("dumbnailer listening on port %s\n", port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}

func canary(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok\n"))
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
	var meta params.Meta
	if err := json.Unmarshal([]byte(r.FormValue("meta")), &meta); err != nil {
		fail(w, http.StatusBadRequest, "error decoding meta JSON: %v", err)
		return
	}
	pdfFileName, err := store(fileReader)
	if err != nil {
		fail(w, http.StatusInternalServerError, "storing PDF: %v", err)
		return
	}
	defer func() {
		if err := os.Remove(pdfFileName); err != nil {
			log.Printf("deleting %s: %v", pdfFileName, err)
		}
	}()
	args, files, err := meta.PrepareCommand(pdfFileName)
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
	resp, err := base64Response(files)
	if err != nil {
		fail(w, http.StatusInternalServerError, "files to base64 response: %v", err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(resp)
}

type base64Junks []string

func (b base64Junks) Len() int           { return len(b) }
func (b base64Junks) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b base64Junks) Less(i, j int) bool { return len(b[i]) < len(b[j]) }

func base64Response(files []*os.File) ([]byte, error) {
	var base64Thumbnails base64Junks
	for _, f := range files {
		content, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, fmt.Errorf("reading from file %s: %v", f.Name(), err)
		}
		if err := os.Remove(f.Name()); err != nil {
			log.Printf("deleting %s: %v", f.Name(), err)
		}
		encoded := base64.StdEncoding.EncodeToString(content)
		base64Thumbnails = append(base64Thumbnails, encoded)
	}

	// ascending order of thumbnail size (resolution)
	sort.Sort(base64Thumbnails)

	var response Response
	response.Images = base64Thumbnails
	base64Response, err := json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("encoding thumbnails to base64 JSON: %v", err)
	}
	return base64Response, nil
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
