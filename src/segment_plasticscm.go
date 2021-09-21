package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
)

type plasticscm struct {
	props *properties
	env   environmentInfo
}

const (
	// Binary to be used
	CmProperty        Property = "cmPath"
	ScmFolderProperty Property = "scmFolder"
)

func (scm *plasticscm) enabled() bool {
	binary := scm.props.getString(CmProperty, "cm")
	if !scm.env.hasCommand(binary) {
		return false
	}

	_, err := scm.getHomeDir()

	return err == nil
}

func (scm *plasticscm) string() string {
	homeDir, err := scm.getHomeDir()
	if err != nil {
		return "ERR"
	}

	selectorPath := fmt.Sprintf("%s/plastic.selector", homeDir.path)
	content, err := ioutil.ReadFile(selectorPath)
	if err != nil {
		return "ERR"
	}

	str1 := string(content)

	branchName := scm.getBranchName(str1)

	return branchName
}

func (scm *plasticscm) init(props *properties, env environmentInfo) {
	scm.props = props
	scm.env = env
}

func (scm *plasticscm) getHomeDir() (*fileInfo, error) {
	scmFolder := scm.props.getString(ScmFolderProperty, ".plastic")
	homeDir, err := scm.env.hasParentFilePath(scmFolder)

	return homeDir, err
}

func (scm *plasticscm) getBranchName(content string) string {
	re := regexp.MustCompile(`.*(smartbranch|br|branch) "(?P<branchname>.*)"`)
	branchName := re.SubexpNames()

	result := re.FindAllStringSubmatch(content, -1)

	m := map[string]string{}
	for i, n := range result[0] {
		m[branchName[i]] = n
	}

	return m["branchname"]
}
