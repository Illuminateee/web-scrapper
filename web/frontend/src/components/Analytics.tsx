import React from 'react';
import {
  Box,
  Card,
  CardContent,
  Grid,
  Typography,
  Chip,
  List,
  ListItem,
  ListItemText,
  LinearProgress,
} from '@mui/material';
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  PieChart,
  Pie,
  Cell,
  ResponsiveContainer,
} from 'recharts';
import { JobAnalytics } from '../types';

interface AnalyticsProps {
  analytics: JobAnalytics;
}

const COLORS = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042', '#8884D8', '#82CA9D'];

const Analytics: React.FC<AnalyticsProps> = ({ analytics }) => {
  // Prepare data for charts
  const experienceData = Object.entries(analytics.experience_levels).map(([level, count]) => ({
    name: level.charAt(0).toUpperCase() + level.slice(1),
    count,
  }));

  const remoteData = Object.entries(analytics.remote_options).map(([option, count]) => ({
    name: option.charAt(0).toUpperCase() + option.slice(1),
    count,
  }));

  const degreeData = Object.entries(analytics.degree_requirements).map(([req, count]) => ({
    name: req,
    count,
  }));

  const locationData = Object.entries(analytics.location_distribution)
    .sort(([, a], [, b]) => b - a)
    .slice(0, 10)
    .map(([location, count]) => ({
      name: location.length > 20 ? location.substring(0, 20) + '...' : location,
      count,
    }));

  const formatSalary = (salary: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      maximumFractionDigits: 0,
    }).format(salary);
  };

  return (
    <Box sx={{ mt: 3 }}>
      <Typography variant="h5" component="h2" gutterBottom>
        Job Market Analytics
      </Typography>

      <Grid container spacing={3}>
        {/* Summary Cards */}
        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Typography color="textSecondary" gutterBottom>
                Total Jobs
              </Typography>
              <Typography variant="h4">
                {analytics.total_jobs.toLocaleString()}
              </Typography>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Typography color="textSecondary" gutterBottom>
                Average Salary
              </Typography>
              <Typography variant="h4">
                {formatSalary(analytics.average_salary)}
              </Typography>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Typography color="textSecondary" gutterBottom>
                Salary Range
              </Typography>
              <Typography variant="h6">
                {formatSalary(analytics.salary_range.min)} - {formatSalary(analytics.salary_range.max)}
              </Typography>
              <Typography variant="body2" color="textSecondary">
                Median: {formatSalary(analytics.salary_range.median)}
              </Typography>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Typography color="textSecondary" gutterBottom>
                Top Skills Found
              </Typography>
              <Typography variant="h4">
                {analytics.top_skills.length}
              </Typography>
            </CardContent>
          </Card>
        </Grid>

        {/* Top Skills */}
        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Most In-Demand Skills
              </Typography>
              <List dense>
                {analytics.top_skills.slice(0, 10).map((skill, index) => (
                  <ListItem key={skill.skill} sx={{ px: 0 }}>
                    <ListItemText
                      primary={
                        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                          <Typography variant="body2">
                            {index + 1}. {skill.skill}
                          </Typography>
                          <Chip
                            label={skill.count}
                            size="small"
                            color="primary"
                            variant="outlined"
                          />
                        </Box>
                      }
                    />
                    <Box sx={{ width: '100%', mt: 1 }}>
                      <LinearProgress
                        variant="determinate"
                        value={(skill.count / analytics.top_skills[0].count) * 100}
                        sx={{ height: 4, borderRadius: 2 }}
                      />
                    </Box>
                  </ListItem>
                ))}
              </List>
            </CardContent>
          </Card>
        </Grid>

        {/* Top Companies */}
        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Top Hiring Companies
              </Typography>
              <List dense>
                {analytics.top_companies.slice(0, 10).map((company, index) => (
                  <ListItem key={company.company} sx={{ px: 0 }}>
                    <ListItemText
                      primary={
                        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                          <Typography variant="body2">
                            {index + 1}. {company.company}
                          </Typography>
                          <Chip
                            label={`${company.count} jobs`}
                            size="small"
                            color="secondary"
                            variant="outlined"
                          />
                        </Box>
                      }
                    />
                    <Box sx={{ width: '100%', mt: 1 }}>
                      <LinearProgress
                        variant="determinate"
                        value={(company.count / analytics.top_companies[0].count) * 100}
                        color="secondary"
                        sx={{ height: 4, borderRadius: 2 }}
                      />
                    </Box>
                  </ListItem>
                ))}
              </List>
            </CardContent>
          </Card>
        </Grid>

        {/* Experience Levels Chart */}
        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Experience Levels Distribution
              </Typography>
              <ResponsiveContainer width="100%" height={250}>
                <PieChart>
                  <Pie
                    data={experienceData}
                    cx="50%"
                    cy="50%"
                    labelLine={false}
                    label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
                    outerRadius={80}
                    fill="#8884d8"
                    dataKey="count"
                  >
                    {experienceData.map((entry, index) => (
                      <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                    ))}
                  </Pie>
                  <Tooltip />
                </PieChart>
              </ResponsiveContainer>
            </CardContent>
          </Card>
        </Grid>

        {/* Remote Options Chart */}
        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Work Arrangement Options
              </Typography>
              <ResponsiveContainer width="100%" height={250}>
                <BarChart data={remoteData}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="name" />
                  <YAxis />
                  <Tooltip />
                  <Bar dataKey="count" fill="#8884d8" />
                </BarChart>
              </ResponsiveContainer>
            </CardContent>
          </Card>
        </Grid>

        {/* Degree Requirements */}
        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Degree Requirements
              </Typography>
              <ResponsiveContainer width="100%" height={250}>
                <PieChart>
                  <Pie
                    data={degreeData}
                    cx="50%"
                    cy="50%"
                    labelLine={false}
                    label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
                    outerRadius={80}
                    fill="#8884d8"
                    dataKey="count"
                  >
                    {degreeData.map((entry, index) => (
                      <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                    ))}
                  </Pie>
                  <Tooltip />
                </PieChart>
              </ResponsiveContainer>
            </CardContent>
          </Card>
        </Grid>

        {/* Top Locations */}
        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Top Locations
              </Typography>
              <ResponsiveContainer width="100%" height={250}>
                <BarChart data={locationData} layout="horizontal">
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis type="number" />
                  <YAxis dataKey="name" type="category" width={100} />
                  <Tooltip />
                  <Bar dataKey="count" fill="#82ca9d" />
                </BarChart>
              </ResponsiveContainer>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Box>
  );
};

export default Analytics;