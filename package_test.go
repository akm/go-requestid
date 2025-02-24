package requestid

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefault(t *testing.T) {
	assert.Equal(t, defaultNamespace, Default())
	assert.Equal(t, defaultNamespace, Default(), "DefaultNamespace should return the same instance")

	ns1 := New(ResponseHeader("X-REQ-ID"))
	SetDefault(ns1)
	defer SetDefault(defaultNamespace)

	assert.Equal(t, ns1, Default())
	assert.Equal(t, ns1, Default(), "DefaultNamespace should return the same instance")
}

var testOptions = &slog.HandlerOptions{
	ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			return slog.Attr{}
		}
		return a
	},
}

func TestDefaultSimpleWrap(t *testing.T) {
	backupDefaultNamespace := defaultNamespace
	defer func() { defaultNamespace = backupDefaultNamespace }()

	buf := bytes.NewBuffer(nil)
	jsonHandler := slog.NewJSONHandler(buf, testOptions)
	logger := NewLogger(jsonHandler)

	baseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "working")
		_, err := w.Write([]byte("Hello, world!"))
		assert.NoError(t, err)
		w.WriteHeader(http.StatusOK)
	})

	ts := httptest.NewServer(Wrap(baseHandler))
	defer ts.Close()

	t.Run("request", func(t *testing.T) {
		req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, ts.URL, nil)
		require.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Equal(t, "Hello, world!", string(body))

		rec := map[string]interface{}{}
		err = json.Unmarshal(buf.Bytes(), &rec)
		require.NoError(t, err)

		assert.Equal(t, "INFO", rec["level"])
		assert.Equal(t, "working", rec["msg"])
	})
}

func TestWrapSlogHandler(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	handler := slog.NewJSONHandler(buf, testOptions)
	logger := slog.New(WrapSlogHandler(handler))

	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		logger.InfoContext(ctx, "Start")
		defer logger.InfoContext(ctx, "End")
		io.WriteString(w, "Hello, world!\n") //nolint:errcheck
	}

	ts := httptest.NewServer(
		Wrap(http.HandlerFunc(helloHandler)),
	)
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	require.NoError(t, err)
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	lines := bytes.Split(buf.Bytes(), []byte("\n"))
	if len(lines) < 2 {
		t.Fatalf("unexpected response: %s", b)
	}
	for i, line := range lines[:len(lines)-1] {
		if len(line) == 0 {
			continue
		}
		t.Logf("line: %s", line)
		d := map[string]interface{}{}
		err := json.Unmarshal(line, &d)
		require.NoError(t, err)
		require.Contains(t, d, "req_id")
		require.NotEmpty(t, d["req_id"])
		require.Contains(t, d, "level")
		require.Equal(t, "INFO", d["level"])
		require.Contains(t, d, "msg")
		switch i {
		case 0:
			require.Equal(t, "Start", d["msg"])
		case 1:
			require.Equal(t, "End", d["msg"])
		default:
			t.Fatalf("unexpected line: %s", line)
		}
	}
}
