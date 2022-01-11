package config

func InitializeConfig() GenerateConfig {
	return GenerateConfig{
		RarityDelimiter:   "#",
		DnaDelimiter:      "-",
		UniqueDnaTorrence: 10000,
	}
}

func Layers() (conf LayerConfigurations) {
	layerConfig1 := LayerConfig{
		EditionSize: 2,
		LayerOrder: append(make([]Layer, 0),
			Layer{"Background Color"},
			Layer{"Background Noise"},
			Layer{"Background Aura"},
			Layer{"Face Color"},
			Layer{"Face Texture"},
			Layer{"Face Shading"},
			Layer{"Eye Type"},
			Layer{"Eye Color"},
			Layer{"Pupil Type"},
			Layer{"Upper Eyelids"},
			Layer{"Lower Eyelids"},
			Layer{"Nose"},
			Layer{"Mouth"},
			Layer{"Head Accessory"},
			Layer{"Face Accessory"},
			Layer{"Mouth Accessory"},
		),
	}
	conf.All = append(conf.All, layerConfig1)

	return conf
}

func Metadata() MetadataConfig {
	return MetadataConfig{
		Name:        "Spooky Head Club",
		Description: "Collection of generative NFTs",
		FileUrl:     "https://test.com",
		Creator:     "Caleb T",
	}
}

func IgnoredData() Ignored {
	return Ignored{
		Traits: []string{
			"Face Shading",
		},
		Values: []string{
			"blank",
		},
	}
}
