// import React from 'react';

type Props = {
  summary: {
    protein: number;
    carbohydrate: number;
    fat: number;
  };
  goals?: {
    protein: number;
    carbohydrate: number;
    fat: number;
  };
};

export const PFCProgressBars = ({ summary, goals }: Props) => {
  // デフォルトの目標値（体重70kgの成人男性を想定）
  const defaultGoals = {
    protein: 140, // 体重×2g
    carbohydrate: 280, // 体重×4g
    fat: 70, // 体重×1g
  };

  const targetGoals = goals || defaultGoals;

  const calculateProgress = (current: number, target: number) => {
    return Math.min((current / target) * 100, 100);
  };

  const getProgressBarStyle = (progress: number) => ({
    width: `${progress}%`,
    backgroundColor: progress >= 100 ? '#4CAF50' : progress >= 80 ? '#FF9800' : '#2196F3',
    height: '20px',
    borderRadius: '10px',
    transition: 'width 0.3s ease',
  });

  return (
    <div>
      <h3>PFCバランス</h3>
      
      <div style={{ marginBottom: '15px' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '5px' }}>
          <span>タンパク質(P): {summary.protein.toFixed(1)} / {targetGoals.protein} g</span>
          <span>{calculateProgress(summary.protein, targetGoals.protein).toFixed(0)}%</span>
        </div>
        <div style={{ border: '1px solid #ddd', borderRadius: '10px', overflow: 'hidden' }}>
          <div style={getProgressBarStyle(calculateProgress(summary.protein, targetGoals.protein))}></div>
        </div>
      </div>

      <div style={{ marginBottom: '15px' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '5px' }}>
          <span>炭水化物(C): {summary.carbohydrate.toFixed(1)} / {targetGoals.carbohydrate} g</span>
          <span>{calculateProgress(summary.carbohydrate, targetGoals.carbohydrate).toFixed(0)}%</span>
        </div>
        <div style={{ border: '1px solid #ddd', borderRadius: '10px', overflow: 'hidden' }}>
          <div style={getProgressBarStyle(calculateProgress(summary.carbohydrate, targetGoals.carbohydrate))}></div>
        </div>
      </div>

      <div style={{ marginBottom: '15px' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '5px' }}>
          <span>脂質(F): {summary.fat.toFixed(1)} / {targetGoals.fat} g</span>
          <span>{calculateProgress(summary.fat, targetGoals.fat).toFixed(0)}%</span>
        </div>
        <div style={{ border: '1px solid #ddd', borderRadius: '10px', overflow: 'hidden' }}>
          <div style={getProgressBarStyle(calculateProgress(summary.fat, targetGoals.fat))}></div>
        </div>
      </div>
    </div>
  );
};