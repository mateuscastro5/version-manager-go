package config

type Config struct {
	Remote string

	SourceBranch string

	DestinationBranch string

	Push bool

	RemoveBranch bool

	Tag string

	CreateRelease bool

	ReleaseTitle string

	ReleaseNotes string

	RepoType string
}

func NewConfig() *Config {
	return &Config{
		Push:         false,
		RemoveBranch: false,
	}
}
