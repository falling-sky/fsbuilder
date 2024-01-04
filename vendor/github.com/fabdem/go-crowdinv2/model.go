package crowdin

import (
	"time"
)

// ListProjectBuilds - List Project Build API call
type ListProjectBuildsOptions struct {
	// ProjectId 			int
	// Body struct {
	BranchId int
	Limit    int
	Offset   int
	//}
}

type ResponseListProjectBuilds struct {
	Data []struct {
		Data struct {
			Id         int    `json:"id"`
			ProjectId  int    `json:"projectId"`
			Status     string `json:"status"`
			Progress   int    `json:"progress"`
			Attributes struct {
				BranchId             		int   		`json:"branchId,omitempty"`
				TargetLanguageIds    		[]string 	`json:"targetLanguageIds,omitempty"`
				SkipUntranslatedStrings 	bool		`json:"SkipUntranslatedStrings,omitempty"`
				SkipUntranslatedFiles		bool		`json:"SkipUntranslatedFiles,omitempty"`
				ExportWithMinApprovalsCount	int			`json:"ExportWithMinApprovalsCount,omitempty"`
				ExportTranslatedOnly 		bool  		`json:"exportTranslatedOnly,omitempty"`
				ExportApprovedOnly   		bool  		`json:"exportApprovedOnly,omitempty"`
			} `json:"attributes"`
		} `json:"data"`
	} `json:"data"`
	Pagination struct {
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
	} `json:"pagination"`
}

// ListProjects - List Projects API call
type ListProjectsOptions struct {
	GroupId          int
	HasManagerAccess int
	Limit            int
	Offset           int
}

type ResponseListProjects struct {
	Data []struct {
		Data struct {
			ID                   int       `json:"id"`
			UserID               int       `json:"userId"`
			SourceLanguageID     string    `json:"sourceLanguageId"`
			TargetLanguageIds    []string  `json:"targetLanguageIds"`
			LanguageAccessPolicy string    `json:"languageAccessPolicy"`
			Name                 string    `json:"name"`
			Cname                string    `json:"cname"`
			Identifier           string    `json:"identifier"`
			Description          string    `json:"description"`
			Visibility           string    `json:"visibility"`
			Logo                 string    `json:"logo"`
			PublicDownloads      bool      `json:"publicDownloads"`
			CreatedAt            time.Time `json:"createdAt"`
			UpdatedAt            time.Time `json:"updatedAt"`
			LastActivity         time.Time `json:"lastActivity"`
			TargetLanguages      []struct {
				ID                  string   `json:"id"`
				Name                string   `json:"name"`
				EditorCode          string   `json:"editorCode"`
				TwoLettersCode      string   `json:"twoLettersCode"`
				ThreeLettersCode    string   `json:"threeLettersCode"`
				Locale              string   `json:"locale"`
				AndroidCode         string   `json:"androidCode"`
				OsxCode             string   `json:"osxCode"`
				OsxLocale           string   `json:"osxLocale"`
				PluralCategoryNames []string `json:"pluralCategoryNames"`
				PluralRules         string   `json:"pluralRules"`
				PluralExamples      []string `json:"pluralExamples"`
				TextDirection       string   `json:"textDirection"`
				DialectOf           int      `json:"dialectOf"`
			} `json:"targetLanguages"`
		} `json:"data"`
	} `json:"data"`
	Pagination struct {
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
	} `json:"pagination"`
}

// UpdateFile - Update a file API call
type UpdateFileOptions struct {
	StorageId     int    `json:"storageId"`
	UpdateOption  string `json:"updateOption,omitempty"` // needs to be either: "clear_translations_and_approvals" "keep_translations" "keep_translations_and_approvals"
	ImportOptions struct {
		ContentSegmentation  bool     `json:"contentSegmentation,omitempty"`
		TranslateContent     bool     `json:"translateContent,omitempty"`
		TranslateAttributes  bool     `json:"translateAttributes,omitempty"`
		TranslatableElements []string `json:"translatableElements,omitempty"`
	} `json:"importOptions,omitempty"`
	ExportOptions struct {
		ContentSegmentation bool `json:"contentSegmentation,omitempty"`
	} `json:"exportOptions,omitempty"`
}

