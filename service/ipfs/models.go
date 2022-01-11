package ipfs

import (
	"io"
)

type UploadIpfsImageRequest struct {
	Name   string
	Reader io.Reader
}

type UploadIpfsImageResponse struct {
	Response    string  `json:"response"`
	IpfsUrl     string  `json:"ipfs_url"`
	FileName    string  `json:"file_name"`
	ContentType string  `json:"content_type"`
	FileSize    int     `json:"file_size"`
	FileSizeMb  float64 `json:"file_size_mb"`
}

type UploadIpfsMetadataResponse struct {
	Response    string `json:"response"`
	MetadataUri string `json:"metadata_uri"`
	Name        string `json:"name"`
	Description string `json:"description"`
	FileUrl     string `json:"file_url"`
}

type UploadIpfsMetadataRequest struct {
	Name         string       `json:"Name"`
	Description  string       `json:"Description"`
	FileUrl      string       `json:"FileUrl"`
	CustomFields CustomFields `json:"CustomFields"`
	Attributes   []Attributes `json:"Attributes"`
}

type CustomFields struct {
	DNA        string `json:"Genetic Code"`
	Generation int    `json:"Generation"`
	Edition    int    `json:"Mint Number"`
	Date       string `json:"Date of Production"`
	Compiler   string `json:"Compiled with"`
}

type Attributes struct {
	TraitType string `json:"Trait type"`
	Value     string `json:"Value"`
}
