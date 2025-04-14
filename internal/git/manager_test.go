package git

import (
	"fmt"
	"testing"

	"github.com/be-tech/version-manager/pkg/config"
)

// MockCommandRunner é uma implementação simulada de CommandRunner para testes
type MockCommandRunner struct {
	// Mapa de comando para resultado simulado
	commandResults map[string]struct {
		output []byte
		err    error
	}
}

// NewMockCommandRunner cria uma nova instância de MockCommandRunner
func NewMockCommandRunner() *MockCommandRunner {
	return &MockCommandRunner{
		commandResults: make(map[string]struct {
			output []byte
			err    error
		}),
	}
}

// AddMockResult adiciona um resultado simulado para um comando específico
func (m *MockCommandRunner) AddMockResult(cmd string, output []byte, err error) {
	m.commandResults[cmd] = struct {
		output []byte
		err    error
	}{output, err}
}

// Run simula a execução de um comando e retorna um resultado pré-configurado
func (m *MockCommandRunner) Run(name string, args ...string) ([]byte, error) {
	// Constrói uma string representando o comando para pesquisa no mapa
	cmd := name
	for _, arg := range args {
		cmd += " " + arg
	}

	// Verifica se o comando está no mapa de resultados simulados
	if result, ok := m.commandResults[cmd]; ok {
		return result.output, result.err
	}

	// Se o comando não foi configurado, retorna erro padrão
	return []byte("Mock command not configured: " + cmd), fmt.Errorf("mock command not configured: %s", cmd)
}

// Output simula a saída stdout de um comando
func (m *MockCommandRunner) Output(name string, args ...string) ([]byte, error) {
	// Implementação similar a Run
	return m.Run(name, args...)
}

// TestNewManager verifica se o NewManager está criando uma instância corretamente
func TestNewManager(t *testing.T) {
	cfg := &config.Config{
		Remote:            "origin",
		SourceBranch:      "feature",
		DestinationBranch: "main",
		Push:              true,
		RemoveBranch:      false,
		Tag:               "minor",
	}

	manager := NewManager(cfg)

	if manager == nil {
		t.Fatal("NewManager deve retornar uma instância não nula")
	}

	if manager.config != cfg {
		t.Error("A configuração armazenada não corresponde à configuração fornecida")
	}
}

// TestCheckoutDestinationBranch testa o método checkoutDestinationBranch usando mock
func TestCheckoutDestinationBranch(t *testing.T) {
	// Configuração para teste
	cfg := &config.Config{
		DestinationBranch: "main",
	}

	// Criar o mock runner
	mockRunner := NewMockCommandRunner()
	mockRunner.AddMockResult("git checkout main", []byte("Switched to branch 'main'"), nil)

	// Criar manager com o mock
	manager := NewManagerWithRunner(cfg, mockRunner)

	// Executar o método a testar
	err := manager.checkoutDestinationBranch()

	// Verificar resultado
	if err != nil {
		t.Errorf("checkoutDestinationBranch falhou com erro: %v", err)
	}
}

// TestCheckoutDestinationBranchError testa cenário de erro no checkout
func TestCheckoutDestinationBranchError(t *testing.T) {
	// Configuração para teste
	cfg := &config.Config{
		DestinationBranch: "main",
	}

	// Criar o mock runner com erro
	mockRunner := NewMockCommandRunner()
	mockRunner.AddMockResult(
		"git checkout main",
		[]byte("error: pathspec 'main' did not match any file(s) known to git"),
		fmt.Errorf("exit status 1"),
	)

	// Criar manager com o mock
	manager := NewManagerWithRunner(cfg, mockRunner)

	// Executar o método a testar
	err := manager.checkoutDestinationBranch()

	// Verificar resultado - espera-se um erro
	if err == nil {
		t.Error("checkoutDestinationBranch deveria ter falhado com um erro, mas retornou nil")
	}
}

// TestMergeBranches testa o método mergeBranches usando mock
func TestMergeBranches(t *testing.T) {
	// Configuração para teste
	cfg := &config.Config{
		SourceBranch:      "feature",
		DestinationBranch: "main",
	}

	// Criar o mock runner
	mockRunner := NewMockCommandRunner()
	mockRunner.AddMockResult("git merge feature", []byte("Updating abc123..def456\nFast-forward"), nil)

	// Criar manager com o mock
	manager := NewManagerWithRunner(cfg, mockRunner)

	// Executar o método a testar
	err := manager.mergeBranches(false)

	// Verificar resultado
	if err != nil {
		t.Errorf("mergeBranches falhou com erro: %v", err)
	}
}

// TestCreateVersionTag_NoTag verifica se createVersionTag funciona quando não há tag
func TestCreateVersionTag_NoTag(t *testing.T) {
	// Configuração para o teste com Tag vazia
	cfg := &config.Config{
		Tag: "",
	}

	// Criar o mock runner
	mockRunner := NewMockCommandRunner()

	// Criar manager com o mock
	manager := NewManagerWithRunner(cfg, mockRunner)

	// O método deve retornar nil quando não há tag
	err := manager.createVersionTag()
	if err != nil {
		t.Errorf("createVersionTag com tag vazia deve retornar nil, obteve erro: %v", err)
	}
}

// TestCreateVersionTag testa a criação de uma nova tag de versão
func TestCreateVersionTag(t *testing.T) {
	// Configuração para o teste
	cfg := &config.Config{
		Tag: "minor",
	}

	// Criar o mock runner
	mockRunner := NewMockCommandRunner()
	mockRunner.AddMockResult("git describe --tags --abbrev=0", []byte("v1.0.0"), nil)
	mockRunner.AddMockResult(
		"git tag -a v1.1.0 -m Version v1.1.0",
		[]byte(""),
		nil,
	)

	// Criar manager com o mock
	manager := NewManagerWithRunner(cfg, mockRunner)

	// Executar o método a testar
	err := manager.createVersionTag()

	// Verificar resultado
	if err != nil {
		t.Errorf("createVersionTag falhou com erro: %v", err)
	}
}