type ResponseUpdateFile struct {
	Data struct {
		Id            int    `json:"id"`
		ProjectId     int    `json:"projectId"`
		BranchId      int    `json:"branchId"`
		DirectoryId   int    `json:"directoryId"`
		Name          string `json:"name"`
		Title         string `json:"title"`
		Type          string `json:"type"`
		RevisionId    int    `json:"revisionId"`
		Status        string `json:"status"`
		Priority      string `json:"priority"`
		ImportOptions struct {
			FirstLineContainsHeader bool `json:"firstLineContainsHeader"`
			ImportTranslations      bool `json:"importTranslations"`
			Scheme                  struct {
				Identifier   int `json:"identifier"`
				SourcePhrase int `json:"sourcePhrase"`
				En           int `json:"en"`
				De           int `json:"de"`
			} `json:"scheme"`
		} `json:"importOptions"`
		ExportOptions struct {
			ExportPattern string `json:"exportPattern"`
		} `json:"exportOptions"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	} `json:"data"`
}

type ResponseUpdate struct {
	Data struct {
		Id            int    `json:"id"`
		ProjectId     int    `json:"projectId"`
		BranchId      int    `json:"branchId"`
		DirectoryId   int    `json:"directoryId"`
		Name          string `json:"name"`
		Title         string `json:"title"`
		Type          string `json:"type"`
		RevisionId    int    `json:"revisionId"`
		Status        string `json:"status"`
		Priority      string `json:"priority"`
		ImportOptions struct {
			FirstLineContainsHeader bool `json:"firstLineContainsHeader"`
			ImportTranslations      bool `json:"importTranslations"`
			Scheme                  struct {
				Identifier   int `json:"identifier"`
				SourcePhrase int `json:"sourcePhrase"`
				En           int `json:"en"`
				De           int `json:"de"`
			} `json:"scheme"`
		} `json:"importOptions"`
		ExportOptions struct {
			ExportPattern string `json:"exportPattern"`
		} `json:"exportOptions"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	} `json:"data"`
}

// GetProjectBuildProgressOptions are options for Check Project Build Status api call
type CheckProjectBuildStatusOptions struct {
	// Project Identifier.
	// ProjectId int
	BuildId int
}

// {"data":{"id":47,"projectId":17,"status":"inProgress","progress":11,"attributes":{"branchId":null,"targetLanguageIds":[],"exportTranslatedOnly":false,"exportApprovedOnly":false}}}
/* type ResponseCheckProjectBuildStatus struct {
	Data struct {
		Id         int    `json:"id"`
		ProjectId  int    `json:"projectId"`
		Status     string `json:"status"`
		Progress   int    `json:"progress"`
		Attributes struct {
			BranchId             int   `json:"branchId,omitempty"`
			TargetLanguageIds    []int `json:"targetLanguageIds,omitempty"`
			ExportTranslatedOnly bool  `json:"exportTranslatedOnly"`
			ExportApprovedOnly   bool  `json:"exportApprovedOnly"`
		} `json:"attributes"`
	} `json:"data"`
}
*/

type ResponseCheckProjectBuildStatus struct {
	Data struct {
		Id         int    `json:"id"`
		ProjectId  int    `json:"projectId"`
		Status     string `json:"status"`
		Progress   int    `json:"progress"`
		Attributes struct {
			BranchId                    int      `json:"branchId,omitempty"`
			TargetLanguageIds           []string `json:"targetLanguageIds,omitempty"`
			SkipUntranslatedStrings     bool     `json:"skipUntranslatedStrings"`
			SkipUntranslatedFiles       bool     `json:"skipUntranslatedFiles"`
			ExportApprovedOnly          bool     `json:"exportApprovedOnly,omitempty"`          // crowdin.com
			ExportWithMinApprovalsCount int      `json:"exportWithMinApprovalsCount,omitempty"` // Enterprise
		} `json:"attributes"`
	} `json:"data"`
}

