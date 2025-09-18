import React, { useState } from 'react';
import {
  Box,
  Container,
  Typography,
  AppBar,
  Toolbar,
  Button,
  Alert,
  CircularProgress,
  Tabs,
  Tab,
  Paper,
  Snackbar,
} from '@mui/material';
import { Search, Analytics as AnalyticsIcon, Refresh } from '@mui/icons-material';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';

import SearchForm from './components/SearchForm';
import JobCard from './components/JobCard';
import Analytics from './components/Analytics';
import { JobService, handleApiError } from './services/jobService';
import { SearchFilters, SearchResponse, Job } from './types';

const theme = createTheme({
  palette: {
    mode: 'light',
    primary: {
      main: '#1976d2',
    },
    secondary: {
      main: '#dc004e',
    },
  },
});

interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

function TabPanel({ children, value, index }: TabPanelProps) {
  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`simple-tabpanel-${index}`}
      aria-labelledby={`simple-tab-${index}`}
    >
      {value === index && <Box sx={{ p: 3 }}>{children}</Box>}
    </div>
  );
}

function App() {
  const [searchResponse, setSearchResponse] = useState<SearchResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [tabValue, setTabValue] = useState(0);
  const [snackbarOpen, setSnackbarOpen] = useState(false);
  const [snackbarMessage, setSnackbarMessage] = useState('');

  const handleSearch = async (filters: SearchFilters) => {
    setLoading(true);
    setError(null);
    
    try {
      const response = await JobService.searchJobs(filters);
      setSearchResponse(response);
      setSnackbarMessage(`Found ${response.total} jobs matching your criteria`);
      setSnackbarOpen(true);
    } catch (err) {
      const errorMessage = handleApiError(err);
      setError(errorMessage);
      console.error('Search error:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleClear = () => {
    setSearchResponse(null);
    setError(null);
  };

  const handleClearCache = async () => {
    try {
      await JobService.clearCache();
      setSnackbarMessage('Cache cleared successfully');
      setSnackbarOpen(true);
      if (searchResponse) {
        setSearchResponse(null);
      }
    } catch (err) {
      const errorMessage = handleApiError(err);
      setError(errorMessage);
    }
  };

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue);
  };

  const handleViewJobDetails = (job: Job) => {
    // This could open a modal or navigate to a detail page
    window.open(job.url, '_blank');
  };

  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Box sx={{ flexGrow: 1 }}>
        {/* Header */}
        <AppBar position="static">
          <Toolbar>
            <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
              Job Market Scraper
            </Typography>
            <Button
              color="inherit"
              startIcon={<Refresh />}
              onClick={handleClearCache}
              sx={{ mr: 1 }}
            >
              Clear Cache
            </Button>
          </Toolbar>
        </AppBar>

        {/* Main Content */}
        <Container maxWidth="xl" sx={{ mt: 3, mb: 3 }}>
          {/* Search Form */}
          <SearchForm
            onSearch={handleSearch}
            loading={loading}
            onClear={handleClear}
          />

          {/* Error Display */}
          {error && (
            <Alert severity="error" sx={{ mb: 3 }} onClose={() => setError(null)}>
              {error}
            </Alert>
          )}

          {/* Loading Indicator */}
          {loading && (
            <Box sx={{ display: 'flex', justifyContent: 'center', my: 4 }}>
              <CircularProgress size={60} />
              <Box sx={{ ml: 2, display: 'flex', alignItems: 'center' }}>
                <Typography variant="h6">
                  Searching jobs across multiple sites...
                </Typography>
              </Box>
            </Box>
          )}

          {/* Results */}
          {searchResponse && !loading && (
            <Paper sx={{ width: '100%' }}>
              <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
                <Tabs value={tabValue} onChange={handleTabChange} aria-label="job results tabs">
                  <Tab
                    icon={<Search />}
                    label={`Jobs (${searchResponse.total})`}
                    id="simple-tab-0"
                    aria-controls="simple-tabpanel-0"
                  />
                  <Tab
                    icon={<AnalyticsIcon />}
                    label="Analytics"
                    id="simple-tab-1"
                    aria-controls="simple-tabpanel-1"
                  />
                </Tabs>
              </Box>

              {/* Jobs Tab */}
              <TabPanel value={tabValue} index={0}>
                {searchResponse.jobs.length === 0 ? (
                  <Box sx={{ textAlign: 'center', py: 4 }}>
                    <Typography variant="h6" color="text.secondary">
                      No jobs found matching your criteria.
                    </Typography>
                    <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
                      Try adjusting your search filters and search again.
                    </Typography>
                  </Box>
                ) : (
                  <Box>
                    <Typography variant="h6" sx={{ mb: 2 }}>
                      Found {searchResponse.total} jobs
                    </Typography>
                    {searchResponse.jobs.map((job) => (
                      <JobCard
                        key={job.id}
                        job={job}
                        onViewDetails={handleViewJobDetails}
                      />
                    ))}
                  </Box>
                )}
              </TabPanel>

              {/* Analytics Tab */}
              <TabPanel value={tabValue} index={1}>
                <Analytics analytics={searchResponse.analytics} />
              </TabPanel>
            </Paper>
          )}
        </Container>

        {/* Snackbar for notifications */}
        <Snackbar
          open={snackbarOpen}
          autoHideDuration={4000}
          onClose={() => setSnackbarOpen(false)}
          message={snackbarMessage}
        />
      </Box>
    </ThemeProvider>
  );
}

export default App;