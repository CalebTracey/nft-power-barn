package ipfs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

const (
	testAPI       = "6d68a19d-3f75-442f-a593-866763c9d6e2"
	uploadUrl     = "https://api.nftport.xyz/v0/files"
	uploadMetaUrl = "https://api.nftport.xyz/v0/metadata"
)

//go:generate mockgen -destination=mockService.go -package=ipfs . ServiceI
type ServiceI interface {
	UploadImageIpfs(req UploadIpfsImageRequest) (res UploadIpfsImageResponse, err error)
	UploadMetaIpfs(req UploadIpfsMetadataRequest) (res UploadIpfsMetadataResponse, err error)
}

type Service struct {
	Client *http.Client
}

func InitializeService() *Service {
	c := &http.Client{}

	return &Service{
		Client: c,
	}
}

func (s Service) UploadImageIpfs(req UploadIpfsImageRequest) (res UploadIpfsImageResponse, err error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	var fw io.Writer
	if x, ok := req.Reader.(io.Closer); ok {
		defer func(x io.Closer) {
			err := x.Close()
			if err != nil {
				log.Fatalln(err)
			}
		}(x)
	}
	if x, ok := req.Reader.(*os.File); ok {
		if fw, err = w.CreateFormFile(req.Name, x.Name()); err != nil {
			return res, err
		}
	} else {
		if fw, err = w.CreateFormField(req.Name); err != nil {
			return res, err
		}
	}
	if _, err = io.Copy(fw, req.Reader); err != nil {
		return res, err
	}
	r, err := http.NewRequest("POST", uploadUrl, &b)
	if err != nil {
		return res, err
	}
	r.Header.Set("Content-Type", w.FormDataContentType())
	r.Header.Add("Authorization", testAPI)
	response, err := s.Client.Do(r)
	if err != nil {
		return res, err
	}
	err = json.NewDecoder(response.Body).Decode(&res)

	if err != nil {
		return res, err
	}

	log.Println(res)
	return res, nil
}

func (s Service) UploadMetaIpfs(upReq UploadIpfsMetadataRequest) (res UploadIpfsMetadataResponse, err error) {
	//payload := strings.NewReader("{\n  \"name\": \"My Art\",\n  \"description\": \"This is my custom art piece\",\n  \"file_url\": \"https://ipfs.io/ipfs/QmRModSr9gQTSZrbfbLis6mw21HycZyqMA3j8YMRD11nAQ\"\n}")
	buf := new(bytes.Buffer)
	err = json.NewDecoder(buf).Decode(&upReq)
	if err != nil {
		return res, err
	}
	//payload := ioutil.ReadAll(req)
	req, _ := http.NewRequest("POST", uploadMetaUrl, buf)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "6d68a19d-3f75-442f-a593-866763c9d6e2")

	response, _ := s.Client.Do(req)

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Panic(err)
		}
	}(response.Body)

	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println(res)
	fmt.Println(string(body))

	err = json.Unmarshal(body, &res)
	return res, nil
}
