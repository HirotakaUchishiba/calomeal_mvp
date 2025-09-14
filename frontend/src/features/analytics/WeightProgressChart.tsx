import React from 'react';
import { useQuery } from '@apollo/client';
import { GET_WEIGHT_PROGRESS_QUERY } from '../../graphql/queries';

interface WeightProgressChartProps {
  startDate: string;
  endDate: string;
}

interface WeightDataPoint {
  date: string;
  weight: number;
}

interface WeightProgress {
  startDate: string;
  endDate: string;
  startWeight: number;
  endWeight: number;
  weightChange: number;
  weightChangePercentage: number;
  avgWeeklyChange: number;
  trend: string;
  dataPoints: WeightDataPoint[];
}

const WeightProgressChart: React.FC<WeightProgressChartProps> = ({ startDate, endDate }) => {
  const { loading, error, data } = useQuery<{ weightProgress: WeightProgress }>(
    GET_WEIGHT_PROGRESS_QUERY,
    {
      variables: { startDate, endDate },
      errorPolicy: 'all'
    }
  );

  if (loading) return <div className="p-4 bg-white rounded-lg shadow">読み込み中...</div>;
  if (error) return <div className="p-4 bg-white rounded-lg shadow text-red-500">エラー: {error.message}</div>;
  if (!data?.weightProgress) return <div className="p-4 bg-white rounded-lg shadow">データがありません</div>;

  const progress = data.weightProgress;

  const getTrendColor = (trend: string) => {
    switch (trend.toLowerCase()) {
      case 'increasing': return 'text-red-600';
      case 'decreasing': return 'text-green-600';
      case 'stable': return 'text-blue-600';
      default: return 'text-gray-600';
    }
  };

  const getTrendIcon = (trend: string) => {
    switch (trend.toLowerCase()) {
      case 'increasing': return '↗️';
      case 'decreasing': return '↘️';
      case 'stable': return '→';
      default: return '📊';
    }
  };

  return (
    <div className="bg-white rounded-lg shadow-md p-6">
      <h3 className="text-lg font-semibold mb-4 text-gray-800">
        体重進捗 ({startDate} - {endDate})
      </h3>

      {/* 進捗サマリー */}
      <div className="grid grid-cols-2 gap-4 mb-6">
        <div className="bg-blue-50 p-4 rounded-lg">
          <div className="text-sm text-gray-600">開始体重</div>
          <div className="text-2xl font-bold text-blue-600">{progress.startWeight.toFixed(1)}kg</div>
        </div>
        <div className="bg-green-50 p-4 rounded-lg">
          <div className="text-sm text-gray-600">現在体重</div>
          <div className="text-2xl font-bold text-green-600">{progress.endWeight.toFixed(1)}kg</div>
        </div>
        <div className="bg-purple-50 p-4 rounded-lg">
          <div className="text-sm text-gray-600">体重変化</div>
          <div className={`text-2xl font-bold ${progress.weightChange >= 0 ? 'text-red-600' : 'text-green-600'}`}>
            {progress.weightChange >= 0 ? '+' : ''}{progress.weightChange.toFixed(1)}kg
          </div>
        </div>
        <div className="bg-yellow-50 p-4 rounded-lg">
          <div className="text-sm text-gray-600">変化率</div>
          <div className={`text-2xl font-bold ${progress.weightChangePercentage >= 0 ? 'text-red-600' : 'text-green-600'}`}>
            {progress.weightChangePercentage >= 0 ? '+' : ''}{progress.weightChangePercentage.toFixed(1)}%
          </div>
        </div>
      </div>

      {/* トレンド情報 */}
      <div className="bg-gray-50 p-4 rounded-lg mb-6">
        <div className="flex items-center justify-between">
          <div>
            <div className="text-sm text-gray-600">トレンド</div>
            <div className={`text-lg font-semibold ${getTrendColor(progress.trend)}`}>
              {getTrendIcon(progress.trend)} {progress.trend}
            </div>
          </div>
          <div>
            <div className="text-sm text-gray-600">週平均変化</div>
            <div className={`text-lg font-semibold ${progress.avgWeeklyChange >= 0 ? 'text-red-600' : 'text-green-600'}`}>
              {progress.avgWeeklyChange >= 0 ? '+' : ''}{progress.avgWeeklyChange.toFixed(2)}kg/週
            </div>
          </div>
        </div>
      </div>

      {/* データポイント */}
      <div>
        <h4 className="text-md font-semibold mb-3 text-gray-700">体重記録</h4>
        <div className="space-y-2 max-h-64 overflow-y-auto">
          {progress.dataPoints.map((point, index) => (
            <div key={index} className="bg-gray-50 p-3 rounded flex justify-between items-center">
              <span className="font-medium text-gray-800">{point.date}</span>
              <span className="text-lg font-semibold text-gray-700">{point.weight.toFixed(1)}kg</span>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default WeightProgressChart;
