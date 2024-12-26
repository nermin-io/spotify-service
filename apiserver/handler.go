package apiserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nermin-io/spotify-service/apiserver/middleware"
	"github.com/nermin-io/spotify-service/spotify"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

func NewHandler(logger *zap.Logger, spotifyClient *spotify.Client) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /currently-playing", handleGetCurrentTrack(logger, spotifyClient))

	return middleware.Apply(mux, middleware.NewLoggingMiddleware(logger))
}

type currentlyPlayingResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Artists  string `json:"artists"`
	URL      string `json:"url"`
	ImageURL string `json:"image_url"`
	Playing  bool   `json:"playing"`
}

func handleGetCurrentTrack(logger *zap.Logger, spotifyClient *spotify.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		track, err := spotifyClient.CurrentlyPlayingTrack(r.Context())
		if err != nil {
			logger.Warn("unable to get current playing track", zap.Error(err))
			w.WriteHeader(http.StatusNoContent)
			return
		}
		imageURL, err := getImageURLByDimensions(track.Item.Album.Images, 300, 300)
		if err != nil {
			logger.Warn("unable to get image URL", zap.Error(err))
		}
		resp := currentlyPlayingResponse{
			ID:       track.Item.ID,
			Name:     track.Item.Name,
			Artists:  getArtistNamesAsString(track.Item.Artists),
			URL:      track.Item.ExternalURLs["spotify"],
			Playing:  track.IsPlaying,
			ImageURL: imageURL,
		}
		if err := encode(w, r, http.StatusOK, &resp); err != nil {
			logger.Warn("unable to encode response", zap.Error(err))
		}
	})
}

func encode[T any](w http.ResponseWriter, _ *http.Request, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(&v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

func decode[T any](r *http.Request) (*T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return nil, fmt.Errorf("decode json: %w", err)
	}
	return &v, nil
}

func getImageURLByDimensions(images []*spotify.Image, w int, h int) (string, error) {
	for _, image := range images {
		if image.Width == w && image.Height == h {
			return image.URL, nil
		}
	}
	return "", errors.New("no images with those dimensions")
}

func getArtistNamesAsString(artists []*spotify.Artist) string {
	var sb strings.Builder
	for idx, artist := range artists {
		sb.WriteString(artist.Name)
		if idx < len(artists)-1 {
			sb.WriteString(", ")
		}
	}
	return sb.String()
}
