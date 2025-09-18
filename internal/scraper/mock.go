package scraper

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/Illuminateee/web-scrapper.git/internal/models"
)

// MockJobScraper simulates job scraping for development and testing
type MockJobScraper struct {
	*BaseScraper
	jobTemplates []jobTemplate
}

type jobTemplate struct {
	titleTemplates   []string
	companies        []string
	locations        []string
	salaryRanges     []salaryRange
	skillSets        [][]string
	experienceLevels []string
	remoteOptions    []string
	industries       []string
}

type salaryRange struct {
	min, max int
}

// NewMockJobScraper creates a new mock job scraper
func NewMockJobScraper(name string) *MockJobScraper {
	scraper := &MockJobScraper{
		BaseScraper: &BaseScraper{
			name:    name,
			baseURL: "https://mock-job-site.com",
		},
	}

	scraper.initializeTemplates()
	return scraper
}

func (m *MockJobScraper) initializeTemplates() {
	m.jobTemplates = []jobTemplate{
		// Technology Jobs
		{
			titleTemplates: []string{
				"Backend Developer", "Senior Backend Developer", "Junior Backend Developer",
				"Go Developer", "Senior Go Developer", "Golang Developer",
				"Full Stack Developer", "Senior Full Stack Developer",
				"Software Engineer", "Senior Software Engineer",
				"DevOps Engineer", "Cloud Engineer", "API Developer",
				"Frontend Developer", "React Developer", "Node.js Developer",
				"Data Scientist", "Machine Learning Engineer", "AI Engineer",
				"Mobile Developer", "iOS Developer", "Android Developer",
				"QA Engineer", "Test Automation Engineer", "Security Engineer",
			},
			companies: []string{
				"TechCorp Inc", "InnovateSoft", "CloudTech Solutions", "DataFlow Systems",
				"NextGen Software", "StreamlineIT", "ScaleUp Technologies", "DevMasters LLC",
				"AgileWorks", "CodeCraft Solutions", "TechPioneer", "DigitalForge",
				"CloudNative Systems", "MicroServices Co", "APIBuilders Inc",
			},
			locations: []string{
				"San Francisco, CA", "New York, NY", "Austin, TX", "Seattle, WA",
				"Remote", "Boston, MA", "Denver, CO", "Chicago, IL", "Los Angeles, CA",
				"Portland, OR", "Atlanta, GA", "Remote - US", "Hybrid - New York",
				"Remote - Global", "Toronto, ON", "London, UK", "Berlin, Germany",
			},
			salaryRanges: []salaryRange{
				{60000, 90000},   // Junior
				{80000, 120000},  // Mid-level
				{110000, 160000}, // Senior
				{140000, 200000}, // Lead/Principal
			},
			skillSets: [][]string{
				{"Go", "PostgreSQL", "Docker", "Kubernetes", "REST APIs"},
				{"Golang", "MongoDB", "Redis", "AWS", "Microservices"},
				{"Go", "MySQL", "Git", "Linux", "CI/CD"},
				{"Golang", "React", "Node.js", "TypeScript", "GraphQL"},
				{"Go", "Python", "Docker", "Jenkins", "Terraform"},
				{"Golang", "Vue.js", "PostgreSQL", "Azure", "DevOps"},
				{"Go", "Java", "Spring", "Elasticsearch", "Kafka"},
				{"Golang", "Angular", "MongoDB", "GCP", "Prometheus"},
			},
			experienceLevels: []string{"entry", "mid", "senior", "lead"},
			remoteOptions:    []string{"onsite", "remote", "hybrid"},
			industries: []string{
				"Technology", "Finance", "Healthcare", "E-commerce", "Gaming",
				"EdTech", "FinTech", "SaaS", "Startups", "Enterprise Software",
			},
		},
		// Healthcare Jobs
		{
			titleTemplates: []string{
				"Registered Nurse", "Nurse Practitioner", "Medical Assistant",
				"Physical Therapist", "Occupational Therapist", "Pharmacist",
				"Medical Technologist", "Healthcare Administrator", "Medical Scribe",
				"Clinical Research Coordinator", "Healthcare Data Analyst",
				"Medical Device Sales", "Healthcare IT Specialist",
			},
			companies: []string{
				"General Hospital", "MedCenter Health", "Healthcare Solutions Inc",
				"Regional Medical Center", "Community Health Network", "MedTech Corp",
				"Healthcare Analytics", "Medical Innovations", "PharmaCare",
			},
			locations: []string{
				"New York, NY", "Los Angeles, CA", "Chicago, IL", "Houston, TX",
				"Phoenix, AZ", "Philadelphia, PA", "San Antonio, TX", "Dallas, TX",
				"Remote", "Boston, MA", "Atlanta, GA", "Miami, FL",
			},
			salaryRanges: []salaryRange{
				{45000, 65000},   // Entry
				{60000, 85000},   // Mid
				{80000, 110000},  // Senior
				{100000, 140000}, // Lead
			},
			skillSets: [][]string{
				{"Patient Care", "Medical Records", "EMR", "HIPAA"},
				{"Clinical Skills", "Patient Assessment", "Medical Terminology"},
				{"Healthcare Technology", "Medical Devices", "Data Analysis"},
				{"Pharmacy Knowledge", "Drug Interactions", "Patient Counseling"},
			},
			experienceLevels: []string{"entry", "mid", "senior", "lead"},
			remoteOptions:    []string{"onsite", "hybrid"},
			industries:       []string{"Healthcare", "Medical", "Pharmaceuticals", "Biotechnology"},
		},
		// Finance Jobs
		{
			titleTemplates: []string{
				"Financial Analyst", "Senior Financial Analyst", "Investment Advisor",
				"Accountant", "Senior Accountant", "Tax Specialist", "Auditor",
				"Banking Associate", "Loan Officer", "Financial Planner",
				"Risk Analyst", "Credit Analyst", "Treasury Analyst",
				"Compliance Officer", "Investment Banking Analyst",
			},
			companies: []string{
				"First National Bank", "Investment Partners", "Financial Services Corp",
				"Capital Management", "Trust & Investment", "Financial Advisors Inc",
				"Banking Solutions", "Credit Union", "Wealth Management",
			},
			locations: []string{
				"New York, NY", "Chicago, IL", "San Francisco, CA", "Boston, MA",
				"Remote", "Charlotte, NC", "Dallas, TX", "Los Angeles, CA",
				"Miami, FL", "Denver, CO", "Seattle, WA",
			},
			salaryRanges: []salaryRange{
				{50000, 70000},   // Entry
				{70000, 95000},   // Mid
				{90000, 125000},  // Senior
				{120000, 170000}, // Lead
			},
			skillSets: [][]string{
				{"Excel", "Financial Modeling", "SQL", "Python"},
				{"QuickBooks", "SAP", "Financial Reporting", "GAAP"},
				{"Risk Management", "Compliance", "Auditing", "Taxation"},
				{"Investment Analysis", "Portfolio Management", "Bloomberg Terminal"},
			},
			experienceLevels: []string{"entry", "mid", "senior", "lead"},
			remoteOptions:    []string{"onsite", "remote", "hybrid"},
			industries:       []string{"Finance", "Banking", "Investment", "Insurance"},
		},
		// Retail & Sales Jobs
		{
			titleTemplates: []string{
				"Sales Associate", "Store Manager", "Assistant Manager",
				"Customer Service Representative", "Cashier", "Sales Manager",
				"Retail Specialist", "Visual Merchandiser", "Inventory Coordinator",
				"Regional Sales Manager", "Account Manager", "Sales Director",
			},
			companies: []string{
				"RetailMart", "Fashion Forward", "Electronics Plus", "Home & Garden",
				"SportsCorp", "TechRetail", "Department Store Co", "Specialty Retail",
			},
			locations: []string{
				"Nationwide", "New York, NY", "Los Angeles, CA", "Chicago, IL",
				"Houston, TX", "Miami, FL", "Atlanta, GA", "Phoenix, AZ",
				"Dallas, TX", "Philadelphia, PA",
			},
			salaryRanges: []salaryRange{
				{25000, 35000},  // Entry
				{35000, 50000},  // Mid
				{50000, 70000},  // Senior
				{70000, 100000}, // Lead
			},
			skillSets: [][]string{
				{"Customer Service", "POS Systems", "Inventory Management"},
				{"Sales Techniques", "Product Knowledge", "Team Leadership"},
				{"Visual Merchandising", "Store Operations", "Cash Handling"},
				{"CRM Software", "Sales Analytics", "Training & Development"},
			},
			experienceLevels: []string{"entry", "mid", "senior", "lead"},
			remoteOptions:    []string{"onsite", "hybrid"},
			industries:       []string{"Retail", "Consumer Goods", "Fashion", "Electronics"},
		},
	}
}

