package crowdin

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// Publicly available high level functions generally combining several API calls

const polldelaysec = 5 // Defines delay in seconds between each api call when polling a progress status

// Lookup buildId for current project
func (crowdin *Crowdin) GetBuildId() (buildId int, err error) {

	crowdin.log("GetBuildId()")

	var opt ListProjectBuildsOptions
	rl, err := crowdin.ListProjectBuilds(&opt)
	if err != nil {
		return 0, err
	}
	for _, v := range rl.Data {
		if (v.Data.ProjectId == crowdin.config.projectId) && (v.Data.Status == "finished") {
			buildId = v.Data.Id
		}
	}
	if buildId == 0 {
		return 0, errors.New("Can't find a build for this project or build is in progress.")
	}
	return buildId, nil
}

// Lookup projectId
func (crowdin *Crowdin) GetProjectId(projectName string) (projectId int, err error) {

	crowdin.log("GetProjectId()")

	var opt ListProjectsOptions
	rl, err := crowdin.ListProjects(&opt)
	if err != nil {
		return 0, err
	}

	for _, v := range rl.Data {
		if v.Data.Name == projectName {
			projectId = v.Data.ID
		}
	}
	if projectId == 0 {
		return 0, errors.New("Can't find project.")
	}
	return projectId, nil
}


// BuildTranslationsAllLg - Build transaltions for a project or folder for all languages
// Options to export:
//   - project Id and optionaly folder Id
//   - translated strings only Y/N
//   - approved strings only integer
//		Enterprise: min nb of approval steps required to export a string
//		crowdin.com: 0 means approval not required, diff from 0: approval required 
//   - fully translated files only Y/N
//   - only strings that have passed their workflow Y/N
//	"translated strings only" and fully "translated files only" are mutually exclusive.
// Update buildProgress

type BuildTranslationAllLgOptions struct {
	BuildTO						time.Duration
	TranslatedOnly				bool
	MinApprovalSteps 			int
	FullyTranslatedFilesOnly	bool
	ExportStringsThatPassedWkfl	bool
	FolderName					string		// Optional - if not empty and valid, build a folder
}

func (crowdin *Crowdin) BuildTranslationAllLg(opt BuildTranslationAllLgOptions) (buildId int, err error) {
	crowdin.log("BuildTranslationAllLg()")

	// Invoke build
	var dirId int
	var status string
	
	if opt.TranslatedOnly && opt.FullyTranslatedFilesOnly {
		return buildId, errors.New("\nOption error - Can't have both TranslatedOnly and FullyTranslatedFilesOnly set to true.")
	}

	// Look up the dir Id  if we need to do a folder build
	if len(opt.FolderName) > 0 {
		// ---- DIRECTORY BUILD ------
		crowdin.log(fmt.Sprintf("Building folder: %s",opt.FolderName))
		dirId, _, err = crowdin.LookupDirId(opt.FolderName)
		if err != nil {
			return buildId, err
		}			
		var bo BuildDirectoryTranslationOptions
		bo.TargetLanguageIds = nil
		bo.SkipUntranslatedFiles = opt.FullyTranslatedFilesOnly
		bo.SkipUntranslatedStrings = opt.TranslatedOnly
		bo.ExportStringsThatPassedWorkflow = opt.ExportStringsThatPassedWkfl
		if crowdin.config.apiBaseURL == API_CROWDINDOTCOM {
			bo.ExportApprovedOnly = (opt.MinApprovalSteps != 0) // crowdin.com
		} else {
			bo.ExportWithMinApprovalsCount = opt.MinApprovalSteps // Enterprise
		}
		rb, err := crowdin.BuildDirectoryTranslation(dirId, &bo)
		if err != nil {
			return buildId, errors.New("\nBuild Err.")
		}
		buildId = rb.Data.ID
		status = rb.Data.Status

	} else {
		// ---- PROJECT BUILD ------
		crowdin.log(fmt.Sprintf("Building project: %d", crowdin.config.projectId))
		var bo BuildProjectTranslationOptions
		bo.Languages = nil
		bo.SkipUntranslatedFiles = opt.FullyTranslatedFilesOnly
		bo.SkipUntranslatedStrings = opt.TranslatedOnly
		bo.ExportStringsThatPassedWorkflow = opt.ExportStringsThatPassedWkfl
		if crowdin.config.apiBaseURL == API_CROWDINDOTCOM {
			bo.ExportApprovedOnly = (opt.MinApprovalSteps != 0) // crowdin.com
		} else {
			bo.ExportWithMinApprovalsCount = opt.MinApprovalSteps // Enterprise
		}
		rb, err := crowdin.BuildProjectTranslation(&bo)
		if err != nil {
			return buildId, errors.New("\nBuild Err.")
		}
		buildId = rb.Data.Id
		status = rb.Data.Status
	}

	crowdin.log(fmt.Sprintf("	BuildId=%d", buildId))

	// Poll build status with a timeout
	crowdin.log("	Poll build status crowdin.CheckProjectBuildStatus()")
	timer := time.NewTimer(opt.BuildTO)
	defer timer.Stop()
	rp := &ResponseCheckProjectBuildStatus{}
	for rp.Data.Status = status; rp.Data.Status == "inProgress"; { // initial value is read from previous API call
		time.Sleep(polldelaysec * time.Second) // delay between each call
		rp, err = crowdin.CheckProjectBuildStatus(&CheckProjectBuildStatusOptions{BuildId: buildId})
		if err != nil {
			crowdin.log(fmt.Sprintf(" Error CheckProjectBuildStatus()=%s", err))
			return buildId, err
			// break
		}
		select {
		case <-timer.C:
			err = errors.New("Build Timeout.")
			return buildId, err
		default:
		}
	}

	if rp.Data.Status != "finished" {
		err = errors.New(fmt.Sprintf("	Build Error:%s", rp.Data.Status))
	}
	return buildId, err
}
	

