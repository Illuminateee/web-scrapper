# 🔍 Mock Concurrent Job Market Web Scraper

A powerful, concurrent web scraper built with **Go** and **React.js** that aggregates job postings from multiple sources, providing insights into the IT job market with advanced filtering capabilities.

## ✨ Features

- **🚀 Concurrent Scraping**: Multi-threaded scraping from multiple job sites
- **📊 Real-time Analytics**: Salary trends, skill analysis, and company hiring patterns  
- **🔍 Advanced Filtering**: Filter by degree requirements, experience level, skills, salary range
- **💼 IT Focus**: Specialized for Backend Developer, Golang Developer, Full-stack roles
- **🌐 Dynamic Search**: User-configurable search parameters
- **📈 Market Insights**: Comprehensive job market analytics
- **⚡ Fast API**: RESTful API with in-memory caching

## 🏗️ Architecture

```
web-scrapper/
├── cmd/server/           # Main application entry point
├── internal/
│   ├── api/             # REST API handlers
│   ├── models/          # Data structures
│   ├── scraper/         # Concurrent scraping logic
│   └── storage/         # In-memory data storage
├── web/frontend/        # React.js frontend
└── api-tester.html      # API testing interface
```

## 🚀 Quick Start

### Prerequisites

- Go 1.19+ 
- Node.js 16+ (optional, for React frontend)
- Git

### Backend Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/Illuminateee/web-scrapper.git
   cd web-scrapper
   ```

2. **Install Go dependencies**
   ```bash
   go mod tidy
   ```

3. **Start the server**
   ```bash
   go run cmd/server/main.go
   ```

   The server will start on `http://localhost:8080`

### Testing the API

#### Option 1: HTML Tester (Recommended)

Open `api-tester.html` in your browser to use the interactive testing interface.

#### Option 2: Command Line

```bash
# Health check
curl -X GET "http://localhost:8080/api/v1/health"

# Search for jobs
curl -X GET "http://localhost:8080/api/v1/jobs/search?title=golang+developer&experience_level=mid"

# Get analytics
curl -X GET "http://localhost:8080/api/v1/analytics"
```

#### Option 3: PowerShell (Windows)

```powershell
# Health check
Invoke-WebRequest -Uri "http://localhost:8080/api/v1/health" -Method GET

# Job search  
Invoke-WebRequest -Uri "http://localhost:8080/api/v1/jobs/search?title=backend+developer" -Method GET
```

## 📡 API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

### Endpoints

#### `GET /health`
Health check endpoint
```json
{
  "status": "healthy",
  "time": "2025-09-18T10:30:00Z"
}
```

#### `GET /jobs/search`
Search for jobs with filters

**Query Parameters:**
- `title` (string): Job title search term
- `keywords` (string): Comma-separated keywords
- `location` (string): Location filter
- `remote_only` (boolean): Filter for remote jobs only
- `min_salary` (integer): Minimum salary
- `max_salary` (integer): Maximum salary  
- `experience_level` (string): `entry`, `mid`, `senior`, `lead`
- `degree_required` (boolean): Filter by degree requirement
- `skills` (string): Comma-separated required skills
- `limit` (integer): Results per page (default: 50)
- `offset` (integer): Pagination offset

**Example Request:**
```
GET /jobs/search?title=golang&experience_level=mid&degree_required=false&min_salary=80000&remote_only=true
```

**Response:**
```json
{
  "jobs": [
    {
      "id": "mock-123",
      "title": "Senior Go Developer",
      "company": "TechCorp Inc",
      "location": "Remote",
      "description": "We are looking for a talented Senior Go Developer...",
      "skills": ["Go", "PostgreSQL", "Docker", "Kubernetes"],
      "salary_min": 110000,
      "salary_max": 160000,
      "salary_currency": "USD",
      "degree_required": false,
      "experience_level": "senior",
      "remote_option": "remote",
      "posted_date": "2025-09-15T00:00:00Z",
      "source": "MockJobSite1"
    }
  ],
  "total": 25,
  "analytics": {
    "total_jobs": 25,
    "average_salary": 125000,
    "top_skills": [
      {"skill": "Go", "count": 20},
      {"skill": "Docker", "count": 18}
    ]
  }
}
```

#### `GET /jobs/{id}`
Get specific job by ID

