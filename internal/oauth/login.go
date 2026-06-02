package oauth

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"

	"golang.org/x/term"
)

const (
	// ClientID is the OAuth client identifier registered with CircleCI.
	ClientID = "chunk-cli"
	// CallbackTimeout is how long to wait for the browser callback before giving up.
	CallbackTimeout = 5 * time.Minute
	callbackPath    = "/callback"
)

type LoginConfig struct {
	BaseURL   string
	NoBrowser bool
	Signup    bool
}

func Login(ctx context.Context, cfg LoginConfig, status io.Writer) (string, error) {
	deviceID, err := LoadOrCreateDeviceID()
	if err != nil {
		return "", fmt.Errorf("device id: %w", err)
	}

	verifier, err := GenerateVerifier()
	if err != nil {
		return "", fmt.Errorf("pkce verifier: %w", err)
	}
	challenge := S256Challenge(verifier)

	state, err := GenerateState()
	if err != nil {
		return "", fmt.Errorf("state: %w", err)
	}

	port, resultCh, cleanup, err := ListenForCallback(ctx)
	if err != nil {
		return "", err
	}
	defer cleanup()

	redirectURI := fmt.Sprintf("http://127.0.0.1:%d%s", port, callbackPath)
	authorizeURL := buildAuthorizeURL(cfg.BaseURL, redirectURI, challenge, state, deviceID, cfg.Signup)

	w := func(s string) { _, _ = fmt.Fprintln(status, s) }

	if cfg.NoBrowser {
		w("Open this URL in your browser to log in:")
		w("")
		w("  " + authorizeURL)
		w("")
	} else {
		w(authorizeURL)
		if waitForEnter(ctx, status) {
			if err := OpenBrowser(authorizeURL); err != nil {
				w("Could not open browser. Open this URL manually:")
				w("")
				w("  " + authorizeURL)
				w("")
			}
		}
	}
	w("Waiting for login (up to 5 minutes)...")

	var res CallbackResult
	select {
	case res = <-resultCh:
	case <-time.After(CallbackTimeout):
		return "", fmt.Errorf("timed out waiting for browser callback after %s", CallbackTimeout)
	case <-ctx.Done():
		return "", ctx.Err()
	}

	if res.Error != "" {
		return "", fmt.Errorf("authorization denied: %s", res.Error)
	}
	if res.State != state {
		return "", fmt.Errorf("state mismatch (possible CSRF)")
	}
	if res.Code == "" {
		return "", fmt.Errorf("callback contained no authorization code")
	}

	tok, err := exchangeToken(ctx, cfg.BaseURL, res.Code, redirectURI, verifier)
	if err != nil {
		return "", err
	}
	return tok.AccessToken, nil
}

// waitForEnter prompts the user to press Enter before opening the browser.
// Returns true if the browser should be opened, false if cancelled or non-interactive.
func waitForEnter(ctx context.Context, status io.Writer) bool {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return true
	}
	_, _ = fmt.Fprint(status, "Press Enter to open CircleCI in your browser...")
	errCh := make(chan error, 1)
	go func() {
		_, err := bufio.NewReader(os.Stdin).ReadString('\n')
		errCh <- err
	}()
	select {
	case err := <-errCh:
		return err == nil || errors.Is(err, io.EOF)
	case <-ctx.Done():
		return false
	}
}

func buildAuthorizeURL(baseURL, redirectURI, challenge, state, deviceID string, signup bool) string {
	params := url.Values{
		"client_id":             {ClientID},
		"response_type":         {"code"},
		"redirect_uri":          {redirectURI},
		"code_challenge":        {challenge},
		"code_challenge_method": {"S256"},
		"state":                 {state},
		"os":                    {runtime.GOOS},
		"device_id":             {deviceID},
	}
	if signup {
		params.Set("signup", "true")
	}
	return strings.TrimRight(baseURL, "/") + "/oauth/authorize?" + params.Encode()
}
