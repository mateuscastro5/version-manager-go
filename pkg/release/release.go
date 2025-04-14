package release

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/be-tech/version-manager/internal/utils"
	"github.com/be-tech/version-manager/pkg/config"
	"github.com/joho/godotenv"
)

const (
	githubAPIBaseURL = "https://api.github.com"
	gitlabAPIBaseURL = "https://gitlab.com/api/v4"
)

type ReleaseManager struct {
	config *config.Config
	logger *utils.Logger
	gitCmd *utils.GitCommands
}

type GitHubReleaseRequest struct {
	TagName         string `json:"tag_name"`
	TargetCommitish string `json:"target_commitish"`
	Name            string `json:"name"`
	Body            string `json:"body"`
	Draft           bool   `json:"draft"`
	Prerelease      bool   `json:"prerelease"`
}

type GitLabReleaseRequest struct {
	Name        string `json:"name"`
	TagName     string `json:"tag_name"`
	Description string `json:"description"`
}

func NewReleaseManager(config *config.Config) *ReleaseManager {
	return &ReleaseManager{
		config: config,
		logger: utils.NewLogger(),
		gitCmd: utils.NewGitCommands(&defaultCommandRunner{}),
	}
}

// Implementação de CommandRunner para o ReleaseManager
type defaultCommandRunner struct{}

func (r *defaultCommandRunner) Run(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	return cmd.CombinedOutput()
}

func (r *defaultCommandRunner) Output(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	return cmd.Output()
}

func (r *ReleaseManager) CreateRelease(tagVersion string) error {
	switch r.config.RepoType {
	case "github":
		return r.createGitHubRelease(tagVersion)
	case "gitlab":
		return r.createGitLabRelease(tagVersion)
	default:
		return fmt.Errorf("tipo de repositório não suportado: %s", r.config.RepoType)
	}
}

func (r *ReleaseManager) getRepoFullName() (string, error) {
	// Executando o comando git diretamente sem tentar acessar o campo privado runner
	cmd := exec.Command("git", "remote", "get-url", r.config.Remote)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("falha ao obter URL do repositório: %v", err)
	}

	repoURL := strings.TrimSpace(string(output))
	repoURL = strings.TrimSuffix(repoURL, ".git")

	var repoFullName string

	if strings.Contains(repoURL, "github.com") {
		parts := strings.Split(repoURL, "github.com/")
		if len(parts) > 1 {
			repoFullName = parts[1]
		}
	} else if strings.Contains(repoURL, "gitlab.com") {
		parts := strings.Split(repoURL, "gitlab.com/")
		if len(parts) > 1 {
			repoFullName = parts[1]
		}
	} else {
		return "", fmt.Errorf("formato de URL do repositório não reconhecido: %s", repoURL)
	}

	if repoFullName == "" {
		return "", fmt.Errorf("não foi possível extrair o nome completo do repositório da URL: %s", repoURL)
	}

	return repoFullName, nil
}

func (r *ReleaseManager) getToken() (string, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	var envVar string
	if r.config.RepoType == "github" {
		envVar = "GITHUB_TOKEN"
	} else {
		envVar = "GITLAB_TOKEN"
	}

	token := os.Getenv(envVar)
	if token == "" {
		return "", fmt.Errorf("token de acesso não encontrado. Configure a variável de ambiente %s", envVar)
	}

	return token, nil
}

func (r *ReleaseManager) createGitHubRelease(tagVersion string) error {
	r.logger.Info("Criando release no GitHub...")

	token, err := r.getToken()
	if err != nil {
		return err
	}

	repoFullName, err := r.getRepoFullName()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/repos/%s/releases", githubAPIBaseURL, repoFullName)

	isPrerelease := strings.Contains(tagVersion, "-")

	releaseData := GitHubReleaseRequest{
		TagName:         "v" + tagVersion,
		TargetCommitish: r.config.DestinationBranch,
		Name:            r.config.ReleaseTitle,
		Body:            r.config.ReleaseNotes,
		Draft:           false,
		Prerelease:      isPrerelease,
	}

	jsonData, err := json.Marshal(releaseData)
	if err != nil {
		return fmt.Errorf("erro ao serializar dados da release: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("erro ao criar requisição: %v", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao enviar requisição: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("falha ao criar release (código %d): %s", resp.StatusCode, string(bodyBytes))
	}

	r.logger.Success("Release criada com sucesso no GitHub!")
	return nil
}

func (r *ReleaseManager) createGitLabRelease(tagVersion string) error {
	r.logger.Info("Criando release no GitLab...")

	token, err := r.getToken()
	if err != nil {
		return err
	}

	repoFullName, err := r.getRepoFullName()
	if err != nil {
		return err
	}

	repoFullName = strings.Replace(repoFullName, "/", "%2F", -1)
	url := fmt.Sprintf("%s/projects/%s/releases", gitlabAPIBaseURL, repoFullName)

	releaseData := GitLabReleaseRequest{
		Name:        r.config.ReleaseTitle,
		TagName:     "v" + tagVersion,
		Description: r.config.ReleaseNotes,
	}

	jsonData, err := json.Marshal(releaseData)
	if err != nil {
		return fmt.Errorf("erro ao serializar dados da release: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("erro ao criar requisição: %v", err)
	}

	req.Header.Set("PRIVATE-TOKEN", token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao enviar requisição: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("falha ao criar release (código %d): %s", resp.StatusCode, string(bodyBytes))
	}

	r.logger.Success("Release criada com sucesso no GitLab!")
	return nil
}
