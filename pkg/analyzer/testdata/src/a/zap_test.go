package a

import "go.uber.org/zap"

func zapExamples() {
	logger, _ := zap.NewProduction()
	sugar := logger.Sugar()

	logger.Info("this is fine")
	logger.Info("Uppercase start")    // want `log message should start with a lowercase letter`
	logger.Warn("what happened?")     // want `log message contains disallowed character`
	logger.Error("something broke!")  // want `log message contains disallowed character`
	logger.Debug("loading...")        // want `log message contains ellipsis`
	logger.Error("привет")            // want `log message contains non-ASCII character`

	apiKey := "k"
	logger.Info("api_key " + apiKey)  // want `log message may contain sensitive data` `log message may contain sensitive data`

	sugar.Info("this is fine")
	sugar.Info("Bad message")         // want `log message should start with a lowercase letter`
	sugar.Infow("token " + apiKey)    // want `log message may contain sensitive data`

	logger.Fatal("clean shutdown")
	logger.DPanic("clean debug panic")
}
