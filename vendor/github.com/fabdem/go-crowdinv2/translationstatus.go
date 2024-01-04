package crowdin

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// GetFileProgress() - Get file progress info
// {protocol}://{host}/api/v2/projects/{projectId}/files/{fileId}/languages/progress
func (crowdin *Crowdin) GetFileProgress(options *GetFileProgressOptions) (*ResponseGetFileProgress, error) {
	crowdin.log(fmt.Sprintf("GetFileProgress()\n"))

	var limit string
	if options.Limit > 0 {
		limit = strconv.Itoa(options.Limit)
	}

	var offset string
	if options.Offset > 0 {
		offset = strconv.Itoa(options.Offset)
	}

	response, err := crowdin.get(&getOptions{
		urlStr: fmt.Sprintf(crowdin.config.apiBaseURL+"projects/%v/files/%v/languages/progress", crowdin.config.projectId, options.FileId),
		params: map[string]string{
			"limit":  limit,
			"offset": offset,
		},
	})

	if err != nil {
		crowdin.log(err)
		return nil, err
	}

	crowdin.log(string(response))

	var responseAPI ResponseGetFileProgress
	err = json.Unmarshal(response, &responseAPI)
	if err != nil {
		crowdin.log(err)
		return nil, err
	}

	return &responseAPI, nil
}
