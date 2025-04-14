package config

// Config armazena todas as configurações coletadas do usuário
type Config struct {
	// Repositório remoto selecionado
	Remote string

	// Branch de origem para o merge
	SourceBranch string

	// Branch de destino para o merge
	DestinationBranch string

	// Se deve fazer push para o repositório remoto
	Push bool

	// Se deve remover a branch de origem após o merge
	RemoveBranch bool

	// Tipo de tag de versão a ser criada (major, minor, patch, etc.)
	Tag string
}

// NewConfig cria uma nova instância de Config
func NewConfig() *Config {
	return &Config{}
}
