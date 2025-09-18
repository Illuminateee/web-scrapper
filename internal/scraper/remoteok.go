package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Illuminateee/web-scrapper.git/internal/models"
)

// RemoteOKScraper scrapes jobs from RemoteOK.io
type RemoteOKScraper struct {
	*BaseScraper
}

// NewRemoteOKScraper creates a new RemoteOK scraper
func NewRemoteOKScraper(client *http.Client) *RemoteOKScraper {
	return &RemoteOKScraper{
		BaseScraper: NewBaseScraper("RemoteOK", "https://remoteok.io", client),
	}
}

// RemoteOKJob represents a job from RemoteOK API
type RemoteOKJob struct {
	ID          string   `json:"id"`
	URL         string   `json:"url"`
	Position    string   `json:"position"`
	Company     string   `json:"company"`
	CompanyLogo string   `json:"company_logo"`
	Location    string   `json:"location"`
	Tags        []string `json:"tags"`
	Date        string   `json:"date"`
	Description string   `json:"description"`
	Apply       string   `json:"apply"`
	Salary      string   `json:"salary"`
}

// Scrape implements the JobScraper interface
func (r *RemoteOKScraper) Scrape(ctx context.Context, filters models.SearchFilters) ([]models.Job, error) {
	// RemoteOK has a public API
	apiURL := "https://remoteok.io/api"

	// Add query parameters if available
	if filters.JobTitle != "" {
		// RemoteOK doesn't support direct filtering, we'll filter results
	}

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers to mimic a browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch RemoteOK API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("RemoteOK API returned status: %d", resp.StatusCode)
	}

	var remoteJobs []RemoteOKJob
	if err := json.NewDecoder(resp.Body).Decode(&remoteJobs); err != nil {
		return nil, fmt.Errorf("failed to decode RemoteOK response: %w", err)
	}

	var jobs []models.Job
	for _, rJob := range remoteJobs {
		// Skip the first item which is usually metadata
		if rJob.ID == "" {
			continue
		}

		job := r.convertRemoteOKJob(rJob)

		// Apply basic filtering
		if r.matchesFilters(job, filters) {
			jobs = append(jobs, job)
		}
	}

	return jobs, nil
}

func (r *RemoteOKScraper) convertRemoteOKJob(rJob RemoteOKJob) models.Job {
	// Parse salary from string
	var salaryMin, salaryMax int
	if rJob.Salary != "" {
		salaryMin, salaryMax = r.parseSalary(rJob.Salary)
	}

	// Determine experience level from tags and position
	expLevel := r.DetermineExperienceLevel(rJob.Position, strings.Join(rJob.Tags, " "))

	// Extract skills from tags
	skills := r.extractSkillsFromTags(rJob.Tags)

	// Parse date
	postedDate := time.Now() // Default to now if parsing fails
	if rJob.Date != "" {
		if parsed, err := time.Parse("2006-01-02T15:04:05Z", rJob.Date); err == nil {
			postedDate = parsed
		}
	}

	return models.Job{
		ID:              fmt.Sprintf("remoteok-%s", rJob.ID),
		Title:           rJob.Position,
		Company:         rJob.Company,
		Location:        "Remote", // RemoteOK is all remote jobs
		Description:     rJob.Description,
		Skills:          skills,
		SalaryMin:       salaryMin,
		SalaryMax:       salaryMax,
		SalaryCurrency:  "USD",
		DegreeRequired:  r.CheckDegreeRequirement(rJob.Description),
		ExperienceLevel: expLevel,
		RemoteOption:    "remote",
		PostedDate:      postedDate,
		URL:             fmt.Sprintf("https://remoteok.io/remote-jobs/%s", rJob.ID),
		Source:          r.Name(),
		Industry:        "Technology",
	}
}

func (r *RemoteOKScraper) parseSalary(salaryStr string) (int, int) {
	// Remove common currency symbols and formatting
	cleaned := strings.ReplaceAll(salaryStr, "$", "")
	cleaned = strings.ReplaceAll(cleaned, ",", "")
	cleaned = strings.ReplaceAll(cleaned, "k", "000")
	cleaned = strings.ReplaceAll(cleaned, "K", "000")

	// Look for patterns like "80000-120000" or "80k-120k"
	if strings.Contains(cleaned, "-") {
		parts := strings.Split(cleaned, "-")
		if len(parts) == 2 {
			if min, err := strconv.Atoi(strings.TrimSpace(parts[0])); err == nil {
				if max, err := strconv.Atoi(strings.TrimSpace(parts[1])); err == nil {
					return min, max
				}
			}
		}
	}

	// Single salary value
	if val, err := strconv.Atoi(strings.TrimSpace(cleaned)); err == nil {
		return val, val
	}

	return 0, 0
}

func (r *RemoteOKScraper) extractSkillsFromTags(tags []string) []string {
	var skills []string
	skillSet := make(map[string]bool)

	for _, tag := range tags {
		tagLower := strings.ToLower(tag)
		// Add relevant technical tags as skills
		if r.isTechnicalSkill(tagLower) && !skillSet[tagLower] {
			skills = append(skills, tag)
			skillSet[tagLower] = true
		}
	}

	return skills
}

func (r *RemoteOKScraper) isTechnicalSkill(tag string) bool {
	technicalSkills := []string{
		"javascript", "python", "java", "go", "golang", "rust", "php", "ruby", "swift", "kotlin",
		"react", "angular", "vue", "node", "django", "flask", "spring", "express",
		"mysql", "postgresql", "mongodb", "redis", "elasticsearch",
		"aws", "azure", "gcp", "docker", "kubernetes", "jenkins",
		"git", "linux", "sql", "nosql", "api", "rest", "graphql",
	}

	for _, skill := range technicalSkills {
		if strings.Contains(tag, skill) {
			return true
		}
	}
	return false
}

func (r *RemoteOKScraper) matchesFilters(job models.Job, filters models.SearchFilters) bool {
	// Job title filter
	if filters.JobTitle != "" {
		if !strings.Contains(strings.ToLower(job.Title), strings.ToLower(filters.JobTitle)) {
			return false
		}
	}

	// Keywords filter
	if len(filters.Keywords) > 0 {
		jobText := strings.ToLower(job.Title + " " + job.Description + " " + strings.Join(job.Skills, " "))
		found := false
		for _, keyword := range filters.Keywords {
			if strings.Contains(jobText, strings.ToLower(keyword)) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Salary filters
	if filters.MinSalary > 0 && (job.SalaryMin == 0 || job.SalaryMin < filters.MinSalary) {
		return false
	}
	if filters.MaxSalary > 0 && (job.SalaryMax == 0 || job.SalaryMax > filters.MaxSalary) {
		return false
	}

	// Experience level filter
	if filters.ExperienceLevel != "" {
		if !strings.EqualFold(job.ExperienceLevel, filters.ExperienceLevel) {
			return false
		}
	}

	// Degree requirement filter
	if filters.DegreeRequired != nil {
		if job.DegreeRequired != *filters.DegreeRequired {
			return false
		}
	}

	// Skills filter
	if len(filters.Skills) > 0 {
		jobSkills := make(map[string]bool)
		for _, skill := range job.Skills {
			jobSkills[strings.ToLower(skill)] = true
		}

		for _, requiredSkill := range filters.Skills {
			if !jobSkills[strings.ToLower(requiredSkill)] {
				return false
			}
		}
	}

	return true
}
