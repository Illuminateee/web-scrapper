package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Illuminateee/web-scrapper.git/internal/models"
	"github.com/Illuminateee/web-scrapper.git/internal/scraper"
	"github.com/Illuminateee/web-scrapper.git/internal/storage"
	"github.com/gorilla/mux"
)

// JobHandler handles job-related API requests
type JobHandler struct {
	scraperManager *scraper.ScraperManager
	storage        storage.JobStorage
}

// NewJobHandler creates a new job handler
func NewJobHandler() *JobHandler {
	// Initialize scraper manager
	scraperManager := scraper.NewScraperManager()

	// Create HTTP client for scrapers
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    30 * time.Second,
			DisableCompression: true,
		},
	}

	// Add real job site scrapers
	scraperManager.AddScraper(scraper.NewRemoteOKScraper(client))
	scraperManager.AddScraper(scraper.NewWeWorkRemotelyScraper(client))

	// Add authenticated scrapers (demo mode)
	linkedinScraper := scraper.NewLinkedInScraper(client)
	jobstreetScraper := scraper.NewJobStreetScraper(client, "id") // Indonesia
	scraperManager.AddScraper(linkedinScraper)
	scraperManager.AddScraper(jobstreetScraper)

	// Keep one mock scraper for additional data
	scraperManager.AddScraper(scraper.NewMockJobScraper("MockJobSite"))

	// Initialize storage
	storage := storage.NewInMemoryStorage()

	return &JobHandler{
		scraperManager: scraperManager,
		storage:        storage,
	}
}

// SetupRoutes sets up the API routes
func SetupRoutes(router *mux.Router) {
	handler := NewJobHandler()

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// Job search endpoint
	api.HandleFunc("/jobs/search", handler.SearchJobs).Methods("GET", "OPTIONS")

	// Advanced search with custom job sites
	api.HandleFunc("/jobs/search/advanced", handler.AdvancedSearch).Methods("POST", "OPTIONS")

	// Get job by ID
	api.HandleFunc("/jobs/{id}", handler.GetJob).Methods("GET", "OPTIONS")

	// Health check
	api.HandleFunc("/health", handler.HealthCheck).Methods("GET")

	// Analytics endpoint
	api.HandleFunc("/analytics", handler.GetAnalytics).Methods("GET", "OPTIONS")

	// Clear cache endpoint
	api.HandleFunc("/cache/clear", handler.ClearCache).Methods("POST", "OPTIONS")

	// Scraper management endpoints
	registry := scraper.NewScraperRegistry()
	scraperHandler := NewScraperHandler(registry)

	api.HandleFunc("/scrapers", scraperHandler.ListScrapers).Methods("GET", "OPTIONS")
	api.HandleFunc("/scrapers/{name}", scraperHandler.GetScraperConfig).Methods("GET", "OPTIONS")
	api.HandleFunc("/scrapers/{name}/enable", scraperHandler.EnableScraper).Methods("POST", "OPTIONS")
	api.HandleFunc("/scrapers/{name}/disable", scraperHandler.DisableScraper).Methods("POST", "OPTIONS")
}

