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

// ListDirectories - List directories in a given project
// {protocol}://{host}/api/v2/projects/{projectId}/files
func (crowdin *Crowdin) ListDirectories(options *ListDirectoriesOptions) (*ResponseListDirectories, error) {

	crowdin.log(fmt.Sprintf("ListDirectories()\n"))

	var branchId string
	if options.BranchId > 0 {
		branchId = strconv.Itoa(options.BranchId)
	}

	var directoryId string
	if options.DirectoryId > 0 {
		directoryId = strconv.Itoa(options.DirectoryId)
	}

	var recursion string
	if options.Recursion > 0 {
		recursion = strconv.Itoa(options.Recursion)
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
		urlStr: fmt.Sprintf(crowdin.config.apiBaseURL+"projects/%v/directories", crowdin.config.projectId),
		params: map[string]string{
			"branchId":    branchId,
			"directoryId": directoryId,
			"recursion":   recursion,
			"limit":       limit,
			"offset":      offset,
		},
	})

	if err != nil {
		crowdin.log(fmt.Sprintf("	Error - response:%s\n%s\n", response, err))
		return nil, err
	}

	var responseAPI ResponseListDirectories
	err = json.Unmarshal(response, &responseAPI)
	if err != nil {
		crowdin.log(fmt.Sprintf("	Error - unmarshalling:%s\n%s\n", response, err))
		return nil, err
	}

	return &responseAPI, nil
}


// ListAllDirectories - Helper function: list all directories in a given project (all pages)
// Ignore offset and limit in options.
//
func (crowdin *Crowdin) ListAllDirectories(options *ListDirectoriesOptions) (*ResponseListDirectories, error) {

	crowdin.log(fmt.Sprintf("ListAllDirectories()\n"))

	limit := MAX_RES_PER_PAGE // nb max results returned by call per page.
	page := 0
	var listDirs ResponseListDirectories
	for offset := 0; offset < MAX_RESULTS; offset += limit {
		lst, err := crowdin.ListDirectories(&ListDirectoriesOptions{Offset: offset, Limit: limit})
		if err != nil {
			return &listDirs, errors.New(fmt.Sprintf("LookupFileId() - Error listing project directories. Page %d", page))
		}
		
		if len(lst.Data) <= 0 {  // Reached the end
			break
		}
		
		page++
		listDirs.Data = append(listDirs.Data, lst.Data...)

		crowdin.log(fmt.Sprintf(" - Page of results #%d\n", page))
	}

	return &listDirs, nil
}


// ListFiles - List files in a given project
// {protocol}://{host}/api/v2/projects/{projectId}/files
func (crowdin *Crowdin) ListFiles(options *ListFilesOptions) (*ResponseListFiles, error) {

	crowdin.log(fmt.Sprintf("ListFiles()\n"))

	var branchId string
	if options.BranchId > 0 {
		branchId = strconv.Itoa(options.BranchId)
	}

	var directoryId string
	if options.DirectoryId > 0 {
		directoryId = strconv.Itoa(options.DirectoryId)
	}

	var recursion string
	if options.Recursion > 0 {
		recursion = strconv.Itoa(options.Recursion)
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
		urlStr: fmt.Sprintf(crowdin.config.apiBaseURL+"projects/%v/files", crowdin.config.projectId),
		params: map[string]string{
			"branchId":    branchId,
			"directoryId": directoryId,
			"recursion":   recursion,
			"limit":       limit,
			"offset":      offset,
		},
	})

	if err != nil {
		crowdin.log(fmt.Sprintf("	Error - response:%s\n%s\n", response, err))
		return nil, err
	}

	var responseAPI ResponseListFiles
	err = json.Unmarshal(response, &responseAPI)
	if err != nil {
		crowdin.log(fmt.Sprintf("	Error - unmarshalling:%s\n%s\n", response, err))
		return nil, err
	}

	return &responseAPI, nil
}

