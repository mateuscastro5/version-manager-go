package config

import (
	"testing"
)

func TestNewConfig(t *testing.T) {
	config := NewConfig()

	// Verificar se a instância foi criada corretamente
	if config == nil {
		t.Fatal("NewConfig deve retornar uma instância não nula")
	}

	// Verificar valores iniciais
	if config.Remote != "" {
		t.Errorf("Valor inicial de Remote deve ser string vazia, obteve: %s", config.Remote)
	}

	if config.SourceBranch != "" {
		t.Errorf("Valor inicial de SourceBranch deve ser string vazia, obteve: %s", config.SourceBranch)
	}

	if config.DestinationBranch != "" {
		t.Errorf("Valor inicial de DestinationBranch deve ser string vazia, obteve: %s", config.DestinationBranch)
	}

	if config.Push != false {
		t.Errorf("Valor inicial de Push deve ser false, obteve: %t", config.Push)
	}

	if config.RemoveBranch != false {
		t.Errorf("Valor inicial de RemoveBranch deve ser false, obteve: %t", config.RemoveBranch)
	}

	if config.Tag != "" {
		t.Errorf("Valor inicial de Tag deve ser string vazia, obteve: %s", config.Tag)
	}
}

func TestConfigProperties(t *testing.T) {
	// Criar e preencher uma configuração
	config := NewConfig()
	config.Remote = "origin"
	config.SourceBranch = "feature"
	config.DestinationBranch = "main"
	config.Push = true
	config.RemoveBranch = true
	config.Tag = "major"

	// Verificar se os valores foram atribuídos corretamente
	if config.Remote != "origin" {
		t.Errorf("Remote: esperava 'origin', obteve '%s'", config.Remote)
	}

	if config.SourceBranch != "feature" {
		t.Errorf("SourceBranch: esperava 'feature', obteve '%s'", config.SourceBranch)
	}

	if config.DestinationBranch != "main" {
		t.Errorf("DestinationBranch: esperava 'main', obteve '%s'", config.DestinationBranch)
	}

	if !config.Push {
		t.Errorf("Push: esperava true, obteve %t", config.Push)
	}

	if !config.RemoveBranch {
		t.Errorf("RemoveBranch: esperava true, obteve %t", config.RemoveBranch)
	}

	if config.Tag != "major" {
		t.Errorf("Tag: esperava 'major', obteve '%s'", config.Tag)
	}
}
