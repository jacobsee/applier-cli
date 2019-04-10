package cmd

import (
	"testing"

	fileinterface "github.com/jacobsee/applier-cli/pkg/file_interface"
	githubapi "github.com/jacobsee/applier-cli/pkg/github_api"
)

func TestInit(t *testing.T) {

	releaseAPI := githubapi.MockReleaseAPI{}
	fileInterface := fileinterface.MockFileInterface{}
	initRun(&fileInterface, &releaseAPI)

}
