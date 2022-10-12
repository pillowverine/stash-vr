package config

import (
	"github.com/rs/zerolog/log"
	"os"
	"strings"
	"sync"
)

const (
	envKeyStashGraphQLUrl  = "STASH_GRAPHQL_URL"
	envKeyStashApiKey      = "STASH_API_KEY"
	envKeyFavoriteTag      = "FAVORITE_TAG"
	envKeyFilters          = "FILTERS"
	envKeyLogLevel         = "LOG_LEVEL"
	envKeyDisableRedact    = "DISABLE_REDACT"
	envKeyForceHTTPS       = "FORCE_HTTPS"
	envKeyDisableHeatmap   = "DISABLE_HEATMAP"
	envKeyAllowSyncMarkers = "ALLOW_SYNC_MARKERS"
)

type Application struct {
	StashGraphQLUrl      string
	StashApiKey          string
	FavoriteTag          string
	Filters              string
	IsSyncMarkersAllowed bool
	LogLevel             string
	IsRedactDisabled     bool
	ForceHTTPS           bool
	IsHeatmapDisabled    bool
}

var cfg Application

var once sync.Once

func Get() Application {
	once.Do(func() {
		cfg = Application{
			StashGraphQLUrl:      getEnvOrDefault(envKeyStashGraphQLUrl, "http://localhost:9999/graphql"),
			StashApiKey:          getEnvOrDefault(envKeyStashApiKey, ""),
			FavoriteTag:          getEnvOrDefault(envKeyFavoriteTag, "FAVORITE"),
			Filters:              getEnvOrDefault(envKeyFilters, ""),
			IsSyncMarkersAllowed: getEnvOrDefault(envKeyAllowSyncMarkers, "false") == "true",
			LogLevel:             strings.ToLower(getEnvOrDefault(envKeyLogLevel, "info")),
			IsRedactDisabled:     getEnvOrDefault(envKeyDisableRedact, "false") == "true",
			ForceHTTPS:           getEnvOrDefault(envKeyForceHTTPS, "false") == "true",
			IsHeatmapDisabled:    getEnvOrDefault(envKeyDisableHeatmap, "false") == "true",
		}
	})
	return cfg
}

func getEnvOrDefault(key string, defaultValue string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return val

}

func findEnvOrDefault(keys []string, defaultValue string) string {
	for i, key := range keys {
		v, ok := os.LookupEnv(key)
		if ok {
			if i > 0 {
				log.Warn().Str("deprecated", key).Str("replace with", keys[0]).Msg("Deprecated env. var. found")
			}
			return v
		}
	}
	return defaultValue
}

func (a Application) Redacted() Application {
	a.StashGraphQLUrl = Redacted(a.StashGraphQLUrl)
	a.StashApiKey = Redacted(a.StashApiKey)
	return a
}
