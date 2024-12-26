package apiserver

import (
	"encoding/json"
	"errors"
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
			Artists:  artistNamesToString(track.Item.Artists),
			URL:      track.Item.ExternalURLs["spotify"],
			Playing:  track.IsPlaying,
			ImageURL: imageURL,
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(&resp); err != nil {
			logger.Error("unable to encode response", zap.Error(err))
		}
	})
}

func getImageURLByDimensions(images []*spotify.Image, w int, h int) (string, error) {
	for _, image := range images {
		if image.Width == w && image.Height == h {
			return image.URL, nil
		}
	}
	return "", errors.New("no images with those dimensions")
}

func artistNamesToString(artists []*spotify.Artist) string {
	var sb strings.Builder
	for idx, artist := range artists {
		sb.WriteString(artist.Name)
		if idx < len(artists)-1 {
			sb.WriteString(", ")
		}
	}
	return sb.String()
}
