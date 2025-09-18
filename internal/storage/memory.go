package storage

import (
	"sort"
	"strings"
	"sync"

	"github.com/Illuminateee/web-scrapper.git/internal/models"
)

// JobStorage interface defines methods for job storage
type JobStorage interface {
	Store(jobs []models.Job) error
	Search(filters models.SearchFilters) (*models.SearchResponse, error)
	Clear() error
	GetAnalytics(jobs []models.Job) models.JobAnalytics
}

// InMemoryStorage implements JobStorage using in-memory storage
type InMemoryStorage struct {
	jobs []models.Job
	mu   sync.RWMutex
}

// NewInMemoryStorage creates a new in-memory storage instance
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		jobs: make([]models.Job, 0),
	}
}

// Store saves jobs to memory
func (s *InMemoryStorage) Store(jobs []models.Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Add new jobs, avoiding duplicates by URL
	existingURLs := make(map[string]bool)
	for _, job := range s.jobs {
		existingURLs[job.URL] = true
	}

	for _, job := range jobs {
		if !existingURLs[job.URL] {
			s.jobs = append(s.jobs, job)
			existingURLs[job.URL] = true
		}
	}

	return nil
}

// Search filters and returns jobs based on criteria
func (s *InMemoryStorage) Search(filters models.SearchFilters) (*models.SearchResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var filteredJobs []models.Job

	for _, job := range s.jobs {
		if s.matchesFilters(job, filters) {
			filteredJobs = append(filteredJobs, job)
		}
	}

	// Sort by posted date (newest first)
	sort.Slice(filteredJobs, func(i, j int) bool {
		return filteredJobs[i].PostedDate.After(filteredJobs[j].PostedDate)
	})

	total := len(filteredJobs)

	// Apply pagination
	start := filters.Offset
	end := start + filters.Limit
	if filters.Limit == 0 {
		end = total
	}
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	paginatedJobs := filteredJobs[start:end]
	analytics := s.GetAnalytics(filteredJobs)

	return &models.SearchResponse{
		Jobs:      paginatedJobs,
		Total:     total,
		Analytics: analytics,
		Filters:   filters,
	}, nil
}

// Clear removes all jobs from storage
func (s *InMemoryStorage) Clear() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobs = make([]models.Job, 0)
	return nil
}

// matchesFilters checks if a job matches the search filters
func (s *InMemoryStorage) matchesFilters(job models.Job, filters models.SearchFilters) bool {
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

	// Location filter (single location)
	if filters.Location != "" {
		if !strings.Contains(strings.ToLower(job.Location), strings.ToLower(filters.Location)) {
			return false
		}
	}

	// Multiple locations filter
	if len(filters.Locations) > 0 {
		locationMatch := false
		for _, filterLocation := range filters.Locations {
			if strings.Contains(strings.ToLower(job.Location), strings.ToLower(filterLocation)) {
				locationMatch = true
				break
			}
		}
		if !locationMatch {
			return false
		}
	}

	// Remote filter
	if filters.RemoteOnly && !strings.Contains(strings.ToLower(job.RemoteOption), "remote") {
		return false
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

// GetAnalytics calculates analytics from job data
func (s *InMemoryStorage) GetAnalytics(jobs []models.Job) models.JobAnalytics {
	if len(jobs) == 0 {
		return models.JobAnalytics{}
	}

	analytics := models.JobAnalytics{
		TotalJobs:            len(jobs),
		ExperienceLevels:     make(map[string]int),
		RemoteOptions:        make(map[string]int),
		DegreeRequirements:   make(map[string]int),
		LocationDistribution: make(map[string]int),
		IndustryDistribution: make(map[string]int),
	}

	skillCounts := make(map[string]int)
	companyCounts := make(map[string]int)
	salaries := make([]float64, 0)

	for _, job := range jobs {
		// Experience levels
		if job.ExperienceLevel != "" {
			analytics.ExperienceLevels[job.ExperienceLevel]++
		}

		// Remote options
		if job.RemoteOption != "" {
			analytics.RemoteOptions[job.RemoteOption]++
		}

		// Degree requirements
		if job.DegreeRequired {
			analytics.DegreeRequirements["Required"]++
		} else {
			analytics.DegreeRequirements["Not Required"]++
		}

		// Location distribution
		if job.Location != "" {
			analytics.LocationDistribution[job.Location]++
		}

		// Industry distribution
		if job.Industry != "" {
			analytics.IndustryDistribution[job.Industry]++
		}

		// Skills counting
		for _, skill := range job.Skills {
			skillCounts[skill]++
		}

		// Company counting
		if job.Company != "" {
			companyCounts[job.Company]++
		}

		// Salary calculation
		if job.SalaryMin > 0 && job.SalaryMax > 0 {
			avgSalary := float64(job.SalaryMin+job.SalaryMax) / 2
			salaries = append(salaries, avgSalary)
		} else if job.SalaryMin > 0 {
			salaries = append(salaries, float64(job.SalaryMin))
		} else if job.SalaryMax > 0 {
			salaries = append(salaries, float64(job.SalaryMax))
		}
	}

	// Calculate salary statistics
	if len(salaries) > 0 {
		sort.Float64s(salaries)

		sum := 0.0
		for _, salary := range salaries {
			sum += salary
		}
		analytics.AverageSalary = sum / float64(len(salaries))

		analytics.SalaryRange = models.SalaryRange{
			Min:    int(salaries[0]),
			Max:    int(salaries[len(salaries)-1]),
			Median: percentile(salaries, 50),
			P25:    percentile(salaries, 25),
			P75:    percentile(salaries, 75),
		}
	}

	// Top skills
	analytics.TopSkills = getTopCounts(skillCounts, 10)

	// Top companies
	analytics.TopCompanies = getTopCompanies(companyCounts, 10)

	return analytics
}

// Helper function to calculate percentiles
func percentile(sortedData []float64, p float64) float64 {
	n := len(sortedData)
	if n == 0 {
		return 0
	}
	if n == 1 {
		return sortedData[0]
	}

	index := (p / 100.0) * float64(n-1)
	lower := int(index)
	upper := lower + 1

	if upper >= n {
		return sortedData[n-1]
	}

	weight := index - float64(lower)
	return sortedData[lower]*(1-weight) + sortedData[upper]*weight
}

// Helper function to get top skills
func getTopCounts(counts map[string]int, limit int) []models.SkillCount {
	type skillCountPair struct {
		skill string
		count int
	}

	var pairs []skillCountPair
	for skill, count := range counts {
		pairs = append(pairs, skillCountPair{skill, count})
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].count > pairs[j].count
	})

	result := make([]models.SkillCount, 0, limit)
	for i, pair := range pairs {
		if i >= limit {
			break
		}
		result = append(result, models.SkillCount{
			Skill: pair.skill,
			Count: pair.count,
		})
	}

	return result
}

// Helper function to get top companies
func getTopCompanies(counts map[string]int, limit int) []models.CompanyCount {
	type companyCountPair struct {
		company string
		count   int
	}

	var pairs []companyCountPair
	for company, count := range counts {
		pairs = append(pairs, companyCountPair{company, count})
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].count > pairs[j].count
	})

	result := make([]models.CompanyCount, 0, limit)
	for i, pair := range pairs {
		if i >= limit {
			break
		}
		result = append(result, models.CompanyCount{
			Company: pair.company,
			Count:   pair.count,
		})
	}

	return result
}
