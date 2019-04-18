package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

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
