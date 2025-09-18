import React, { useState } from 'react';
import {
  Box,
  Button,
  Card,
  CardContent,
  Chip,
  FormControl,
  FormControlLabel,
  Grid,
  InputLabel,
  MenuItem,
  Select,
  Switch,
  TextField,
  Typography,
  Autocomplete,
  Slider,
} from '@mui/material';
import { Search, Clear } from '@mui/icons-material';
import {
  SearchFilters,
  EXPERIENCE_LEVELS,
  DEGREE_OPTIONS,
  COMMON_SKILLS,
  JOB_TITLE_SUGGESTIONS,
} from '../types';

interface SearchFormProps {
  onSearch: (filters: SearchFilters) => void;
  loading: boolean;
  onClear: () => void;
}

const SearchForm: React.FC<SearchFormProps> = ({ onSearch, loading, onClear }) => {
  const [filters, setFilters] = useState<SearchFilters>({
    limit: 50,
    offset: 0,
  });

  const [salaryRange, setSalaryRange] = useState<number[]>([0, 200000]);
  const [selectedSkills, setSelectedSkills] = useState<string[]>([]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    const searchFilters: SearchFilters = {
      ...filters,
      min_salary: salaryRange[0] > 0 ? salaryRange[0] : undefined,
      max_salary: salaryRange[1] < 200000 ? salaryRange[1] : undefined,
      skills: selectedSkills.length > 0 ? selectedSkills : undefined,
    };

    onSearch(searchFilters);
  };

  const handleClear = () => {
    setFilters({ limit: 50, offset: 0 });
    setSalaryRange([0, 200000]);
    setSelectedSkills([]);
    onClear();
  };

  const handleInputChange = (field: keyof SearchFilters) => (
    event: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement> | any
  ) => {
    const value = event.target.value;
    setFilters(prev => ({
      ...prev,
      [field]: value === '' ? undefined : value,
    }));
  };

  const handleSwitchChange = (field: keyof SearchFilters) => (
    event: React.ChangeEvent<HTMLInputElement>
  ) => {
    setFilters(prev => ({
      ...prev,
      [field]: event.target.checked,
    }));
  };

  const handleDegreeRequiredChange = (event: any) => {
    const value = event.target.value;
    setFilters(prev => ({
      ...prev,
      degree_required: value === '' ? undefined : value === 'true',
    }));
  };

  const salaryMarks = [
    { value: 0, label: '$0' },
    { value: 50000, label: '$50K' },
    { value: 100000, label: '$100K' },
    { value: 150000, label: '$150K' },
    { value: 200000, label: '$200K+' },
  ];

  return (
    <Card sx={{ mb: 3 }}>
      <CardContent>
        <Typography variant="h5" component="h2" gutterBottom>
          Search Jobs
        </Typography>
        
        <Box component="form" onSubmit={handleSubmit}>
          <Grid container spacing={3}>
            {/* Job Title */}
            <Grid item xs={12} md={6}>
              <Autocomplete
                freeSolo
                options={JOB_TITLE_SUGGESTIONS}
                value={filters.job_title || ''}
                onInputChange={(_, value) => {
                  setFilters(prev => ({ ...prev, job_title: value || undefined }));
                }}
                renderInput={(params) => (
                  <TextField
                    {...params}
                    label="Job Title"
                    placeholder="e.g., Backend Developer, Go Developer"
                    fullWidth
                  />
                )}
              />
            </Grid>

            {/* Location */}
            <Grid item xs={12} md={6}>
              <TextField
                label="Location"
                value={filters.location || ''}
                onChange={handleInputChange('location')}
                placeholder="e.g., San Francisco, Remote"
                fullWidth
              />
            </Grid>

            {/* Skills */}
            <Grid item xs={12}>
              <Autocomplete
                multiple
                options={COMMON_SKILLS}
                value={selectedSkills}
                onChange={(_, value) => setSelectedSkills(value)}
                renderTags={(value, getTagProps) =>
                  value.map((option, index) => (
                    <Chip
                      variant="outlined"
                      label={option}
                      {...getTagProps({ index })}
                      key={option}
                    />
                  ))
                }
                renderInput={(params) => (
                  <TextField
                    {...params}
                    label="Skills"
                    placeholder="Select skills (e.g., Go, React, AWS)"
                  />
                )}
              />
            </Grid>

            {/* Experience Level */}
            <Grid item xs={12} md={6}>
              <FormControl fullWidth>
                <InputLabel>Experience Level</InputLabel>
                <Select
                  value={filters.experience_level || ''}
                  onChange={handleInputChange('experience_level')}
                  label="Experience Level"
                >
                  {EXPERIENCE_LEVELS.map(option => (
                    <MenuItem key={option.value} value={option.value}>
                      {option.label}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>

            {/* Degree Requirement */}
            <Grid item xs={12} md={6}>
              <FormControl fullWidth>
                <InputLabel>Degree Requirement</InputLabel>
                <Select
                  value={
                    filters.degree_required === undefined 
                      ? '' 
                      : filters.degree_required.toString()
                  }
                  onChange={handleDegreeRequiredChange}
                  label="Degree Requirement"
                >
                  {DEGREE_OPTIONS.map(option => (
                    <MenuItem key={option.value} value={option.value}>
                      {option.label}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>

            {/* Salary Range */}
            <Grid item xs={12}>
              <Typography gutterBottom>
                Salary Range: ${salaryRange[0].toLocaleString()} - ${salaryRange[1] >= 200000 ? '200K+' : salaryRange[1].toLocaleString()}
              </Typography>
              <Slider
                value={salaryRange}
                onChange={(_, value) => setSalaryRange(value as number[])}
                valueLabelDisplay="auto"
                valueLabelFormat={(value) => `$${value.toLocaleString()}`}
                min={0}
                max={200000}
                step={5000}
                marks={salaryMarks}
                sx={{ mt: 2 }}
              />
            </Grid>

            {/* Remote Only Switch */}
            <Grid item xs={12}>
              <FormControlLabel
                control={
                  <Switch
                    checked={filters.remote_only || false}
                    onChange={handleSwitchChange('remote_only')}
                  />
                }
                label="Remote Only"
              />
            </Grid>

            {/* Action Buttons */}
            <Grid item xs={12}>
              <Box sx={{ display: 'flex', gap: 2, mt: 2 }}>
                <Button
                  type="submit"
                  variant="contained"
                  size="large"
                  startIcon={<Search />}
                  disabled={loading}
                  sx={{ minWidth: 140 }}
                >
                  {loading ? 'Searching...' : 'Search Jobs'}
                </Button>
                
                <Button
                  variant="outlined"
                  size="large"
                  startIcon={<Clear />}
                  onClick={handleClear}
                  disabled={loading}
                >
                  Clear Filters
                </Button>
              </Box>
            </Grid>
          </Grid>
        </Box>
      </CardContent>
    </Card>
  );
};

export default SearchForm;