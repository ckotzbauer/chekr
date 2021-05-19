package internal

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	introducedStr  = "// +k8s:prerelease-lifecycle-gen:introduced="
	deprecatedStr  = "// +k8s:prerelease-lifecycle-gen:deprecated="
	removedStr     = "// +k8s:prerelease-lifecycle-gen:removed="
	replacementStr = "// +k8s:prerelease-lifecycle-gen:replacement="
	structRegex    = `type ([A-Za-z]*) struct \{`
)

func ParseGroups(groups []FileGroup) ([]GroupVersion, error) {
	var groupVersions []GroupVersion

	for _, g := range groups {
		group, version, err := parseRegister(g.registerFile)

		if err != nil {
			return nil, err
		}

		groupVersion, err := parseTypes(g.typesFile, group, version)

		if err != nil {
			return nil, err
		}

		if groupVersion.Resources != nil {
			groupVersions = append(groupVersions, groupVersion)
		}
	}

	return groupVersions, nil
}

func parseRegister(file string) (string, string, error) {
	content, err := ReadFile(file)

	if err != nil {
		return "", "", err
	}

	groupRegex := regexp.MustCompile(`const GroupName = "([a-z0-9\.]*)"`)
	versionRegex := regexp.MustCompile(`Version: "([a-z0-9]*)"}`)

	groupResult := groupRegex.FindStringSubmatch(content)
	versionResult := versionRegex.FindStringSubmatch(content)

	if len(versionResult) != 2 {
		return "", "", fmt.Errorf("wrong match count")
	}

	if groupResult[1] == "" {
		return "", versionResult[1], nil
	} else {
		return groupResult[1], versionResult[1], nil
	}
}

func parseTypes(file, group, version string) (GroupVersion, error) {
	content, err := ReadFile(file)

	if err != nil {
		return GroupVersion{}, err
	}

	lines := strings.Split(content, "\n")
	typeRegex := regexp.MustCompile(structRegex)
	var foundResources []Resource

	lastIntroduced := ""
	lastDeprecated := ""
	lastRemoved := ""
	lastReplacement := ""

	for _, line := range lines {
		if strings.Contains(line, introducedStr) {
			lastIntroduced = strings.TrimSpace(strings.Replace(line, introducedStr, "", 1))
		}

		if strings.Contains(line, deprecatedStr) {
			lastDeprecated = strings.TrimSpace(strings.Replace(line, deprecatedStr, "", 1))
		}

		if strings.Contains(line, removedStr) {
			lastRemoved = strings.TrimSpace(strings.Replace(line, removedStr, "", 1))
		}

		if strings.Contains(line, replacementStr) {
			lastReplacement = strings.TrimSpace(strings.Replace(line, replacementStr, "", 1))
		}

		matches := typeRegex.FindStringSubmatch(line)
		if len(matches) > 1 {
			if lastIntroduced != "" && lastDeprecated != "" {
				replacement := GroupVersionKind{Group: "", Version: "", Name: ""}

				if lastReplacement != "" {
					parts := strings.Split(lastReplacement, ",")
					replacement = GroupVersionKind{Group: parts[0], Version: parts[1], Name: parts[2]}
				}

				foundResources = append(foundResources, Resource{
					Name:        matches[1],
					Introduced:  lastIntroduced,
					Deprecated:  lastDeprecated,
					Removed:     lastRemoved,
					Replacement: replacement,
				})
			}

			lastIntroduced = ""
			lastDeprecated = ""
			lastRemoved = ""
			lastReplacement = ""
		}
	}

	return GroupVersion{Group: group, Version: version, Resources: foundResources}, nil
}
