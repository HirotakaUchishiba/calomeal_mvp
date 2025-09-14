import React from 'react';
import { useQuery } from '@apollo/client';
import { GET_NUTRITION_SUMMARY_QUERY } from '../../graphql/queries';

interface NutritionSummaryCardProps {
  date: string;
}

interface MealSummary {
  mealType: string;
  calories: number;
  protein: number;
  carbohydrate: number;
  fat: number;
  foodItems: string[];
}

interface NutritionSummary {
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

const NutritionSummaryCard: React.FC<NutritionSummaryCardProps> = ({ date }) => {
  const { loading, error, data } = useQuery<{ nutritionSummary: NutritionSummary }>(
    GET_NUTRITION_SUMMARY_QUERY,
    {
      variables: { date },
      errorPolicy: 'all'
    }
  );

  if (loading) return <div className="p-4 bg-white rounded-lg shadow">読み込み中...</div>;
  if (error) return <div className="p-4 bg-white rounded-lg shadow text-red-500">エラー: {error.message}</div>;
  if (!data?.nutritionSummary) return <div className="p-4 bg-white rounded-lg shadow">データがありません</div>;

  const summary = data.nutritionSummary;

  return (
    <div className="bg-white rounded-lg shadow-md p-6">
      <h3 className="text-lg font-semibold mb-4 text-gray-800">
        栄養サマリー - {summary.date}
      </h3>
      
      {/* カロリーバランス */}
      <div className="grid grid-cols-3 gap-4 mb-6">
        <div className="text-center">
          <div className="text-2xl font-bold text-blue-600">{summary.caloriesIntake}</div>
          <div className="text-sm text-gray-600">摂取カロリー</div>
        </div>
        <div className="text-center">
          <div className="text-2xl font-bold text-green-600">{summary.caloriesBurned}</div>
          <div className="text-sm text-gray-600">消費カロリー</div>
        </div>
        <div className="text-center">
          <div className={`text-2xl font-bold ${summary.caloriesBalance >= 0 ? 'text-red-600' : 'text-green-600'}`}>
            {summary.caloriesBalance}
          </div>
          <div className="text-sm text-gray-600">バランス</div>
        </div>
      </div>

      {/* 栄養素 */}
      <div className="grid grid-cols-2 gap-4 mb-6">
        <div className="bg-gray-50 p-3 rounded">
          <div className="text-sm text-gray-600">タンパク質</div>
          <div className="text-lg font-semibold">{summary.protein}g</div>
        </div>
        <div className="bg-gray-50 p-3 rounded">
          <div className="text-sm text-gray-600">炭水化物</div>
          <div className="text-lg font-semibold">{summary.carbohydrate}g</div>
        </div>
        <div className="bg-gray-50 p-3 rounded">
          <div className="text-sm text-gray-600">脂質</div>
          <div className="text-lg font-semibold">{summary.fat}g</div>
        </div>
        <div className="bg-gray-50 p-3 rounded">
          <div className="text-sm text-gray-600">食物繊維</div>
          <div className="text-lg font-semibold">{summary.fiber}g</div>
        </div>
      </div>

      {/* 食事サマリー */}
      <div>
        <h4 className="text-md font-semibold mb-3 text-gray-700">食事サマリー</h4>
        <div className="space-y-2">
          {summary.meals.map((meal, index) => (
            <div key={index} className="bg-gray-50 p-3 rounded">
              <div className="flex justify-between items-center mb-2">
                <span className="font-medium text-gray-800">{meal.mealType}</span>
                <span className="text-sm text-gray-600">{meal.calories} kcal</span>
              </div>
              <div className="grid grid-cols-3 gap-2 text-sm text-gray-600">
                <div>P: {meal.protein}g</div>
                <div>C: {meal.carbohydrate}g</div>
                <div>F: {meal.fat}g</div>
              </div>
              {meal.foodItems.length > 0 && (
                <div className="mt-2 text-sm text-gray-500">
                  食品: {meal.foodItems.join(', ')}
                </div>
              )}
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default NutritionSummaryCard;