// Scrape simulates scraping jobs based on search filters
func (m *MockJobScraper) Scrape(ctx context.Context, filters models.SearchFilters) ([]models.Job, error) {
	// Simulate network delay
	time.Sleep(time.Millisecond * time.Duration(500+rand.Intn(1000)))

	var jobs []models.Job

	// Select template based on job category or use all templates
	var templatesToUse []jobTemplate
	if filters.JobCategory != "" {
		switch strings.ToLower(filters.JobCategory) {
		case "technology", "tech", "it":
			templatesToUse = []jobTemplate{m.jobTemplates[0]}
		case "healthcare", "medical":
			if len(m.jobTemplates) > 1 {
				templatesToUse = []jobTemplate{m.jobTemplates[1]}
			}
		case "finance", "banking":
			if len(m.jobTemplates) > 2 {
				templatesToUse = []jobTemplate{m.jobTemplates[2]}
			}
		case "retail", "sales":
			if len(m.jobTemplates) > 3 {
				templatesToUse = []jobTemplate{m.jobTemplates[3]}
			}
		default:
			templatesToUse = m.jobTemplates
		}
	} else {
		templatesToUse = m.jobTemplates
	}

	// Generate jobs from selected templates
	for _, template := range templatesToUse {
		// Generate 10-20 jobs per template
		jobCount := 10 + rand.Intn(11)

		for i := 0; i < jobCount; i++ {
			job := m.generateJob(template, filters, i)

			// Apply basic filtering
			if m.matchesBasicFilters(job, filters) {
				jobs = append(jobs, job)
			}
		}
	}

	return jobs, nil
}

