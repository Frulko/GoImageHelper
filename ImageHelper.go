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
	"strings"
)

const version string = "0.1.0"

type Output struct {
	Version string `json:"version"`
	Result bool `json:"result"`
	Message string `json:"message"`
	Code string `json:"code"`
}

type Assets struct {
    Assets []string `json:"assets"`
}

type ExchangeGenerate struct {
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
		log(false, "Specify a command", "50")
        os.Exit(0);
    }



    cmd := os.Args[1]
	available_cmd := []string{"infos","generate", "version"}
	if(!stringInSlice(cmd, available_cmd)) {
		var commands string = "";
		for index,el := range available_cmd {
			sep := ","
			if(len(available_cmd)-1 <= index){
				sep = ""
			}
			commands += el+sep+" "
		}
		log(false, "Specify a valid command (valid command : "+commands+")", "51")
		os.Exit(0);
	}else{


		switch cmd {
		case "infos":
			getInformations()
		case "generate":
			generateCanvas()
		case "version":
			getVersion()
		default:
			log(true, "Do nothing", "20")
		}


	}



}

func getImageDimension(imagePath string) (int, int, bool) {
	validate := true
    file, err := os.Open(imagePath)
    if(err != nil) {validate = false}

    image, _, err := image.DecodeConfig(file)
    if(err != nil) {validate = false}

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
		log(false, "Specify a JSON file", "52")
		os.Exit(0)
	}

	file, err := ioutil.ReadFile(filename)
	crash(err, "error loading xchange file", "55")
	str := string(file)
	res := &Assets{}
	json.Unmarshal([]byte(str), res)

	return res
}

func log(state bool, message string, optional ...string){


	code := "00"

	if(len(optional) > 1){
		log(false, "Log too much args", "50")
		os.Exit(0)
	}else if len(optional) == 1 {
		code = optional[0]
	}

	json_result, err := json.Marshal(Output{version, state, message, code})
	crash(err, "Error while encode JSON", "54")
	fmt.Println(string(json_result))
}

func getVersion(){
	log(true, "GoImageHelper", "20")
}

func writeOutImagesInformations(assets ImageAssets, filename string){

	json_encoded, err := json.Marshal(assets)
	crash(err, "Error while encode JSON for ImageInformations", "54")

	converted_json := []byte(json_encoded);
	var file_base string = filepath.Base(filename)
	crash(ioutil.WriteFile(filename, converted_json, 0644), "Error on write JSON file for : "+file_base, "53")

	log(true, "task complete", "21")
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func getInformations(){
	if(hasAJSONFile()){
		exchange_filename := os.Args[2]
		res := loadExchangeFileJsonData(exchange_filename)
		glob := getImagesInformations(res.Assets)
		writeOutImagesInformations(glob, exchange_filename)
	}
}

func generateCanvas() {
	if(hasAJSONFile()){
		exchange_filename := os.Args[2]
		file, err := ioutil.ReadFile(exchange_filename)
		crash(err, "error loading xchange file", "55")

		str := string(file)

		res := &ExchangeGenerate{}

		crash(json.Unmarshal([]byte(str), &res), "error JSON Parse", "40")



		if(len(res.Canvas.Name) < 1 && len(res.Assets) < 1){
			log(false, "Error in JSON structure", "47")
			os.Exit(0)
		}

		if(len(res.Canvas.Name) < 1){
			log(false, "Specify a valid output name", "41")
			os.Exit(0)
		}

		if(res.Canvas.Width < 1 || res.Canvas.Height < 1){
			log(false, "Specify a valid size of output file", "42")
			os.Exit(0)
		}

		if(len(res.Assets) < 1){
			log(false, "Error no assets was specified", "46")
			os.Exit(0)
		}



		canvas := image.NewRGBA(image.Rect(0, 0, res.Canvas.Width, res.Canvas.Height))

		ext := []string{".png",".jpg",".jpeg"}

		for _,el := range res.Assets {

			var file_extension string = strings.ToLower(filepath.Ext(el.Image))
			var file_base string = filepath.Base(el.Image)

			if(len(el.Image) < 1){
				log(false, "There is an error in your JSON 'assets' structure", "47")
				os.Exit(0)
			}

			if(!stringInSlice(file_extension, ext)) {
				log(false, "error extension '"+ file_extension +"' of '" + file_base + "' is not a valid image", "43")
				os.Exit(0)
			}

			file, err := os.Open(el.Image)
			crash(err, "error when opening file " + file_base, "44")

			defer file.Close()

			loaded_image, _, err := image.Decode(file)
			crash(err, "error when decoding image " +file_base, "45")

			posX := el.Top * -1
			posY := el.Left * -1

			draw.Draw(canvas, canvas.Bounds(), loaded_image, image.Point{posY,posX}, draw.Src)
		}



		output_canvas, err := os.Create(res.Canvas.Name)
		crash(err, "error save PNG Canvas", "56")

		defer output_canvas.Close()
		err = png.Encode(output_canvas, canvas)
		crash(err, "error write image to PNG format", "57")

		log(true, "Generate Canvas", "22")

	}
}

func hasAJSONFile() bool{
	if len(os.Args) == 2  {
		log(false, "Specify a JSON file")
		os.Exit(0);
	}

	return true
}

func crash(err error, msg string, code string){
	if err != nil {
		log(false, msg, code)
		os.Exit(0)
	}
}
