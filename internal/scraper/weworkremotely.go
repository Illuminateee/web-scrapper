package scraper

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Illuminateee/web-scrapper.git/internal/models"
	"github.com/PuerkitoBio/goquery"
)

// WeWorkRemotelyScraper scrapes jobs from WeWorkRemotely.com
type WeWorkRemotelyScraper struct {
	*BaseScraper
}

// NewWeWorkRemotelyScraper creates a new WeWorkRemotely scraper
func NewWeWorkRemotelyScraper(client *http.Client) *WeWorkRemotelyScraper {
	return &WeWorkRemotelyScraper{
		BaseScraper: NewBaseScraper("WeWorkRemotely", "https://weworkremotely.com", client),
	}
}

// Scrape implements the JobScraper interface
func (w *WeWorkRemotelyScraper) Scrape(ctx context.Context, filters models.SearchFilters) ([]models.Job, error) {
	var jobs []models.Job

	// WeWorkRemotely has different categories, let's scrape programming jobs
	categories := []string{
		"remote-programming-jobs",
		"remote-devops-sysadmin-jobs",
		"remote-customer-support-jobs",
		"remote-design-jobs",
		"remote-sales-marketing-jobs",
	}

	for _, category := range categories {
		categoryJobs, err := w.scrapeCategory(ctx, category, filters)
		if err != nil {
			// Log error but continue with other categories
			continue
		}
		jobs = append(jobs, categoryJobs...)
	}

	return jobs, nil
}

func (w *WeWorkRemotelyScraper) scrapeCategory(ctx context.Context, category string, filters models.SearchFilters) ([]models.Job, error) {
	url := fmt.Sprintf("%s/remote-jobs/%s", w.baseURL, category)

	doc, err := w.FetchDocument(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %w", url, err)
	}

	var jobs []models.Job

	// Parse job listings
	doc.Find("article.job").Each(func(i int, s *goquery.Selection) {
		job := w.parseJobListing(s, category)
		if job.Title != "" && w.matchesFilters(job, filters) {
			jobs = append(jobs, job)
		}
	})

	return jobs, nil
}

func (w *WeWorkRemotelyScraper) parseJobListing(s *goquery.Selection, category string) models.Job {
	// Extract job title
	title := w.CleanText(s.Find("h2 a").Text())

	// Extract company name
	company := w.CleanText(s.Find(".company a").Text())

	// Extract job URL
	jobURL, _ := s.Find("h2 a").Attr("href")
	if jobURL != "" && !strings.HasPrefix(jobURL, "http") {
		jobURL = w.baseURL + jobURL
	}

	// Extract location (usually remote)
	location := "Remote"

	// Extract salary if available
	salaryText := w.CleanText(s.Find(".salary").Text())
	salaryMin, salaryMax := w.parseSalary(salaryText)

	// Extract posted date
	dateText := w.CleanText(s.Find("time").Text())
	postedDate := w.parseDate(dateText)

	// Determine experience level from title
	expLevel := w.DetermineExperienceLevel(title, "")

	// Extract skills from title and category
	skills := w.extractSkillsFromText(title + " " + category)

	// Create job description placeholder
	description := fmt.Sprintf("Remote %s position at %s", title, company)

	return models.Job{
		ID:              fmt.Sprintf("wwr-%s-%d", strings.ReplaceAll(strings.ToLower(title), " ", "-"), time.Now().Unix()),
		Title:           title,
		Company:         company,
		Location:        location,
		Description:     description,
		Skills:          skills,
		SalaryMin:       salaryMin,
		SalaryMax:       salaryMax,
		SalaryCurrency:  "USD",
		DegreeRequired:  false, // WeWorkRemotely jobs often don't require degrees
		ExperienceLevel: expLevel,
		RemoteOption:    "remote",
		PostedDate:      postedDate,
		URL:             jobURL,
		Source:          w.Name(),
		Industry:        w.categoryToIndustry(category),
	}
}

