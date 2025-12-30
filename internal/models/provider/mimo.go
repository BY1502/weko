package provider

import (
	"fmt"

	"github.com/Tencent/WeKnora/internal/types"
)

const (
	// MimoBaseURL Xiaomi MiMo API BaseURL
	MimoBaseURL = "https://api.xiaomimimo.com/v1"
)

// MimoProvider implements the Xiaomi MiMo provider interface
type MimoProvider struct{}

func init() {
	Register(&MimoProvider{})
}

// Info returns metadata for MiMo provider
func (p *MimoProvider) Info() ProviderInfo {
	return ProviderInfo{
		Name:        ProviderMimo,
		DisplayName: "Xiaomi MiMo",
		Description: "mimo-v2-flash",
		DefaultURLs: map[types.ModelType]string{
			types.ModelTypeKnowledgeQA: MimoBaseURL,
		},
		ModelTypes: []types.ModelType{
			types.ModelTypeKnowledgeQA,
		},
		RequiresAuth: true,
	}
}

// ValidateConfig validates MiMo provider config
func (p *MimoProvider) ValidateConfig(config *Config) error {
	if config.APIKey == "" {
		return fmt.Errorf("API key is required for Mimo provider")
	}
	if config.ModelName == "" {
		return fmt.Errorf("model name is required")
	}
	return nil
}
