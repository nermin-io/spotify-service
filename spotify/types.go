package spotify

type Context struct {
	Type         string            `json:"type"`
	URI          string            `json:"uri"`
	Href         string            `json:"href"`
	ExternalURLs map[string]string `json:"external_urls"`
}

type Track struct {
	ID           string            `json:"id"`
	Type         string            `json:"type"`
	URI          string            `json:"uri"`
	Name         string            `json:"name"`
	Href         string            `json:"href"`
	DiscNumber   int               `json:"disc_number"`
	TrackNumber  int               `json:"track_number"`
	Popularity   int               `json:"popularity"`
	IsLocal      bool              `json:"is_local"`
	IsPlayable   bool              `json:"is_playable"`
	Explicit     bool              `json:"explicit"`
	Duration     int               `json:"duration_ms"`
	Album        *Album            `json:"album"`
	Artists      []*Artist         `json:"artists"`
	ExternalURLs map[string]string `json:"external_urls"`
	ExternalIDs  map[string]string `json:"external_ids"`
}

type Album struct {
	ID                   string            `json:"id"`
	Name                 string            `json:"name"`
	Href                 string            `json:"href"`
	AlbumType            string            `json:"album_type"`
	ExternalURLs         map[string]string `json:"external_urls"`
	Images               []*Image          `json:"images"`
	Artists              []*Artist         `json:"artists"`
	IsPlayable           bool              `json:"is_playable"`
	ReleaseDate          string            `json:"release_date"`
	ReleaseDatePrecision string            `json:"release_date_precision"`
	TotalTracks          int               `json:"total_tracks"`
	Type                 string            `json:"type"`
	URI                  string            `json:"uri"`
}

type Image struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type Artist struct {
	ID           string            `json:"id"`
	Href         string            `json:"href"`
	Name         string            `json:"name"`
	Type         string            `json:"type"`
	URI          string            `json:"uri"`
	ExternalURLs map[string]string `json:"external_urls"`
}
