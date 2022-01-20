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
