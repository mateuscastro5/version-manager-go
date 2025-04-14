package version

import (
	"fmt"
	"strconv"
	"strings"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) GenerateNewTag(currentTag, versionType string) (string, error) {
	currentTag = strings.TrimSpace(currentTag)

	if currentTag == "" {
		if versionType == "premajor" {
			return "v1.0.0-pre.0", nil
		} else if versionType == "preminor" {
			return "v0.1.0-pre.0", nil
		} else if versionType == "prepatch" || versionType == "prerelease" {
			return "v0.0.1-pre.0", nil
		}
		return "v1.0.0", nil
	}

	tagVersion := currentTag
	hasV := false
	if strings.HasPrefix(currentTag, "v") {
		tagVersion = currentTag[1:]
		hasV = true
	}

	isPre := false
	preVersion := 0
	mainVersion := tagVersion
	if strings.Contains(tagVersion, "-pre") {
		isPre = true
		parts := strings.Split(tagVersion, "-pre.")
		if len(parts) >= 2 {
			mainVersion = parts[0]
			preStr := parts[1]
			preVersion, _ = strconv.Atoi(preStr)
		}
	}

	parts := strings.Split(mainVersion, ".")
	if len(parts) < 3 {
		return "v1.0.0", nil
	}

	major, _ := strconv.Atoi(parts[0])
	minor, _ := strconv.Atoi(parts[1])
	patch, _ := strconv.Atoi(parts[2])

	var newTag string

	switch versionType {
	case "major":
		major++
		minor = 0
		patch = 0
		newTag = fmt.Sprintf("%d.%d.%d", major, minor, patch)

	case "minor":
		minor++
		patch = 0
		newTag = fmt.Sprintf("%d.%d.%d", major, minor, patch)

	case "patch":
		patch++
		newTag = fmt.Sprintf("%d.%d.%d", major, minor, patch)

	case "premajor":
		if isPre && strings.HasPrefix(versionType, "pre") {
			preVersion++
			newTag = fmt.Sprintf("%d.%d.%d-pre.%d", major, minor, patch, preVersion)
		} else {
			major++
			minor = 0
			patch = 0
			newTag = fmt.Sprintf("%d.%d.%d-pre.0", major, minor, patch)
		}

	case "preminor":
		if isPre && strings.HasPrefix(versionType, "pre") {
			preVersion++
			newTag = fmt.Sprintf("%d.%d.%d-pre.%d", major, minor, patch, preVersion)
		} else {
			minor++
			patch = 0
			newTag = fmt.Sprintf("%d.%d.%d-pre.0", major, minor, patch)
		}

	case "prepatch":
		if isPre && strings.HasPrefix(versionType, "pre") {
			preVersion++
			newTag = fmt.Sprintf("%d.%d.%d-pre.%d", major, minor, patch, preVersion)
		} else {
			patch++
			newTag = fmt.Sprintf("%d.%d.%d-pre.0", major, minor, patch)
		}

	case "prerelease":
		if isPre {
			preVersion++
			newTag = fmt.Sprintf("%d.%d.%d-pre.%d", major, minor, patch, preVersion)
		} else {
			newTag = fmt.Sprintf("%d.%d.%d-pre.0", major, minor, patch)
		}

	default:
		newTag = fmt.Sprintf("%d.%d.%d", major, minor, patch)
	}

	if hasV {
		newTag = "v" + newTag
	}

	return newTag, nil
}