// BuildAllLg - Build a project for all languages. Kept to maintain compatibility with older versions.
// Options to export:
//   - translated strings only Y/N
//   - approved strings only integer
//		Enterprise: min nb of approval steps required to export a string
//		crowdin.com: 0 means approval not required, diff from 0: approval required 
//   - fully translated files only Y/N
//	"translated strings only" and fully "translated files only" are mutually exclusive.
// Update buildProgress
func (crowdin *Crowdin) BuildAllLg(buildTO time.Duration, translatedOnly bool, minApprovalSteps int, fullyTranslatedFilesOnly bool) (buildId int, err error) {
	crowdin.log("BuildAllLg()")

	opt := BuildTranslationAllLgOptions{
		BuildTO	:				    buildTO,
		TranslatedOnly:				translatedOnly,
		MinApprovalSteps: 			minApprovalSteps,
		FullyTranslatedFilesOnly:	fullyTranslatedFilesOnly,
		FolderName:					"",
	}

	buildId, err =  crowdin.BuildTranslationAllLg(opt)
	if err != nil {
		return buildId, err
	}

	return buildId, err
}

// Download a build of the current project
//    outputFileNamePath  required
//    projectId           required if projectName is not provided
//    buildId             optional
// limitation: total number of project directories needs to be 500 max
func (crowdin *Crowdin) DownloadBuild(outputFileNamePath string, buildId int) (err error) {

	// Get URL for downloading
	rd, err := crowdin.DownloadProjectTranslations(&DownloadProjectTranslationsOptions{buildId})
	if err != nil {
		return errors.New("DownloadBuild() - Error getting URL for download.")
	}
	url := rd.Data.Url

	// Actual downloading
	err = crowdin.DownloadFile(url, outputFileNamePath)

	return err
}

