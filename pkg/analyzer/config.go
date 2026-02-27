package analyzer

import (
	"os"
	"strings"

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

// buildSensitiveKeywordsSet normalizes and expands keywords for matching.
// Examples:
// - "api_key" -> "api_key", "api key", "api", "key", "apikey"
// - "access-token" -> "access-token", "access token", "access", "token", "accesstoken"
func buildSensitiveKeywordsSet(raw []string) map[string]bool {
	out := make(map[string]bool)

	add := func(s string) {
		s = strings.ToLower(strings.TrimSpace(s))
		if s != "" {
			out[s] = true
		}
	}

	repl := strings.NewReplacer("_", " ", "-", " ")

	for _, kw := range raw {
		kw = strings.TrimSpace(kw)
		if kw == "" {
			continue
		}

		// Original (lowercased later in add)
		add(kw)

		// Normalized separators to spaces
		kw2 := repl.Replace(kw)
		add(kw2)

		// Tokens (api_key -> api, key)
		for _, t := range strings.Fields(kw2) {
			add(t)
		}

		// Joined form (api key -> apikey)
		add(strings.ReplaceAll(kw2, " ", ""))
	}

	return out
}

func (c *Config) init() {
	// disabled rules: normalize to lower-case + trim
	c.disabledRulesSet = make(map[string]bool, len(c.DisabledRules))
	for _, r := range c.DisabledRules {
		r = strings.ToLower(strings.TrimSpace(r))
		if r == "" {
			continue
		}
		c.disabledRulesSet[r] = true
	}

	// sensitive keywords: normalize + expand
	c.sensitiveKeywordsSet = buildSensitiveKeywordsSet(c.SensitiveKeywords)
}

// IsRuleDisabled returns true if the given rule name is disabled.
func (c *Config) IsRuleDisabled(rule string) bool {
	rule = strings.ToLower(strings.TrimSpace(rule))
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
