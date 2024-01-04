package crowdin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type postOptions struct {
	urlStr   string
	body     interface{}
	fileName string
}

type delOptions struct {
	urlStr string
	body   interface{}
}

type patchOptions struct {
	urlStr string
	body   interface{}
}

type putOptions struct {
	urlStr string
	body   interface{}
}

type getOptions struct {
	urlStr string
	params map[string]string
	//	body   interface{}
}

// POST request
func (crowdin *Crowdin) post(options *postOptions) ([]byte, error) {

	crowdin.log(fmt.Sprintf("Create POST http request\nBody: %s", options.body))

	var req *http.Request
	var err error

	buf := new(bytes.Buffer)

	if options.fileName == "" { // Doesn't include a file to upload
		json.NewEncoder(buf).Encode(options.body)
		req, err = http.NewRequest("POST", options.urlStr, buf)
		if err != nil {
			crowdin.log(fmt.Sprintf("Post() - can't create a http request %s", req))
			return nil, err
		}

		// Set headers
		req.Header.Set("Authorization", "Bearer "+crowdin.config.token)
		req.Header.Set("Content-Type", "application/json")

		crowdin.log(fmt.Sprintf("Headers: %s", req.Header))
		// DEBUG
		// dump, err := httputil.DumpRequestOut(req, true)
		// crowdin.log(dump)

	} else { // There is a file to upload
		openfile, err := os.Open(options.fileName)
		defer openfile.Close()
		if err != nil {
			crowdin.log(fmt.Sprintf("Post() - can't open %s", options.fileName))
			return nil, err
		}
		fileStat, _ := openfile.Stat() //Get info from file
		// fileSize := strconv.FormatInt(fileStat.Size(), 10) //Get file size as a string
		crowdin.log(fmt.Sprintf("post() - %s size of the file to upload: %d", options.fileName, fileStat.Size()))

		req, err = http.NewRequest("POST", options.urlStr, ioutil.NopCloser(openfile))
		if err != nil {
			crowdin.log(fmt.Sprintf("post() - can't create a http request %s", req))
			return nil, err
		}
		req.ContentLength = fileStat.Size()

		// Set headers
		req.Header.Set("Authorization", "Bearer "+crowdin.config.token)
		req.Header.Set("Content-Type", "application/octet-stream")
		req.Header.Set("Crowdin-API-FileName", filepath.Base(options.fileName)) // Extract file name - at this point there has to be one...
		// req.Header.Set("Content-Length", FileSize)
	}

	// Run the  request
	response, err := crowdin.config.client.Do(req)
	if err != nil {
		crowdin.log(fmt.Sprintf("Post() - Do() returned an error %s", response))
		return nil, err
	}
	defer response.Body.Close()

	bodyResponse, err := ioutil.ReadAll(response.Body)
	if err != nil {
		crowdin.log(fmt.Sprintf("Post() - Error while reading request response %s", err))
		return nil, err
	}

	if response.StatusCode < http.StatusOK || response.StatusCode > http.StatusIMUsed {
		return bodyResponse, APIError{What: fmt.Sprintf("Status code: %v", response.StatusCode)}
	}

	return bodyResponse, nil
}

// PUT request
func (crowdin *Crowdin) put(options *putOptions) ([]byte, error) {

	crowdin.log(fmt.Sprintf("Create PUT http request\nBody: %s", options.body))

	var req *http.Request
	var err error

	buf := new(bytes.Buffer)

	json.NewEncoder(buf).Encode(options.body)
	req, err = http.NewRequest("PUT", options.urlStr, buf)
	if err != nil {
		crowdin.log(fmt.Sprintf("Put() - can't create a http request %s", req))
		return nil, err
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+crowdin.config.token)
	req.Header.Set("Content-Type", "application/json")
	crowdin.log(fmt.Sprintf("Headers: %s", req.Header))

	// DEBUG
	// dump, err := httputil.DumpRequestOut(req, true)
	// crowdin.log(dump)

	// Run the  request
	response, err := crowdin.config.client.Do(req)
	if err != nil {
		crowdin.log(fmt.Sprintf("Put() - Do() returned an error %s", response))
		return nil, err
	}
	defer response.Body.Close()

	bodyResponse, err := ioutil.ReadAll(response.Body)
	if err != nil {
		crowdin.log(fmt.Sprintf("Put() - Error while reading request response %s", err))
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return bodyResponse, APIError{What: fmt.Sprintf("Status code: %v", response.StatusCode)}
	}

	return bodyResponse, nil
}

// PATCH request
func (crowdin *Crowdin) patch(options *patchOptions) ([]byte, error) {

	crowdin.log(fmt.Sprintf("Create PATCH http request\nBody: %s", options.body))

	var req *http.Request
	var err error

	buf := new(bytes.Buffer)

	json.NewEncoder(buf).Encode(options.body)
	req, err = http.NewRequest("PATCH", options.urlStr, buf)
	if err != nil {
		crowdin.log(fmt.Sprintf("patch() - can't create a http request %s", req))
		return nil, err
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+crowdin.config.token)
	req.Header.Set("Content-Type", "application/json")
	crowdin.log(fmt.Sprintf("Headers: %s", req.Header))

	// DEBUG
	// dump, err := httputil.DumpRequestOut(req, true)
	// crowdin.log(dump)

	// Run the  request
	response, err := crowdin.config.client.Do(req)
	if err != nil {
		crowdin.log(fmt.Sprintf("patch() - Do() returned an error %s", response))
		return nil, err
	}
	defer response.Body.Close()

	bodyResponse, err := ioutil.ReadAll(response.Body)
	if err != nil {
		crowdin.log(fmt.Sprintf("patch() - Error while reading request response %s", err))
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return bodyResponse, APIError{What: fmt.Sprintf("Status code: %v", response.StatusCode)}
	}

	return bodyResponse, nil
}

