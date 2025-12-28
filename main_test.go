package main

import (
	"os/exec"
	"regexp"
	"testing"

	ce "github.com/foursixnine/imdblookup/internal/errors"
	"github.com/foursixnine/imdblookup/tests"
)

func TestMain(t *testing.T) {
	testCases := map[string]struct {
		params   string
		expected string
		exitcode int
	}{
		"with empty query": {
			expected: ".*Search title cannot be empty.*",
			params:   "",
			exitcode: ce.EMPTYQUERYERROR,
		},
		"with Stranger Things": {
			expected: `\(foobar\).*"Stranger Things"`,
			params:   "Stranger Things",
			// exitcode: ce.SUCCESS, Success should not pupulate err
		},
		"with broken api": {
			expected: `api url does not have scheme: 'localhost:22/'`,
			params:   "Stranger Things",
			exitcode: ce.GENERICERROR,
		},
	}
	server := tests.SetupServer(t)
	defer server.Close()
	apiurl := server.URL

	cmd := exec.Command("go", "build", "-o", "test_binary", ".")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}
	defer exec.Command("rm", "test_binary").Run()

	for testName, testCase := range testCases {

		if testName == "with broken api" {
			apiurl = "localhost:22/"
		}

		cmd = exec.Command("./test_binary", "--query", testCase.params, "--api", apiurl)
		output, err := cmd.CombinedOutput()
		exitCode := cmd.ProcessState.ExitCode()

		if err != nil {
			if exitCode != testCase.exitcode {
				t.Fatalf("Unexpected exit code, expected 0, '%s' got %d, expected %d. Output:\n%serror:\n%v\n", testName, exitCode, testCase.exitcode, string(output), err)
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
