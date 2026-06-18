package root

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestRootCommandDoesNotAdvertiseRetiredAlias(t *testing.T) {
	if strings.Contains(rootCmd.Long, "manus skill") || strings.Contains(rootCmd.Long, "Compatibility alias") {
		t.Fatalf("root command long help still advertises retired alias")
	}

	output := captureStdout(t, func() {
		rootCmd.Run(rootCmd, nil)
	})
	if strings.Contains(output, "Compatibility alias") || strings.Contains(output, "manus skill") {
		t.Fatalf("root command output still advertises retired alias: %q", output)
	}
	if !strings.Contains(output, "provider-api") {
		t.Fatalf("root command output should advertise the provider API adapter: %q", output)
	}
	if strings.Contains(output, "Manus API") {
		t.Fatalf("root command output should not advertise Manus as the primary API adapter: %q", output)
	}
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w

	fn()

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	os.Stdout = old

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatal(err)
	}
	return buf.String()
}
