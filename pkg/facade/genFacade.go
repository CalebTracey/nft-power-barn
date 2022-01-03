package facade

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"gitlab.com/CalebTracey/nft-power-barn/pkg/config"
	"image"
	"image/draw"
	"image/png"
	"sync"
	"time"

	"io/fs"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	fmCode         = int(0777)
	defaultBufSize = 4096
)

var buildDir, layersDir, rarityDelimiter, dnaDelimiter string
var editionCount, failedCount, uniqueDnaTorrence, layerConfigIdx int
var metadataSlice []Metadata
var attributesList map[int][]Attributes
var dnaList map[int]string

type GenFacade interface {
	StartCreating()
}

type GenService struct {
	Elements    Elements
	Image       *image.RGBA
	ImageLayers []image.Image
}

type Elements struct {
	layers      []LayerElement
	layerImages []ImageResponse
}

func init() {
	conf := config.InitializeConfig()
	path, err := os.Getwd()
	if err != nil {
		log.Println(err.Error())
	}
	buildDir = fmt.Sprintf("%v/build", path)
	layersDir = fmt.Sprintf("%v/layers", path)
	rand.Seed(time.Now().UnixNano())

	metadataSlice = make([]Metadata, 0)
	dnaList = make(map[int]string)
	attributesList = make(map[int][]Attributes)

	rarityDelimiter = conf.RarityDelimiter
	dnaDelimiter = conf.DnaDelimiter
	uniqueDnaTorrence = conf.UniqueDnaTorrence

	failedCount = 0
	layerConfigIdx = 0
	editionCount = 1
}

func NewService() GenService {
	return GenService{
		Elements: Elements{},
	}
}

func (s *GenService) StartCreating() {
	start := time.Now()
	log.Println("And we're off!")
	layerConfigs := config.Layers()
	allConfigs := layerConfigs.All
	abstractedIndexes := make([]int, 0)
	for i := 1; i <= allConfigs[0].EditionSize; i++ {
		abstractedIndexes = append(abstractedIndexes, i)
	}
	log.Println("Processing configs...")

	for _, conf := range allConfigs {
		setupErr := s.layersSetup(conf.LayerOrder)

		if setupErr != nil {
			log.Panicf(setupErr.Error())
			return
		}
		for editionCount <= conf.EditionSize {
			newDna := s.createDna()
			if isDnaUnique(dnaList, newDna) {
				src := image.NewRGBA(image.Rect(0, 0, 1600, 1600))
				bounds := src.Bounds()
				newImage := image.NewRGBA(bounds)
				s.Image = newImage

				results := constructLayerToDna(newDna, s.Elements.layers)
				dna, err := dnaHash(newDna)
				fatality(err)

				addAttributes(results, editionCount)
				addMetadata(dna, abstractedIndexes[0])
				var wg sync.WaitGroup
				work := make(chan []image.Image)

				wg.Add(1)
				go func() {
					s.loadImages(work, results)

				}()

				consumer := <-work
				go func() {
					defer wg.Done()
					for i := range consumer {
						s.drawLayer(consumer[i])
					}
				}()

				wg.Wait()
				go func() {
					imgErr := saveImageFile(s.Image, editionCount)
					fatality(imgErr)
					metaErr := saveMetadata(editionCount)
					fatality(metaErr)
				}()

				log.Printf("Created edition: %v, with DNA: %v", abstractedIndexes[0], dna)
				dnaList[abstractedIndexes[0]] = dna
				editionCount = editionCount + 1
				abstractedIndexes = abstractedIndexes[1:]
			} else {
				log.Println("DNA Exists!")
				failedCount++
				if failedCount >= uniqueDnaTorrence {
					log.Fatalf("You need more layers or elements to grow your edition to %v works", conf.EditionSize)
				}
			}
		}
		layerConfigIdx = layerConfigIdx + 1
	}
	writeFullMetadata()
	dur := fmt.Sprintf("%.2f", time.Since(start).Seconds())
	log.Printf("Finished in %v seconds", dur)
}

//
//func (s *GenService) DoAll(res []LayerToDnaResults, images chan image.Image) {
//	idx := 0
//
//	c := make(chan int)
//	for i := 0; i < 4; i++ {
//		go s.DoSome(i*len(s.ImageLayers)/4, (i+1)*4/4, res, c, images)
//	}
//	//for i := 0; i < 4; i++ {
//		idx += <-c
//
//	//}
//}
//
//func (s *GenService) DoSome(i, n int, res []LayerToDnaResults, c chan int, images chan image.Image) {
//	for ; i < n; i++ {
//		img, err := s.loadLayerImage(res[i])
//		if err != nil {
//			fatality(err)
//		}
//		if i < len(res) {
//			images <- img
//			<- c
//		}
//	}
//
//	c <- 1 // signal that this piece is done
//}

