package api

import (
	"encoding/json"
	"net/http"

	"github.com/Illuminateee/web-scrapper.git/internal/scraper"
	"github.com/gorilla/mux"
)

// ScraperHandler handles scraper configuration requests
type ScraperHandler struct {
	registry *scraper.ScraperRegistry
}

// NewScraperHandler creates a new scraper handler
func NewScraperHandler(registry *scraper.ScraperRegistry) *ScraperHandler {
	return &ScraperHandler{
		registry: registry,
	}
}

// ListScrapers returns all available scrapers and their status
func (h *ScraperHandler) ListScrapers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	configs := h.registry.ListScrapers()

	// Convert to status format for frontend
	statuses := make([]scraper.ScraperStatus, 0, len(configs))
	for _, config := range configs {
		status := scraper.ScraperStatus{
			Name:         config.Name,
			Enabled:      config.Enabled,
			Type:         config.Type,
			RequiresAuth: config.RequiresAuth,
			Status:       getScraperStatus(config),
		}
		statuses = append(statuses, status)
	}

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    statuses,
	}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// EnableScraper enables a specific scraper
func (h *ScraperHandler) EnableScraper(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	scraperName := vars["name"]

	if err := h.registry.EnableScraper(scraperName); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Scraper enabled successfully",
	})
}

// DisableScraper disables a specific scraper
func (h *ScraperHandler) DisableScraper(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	scraperName := vars["name"]

	if err := h.registry.DisableScraper(scraperName); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Scraper disabled successfully",
	})
}

// GetScraperConfig returns configuration for a specific scraper
func (h *ScraperHandler) GetScraperConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	scraperName := vars["name"]

	config, err := h.registry.GetScraperConfig(scraperName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    config,
	})
}

// getScraperStatus determines the status of a scraper based on its configuration
func getScraperStatus(config scraper.ScraperConfig) string {
	if !config.Enabled {
		return "disabled"
	}

	if config.RequiresAuth {
		return "requires_auth"
	}

	return "active"
}
