'use client';

import { useEffect, useState } from 'react';
import { 
  Statistics, 
  NationalityStats, 
  AreaStat, 
  GrowthStat,
  statisticsAPI 
} from '@/lib/api';
import { Navbar } from '@/components/layout/navbar';
import RequireAuth from '@/components/auth/require-auth';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { 
  BarChart,
  PieChart,
  LineChart,
  Line,
  Bar,
  Pie,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
  Cell,
} from 'recharts';

// Define a color palette
const COLORS = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042', '#8884d8', '#82ca9d', '#ffc658'];
const BAR_COLORS = ['#3b82f6', '#0369a1', '#0284c7', '#0ea5e9', '#38bdf8', '#7dd3fc'];

export default function StatisticsPage() {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [statistics, setStatistics] = useState<Statistics | null>(null);
  const [nationalityStats, setNationalityStats] = useState<NationalityStats | null>(null);
  const [iscedStats, setIscedStats] = useState<AreaStat[] | null>(null);
  const [engagementStats, setEngagementStats] = useState<AreaStat[] | null>(null);
  const [growthStats, setGrowthStats] = useState<GrowthStat[] | null>(null);
  const [activeTab, setActiveTab] = useState('overview');

  useEffect(() => {
    const fetchAllData = async () => {
      try {
        setLoading(true);
        const [statsData, nationalityData, iscedData, engagementData, growthData] = await Promise.all([
          statisticsAPI.getAllStatistics(),
          statisticsAPI.getNationalityStats(),
          statisticsAPI.getISCEDStats(),
          statisticsAPI.getEngagementStats(),
          statisticsAPI.getGrowthStats(12) // Last 12 months
        ]);
        setStatistics(statsData);
        setNationalityStats(nationalityData);
        setIscedStats(iscedData);
        setEngagementStats(engagementData);
        setGrowthStats(growthData);
      } catch (err) {
        setError('Failed to load statistics. Please try again later.');
        console.error('Error fetching statistics:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchAllData();
  }, []);

  // Prepare nationality data for pie chart
  const nationalityData = nationalityStats
    ? [
        { name: 'Bahraini', value: nationalityStats.bahraini.count },
        { name: 'Non-Bahraini', value: nationalityStats.nonBahraini.count }
      ]
    : [];

  // Format growth data for time series chart
  const formattedGrowthData = growthStats?.map(stat => ({
    month: stat.period,
    experts: stat.count,
    growth: stat.growthRate
  }));

  // Prepare ISCED field data for bar chart
  const iscedChartData = iscedStats
    ? iscedStats.slice(0, 8).map(stat => ({
        name: stat.name,
        count: stat.count,
        percentage: stat.percentage
      }))
    : [];

  if (loading) {
    return (
      <div className="container p-4 mx-auto">
        <div className="flex justify-center items-center min-h-[400px]">
          <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary"></div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="container p-4 mx-auto">
        <div className="bg-red-50 text-red-500 p-4 rounded-md">
          <h2 className="text-lg font-bold">Error</h2>
          <p>{error}</p>
        </div>
      </div>
    );
  }

  return (
    <RequireAuth>
      <>
        <Navbar />
        <div className="container p-4 mx-auto">
          <div className="mb-6">
            <h1 className="text-3xl font-bold">Statistics Dashboard</h1>
            <p className="text-muted-foreground">
              Overview of expert database metrics and trends
            </p>
          </div>

      <Tabs defaultValue="overview" className="space-y-6" onValueChange={setActiveTab}>
        <TabsList>
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="nationality">Nationality</TabsTrigger>
          <TabsTrigger value="isced">ISCED Fields</TabsTrigger>
          <TabsTrigger value="growth">Growth Trends</TabsTrigger>
          <TabsTrigger value="engagements">Engagements</TabsTrigger>
        </TabsList>

        {/* Overview */}
        <TabsContent value="overview" className="space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            <StatCard 
              title="Total Experts" 
              value={statistics?.totalExperts || 0} 
            />
            <StatCard 
              title="Bahraini Experts" 
              value={nationalityStats?.bahraini.count || 0} 
              detail={`${nationalityStats?.bahraini.percentage.toFixed(1) || 0}%`}
            />
            <StatCard 
              title="Non-Bahraini Experts" 
              value={nationalityStats?.nonBahraini.count || 0} 
              detail={`${nationalityStats?.nonBahraini.percentage.toFixed(1) || 0}%`}
            />
            <StatCard 
              title="Available Experts" 
              value={statistics?.totalExperts || 0} 
              detail="Ready for engagements"
            />
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <Card>
              <CardHeader>
                <CardTitle>Nationality Distribution</CardTitle>
              </CardHeader>
              <CardContent className="h-[300px]">
                {nationalityData.length > 0 && (
                  <ResponsiveContainer width="100%" height="100%">
                    <PieChart>
                      <Pie
                        data={nationalityData}
                        cx="50%"
                        cy="50%"
                        innerRadius={60}
                        outerRadius={100}
                        fill="#8884d8"
                        paddingAngle={5}
                        dataKey="value"
                        label={({ name, percent }) => `${name}: ${(percent * 100).toFixed(1)}%`}
                      >
                        {nationalityData.map((_, index) => (
                          <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                        ))}
                      </Pie>
                      <Tooltip />
                      <Legend />
                    </PieChart>
                  </ResponsiveContainer>
                )}
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Top ISCED Fields</CardTitle>
              </CardHeader>
              <CardContent className="h-[300px]">
                {iscedChartData.length > 0 && (
                  <ResponsiveContainer width="100%" height="100%">
                    <BarChart
                      data={iscedChartData.slice(0, 5)}
                      layout="vertical"
                      margin={{ top: 0, right: 0, left: 70, bottom: 0 }}
                    >
                      <CartesianGrid strokeDasharray="3 3" />
                      <XAxis type="number" />
                      <YAxis dataKey="name" type="category" width={70} />
                      <Tooltip
                        formatter={(value, name) => [`${value} experts`, name]}
                      />
                      <Bar dataKey="count" fill="#3b82f6">
                        {iscedChartData.map((_, index) => (
                          <Cell key={`cell-${index}`} fill={BAR_COLORS[index % BAR_COLORS.length]} />
                        ))}
                      </Bar>
                    </BarChart>
                  </ResponsiveContainer>
                )}
              </CardContent>
            </Card>
          </div>

          <Card>
            <CardHeader>
              <CardTitle>Expert Growth Trend</CardTitle>
            </CardHeader>
            <CardContent className="h-[300px]">
              {formattedGrowthData && formattedGrowthData.length > 0 && (
                <ResponsiveContainer width="100%" height="100%">
                  <LineChart
                    data={formattedGrowthData}
                    margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
                  >
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="month" />
                    <YAxis />
                    <Tooltip />
                    <Legend />
                    <Line type="monotone" dataKey="experts" stroke="#3b82f6" name="Expert Count" />
                  </LineChart>
                </ResponsiveContainer>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        {/* Nationality Tab */}
        <TabsContent value="nationality" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Nationality Distribution</CardTitle>
            </CardHeader>
            <CardContent className="h-[400px]">
              {nationalityData.length > 0 && (
                <ResponsiveContainer width="100%" height="100%">
                  <PieChart>
                    <Pie
                      data={nationalityData}
                      cx="50%"
                      cy="50%"
                      innerRadius={80}
                      outerRadius={150}
                      fill="#8884d8"
                      paddingAngle={5}
                      dataKey="value"
                      label={({ name, percent }) => `${name}: ${(percent * 100).toFixed(1)}%`}
                    >
                      {nationalityData.map((_, index) => (
                        <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                      ))}
                    </Pie>
                    <Tooltip />
                    <Legend />
                  </PieChart>
                </ResponsiveContainer>
              )}
            </CardContent>
          </Card>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <Card>
              <CardHeader>
                <CardTitle>Bahraini Experts</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="text-4xl font-bold">{nationalityStats?.bahraini.count || 0}</div>
                <p className="text-sm text-muted-foreground mt-2">
                  {nationalityStats?.bahraini.percentage.toFixed(1) || 0}% of total experts
                </p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Non-Bahraini Experts</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="text-4xl font-bold">{nationalityStats?.nonBahraini.count || 0}</div>
                <p className="text-sm text-muted-foreground mt-2">
                  {nationalityStats?.nonBahraini.percentage.toFixed(1) || 0}% of total experts
                </p>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        {/* ISCED Fields Tab */}
        <TabsContent value="isced" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Expert Distribution by ISCED Field</CardTitle>
            </CardHeader>
            <CardContent className="h-[450px]">
              {iscedChartData.length > 0 && (
                <ResponsiveContainer width="100%" height="100%">
                  <BarChart
                    data={iscedChartData}
                    layout="vertical"
                    margin={{ top: 5, right: 30, left: 120, bottom: 5 }}
                  >
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis type="number" />
                    <YAxis dataKey="name" type="category" width={120} />
                    <Tooltip
                      formatter={(value, name, props) => [`${value} experts (${props.payload.percentage.toFixed(1)}%)`, 'Count']}
                    />
                    <Bar dataKey="count" fill="#3b82f6">
                      {iscedChartData.map((_, index) => (
                        <Cell key={`cell-${index}`} fill={BAR_COLORS[index % BAR_COLORS.length]} />
                      ))}
                    </Bar>
                  </BarChart>
                </ResponsiveContainer>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        {/* Growth Trends Tab */}
        <TabsContent value="growth" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Expert Growth Trend (Last 12 Months)</CardTitle>
            </CardHeader>
            <CardContent className="h-[400px]">
              {formattedGrowthData && formattedGrowthData.length > 0 && (
                <ResponsiveContainer width="100%" height="100%">
                  <LineChart
                    data={formattedGrowthData}
                    margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
                  >
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="month" />
                    <YAxis yAxisId="left" orientation="left" />
                    <YAxis yAxisId="right" orientation="right" />
                    <Tooltip />
                    <Legend />
                    <Line yAxisId="left" type="monotone" dataKey="experts" stroke="#3b82f6" name="Expert Count" />
                    <Line yAxisId="right" type="monotone" dataKey="growth" stroke="#10b981" name="Growth Rate (%)" />
                  </LineChart>
                </ResponsiveContainer>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        {/* Engagements Tab */}
        <TabsContent value="engagements" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Engagements by Type</CardTitle>
            </CardHeader>
            <CardContent className="h-[400px]">
              {engagementStats && engagementStats.length > 0 && (
                <ResponsiveContainer width="100%" height="100%">
                  <PieChart>
                    <Pie
                      data={engagementStats}
                      cx="50%"
                      cy="50%"
                      outerRadius={140}
                      fill="#8884d8"
                      dataKey="count"
                      label={({ name, percent }) => `${name}: ${(percent * 100).toFixed(1)}%`}
                    >
                      {engagementStats.map((_, index) => (
                        <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                      ))}
                    </Pie>
                    <Tooltip formatter={(value, name, props) => [`${value} (${props.payload.percentage.toFixed(1)}%)`, props.payload.name]} />
                    <Legend />
                  </PieChart>
                </ResponsiveContainer>
              )}
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
    </>
    </RequireAuth>
  );
}

// Stat Card Component
function StatCard({ title, value, detail }: { title: string; value: number; detail?: string }) {
  return (
    <Card>
      <CardHeader className="pb-2">
        <CardTitle className="text-sm font-medium">{title}</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="text-2xl font-bold">{value}</div>
        {detail && <p className="text-xs text-muted-foreground">{detail}</p>}
      </CardContent>
    </Card>
  );
}