package version

import (
	"testing"
)

func TestGenerateNewTag(t *testing.T) {
	handler := NewHandler()

	testCases := []struct {
		name        string
		currentTag  string
		versionType string
		expected    string
	}{
		// Casos de teste para quando não há tag anterior
		{"No Previous Tag - Major", "", "major", "v1.0.0"},
		{"No Previous Tag - Minor", "", "minor", "v1.0.0"},
		{"No Previous Tag - Patch", "", "patch", "v1.0.0"},
		{"No Previous Tag - Premajor", "", "premajor", "v1.0.0-pre.0"},
		{"No Previous Tag - Preminor", "", "preminor", "v0.1.0-pre.0"},
		{"No Previous Tag - Prepatch", "", "prepatch", "v0.0.1-pre.0"},

		// Casos para incremento normal de versão
		{"Major Increment", "v1.0.0", "major", "v2.0.0"},
		{"Minor Increment", "v1.0.0", "minor", "v1.1.0"},
		{"Patch Increment", "v1.0.0", "patch", "v1.0.1"},

		// Casos para pré-lançamentos
		{"Premajor New", "v1.0.0", "premajor", "v2.0.0-pre.0"},
		{"Preminor New", "v1.0.0", "preminor", "v1.1.0-pre.0"},
		{"Prepatch New", "v1.0.0", "prepatch", "v1.0.1-pre.0"},
		{"Prerelease New", "v1.0.0", "prerelease", "v1.0.0-pre.0"},

		// Casos para incremento de pré-lançamentos
		{"Premajor Increment", "v1.0.0-pre.0", "premajor", "v1.0.0-pre.1"},
		{"Preminor Increment", "v1.0.0-pre.0", "preminor", "v1.0.0-pre.1"},
		{"Prepatch Increment", "v1.0.0-pre.0", "prepatch", "v1.0.0-pre.1"},
		{"Prerelease Increment", "v1.0.0-pre.0", "prerelease", "v1.0.0-pre.1"},

		// Casos sem prefixo 'v'
		{"Major No V", "1.0.0", "major", "2.0.0"},
		{"Minor No V", "1.0.0", "minor", "1.1.0"},
		{"Patch No V", "1.0.0", "patch", "1.0.1"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := handler.GenerateNewTag(tc.currentTag, tc.versionType)
			if err != nil {
				t.Fatalf("Erro inesperado: %v", err)
			}
			if result != tc.expected {
				t.Errorf("Para tag atual '%s' e tipo de versão '%s', esperava '%s' mas obteve '%s'",
					tc.currentTag, tc.versionType, tc.expected, result)
			}
		})
	}
}

func TestGenerateNewTagWithInvalidInput(t *testing.T) {
	handler := NewHandler()

	testCases := []struct {
		name        string
		currentTag  string
		versionType string
		expected    string
	}{
		// Casos com formatos inválidos
		{"Invalid Format", "v1.0", "major", "v1.0.0"},
		{"Very Invalid Format", "not-a-version", "major", "v1.0.0"},
		{"Empty Version Type", "v1.0.0", "", "v1.0.0"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := handler.GenerateNewTag(tc.currentTag, tc.versionType)
			if err != nil {
				t.Fatalf("Erro inesperado: %v", err)
			}
			// Para entradas inválidas, esperamos que o sistema ainda retorne algo razoável
			if result != tc.expected {
				t.Errorf("Para entrada inválida, esperava '%s' mas obteve '%s'",
					tc.expected, result)
			}
		})
	}
}
