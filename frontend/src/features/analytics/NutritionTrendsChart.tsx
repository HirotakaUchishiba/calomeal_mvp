import React from 'react';
import { useQuery } from '@apollo/client';
import { GET_NUTRITION_TRENDS_QUERY } from '../../graphql/queries';

interface NutritionTrendsChartProps {
  startDate: string;
  endDate: string;
}

interface MealSummary {
  mealType: string;
  calories: number;
  protein: number;
  carbohydrate: number;
  fat: number;
  foodItems: string[];
}

interface DailySummary {
  date: string;
  caloriesIntake: number;
  caloriesBurned: number;
  caloriesBalance: number;
  protein: number;
  carbohydrate: number;
  fat: number;
  fiber: number;
  sugar: number;
  sodium: number;
  meals: MealSummary[];
}

interface NutritionTrends {
  dailySummaries: DailySummary[];
  caloriesAvg: number;
  proteinAvg: number;
  carbohydrateAvg: number;
  fatAvg: number;
  caloriesTrend: number;
  proteinTrend: number;
  carbohydrateTrend: number;
  fatTrend: number;
}

const NutritionTrendsChart: React.FC<NutritionTrendsChartProps> = ({ startDate, endDate }) => {
  const { loading, error, data } = useQuery<{ nutritionTrends: NutritionTrends }>(
    GET_NUTRITION_TRENDS_QUERY,
    {
      variables: { startDate, endDate },
      errorPolicy: 'all'
    }
  );

  if (loading) return <div className="p-4 bg-white rounded-lg shadow">読み込み中...</div>;
  if (error) return <div className="p-4 bg-white rounded-lg shadow text-red-500">エラー: {error.message}</div>;
  if (!data?.nutritionTrends) return <div className="p-4 bg-white rounded-lg shadow">データがありません</div>;

  const trends = data.nutritionTrends;

  const formatTrend = (trend: number) => {
    const sign = trend >= 0 ? '+' : '';
    return `${sign}${trend.toFixed(1)}%`;
  };

  const getTrendColor = (trend: number) => {
    if (trend > 0) return 'text-red-600';
    if (trend < 0) return 'text-green-600';
    return 'text-gray-600';
  };

  return (
    <div className="bg-white rounded-lg shadow-md p-6">
      <h3 className="text-lg font-semibold mb-4 text-gray-800">
        栄養トレンド ({startDate} - {endDate})
      </h3>

      {/* 平均値とトレンド */}
      <div className="grid grid-cols-2 gap-4 mb-6">
        <div className="bg-blue-50 p-4 rounded-lg">
          <div className="text-sm text-gray-600">平均カロリー</div>
          <div className="text-2xl font-bold text-blue-600">{trends.caloriesAvg.toFixed(0)}</div>
          <div className={`text-sm ${getTrendColor(trends.caloriesTrend)}`}>
            {formatTrend(trends.caloriesTrend)}
          </div>
        </div>
        <div className="bg-green-50 p-4 rounded-lg">
          <div className="text-sm text-gray-600">平均タンパク質</div>
          <div className="text-2xl font-bold text-green-600">{trends.proteinAvg.toFixed(1)}g</div>
          <div className={`text-sm ${getTrendColor(trends.proteinTrend)}`}>
            {formatTrend(trends.proteinTrend)}
          </div>
        </div>
        <div className="bg-yellow-50 p-4 rounded-lg">
          <div className="text-sm text-gray-600">平均炭水化物</div>
          <div className="text-2xl font-bold text-yellow-600">{trends.carbohydrateAvg.toFixed(1)}g</div>
          <div className={`text-sm ${getTrendColor(trends.carbohydrateTrend)}`}>
            {formatTrend(trends.carbohydrateTrend)}
          </div>
        </div>
        <div className="bg-purple-50 p-4 rounded-lg">
          <div className="text-sm text-gray-600">平均脂質</div>
          <div className="text-2xl font-bold text-purple-600">{trends.fatAvg.toFixed(1)}g</div>
          <div className={`text-sm ${getTrendColor(trends.fatTrend)}`}>
            {formatTrend(trends.fatTrend)}
          </div>
        </div>
      </div>

      {/* 日次サマリー */}
      <div>
        <h4 className="text-md font-semibold mb-3 text-gray-700">日次サマリー</h4>
        <div className="space-y-2 max-h-64 overflow-y-auto">
          {trends.dailySummaries.map((summary, index) => (
            <div key={index} className="bg-gray-50 p-3 rounded">
              <div className="flex justify-between items-center mb-2">
                <span className="font-medium text-gray-800">{summary.date}</span>
                <span className="text-sm text-gray-600">
                  {summary.caloriesIntake} / {summary.caloriesBurned} kcal
                </span>
              </div>
              <div className="grid grid-cols-4 gap-2 text-sm text-gray-600">
                <div>P: {summary.protein}g</div>
                <div>C: {summary.carbohydrate}g</div>
                <div>F: {summary.fat}g</div>
                <div className={`${summary.caloriesBalance >= 0 ? 'text-red-600' : 'text-green-600'}`}>
                  バランス: {summary.caloriesBalance}
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default NutritionTrendsChart;