// Lookup fileId in current project
//    CrowdinFileName required - full Crowdin path to file.
//		Returns Id and crowdin file name
func (crowdin *Crowdin) LookupFileId(CrowdinFileName string) (id int, name string, err error) {

	crowdin.log(fmt.Sprintf("LookupFileId()\n"))

	// Lookup fileId in Crowdin
	dirId := 0
	crowdinFile := strings.Split(CrowdinFileName, "/")

	crowdin.log(fmt.Sprintf("  len=%d\n", len(crowdinFile)))
	crowdin.log(fmt.Sprintf("  crowdinFile %v\n", crowdinFile))
	// crowdin.log(fmt.Sprintf("  crowdinFile[1] %s\n", crowdinFile[1] ))

	switch l := len(crowdinFile); l {
	case 0:
		return 0, "", errors.New("LookupFileId() - Crowdin file name should not be null.")
	case 1: // no directory so dirId is 0 - value is like "a_file_name"
	case 2: // no directory so dirId is 0 - value is like "/a_file_name"
	default: // l > 1
		// Lookup end directoryId
		// Get a list of all the project folders
		listDirs, err := crowdin.ListAllDirectories(&ListDirectoriesOptions{})
		if err != nil {
			return 0, "", errors.New("LookupFileId() - Error listing project directories.")
		}
		if len(listDirs.Data) > 0 {
			// Lookup last directory's Id
			dirId = 0
			for i, dirName := range crowdinFile { // Go down the directory branch
				crowdin.log(fmt.Sprintf("  idx %d dirName %s len %d dirId %d", i, dirName, len(crowdinFile), dirId))
				if i > 0 && i < len(crowdinFile)-1 { // 1st entry is empty and we're done once we reach the file name (last item of the slice).
					for _, crwdPrjctDirName := range listDirs.Data { // Look up in list of project dirs the right one
						crowdin.log(fmt.Sprintf("  check -> crwdPrjctDirName.Data.DirectoryId %d crwdPrjctDirName.Data.Name %s", crwdPrjctDirName.Data.DirectoryId, crwdPrjctDirName.Data.Name))
						if crwdPrjctDirName.Data.DirectoryId == dirId && crwdPrjctDirName.Data.Name == dirName {
							dirId = crwdPrjctDirName.Data.Id // Bingo get that Id
							crowdin.log(fmt.Sprintf("  BINGO dirId=%d Crowdin dir name %s", dirId, crwdPrjctDirName.Data.Name))
							break // Done for that one
						}
					}
					if dirId == 0 {
						return 0, "", errors.New(fmt.Sprintf("LookupFileId() - Error: can't match directory names with Crowdin path."))
					}
				}
			}
			if dirId == 0 {
				return 0, "", errors.New(fmt.Sprintf("LookupFileId() - Error: can't match directory names with Crowdin path."))
			}
		} else {
			return 0, "", errors.New("LookupFileId() - Error: mismatch between # of folder found and # of folder expected.")
		}
	}

	crowdinFilename := crowdinFile[len(crowdinFile)-1] // Get file name
	crowdin.log(fmt.Sprintf("  crowdinFilename %s\n", crowdinFilename))

	// Look up file
	listFiles, err := crowdin.ListFiles(&ListFilesOptions{DirectoryId: dirId, Limit: 500})
	if err != nil {
		return 0, "", errors.New("LookupFileId() - Error listing files.")
	}

	fileId := 0
	for _, list := range listFiles.Data {
		crowdin.log(fmt.Sprintf("  check -> list.Data.Name %s", list.Data.Name))
		if list.Data.Name == crowdinFilename {
			fileId = list.Data.Id
			crowdin.log(fmt.Sprintf("  BINGO fileId=%d File name %s", fileId, crowdinFilename))
			break // found it
		}
	}

	if fileId == 0 {
		return 0, "", errors.New(fmt.Sprintf("LookupFileId() - Can't find file %s in Crowdin.", crowdinFilename))
	}

	crowdin.log(fmt.Sprintf("  fileId=%d\n", fileId))
	return fileId, crowdinFilename, nil
}


