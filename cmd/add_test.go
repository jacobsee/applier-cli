package cmd

import (
	"testing"

	clusterinterface "github.com/jacobsee/applier-cli/pkg/cluster_interface"
	fileinterface "github.com/jacobsee/applier-cli/pkg/file_interface"
)

func TestAddFromClusterMakeTemplate(t *testing.T) {

	// Test adding a resource from a cluster, turning it into a template
	testFlags := runFlags{
		fromCluster:  true,
		fromFile:     false,
		makeTemplate: true,
		edit:         false,
	}
	clusterInterface := clusterinterface.MockClusterInterface{}
	fileInterface := fileinterface.MockFileInterface{}
	add(testFlags, []string{"test_pod"}, &clusterInterface, &fileInterface)

}
