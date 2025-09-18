package scraper

import (
	"fmt"
	"net/http"
	"time"
)

// ScraperConfig represents configuration for a job scraper
type ScraperConfig struct {
	Name         string            `json:"name"`
	Enabled      bool              `json:"enabled"`
	URL          string            `json:"url"`
	Type         string            `json:"type"`       // "public", "authenticated", "api"
	RateLimit    int               `json:"rate_limit"` // requests per minute
	Timeout      time.Duration     `json:"timeout"`
	Headers      map[string]string `json:"headers"`
	RequiresAuth bool              `json:"requires_auth"`
	Credentials  map[string]string `json:"credentials"`
}

// ScraperRegistry manages available scrapers
type ScraperRegistry struct {
	configs    map[string]ScraperConfig
	httpClient *http.Client
}

// NewScraperRegistry creates a new scraper registry
func NewScraperRegistry() *ScraperRegistry {
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    30 * time.Second,
			DisableCompression: true,
		},
	}

	registry := &ScraperRegistry{
		configs:    make(map[string]ScraperConfig),
		httpClient: client,
	}

	registry.loadDefaultConfigs()
	return registry
}

// loadDefaultConfigs loads default scraper configurations
func (sr *ScraperRegistry) loadDefaultConfigs() {
	// RemoteOK configuration
	sr.configs["remoteok"] = ScraperConfig{
		Name:         "RemoteOK",
		Enabled:      true,
		URL:          "https://remoteok.io",
		Type:         "api",
		RateLimit:    30, // 30 requests per minute
		Timeout:      30 * time.Second,
		RequiresAuth: false,
		Headers: map[string]string{
			"Accept": "application/json",
		},
	}

	// WeWorkRemotely configuration
	sr.configs["weworkremotely"] = ScraperConfig{
		Name:         "WeWorkRemotely",
		Enabled:      true,
		URL:          "https://weworkremotely.com",
		Type:         "public",
		RateLimit:    20, // 20 requests per minute
		Timeout:      30 * time.Second,
		RequiresAuth: false,
		Headers: map[string]string{
			"Accept": "text/html",
		},
	}

	// LinkedIn configuration (requires authentication)
	sr.configs["linkedin"] = ScraperConfig{
		Name:         "LinkedIn",
		Enabled:      false, // Disabled by default due to auth requirements
		URL:          "https://linkedin.com",
		Type:         "authenticated",
		RateLimit:    10, // 10 requests per minute
		Timeout:      30 * time.Second,
		RequiresAuth: true,
		Headers: map[string]string{
			"Accept": "text/html",
		},
	}

	// JobStreet configuration (requires authentication)
	sr.configs["jobstreet"] = ScraperConfig{
		Name:         "JobStreet",
		Enabled:      false, // Disabled by default due to auth requirements
		URL:          "https://id.jobstreet.com",
		Type:         "authenticated",
		RateLimit:    15, // 15 requests per minute
		Timeout:      30 * time.Second,
		RequiresAuth: true,
		Headers: map[string]string{
			"Accept": "text/html",
		},
	}

	// Indeed configuration (public RSS feeds only)
	sr.configs["indeed"] = ScraperConfig{
		Name:         "Indeed",
		Enabled:      true,
		URL:          "https://indeed.com",
		Type:         "public",
		RateLimit:    10, // 10 requests per minute
		Timeout:      30 * time.Second,
		RequiresAuth: false,
		Headers: map[string]string{
			"Accept": "application/rss+xml",
		},
	}
}

// CreateScraper creates a scraper instance based on configuration
func (sr *ScraperRegistry) CreateScraper(name string) (JobScraper, error) {
	config, exists := sr.configs[name]
	if !exists {
		return nil, fmt.Errorf("scraper configuration not found: %s", name)
	}

	if !config.Enabled {
		return nil, fmt.Errorf("scraper is disabled: %s", name)
	}

	switch name {
	case "remoteok":
		return NewRemoteOKScraper(sr.httpClient), nil
	case "weworkremotely":
		return NewWeWorkRemotelyScraper(sr.httpClient), nil
	case "linkedin":
		return NewLinkedInScraper(sr.httpClient), nil
	case "jobstreet":
		return NewJobStreetScraper(sr.httpClient, "id"), nil
	default:
		return nil, fmt.Errorf("scraper implementation not found: %s", name)
	}
}

// GetEnabledScrapers returns all enabled scrapers
func (sr *ScraperRegistry) GetEnabledScrapers() []JobScraper {
	var scrapers []JobScraper

	for name, config := range sr.configs {
		if config.Enabled {
			if scraper, err := sr.CreateScraper(name); err == nil {
				scrapers = append(scrapers, scraper)
			}
		}
	}

	return scrapers
}

// EnableScraper enables a scraper
func (sr *ScraperRegistry) EnableScraper(name string) error {
	config, exists := sr.configs[name]
	if !exists {
		return fmt.Errorf("scraper not found: %s", name)
	}

	config.Enabled = true
	sr.configs[name] = config
	return nil
}

// DisableScraper disables a scraper
func (sr *ScraperRegistry) DisableScraper(name string) error {
	config, exists := sr.configs[name]
	if !exists {
		return fmt.Errorf("scraper not found: %s", name)
	}

	config.Enabled = false
	sr.configs[name] = config
	return nil
}

// GetScraperConfig returns configuration for a scraper
func (sr *ScraperRegistry) GetScraperConfig(name string) (ScraperConfig, error) {
	config, exists := sr.configs[name]
	if !exists {
		return ScraperConfig{}, fmt.Errorf("scraper not found: %s", name)
	}
	return config, nil
}

// ListScrapers returns all available scraper configurations
func (sr *ScraperRegistry) ListScrapers() map[string]ScraperConfig {
	return sr.configs
}

// AddCustomScraper adds a custom scraper configuration
func (sr *ScraperRegistry) AddCustomScraper(name string, config ScraperConfig) {
	sr.configs[name] = config
}

// ScraperStatus represents the status of a scraper
type ScraperStatus struct {
	Name         string `json:"name"`
	Enabled      bool   `json:"enabled"`
	Type         string `json:"type"`
	RequiresAuth bool   `json:"requires_auth"`
	LastUsed     string `json:"last_used"`
	Status       string `json:"status"` // "active", "error", "disabled"
}
