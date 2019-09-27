package main

import "testing"

func TestRequiresRun(t *testing.T) {
	tableTestData := []struct {
		changedFiles   []string
		dependencies   []string
		expectedResult bool
	}{
		{[]string{"package.json"}, []string{"services/auth"}, false},
		{[]string{"services/auth/package.json"}, []string{"services/auth"}, true},
		{[]string{"services/auth/src/index.ts"}, []string{"services/auth"}, true},
		{[]string{"package.json", "services/auth/package.json"}, []string{"services/auth"}, true},
		{[]string{"package.json", "packages/core/package.json"}, []string{"services/auth", "packages/core"}, true},
	}
	for _, testData := range tableTestData {
		result := requiresRun(testData.changedFiles, testData.dependencies)
		if result != testData.expectedResult {
			t.Errorf("error: changed files %v for dependencies %v should have been '%v', but was '%v'", testData.changedFiles, testData.dependencies, testData.expectedResult, result)
		}
	}
}
