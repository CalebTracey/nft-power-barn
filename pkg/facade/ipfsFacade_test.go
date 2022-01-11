package facade

import (
	"gitlab.com/CalebTracey/nft-power-barn/service/ipfs"
	"image"
	"os"
	"reflect"
	"testing"
)

func TestIpfsService_UploadCollectionIpfs(t *testing.T) {
	type fields struct {
		Service ipfs.Service
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := IpfsService{
				Service: tt.fields.Service,
			}
			if err := s.UploadCollectionIpfs(); (err != nil) != tt.wantErr {
				t.Errorf("UploadCollectionIpfs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIpfsService_loadImage(t *testing.T) {
	type fields struct {
		Service ipfs.Service
	}
	type args struct {
		idx int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    image.Image
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := IpfsService{
				Service: tt.fields.Service,
			}
			got, err := s.loadImage(tt.args.idx)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadImage() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIpfsService_loadMetadata(t *testing.T) {
	type fields struct {
		Service ipfs.Service
	}
	tests := []struct {
		name   string
		fields fields
		want   Metadata
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := IpfsService{
				Service: tt.fields.Service,
			}
			if got := s.loadMetadata(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadMetadata() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewIpfsService(t *testing.T) {
	tests := []struct {
		name string
		want IpfsService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewIpfsService(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewIpfsService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readDirGetFileInfo(t *testing.T) {
	type args struct {
		dir string
	}
	tests := []struct {
		name    string
		args    args
		want    []os.FileInfo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readDirGetFileInfo(tt.args.dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("readDirGetFileInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readDirGetFileInfo() got = %v, want %v", got, tt.want)
			}
		})
	}
}