func (m *MockJobScraper) generateJob(template jobTemplate, filters models.SearchFilters, index int) models.Job {
	// Create a new random source for this job generation
	source := rand.NewSource(time.Now().UnixNano() + int64(index))
	rng := rand.New(source)

	// Select random elements
	title := template.titleTemplates[rng.Intn(len(template.titleTemplates))]
	company := template.companies[rng.Intn(len(template.companies))]
	location := template.locations[rng.Intn(len(template.locations))]
	skillSet := template.skillSets[rng.Intn(len(template.skillSets))]
	expLevel := template.experienceLevels[rng.Intn(len(template.experienceLevels))]
	remoteOption := template.remoteOptions[rng.Intn(len(template.remoteOptions))]
	industry := template.industries[rng.Intn(len(template.industries))]

	// Generate salary based on experience level
	var salaryRange salaryRange
	switch expLevel {
	case "entry":
		salaryRange = template.salaryRanges[0]
	case "mid":
		salaryRange = template.salaryRanges[1]
	case "senior":
		salaryRange = template.salaryRanges[2]
	case "lead":
		salaryRange = template.salaryRanges[3]
	default:
		salaryRange = template.salaryRanges[1]
	}

	// Add some variance to salary
	variance := (salaryRange.max - salaryRange.min) / 10
	minSalary := salaryRange.min + rng.Intn(variance) - variance/2
	maxSalary := salaryRange.max + rng.Intn(variance) - variance/2

	// Generate description
	description := m.generateJobDescription(title, skillSet, expLevel)

	// Generate requirements
	requirements := m.generateRequirements(expLevel, skillSet)

	// Determine if degree is required (60% chance for mock data)
	degreeRequired := rng.Float32() < 0.6

	// Generate posted date (within last 30 days)
	daysAgo := rng.Intn(30)
	postedDate := time.Now().AddDate(0, 0, -daysAgo)

	return models.Job{
		ID:              fmt.Sprintf("%s-%d-%d", m.Name(), time.Now().Unix(), index),
		Title:           title,
		Company:         company,
		Location:        location,
		Description:     description,
		Requirements:    requirements,
		Skills:          skillSet,
		SalaryMin:       minSalary,
		SalaryMax:       maxSalary,
		SalaryCurrency:  "USD",
		DegreeRequired:  degreeRequired,
		ExperienceLevel: expLevel,
		RemoteOption:    remoteOption,
		PostedDate:      postedDate,
		URL:             fmt.Sprintf("https://%s.com/jobs/%d", strings.ToLower(m.Name()), index),
		Source:          m.Name(),
		Industry:        industry,
		Benefits:        m.generateBenefits(),
	}
}

