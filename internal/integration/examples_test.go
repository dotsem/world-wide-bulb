package integration_test

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type exampleTestCase struct {
	name        string
	dir         string
	command     string
	args        []string
	extraEnv    []string
	stateMatch  string
	reasonMatch string
}

func TestExamples_SSEStreams(t *testing.T) {
	repoRoot, err := filepath.Abs("../../")
	require.NoError(t, err)

	cases := []exampleTestCase{
		{
			name:        "Go SSE Client",
			dir:         filepath.Join(repoRoot, "examples", "sse", "go"),
			command:     "go",
			args:        []string{"run", "main.go"},
			stateMatch:  "State: ON",
			reasonMatch: "integration test reason",
		},
		{
			name:        "Python SSE Client",
			dir:         filepath.Join(repoRoot, "examples", "sse", "python"),
			command:     "python3",
			args:        []string{"-u", "client.py"},
			extraEnv:    []string{"PYTHONUNBUFFERED=1"},
			stateMatch:  "Lamp state is now: ON",
			reasonMatch: "integration test reason",
		},
		{
			name:        "JavaScript SSE Client",
			dir:         filepath.Join(repoRoot, "examples", "sse", "javascript"),
			command:     "node",
			args:        []string{"--experimental-eventsource", "client.js"},
			stateMatch:  "Lamp state changed: ON",
			reasonMatch: "integration test reason",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			env := setupIntegrationServer(t)
			eventsURL := env.baseURL + "/api/v1/events"

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			cmd := exec.CommandContext(ctx, tc.command, tc.args...) //nolint:gosec // Test runs static local client examples.
			cmd.Dir = tc.dir
			cmd.Env = append(os.Environ(), tc.extraEnv...)
			cmd.Env = append(cmd.Env, "WWB_URL="+eventsURL)

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

			waitForToken(t, outputLines, "Connecting to", 4*time.Second)
			time.Sleep(150 * time.Millisecond)

			toggleResp, err := env.server.Client().Post(env.baseURL+"/api/v1/toggle", "application/json", nil)
			require.NoError(t, err)
			defer func() { _ = toggleResp.Body.Close() }()
			require.Equal(t, http.StatusOK, toggleResp.StatusCode)

			var toggleRes struct {
				ID string `json:"id"`
			}
			err = json.NewDecoder(toggleResp.Body).Decode(&toggleRes)
			require.NoError(t, err)
			require.NotEmpty(t, toggleRes.ID)

			waitForToken(t, outputLines, tc.stateMatch, 4*time.Second)

			reasonPayload := fmt.Sprintf(`{"id":"%s","reason":"integration test reason"}`, toggleRes.ID)
			reasonResp, err := env.server.Client().Post(env.baseURL+"/api/v1/reason", "application/json", strings.NewReader(reasonPayload))
			require.NoError(t, err)
			defer func() { _ = reasonResp.Body.Close() }()
			require.Equal(t, http.StatusOK, reasonResp.StatusCode)

			waitForToken(t, outputLines, tc.reasonMatch, 4*time.Second)
		})
	}
}

func waitForToken(t *testing.T, lines <-chan string, token string, timeout time.Duration) {
	t.Helper()
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case line, ok := <-lines:
			if !ok {
				t.Fatalf("stdout closed before finding token: %q", token)
				return
			}
			t.Logf("stdout: %s", line)
			if strings.Contains(line, token) {
				return
			}
		case <-timer.C:
			t.Fatalf("timed out after %v waiting for token: %q", timeout, token)
			return
		}
	}
}
