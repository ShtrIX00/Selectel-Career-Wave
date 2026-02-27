package analyzer

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the logcheck linter configuration.
type Config struct {
	SensitiveKeywords    []string `yaml:"sensitive_keywords"`
	AllowedSpecialChars  string   `yaml:"allowed_special_chars"`
	DisabledRules        []string `yaml:"disabled_rules"`
	disabledRulesSet     map[string]bool
	sensitiveKeywordsSet map[string]bool
}

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	c := &Config{
		SensitiveKeywords: []string{
			"password", "passwd", "secret", "token",
			"api_key", "apikey", "api-key",
			"credential", "private_key", "access_token", "auth",
		},
		AllowedSpecialChars: "",
	}
	c.init()
	return c
}

func (c *Config) init() {
	c.disabledRulesSet = make(map[string]bool, len(c.DisabledRules))
	for _, r := range c.DisabledRules {
		c.disabledRulesSet[r] = true
	}
	c.sensitiveKeywordsSet = make(map[string]bool, len(c.SensitiveKeywords))
	for _, k := range c.SensitiveKeywords {
		c.sensitiveKeywordsSet[k] = true
	}
}

// IsRuleDisabled returns true if the given rule name is disabled.
func (c *Config) IsRuleDisabled(rule string) bool {
	return c.disabledRulesSet[rule]
}

// LoadConfig loads configuration from a YAML file at the given path.
// If the file does not exist, it returns the default config.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return nil, err
	}
	c := DefaultConfig()
	if err := yaml.Unmarshal(data, c); err != nil {
		return nil, err
	}
	c.init()
	return c, nil
}

// LoadConfigFromWorkDir tries to load .logcheck.yml from the working directory.
func LoadConfigFromWorkDir() *Config {
	c, err := LoadConfig(".logcheck.yml")
	if err != nil {
		return DefaultConfig()
	}
	return c
}
