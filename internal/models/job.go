package models

import "time"

// Job represents a single job posting
type Job struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	Company         string    `json:"company"`
	Location        string    `json:"location"`
	Description     string    `json:"description"`
	Requirements    []string  `json:"requirements"`
	Skills          []string  `json:"skills"`
	SalaryMin       int       `json:"salary_min,omitempty"`
	SalaryMax       int       `json:"salary_max,omitempty"`
	SalaryCurrency  string    `json:"salary_currency,omitempty"`
	DegreeRequired  bool      `json:"degree_required"`
	ExperienceLevel string    `json:"experience_level"` // entry, mid, senior, lead
	RemoteOption    string    `json:"remote_option"`    // onsite, remote, hybrid
	PostedDate      time.Time `json:"posted_date"`
	URL             string    `json:"url"`
	Source          string    `json:"source"` // indeed, linkedin, glassdoor
	CompanySize     string    `json:"company_size,omitempty"`
	Industry        string    `json:"industry,omitempty"`
	Benefits        []string  `json:"benefits,omitempty"`
}

// SearchFilters represents the search criteria
type SearchFilters struct {
	JobTitle        string   `json:"job_title"`
	Keywords        []string `json:"keywords"`
	Location        string   `json:"location"`
	Locations       []string `json:"locations"` // Multiple locations support
	RemoteOnly      bool     `json:"remote_only"`
	MinSalary       int      `json:"min_salary"`
	MaxSalary       int      `json:"max_salary"`
	ExperienceLevel string   `json:"experience_level"`
	DegreeRequired  *bool    `json:"degree_required"` // nil = any, true = required, false = not required
	Skills          []string `json:"skills"`
	CompanySize     string   `json:"company_size"`
	Industry        string   `json:"industry"`
	JobCategory     string   `json:"job_category"` // healthcare, finance, retail, etc.
	JobSites        []string `json:"job_sites"`    // Custom job sites URLs
	Limit           int      `json:"limit"`
	Offset          int      `json:"offset"`
}

// SearchResponse represents the response from job search
type SearchResponse struct {
	Jobs      []Job         `json:"jobs"`
	Total     int           `json:"total"`
	Analytics JobAnalytics  `json:"analytics"`
	Filters   SearchFilters `json:"filters"`
}

// JobAnalytics provides insights from the job search results
type JobAnalytics struct {
	TotalJobs            int            `json:"total_jobs"`
	AverageSalary        float64        `json:"average_salary"`
	SalaryRange          SalaryRange    `json:"salary_range"`
	TopSkills            []SkillCount   `json:"top_skills"`
	TopCompanies         []CompanyCount `json:"top_companies"`
	ExperienceLevels     map[string]int `json:"experience_levels"`
	RemoteOptions        map[string]int `json:"remote_options"`
	DegreeRequirements   map[string]int `json:"degree_requirements"`
	LocationDistribution map[string]int `json:"location_distribution"`
	IndustryDistribution map[string]int `json:"industry_distribution"`
}

// SalaryRange represents salary statistics
type SalaryRange struct {
	Min    int     `json:"min"`
	Max    int     `json:"max"`
	Median float64 `json:"median"`
	P25    float64 `json:"p25"`
	P75    float64 `json:"p75"`
}

// SkillCount represents skill frequency
type SkillCount struct {
	Skill string `json:"skill"`
	Count int    `json:"count"`
}

// CompanyCount represents company hiring frequency
type CompanyCount struct {
	Company string `json:"company"`
	Count   int    `json:"count"`
}

// ScrapingResult represents the result from a single scraping operation
type ScrapingResult struct {
	Jobs   []Job  `json:"jobs"`
	Source string `json:"source"`
	Error  error  `json:"error,omitempty"`
}

// JobSiteConfig represents configuration for a custom job site
type JobSiteConfig struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description,omitempty"`
	Active      bool   `json:"active"`
}

// SearchRequest represents a complete search request with job sites
type SearchRequest struct {
	Filters  SearchFilters   `json:"filters"`
	JobSites []JobSiteConfig `json:"job_sites"`
}
