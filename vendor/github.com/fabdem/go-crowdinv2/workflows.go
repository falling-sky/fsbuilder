package crowdin

import (
	"encoding/json"
	// "errors"
	"fmt"
	"strconv"
)

// ListWorkflowsSteps - List workflow steps
// {protocol}://{host}/api/v2/projects/{projectId}/workflow-steps
//
// 
// 
func (crowdin *Crowdin) ListWorkflowsSteps(options *ListWorkflowsStepsOptions) (*ResponseListWorkflowsSteps, error) {
	crowdin.log(fmt.Sprintf("ListWorkflowsSteps(%d)\n", crowdin.config.projectId))

	var limit string
	if options.Limit > 0 {
		limit = strconv.Itoa(options.Limit)
	}

	var offset string
	if options.Offset > 0 {
		offset = strconv.Itoa(options.Offset)
	}

	response, err := crowdin.get(&getOptions{
		urlStr: fmt.Sprintf(crowdin.config.apiBaseURL+"projects/%v/workflow-steps", crowdin.config.projectId),
		params: map[string]string{
			"limit":  limit,
			"offset": offset,
		},
	})

	if err != nil {
		crowdin.log(fmt.Sprintf("	Error - response:%s\n%s\n", response, err))
		return nil, err
	}

	var responseAPI ResponseListWorkflowsSteps
	err = json.Unmarshal(response, &responseAPI)
	if err != nil {
		crowdin.log(fmt.Sprintf("	Error - unmarshalling:%s\n%s\n", response, err))
		return nil, err
	}

	return &responseAPI, nil

}
