// Job-related types
export interface Job {
  id: string;
  title: string;
  company: string;
  location: string;
  description: string;
  requirements: string[];
  skills: string[];
  salary_min?: number;
  salary_max?: number;
  salary_currency?: string;
  degree_required: boolean;
  experience_level: string;
  remote_option: string;
  posted_date: string;
  url: string;
  source: string;
  company_size?: string;
  industry?: string;
  benefits?: string[];
}

// Search filter types
export interface SearchFilters {
  job_title?: string;
  keywords?: string[];
  location?: string;
  remote_only?: boolean;
  min_salary?: number;
  max_salary?: number;
  experience_level?: string;
  degree_required?: boolean;
  skills?: string[];
  company_size?: string;
  industry?: string;
  limit?: number;
  offset?: number;
}

// Analytics types
export interface JobAnalytics {
  total_jobs: number;
  average_salary: number;
  salary_range: SalaryRange;
  top_skills: SkillCount[];
  top_companies: CompanyCount[];
  experience_levels: Record<string, number>;
  remote_options: Record<string, number>;
  degree_requirements: Record<string, number>;
  location_distribution: Record<string, number>;
  industry_distribution: Record<string, number>;
}

export interface SalaryRange {
  min: number;
  max: number;
  median: number;
  p25: number;
  p75: number;
}

export interface SkillCount {
  skill: string;
  count: number;
}

export interface CompanyCount {
  company: string;
  count: number;
}

// API response types
export interface SearchResponse {
  jobs: Job[];
  total: number;
  analytics: JobAnalytics;
  filters: SearchFilters;
}

// Experience level options
export const EXPERIENCE_LEVELS = [
  { value: '', label: 'Any Experience Level' },
  { value: 'entry', label: 'Entry Level (0-2 years)' },
  { value: 'mid', label: 'Mid Level (2-5 years)' },
  { value: 'senior', label: 'Senior Level (5+ years)' },
  { value: 'lead', label: 'Lead/Principal (7+ years)' },
];

// Remote options
export const REMOTE_OPTIONS = [
  { value: '', label: 'Any Location' },
  { value: 'remote', label: 'Remote Only' },
  { value: 'hybrid', label: 'Hybrid' },
  { value: 'onsite', label: 'On-site Only' },
];

// Degree requirement options
export const DEGREE_OPTIONS = [
  { value: '', label: 'Any Degree Requirement' },
  { value: 'false', label: 'No Degree Required' },
  { value: 'true', label: 'Degree Required' },
];

// Common IT skills for autocomplete
export const COMMON_SKILLS = [
  'Go', 'Golang', 'Python', 'JavaScript', 'TypeScript', 'Java', 'C++', 'C#', 'Rust', 'PHP', 'Ruby',
  'React', 'Angular', 'Vue.js', 'Node.js', 'Express', 'Django', 'Flask', 'Spring', 'Gin', 'Fiber',
  'MySQL', 'PostgreSQL', 'MongoDB', 'Redis', 'Elasticsearch', 'Cassandra', 'SQLite',
  'AWS', 'Azure', 'GCP', 'Docker', 'Kubernetes', 'Jenkins', 'GitLab', 'GitHub', 'Terraform', 'Ansible',
  'Linux', 'Unix', 'Git', 'SQL', 'NoSQL', 'REST', 'GraphQL', 'Microservices', 'API',
  'Agile', 'Scrum', 'DevOps', 'CI/CD', 'TDD', 'BDD'
];

// Job title suggestions
export const JOB_TITLE_SUGGESTIONS = [
  'Backend Developer',
  'Frontend Developer',
  'Full Stack Developer',
  'Go Developer',
  'Golang Developer',
  'Software Engineer',
  'DevOps Engineer',
  'Cloud Engineer',
  'API Developer',
  'Microservices Developer',
  'Senior Developer',
  'Junior Developer',
  'Lead Developer',
  'Principal Engineer'
];