#### `GET /analytics`
Get job market analytics

**Response:**
```json
{
  "total_jobs": 150,
  "average_salary": 118500,
  "salary_range": {
    "min": 60000,
    "max": 200000,
    "median": 115000,
    "p25": 85000,
    "p75": 145000
  },
  "top_skills": [
    {"skill": "Go", "count": 75},
    {"skill": "Docker", "count": 68},
    {"skill": "Kubernetes", "count": 52}
  ],
  "top_companies": [
    {"company": "TechCorp Inc", "count": 8},
    {"company": "InnovateSoft", "count": 6}
  ],
  "experience_levels": {
    "entry": 25,
    "mid": 65,
    "senior": 50,
    "lead": 10
  }
}
```

#### `POST /cache/clear`
Clear the job cache

## 🔧 Configuration

### Environment Variables

- `PORT`: Server port (default: 8080)

### Scraper Configuration

The scraper uses mock data for demonstration. To integrate real job sites:

1. Implement the `JobScraper` interface in `internal/scraper/`
2. Add your scraper to the `ScraperManager` in `internal/api/handlers.go`
3. Configure rate limiting and request headers appropriately

## 📊 Mock Data

The application includes realistic mock data generators that simulate:
- **15-30 jobs per scraper** with realistic variations
- **IT-focused job titles**: Backend Developer, Go Developer, Full-stack Developer
- **Realistic salaries** based on experience level
- **Common tech skills**: Go, Docker, Kubernetes, PostgreSQL, etc.
- **Various experience levels** and degree requirements
- **Remote/hybrid/onsite** options

## 🛠️ Development

### Project Structure

```
internal/
├── api/
│   └── handlers.go      # HTTP request handlers
├── models/
│   └── job.go          # Data structures
├── scraper/
│   ├── base.go         # Base scraper functionality  
│   ├── mock.go         # Mock data generator
│   └── rate_limiter.go # Request rate limiting
└── storage/
    └── memory.go       # In-memory storage
```

### Adding New Scrapers

1. **Implement the interface:**
   ```go
   type JobScraper interface {
       Name() string
       Scrape(ctx context.Context, filters models.SearchFilters) ([]models.Job, error)
       GetBaseURL() string
   }
   ```

2. **Add to scraper manager:**
   ```go
   scraperManager.AddScraper(NewYourScraper())
   ```

### Key Features

- **Concurrent Processing**: All scrapers run in parallel using goroutines
- **Rate Limiting**: Configurable request rate limiting per scraper
- **Error Handling**: Graceful error handling with detailed logging
- **Caching**: In-memory storage with search and analytics capabilities
- **CORS Support**: Cross-origin resource sharing for frontend integration

## 🎯 Use Cases

### For Job Seekers
- Find IT positions without degree requirements
- Compare salaries across different experience levels
- Identify in-demand skills for career planning
- Discover companies actively hiring

### For Recruiters  
- Analyze market salary trends
- Identify skill gaps in the market
- Monitor competitor hiring patterns
- Track job posting trends

### For Market Research
- IT job market analysis
- Skill demand forecasting
- Salary benchmarking
- Regional job market insights

## 🚦 Testing Examples

### Basic Search
```bash
curl "http://localhost:8080/api/v1/jobs/search?title=developer"
```

### Advanced Filter
```bash
curl "http://localhost:8080/api/v1/jobs/search?title=golang&experience_level=senior&degree_required=false&min_salary=100000&skills=docker,kubernetes"
```

### Remote Jobs Only
```bash
curl "http://localhost:8080/api/v1/jobs/search?remote_only=true&experience_level=mid"
```

### Market Analytics
```bash
curl "http://localhost:8080/api/v1/analytics"
```

## 🔄 Future Enhancements

- **Database Integration**: PostgreSQL/MongoDB for persistent storage
- **Real Scraper Integration**: Indeed, LinkedIn, Glassdoor scrapers
- **Advanced Analytics**: Trend analysis, prediction models
- **User Authentication**: Saved searches, job alerts
- **Email Notifications**: Job alert system
- **API Rate Limiting**: Per-user rate limiting
- **Caching**: Redis integration for distributed caching

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License.

## 👨‍💻 Author

Created by Illuminateee - A concurrent job market scraper focused on IT roles without degree barriers.

---

**Happy Job Hunting! 🎯**
