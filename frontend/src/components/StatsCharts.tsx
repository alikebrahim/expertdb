import { 
  PieChart, Pie, Cell, BarChart, Bar, XAxis, YAxis, 
  CartesianGrid, Tooltip, Legend, LineChart, Line, 
  ResponsiveContainer
} from 'recharts';
import { NationalityStats, GrowthStats, IscedStats } from '../types';

// Custom colors
const COLORS = ['#003366', '#0055a4', '#e63946', '#457b9d', '#1d3557', '#a8dadc'];

interface NationalityChartProps {
  data: NationalityStats[];
  isLoading: boolean;
}

export const NationalityChart = ({ data, isLoading }: NationalityChartProps) => {
  if (isLoading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin h-8 w-8 border-4 border-primary border-t-transparent rounded-full"></div>
      </div>
    );
  }
  
  if (!data || data.length === 0) {
    return (
      <div className="flex justify-center items-center h-64 bg-accent rounded">
        <p className="text-neutral-600">No nationality data available</p>
      </div>
    );
  }
  
  return (
    <div className="h-64">
      <ResponsiveContainer width="100%" height="100%">
        <PieChart>
          <Pie
            data={data}
            cx="50%"
            cy="50%"
            labelLine={false}
            outerRadius={80}
            fill="#8884d8"
            dataKey="count"
            nameKey="nationality"
            label={({ nationality, percent }) => `${nationality}: ${(percent * 100).toFixed(0)}%`}
          >
            {data.map((_, index) => (
              <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
            ))}
          </Pie>
          <Tooltip 
            formatter={(value) => [`${value} experts`, 'Count']}
          />
          <Legend />
        </PieChart>
      </ResponsiveContainer>
    </div>
  );
};

interface GrowthChartProps {
  data: GrowthStats[];
  isLoading: boolean;
}

export const GrowthChart = ({ data, isLoading }: GrowthChartProps) => {
  if (isLoading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin h-8 w-8 border-4 border-primary border-t-transparent rounded-full"></div>
      </div>
    );
  }
  
  if (!data || data.length === 0) {
    return (
      <div className="flex justify-center items-center h-64 bg-accent rounded">
        <p className="text-neutral-600">No growth data available</p>
      </div>
    );
  }
  
  return (
    <div className="h-64">
      <ResponsiveContainer width="100%" height="100%">
        <LineChart
          data={data}
          margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
        >
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis dataKey="year" />
          <YAxis />
          <Tooltip formatter={(value) => [`${value} experts`, 'Count']} />
          <Legend />
          <Line type="monotone" dataKey="count" stroke="#003366" activeDot={{ r: 8 }} />
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
};

interface IscedChartProps {
  data: IscedStats[];
  isLoading: boolean;
}

export const IscedChart = ({ data, isLoading }: IscedChartProps) => {
  if (isLoading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin h-8 w-8 border-4 border-primary border-t-transparent rounded-full"></div>
      </div>
    );
  }
  
  if (!data || data.length === 0) {
    return (
      <div className="flex justify-center items-center h-64 bg-accent rounded">
        <p className="text-neutral-600">No ISCED data available</p>
      </div>
    );
  }
  
  // Sort by count in descending order
  const sortedData = [...data].sort((a, b) => b.count - a.count);
  
  return (
    <div className="h-64">
      <ResponsiveContainer width="100%" height="100%">
        <BarChart
          data={sortedData}
          layout="vertical"
          margin={{ top: 5, right: 30, left: 150, bottom: 5 }}
        >
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis type="number" />
          <YAxis type="category" dataKey="isced" width={120} />
          <Tooltip formatter={(value) => [`${value} experts`, 'Count']} />
          <Legend />
          <Bar dataKey="count" fill="#003366" />
        </BarChart>
      </ResponsiveContainer>
    </div>
  );
};

export default {
  NationalityChart,
  GrowthChart,
  IscedChart
};