package config

type GenerateConfig struct {
	RarityDelimiter   string
	DnaDelimiter      string
	UniqueDnaTorrence int
}

type LayerConfigurations struct {
	All []LayerConfig
}

type LayerConfig struct {
	EditionSize int
	LayerOrder  []Layer
}

type Layer struct {
	Name string
}

type MetadataConfig struct {
	Name        string
	Description string
	FileUrl     string
	Creator     string
}

type Ignored struct {
	Traits []string
	Values []string
}
