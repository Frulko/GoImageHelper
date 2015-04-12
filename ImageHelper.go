package main

import (
    "fmt"
    "image"
    "os"
	"image/draw"
    _ "image/jpeg"
	"image/png"
    "encoding/json"
    "path/filepath"
	"io/ioutil"
)

type Output struct {
	Result bool `json:"result"`
	Message string `json:"message"`
}

type Assets struct {
    Assets []string `json:"assets"`
}

type test struct {
	Canvas Canvas `json:"canvas"`
	Assets []ImageAsset `json:"assets"`
}

type Canvas struct {
	Name   string `json:"name"`
	Width   int `json:"width"`
	Height   int `json:"height"`
}

type ImageAsset struct {
    Image   string `json:"filename"`
    Width   int `json:"width"`
    Height  int `json:"height"`
    Valid   bool `json:"validate"`
	Top     int  `json:"top"`
	Left     int  `json:"left"`
}

type ImageAssets []* ImageAsset

func main() {
    //img_path := "rainy.png";
	//



    if len(os.Args) <= 1  {
		log(false, "Specify a command")
        os.Exit(0);
    }

    cmd := os.Args[1]
	if cmd == "infos" {
		if len(os.Args) == 2  {
			log(false, "Specify a JSON file")
			os.Exit(0);
		}

		exchange_filename := os.Args[2]
		res := loadExchangeFileJsonData(exchange_filename)
		glob := getImagesInformations(res.Assets)
		writeOutImagesInformations(glob, exchange_filename)


	} else if cmd == "generate" {


		if len(os.Args) == 2  {
			log(false, "Specify a JSON file")
			os.Exit(0);
		}

		exchange_filename := os.Args[2]
		//res := loadExchangeFileJsonData(exchange_filename)
		file, e := ioutil.ReadFile(exchange_filename)
		if e != nil {
			log(false, "error loading xchange file")
			os.Exit(0)
		}
		str := string(file)

		res := &test{}
		json.Unmarshal([]byte(str), &res)


		m := image.NewRGBA(image.Rect(0, 0, res.Canvas.Width, res.Canvas.Height))

		for _,el := range res.Assets {
			fImg2, _ := os.Open(el.Image)
			defer fImg2.Close()
			img2, _, _ := image.Decode(fImg2)

			posX := el.Top * -1
			posY := el.Left * -1

			draw.Draw(m, m.Bounds(), img2, image.Point{posY,posX}, draw.Src)
		}



		toimg, _ := os.Create(res.Canvas.Name)
		defer toimg.Close()
		png.Encode(toimg, m)
		log(true, "Generate Canvas")

	}





}

func getImageDimension(imagePath string) (int, int, bool) {
	validate := true
    file, err := os.Open(imagePath)
    if err != nil {
        validate = false
    }
    image, _, err := image.DecodeConfig(file)
    if err != nil {
		validate = false
    }
    return image.Width, image.Height, validate
}


func getImagesInformations(images []string) (ImageAssets) {

	t_images := ImageAssets{}
	for _,image_path := range images {
		img_path, _ := filepath.Abs(image_path)
		width, height, validate := getImageDimension(img_path)
		asset := &ImageAsset{ Image: img_path, Width: width, Height: height, Valid: validate, Top: 0, Left: 0}
		t_images = append(t_images , asset)
	}

	return t_images
}

func loadExchangeFileJsonData(filename string) (*Assets){
	if filepath.Ext(filename) != ".json"{
		log(false, "Specify a JSON file")
		os.Exit(0)
	}

	file, e := ioutil.ReadFile(filename)
	if e != nil {
		log(false, "error loading xchange file")
		os.Exit(0)
	}
	str := string(file)
	res := &Assets{}
	json.Unmarshal([]byte(str), res)

	return res
}

func log(state bool, message string){
	res1B, e := json.Marshal(Output{state, message})
	if e != nil {
		log(false, "Error while encode JSON")
		os.Exit(0)
	}
	fmt.Println(string(res1B))
}

func writeOutImagesInformations(assets ImageAssets, filename string){

	json_encoded, e := json.Marshal(assets)
	if e != nil {
		log(false, "Error while encode JSON")
		os.Exit(0)
	}

	converted_json := []byte(json_encoded);
	err := ioutil.WriteFile(filename, converted_json, 0644)
	if err != nil {
		log(false, "Error on write JSON file")
		os.Exit(0)
	}

	log(true, "task complete")
}
