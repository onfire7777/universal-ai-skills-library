package api

import "testing"

func TestNormalizeAPIBase(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "https", input: "https://api.example.test/v2/", want: "https://api.example.test/v2"},
		{name: "localhost http", input: "http://localhost:8080/v2/", want: "http://localhost:8080/v2"},
		{name: "loopback http", input: "http://127.0.0.1:8080/v2", want: "http://127.0.0.1:8080/v2"},
		{name: "missing scheme", input: "api.example.test/v2", wantErr: true},
		{name: "unsupported scheme", input: "file:///tmp/api", wantErr: true},
		{name: "nonlocal http", input: "http://api.example.test/v2", wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := normalizeAPIBase(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("normalizeAPIBase(%q) succeeded, want error", tc.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("normalizeAPIBase(%q) error = %v", tc.input, err)
			}
			if got != tc.want {
				t.Fatalf("normalizeAPIBase(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestAPIURLUsesNormalizedBase(t *testing.T) {
	old := apiBaseURL
	t.Cleanup(func() { apiBaseURL = old })

	apiBaseURL = "https://api.example.test/v2/"
	got, err := apiURL("/tasks")
	if err != nil {
		t.Fatalf("apiURL returned error: %v", err)
	}
	if want := "https://api.example.test/v2/tasks"; got != want {
		t.Fatalf("apiURL() = %q, want %q", got, want)
	}
}
