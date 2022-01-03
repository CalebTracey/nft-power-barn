package facade

import (
	"image"
)

type ImageResponse struct {
	Layer LayerToDnaResults
	Image image.Image
}

type Metadata struct {
	Name         string       `json:"Name"`
	Description  string       `json:"Description"`
	FileUrl      string       `json:"FileUrl"`
	CustomFields CustomFields `json:"CustomFields"`
	Creator      string       `json:"Creator"`
	Attributes   []Attributes `json:"Attributes"`
}

type CustomFields struct {
	DNA      string `json:"DNA"`
	Edition  int    `json:"Edition"`
	Date     string `json:"Date"`
	Compiler string `json:"Compiler"`
}

type Attributes struct {
	TraitType string `json:"Trait type"`
	Value     string `json:"Value"`
}

type LayerElement struct {
	Id        int
	Elements  []Element
	Name      string
	Blend     string
	Opacity   int
	BypassDNA bool
}

type Element struct {
	Id       int
	Name     string
	FileName string
	Weight   int
	Path     string
}

type LayerToDnaResults struct {
	Name            string
	Blend           string
	Opacity         int
	SelectedElement Element
}
