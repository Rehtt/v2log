package main

import "os"

type CVS struct {
	title    []string
	text     [][]string
	fileName string
	flash    bool
}

func NewCSV(fileName string) (*CVS, error) {
	_, err := os.OpenFile(fileName, os.O_RDONLY, 644)
	if err != nil {
		return nil, err
	}
	return &CVS{
		flash: false,
	}, err
}

func (csv *CVS) Title(title []string) {
	csv.flash = false
	csv.title = title
}
func (csv *CVS) Text(text [][]string) {
	csv.text = text
}

func (csv *CVS) Flash() error {
	if !csv.flash {
		csv.flash = true

	}
}
