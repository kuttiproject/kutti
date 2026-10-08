package node

import (
	"fmt"
	"os"
	"testing"
)

func TestCPArg(t *testing.T) {

	testCases := []struct {
		arg           string
		errorexpected bool
		expected      *cparg
	}{
		{
			arg:           "justafile",
			errorexpected: false,
			expected: &cparg{
				nodename:          "",
				filepath:          "justafile",
				hasnodename:       false,
				iswindowsfilepath: false,
				localfileexists:   false,
				localisdirectory:  false,
			},
		},
		{
			arg:           "node1:file1",
			errorexpected: false,
			expected: &cparg{
				nodename:          "node1",
				filepath:          "file1",
				hasnodename:       true,
				iswindowsfilepath: false,
				localfileexists:   false,
				localisdirectory:  false,
			},
		},
		{
			arg:           "c:file1",
			errorexpected: false,
			expected: &cparg{
				nodename:          "",
				filepath:          "c:file1",
				hasnodename:       false,
				iswindowsfilepath: true,
				localfileexists:   false,
				localisdirectory:  false,
			},
		},
		{
			arg:           "C:file1",
			errorexpected: false,
			expected: &cparg{
				nodename:          "",
				filepath:          "C:file1",
				hasnodename:       false,
				iswindowsfilepath: true,
				localfileexists:   false,
				localisdirectory:  false,
			},
		},
		{
			arg:           "invalid:because:multiple:colons",
			errorexpected: true,
			expected: &cparg{
				nodename:          "",
				filepath:          "",
				hasnodename:       false,
				iswindowsfilepath: false,
				localfileexists:   false,
				localisdirectory:  false,
			},
		},
		{
			arg:           "Invalid:becausenocapsinnodename",
			errorexpected: true,
			expected: &cparg{
				nodename:          "",
				filepath:          "",
				hasnodename:       false,
				iswindowsfilepath: false,
				localfileexists:   false,
				localisdirectory:  false,
			},
		},
		{
			arg:           "node1:testfile.tmp",
			errorexpected: false,
			expected: &cparg{
				nodename:          "node1",
				filepath:          "testfile.tmp",
				hasnodename:       true,
				iswindowsfilepath: false,
				localfileexists:   true,
				localisdirectory:  false,
			},
		},
		{
			arg:           "node1:testdir",
			errorexpected: false,
			expected: &cparg{
				nodename:          "node1",
				filepath:          "testdir",
				hasnodename:       true,
				iswindowsfilepath: false,
				localfileexists:   true,
				localisdirectory:  true,
			},
		},
		{
			arg:           "node-1:file1",
			errorexpected: false,
			expected: &cparg{
				nodename:          "node-1",
				filepath:          "file1",
				hasnodename:       true,
				iswindowsfilepath: false,
				localfileexists:   false,
				localisdirectory:  false,
			},
		},
		{
			arg:           "k8s_worker_node-01:file1",
			errorexpected: false,
			expected: &cparg{
				nodename:          "k8s_worker_node-01",
				filepath:          "file1",
				hasnodename:       true,
				iswindowsfilepath: false,
				localfileexists:   false,
				localisdirectory:  false,
			},
		},
		{
			arg:           "n1:file1",
			errorexpected: false,
			expected: &cparg{
				nodename:          "n1",
				filepath:          "file1",
				hasnodename:       true,
				iswindowsfilepath: false,
				localfileexists:   false,
				localisdirectory:  false,
			},
		},
		{
			arg:           "ab:file1",
			errorexpected: false,
			expected: &cparg{
				nodename:          "ab",
				filepath:          "file1",
				hasnodename:       true,
				iswindowsfilepath: false,
				localfileexists:   false,
				localisdirectory:  false,
			},
		},
		{
			// 63 characters (1 letter + 62 chars) - maximum valid length
			arg:           "a012345678901234567890123456789012345678901234567890123456789_-:file1",
			errorexpected: false,
			expected: &cparg{
				nodename:          "a012345678901234567890123456789012345678901234567890123456789_-",
				filepath:          "file1",
				hasnodename:       true,
				iswindowsfilepath: false,
				localfileexists:   false,
				localisdirectory:  false,
			},
		},
		{
			// 64 characters - exceeds maximum allowed length
			arg:           "a012345678901234567890123456789012345678901234567890123456789_--:file1",
			errorexpected: true,
			expected: &cparg{
				nodename:          "",
				filepath:          "",
				hasnodename:       false,
				iswindowsfilepath: false,
				localfileexists:   false,
				localisdirectory:  false,
			},
		},
		{
			// Starts with hyphen - invalid
			arg:           "-node1:file1",
			errorexpected: true,
			expected: &cparg{
				nodename:          "",
				filepath:          "",
				hasnodename:       false,
				iswindowsfilepath: false,
				localfileexists:   false,
				localisdirectory:  false,
			},
		},
		{
			// Starts with underscore - invalid
			arg:           "_node1:file1",
			errorexpected: true,
			expected: &cparg{
				nodename:          "",
				filepath:          "",
				hasnodename:       false,
				iswindowsfilepath: false,
				localfileexists:   false,
				localisdirectory:  false,
			},
		},
		{
			// Starts with digit - invalid
			arg:           "1node:file1",
			errorexpected: true,
			expected: &cparg{
				nodename:          "",
				filepath:          "",
				hasnodename:       false,
				iswindowsfilepath: false,
				localfileexists:   false,
				localisdirectory:  false,
			},
		},
		{
			// Contains dot - invalid character
			arg:           "node.1:file1",
			errorexpected: true,
			expected: &cparg{
				nodename:          "",
				filepath:          "",
				hasnodename:       false,
				iswindowsfilepath: false,
				localfileexists:   false,
				localisdirectory:  false,
			},
		},
		{
			arg:           "node1:ud/\\/\\|\000[]*?",
			errorexpected: true,
			expected: &cparg{
				nodename:          "node1",
				filepath:          "testdir",
				hasnodename:       true,
				iswindowsfilepath: false,
				localfileexists:   true,
				localisdirectory:  true,
			},
		},
	}

	f, _ := os.Create("testfile.tmp")
	defer os.Remove("testfile.tmp")
	fmt.Fprintln(f, "Testing")
	f.Close()

	_ = os.MkdirAll("testdir", 0755)
	defer os.RemoveAll("testdir")

	for _, tc := range testCases {
		result, err := parseCPArg(tc.arg)
		if err == nil {
			if tc.errorexpected {
				t.Fatalf("case '%v' expected an error. Didn't happen", tc.arg)
			}
		} else {
			if !tc.errorexpected {
				t.Fatalf("case '%v' failed with unexpected error: %v", tc.arg, err)
			}

			continue
		}

		if result.nodename != tc.expected.nodename ||
			result.filepath != tc.expected.filepath ||
			result.hasnodename != tc.expected.hasnodename ||
			result.iswindowsfilepath != tc.expected.iswindowsfilepath ||
			result.localfileexists != tc.expected.localfileexists ||
			result.localisdirectory != tc.expected.localisdirectory {

			wd, _ := os.Getwd()
			t.Logf("Directory: %v", wd)
			t.Fatalf("case '%v': expected '%+v', got '%+v'", tc.arg, tc.expected, result)
		}

	}

}