func (s *GenService) layersSetup(layers []config.Layer) error {
	log.Println("Setting up layers...")
	for i, layer := range layers {
		elements, err := s.getElements(fmt.Sprintf("%v/%v", layersDir, layer.Name))
		if err != nil {
			return fmt.Errorf("error fetching elements for layer: %v ||| %v", layer.Name, err)
		}
		s.Elements.layers = append(s.Elements.layers, LayerElement{
			Id:        i,
			Elements:  elements,
			Name:      layer.Name,
			Blend:     "source-over",
			Opacity:   1,
			BypassDNA: false,
		})
	}
	return nil
}

func (s *GenService) getElements(path string) (res []Element, err error) {
	var fileInfo = make(map[string]fs.FileInfo)
	items, err := ioutil.ReadDir(path)
	if err != nil {
		return res, err
	}
	err = filepath.Walk(buildDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fileInfo[path] = info
		}
		return nil
	})

	if err != nil {
		return res, err
	}
	str1, _ := regexp.Compile("(^|\\/)\\.[^\\/\\.]/g")
	str2, _ := regexp.Compile(".DS_Store")
	for i, item := range items {
		if !str1.MatchString(item.Name()) && !str2.MatchString(item.Name()) {
			res = append(res, Element{
				Id:       i,
				Name:     cleanName(item.Name()),
				FileName: item.Name(),
				Weight:   getRarityWeight(item.Name()),
				Path:     fmt.Sprintf("%v/%v", path, item.Name()),
			})
		}
	}
	return res, nil
}

func (s *GenService) loadImages(work chan []image.Image, results []LayerToDnaResults) {
	var temp []image.Image
	for _, layer := range results {
		decoded, err := s.loadLayerImage(layer)
		fatality(err)
		temp = append(temp, decoded)
	}
	work <- temp
}

func (s *GenService) loadLayerImage(layer LayerToDnaResults) (image.Image, error) {
	img, err := os.OpenFile(layer.SelectedElement.Path, os.O_RDWR|os.O_CREATE, defaultBufSize)
	defer func(img *os.File) {
		e := img.Close()
		if e != nil {
			log.Panic(e)
			return
		}
	}(img)
	type Result struct {
		Image image.Image
	}
	if err != nil {
		return nil, fmt.Errorf("failed to load file at path: %v. %v", layer.SelectedElement.Path, err)
	}
	decoded, err := png.Decode(img)
	if err != nil {
		return nil, fmt.Errorf("failed to decode file opened at path: %v. %v", layer.SelectedElement.Path, err)
	}
	return decoded, nil
}

func getRarityWeight(str string) int {
	nameWithoutExtension := str[0 : len(str)-4]
	nameSlice := strings.Split(nameWithoutExtension, rarityDelimiter)
	weight := nameSlice[len(nameSlice)-1]
	res, err := strconv.Atoi(weight)
	if err != nil && weight != "" {
		return 1
	}
	return res
}

func (s *GenService) createDna() (newDna string) {
	randNum := make([]string, 0)
	for _, layer := range s.Elements.layers {
		totalWeight := 0
		for _, e := range layer.Elements {
			totalWeight = totalWeight + e.Weight
		}
		random := rand.Intn(totalWeight)
		for i, e := range layer.Elements {
			random = random - e.Weight
			if random < 0 {
				bypassStr := ""
				if layer.BypassDNA {
					bypassStr = "?bypassDNA=true"
				}
				data := fmt.Sprintf("%v:%v%v", layer.Elements[i].Id, layer.Elements[i].FileName, bypassStr)
				randNum = append(randNum, data)
				break
			}
		}
	}
	return strings.Join(randNum, dnaDelimiter)
}

func isDnaUnique(list map[int]string, newDna string) bool {
	for _, dna := range list {
		if strings.EqualFold(dna, newDna) {
			return false
		}
	}
	return true
}

func constructLayerToDna(dna string, layers []LayerElement) (res []LayerToDnaResults) {
	for i, layer := range layers {
		elementId, err := cleanDna(&strings.Split(dna, dnaDelimiter)[i])

		for _, e := range layer.Elements {
			if err != nil {
				log.Printf("dna error: %v", dna)
				log.Fatalln(err)
			}
			if e.Id == elementId {
				res = append(res, LayerToDnaResults{
					Name:            layer.Name,
					Blend:           layer.Blend,
					Opacity:         layer.Opacity,
					SelectedElement: e,
				})
			}
		}
	}
	return res
}

func cleanName(name string) string {
	clean := name[0 : len(name)-4]
	split := strings.Split(clean, rarityDelimiter)
	return split[0]
}

func cleanDna(dna *string) (int, error) {
	clean := removeQueryStrings(*dna)
	split := strings.Split(clean, ":")
	if split[0] == "" {
		return -1, fmt.Errorf("bad element at index 0 of dna string: %v", *dna)
	}
	x := split[0]
	split = split[1:]

	return strconv.Atoi(x)
}

