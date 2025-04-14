package version

import (
	"fmt"
	"strconv"
	"strings"
)

// Handler gerencia operações de versionamento semântico
type Handler struct{}

// NewHandler cria uma nova instância de Handler
func NewHandler() *Handler {
	return &Handler{}
}

// GenerateNewTag gera uma nova tag baseada na tag anterior e no tipo de versão
func (h *Handler) GenerateNewTag(currentTag, versionType string) (string, error) {
	// Remove nova linha e quaisquer caracteres de controle
	currentTag = strings.TrimSpace(currentTag)

	// Se não há tag anterior, começamos com v1.0.0
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

	// Remove o 'v' inicial se houver
	tagVersion := currentTag
	hasV := false
	if strings.HasPrefix(currentTag, "v") {
		tagVersion = currentTag[1:]
		hasV = true
	}

	// Identificar se é uma versão pre-release
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

	// Divida a versão em partes
	parts := strings.Split(mainVersion, ".")
	if len(parts) < 3 {
		// Se a versão não estiver no formato esperado, use v1.0.0
		return "v1.0.0", nil
	}

	// Analisa as partes numéricas
	major, _ := strconv.Atoi(parts[0])
	minor, _ := strconv.Atoi(parts[1])
	patch, _ := strconv.Atoi(parts[2])

	// Gera a nova versão com base no tipo solicitado
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
			// Se já é uma pre-release do mesmo tipo, incrementa apenas o preVersion
			preVersion++
			newTag = fmt.Sprintf("%d.%d.%d-pre.%d", major, minor, patch, preVersion)
		} else {
			// Se não é uma pre-release ou é de outro tipo, incrementa a versão e adiciona -pre.0
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

	// Adicionar o 'v' novamente se existia
	if hasV {
		newTag = "v" + newTag
	}

	return newTag, nil
}
