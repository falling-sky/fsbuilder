package crowdin

import (
	"encoding/json"
	"errors"
	"fmt"
	// "io"
	// "net/http"
	// "net/url"
	//"os"
	"strconv"
	// "time"
	// "github.com/mreiferson/go-httpclient"
)

// EditStrings - Edit Source Strings
// {protocol}://{host}/api/v2/projects/{projectId}/strings/{stringId}
//
// Validate EditStringOptions.Value type to prevent panic
// but relies on the API for the validation of the other parameters.

func (crowdin *Crowdin) EditStrings(options *EditStringsOptions, stringId int) (*ResponseEditStrings, error) {

	crowdin.log(fmt.Sprintf("EditString()\n"))

	if len(*options) > 0 { // Need at least 1 set of parameters
		// Check that the interface underlying type is string, int or boolean.
		for _, val := range *options {
			switch t := val.Value.(type) {
			case bool:
			case int:
			case string:
			default:
				crowdin.log(fmt.Sprintf("	Error - param type not allowed:%v\n", t))
				return nil, errors.New("Parameters type not allowed.")
			}
		}
	} else { // No params?!
		crowdin.log(fmt.Sprintf("	Error - at least one set of parameters is needed\n"))
		return nil, errors.New("No parameters found.")
	}

	response, err := crowdin.patch(&patchOptions{
		urlStr: fmt.Sprintf(crowdin.config.apiBaseURL+"projects/%v/strings/%v", crowdin.config.projectId, stringId),
		body:   options})

	if err != nil {
		crowdin.log(fmt.Sprintf("	Error - response:%s\n%s\n", response, err))
		return nil, err
	}

	var responseAPI ResponseEditStrings
	err = json.Unmarshal(response, &responseAPI)
	if err != nil {
		crowdin.log(fmt.Sprintf("	Error - unmarshalling:%s\n%s\n", response, err))
		return nil, err
	}

	return &responseAPI, nil

}

// ListStrings - List Source Strings
// {protocol}://{host}/api/v2/projects/{projectId}/strings
func (crowdin *Crowdin) ListStrings(options *ListStringsOptions) (*ResponseListStrings, error) {

	crowdin.log(fmt.Sprintf("ListDirectories()\n"))

	var fileId string
	if options.FileId > 0 {
		fileId = strconv.Itoa(options.FileId)
	}

	var denormalizePlaceholders string
	if options.DenormalizePlaceholders > 0 {
		denormalizePlaceholders = strconv.Itoa(options.DenormalizePlaceholders)
	}

	var limit string
	if options.Limit > 0 {
		limit = strconv.Itoa(options.Limit)
	}

	var offset string
	if options.Offset > 0 {
		offset = strconv.Itoa(options.Offset)
	}

	response, err := crowdin.get(&getOptions{
		urlStr: fmt.Sprintf(crowdin.config.apiBaseURL+"projects/%v/strings", crowdin.config.projectId),
		params: map[string]string{
			"fileId":                  fileId,
			"denormalizePlaceholders": denormalizePlaceholders,
			"labelIds":                options.LabelIds,
			"filter":                  options.Filter,
			"scope":                   options.Scope,
			"limit":                   limit,
			"offset":                  offset,
		},
	})

	if err != nil {
		crowdin.log(fmt.Sprintf("	Error - response:%s\n%s\n", response, err))
		return nil, err
	}

	var responseAPI ResponseListStrings
	err = json.Unmarshal(response, &responseAPI)
	if err != nil {
		crowdin.log(fmt.Sprintf("	Error - unmarshalling:%s\n%s\n", response, err))
		return nil, err
	}

	return &responseAPI, nil
}
