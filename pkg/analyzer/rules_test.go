package analyzer

import "testing"

func TestCheckLowercase(t *testing.T) {
	tests := []struct {
		name    string
		msg     string
		wantNil bool
	}{
		{"lowercase start", "hello world", true},
		{"uppercase start", "Hello world", false},
		{"digit start", "123 items", true},
		{"empty string", "", true},
		{"single uppercase", "A", false},
		{"single lowercase", "a", true},
		{"unicode lowercase", "über cool", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := checkLowercase(tt.msg)
			if tt.wantNil && f != nil {
				t.Errorf("expected nil, got %+v", f)
			}
			if !tt.wantNil && f == nil {
				t.Error("expected finding, got nil")
			}
			if !tt.wantNil && f != nil && f.Rule != "lowercase" {
				t.Errorf("expected rule 'lowercase', got %q", f.Rule)
			}
		})
	}
}

func TestCheckLowercaseSuggestion(t *testing.T) {
	f := checkLowercase("Hello world")
	if f == nil {
		t.Fatal("expected finding")
	}
	if f.Suggestion != "hello world" {
		t.Errorf("expected suggestion 'hello world', got %q", f.Suggestion)
	}
}

func TestCheckEnglishOnly(t *testing.T) {
	tests := []struct {
		name    string
		msg     string
		wantNil bool
	}{
		{"ascii only", "hello world 123", true},
		{"with tab", "hello\tworld", true},
		{"with newline", "hello\nworld", true},
		{"cyrillic", "привет", false},
		{"chinese", "你好", false},
		{"mixed", "hello мир", false},
		{"empty", "", true},
		{"punctuation is still ASCII", "hello, world", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := checkEnglishOnly(tt.msg)
			if tt.wantNil && f != nil {
				t.Errorf("expected nil, got %+v", f)
			}
			if !tt.wantNil && f == nil {
				t.Error("expected finding, got nil")
			}
		})
	}
}

func TestCheckSpecialChars_DefaultIsStrict(t *testing.T) {
	allowed := DefaultConfig().AllowedSpecialChars // default is strict now (empty)

	tests := []struct {
		name    string
		msg     string
		wantNil bool
	}{
		{"plain text", "hello world", true},
		{"with colon", "server: started", false},
		{"with dash", "key-value", false},
		{"with period", "done.", false},
		{"with comma", "hello, world", false},
		{"with parens", "value (default)", false},
		{"with slash", "path/to/file", false},
		{"with exclamation", "error!", false},
		{"with question", "what?", false},
		{"with ellipsis", "loading...", false},
		{"with emoji", "hello 🎉", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := checkSpecialChars(tt.msg, allowed)
			if tt.wantNil && f != nil {
				t.Errorf("expected nil, got %+v", f)
			}
			if !tt.wantNil && f == nil {
				t.Error("expected finding, got nil")
			}
		})
	}
}

func TestCheckSensitiveText(t *testing.T) {
	keywords := DefaultConfig().sensitiveKeywordsSet

	tests := []struct {
		name    string
		msg     string
		wantNil bool
	}{
		{"no sensitive", "user authenticated successfully", true},
		{"contains password", "user password reset requested", false},
		{"contains token", "token validated", false},
		{"contains api key", "api_key rotated", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := checkSensitiveText(tt.msg, keywords)
			if tt.wantNil && f != nil {
				t.Errorf("expected nil, got %+v", f)
			}
			if !tt.wantNil && f == nil {
				t.Error("expected finding, got nil")
			}
		})
	}
}

func TestCheckSensitiveData(t *testing.T) {
	keywords := DefaultConfig().sensitiveKeywordsSet

	tests := []struct {
		name    string
		parts   []string
		wantNil bool
	}{
		{"no sensitive", []string{"username", "count"}, true},
		{"password", []string{"userPassword"}, false},
		{"token", []string{"authToken"}, false},
		{"api key", []string{"apiKey"}, false},
		{"secret", []string{"clientSecret"}, false},
		{"safe", []string{"name", "age", "city"}, true},
		{"empty", []string{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := checkSensitiveData(tt.parts, keywords)
			if tt.wantNil && f != nil {
				t.Errorf("expected nil, got %+v", f)
			}
			if !tt.wantNil && f == nil {
				t.Error("expected finding, got nil")
			}
		})
	}
}
