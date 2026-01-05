package review

import (
	"os"
	"path/filepath"
	"testing"
)

// Helper function to create a temporary test file
func createTestFile(t *testing.T, dir, filename, content string) string {
	t.Helper()
	filePath := filepath.Join(dir, filename)
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	return filePath
}

// Helper to check if an issue with specific properties exists
func hasIssue(report *Report, issueType, severity, messageContains string) bool {
	for _, issue := range report.Issues {
		if issue.Type == issueType &&
			issue.Severity == severity &&
			(messageContains == "" || contains(issue.Message, messageContains)) {
			return true
		}
	}
	return false
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// ============== Python Analyzer Tests ==============

func TestPythonQuality_PrintStatement(t *testing.T) {
	tmpDir := t.TempDir()
	createTestFile(t, tmpDir, "test.py", `
def hello():
    print("Hello World")
`)
	analyzer := NewAnalyzer(tmpDir)
	report := NewReport()
	report.ChangedFiles = []string{"test.py"}
	analyzer.checkPythonQuality("test.py", report)

	if !hasIssue(report, "quality", "low", "print()") {
		t.Error("Expected print statement warning")
	}
}

func TestPythonQuality_DebuggerStatement(t *testing.T) {
	tmpDir := t.TempDir()
	createTestFile(t, tmpDir, "test.py", `
import pdb
pdb.set_trace()
breakpoint()
`)
	analyzer := NewAnalyzer(tmpDir)
	report := NewReport()
	report.ChangedFiles = []string{"test.py"}
	analyzer.checkPythonQuality("test.py", report)

	if !hasIssue(report, "quality", "medium", "Debugger") {
		t.Error("Expected debugger statement warning")
	}
}

func TestPythonSecurity_EvalUsage(t *testing.T) {
	tmpDir := t.TempDir()
	createTestFile(t, tmpDir, "test.py", `
result = eval(user_input)
exec(code)
`)
	analyzer := NewAnalyzer(tmpDir)
	report := NewReport()
	report.ChangedFiles = []string{"test.py"}
	analyzer.checkPythonQuality("test.py", report)

	if !hasIssue(report, "security", "high", "eval") {
		t.Error("Expected eval/exec security warning")
	}
}

func TestPythonSecurity_SubprocessShell(t *testing.T) {
	tmpDir := t.TempDir()
	createTestFile(t, tmpDir, "test.py", `
import subprocess
subprocess.run(cmd, shell=True)
`)
	analyzer := NewAnalyzer(tmpDir)
	report := NewReport()
	report.ChangedFiles = []string{"test.py"}
	analyzer.checkPythonQuality("test.py", report)

	if !hasIssue(report, "security", "medium", "shell=True") {
		t.Error("Expected subprocess shell=True warning")
	}
}

func TestPythonQuality_BareExcept(t *testing.T) {
	tmpDir := t.TempDir()
	createTestFile(t, tmpDir, "test.py", `
try:
    do_something()
except:
    pass
`)
	analyzer := NewAnalyzer(tmpDir)
	report := NewReport()
	report.ChangedFiles = []string{"test.py"}
	analyzer.checkPythonQuality("test.py", report)

	if !hasIssue(report, "quality", "medium", "Bare except") {
		t.Error("Expected bare except clause warning")
	}
}

func TestPythonSecurity_PickleLoad(t *testing.T) {
	tmpDir := t.TempDir()
	createTestFile(t, tmpDir, "test.py", `
import pickle
data = pickle.load(file)
`)
	analyzer := NewAnalyzer(tmpDir)
	report := NewReport()
	report.ChangedFiles = []string{"test.py"}
	analyzer.checkPythonQuality("test.py", report)

	if !hasIssue(report, "security", "high", "pickle") {
		t.Error("Expected pickle.load security warning")
	}
}

func TestPythonSecurity_SQLInjection(t *testing.T) {
	tmpDir := t.TempDir()
	createTestFile(t, tmpDir, "test.py", `
cursor.execute("SELECT * FROM users WHERE id = %s" % user_id)
cursor.execute(f"SELECT * FROM users WHERE name = '{name}'")
`)
	analyzer := NewAnalyzer(tmpDir)
	report := NewReport()
	report.ChangedFiles = []string{"test.py"}
	analyzer.checkPythonQuality("test.py", report)

	if !hasIssue(report, "security", "high", "SQL") {
		t.Error("Expected SQL injection warning")
	}
}

// ============== JavaScript Analyzer Tests ==============

func TestJavaScriptQuality_ConsoleLog(t *testing.T) {
	tmpDir := t.TempDir()
	createTestFile(t, tmpDir, "test.js", `
function hello() {
    console.log("Hello");
}
`)
	analyzer := NewAnalyzer(tmpDir)
	report := NewReport()
	report.ChangedFiles = []string{"test.js"}
	analyzer.checkJavaScriptQuality("test.js", report)

	if !hasIssue(report, "quality", "low", "console.log") {
		t.Error("Expected console.log warning")
	}
}

func TestJavaScriptQuality_Debugger(t *testing.T) {
	tmpDir := t.TempDir()
	createTestFile(t, tmpDir, "test.js", `
function test() {
    debugger
    return true;
}
`)
	analyzer := NewAnalyzer(tmpDir)
	report := NewReport()
	report.ChangedFiles = []string{"test.js"}
	analyzer.checkJavaScriptQuality("test.js", report)

	if !hasIssue(report, "quality", "medium", "debugger") {
		t.Error("Expected debugger statement warning")
	}
}

func TestJavaScriptSecurity_Eval(t *testing.T) {
	tmpDir := t.TempDir()
	createTestFile(t, tmpDir, "test.js", `
eval(userInput);
`)
	analyzer := NewAnalyzer(tmpDir)
	report := NewReport()
	report.ChangedFiles = []string{"test.js"}
	analyzer.checkJavaScriptQuality("test.js", report)

	if !hasIssue(report, "security", "high", "eval") {
		t.Error("Expected eval security warning")
	}
}

func TestJavaScriptSecurity_InnerHTML(t *testing.T) {
	tmpDir := t.TempDir()
	createTestFile(t, tmpDir, "test.js", `
element.innerHTML = userContent;
`)
	analyzer := NewAnalyzer(tmpDir)
	report := NewReport()
	report.ChangedFiles = []string{"test.js"}
	analyzer.checkJavaScriptQuality("test.js", report)

	if !hasIssue(report, "security", "high", "innerHTML") {
		t.Error("Expected innerHTML XSS warning")
	}
}

func TestJavaScriptSecurity_SSLDisabled(t *testing.T) {
	tmpDir := t.TempDir()
	createTestFile(t, tmpDir, "test.js", `
const options = { rejectUnauthorized: false };
`)
	analyzer := NewAnalyzer(tmpDir)
	report := NewReport()
	report.ChangedFiles = []string{"test.js"}
	analyzer.checkJavaScriptQuality("test.js", report)

	if !hasIssue(report, "security", "high", "SSL verification") {
		t.Error("Expected SSL verification disabled warning")
	}
}

// ============== TypeScript Analyzer Tests ==============

func TestTypeScriptQuality_AnyType(t *testing.T) {
	tmpDir := t.TempDir()
	createTestFile(t, tmpDir, "test.ts", `
function process(data: any): any {
    return data;
}
`)
	analyzer := NewAnalyzer(tmpDir)
	report := NewReport()
	report.ChangedFiles = []string{"test.ts"}
	analyzer.checkTypeScriptQuality("test.ts", report)

	if !hasIssue(report, "quality", "medium", "any") {
		t.Error("Expected 'any' type usage warning")
	}
}

func TestTypeScriptQuality_TsIgnore(t *testing.T) {
	tmpDir := t.TempDir()
	createTestFile(t, tmpDir, "test.ts", `
// @ts-ignore
const x: string = 123;
`)
	analyzer := NewAnalyzer(tmpDir)
	report := NewReport()
	report.ChangedFiles = []string{"test.ts"}
	analyzer.checkTypeScriptQuality("test.ts", report)

	if !hasIssue(report, "quality", "medium", "ignore") {
		t.Error("Expected @ts-ignore warning")
	}
}

func TestTypeScriptSecurity_FunctionConstructor(t *testing.T) {
	tmpDir := t.TempDir()
	createTestFile(t, tmpDir, "test.ts", `
const fn = new Function(userCode);
`)
	analyzer := NewAnalyzer(tmpDir)
	report := NewReport()
	report.ChangedFiles = []string{"test.ts"}
	analyzer.checkTypeScriptQuality("test.ts", report)

	if !hasIssue(report, "security", "high", "Function") {
		t.Error("Expected Function constructor warning")
	}
}

// ============== Ruby Analyzer Tests ==============

func TestRubyQuality_DebuggerStatement(t *testing.T) {
	tmpDir := t.TempDir()
	createTestFile(t, tmpDir, "test.rb", `
def debug_method
  binding.pry
  byebug
end
`)
	analyzer := NewAnalyzer(tmpDir)
	report := NewReport()
	report.ChangedFiles = []string{"test.rb"}
	analyzer.checkRubyQuality("test.rb", report)

	if !hasIssue(report, "quality", "medium", "Debugger") {
		t.Error("Expected debugger statement warning")
	}
}

func TestRubySecurity_Eval(t *testing.T) {
	tmpDir := t.TempDir()
	createTestFile(t, tmpDir, "test.rb", `
result = eval(user_input)
instance_eval(code)
`)
	analyzer := NewAnalyzer(tmpDir)
	report := NewReport()
	report.ChangedFiles = []string{"test.rb"}
	analyzer.checkRubyQuality("test.rb", report)

	if !hasIssue(report, "security", "high", "eval") {
		t.Error("Expected eval security warning")
	}
}

func TestRubySecurity_UnsafeYAML(t *testing.T) {
	tmpDir := t.TempDir()
	createTestFile(t, tmpDir, "test.rb", `
data = YAML.load(user_input)
`)
	analyzer := NewAnalyzer(tmpDir)
	report := NewReport()
	report.ChangedFiles = []string{"test.rb"}
	analyzer.checkRubyQuality("test.rb", report)

	if !hasIssue(report, "security", "high", "YAML") {
		t.Error("Expected unsafe YAML.load warning")
	}
}

func TestRubySecurity_HTMLSafe(t *testing.T) {
	tmpDir := t.TempDir()
	createTestFile(t, tmpDir, "test.rb", `
<%= user_input.html_safe %>
`)
	analyzer := NewAnalyzer(tmpDir)
	report := NewReport()
	report.ChangedFiles = []string{"test.rb"}
	analyzer.checkRubyQuality("test.rb", report)

	if !hasIssue(report, "security", "high", "XSS") {
		t.Error("Expected XSS warning for html_safe")
	}
}

// ============== Dart Analyzer Tests ==============

func TestDartQuality_PrintStatement(t *testing.T) {
	tmpDir := t.TempDir()
	createTestFile(t, tmpDir, "test.dart", `
void main() {
  print("Hello");
  debugPrint("Debug");
}
`)
	analyzer := NewAnalyzer(tmpDir)
	report := NewReport()
	report.ChangedFiles = []string{"test.dart"}
	analyzer.checkDartQuality("test.dart", report)

	if !hasIssue(report, "quality", "low", "print()") {
		t.Error("Expected print statement warning")
	}
}

func TestDartQuality_DynamicType(t *testing.T) {
	tmpDir := t.TempDir()
	createTestFile(t, tmpDir, "test.dart", `
dynamic data = fetchData();
List<dynamic> items = [];
`)
	analyzer := NewAnalyzer(tmpDir)
	report := NewReport()
	report.ChangedFiles = []string{"test.dart"}
	analyzer.checkDartQuality("test.dart", report)

	if !hasIssue(report, "quality", "medium", "dynamic") {
		t.Error("Expected dynamic type warning")
	}
}

func TestDartSecurity_HardcodedCredentials(t *testing.T) {
	tmpDir := t.TempDir()
	createTestFile(t, tmpDir, "test.dart", `
const apiKey = "sk_live_12345";
const password = "secret123";
`)
	analyzer := NewAnalyzer(tmpDir)
	report := NewReport()
	report.ChangedFiles = []string{"test.dart"}
	analyzer.checkDartQuality("test.dart", report)

	if !hasIssue(report, "security", "high", "credential") {
		t.Error("Expected hardcoded credential warning")
	}
}

// ============== PHP Analyzer Tests ==============

func TestPHPQuality_VarDump(t *testing.T) {
	tmpDir := t.TempDir()
	createTestFile(t, tmpDir, "test.php", `<?php
var_dump($data);
print_r($array);
?>`)
	analyzer := NewAnalyzer(tmpDir)
	report := NewReport()
	report.ChangedFiles = []string{"test.php"}
	analyzer.checkPHPQuality("test.php", report)

	if !hasIssue(report, "quality", "low", "var_dump") {
		t.Error("Expected var_dump warning")
	}
}

func TestPHPSecurity_Eval(t *testing.T) {
	tmpDir := t.TempDir()
	createTestFile(t, tmpDir, "test.php", `<?php
eval($_POST['code']);
?>`)
	analyzer := NewAnalyzer(tmpDir)
	report := NewReport()
	report.ChangedFiles = []string{"test.php"}
	analyzer.checkPHPQuality("test.php", report)

	if !hasIssue(report, "security", "high", "eval") {
		t.Error("Expected eval security warning")
	}
}

func TestPHPSecurity_SQLInjection(t *testing.T) {
	tmpDir := t.TempDir()
	createTestFile(t, tmpDir, "test.php", `<?php
$result = mysql_query("SELECT * FROM users WHERE id = " . $_GET['id']);
?>`)
	analyzer := NewAnalyzer(tmpDir)
	report := NewReport()
	report.ChangedFiles = []string{"test.php"}
	analyzer.checkPHPQuality("test.php", report)

	if !hasIssue(report, "security", "high", "SQL injection") {
		t.Error("Expected SQL injection warning")
	}
}

func TestPHPSecurity_XSS(t *testing.T) {
	tmpDir := t.TempDir()
	createTestFile(t, tmpDir, "test.php", `<?php
echo $_GET['name'];
?>`)
	analyzer := NewAnalyzer(tmpDir)
	report := NewReport()
	report.ChangedFiles = []string{"test.php"}
	analyzer.checkPHPQuality("test.php", report)

	if !hasIssue(report, "security", "high", "XSS") {
		t.Error("Expected XSS warning")
	}
}

// ============== Java/Kotlin Analyzer Tests ==============

func TestJavaQuality_SystemOut(t *testing.T) {
	tmpDir := t.TempDir()
	createTestFile(t, tmpDir, "Test.java", `
public class Test {
    public void log() {
        System.out.println("Debug");
    }
}
`)
	analyzer := NewAnalyzer(tmpDir)
	report := NewReport()
	report.ChangedFiles = []string{"Test.java"}
	analyzer.checkJavaKotlinQuality("Test.java", report)

	if !hasIssue(report, "quality", "low", "System.out") {
		t.Error("Expected System.out.println warning")
	}
}

func TestJavaQuality_PrintStackTrace(t *testing.T) {
	tmpDir := t.TempDir()
	createTestFile(t, tmpDir, "Test.java", `
try {
    doSomething();
} catch (Exception e) {
    e.printStackTrace();
}
`)
	analyzer := NewAnalyzer(tmpDir)
	report := NewReport()
	report.ChangedFiles = []string{"Test.java"}
	analyzer.checkJavaKotlinQuality("Test.java", report)

	if !hasIssue(report, "quality", "medium", "printStackTrace") {
		t.Error("Expected printStackTrace warning")
	}
}

func TestJavaSecurity_ProcessExecution(t *testing.T) {
	tmpDir := t.TempDir()
	createTestFile(t, tmpDir, "Test.java", `
Runtime.getRuntime().exec(command);
`)
	analyzer := NewAnalyzer(tmpDir)
	report := NewReport()
	report.ChangedFiles = []string{"Test.java"}
	analyzer.checkJavaKotlinQuality("Test.java", report)

	if !hasIssue(report, "security", "medium", "Process") {
		t.Error("Expected process execution warning")
	}
}

func TestJavaSecurity_WeakCrypto(t *testing.T) {
	tmpDir := t.TempDir()
	createTestFile(t, tmpDir, "Test.java", `
MessageDigest md = MessageDigest.getInstance("MD5");
`)
	analyzer := NewAnalyzer(tmpDir)
	report := NewReport()
	report.ChangedFiles = []string{"Test.java"}
	analyzer.checkJavaKotlinQuality("Test.java", report)

	if !hasIssue(report, "security", "medium", "Weak") {
		t.Error("Expected weak cryptography warning")
	}
}

func TestKotlinQuality_ForceUnwrap(t *testing.T) {
	tmpDir := t.TempDir()
	createTestFile(t, tmpDir, "Test.kt", `
val name = user!!.name
val length = text!!.length
`)
	analyzer := NewAnalyzer(tmpDir)
	report := NewReport()
	report.ChangedFiles = []string{"Test.kt"}
	analyzer.checkJavaKotlinQuality("Test.kt", report)

	if !hasIssue(report, "quality", "medium", "!!") {
		t.Error("Expected force unwrap warning")
	}
}

// ============== Core Analyzer Tests ==============

func TestAnalyzer_IgnoreFile(t *testing.T) {
	tmpDir := t.TempDir()
	// Create .autoreview-ignore file
	createTestFile(t, tmpDir, ".autoreview-ignore", `
vendor/
*.min.js
test_data/
`)
	analyzer := NewAnalyzer(tmpDir)

	tests := []struct {
		path     string
		expected bool
	}{
		{"vendor/package/file.go", true},
		{"src/main.go", false},
		{"bundle.min.js", true},
		{"test_data/sample.json", true},
		{"app/controller.rb", false},
	}

	for _, tt := range tests {
		result := analyzer.shouldIgnoreFile(tt.path)
		if result != tt.expected {
			t.Errorf("shouldIgnoreFile(%q) = %v, want %v", tt.path, result, tt.expected)
		}
	}
}

func TestReport_AddIssue(t *testing.T) {
	report := NewReport()

	report.AddIssue(Issue{Type: "security", Severity: "high", Message: "Test high"})
	report.AddIssue(Issue{Type: "quality", Severity: "medium", Message: "Test medium"})
	report.AddIssue(Issue{Type: "quality", Severity: "low", Message: "Test low"})

	if report.Summary.TotalIssues != 3 {
		t.Errorf("Expected 3 total issues, got %d", report.Summary.TotalIssues)
	}
	if report.Summary.HighSeverity != 1 {
		t.Errorf("Expected 1 high severity, got %d", report.Summary.HighSeverity)
	}
	if report.Summary.MediumSeverity != 1 {
		t.Errorf("Expected 1 medium severity, got %d", report.Summary.MediumSeverity)
	}
	if report.Summary.LowSeverity != 1 {
		t.Errorf("Expected 1 low severity, got %d", report.Summary.LowSeverity)
	}
}
