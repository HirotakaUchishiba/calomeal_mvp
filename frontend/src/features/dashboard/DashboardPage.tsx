import React, { useState } from 'react';
import { useQuery } from '@apollo/client';
import { FoodLogModal } from '../records/FoodLogModal';
import { ExerciseLogModal } from '../records/ExerciseLogModal';
import { DailySummaryNumbers } from '../../components/DailySummaryNumbers';
import { PFCProgressBars } from '../../components/PFCProgressBars';
import { GET_DAILY_SUMMARY_QUERY } from '../../graphql/queries'; 

export const DashboardPage = () => {
  const [isFoodModalOpen, setFoodModalOpen] = useState(false);
  const [isExerciseModalOpen, setExerciseModalOpen] = useState(false);
  const today = new Date().toISOString().split('T')[0];

  const { data, loading, error } = useQuery(GET_DAILY_SUMMARY_QUERY, {
    variables: { date: today }
  });

  if (loading) return <p>読み込み中...</p>;
  if (error) return <p>エラーが発生しました: {error.message}</p>;

  return (
    <div>
      <h1>ダッシュボード</h1>

      {/* サマリー表示エリア */}
      <DailySummaryNumbers summary={data.dailySummary} />
      <PFCProgressBars summary={data.dailySummary} />

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