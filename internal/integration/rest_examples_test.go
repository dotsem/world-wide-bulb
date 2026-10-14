//go:build examples

package integration_test

import (
	"bufio"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type restExampleCase struct {
	name     string
	dir      string
	command  string
	args     []string
	extraEnv []string
}

func TestExamples_RESTClients(t *testing.T) {
	repoRoot, err := filepath.Abs("../../")
	require.NoError(t, err)

	cases := []restExampleCase{
		{
			name:    "Go REST Client",
			dir:     filepath.Join(repoRoot, "examples", "rest", "go"),
			command: "go",
			args:    []string{"run", "main.go"},
		},
		{
			name:     "Python REST Client",
			dir:      filepath.Join(repoRoot, "examples", "rest", "python"),
			command:  "python3",
			args:     []string{"-u", "client.py"},
			extraEnv: []string{"PYTHONUNBUFFERED=1"},
		},
		{
			name:    "JavaScript REST Client",
			dir:     filepath.Join(repoRoot, "examples", "rest", "javascript"),
			command: "node",
			args:    []string{"client.js"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			env := setupIntegrationServer(t)

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			cmd := exec.CommandContext(ctx, tc.command, tc.args...) //nolint:gosec // Test runs static local client examples.
			cmd.Dir = tc.dir
			cmd.Env = append(os.Environ(), tc.extraEnv...)
			cmd.Env = append(cmd.Env, "WWB_URL="+env.baseURL)

			cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

			stdout, err := cmd.StdoutPipe()
			require.NoError(t, err)
			stderr, err := cmd.StderrPipe()
			require.NoError(t, err)

			require.NoError(t, cmd.Start())
			defer func() {
				if cmd.Process != nil {
					_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
				}
				_ = cmd.Wait()
			}()

			outputLines := make(chan string, 100)
			go func() {
				scanner := bufio.NewScanner(stdout)
				for scanner.Scan() {
					outputLines <- scanner.Text()
				}
			}()
			go func() {
				scanner := bufio.NewScanner(stderr)
				for scanner.Scan() {
					t.Logf("[%s stderr] %s", tc.name, scanner.Text())
				}
			}()

			waitForToken(t, outputLines, "Connecting to REST API", 4*time.Second)
			waitForToken(t, outputLines, "State: OFF", 4*time.Second)
			waitForToken(t, outputLines, "History entries:", 4*time.Second)
			waitForToken(t, outputLines, "Toggled bulb, ID:", 4*time.Second)
			waitForToken(t, outputLines, "Reason updated for", 4*time.Second)
		})
	}
}
