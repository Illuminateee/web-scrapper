package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Illuminateee/web-scrapper.git/internal/models"
	"github.com/Illuminateee/web-scrapper.git/internal/scraper"
)

func main() {
	fmt.Println("🔍 Testing Job Scrapers")
	fmt.Println("=" + fmt.Sprintf("%*s", 50, "="))

	// Create HTTP client
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    30 * time.Second,
			DisableCompression: true,
		},
	}

	// Create test search filters
	filters := models.SearchFilters{
		Keywords:        []string{"developer", "engineer", "programmer"},
		Locations:       []string{"remote", "usa", "europe"},
		JobCategory:     "technology",
		ExperienceLevel: "mid",
		MinSalary:       50000,
		MaxSalary:       150000,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Test Mock Scraper
	fmt.Println("\n📝 Testing Mock Scraper...")
	testMockScraper(ctx, filters)

	// Test RemoteOK Scraper
	fmt.Println("\n🌐 Testing RemoteOK Scraper...")
	testRemoteOKScraper(ctx, filters, client)

	// Test WeWorkRemotely Scraper
	fmt.Println("\n💼 Testing WeWorkRemotely Scraper...")
	testWeWorkRemotelyScraper(ctx, filters, client)

	// Test Scraper Registry
	fmt.Println("\n🏭 Testing Scraper Registry...")
	testScraperRegistry(ctx, filters)

	fmt.Println("\n✅ All tests completed!")
}

func testMockScraper(ctx context.Context, filters models.SearchFilters) {
	mockScraper := scraper.NewMockJobScraper("TestMockScraper")

	jobs, err := mockScraper.Scrape(ctx, filters)
	if err != nil {
		log.Printf("❌ Mock scraper error: %v", err)
		return
	}

	fmt.Printf("   ✅ Mock scraper found %d jobs\n", len(jobs))
	if len(jobs) > 0 {
		job := jobs[0]
		fmt.Printf("   📋 Sample job: %s at %s (Salary: $%d-%d)\n",
			job.Title, job.Company, job.SalaryMin, job.SalaryMax)
	}
}

func testRemoteOKScraper(ctx context.Context, filters models.SearchFilters, client *http.Client) {
	remoteScraper := scraper.NewRemoteOKScraper(client)

	jobs, err := remoteScraper.Scrape(ctx, filters)
	if err != nil {
		log.Printf("❌ RemoteOK scraper error: %v", err)
		return
	}

	fmt.Printf("   ✅ RemoteOK scraper found %d jobs\n", len(jobs))
	if len(jobs) > 0 {
		job := jobs[0]
		fmt.Printf("   📋 Sample job: %s at %s\n", job.Title, job.Company)
		fmt.Printf("   🔗 URL: %s\n", job.URL)
		if len(job.Skills) > 0 {
			fmt.Printf("   🏷️  Skills: %v\n", job.Skills[:min(3, len(job.Skills))])
		}
	}
}

func testWeWorkRemotelyScraper(ctx context.Context, filters models.SearchFilters, client *http.Client) {
	wwrScraper := scraper.NewWeWorkRemotelyScraper(client)

	jobs, err := wwrScraper.Scrape(ctx, filters)
	if err != nil {
		log.Printf("❌ WeWorkRemotely scraper error: %v", err)
		return
	}

	fmt.Printf("   ✅ WeWorkRemotely scraper found %d jobs\n", len(jobs))
	if len(jobs) > 0 {
		job := jobs[0]
		fmt.Printf("   📋 Sample job: %s at %s\n", job.Title, job.Company)
		fmt.Printf("   🔗 URL: %s\n", job.URL)
		if job.Description != "" {
			desc := job.Description
			if len(desc) > 100 {
				desc = desc[:100] + "..."
			}
			fmt.Printf("   📄 Description: %s\n", desc)
		}
	}
}

func testScraperRegistry(ctx context.Context, filters models.SearchFilters) {
	registry := scraper.NewScraperRegistry()

	// List all scrapers
	configs := registry.ListScrapers()
	fmt.Printf("   📊 Registry has %d scrapers configured:\n", len(configs))

	for _, config := range configs {
		status := "disabled"
		if config.Enabled {
			status = "enabled"
		}
		if config.RequiresAuth {
			status += " (auth required)"
		}
		fmt.Printf("   • %s: %s [%s]\n", config.Name, status, config.Type)
	}

	// Test enabled scrapers
	enabledScrapers := registry.GetEnabledScrapers()
	fmt.Printf("\n   🚀 Testing %d enabled scrapers:\n", len(enabledScrapers))

	totalJobs := 0
	for _, scraper := range enabledScrapers {
		jobs, err := scraper.Scrape(ctx, filters)
		if err != nil {
			log.Printf("   ❌ Scraper error: %v", err)
			continue
		}

		fmt.Printf("   ✅ Found %d jobs from this scraper\n", len(jobs))
		totalJobs += len(jobs)
	}

	fmt.Printf("   🎯 Total jobs found: %d\n", totalJobs)
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
