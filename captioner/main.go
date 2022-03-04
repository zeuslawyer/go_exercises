package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

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
	FILE_NAME := getFilePathFromUser()
	fileData, err := ioutil.ReadFile(FILE_NAME)
	if err != nil {
		log.Fatalln("error loading file to be transcribed: ", err)
	}

	// Setup HTTP client and set header
	client := &http.Client{}
	var responses map[string]interface{}

	// Make HTTP request to upload file.
	// res := requestEndpoint(client, "POST", UPLOAD_ENDPOINT, false, bytes.NewBuffer(fileData))
	// updateResponses(*res, &responses)

	// fileUrl := responses["upload_url"]
	// fileEndpoint, ok := fileUrl.(string)
	// if !ok {
	// 	log.Fatalln("file endpoint is not of type string")
	// }
	// fmt.Println("File is uploaded to ", fileEndpoint)

	// // Construct JSON object to send to transcription API.
	// sendData := map[string]string{"audio_url": fileEndpoint}
	// sendDataJson, err := json.Marshal(sendData)
	// if err != nil {
	// 	log.Fatalln("error marshalling data to be sent into json: ", err)
	// }

	// Make HTTP request to initiate transcription and get id.
	// res = requestEndpoint(client, "POST", TRANSCRIPT_ENDPOINT, true, bytes.NewBuffer(sendDataJson))
	// updateResponses(*res, &responses)

	// id := responses["id"]
	// idString, ok := id.(string)
	// if !ok {
	// 	log.Fatalln("transcription id is not a string")
	// }
	// fmt.Println("Transcription succeeded - id is: ", id)

	// Make HTTP request to get transcription text.
	STATUS_POLLING_URL := TRANSCRIPT_ENDPOINT + "/" + "ogfuflk54l-ec9f-48b6-b9c6-b8278e23eeac" // idString  TODO reinstate
	_ = fileData                                                                               // TODO remove

	done := false
	for !done {
		res := requestEndpoint(client, "GET", STATUS_POLLING_URL, false, nil)

		status, ok := responses["status"].(string)
		if !ok {
			fmt.Printf("Status field value is: %v\n", responses["status"])
		}

		fmt.Printf("ZUBIN the res status is : %v\n", status)

		status = "completed" // TODO reemove

		if status == "completed" {
			done = true
			parseJsonResp(*res, &responses)
			break
		}
		// Processing not complete.
		fmt.Println("Transcription not complete. Trying again in 1 minute...")
		time.Sleep(1 * time.Second)
		// time.Sleep(1 * time.Minute) // TODO reinstate
	}

	// Get the SRT text.
	SRT_ENDPOINT := TRANSCRIPT_ENDPOINT + "/" + "ogfuflk54l-ec9f-48b6-b9c6-b8278e23eeac" + "/" + "srt" //   TODO replace with idString
	res := requestEndpoint(client, "GET", SRT_ENDPOINT, false, nil)
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("error reading SRT text from response: ", err)
	}

	fmt.Println("Here is the transcribed text: ", string(body))

	writeSrtFile(FILE_NAME, string(body))

}

func getApiKey(key string) string {
	env, ok := viper.Get(key).(string)
	if !ok {
		log.Fatal("Env Variable not of type string")
	}
	return env
}

func prettyPrintMap(m map[string]interface{}) string {
	s, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		fmt.Println("error marshalling map to indented json:", err)
	}
	return string(s)
}

func getFilePathFromUser() string {
	return "/Users/zubinpratap/Downloads/clip.mp4"
}

func requestEndpoint(client *http.Client, method, url string, contentTypeJson bool, body io.Reader) *http.Response {
	fmt.Printf("making %s request to %s\n", method, url)
	API_KEY := getApiKey("API_KEY")

	// Construct request to the endpoint.
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Fatalln("error constructing new request to  transcription endpoint: ", err)
	}

	// Set headers.
	req.Header.Set("authorization", API_KEY)
	if contentTypeJson {
		req.Header.Set("content-type", "application/json")
	}

	// Make request.
	res, err := client.Do(req)
	if err != nil {
		log.Fatalln("error making http request to the transcription endpoint: ", err)
	}

	if res.StatusCode != http.StatusOK {
		fmt.Printf("Request status code not OK - %v", res.StatusCode)
	}

	return res
}

// Takes the given http responses and decodes it's JSON into the the responses map pointer.
func parseJsonResp(res http.Response, responses *map[string]interface{}) {
	if err := json.NewDecoder(res.Body).Decode(responses); err != nil {
		log.Fatalln("error decoding result to JSON: ", err)
	}

	resps := *responses
	if resps["error"] != nil {
		log.Fatalln("error string found in responses: ", resps["error"])
	}

	fmt.Printf("Just updated responses with %v", prettyPrintMap(*responses))
}

func writeSrtFile(transcribedFilePath, srtData string) {
	dir, sourceFile := filepath.Split(transcribedFilePath)
	fmt.Printf("REMOVE ME:  %s", dir+sourceFile+".srt") // TODO remove

	f, err := os.Create(dir + sourceFile + ".srt")
	if err != nil {
		log.Fatalln("error creating srt file: ", err)
	}
	defer f.Close()

	_, err2 := f.WriteString(srtData)

	if err2 != nil {
		log.Fatalln("error writing srt file: ", err2)
	}
}
