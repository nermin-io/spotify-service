package spotify

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
)

type CurrentlyPlaying struct {
	Timestamp            int                    `json:"timestamp"`
	Progress             int                    `json:"progress_ms"`
	Context              *Context               `json:"context"`
	CurrentlyPlayingType string                 `json:"currently_playing_type"`
	IsPlaying            bool                   `json:"is_playing"`
	Actions              map[string]interface{} `json:"actions"`
	Item                 *Track                 `json:"item"`
}

func (sc *Client) CurrentlyPlayingTrack(ctx context.Context) (*CurrentlyPlaying, error) {
	if err := sc.ensureValidAccessToken(ctx); err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}
	url := fmt.Sprintf("%s/v1/me/player/currently-playing?market=AU", sc.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", sc.accessToken))
	resp, err := sc.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("spotify request failed: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			sc.logger.Error("could not close response body", zap.Error(err))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unable to get currently playing track: %s", resp.Status)
	}
	var playing CurrentlyPlaying
	if err := json.NewDecoder(resp.Body).Decode(&playing); err != nil {
		return nil, fmt.Errorf("could not decode response: %w", err)
	}

	if playing.Item == nil {
		return nil, fmt.Errorf("no track information found: %s", playing.CurrentlyPlayingType)
	}

	return &playing, nil
}
