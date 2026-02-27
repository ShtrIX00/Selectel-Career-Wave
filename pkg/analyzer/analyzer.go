package analyzer

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer is the logcheck analyzer that checks log messages for style issues.
var Analyzer = &analysis.Analyzer{
	Name:     "logcheck",
	Doc:      "checks log messages for lowercase start, English-only text, no special characters, and no sensitive data",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

var cfg *Config

func init() {
	cfg = LoadConfigFromWorkDir()
}

// SetConfig allows overriding the config (used in tests).
func SetConfig(c *Config) {
	cfg = c
}

// slog methods that take a message as first argument
var slogFuncs = map[string]bool{
	"Info":    true,
	"Warn":    true,
	"Error":   true,
	"Debug":   true,
	"InfoContext":  true,
	"WarnContext":  true,
	"ErrorContext": true,
	"DebugContext": true,
}

// zap Logger methods that take a message as first argument
var zapFuncs = map[string]bool{
	"Info":   true,
	"Warn":   true,
	"Error":  true,
	"Debug":  true,
	"DPanic": true,
	"Panic":  true,
	"Fatal":  true,

	// Sugared logger variants (message is still the first argument)
	"Infof":   true,
	"Warnf":   true,
	"Errorf":  true,
	"Debugf":  true,
	"DPanicf": true,
	"Panicf":  true,
	"Fatalf":  true,

	"Infow":   true,
	"Warnw":   true,
	"Errorw":  true,
	"Debugw":  true,
	"DPanicw": true,
	"Panicw":  true,
	"Fatalw":  true,
}

func run(pass *analysis.Pass) (interface{}, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)
		checkCall(pass, call)
	})

	return nil, nil
}

func checkCall(pass *analysis.Pass, call *ast.CallExpr) {
	funcName, isSlog, isZap := classifyCall(call)
	if !isSlog && !isZap {
		return
	}
	_ = funcName

	// Determine the message argument index
	msgIdx := 0
	if isSlog {
		// slog.*Context functions have context as first arg, message as second
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if ok && strings.HasSuffix(sel.Sel.Name, "Context") {
			msgIdx = 1
		}
	}

	if len(call.Args) <= msgIdx {
		return
	}

	msgArg := call.Args[msgIdx]
	msg, parts, ok := extractMessage(msgArg)
	if !ok {
		return
	}

	// Run rule checks
	runChecks(pass, call, msgArg, msg, parts)
}

// classifyCall determines if a call expression is a slog or zap logging call.
func classifyCall(call *ast.CallExpr) (funcName string, isSlog bool, isZap bool) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return "", false, false
	}

	methodName := sel.Sel.Name

	// Check for package-level slog calls: slog.Info(...), etc.
	if ident, ok := sel.X.(*ast.Ident); ok {
		if ident.Name == "slog" && slogFuncs[methodName] {
			return methodName, true, false
		}
	}

	// For method calls on a variable (e.g., logger.Info(...)):
	// We check the method name against both slog and zap patterns.
	// This is a heuristic — we look at the method name.
	if slogFuncs[methodName] {
		// Could be slog.Logger method call or zap
		// Check if receiver looks like it could be zap (has Sugar methods, etc.)
		// For simplicity, we check both
		if isZapReceiver(sel.X) {
			return methodName, false, true
		}
		return methodName, true, false
	}

	if zapFuncs[methodName] && !slogFuncs[methodName] {
		// DPanic, Panic, Fatal are zap-only
		return methodName, false, true
	}

	return "", false, false
}

// isZapReceiver tries to determine if the receiver expression is a zap logger.
// This is a simple heuristic based on common variable names.
func isZapReceiver(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.Ident:
		name := strings.ToLower(e.Name)
		return strings.Contains(name, "zap") || name == "sugar"
	case *ast.SelectorExpr:
		return e.Sel.Name == "Logger" || e.Sel.Name == "SugaredLogger"
	case *ast.CallExpr:
		// e.g., zap.NewProduction().Sugar()
		return true
	}
	return false
}

