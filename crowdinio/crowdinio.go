package crowdinio

import (
	"errors"
	"fmt"
	crowdin "github.com/fabdem/go-crowdinv2"
	"log"
	"os"
	"path/filepath"
	"time"
)

func DownloadAndExit(fn string) {
	c, err := crowdinInit()
	if err != nil {
		log.Fatal(err)
	}
	buildId, err := c.BuildTranslationAllLg(crowdin.BuildTranslationAllLgOptions{
		BuildTO:                     10 * time.Minute, // Timeout
		TranslatedOnly:              false,
		MinApprovalSteps:            0,
		FullyTranslatedFilesOnly:    false,
		ExportStringsThatPassedWkfl: false,
		FolderName:                  "",
	})
	if err != nil {
		log.Printf("WARNING: c.BuildTranslationAllLg(%q): %s", fn, err)
		log.Printf("WARNING: will use latest build, which MAY be stale")
		buildId = 0
	}

	if err = c.DownloadBuild(fn, buildId); err != nil {
		log.Fatalf("c.DownloadBuild(%q): %s", fn, err)
	}

	os.Exit(0)
}

func UploadAndExit(fn string) {
	c, err := crowdinInit()
	if err != nil {
		log.Fatal(err)
	}

	fileID, revID, err := c.Update(filepath.Base(fn), fn, "keep_translations_and_approvals")
	if err != nil {
		log.Fatalf("c.Update(%q...): %s", fn, err)
	}
	log.Printf("file ID: %d", fileID)
	log.Printf("revision ID: %d", revID)
	os.Exit(0)
}

func crowdinInit() (*crowdin.Crowdin, error) {

	config, err := load("crowdin.json")
	if err != nil {
		return nil, fmt.Errorf("crowdin.json: %w", err)
	}
	if config == nil {
		return nil, errors.New("config nil")
	}
	if config.Token == "" {
		return nil, errors.New("config missing crowdin.token")
	}
	if config.ProjectID == 0 {
		return nil, errors.New("config missing crowdin.project_id")
	}

	c, err := crowdin.New(config.Token, config.ProjectID, "", "")
	if err != nil {
		return nil, err
	}

	c.SetDebug(true, os.Stderr)
	return c, nil
}
