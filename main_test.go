package main

import (
	"os/exec"
	"testing"
)

type TestCase struct {
	name       string
	cmd        string
	args       []string
	unixCmd    string
	returnCode int
}

func TestGrepCommand(t *testing.T) {
	testCases := []TestCase{

		{
			name:       "Recurse",
			cmd:        "./ccgrep",
			unixCmd:    "grep",
			args:       []string{"-r", "Nirvana", "./testdata"},
			returnCode: 0,
		},

		{
			name:       "One letter pattern",
			cmd:        "./ccgrep",
			unixCmd:    "grep",
			args:       []string{"J", "./testdata/rockbands.txt"},
			returnCode: 0,
		},
		{
			name:       "Case insensitive",
			cmd:        "./ccgrep",
			unixCmd:    "grep",
			args:       []string{"-i", "A", "./testdata/rockbands.txt"},
			returnCode: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command(tc.cmd, tc.args...)
			got, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("Command %v failed with error: %v, output: %s", cmd.Args, err, string(got))
			}
			unixCmd := exec.Command(tc.unixCmd, tc.args...)
			want, err := unixCmd.CombinedOutput()
			if err != nil {
				t.Fatalf("Command %v failed with error: %v, output: %s", unixCmd.Args, err, string(want))
			}
			if string(got) != string(want) {
				t.Errorf("Testcase %s failed, got %v but wanted %v", tc.name, string(got), string(want))
			}
		})

	}

}
