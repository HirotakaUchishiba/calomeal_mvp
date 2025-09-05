import React from 'react';

type Props = {
  summary: {
    caloriesIntake: number;
    caloriesBurned: number;
  };
};

export const DailySummaryNumbers = ({ summary }: Props) => {
  return (
    <div>
      <h2>今日のサマリー</h2>
      <p>摂取カロリー: {summary.caloriesIntake} kcal</p>
      <p>消費カロリー: {summary.caloriesBurned} kcal</p>
    </div>
  );
};