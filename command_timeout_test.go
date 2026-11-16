package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/adnanh/webhook/internal/hook"
	"github.com/gorilla/mux"
)

// TestCommandTimeoutHelper runs only in short, synthetic child processes.
func TestCommandTimeoutHelper(t *testing.T) {
	mode := os.Getenv("WEBHOOK_TIMEOUT_HELPER")
	if mode == "" {
		return
	}
	fmt.Println("started")
	switch mode {
	case "sleep":
		time.Sleep(3 * time.Second)
	case "exit":
		os.Exit(7)
	case "pipe":
		child := exec.Command(os.Args[0], "-test.run=^TestCommandTimeoutHelper$")
		child.Env = append(os.Environ(), "WEBHOOK_TIMEOUT_HELPER=sleep")
		child.Stdout, child.Stderr = os.Stdout, os.Stderr
		if err := child.Start(); err != nil {
			panic(err)
		}
		// The child is intentionally finite and is cleaned up by the test parent.
		fmt.Printf("child-pid=%d\n", child.Process.Pid)
	}
	fmt.Println("finished")
	os.Exit(0)
}

func timeoutTestHook(t *testing.T, timeout, mode string) *hook.Hook {
	t.Helper()
	var h hook.Hook
	data, _ := json.Marshal(map[string]string{"command-timeout": timeout})
	if err := json.Unmarshal(data, &h); err != nil {
		t.Fatal(err)
	}
	h.ID = "timeout-test"
	h.ExecuteCommand = os.Args[0]
	h.PassArgumentsToCommand = []hook.Argument{{Source: "string", Name: "-test.run=^TestCommandTimeoutHelper$"}}
	h.PassEnvironmentToCommand = []hook.Argument{
		{Source: "string", Name: mode, EnvName: "WEBHOOK_TIMEOUT_HELPER"},
		{Source: "string", Name: "atexit_sleep_ms=0", EnvName: "GORACE"},
	}
	return &h
}

func TestCommandTimeout(t *testing.T) {
	h := timeoutTestHook(t, "500ms", "sleep")
	h.CommandWorkingDirectory = t.TempDir()
	h.PassFileToCommand = []hook.Argument{{Source: "string", Name: "synthetic payload", EnvName: "PAYLOAD_FILE"}}
	start := time.Now()
	out, err := handleHook(h, &hook.Request{ID: "timeout-test"})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected deadline error, got %v; output=%q", err, out)
	}
	if elapsed := time.Since(start); elapsed > 2*time.Second {
		t.Errorf("execution exceeded timeout budget: %v", elapsed)
	}
	if !strings.Contains(out, "started") || strings.Contains(out, "finished") {
		t.Errorf("unexpected partial output %q", out)
	}
	files, err := os.ReadDir(h.CommandWorkingDirectory)
	if err != nil || len(files) != 0 {
		t.Fatalf("temporary files not cleaned: %v, %v", files, err)
	}
}

func TestCommandTimeoutCompatibility(t *testing.T) {
	for _, timeout := range []string{"", "0", "0s", "5s"} {
		t.Run(timeout, func(t *testing.T) {
			out, err := handleHook(timeoutTestHook(t, timeout, "quick"), &hook.Request{})
			if err != nil || !strings.Contains(out, "finished") {
				t.Fatalf("output=%q err=%v", out, err)
			}
		})
	}
}

func TestCommandTimeoutInvalid(t *testing.T) {
	for _, timeout := range []string{"bad", "-1s", "999999999999999999999h"} {
		t.Run(timeout, func(t *testing.T) {
			out, err := handleHook(timeoutTestHook(t, timeout, "quick"), &hook.Request{})
			if err == nil || !strings.Contains(err.Error(), "invalid command-timeout") || out != "" {
				t.Fatalf("output=%q err=%v", out, err)
			}
		})
	}
}

func TestCommandTimeoutPreservesExitError(t *testing.T) {
	out, err := handleHook(timeoutTestHook(t, "5s", "exit"), &hook.Request{})
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) || exitErr.ExitCode() != 7 || !strings.Contains(out, "started") {
		t.Fatalf("output=%q err=%v", out, err)
	}
}

func TestCommandTimeoutInheritedPipe(t *testing.T) {
	start := time.Now()
	out, err := handleHook(timeoutTestHook(t, "10s", "pipe"), &hook.Request{})
	for _, line := range strings.Split(out, "\n") {
		var pid int
		if _, scanErr := fmt.Sscanf(line, "child-pid=%d", &pid); scanErr == nil {
			if child, findErr := os.FindProcess(pid); findErr == nil {
				_ = child.Kill()
				_ = child.Release()
			}
		}
	}
	if !errors.Is(err, exec.ErrWaitDelay) {
		t.Fatalf("expected inherited pipe deadline, got %v, output=%q", err, out)
	}
	if elapsed := time.Since(start); elapsed > 2500*time.Millisecond {
		t.Errorf("pipe wait not bounded: %v", elapsed)
	}
}

func TestCommandTimeoutHTTPError(t *testing.T) {
	previous := loadedHooksFromFiles
	defer func() { loadedHooksFromFiles = previous }()
	for _, capture := range []bool{false, true} {
		h := timeoutTestHook(t, "500ms", "sleep")
		h.CaptureCommandOutput = true
		h.CaptureCommandOutputOnError = capture
		loadedHooksFromFiles = map[string]hook.Hooks{"test": {*h}}
		r := mux.SetURLVars(httptest.NewRequest(http.MethodPost, "/hooks/timeout-test", nil), map[string]string{"id": h.ID})
		w := httptest.NewRecorder()
		hookHandler(w, r)
		if w.Code != http.StatusInternalServerError {
			t.Fatalf("status=%d", w.Code)
		}
		if capture && !strings.Contains(w.Body.String(), "started") {
			t.Fatalf("partial output missing: %q", w.Body.String())
		}
		if !capture && strings.Contains(w.Body.String(), "started") {
			t.Fatal("partial output exposed without opt-in")
		}
	}
}

func TestCommandTimeoutConfigFormats(t *testing.T) {
	for name, data := range map[string]string{"json": `[{"id":"test","command-timeout":"500ms"}]`, "yaml": "- id: test\n  command-timeout: 500ms\n"} {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "hooks."+name)
			if err := os.WriteFile(path, []byte(data), 0600); err != nil {
				t.Fatal(err)
			}
			var hooks hook.Hooks
			if err := hooks.LoadFromFile(path, false); err != nil {
				t.Fatal(err)
			}
			if len(hooks) != 1 || hooks[0].CommandTimeout != "500ms" {
				t.Fatalf("unexpected config: %+v", hooks)
			}
		})
	}
}
