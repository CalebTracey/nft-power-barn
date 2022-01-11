package facade

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gitlab.com/CalebTracey/nft-power-barn/service/ipfs"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type IpfsFacade interface {
	UploadCollectionIpfs() error
}

type IpfsService struct {
	ipfs.Service
}

func NewIpfsService() IpfsService {
	ifpsSvc := ipfs.InitializeService()

	return IpfsService{
		*ifpsSvc,
	}
}

type NftInfo struct {
	Image os.FileInfo
	Json  os.FileInfo
}

type NftFiles struct {
	Image *os.File
	Json  *os.File
}

type ImageInfo struct {
	File os.FileInfo
	Path string
}

func init() {

}
func (s IpfsService) UploadCollectionIpfs() error {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err.Error())
	}
	buildDir = fmt.Sprintf("%v/build", path)
	fileCountInfo, _ := readDirGetFileInfo()

	for i := range fileCountInfo {
		img := s.loadImage(i)
		meta := s.loadMetadata(i)
		buf := new(bytes.Buffer)

		imageIpfs, err := s.UploadImageIpfs(strings.NewReader())
		if err != nil {
			return err
		}
	}
}

//
//func (s IpfsService) loadMetadataInfo() []os.FileInfo{
//	metaPath := fmt.Sprintf("%v/json", buildDir)
//	metaData, err := readDirGetFileInfo(metaPath)
//	if err != nil {
//		panic(err)
//	}
//	return metaData
//}

//func (s IpfsService) loadImageFileData(imgInfo os.FileInfo) *os.File {
//	imgPath := fmt.Sprintf("%v/images", buildDir)
//	imgFile, err := os.OpenFile(fmt.Sprintf("%v/%v", imgPath, imgInfo.Name()), os.O_RDWR|os.O_CREATE, defaultBufSize)
//	if err != nil {
//		panic(err)
//	}
//	return imgFile
//}

func (s IpfsService) loadMetadata(idx int) Metadata {
	metaPath := fmt.Sprintf("%v/json/%v.json", buildDir, idx)
	f, err := os.OpenFile(metaPath, os.O_RDWR|os.O_CREATE, defaultBufSize)
	defer func(img *os.File) {
		e := img.Close()
		if e != nil {
			log.Panic(e)
			return
		}
	}(f)
	if err != nil {
		panic(err)
	}
	bytes, _ := ioutil.ReadAll(f)
	var metaData Metadata
	err = json.Unmarshal(bytes, &metaData)
	if err != nil {
		panic(err)
	}
	return metaData
}

func (s IpfsService) loadImage(idx int) image.Image {
	imgPath := fmt.Sprintf("%v/images/%v.png", buildDir, idx)
	f, err := os.OpenFile(imgPath, os.O_RDWR|os.O_CREATE, defaultBufSize)
	defer func(img *os.File) {
		e := img.Close()
		if e != nil {
			panic(err)
		}
	}(f)
	if err != nil {
		panic(err)
	}
	img, err := png.Decode(f)
	if err != nil {
		panic(err)
	}
	return img
}

func readDirGetFileInfo() ([]os.FileInfo, error) {
	imgPath := fmt.Sprintf("%v/images", buildDir)

	var temp []os.FileInfo
	f, err := os.Open(imgPath)
	if err != nil {
		return nil, err
	}
	fileInfo, err := f.Readdir(-1)
	err = f.Close()
	if err != nil {
		return nil, err
	}
	for _, file := range fileInfo {
		temp = append(temp, file)
	}
	return temp, nil
}