// Lookup directoryId in current project
//    CrowdinDirName required - full Crowdin path to directory.
//		Returns Id and crowdin dir name
func (crowdin *Crowdin) LookupDirId(CrowdinDirName string) (id int, name string, err error) {

	crowdin.log(fmt.Sprintf("LookupDirId()\n"))

	// Lookup directoryId in Crowdin
	dirId := 0
	crowdinDirs := strings.Split(CrowdinDirName, "/")

	crowdin.log(fmt.Sprintf("  len=%d\n", len(crowdinDirs)))
	crowdin.log(fmt.Sprintf("  crowdinDir %v\n", crowdinDirs))
	// crowdin.log(fmt.Sprintf("  crowdinDirs[1] %s\n", crowdinDirs[1] ))

	switch l := len(crowdinDirs); l {
	case 0:
		return 0, "", errors.New("LookupDirId() - Crowdin directory name should not be null.")
	// case 1: // no directory so dirId is 0 - value is like "a_file_name"
	// case 2: // no directory so dirId is 0 - value is like "/a_file_name"
	default: // l >= 1
		// Lookup end directoryId
		// Get a list of all the project's directories
		listDirs, err := crowdin.ListAllDirectories(&ListDirectoriesOptions{})
		if err != nil {
			return 0, "", errors.New("LookupDirId() - Error listing project directories.")
		}
		
		if len(listDirs.Data) > 0 {
			// Lookup last directory's Id
			dirId = 0
			for i, dirName := range crowdinDirs { // Go down the directory branch
				crowdin.log(fmt.Sprintf("  idx %d dirName %s len %d dirId %d", i, dirName, len(crowdinDirs), dirId))
				found := false
				if i > 0 && i < len(crowdinDirs) { // 1st entry is empty and we're done once we reach the last directory name.
					for _, crwdPrjctDirName := range listDirs.Data { // Look up in list of project dirs the right one
						crowdin.log(fmt.Sprintf("  check -> crwdPrjctDirName.Data.DirectoryId %d crwdPrjctDirName.Data.Name |%s|", crwdPrjctDirName.Data.DirectoryId, crwdPrjctDirName.Data.Name))
						if crwdPrjctDirName.Data.DirectoryId == dirId && crwdPrjctDirName.Data.Name == dirName {
							dirId = crwdPrjctDirName.Data.Id // Bingo get that Id
							name = dirName
							found = true
							crowdin.log(fmt.Sprintf("  BINGO dirId=%d Crowdin dir name %s", dirId, crwdPrjctDirName.Data.Name))
							break // Done for that one
						}
					}
					if !found {
						return 0, "", errors.New(fmt.Sprintf("LookupDirId() - Error: can't match directory names with Crowdin path."))
					}
				}
			}
		} else {
			return 0, "", errors.New("LookupDirId() - Error: mismatch between # of folder found and # of folder expected.")
		}
	}

	return dirId, name, nil
}

// Update a file of the current project
//    LocalFileName  required
//    CrowdinFileName required
//    updateOption required needs to be either: clear_translations_and_approvals, keep_translations or keep_translations_and_approvals
//		Returns file Id and rev
func (crowdin *Crowdin) Update(CrowdinFileName string, LocalFileName string, updateOption string) (fileId int, revId int, err error) {

	crowdin.log(fmt.Sprintf("Update()\n"))

	// Lookup fileId in Crowdin
	fileId, crowdinFilename, err := crowdin.LookupFileId(CrowdinFileName)
	if err != nil {
		crowdin.log(fmt.Sprintf("  err=%s\n", err))
		return 0, 0, err
	}

	crowdin.log(fmt.Sprintf("Update() fileId=%d fileName=%s\n", fileId, crowdinFilename))

	// Send local file to storageId
	addStor, err := crowdin.AddStorage(&AddStorageOptions{FileName: LocalFileName})
	if err != nil {
		return 0, 0, errors.New("Update() - Error adding file to storage.")
	}
	storageId := addStor.Data.Id

	// fmt.Printf("Directory Id = %d, filename= %s, fileId %d storageId= %d\n", dirId, crowdinFilename, fileId, storageId)

	// Update file
	updres, err := crowdin.UpdateFile(fileId, &UpdateFileOptions{StorageId: storageId, UpdateOption: updateOption})

	// Delete storage
	err1 := crowdin.DeleteStorage(&DeleteStorageOptions{StorageId: storageId})

	if err != nil {
		crowdin.log(fmt.Sprintf("Update() - error updating file %v", updres))
		return 0, 0, errors.New("Update() - Error updating file.") //
	}

	if err1 != nil {
		crowdin.log(fmt.Sprintf("Update() - error deleting storage %v", err1))
	}

	revId = updres.Data.RevisionId

	crowdin.log(fmt.Sprintf("Update() - result %v", updres))

	return fileId, revId, nil
}