// extractMessage extracts the string message from an AST expression.
// It returns the message string, any identifier parts found in concatenation,
// and whether extraction was successful.
func extractMessage(expr ast.Expr) (msg string, parts []string, ok bool) {
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind == token.STRING {
			// Remove quotes
			msg = strings.Trim(e.Value, "\"`")
			return msg, nil, true
		}
	case *ast.BinaryExpr:
		if e.Op == token.ADD {
			leftMsg, leftParts, leftOk := extractMessage(e.X)
			rightMsg, rightParts, rightOk := extractMessage(e.Y)

			var allParts []string
			allParts = append(allParts, leftParts...)
			allParts = append(allParts, rightParts...)

			// Collect identifiers for sensitive data check
			collectIdents(e, &allParts)

			if leftOk && rightOk {
				return leftMsg + rightMsg, allParts, true
			}
			if leftOk {
				return leftMsg, allParts, true
			}
			if rightOk {
				return rightMsg, allParts, true
			}
			// Neither side is a literal, but we found idents
			if len(allParts) > 0 {
				return "", allParts, false
			}
		}
	}
	return "", nil, false
}

// collectIdents collects all identifier names from an expression tree.
func collectIdents(expr ast.Expr, parts *[]string) {
	switch e := expr.(type) {
	case *ast.Ident:
		*parts = append(*parts, e.Name)
	case *ast.BinaryExpr:
		collectIdents(e.X, parts)
		collectIdents(e.Y, parts)
	case *ast.SelectorExpr:
		collectIdents(e.X, parts)
		*parts = append(*parts, e.Sel.Name)
	}
}

func runChecks(pass *analysis.Pass, call *ast.CallExpr, msgArg ast.Expr, msg string, parts []string) {
	if msg != "" {
		if !cfg.IsRuleDisabled("lowercase") {
			if f := checkLowercase(msg); f != nil {
				diag := analysis.Diagnostic{
					Pos:     msgArg.Pos(),
					End:     msgArg.End(),
					Message: f.Message,
				}
				if f.Suggestion != "" {
					diag.SuggestedFixes = []analysis.SuggestedFix{
						{
							Message: "lowercase the first letter",
							TextEdits: []analysis.TextEdit{
								{
									Pos:     msgArg.Pos(),
									End:     msgArg.End(),
									NewText: []byte(`"` + f.Suggestion + `"`),
								},
							},
						},
					}
				}
				pass.Report(diag)
			}
		}

		if !cfg.IsRuleDisabled("english-only") {
			if f := checkEnglishOnly(msg); f != nil {
				pass.Report(analysis.Diagnostic{
					Pos:     msgArg.Pos(),
					End:     msgArg.End(),
					Message: f.Message,
				})
			}
		}

		if !cfg.IsRuleDisabled("no-special-chars") {
			if f := checkSpecialChars(msg, cfg.AllowedSpecialChars); f != nil {
				pass.Report(analysis.Diagnostic{
					Pos:     msgArg.Pos(),
					End:     msgArg.End(),
					Message: f.Message,
				})
			}
		}

		// Sensitive keyword check inside the literal message text
		if !cfg.IsRuleDisabled("no-sensitive-data") {
			if f := checkSensitiveText(msg, cfg.sensitiveKeywordsSet); f != nil {
				pass.Report(analysis.Diagnostic{
					Pos:     msgArg.Pos(),
					End:     msgArg.End(),
					Message: f.Message,
				})
			}
		}
	}

	// Sensitive keyword check for identifiers involved in string concatenation
	if !cfg.IsRuleDisabled("no-sensitive-data") && len(parts) > 0 {
		if f := checkSensitiveData(parts, cfg.sensitiveKeywordsSet); f != nil {
			pass.Report(analysis.Diagnostic{
				Pos:     msgArg.Pos(),
				End:     msgArg.End(),
				Message: f.Message,
			})
		}
	}
}