// DownloadProjectTranslationsOptions are options for  DownloadProjectTranslations api call
type DownloadProjectTranslationsOptions struct {
	// Project Identifier.
	// ProjectId int
	// Build Identifier.
	BuildId int
}

type ResponseDownloadProjectTranslations struct {
	Data struct {
		Url      string `json:"url"`
		ExpireIn string `json:"expireIn"`
	} `json: "data"`
}

// GetFileProgress - options for Language Progress API call
type GetFileProgressOptions struct {
	FileId int
	Limit  int
	Offset int
}

// GetFileProgress api call
type ResponseGetFileProgress struct {
	Data []struct {
		Data struct {
			Words struct {
				Total      int `json:"total"`
				Translated int `json:"translated"`
				Approved   int `json:"approved"`
			} `json:"words"`
			Phrases struct {
				Total      int `json:"total"`
				Translated int `json:"translated"`
				Approved   int `json:"approved"`
			} `json:"phrases"`
			TranslationProgress int    `json:"translationProgress"`
			ApprovalProgress    int    `json:"approvalProgress"`
			LanguageId          string `json:"languageId"`
			ETag                string `json:"eTag"`
		} `json:"data"`
	} `json:"data"`
	Pagination struct {
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
	} `json:"pagination"`
}

// BuildProjectTranslationOptions are options for BuildProjectTranslation api call
type BuildProjectTranslationOptions struct {
	BranchId int `json:"branchId,omitempty"` // Branch Identifier. - optional
	// Specify target languages for build.
	// Leave this field empty to build all target languages
	Languages                   	[]string `json:"targetLanguageIds,omitempty"`
	SkipUntranslatedStrings     	bool     `json:"skipUntranslatedStrings,omitempty"`
	SkipUntranslatedFiles       	bool     `json:"skipUntranslatedFiles,omitempty"`
	ExportApprovedOnly          	bool     `json:"exportApprovedOnly,omitempty"`          // crowdin.com specific
	ExportWithMinApprovalsCount 	int      `json:"exportWithMinApprovalsCount,omitempty"` // Enterprise specific
	ExportStringsThatPassedWorkflow bool     `json:"exportStringsThatPassedWorkflow,omitempty"`
}

type ResponseBuildProjectTranslation struct {
	Data struct {
		Id         int    `json:"id"`
		ProjectID  int    `json:"projectId"`
		Status     string `json:"status"`
		Progress   int    `json:"progress"`
		Attributes struct {
			BranchID                    	int      `json:"branchId"`
			TargetLanguageIDs           	[]string `json:"targetLanguageIds"`
			SkipUntranslatedStrings     	bool     `json:"skipUntranslatedStrings"`
			SkipUntranslatedFiles       	bool     `json:"skipUntranslatedFiles"`
			ExportApprovedOnly          	bool     `json:"exportApprovedOnly"`
			ExportWithMinApprovalsCount 	int      `json:"exportWithMinApprovalsCount"`
			ExportStringsThatPassedWorkflow bool     `json:"exportStringsThatPassedWorkflow"`
		} `json:"attributes"`
	} `json:"data"`
}


// BuildDirectoryTranslationOptions are options for BuildDirectoryTranslation api call
type BuildDirectoryTranslationOptions struct {
	TargetLanguageIds           	[]string `json:"targetLanguageIds,omitempty"`
	SkipUntranslatedStrings     	bool     `json:"skipUntranslatedStrings,omitempty"`
	SkipUntranslatedFiles       	bool     `json:"skipUntranslatedFiles,omitempty"`
	ExportApprovedOnly      		bool     `json:"exportApprovedOnly,omitempty"`			// crowdin.com specific
	ExportWithMinApprovalsCount 	int      `json:"exportWithMinApprovalsCount,omitempty"` // Enterprise specific
	PreserveFolderHierarchy     	bool     `json:"preserveFolderHierarchy,omitempty"`
	ExportStringsThatPassedWorkflow bool     `json:"exportStringsThatPassedWorkflow,omitempty"`
}

