package yamlresources

import (
	"io/ioutil"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// ClusterContentList represents openshift_cluster_content within an OpenShift-Applier inventory
type ClusterContentList struct {
	OpenShiftClusterContent []ClusterContentObject `yaml:"openshift_cluster_content"`
}

// ClusterContentObject represents a single object within ClusterContentList
type ClusterContentObject struct {
	Object  string           `yaml:"object"`
	Content []ClusterContent `yaml:"content"`
}

// ClusterContent represents the actual content of a ClusterContentObject
type ClusterContent struct {
	Name           string `yaml:"name"`
	File           string `yaml:"file,omitempty"`
	Template       string `yaml:"template,omitempty"`
	Params         string `yaml:"params,omitempty"`
	ParamsFromVars string `yaml:"params_from_vars,omitempty"`
	Action         string `yaml:"action,omitempty"`
}

func GetClusterContentsFromFile() ClusterContentList {

	var clusterContents ClusterContentList

	fileContents, err := ioutil.ReadFile("inventory/group_vars/all.yml")
	if err != nil {
		log.Fatal("Could not read the specified file.")
	}
	yaml.Unmarshal(fileContents, &clusterContents)
	if err != nil {
		log.Fatal("Unable to interpret the file as valid YAML.")
	}

	return clusterContents

}

func WriteClusterContents(contents ClusterContentList) {

	outByte := []byte{}
	outByte, _ = yaml.Marshal(contents)

	err := ioutil.WriteFile("inventory/group_vars/all.yml", outByte, 0766)
	if err != nil {
		log.Fatal("Could not write cluster contents to file.")
	}

}

func TouchParamsFile(fileName string) {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		err := ioutil.WriteFile(fileName, []byte("# Use this parameter file as shown:\n# PARAMETER=value"), 0766)
		if err != nil {
			log.Fatal("Could not touch parameter file.")
		}
	}
}