func removeQueryStrings(dna string) string {
	query := regexp.MustCompile("\\?.*$")
	return query.ReplaceAllString(dna, "")
}

func addAttributes(results []LayerToDnaResults, idx int) {
	ig := config.IgnoredData()
	var abs []Attributes
	for _, layer := range results {
		var igTrait bool
		name := strings.TrimSpace(layer.Name)
		trait := strings.TrimSpace(layer.SelectedElement.Name)
		for _, t := range ig.Traits {
			if strings.EqualFold(t, name) {
				igTrait = true
				break
			}
		}
		if strings.EqualFold(trait, "blank") {
			trait = "None"
		}
		if !igTrait {
			abs = append(abs, Attributes{
				TraitType: name,
				Value:     trait,
			})
		}
	}
	attributesList[idx] = abs
}

func findAttributes(idx int) ([]Attributes, error) {
	abs := attributesList[idx]
	if len(abs) <= 0 {
		return []Attributes{}, fmt.Errorf("error finding attributes for index %v, of attributes list", idx)
	}
	return abs, nil
}

func (s *GenService) drawLayer(img image.Image) {
	pointVal := image.Point{X: 1600, Y: 1600}
	mask := &image.Rectangle{
		Min: image.Point{},
		Max: pointVal,
	}
	draw.DrawMask(s.Image, image.Rectangle{
		Min: image.Point{},
		Max: pointVal,
	},
		img, image.Point{},
		mask, image.Point{},
		draw.Over)
}

func saveImageFile(newImg *image.RGBA, edition int) error {
	path := fmt.Sprintf("%v/images", buildDir)
	fileName := filepath.Join(path, filepath.Base(fmt.Sprintf("/%v.png", strconv.Itoa(edition))))

	outFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer func(outFile *os.File) {
		err = outFile.Close()
		if err != nil {
			return
		}
	}(outFile)

	buff := new(bytes.Buffer)
	err = png.Encode(buff, newImg)
	if err != nil {
		return err
	}

	_, err = outFile.Write(buff.Bytes())
	if err != nil {
		return err
	}
	err = outFile.Sync()
	if err != nil {
		return err
	}

	err = outFile.Close()
	if err != nil {
		return err
	}
	return nil
}

func addMetadata(dna string, edition int) {
	conf := config.Metadata()
	attributes, err := findAttributes(edition)
	if err != nil {
		log.Println(err)
		attributes = make([]Attributes, 0)
	}
	metaData := Metadata{
		Name:        conf.Name,
		Description: conf.Description,
		FileUrl:     conf.FileUrl,
		Creator:     conf.Creator,
		CustomFields: CustomFields{
			DNA:      dna,
			Edition:  edition,
			Date:     fmt.Sprint(time.Now().Date()),
			Compiler: "Caleb's NFT Power Barn",
		},
		Attributes: attributes,
	}

	metadataSlice = append(metadataSlice, metaData)
}

func saveMetadata(edition int) error {
	for _, m := range metadataSlice {
		if m.CustomFields.Edition == edition {
			err := saveMetadataFile(m, edition)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func saveMetadataFile(m Metadata, edition int) error {
	path := fmt.Sprintf("%v/json", buildDir)
	outFile, err := os.Create(filepath.Join(path, filepath.Base(fmt.Sprintf("/%v.json", strconv.Itoa(edition)))))
	if err != nil {
		return err
	}
	defer func(outFile *os.File) {
		err = outFile.Close()
		if err != nil {
			return
		}
	}(outFile)

	buff := new(bytes.Buffer)
	enc := json.NewEncoder(buff)
	enc.SetIndent("", "    ")
	if err = enc.Encode(m); err != nil {
		panic(err)
	}
	_, err = outFile.Write(buff.Bytes())
	if err != nil {
		return err
	}
	err = outFile.Sync()
	if err != nil {
		return err
	}
	err = outFile.Close()
	if err != nil {
		return err
	}
	return nil
}

func createMetadataFile(path string, edition int) (*os.File, error) {
	outFile, err := os.Create(filepath.Join(path, filepath.Base(fmt.Sprintf("%v.json", strconv.Itoa(edition)))))
	defer func(outFile *os.File) {
		err = outFile.Close()
		if err != nil {
			return
		}
	}(outFile)
	if err != nil {
		return nil, err
	}
	return outFile, nil
}

func writeFullMetadata() {
	path := fmt.Sprintf("%v/json", buildDir)
	outFile, err := os.Create(filepath.Join(path, filepath.Base("_metadata.json")))
	fatality(err)
	file, err := json.MarshalIndent(metadataSlice, "", " ")
	fatality(err)
	err = ioutil.WriteFile(outFile.Name(), file, os.FileMode(fmCode))
	fatality(err)
}

func dnaHash(dna string) (string, error) {
	if dna == "" {
		return "", fmt.Errorf("missing DNA")
	}
	h := sha1.New()
	h.Write([]byte(dna))
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func fatality(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
