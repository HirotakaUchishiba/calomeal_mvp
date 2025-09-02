import React, { useState } from 'react';

export const OnboardingPage = () => {
  // プロフィール情報の状態
  const [height, setHeight] = useState('');
  const [weight, setWeight] = useState('');
  const [activityLevel, setActivityLevel] = useState('normal'); // 'low', 'normal', 'high'

  // 目標情報の状態
  const [targetWeight, setTargetWeight] = useState('');
  const [targetDate, setTargetDate] = useState('');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    // TODO: タスク5でAPI連携ロジックを実装します
    const profileInput = {
      height: parseFloat(height),
      weight: parseFloat(weight),
      activityLevel,
    };
    const goalInput = {
      targetWeight: parseFloat(targetWeight),
      targetDate,
    };
    console.log('Onboarding Submitted:', { profileInput, goalInput });
  };

  return (
    <div>
      <h1>ようこそ！目標を設定しましょう</h1>
      <form onSubmit={handleSubmit}>
        <h2>プロフィール</h2>
        <div>
          <label htmlFor="height">身長 (cm):</label>
          <input
            id="height"
            type="number"
            value={height}
            onChange={(e) => setHeight(e.target.value)}
            required
          />
        </div>
        <div>
          <label htmlFor="weight">現在の体重 (kg):</label>
          <input
            id="weight"
            type="number"
            value={weight}
            onChange={(e) => setWeight(e.target.value)}
            required
          />
        </div>
        <div>
          <label htmlFor="activityLevel">活動レベル:</label>
          <select
            id="activityLevel"
            value={activityLevel}
            onChange={(e) => setActivityLevel(e.target.value)}
          >
            <option value="low">低い</option>
            <option value="normal">普通</option>
            <option value="high">高い</option>
          </select>
        </div>

        <h2>目標</h2>
        <div>
          <label htmlFor="targetWeight">目標体重 (kg):</label>
          <input
            id="targetWeight"
            type="number"
            value={targetWeight}
            onChange={(e) => setTargetWeight(e.target.value)}
            required
          />
        </div>
        <div>
          <label htmlFor="targetDate">目標期日:</label>
          <input
            id="targetDate"
            type="date"
            value={targetDate}
            onChange={(e) => setTargetDate(e.target.value)}
            required
          />
        </div>

        <button type="submit">はじめる</button>
      </form>
    </div>
  );
};