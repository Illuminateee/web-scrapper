import axios from 'axios';
import { SearchFilters, SearchResponse, Job, JobAnalytics } from './types';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api/v1';

const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 60000, // 60 seconds timeout for scraping operations
});

export class JobService {
  // Search jobs with filters
  static async searchJobs(filters: SearchFilters): Promise<SearchResponse> {
    const params = new URLSearchParams();
    
    if (filters.job_title) params.append('title', filters.job_title);
    if (filters.keywords?.length) params.append('keywords', filters.keywords.join(','));
    if (filters.location) params.append('location', filters.location);
    if (filters.remote_only) params.append('remote_only', 'true');
    if (filters.min_salary) params.append('min_salary', filters.min_salary.toString());
    if (filters.max_salary) params.append('max_salary', filters.max_salary.toString());
    if (filters.experience_level) params.append('experience_level', filters.experience_level);
    if (filters.degree_required !== undefined) {
      params.append('degree_required', filters.degree_required.toString());
    }
    if (filters.skills?.length) params.append('skills', filters.skills.join(','));
    if (filters.company_size) params.append('company_size', filters.company_size);
    if (filters.industry) params.append('industry', filters.industry);
    if (filters.limit) params.append('limit', filters.limit.toString());
    if (filters.offset) params.append('offset', filters.offset.toString());
    
    const response = await api.get(`/jobs/search?${params.toString()}`);
    return response.data;
  }

  // Get specific job by ID
  static async getJob(id: string): Promise<Job> {
    const response = await api.get(`/jobs/${id}`);
    return response.data;
  }

  // Get analytics data
  static async getAnalytics(): Promise<JobAnalytics> {
    const response = await api.get('/analytics');
    return response.data;
  }

  // Clear cache
  static async clearCache(): Promise<void> {
    await api.post('/cache/clear');
  }

  // Health check
  static async healthCheck(): Promise<{ status: string; time: string }> {
    const response = await api.get('/health');
    return response.data;
  }
}

// Error handling utility
export const handleApiError = (error: any): string => {
  if (error.response) {
    // The request was made and the server responded with a status code
    // that falls out of the range of 2xx
    return `Server Error: ${error.response.data?.message || error.response.statusText}`;
  } else if (error.request) {
    // The request was made but no response was received
    return 'Network Error: Unable to connect to the server. Please check if the backend is running.';
  } else {
    // Something happened in setting up the request that triggered an Error
    return `Error: ${error.message}`;
  }
};