type ResponseBuildDirectoryTranslation struct {
	Data struct {
		ID         int       `json:"id"`
		ProjectID  int       `json:"projectId"`
		Status     string    `json:"status"`
		Progress   int       `json:"progress"`
		CreatedAt  time.Time `json:"createdAt"`
		UpdatedAt  time.Time `json:"updatedAt"`
		FinishedAt time.Time `json:"finishedAt"`
	} `json:"data"`
}


// BuildDirectoryFileOptions are options for BuildFileTranslation api call
type BuildFileTranslationOptions struct {
	TargetLanguageID            	string `json:"targetLanguageId"`
	ExportAsXliff               	bool   `json:"exportAsXliff"`
	SkipUntranslatedStrings     	bool   `json:"skipUntranslatedStrings"`
	SkipUntranslatedFiles       	bool   `json:"skipUntranslatedFiles"`
	ExportApprovedOnly          	bool   `json:"exportApprovedOnly,omitempty"`			// crowdin.com specific
	ExportWithMinApprovalsCount 	int    `json:"exportWithMinApprovalsCount,omitempty"`	// Enterprise specific
	ExportStringsThatPassedWorkflow bool   `json:"exportStringsThatPassedWorkflow,omitempty"`
}

type ResponseBuildFileTranslation struct {
	Data struct {
		URL      string    `json:"url"`
		ExpireIn time.Time `json:"expireIn"`
		Etag     string    `json:"etag"`
	} `json:"data"`
}

// ListStoragesOptions are options for ListStorages api call
type ListStoragesOptions struct {
	Limit  int // Maximum number of items to retrieve (25 default, max 500) - optional
	Offset int // Offset in collection - optional
}

// ResponseListStorages are response for ListStorages api call
type ResponseListStorages struct {
	Data []struct {
		Data struct {
			Id       int    `json:"id"`
			FileName string `json:"fileName"`
		} `json:"data"`
	} `json:"data"`
	Pagination []struct {
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
	} `json:"pagination"`
}

// AddStorageOptions are options for AddStorage api call
type AddStorageOptions struct {
	FileName string `json:"filename"` // Filename and path of hte file to upload to storage
}

// ResponseAddStorage are response for AddStorage api call
type ResponseAddStorage struct {
	Data struct {
		Id       int    `json:"id"`
		FileName string `json:"fileName"`
	} `json:"data"`
}

// GetStorageOptions are options for GetStorage api call
type GetStorageOptions struct {
	StorageId int
}

// ResponseGetStorage are response for ListStorages api call
type ResponseGetStorage struct {
	Data struct {
		Id       int    `json:"id"`
		FileName string `json:"fileName"`
	} `json:"data"`
}

// DelStorageOptions are options for DelStorage api call
type DeleteStorageOptions struct {
	StorageId int `json:"storageid"`
}

// ListDirectoriesOptions are options for ListDirectories api call
type ListDirectoriesOptions struct {
	BranchId    int
	DirectoryId int
	Recursion   int
	Limit       int // Maximum number of items to retrieve (25 default, max 500) - optional
	Offset      int // Offset in collection - optional
}

// ResponseListDirectories are response for ListDirectories api call
type ResponseListDirectories struct {
	Data []struct {
		Data struct {
			Id            int       `json:"id"`
			ProjectId     int       `json:"projectId"`
			BranchId      int       `json:"branchId"`
			DirectoryId   int       `json:"directoryId"` // Actually parentId
			Name          string    `json:"name"`
			Title         string    `json:"title"`
			ExportPattern string    `json:"exportPattern"`
			Priority      string    `json:"priority"`
			CreatedAt     time.Time `json:"createdAt"`
			UpdatedAt     time.Time `json:"updatedAt"`
		} `json:"data"`
	} `json:"data"`
	Pagination struct {
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
	} `json:"pagination"`
}