func (m *MockJobScraper) generateJobDescription(title string, skills []string, expLevel string) string {
	descriptions := []string{
		fmt.Sprintf("We are looking for a talented %s to join our growing team. You will be responsible for developing and maintaining our backend services using modern technologies.", title),
		fmt.Sprintf("Join our innovative team as a %s and help build scalable solutions that serve millions of users worldwide.", title),
		fmt.Sprintf("Exciting opportunity for a %s to work with cutting-edge technologies and contribute to our mission of transforming the industry.", title),
	}

	baseDesc := descriptions[rand.Intn(len(descriptions))]

	// Add skills section
	skillsText := fmt.Sprintf(" Key technologies include: %s.", strings.Join(skills, ", "))

	// Add experience requirements
	var expText string
	switch expLevel {
	case "entry":
		expText = " This is a great opportunity for new graduates or developers with 0-2 years of experience."
	case "mid":
		expText = " We're looking for someone with 2-5 years of experience in software development."
	case "senior":
		expText = " This role requires 5+ years of experience and strong leadership skills."
	case "lead":
		expText = " We need an experienced leader with 7+ years of experience to guide our technical direction."
	}

	return baseDesc + skillsText + expText
}

func (m *MockJobScraper) generateRequirements(expLevel string, skills []string) []string {
	baseReqs := []string{
		"Strong problem-solving skills",
		"Excellent communication skills",
		"Ability to work in a team environment",
		"Experience with version control (Git)",
	}

	// Add skill requirements
	for _, skill := range skills {
		baseReqs = append(baseReqs, fmt.Sprintf("Experience with %s", skill))
	}

	// Add experience-specific requirements
	switch expLevel {
	case "entry":
		baseReqs = append(baseReqs, "Bachelor's degree or equivalent experience", "Eagerness to learn new technologies")
	case "mid":
		baseReqs = append(baseReqs, "2-5 years of software development experience", "Experience with agile methodologies")
	case "senior":
		baseReqs = append(baseReqs, "5+ years of software development experience", "Experience mentoring junior developers", "Strong system design skills")
	case "lead":
		baseReqs = append(baseReqs, "7+ years of software development experience", "Proven leadership experience", "Excellent architectural skills")
	}

	return baseReqs
}

func (m *MockJobScraper) generateBenefits() []string {
	allBenefits := []string{
		"Health insurance",
		"Dental insurance",
		"Vision insurance",
		"401(k) matching",
		"Flexible PTO",
		"Remote work options",
		"Professional development budget",
		"Stock options",
		"Life insurance",
		"Gym membership",
		"Free meals",
		"Flexible hours",
	}

	// Return 3-6 random benefits
	benefitCount := 3 + rand.Intn(4)
	selectedBenefits := make([]string, 0, benefitCount)

	// Shuffle and select
	shuffled := make([]string, len(allBenefits))
	copy(shuffled, allBenefits)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	for i := 0; i < benefitCount && i < len(shuffled); i++ {
		selectedBenefits = append(selectedBenefits, shuffled[i])
	}

	return selectedBenefits
}

func (m *MockJobScraper) matchesBasicFilters(job models.Job, filters models.SearchFilters) bool {
	// Job title filter
	if filters.JobTitle != "" {
		if !strings.Contains(strings.ToLower(job.Title), strings.ToLower(filters.JobTitle)) {
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
