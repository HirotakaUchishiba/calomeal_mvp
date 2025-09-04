import React, { useState } from 'react';

type Props = {
  isOpen: boolean;
  onClose: () => void;
};

export const ExerciseLogModal = ({ isOpen, onClose }: Props) => {
  const [exerciseName, setExerciseName] = useState('');
  const [durationMinutes, setDurationMinutes] = useState('');
  const [caloriesBurned, setCaloriesBurned]= useState('');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    // TODO: タスク5でAPI記録ロジックを実装します
    console.log('Logging exercise:', {
      exerciseName,
      durationMinutes: parseInt(durationMinutes),
      caloriesBurned: parseFloat(caloriesBurned),
    });
    onClose();
  };

  if (!isOpen) return null;

  return (
    <div className="modal-backdrop">
      <div className="modal-content">
        <button onClick={onClose}>閉じる</button>
        <h2>運動を記録</h2>
        <form onSubmit={handleSubmit}>
          <div>
            <label htmlFor="exerciseName">運動名:</label>
            <input
              id="exerciseName"
              type="text"
              value={exerciseName}
              onChange={(e) => setExerciseName(e.target.value)}
              required
            />
          </div>
          <div>
            <label htmlFor="duration">実施時間 (分):</label>
            <input
              id="duration"
              type="number"
              value={durationMinutes}
              onChange={(e) => setDurationMinutes(e.target.value)}
              required
            />
          </div>
          <div>
            <label htmlFor="calories">消費カロリー (kcal):</label>
            <input
              id="calories"
              type="number"
              value={caloriesBurned}
              onChange={(e) => setCaloriesBurned(e.target.value)}
              required
            />
          </div>
          <button type="submit">この運動を記録する</button>
        </form>
      </div>
    </div>
  );
};