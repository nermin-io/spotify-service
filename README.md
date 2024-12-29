# Spotify Service

A simple HTTP service that exposes a single endpoint `GET /currently-playing`, which returns information about the song I'm currently listening to on Spotify.

The service is written entirely in Go and is deployed as a container on Cloud Run. You can reach the service via the following URL: [https://spotify.nsehic.com/currently-playing](https://spotify.nsehic.com/currently-playing)
 > **NOTE**: If the service returns the status code 204 No Content, it means I'm currently not listening to any song.

## Motivation
To show what music I'm listening to on my website and other personal sites.

## Build
You can build the project using Docker or the Go Build tool.

### Docker
The `Dockerfile` is located in the root directory. To build an image, run the following command:
```bash
docker build -t spotifyservice .
```
You can then start a container by running the following command (although you'll need to pass the environment variables as described by the configuration section):
```bash
docker run --name spotifyservice-http -p 8080:8080 spotifyservice
```

### Go Build Tool
To build the project using the Go build tool, please make sure you have [Go installed](https://go.dev/dl/).

Once installed, run the following command:
```bash
go build -o <PATH-TO-BINARY> ./cmd/spotifyservice
```

## Configuration

### Flags
- `-debug` Enables debug logging

### Environment Variables
Please note, if you would like to use this service for yourself, you'll need to create your own Spotify app to get access credentials. For more information, see the [Web API](https://developer.spotify.com/documentation/web-api).
- `SPOTIFY_BASE_URL` - The base URL when making requests to Spotify, usually https://api.spotify.com
- `SPOTIFY_CREDENTIALS_URL` - The URL to make requests to get an access token, usually https://accounts.spotify.com/api/token
- `SPOTIFY_CLIENT_ID` - OAuth 2.0 Client ID
- `SPOTIFY_CLIENT_SECRET` - OAuth 2.0 Client Secret
- `SPOTIFY_REFRESH_TOKEN` - The refresh token that is returned after completing the authorization flow.

## Logging
The service outputs structured logs to `stdout` in JSON format. If you plan to also deploy this to Cloud Run, you can optional pass in the `GCP_PROJECT_ID` environment variable, which will enable logging of the automatically generated trace.