// ListFilesOptions are options for ListFiles api call
type ListFilesOptions struct {
	BranchId    int `json:"branchId,omitempty"`
	DirectoryId int `json:"directoryId,omitempty"`
	Recursion   int `json:"recursion,omitempty"`
	Limit       int `json:"limit,omitempty"`  // Maximum number of items to retrieve (25 default, max 500) - optional
	Offset      int `json:"offset,omitempty"` // Offset in collection - optional
}

// ResponseListFiles are response for ListFiles api call
type ResponseListFiles struct {
	Data []struct {
		Data struct {
			Id          int    `json:"id"`
			ProjectId   int    `json:"projectId"`
			BranchId    int    `json:"branchId"`
			DirectoryId int    `json:"directoryId"`
			Name        string `json:"name"`
			Title       string `json:"title"`
			Type        string `json:"type"`
			RevisionId  int    `json:"revisionId"`
			Status      string `json:"status"`
			Priority    string `json:"priority"`
			Attributes  struct {
				MimeType string `json:"mimeType"`
				FileSize int    `json:"fileSize"`
			} `json:"attributes"`
			ImportOptions struct {
				FirstLineContainsHeader bool `json:"firstLineContainsHeader"`
				ImportTranslations      bool `json:"importTranslations"`
				Scheme                  struct {
					Identifier   int `json:"identifier"`
					SourcePhrase int `json:"sourcePhrase"`
					En           int `json:"en"`
					De           int `json:"de"`
				} `json:"scheme"`
			} `json:"importOptions"`
			ExportOptions struct {
				ExportPattern string `json:"exportPattern"`
			} `json:"exportOptions"`
			ExportPattern string    `json:"exportPattern"`
			CreatedAt     time.Time `json:"createdAt"`
			UpdatedAt     time.Time `json:"updatedAt"`
			Revision      int       `json:"revision"`
		} `json:"data"`
	} `json:"data"`
	Pagination struct {
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
	} `json:"pagination"`
}

type responseGeneral struct {
	Success bool `json:"success"`
}

// ListFilesOptions are options for ListFileRevisions api call
type ListFileRevisionsOptions struct {
	Limit  int // Maximum number of items to retrieve (25 default, max 500) - optional
	Offset int // Offset in collection - optional
}

// ResponseListFiles are response for ListFiles api call
type ResponseListFileRevisions struct {
	Data []struct {
		Data struct {
			Id                int `json:"id"`
			ProjectId         int `json:"projectId"`
			FileId            int `json:"fileId"`
			RestoreToRevision int `json:"restoreToRevision"`
			Info              struct {
				Added struct {
					Strings int `json:"strings"`
					Words   int `json:"words"`
				} `json:"added"`
				Deleted struct {
					Strings int `json:"strings"`
					Words   int `json:"words"`
				} `json:"deleted"`
				Updated struct {
					Strings int `json:"strings"`
					Words   int `json:"words"`
				} `json:"updated"`
			} `json:"info"`
			Date time.Time `json:"date"`
		} `json:"data"`
	} `json:"data"`
	Pagination struct {
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
	} `json:"pagination"`
}

// ListStringsOptions are options for ListStrings api call
type ListStringsOptions struct {
	FileId                  int
	DenormalizePlaceholders int
	LabelIds                string
	Filter                  string
	Scope                   string
	Limit                   int // Maximum number of items to retrieve (25 default, max 500) - optional
	Offset                  int // Offset in collection - optional
}

// ResponseGetFileRevision are response for GetFileRevision api call
type ResponseGetFileRevision struct {
	Data struct {
		Id                int `json:"id"`
		ProjectId         int `json:"projectId"`
		FileId            int `json:"fileId"`
		RestoreToRevision int `json:"restoreToRevision"`
		Info              struct {
			Added struct {
				Strings int `json:"strings"`
				Words   int `json:"words"`
			} `json:"added"`
			Deleted struct {
				Strings int `json:"strings"`
				Words   int `json:"words"`
			} `json:"deleted"`
			Updated struct {
				Strings int `json:"strings"`
				Words   int `json:"words"`
			} `json:"updated"`
		} `json:"info"`
		Date time.Time `json:"date"`
	} `json:"data"`
}