// SearchJobs handles job search requests
func (h *JobHandler) SearchJobs(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	filters := h.parseSearchFilters(r)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Printf("Searching jobs with filters: %+v", filters)

	// Scrape jobs from all sources
	results := h.scraperManager.ScrapeAll(ctx, filters)
	allJobs := h.scraperManager.GetAllJobs(results)

	// Store jobs in cache
	if err := h.storage.Store(allJobs); err != nil {
		log.Printf("Error storing jobs: %v", err)
	}

	// Search stored jobs
	response, err := h.storage.Search(filters)
	if err != nil {
		http.Error(w, "Error searching jobs", http.StatusInternalServerError)
		return
	}

	// Add scraping results info
	response.Analytics.TotalJobs = len(allJobs)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// AdvancedSearch handles advanced job search with custom job sites
func (h *JobHandler) AdvancedSearch(w http.ResponseWriter, r *http.Request) {
	var searchRequest models.SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&searchRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Printf("Advanced search with filters: %+v", searchRequest.Filters)
	log.Printf("Custom job sites: %+v", searchRequest.JobSites)

	// For now, we'll use the existing scrapers but could be extended to use custom sites
	results := h.scraperManager.ScrapeAll(ctx, searchRequest.Filters)
	allJobs := h.scraperManager.GetAllJobs(results)

	// Store jobs in cache
	if err := h.storage.Store(allJobs); err != nil {
		log.Printf("Error storing jobs: %v", err)
	}

	// Search stored jobs
	response, err := h.storage.Search(searchRequest.Filters)
	if err != nil {
		http.Error(w, "Error searching jobs", http.StatusInternalServerError)
		return
	}

	// Add scraping results info
	response.Analytics.TotalJobs = len(allJobs)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetJob handles getting a specific job by ID
func (h *JobHandler) GetJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID := vars["id"]

	// Search for the job by ID
	filters := models.SearchFilters{Limit: 1000} // Get all jobs to find by ID
	response, err := h.storage.Search(filters)
	if err != nil {
		http.Error(w, "Error searching jobs", http.StatusInternalServerError)
		return
	}

	// Find job by ID
	for _, job := range response.Jobs {
		if job.ID == jobID {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(job)
			return
		}
	}

	http.Error(w, "Job not found", http.StatusNotFound)
}

// GetAnalytics handles analytics requests
func (h *JobHandler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	filters := models.SearchFilters{Limit: 1000} // Get all jobs for analytics
	response, err := h.storage.Search(filters)
	if err != nil {
		http.Error(w, "Error getting analytics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response.Analytics)
}

// ClearCache handles cache clearing requests
func (h *JobHandler) ClearCache(w http.ResponseWriter, r *http.Request) {
	if err := h.storage.Clear(); err != nil {
		http.Error(w, "Error clearing cache", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Cache cleared successfully"})
}

// HealthCheck handles health check requests
func (h *JobHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// parseSearchFilters parses search filters from query parameters
func (h *JobHandler) parseSearchFilters(r *http.Request) models.SearchFilters {
	filters := models.SearchFilters{
		Limit: 50, // Default limit
	}

	// Job title
	if title := r.URL.Query().Get("title"); title != "" {
		filters.JobTitle = title
	}

	// Keywords
	if keywords := r.URL.Query().Get("keywords"); keywords != "" {
		filters.Keywords = strings.Split(keywords, ",")
		// Trim spaces
		for i, keyword := range filters.Keywords {
			filters.Keywords[i] = strings.TrimSpace(keyword)
		}
	}

	// Location
	if location := r.URL.Query().Get("location"); location != "" {
		filters.Location = location
	}

	// Multiple locations
	if locations := r.URL.Query().Get("locations"); locations != "" {
		filters.Locations = strings.Split(locations, ",")
		// Trim spaces
		for i, loc := range filters.Locations {
			filters.Locations[i] = strings.TrimSpace(loc)
		}
	}

	// Job category
	if category := r.URL.Query().Get("job_category"); category != "" {
		filters.JobCategory = category
	}

	// Custom job sites
	if jobSites := r.URL.Query().Get("job_sites"); jobSites != "" {
		filters.JobSites = strings.Split(jobSites, ",")
		// Trim spaces
		for i, site := range filters.JobSites {
			filters.JobSites[i] = strings.TrimSpace(site)
		}
	}

	// Remote only
	if remoteOnly := r.URL.Query().Get("remote_only"); remoteOnly == "true" {
		filters.RemoteOnly = true
	}

	// Salary filters
	if minSalary := r.URL.Query().Get("min_salary"); minSalary != "" {
		if val, err := strconv.Atoi(minSalary); err == nil {
			filters.MinSalary = val
		}
	}

	if maxSalary := r.URL.Query().Get("max_salary"); maxSalary != "" {
		if val, err := strconv.Atoi(maxSalary); err == nil {
			filters.MaxSalary = val
		}
	}

	// Experience level
	if expLevel := r.URL.Query().Get("experience_level"); expLevel != "" {
		filters.ExperienceLevel = expLevel
	}

	// Degree required
	if degreeReq := r.URL.Query().Get("degree_required"); degreeReq != "" {
		if degreeReq == "true" {
			val := true
			filters.DegreeRequired = &val
		} else if degreeReq == "false" {
			val := false
			filters.DegreeRequired = &val
		}
	}

	// Skills
	if skills := r.URL.Query().Get("skills"); skills != "" {
		filters.Skills = strings.Split(skills, ",")
		// Trim spaces
		for i, skill := range filters.Skills {
			filters.Skills[i] = strings.TrimSpace(skill)
		}
	}

	// Company size
	if companySize := r.URL.Query().Get("company_size"); companySize != "" {
		filters.CompanySize = companySize
	}

	// Industry
	if industry := r.URL.Query().Get("industry"); industry != "" {
		filters.Industry = industry
	}

	// Pagination
	if limit := r.URL.Query().Get("limit"); limit != "" {
		if val, err := strconv.Atoi(limit); err == nil && val > 0 {
			filters.Limit = val
		}
	}

	if offset := r.URL.Query().Get("offset"); offset != "" {
		if val, err := strconv.Atoi(offset); err == nil && val >= 0 {
			filters.Offset = val
		}
	}

	return filters
}
