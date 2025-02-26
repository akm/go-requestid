package requestid

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akm/slogctx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSlog(t *testing.T) {
	slogwNS := slogctx.NewNamespace()

	assertRequestIDInJSON := func(t *testing.T, data []byte, assertion func(actual string)) {
		t.Helper()
		// requestid field is included in log entry from slog after calling RegisterSlogHandle
		type LogEntry struct {
			RequestID string `json:"req_id"` //nolint:tagliatelle
		}
		var logEntry LogEntry
		err := json.Unmarshal(data, &logEntry)
		if assert.NoError(t, err) {
			assertion(logEntry.RequestID)
		} else {
			t.Logf("data: %s", string(data))
		}
	}

	baseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// [Important] Use slog functions with context.Context in your Handler
		slog.InfoContext(r.Context(), "Hello, world!")
		_, err := w.Write([]byte("Hello, world!"))
		assert.NoError(t, err)
		w.WriteHeader(http.StatusOK)
	})
	mw := New(
		SlogwNamespace(slogwNS),
		RequestHeader("X-Request-ID"),
	)
	ts := httptest.NewServer(mw.Wrap(baseHandler))
	defer ts.Close()

	t.Run("request with X-Request-ID header", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		jsonHandler := slog.NewJSONHandler(buf, nil)
		slog.SetDefault(slogwNS.New(jsonHandler))

		req, err := http.NewRequest(http.MethodGet, ts.URL, nil) //nolint:noctx
		require.NoError(t, err)
		req.Header.Set("X-Request-ID", "in-header")

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "in-header", resp.Header.Get("X-Request-ID"))

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Equal(t, "Hello, world!", string(body))

		// t.Logf("buf: %s", buf.String())

		lines := bytes.Split(buf.Bytes(), []byte("\n"))
		assert.Len(t, lines, 3)
		for i, line := range lines {
			if i != 0 { // skip except the first line
				continue
			}
			assertRequestIDInJSON(t, line, func(actual string) {
				assert.Equal(t, "in-header", actual)
			})
		}
	})

	t.Run("request without X-Request-ID header", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		jsonHandler := slog.NewJSONHandler(buf, nil)
		slog.SetDefault(slogwNS.New(jsonHandler))

		req, err := http.NewRequest(http.MethodGet, ts.URL, nil) //nolint:noctx
		require.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NotEmpty(t, resp.Header.Get("X-Request-ID"))

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Equal(t, "Hello, world!", string(body))

		// t.Logf("buf: %s", buf.String())

		lines := bytes.Split(buf.Bytes(), []byte("\n"))
		assert.Len(t, lines, 3)
		for i, line := range lines {
			if i != 0 { // skip except the first line
				continue
			}
			assertRequestIDInJSON(t, line, func(actual string) {
				assert.NotEmpty(t, actual)
				assert.Len(t, actual, 16)
			})
		}
	})
}
