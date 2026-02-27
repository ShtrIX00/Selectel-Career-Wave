package a

import "log/slog"

func slogExamples() {
	slog.Info("this is fine")
	slog.Info("Uppercase start")    // want `log message should start with a lowercase letter`
	slog.Warn("hello world!")       // want `log message contains disallowed character`
	slog.Error("что случилось")     // want `log message contains non-ASCII character`
	slog.Debug("loading...")        // want `log message contains ellipsis`
	slog.Info("hello 🎉")           // want `log message contains non-ASCII character` `log message contains emoji`

	// Sensitive data (no punctuation to avoid triggering rule 3 as noise)
	password := "secret"
	slog.Info("user password " + password) // want `log message may contain sensitive data` `log message may contain sensitive data`

	slog.Info("server started", "port", 8080)
	slog.Info("request completed", "status", 200)
}