// ResponseListStrings are response for ListStrings api call
type ResponseListStrings struct {
	Data []struct {
		Data struct {
			ID         int       `json:"id"`
			ProjectID  int       `json:"projectId"`
			FileID     int       `json:"fileId"`
			Identifier string    `json:"identifier"`
			Text       string    `json:"text"`
			Type       string    `json:"type"`
			Context    string    `json:"context"`
			MaxLength  int       `json:"maxLength"`
			IsHidden   bool      `json:"isHidden"`
			Revision   int       `json:"revision"`
			HasPlurals bool      `json:"hasPlurals"`
			IsIcu      bool      `json:"isIcu"`
			CreatedAt  time.Time `json:"createdAt"`
			UpdatedAt  time.Time `json:"updatedAt"`
		} `json:"data"`
	} `json:"data"`
	Pagination struct {
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
	} `json:"pagination"`
}

// EditString - Edit a source string API call
type EditStringsOptions []struct {
	Value interface{} `json:"value"` // type depends on path value
	Op    string      `json:"op"`
	Path  string      `json:"path"`
}

type ResponseEditStrings struct {
	Data struct {
		ID         int       `json:"id"`
		ProjectID  int       `json:"projectId"`
		FileID     int       `json:"fileId"`
		Identifier string    `json:"identifier"`
		Text       string    `json:"text"`
		Type       string    `json:"type"`
		Context    string    `json:"context"`
		MaxLength  int       `json:"maxLength"`
		IsHidden   bool      `json:"isHidden"`
		Revision   int       `json:"revision"`
		HasPlurals bool      `json:"hasPlurals"`
		IsIcu      bool      `json:"isIcu"`
		CreatedAt  time.Time `json:"createdAt"`
		UpdatedAt  time.Time `json:"updatedAt"`
	} `json:"data"`
}

// Upload translations API call
type UploadTranslationsOptions struct {
	StorageID           int  `json:"storageId"`
	FileID              int  `json:"fileId"`
	ImportEqSuggestions bool `json:"importEqSuggestions,omitempty"`
	AutoApproveImported bool `json:"autoApproveImported,omitempty"`
	TranslateHidden     bool `json:"translateHidden,omitempty"`
}

type ResponseUploadTranslations struct {
	Data struct {
		ProjectID  int    `json:"projectId"`
		StorageID  int    `json:"storageId"`
		LanguageID string `json:"languageId"`
		FileID     int    `json:"fileId"`
	} `json:"data"`
}

//
// List Translation Approvals API call
type ListTranslationApprovalsOptions struct {
	TranslationID       int		`json:"translationId,omitempty"`
	FileID              int		`json:"fileId,omitempty"`
	StringID			int		`json:"stringId,omitempty"`
	LanguageID			string	`json:"languageId,omitempty"`
	Limit               int		`json:"limit,omitempty"`  // Maximum number of items to retrieve (25 default, max 500) - optional
	Offset              int		`json:"offset,omitempty"` // Offset in collection - optional
}

type ResponseListTranslationApprovals struct {
	Data []struct {
		Data struct {
			ID   int `json:"id"`
			User struct {
				ID        int    `json:"id"`
				Username  string `json:"username"`
				FullName  string `json:"fullName"`
				AvatarURL string `json:"avatarUrl"`
			} `json:"user"`
			TranslationID  int       `json:"translationId"`
			StringID       int       `json:"stringId"`
			LanguageID     string    `json:"languageId"`
			WorkflowStepID int       `json:"workflowStepId"`
			CreatedAt      time.Time `json:"createdAt"`
		} `json:"data"`
	} `json:"data"`
	Pagination struct {
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
	} `json:"pagination"`
}

//
// List workflow steps API call
type ListWorkflowsStepsOptions struct {
	Limit               int		`json:"limit,omitempty"`  // Maximum number of items to retrieve (25 default, max 500) - optional
	Offset              int		`json:"offset,omitempty"` // Offset in collection - optional
}


