import React, { useState } from 'react';
import { FoodLogModal } from '../records/FoodLogModal';
import { ExerciseLogModal } from '../records/ExerciseLogModal';
// TODO: フェーズ3で体重記録モーダルもインポートする

export const DashboardPage = () => {
  // 各モーダルの表示状態を管理するためのstate
  const [isFoodModalOpen, setFoodModalOpen] = useState(false);
  const [isExerciseModalOpen, setExerciseModalOpen] = useState(false);

  // 今日の日付を "YYYY-MM-DD" 形式で取得（モーダルに渡すため）
  const today = new Date().toISOString().split('T')[0];

  return (
    <div>
      <h1>ダッシュボード</h1>
      <p>ここにサマリー情報が表示されます（フェーズ3で実装）</p>

      {/* 記録ボタン */}
      <div className="fab-container">
        <button onClick={() => setFoodModalOpen(true)}>食事を記録</button>
        <button onClick={() => setExerciseModalOpen(true)}>運動を記録</button>
        {/* TODO: 体重記録ボタンも追加 */}
      </div>

      {/* モーダルコンポーネントをレンダリング */}
      <FoodLogModal
        isOpen={isFoodModalOpen}
        onClose={() => setFoodModalOpen(false)}
        logDate={today}
      />
      <ExerciseLogModal
        isOpen={isExerciseModalOpen}
        onClose={() => setExerciseModalOpen(false)}
        logDate={today}
      />
    </div>
  );
};