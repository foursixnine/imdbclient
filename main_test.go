package main

import (
	"os/exec"
	"regexp"
	"testing"

	"github.com/foursixnine/imdblookup/tests"
)

func TestMain(t *testing.T) {
	testCases := map[string]struct {
		params   string
		expected string
	}{
		"with empty query": {
			expected: ".*Search title cannot be empty.*",
			params:   "",
		},
		"with Stranger Things": {
			expected: `\(foobar\).*"Stranger Things"`,
			params:   "Stranger Things",
		},
	}
	server := tests.SetupServer(t)
	defer server.Close()

	cmd := exec.Command("go", "build", "-o", "test_binary", ".")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}
	defer exec.Command("rm", "test_binary").Run()

	for testName, testCase := range testCases {
		cmd = exec.Command("./test_binary", "--query", testCase.params, "--api", server.URL)
		output, err := cmd.CombinedOutput()
		exitCode := cmd.ProcessState.ExitCode()

		if err != nil {
			switch exitCode {
			case 3:
				t.Logf("Properly failed when expected: %s", testName)
			default:
				t.Errorf("Unexpected exit code, expected 0, %s got %d. Output: %s, error: %v", testName, exitCode, string(output), err)
			}
		}

		expected_regex := regexp.MustCompile(testCase.expected)
		if !expected_regex.MatchString(string(output)) {
			t.Errorf("%s failed, got \n%s\n expected: %s", testName, output, testCase.expected)
		} else {
			t.Logf("%s passed", testName)
		}
	}

}
