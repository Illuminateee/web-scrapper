package scraper

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Illuminateee/web-scrapper.git/internal/models"
	"github.com/PuerkitoBio/goquery"
)

// JobScraper interface defines methods for job scraping
type JobScraper interface {
	Name() string
	Scrape(ctx context.Context, filters models.SearchFilters) ([]models.Job, error)
	GetBaseURL() string
}

// ScraperManager manages multiple scrapers and coordinates concurrent scraping
type ScraperManager struct {
	scrapers    []JobScraper
	rateLimiter *RateLimiter
	client      *http.Client
}

// NewScraperManager creates a new scraper manager
func NewScraperManager() *ScraperManager {
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    30 * time.Second,
			DisableCompression: true,
		},
	}

	return &ScraperManager{
		scrapers:    make([]JobScraper, 0),
		rateLimiter: NewRateLimiter(5, time.Second), // 5 requests per second
		client:      client,
	}
}

// AddScraper adds a scraper to the manager
func (sm *ScraperManager) AddScraper(scraper JobScraper) {
	sm.scrapers = append(sm.scrapers, scraper)
}

// ScrapeAll scrapes jobs from all registered scrapers concurrently
func (sm *ScraperManager) ScrapeAll(ctx context.Context, filters models.SearchFilters) []models.ScrapingResult {
	var wg sync.WaitGroup
	results := make([]models.ScrapingResult, len(sm.scrapers))

	for i, scraper := range sm.scrapers {
		wg.Add(1)
		go func(index int, s JobScraper) {
			defer wg.Done()

			// Rate limiting
			sm.rateLimiter.Wait(ctx)

			jobs, err := s.Scrape(ctx, filters)
			results[index] = models.ScrapingResult{
				Jobs:   jobs,
				Source: s.Name(),
				Error:  err,
			}

			if err != nil {
				log.Printf("Error scraping %s: %v", s.Name(), err)
			} else {
				log.Printf("Successfully scraped %d jobs from %s", len(jobs), s.Name())
			}
		}(i, scraper)
	}

	wg.Wait()
	return results
}

// GetAllJobs aggregates jobs from all scraping results
func (sm *ScraperManager) GetAllJobs(results []models.ScrapingResult) []models.Job {
	var allJobs []models.Job

	for _, result := range results {
		if result.Error == nil {
			allJobs = append(allJobs, result.Jobs...)
		}
	}

	return allJobs
}

// BaseScraper provides common functionality for scrapers
type BaseScraper struct {
	name       string
	baseURL    string
	client     *http.Client
	userAgents []string
}

// NewBaseScraper creates a new base scraper
func NewBaseScraper(name, baseURL string, client *http.Client) *BaseScraper {
	return &BaseScraper{
		name:    name,
		baseURL: baseURL,
		client:  client,
		userAgents: []string{
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
			"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		},
	}
}

// Name returns the scraper name
func (bs *BaseScraper) Name() string {
	return bs.name
}

// GetBaseURL returns the base URL
func (bs *BaseScraper) GetBaseURL() string {
	return bs.baseURL
}

// FetchDocument fetches and parses an HTML document from the given URL
func (bs *BaseScraper) FetchDocument(ctx context.Context, url string) (*goquery.Document, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set random user agent
	userAgent := bs.userAgents[time.Now().Unix()%int64(len(bs.userAgents))]
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Connection", "keep-alive")

	resp, err := bs.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch URL %s: status code %d", url, resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	return doc, nil
}

// CleanText cleans and normalizes text content
func (bs *BaseScraper) CleanText(text string) string {
	// Remove extra whitespace and normalize
	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\t", " ")

	// Remove multiple spaces
	for strings.Contains(text, "  ") {
		text = strings.ReplaceAll(text, "  ", " ")
	}

	return text
}

// ExtractSkills extracts skills from job description using common IT skills
func (bs *BaseScraper) ExtractSkills(description string) []string {
	commonSkills := []string{
		// Programming Languages
		"go", "golang", "python", "javascript", "typescript", "java", "c++", "c#", "rust", "php", "ruby", "swift", "kotlin",
		// Frameworks
		"react", "angular", "vue", "node.js", "express", "django", "flask", "spring", "gin", "fiber", "echo",
		// Databases
		"mysql", "postgresql", "mongodb", "redis", "elasticsearch", "cassandra", "sqlite",
		// Cloud & DevOps
		"aws", "azure", "gcp", "docker", "kubernetes", "jenkins", "gitlab", "github", "terraform", "ansible",
		// Other Technologies
		"git", "linux", "unix", "sql", "nosql", "rest", "graphql", "microservices", "api", "json", "xml",
		// Methodologies
		"agile", "scrum", "devops", "ci/cd", "tdd", "bdd",
	}

	description = strings.ToLower(description)
	var foundSkills []string
	skillSet := make(map[string]bool)

	for _, skill := range commonSkills {
		if strings.Contains(description, strings.ToLower(skill)) && !skillSet[skill] {
			foundSkills = append(foundSkills, skill)
			skillSet[skill] = true
		}
	}

	return foundSkills
}

// DetermineExperienceLevel determines experience level from job title and description
func (bs *BaseScraper) DetermineExperienceLevel(title, description string) string {
	titleLower := strings.ToLower(title)
	descLower := strings.ToLower(description)

	// Senior/Lead indicators
	if strings.Contains(titleLower, "senior") || strings.Contains(titleLower, "lead") ||
		strings.Contains(titleLower, "principal") || strings.Contains(titleLower, "architect") ||
		strings.Contains(descLower, "5+ years") || strings.Contains(descLower, "7+ years") {
		return "senior"
	}

	// Entry level indicators
	if strings.Contains(titleLower, "junior") || strings.Contains(titleLower, "entry") ||
		strings.Contains(titleLower, "graduate") || strings.Contains(titleLower, "intern") ||
		strings.Contains(descLower, "0-2 years") || strings.Contains(descLower, "no experience") {
		return "entry"
	}

	// Mid-level indicators
	if strings.Contains(descLower, "2-5 years") || strings.Contains(descLower, "3+ years") {
		return "mid"
	}

	return "mid" // default
}

// CheckDegreeRequirement checks if degree is required
func (bs *BaseScraper) CheckDegreeRequirement(description string) bool {
	descLower := strings.ToLower(description)

	// Strong degree requirement indicators
	degreeRequired := []string{
		"bachelor's degree required",
		"master's degree required",
		"degree required",
		"university degree required",
		"college degree required",
	}

	for _, req := range degreeRequired {
		if strings.Contains(descLower, req) {
			return true
		}
	}

	// Look for flexible requirements
	flexibleReqs := []string{
		"degree preferred",
		"or equivalent experience",
		"degree or experience",
		"education or experience",
	}

	for _, req := range flexibleReqs {
		if strings.Contains(descLower, req) {
			return false
		}
	}

	// Default to degree required if "degree" is mentioned without flexibility
	return strings.Contains(descLower, "degree") || strings.Contains(descLower, "bachelor") || strings.Contains(descLower, "master")
}
