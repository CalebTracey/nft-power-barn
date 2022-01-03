package facade

import (
	"fmt"
	"gitlab.com/CalebTracey/nft-power-barn/pkg/config"
	"image"
	"image/png"
	"os"
	"reflect"
	"testing"
)

func Test_getElements(t *testing.T) {
	baseDir := ""
	tests := []struct {
		name     string
		path     string
		elements *Elements
		wantRes  []Element
		wantErr  bool
	}{
		{
			name:     "Happy Path",
			path:     baseDir,
			elements: &Elements{},
			wantErr:  false,
			wantRes: []Element{
				{
					Id:       0,
					Name:     "Blue Green ",
					FileName: "Blue Green #50.png",
					Weight:   50,
					Path:     fmt.Sprintf("%v/Blue Green #50.png", baseDir),
				}, {
					Id:       1,
					Name:     "Cadet Blue ",
					FileName: "Cadet Blue #50.png",
					Weight:   50,
					Path:     fmt.Sprintf("%v/Cadet Blue #50.png", baseDir),
				}, {
					Id:       2,
					Name:     "Cadet Blue ",
					FileName: "Cadet Blue #60.png",
					Weight:   60,
					Path:     fmt.Sprintf("%v/Cadet Blue #60.png", baseDir),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &GenService{
				Elements: *tt.elements,
			}
			gotRes, err := s.getElements(tt.path)

			if (err != nil) != tt.wantErr {
				t.Errorf("getElements() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("getElements() gotRes = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_getRarityWeight(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want int
	}{
		{
			name: "Happy Path",
			str:  "Blue Green #50.png",
			want: 50,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getRarityWeight(tt.str); got != tt.want {
				t.Errorf("getRarityWeight() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenService_layersSetup(t *testing.T) {
	baseDir := "/Users/calebtracey/GolandProjects/generatecollection/pkg/layers/_Test_"
	type fields struct {
		Elements *Elements
	}
	testLayers := []config.Layer{
		{
			"Test 1",
		},
		{
			"Test 2",
		},
	}

	tests := []struct {
		name    string
		fields  fields
		layers  []config.Layer
		wantErr bool
	}{
		{
			name:    "Happy Path",
			wantErr: false,
			fields: fields{
				Elements: &Elements{
					layers: []LayerElement{
						{
							10,
							[]Element{
								{
									Id:       0,
									Name:     "Blue Green ",
									FileName: "Blue Green #50.png",
									Weight:   50,
									Path:     fmt.Sprintf("%v/Blue Green #50.png", baseDir),
								}, {
									Id:       1,
									Name:     "Cadet Blue ",
									FileName: "Cadet Blue #50.png",
									Weight:   50,
									Path:     fmt.Sprintf("%v/Cadet Blue #50.png", baseDir),
								}, {
									Id:       2,
									Name:     "Cadet Blue ",
									FileName: "Cadet Blue #60.png",
									Weight:   60,
									Path:     fmt.Sprintf("%v/Cadet Blue #60.png", baseDir),
								},
							},
							"TEST",
							"TEST",
							10,
							false,
						},
					},
				},
			},
			layers: testLayers,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &GenService{
				Elements: *tt.fields.Elements,
			}

			if err := s.layersSetup(tt.layers); (err != nil) != tt.wantErr {
				t.Errorf("layersSetup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenService_createDna(t *testing.T) {
	baseDir := "/Users/calebtracey/GolandProjects/generatecollection/pkg/layers/_Test_"
	type fields struct {
		Elements *Elements
		Image    *image.RGBA
	}
	tests := []struct {
		name       string
		fields     fields
		wantNewDna string
	}{
		{
			name: "Happy Path",
			fields: fields{
				Elements: &Elements{
					layers: []LayerElement{
						{
							10,
							[]Element{
								{
									Id:       0,
									Name:     "Blue Green ",
									FileName: "Blue Green #50.png",
									Weight:   50,
									Path:     fmt.Sprintf("%v/Blue Green #50.png", baseDir),
								}, {
									Id:       1,
									Name:     "Cadet Blue ",
									FileName: "Cadet Blue #50.png",
									Weight:   50,
									Path:     fmt.Sprintf("%v/Cadet Blue #50.png", baseDir),
								}, {
									Id:       2,
									Name:     "Cadet Blue ",
									FileName: "Cadet Blue #60.png",
									Weight:   60,
									Path:     fmt.Sprintf("%v/Cadet Blue #60.png", baseDir),
								},
							},
							"TEST",
							"TEST",
							10,
							false,
						},
					},
				},
				Image: &image.RGBA{},
			},
			wantNewDna: "1:Cadet Blue #50.png",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &GenService{
				Elements: *tt.fields.Elements,
				Image:    tt.fields.Image,
			}
			if gotNewDna := s.createDna(); gotNewDna != tt.wantNewDna {
				t.Errorf("createDna() = %v, want %v", gotNewDna, tt.wantNewDna)
			}
		})
	}
}

func Test_saveImageFile(t *testing.T) {
	tests := []struct {
		name    string
		newImg  *image.RGBA
		edition int
	}{
		{
			name: "Happy Path",
			newImg: image.NewRGBA(image.Rectangle{
				Min: image.Point{},
				Max: image.Point{},
			}),
			edition: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

//func BenchmarkTest_saveImageFile(b *testing.B) {
//	src := image.NewRGBA(image.Rect(0, 0, 1600, 1600))
//	bounds := src.Bounds()
//	testImg := image.NewRGBA(bounds)
//
//	tests := []struct{
//		name string
//		newImg *image.RGBA
//		edition int
//	}{
//		{
//			name: "Single File Benchmark",
//			newImg: testImg,
//			edition: 1,
//		},
//	}
//	b.Run("Single File Benchmark", func(b *testing.B) {
//
//	})
//}

func readTestFile() (*image.Image, error) {
	testPath := "/Users/calebtracey/GolandProjects/generatecollection/layers/_Test_"
	img, err := os.Open(fmt.Sprintf("%v/Blue Green #50", testPath))
	defer func(img *os.File) {
		e := img.Close()
		if e != nil {
			return
		}
	}(img)
	if err != nil {
		return nil, err
	}
	decoded, err := png.Decode(img)
	if err != nil {
		return nil, err
	}
	return &decoded, nil
}
