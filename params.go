package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
)

type Meta struct {
	Page        uint       `json:"page"`
	Resolutions Dimensions `json:"resolutions"`
}

type Dimensions []Resolution

type Resolution struct {
	Width  uint `json:"width"`
	Height uint `json:"height"`
}

func (r Resolution) Size() uint {
	return r.Width * r.Height
}

func (d Dimensions) Len() int           { return len(d) }
func (d Dimensions) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d Dimensions) Less(i, j int) bool { return d[i].Size() < d[j].Size() }

func (m *Meta) prepareCommand(pdfFileName string) ([]string, []*os.File, error) {
	var args []string
	var files []*os.File
	args = append(args, fmt.Sprintf("%s[%d]", pdfFileName, m.Page-1))
	args = append(args, "-flatten")

	// convert from original to subsequent smaller resolutions
	sort.Sort(sort.Reverse(m.Resolutions))

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