// DEl request
func (crowdin *Crowdin) del(options *delOptions) ([]byte, error) {

	crowdin.log(fmt.Sprintf("Create DEL http request\nBody: %s", options.body))

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(options.body)
	req, err := http.NewRequest("DELETE", options.urlStr, buf)
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+crowdin.config.token)
	req.Header.Set("Content-Type", "application/json")
	crowdin.log(fmt.Sprintf("Headers: %s", req.Header))

	// Run the  request
	response, err := crowdin.config.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	bodyResponse, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusNoContent {
		return bodyResponse, APIError{What: fmt.Sprintf("Status code: %v", response.StatusCode)}
	}

	return bodyResponse, nil
}

// GET request
func (crowdin *Crowdin) get(options *getOptions) ([]byte, error) {

	crowdin.log(fmt.Sprintf("Create GET http request Params: %v", options.params))

	// Get request with authorization
	response, err := crowdin.getResponse(options, true)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	bodyResponse, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return bodyResponse, APIError{What: fmt.Sprintf("Status code: %v", response.StatusCode)}
	}

	return bodyResponse, nil
}

// Get request with or without authorization token depending on flag
func (crowdin *Crowdin) getResponse(options *getOptions, authorization bool) (*http.Response, error) {

	crowdin.log(fmt.Sprintf("getResponse()"))

	if options != nil && options.params != nil {
		addParam := "?"
		for k, v := range options.params {
			if v != "" {
				options.urlStr += addParam + k + "=" + v
				addParam = "&"
			}
		}
	}

	// buf := new(bytes.Buffer)
	// json.NewEncoder(buf).Encode(options.body)
	crowdin.log(fmt.Sprintf("url=%s", options.urlStr))

	req, err := http.NewRequest("GET", options.urlStr, nil)
	if err != nil {
		return nil, err
	}

	if authorization {
		req.Header.Set("Authorization", "Bearer "+crowdin.config.token)
	}

	response, err := crowdin.config.client.Do(req)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// DownloadFile will download a url and store it in local filepath.
// No autorization token required here for this operation.
// Writes to the destination file as it downloads it, without
// loading the entire file into memory.
func (crowdin *Crowdin) DownloadFile(url string, filepath string) error {

	crowdin.log(fmt.Sprintf("DownloadFile() %s", filepath))

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		crowdin.log(fmt.Sprintf("	Download error - open file error:\n %s\n"))
		return err
	}
	defer out.Close()

	// Get request - no authorization
	resp, err := crowdin.getResponse(&getOptions{urlStr: url}, false)
	// resp, err := http.Get(url)
	if err != nil {
		//fmt.Printf("\nDownload error:%s\n",resp)
		crowdin.log(fmt.Sprintf("	Download error - Get retunrs:\n %s \n", err.Error()))
		return err
	}
	defer resp.Body.Close()

	// Write body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		// log.Println("	Download error\n", resp)
		crowdin.log(fmt.Sprintf("	Download error - write to file error:\n %s \n", err.Error()))
		return err
	}
	return nil
}

// Log writer
// Hide keys by overwriting with XXX
func (crowdin *Crowdin) log(a interface{}) {
	if crowdin.debug {
		if crowdin.logWriter != nil {
			timestamp := time.Now().Format(time.RFC3339)
			msg := fmt.Sprintf("%v: %v", timestamp, a)
			token := "Authorization:[Bearer " // prefix for key
			var purged string                 // Build the purged string in here
			list1 := strings.Split(msg, token)
			if len(list1) > 1 {
				for k1, v1 := range list1 {
					if k1 > 0 { // The 1st v1 is empty or doesn't include a key
						list2 := strings.Fields(v1) // Split strings seprated by spaces
						if len(list2) > 0 {
							v2 := list2[0]                                                                   // Supposedly the secret key
							purgedsubstr := v2[0:2] + strings.Repeat("X", len(v2)-7) + v2[len(v2)-5:len(v2)] // Keep the 1st 2 and last 4 digits and ]
							purged += (token + purgedsubstr)
							for i := 1; i < len(list2); i++ { // Add the remaining of the substrings
								purged += (" " + list2[i])
							}
						} else {
							purged += token
						}
					} else {
						purged += v1
					}
				}
			} else {
				purged = msg
			}
			fmt.Fprintln(crowdin.logWriter, purged)
		} else {
			log.Println(a)
		}
	}
}

// APIError holds data of errors returned from the API.
type APIError struct {
	What string
}

func (e APIError) Error() string {
	return fmt.Sprintf("%v", e.What)
}