// ListFileRevisions - List all revisions for a file in current project
// {protocol}://{host}/api/v2/projects/{projectId}/files/{fileId}/revisions
func (crowdin *Crowdin) ListFileRevisions(options *ListFileRevisionsOptions, fileId int) (*ResponseListFileRevisions, error) {

	crowdin.log(fmt.Sprintf("ListFileRevisions()\n"))

	var limit string
	if options.Limit > 0 {
		limit = strconv.Itoa(options.Limit)
	}

	var offset string
	if options.Offset > 0 {
		offset = strconv.Itoa(options.Offset)
	}

	response, err := crowdin.get(&getOptions{
		urlStr: fmt.Sprintf(crowdin.config.apiBaseURL+"projects/%v/files/%v/revisions", crowdin.config.projectId, fileId),
		params: map[string]string{
			"limit":  limit,
			"offset": offset,
		},
	})
	if err != nil {
		crowdin.log(fmt.Sprintf("	Error - response:%s\n%s\n", response, err))
		return nil, err
	}

	var responseAPI ResponseListFileRevisions
	err = json.Unmarshal(response, &responseAPI)
	if err != nil {
		crowdin.log(fmt.Sprintf("	Error - unmarshalling:%s\n%s\n", response, err))
		return nil, err
	}
	crowdin.log(fmt.Sprintf("	Unmarshalled:%s\n", response))

	return &responseAPI, nil
}

// GetFileRevision - List a specific revision details for a file in current project
// {protocol}://{host}/api/v2/projects/{projectId}/files/{fileId}/revisions/{revisionId}
func (crowdin *Crowdin) GetFileRevision(fileId int, revId int) (*ResponseGetFileRevision, error) {

	crowdin.log(fmt.Sprintf("GetFileRevision()\n"))

	response, err := crowdin.get(&getOptions{
		urlStr: fmt.Sprintf(crowdin.config.apiBaseURL+"projects/%v/files/%v/revisions/%v", crowdin.config.projectId, fileId, revId),
	})
	if err != nil {
		crowdin.log(fmt.Sprintf("	Error - response:%s\n%s\n", response, err))
		return nil, err
	}

	var responseAPI ResponseGetFileRevision
	err = json.Unmarshal(response, &responseAPI)
	if err != nil {
		crowdin.log(fmt.Sprintf("	Error - unmarshalling:%s\n%s\n", response, err))
		return nil, err
	}
	crowdin.log(fmt.Sprintf("	Unmarshalled:%s\n", response))

	return &responseAPI, nil
}

// UpdateFile - Update a specific file
// {protocol}://{host}/api/v2/projects/{projectId}/files/{fileId}
// Default update mode is explicitely clear_translations_and_approvals
func (crowdin *Crowdin) UpdateFile(fileId int, options *UpdateFileOptions) (*ResponseUpdateFile, error) {

	crowdin.log(fmt.Sprintf("UpdateFile()\n"))

	if len(options.UpdateOption) > 0 {
		// Check that update options are valid
		if !(options.UpdateOption == "clear_translations_and_approvals" || options.UpdateOption == "keep_translations" || options.UpdateOption == "keep_translations_and_approvals") {
			crowdin.log(fmt.Sprintf("	Error - Update Option is not valid:%s\n", options.UpdateOption))
			return nil, errors.New("Invalid update option.")
		}
	} else {
		options.UpdateOption = "clear_translations_and_approvals" // Default behavior
	}

	response, err := crowdin.put(&putOptions{urlStr: fmt.Sprintf(crowdin.config.apiBaseURL+"projects/%v/files/%v", crowdin.config.projectId, fileId), body: options})

	if err != nil {
		crowdin.log(fmt.Sprintf("	Error - response:%s\n%s\n", response, err))
		return nil, err
	}

	var responseAPI ResponseUpdateFile
	err = json.Unmarshal(response, &responseAPI)
	if err != nil {
		crowdin.log(fmt.Sprintf("	Error - unmarshalling:%s\n%s\n", response, err))
		return nil, err
	}

	return &responseAPI, nil
}
