package spotify

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/nermin-io/spotify-service/logging"
	"go.uber.org/zap"
)

type Client struct {
	httpClient     *http.Client
	baseURL        string
	credentialsURL string
	clientID       string
	clientSecret   string
	accessToken    string
	expiresAt      time.Time
	refreshToken   string
}

func NewClient() *Client {
	return &Client{
		httpClient:     http.DefaultClient,
		baseURL:        os.Getenv("SPOTIFY_BASE_URL"),
		credentialsURL: os.Getenv("SPOTIFY_CREDENTIALS_URL"),
		clientID:       os.Getenv("SPOTIFY_CLIENT_ID"),
		clientSecret:   os.Getenv("SPOTIFY_CLIENT_SECRET"),
		refreshToken:   os.Getenv("SPOTIFY_REFRESH_TOKEN"),
	}
}

func (sc *Client) ensureValidAccessToken(ctx context.Context) error {
	isExpired := time.Now().UTC().After(sc.expiresAt)
	if sc.accessToken == "" || isExpired {
		return sc.refreshAccessToken(ctx)
	}

	return nil
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

func (sc *Client) refreshAccessToken(ctx context.Context) error {
	logger := logging.FromContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", sc.refreshToken)

	req, err := http.NewRequestWithContext(ctx, "POST", sc.credentialsURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}
	basicToken := base64.StdEncoding.EncodeToString(fmt.Appendf(nil, "%s:%s", sc.clientID, sc.clientSecret))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", basicToken))

	resp, err := sc.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("could not send request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logger.Error("could not close response body", zap.Error(err))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to refresh token, status: %d", resp.StatusCode)
	}

	var token tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return fmt.Errorf("could not decode response: %w", err)
	}

	sc.accessToken = token.AccessToken
	sc.expiresAt = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)

	return nil
}