// Obtain a list of string Ids for a given file of the current project.
// Use a filter on "identifier" "text" or "context"
// Parameters:
//  - provide path/filename
//	- a filter string (empty mean "all")
//	- filter on "identifier" "text" or "context"
// Returns:
//	- string IDs in a slice of ints if results found
//	- err (nil if no error)
//
func (crowdin *Crowdin) GetStringIDs(fileName string, filter string, filterType string) (list []int, err error) {

	crowdin.log(fmt.Sprintf("GetStringIDs(%s, %s, %s)\n", fileName, filter, filterType))

	// Lookup fileId in Crowdin
	fileId, _, err := crowdin.LookupFileId(fileName)
	if err != nil {
		crowdin.log(fmt.Sprintf("  err=%s\n", err))
		return list, err
	}

	// Get the string IDs
	limit := 500
	opt := ListStringsOptions{
		FileId: fileId,
		Scope:  filterType,
		Filter: filter,
		Limit:  limit,
	}

	// Pull ListStrings as long as it returns data
	for offset := 0; offset < MAX_RESULTS; offset += limit {
		opt.Offset = offset

		res, err := crowdin.ListStrings(&opt)
		if err != nil {
			crowdin.log(fmt.Sprintf("  err=%s\n", err))
			return list, err
		}

		if len(res.Data) <= 0 {
			break
		}

		crowdin.log(fmt.Sprintf(" - Page of results #%d\n", (offset/limit)+1))

		for _, v := range res.Data {
			list = append(list, v.Data.ID) // Add data to slice
		}
	}

	return list, nil
}

type T_UploadTranslationFileParams struct {
	LocalFileName       string           // File containing the translations to upload
	CrowdinFileName     string           // File in Crowdin where the translations will end up
	LanguageId          string           // Langugage ID as per Crowdin spec and defined as target in the project
	ImportEqSuggestions bool             // Defines whether to add translation if it's the same as the source string
	AutoApproveImported bool             // Mark uploaded translations as approved
	TranslateHidden     bool             // Allow translations upload to hidden source strings
	ResponseTimeOut     time.Duration    // in seconds. The upload operation can take several minutes.
	// The original TO will be restored after operation finishes (ok or not)
}

// Upload a translation file
// Params:
// 	- File containing the translations to upload
// 	- File in Crowdin where the translations will end up
// 	- Language ID as per Crowdin spec and defined as target in the project
// 	- Defines whether to add translation if it's the same as the source string
// 	- Mark uploaded translations as approved
// 	- Allow translations upload to hidden source strings
// 	- in seconds. The upload operation can take several minutes. 0 means no change.
// 		The original TO will be restored after operation finishes (ok or not)
//	Returns the source fileId (0 if error) and err != nil if error
func (crowdin *Crowdin) UploadTranslationFile(params T_UploadTranslationFileParams) (fileId int, err error) {
	crowdin.log(fmt.Sprintf("UploadTranslationFile(%v)\n", params))

	// Lookup fileId in Crowdin
	fileId, crowdinFilename, err := crowdin.LookupFileId(params.CrowdinFileName)
	if err != nil {
		crowdin.log(fmt.Sprintf("  err=%s\n", err))
		return fileId, err
	}

	crowdin.log(fmt.Sprintf("UploadTranslationFile() fileId=%d fileName=%s\n", fileId, crowdinFilename))

	// Send local file to storageId
	addStor, err := crowdin.AddStorage(&AddStorageOptions{FileName: params.LocalFileName})
	if err != nil {
		crowdin.log(fmt.Sprintf("  Error adding file to storage %s\n", err))
		return fileId, errors.New("UploadTranslationFile() - Error adding file to storage.")
	}
	storageId := addStor.Data.Id

	// fmt.Printf("Directory Id = %d, filename= %s, fileId %d storageId= %d\n", dirId, crowdinFilename, fileId, storageId)

	// Upload file
	if params.ResponseTimeOut > 0 { // If a specific to has been defined
		crowdin.PushTimeouts()                         //  Backup comm timeouts
		crowdin.SetTimeouts(0, params.ResponseTimeOut) // Set new TO for this call
	}
	upldres, err := crowdin.UploadTranslations(params.LanguageId,
		&UploadTranslationsOptions{
			StorageID:           storageId,
			FileID:              fileId,
			ImportEqSuggestions: params.ImportEqSuggestions,
			AutoApproveImported: params.AutoApproveImported,
			TranslateHidden:     params.TranslateHidden,
		})
	if params.ResponseTimeOut > 0 { // If a specific to has been defined
		crowdin.PopTimeouts() // Restore current timeouts
	}

	// Delete storage
	err1 := crowdin.DeleteStorage(&DeleteStorageOptions{StorageId: storageId})

	crowdin.log(fmt.Sprintf("UploadTranslationFile() - uploading %s result %v\n", params.LocalFileName, upldres))
	if err != nil {
		crowdin.log(fmt.Sprintf("UploadTranslationFile() - upload error - %s", err))
		return fileId, errors.New("UploadTranslationFile() - Error uploading file.")
	}

	if err1 != nil {
		// Not a fatal err, just log the error
		crowdin.log(fmt.Sprintf("UploadTranslationFile() - error deleting storage %v", err1))
	}

	return fileId, nil
}

