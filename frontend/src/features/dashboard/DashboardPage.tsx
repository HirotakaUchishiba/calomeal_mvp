import React, { useState } from 'react';
import { FoodLogModal } from '../records/FoodLogModal';
import { ExerciseLogModal } from '../records/ExerciseLogModal';
import { DailySummaryNumbers } from '../../components/DailySummaryNumbers';
import { PFCProgressBars } from '../../components/PFCProgressBars';

export const DashboardPage = () => {
  const [isFoodModalOpen, setFoodModalOpen] = useState(false);
  const [isExerciseModalOpen, setExerciseModalOpen] = useState(false);
  const today = new Date().toISOString().split('T')[0];

  // TODO: タスク4でAPIから実際のデータを取得する
  const dummySummary = {
    caloriesIntake: 1500,
    caloriesBurned: 300,
    protein: 80,
    carbohydrate: 200,
    fat: 50,
  };

  return (
    <div>
      <h1>ダッシュボード</h1>

      {/* サマリー表示エリア */}
      <DailySummaryNumbers summary={dummySummary} />
      <PFCProgressBars summary={dummySummary} />

      {/* 記録ボタン */}
      <div className="fab-container">
        <button onClick={() => setFoodModalOpen(true)}>食事を記録</button>
        <button onClick={() => setExerciseModalOpen(true)}>運動を記録</button>
      </div>

      {/* モーダル */}
      <FoodLogModal isOpen={isFoodModalOpen} onClose={() => setFoodModalOpen(false)} logDate={today} />
      <ExerciseLogModal isOpen={isExerciseModalOpen} onClose={() => setExerciseModalOpen(false)} logDate={today} />
    </div>
  );
};