type ResponseListWorkflowsSteps struct {
	Data []struct {
		Data struct {
			ID        int        `json:"id"`
			Title     string     `json:"title"`
			Type      string     `json:"type"`
			Languages []string	 `json:"languages"`
			Config    interface{}`json:"config, omitempty"`
//			Config    struct {
//				MinRelevant			string	`json:"minRelevant,omitempty"`
//				AutoSubstitution	int		`json:"autoSubstitution,omitempty"`
//				Assignees interface{} `json:"assignees,omitempty"`
//			} `json:"config"`
		} `json:"data"`
	} `json:"data"`
	Pagination struct {
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
	} `json:"pagination"`
}

//
// Get Project details API call
type ResponseGetProject struct {
	Data struct {
		ID                int       `json:"id"`
		GroupID           int       `json:"groupId"`
		UserID            int       `json:"userId"`
		SourceLanguageID  string    `json:"sourceLanguageId"`
		TargetLanguageIds []string  `json:"targetLanguageIds"`
		Name              string    `json:"name"`
		Identifier        string    `json:"identifier"`
		Description       string    `json:"description"`
		Logo              string    `json:"logo"`
		Background        string    `json:"background"`
		IsExternal        bool      `json:"isExternal"`
		ExternalType      string    `json:"externalType"`
		WorkflowID        int       `json:"workflowId"`
		HasCrowdsourcing  bool      `json:"hasCrowdsourcing"`
		PublicDownloads   bool      `json:"publicDownloads"`
		CreatedAt         time.Time `json:"createdAt"`
		UpdatedAt         time.Time `json:"updatedAt"`
		LastActivity      time.Time `json:"lastActivity"`
		TargetLanguages   []struct {
			ID                  string   `json:"id"`
			Name                string   `json:"name"`
			EditorCode          string   `json:"editorCode"`
			TwoLettersCode      string   `json:"twoLettersCode"`
			ThreeLettersCode    string   `json:"threeLettersCode"`
			Locale              string   `json:"locale"`
			AndroidCode         string   `json:"androidCode"`
			OsxCode             string   `json:"osxCode"`
			OsxLocale           string   `json:"osxLocale"`
			PluralCategoryNames []string `json:"pluralCategoryNames"`
			PluralRules         string   `json:"pluralRules"`
			PluralExamples      []string `json:"pluralExamples"`
			TextDirection       string   `json:"textDirection"`
			DialectOf           string   `json:"dialectOf"`
		} `json:"targetLanguages"`
	} `json:"data"`
}

//
// List workflow steps API call
type GetTranslationOptions struct {
	TranslationID       	int	`json:"translationID,omitempty"`
	DenormalizePlaceholders int `json:"denormalizePlaceholders, omitempty"`
}

type ResponseGetTranslation struct {
	Data struct {
		ID                 int    `json:"id"`
		Text               string `json:"text"`
		PluralCategoryName string `json:"pluralCategoryName"`
		User               struct {
			ID        int    `json:"id"`
			Username  string `json:"username"`
			FullName  string `json:"fullName"`
			AvatarURL string `json:"avatarUrl"`
		} `json:"user"`
		Rating    int       `json:"rating"`
		CreatedAt time.Time `json:"createdAt"`
	} `json:"data"`
}


//
// Add Approval API call
type AddApprovalOptions struct {
	TranslationID       	int	`json:"translationId,omitempty"`
}

type ResponseAddApproval struct {
	Data struct {
		ID   int `json:"id"`
		User struct {
			ID        int    `json:"id"`
			Username  string `json:"username"`
			FullName  string `json:"fullName"`
			AvatarURL string `json:"avatarUrl"`
		} `json:"user"`
		TranslationID  int       `json:"translationId"`
		StringID       int       `json:"stringId"`
		LanguageID     string    `json:"languageId"`
		WorkflowStepID int       `json:"workflowStepId"`
		CreatedAt      time.Time `json:"createdAt"`
	} `json:"data"`
}