// GetShortLangFileProgress() - Get a simple file completion info for a specific language
//	 Returns a percentage of completion for both translation and approval (0 if error).
func (crowdin *Crowdin) GetShortLangFileProgress(fileId int, langId string) (translationProgress int, approvalProgress int, err error) {
	crowdin.log(fmt.Sprintf("GetShortLangFileProgress()\n"))

	opt := GetFileProgressOptions{FileId: fileId, Limit: 500}
	res, err := crowdin.GetFileProgress(&opt)
	if err == nil {
		// Lookup for language in res
		for _, v := range res.Data {
			if v.Data.LanguageId == langId {
				return v.Data.TranslationProgress, v.Data.ApprovalProgress, nil // found it: done
			}
		}
		crowdin.log(fmt.Sprintf("GetShortLangFileProgress() - language %s not found in %v", langId, res))
		err = errors.New("GetShortLangFileProgress() - Can't find language in result.")
	}
	return 0, 0, err

}

// Get steps of all approved transactions from a given project.
//    
//		Returns a map of all ADMIN approved transactions along with their approval steps.
//
func (crowdin *Crowdin) GetStepsApprovTransId() (listTrans map[int][]int, err error) {
	crowdin.log(fmt.Sprintf("GetStepsApprovTransId() %d\n", crowdin.config.projectId))

	// Get the project language IDs
	listProjDetails, err := crowdin.GetProject()
	if err != nil {
		fmt.Printf("ERREUR: %s\n", err)
		return listTrans, err
	}
	var targetLanguageIDs []string
	targetLanguageIDs = listProjDetails.Data.TargetLanguageIds
	crowdin.log(fmt.Sprintf("Target language Ids: %s\n", targetLanguageIDs))

	// Get all the file ids from a project
	listFiles, err := crowdin.ListFiles(&ListFilesOptions{Limit: 500})
	if err != nil {
		fmt.Printf("ERREUR: %s\n", err)
		return listTrans, err
	}

	crowdin.log(fmt.Sprintf("Files to process:\n"))
	fileIDs := make(map[string]int)
	for _, f := range listFiles.Data {
		fileIDs[f.Data.Name] = f.Data.Id
		crowdin.log(fmt.Sprintf("Target language Ids: %s - %d\n", f.Data.Name, f.Data.Id))
	}

	const REC_PULLED_NB = 500
	listTrans = make(map[int][]int) 
	
	// Process each file
	for _, fileID := range fileIDs {
		crowdin.log(fmt.Sprintf("Processing file %d\n", fileID))
		// Process each language
		for _, langID := range targetLanguageIDs {
			crowdin.log(fmt.Sprintf("Processing lang %s\n", langID))
			// Get all translation approvals for this file/lang
			idx := 0
			for {
				approv, err := crowdin.ListTranslationApprovals(&ListTranslationApprovalsOptions{
				FileID:fileID,
				LanguageID:langID,
				Limit:REC_PULLED_NB,
				Offset:idx})
				if err != nil {
					fmt.Printf("ERREUR: %s\n", err)
					return listTrans, err
				}
				crowdin.log(fmt.Sprintf("Processing %d records\n", len(approv.Data)))
				
				// Store translation IDs in map
				for _, rec := range approv.Data {
					transId := rec.Data.TranslationID
					workflwId := rec.Data.WorkflowStepID
					listTrans[transId] = append(listTrans[transId], workflwId)
				}
				
				// if len(approv.Data) < REC_PULLED_NB || idx > 7250 {
				if len(approv.Data) < REC_PULLED_NB  {
					break	// Nothing left to read
				}
				idx += REC_PULLED_NB  // next page
			}	
		}
	}
	return listTrans, err
}
