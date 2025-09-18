import React from 'react';
import {
  Box,
  Card,
  CardContent,
  Chip,
  Typography,
  Button,
  Grid,
  Divider,
  IconButton,
  Tooltip,
} from '@mui/material';
import {
  LocationOn,
  Business,
  Schedule,
  School,
  Work,
  AttachMoney,
  OpenInNew,
  Star,
} from '@mui/icons-material';
import { Job } from '../types';

interface JobCardProps {
  job: Job;
  onViewDetails?: (job: Job) => void;
}

const JobCard: React.FC<JobCardProps> = ({ job, onViewDetails }) => {
  const formatSalary = (min?: number, max?: number, currency?: string) => {
    if (!min && !max) return 'Salary not specified';
    
    const curr = currency || 'USD';
    const formatter = new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: curr,
      maximumFractionDigits: 0,
    });
    
    if (min && max) {
      return `${formatter.format(min)} - ${formatter.format(max)}`;
    } else if (min) {
      return `${formatter.format(min)}+`;
    } else if (max) {
      return `Up to ${formatter.format(max)}`;
    }
    
    return 'Salary not specified';
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    const now = new Date();
    const diffTime = Math.abs(now.getTime() - date.getTime());
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
    
    if (diffDays === 1) return '1 day ago';
    if (diffDays < 7) return `${diffDays} days ago`;
    if (diffDays < 30) return `${Math.ceil(diffDays / 7)} weeks ago`;
    return `${Math.ceil(diffDays / 30)} months ago`;
  };

  const getExperienceLevelColor = (level: string) => {
    switch (level.toLowerCase()) {
      case 'entry': return 'success';
      case 'mid': return 'primary';
      case 'senior': return 'warning';
      case 'lead': return 'error';
      default: return 'default';
    }
  };

  const getRemoteOptionColor = (option: string) => {
    switch (option.toLowerCase()) {
      case 'remote': return 'success';
      case 'hybrid': return 'warning';
      case 'onsite': return 'default';
      default: return 'default';
    }
  };

  return (
    <Card sx={{ mb: 2, '&:hover': { boxShadow: 4 } }}>
      <CardContent>
        {/* Header */}
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', mb: 2 }}>
          <Box sx={{ flex: 1 }}>
            <Typography variant="h6" component="h3" sx={{ mb: 1, color: 'primary.main' }}>
              {job.title}
            </Typography>
            
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 1 }}>
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
                <Business fontSize="small" color="action" />
                <Typography variant="body2" color="text.secondary">
                  {job.company}
                </Typography>
              </Box>
              
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
                <LocationOn fontSize="small" color="action" />
                <Typography variant="body2" color="text.secondary">
                  {job.location}
                </Typography>
              </Box>
            </Box>
          </Box>
          
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <Tooltip title="View Original Posting">
              <IconButton
                href={job.url}
                target="_blank"
                rel="noopener noreferrer"
                size="small"
              >
                <OpenInNew />
              </IconButton>
            </Tooltip>
            
            <Chip
              label={job.source}
              size="small"
              variant="outlined"
            />
          </Box>
        </Box>

        {/* Key Info */}
        <Grid container spacing={2} sx={{ mb: 2 }}>
          <Grid item xs={12} sm={6} md={3}>
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
              <AttachMoney fontSize="small" color="action" />
              <Typography variant="body2">
                {formatSalary(job.salary_min, job.salary_max, job.salary_currency)}
              </Typography>
            </Box>
          </Grid>
          
          <Grid item xs={12} sm={6} md={3}>
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
              <Work fontSize="small" color="action" />
              <Chip
                label={job.experience_level}
                size="small"
                color={getExperienceLevelColor(job.experience_level)}
                variant="outlined"
              />
            </Box>
          </Grid>
          
          <Grid item xs={12} sm={6} md={3}>
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
              <LocationOn fontSize="small" color="action" />
              <Chip
                label={job.remote_option}
                size="small"
                color={getRemoteOptionColor(job.remote_option)}
                variant="outlined"
              />
            </Box>
          </Grid>
          
          <Grid item xs={12} sm={6} md={3}>
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
              <School fontSize="small" color="action" />
              <Chip
                label={job.degree_required ? 'Degree Required' : 'No Degree Required'}
                size="small"
                color={job.degree_required ? 'warning' : 'success'}
                variant="outlined"
              />
            </Box>
          </Grid>
        </Grid>

        {/* Skills */}
        {job.skills && job.skills.length > 0 && (
          <Box sx={{ mb: 2 }}>
            <Typography variant="subtitle2" sx={{ mb: 1 }}>
              Skills:
            </Typography>
            <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
              {job.skills.slice(0, 8).map((skill, index) => (
                <Chip
                  key={index}
                  label={skill}
                  size="small"
                  variant="outlined"
                  color="primary"
                />
              ))}
              {job.skills.length > 8 && (
                <Chip
                  label={`+${job.skills.length - 8} more`}
                  size="small"
                  variant="outlined"
                />
              )}
            </Box>
          </Box>
        )}

        {/* Description */}
        <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
          {job.description.length > 200 
            ? `${job.description.substring(0, 200)}...` 
            : job.description
          }
        </Typography>

        <Divider sx={{ my: 2 }} />

        {/* Footer */}
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
            <Schedule fontSize="small" color="action" />
            <Typography variant="caption" color="text.secondary">
              Posted {formatDate(job.posted_date)}
            </Typography>
          </Box>
          
          {onViewDetails && (
            <Button
              variant="outlined"
              size="small"
              onClick={() => onViewDetails(job)}
            >
              View Details
            </Button>
          )}
        </Box>
      </CardContent>
    </Card>
  );
};

export default JobCard;