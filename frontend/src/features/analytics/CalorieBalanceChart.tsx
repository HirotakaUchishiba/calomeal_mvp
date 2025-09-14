import React from 'react';
import { useQuery } from '@apollo/client';
import { GET_CALORIE_BALANCE_QUERY } from '../../graphql/queries';

interface CalorieBalanceChartProps {
  startDate: string;
  endDate: string;
}

interface DailyBalance {
  date: string;
  caloriesIntake: number;
  caloriesBurned: number;
  balance: number;
}

interface CalorieBalance {
  startDate: string;
  endDate: string;
  totalCaloriesIntake: number;
  totalCaloriesBurned: number;
  totalCalorieBalance: number;
  avgDailyBalance: number;
  daysInDeficit: number;
  daysInSurplus: number;
  deficitPercentage: number;
  dailyBalances: DailyBalance[];
}

const CalorieBalanceChart: React.FC<CalorieBalanceChartProps> = ({ startDate, endDate }) => {
  const { loading, error, data } = useQuery<{ calorieBalance: CalorieBalance }>(
    GET_CALORIE_BALANCE_QUERY,
    {
      variables: { startDate, endDate },
      errorPolicy: 'all'
    }
  );

  if (loading) return <div className="p-4 bg-white rounded-lg shadow">読み込み中...</div>;
  if (error) return <div className="p-4 bg-white rounded-lg shadow text-red-500">エラー: {error.message}</div>;
  if (!data?.calorieBalance) return <div className="p-4 bg-white rounded-lg shadow">データがありません</div>;

  const balance = data.calorieBalance;

  return (
    <div className="bg-white rounded-lg shadow-md p-6">
      <h3 className="text-lg font-semibold mb-4 text-gray-800">
        カロリーバランス ({startDate} - {endDate})
      </h3>

      {/* サマリー統計 */}
      <div className="grid grid-cols-2 gap-4 mb-6">
        <div className="bg-blue-50 p-4 rounded-lg">
          <div className="text-sm text-gray-600">総摂取カロリー</div>
          <div className="text-2xl font-bold text-blue-600">{balance.totalCaloriesIntake.toLocaleString()}</div>
        </div>
        <div className="bg-green-50 p-4 rounded-lg">
          <div className="text-sm text-gray-600">総消費カロリー</div>
          <div className="text-2xl font-bold text-green-600">{balance.totalCaloriesBurned.toLocaleString()}</div>
        </div>
        <div className="bg-purple-50 p-4 rounded-lg">
          <div className="text-sm text-gray-600">総バランス</div>
          <div className={`text-2xl font-bold ${balance.totalCalorieBalance >= 0 ? 'text-red-600' : 'text-green-600'}`}>
            {balance.totalCalorieBalance >= 0 ? '+' : ''}{balance.totalCalorieBalance.toLocaleString()}
          </div>
        </div>
        <div className="bg-yellow-50 p-4 rounded-lg">
          <div className="text-sm text-gray-600">平均日次バランス</div>
          <div className={`text-2xl font-bold ${balance.avgDailyBalance >= 0 ? 'text-red-600' : 'text-green-600'}`}>
            {balance.avgDailyBalance >= 0 ? '+' : ''}{balance.avgDailyBalance.toFixed(0)}
          </div>
        </div>
      </div>

      {/* 日数統計 */}
      <div className="grid grid-cols-3 gap-4 mb-6">
        <div className="bg-green-50 p-4 rounded-lg text-center">
          <div className="text-2xl font-bold text-green-600">{balance.daysInDeficit}</div>
          <div className="text-sm text-gray-600">赤字日数</div>
        </div>
        <div className="bg-red-50 p-4 rounded-lg text-center">
          <div className="text-2xl font-bold text-red-600">{balance.daysInSurplus}</div>
          <div className="text-sm text-gray-600">黒字日数</div>
        </div>
        <div className="bg-gray-50 p-4 rounded-lg text-center">
          <div className="text-2xl font-bold text-gray-600">{balance.deficitPercentage.toFixed(1)}%</div>
          <div className="text-sm text-gray-600">赤字率</div>
        </div>
      </div>

      {/* 日次バランス */}
      <div>
        <h4 className="text-md font-semibold mb-3 text-gray-700">日次バランス</h4>
        <div className="space-y-2 max-h-64 overflow-y-auto">
          {balance.dailyBalances.map((daily, index) => (
            <div key={index} className="bg-gray-50 p-3 rounded">
              <div className="flex justify-between items-center mb-2">
                <span className="font-medium text-gray-800">{daily.date}</span>
                <span className={`text-sm font-semibold ${daily.balance >= 0 ? 'text-red-600' : 'text-green-600'}`}>
                  {daily.balance >= 0 ? '+' : ''}{daily.balance} kcal
                </span>
              </div>
              <div className="grid grid-cols-2 gap-2 text-sm text-gray-600">
                <div>摂取: {daily.caloriesIntake} kcal</div>
                <div>消費: {daily.caloriesBurned} kcal</div>
              </div>
              {/* バランスバー */}
              <div className="mt-2">
                <div className="w-full bg-gray-200 rounded-full h-2">
                  <div 
                    className={`h-2 rounded-full ${daily.balance >= 0 ? 'bg-red-500' : 'bg-green-500'}`}
                    style={{ 
                      width: `${Math.min(Math.abs(daily.balance) / 1000 * 100, 100)}%` 
                    }}
                  ></div>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default CalorieBalanceChart;
