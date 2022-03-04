package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

// We use viper to read env vars.
func init() {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
}

const (
	UPLOAD_ENDPOINT     = "https://api.assemblyai.com/v2/upload"
	TRANSCRIPT_ENDPOINT = "https://api.assemblyai.com/v2/transcript"
)

func main() {

	API_KEY := getApiKey("API_KEY")
	FILE_PATH := "/Users/zubinpratap/Downloads/clip.mp4"

	fileData, err := ioutil.ReadFile(FILE_PATH)
	if err != nil {
		log.Fatalln("error loading file to be transcribed: ", err)
	}

	// Setup HTTP client and set header
	client := &http.Client{}
	req, err := http.NewRequest("POST", UPLOAD_ENDPOINT, bytes.NewBuffer(fileData))
	if err != nil {
		log.Fatalln("error constructing new request to upload endpoint: ", err)
	}
	req.Header.Set("authorization", API_KEY)

	// Make the HTTP request.
	res, err := client.Do(req)
	if err != nil {
		log.Fatalln("error making http request to the upload endpoint: ", err)
	}
	defer res.Body.Close()

	// Decode JSON response into a map.
	var result map[string]interface{} // result is expected to be a map.
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		log.Fatalln("error decoding result to JSON: ", err)
	}

	fmt.Println(" The file is uploaded and available at: ", result)
	fileUrl := result["upload_url"]
	fileEndpoint, ok := fileUrl.(string)
	if !ok {
		log.Fatalln("file endpoint is not of type string")
	}

	// Construct JSON object to send to transcription API.
	sendData := map[string]string{"audio_url": fileEndpoint}
	sendDataJson, err := json.Marshal(sendData)
	if err != nil {
		log.Fatalln("error marshalling data to be sent into json: ", err)
	}

	// Make HTTP request.
	req, err = http.NewRequest("POST", TRANSCRIPT_ENDPOINT, bytes.NewBuffer(sendDataJson))
	if err != nil {
		log.Fatalln("error constructing new request to  transcription endpoint: ", err)
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("authorization", API_KEY)

	res, err = client.Do(req)
	if err != nil {
		log.Fatalln("error making http request to the transcription endpoint: ", err)
	}
	defer res.Body.Close()

	var transcribeResult map[string]interface{}
	json.NewDecoder(res.Body).Decode(&transcribeResult)

	id := transcribeResult["id"]
	fmt.Println("Id is ", id)

}

func getApiKey(key string) string {
	env, ok := viper.Get(key).(string)
	if !ok {
		log.Fatal("Env Variable not of type string")
	}
	return env
}