func (w *WeWorkRemotelyScraper) parseSalary(salaryText string) (int, int) {
	if salaryText == "" {
		return 0, 0
	}

	// Remove currency symbols and clean up
	cleaned := strings.ReplaceAll(salaryText, "$", "")
	cleaned = strings.ReplaceAll(cleaned, ",", "")
	cleaned = strings.ToLower(cleaned)

	// Look for salary ranges like "80k - 120k" or "80000-120000"
	re := regexp.MustCompile(`(\d+)k?\s*[-–]\s*(\d+)k?`)
	matches := re.FindStringSubmatch(cleaned)

	if len(matches) >= 3 {
		min, _ := strconv.Atoi(matches[1])
		max, _ := strconv.Atoi(matches[2])

		// Convert k to thousands
		if strings.Contains(matches[0], "k") {
			min *= 1000
			max *= 1000
		}

		return min, max
	}

	// Look for single salary values
	re = regexp.MustCompile(`(\d+)k?`)
	match := re.FindStringSubmatch(cleaned)
	if len(match) >= 2 {
		val, _ := strconv.Atoi(match[1])
		if strings.Contains(match[0], "k") {
			val *= 1000
		}
		return val, val
	}

	return 0, 0
}

func (w *WeWorkRemotelyScraper) parseDate(dateText string) time.Time {
	// WeWorkRemotely often uses relative dates like "2 days ago"
	now := time.Now()

	if strings.Contains(dateText, "today") {
		return now
	}

	if strings.Contains(dateText, "yesterday") {
		return now.AddDate(0, 0, -1)
	}

	// Look for "X days ago" pattern
	re := regexp.MustCompile(`(\d+)\s+days?\s+ago`)
	matches := re.FindStringSubmatch(strings.ToLower(dateText))
	if len(matches) >= 2 {
		if days, err := strconv.Atoi(matches[1]); err == nil {
			return now.AddDate(0, 0, -days)
		}
	}

	// Look for "X weeks ago" pattern
	re = regexp.MustCompile(`(\d+)\s+weeks?\s+ago`)
	matches = re.FindStringSubmatch(strings.ToLower(dateText))
	if len(matches) >= 2 {
		if weeks, err := strconv.Atoi(matches[1]); err == nil {
			return now.AddDate(0, 0, -weeks*7)
		}
	}

	// Default to now if we can't parse
	return now
}

func (w *WeWorkRemotelyScraper) extractSkillsFromText(text string) []string {
	var skills []string
	skillSet := make(map[string]bool)

	commonSkills := []string{
		"javascript", "python", "java", "go", "golang", "rust", "php", "ruby", "swift", "kotlin",
		"react", "angular", "vue", "node.js", "django", "flask", "spring", "express",
		"mysql", "postgresql", "mongodb", "redis", "elasticsearch",
		"aws", "azure", "gcp", "docker", "kubernetes", "jenkins",
		"git", "linux", "sql", "nosql", "api", "rest", "graphql",
		"devops", "sysadmin", "frontend", "backend", "fullstack", "full-stack",
		"design", "ui", "ux", "marketing", "sales", "support",
	}

	textLower := strings.ToLower(text)
	for _, skill := range commonSkills {
		if strings.Contains(textLower, skill) && !skillSet[skill] {
			skills = append(skills, skill)
			skillSet[skill] = true
		}
	}

	return skills
}

func (w *WeWorkRemotelyScraper) categoryToIndustry(category string) string {
	switch category {
	case "remote-programming-jobs", "remote-devops-sysadmin-jobs":
		return "Technology"
	case "remote-design-jobs":
		return "Design"
	case "remote-sales-marketing-jobs":
		return "Sales & Marketing"
	case "remote-customer-support-jobs":
		return "Customer Service"
	default:
		return "Technology"
	}
}

func (w *WeWorkRemotelyScraper) matchesFilters(job models.Job, filters models.SearchFilters) bool {
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

	return true
}
