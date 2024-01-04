package crowdin

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// ListProjects - List projects API call. List the projects and their respective details (incl.Id.)
// {protocol}://{host}/api/v2/projects
func (crowdin *Crowdin) ListProjects(options *ListProjectsOptions) (*ResponseListProjects, error) {

	crowdin.log(fmt.Sprintf("ListProjects()"))

	var groupId string
	if options.GroupId > 0 {
		groupId = strconv.Itoa(options.GroupId)
	}

	var hasManagerAccess string
	if options.HasManagerAccess > 0 {
		hasManagerAccess = strconv.Itoa(options.HasManagerAccess)
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
		urlStr: fmt.Sprintf(crowdin.config.apiBaseURL + "projects"),
		params: map[string]string{
			"groupId":          groupId,
			"hasManagerAccess": hasManagerAccess,
			"limit":            limit,
			"offset":           offset,
		},
	})

	if err != nil {
		crowdin.log(err)
		return nil, err
	}

	crowdin.log(string(response))

	var responseAPI ResponseListProjects
	err = json.Unmarshal(response, &responseAPI)
	if err != nil {
		crowdin.log(err)
		return nil, err
	}

	return &responseAPI, nil
}


// GetProject - Read project details
// {protocol}://{host}/api/v2/projects/{projectId}
func (crowdin *Crowdin) GetProject() (*ResponseGetProject, error) {
	crowdin.log(fmt.Sprintf("GetProject(%d)", crowdin.config.projectId))

	response, err := crowdin.get(&getOptions{
		urlStr: fmt.Sprintf(crowdin.config.apiBaseURL + "projects/%v", crowdin.config.projectId),
		},
	)

	if err != nil {
		crowdin.log(err)
		return nil, err
	}
	
	var responseAPI ResponseGetProject
	err = json.Unmarshal(response, &responseAPI)
	if err != nil {
		crowdin.log(err)
		return nil, err
	}

	return &responseAPI, nil
}
