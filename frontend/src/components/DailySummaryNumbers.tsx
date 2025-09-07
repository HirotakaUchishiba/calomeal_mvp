import React from 'react';

type Props = {
  summary: {
    caloriesIntake: number;
    caloriesBurned: number;
  };
  targetCalories?: number;
};

export const DailySummaryNumbers = ({ summary, targetCalories = 2000 }: Props) => {
  const netCalories = summary.caloriesIntake - summary.caloriesBurned;
  const remainingCalories = targetCalories - summary.caloriesIntake;
  
  const getCalorieStatus = () => {
    if (remainingCalories > 0) {
      return { text: '残り', value: remainingCalories, color: '#2196F3' };
    } else if (remainingCalories === 0) {
      return { text: '目標達成', value: 0, color: '#4CAF50' };
    } else {
      return { text: 'オーバー', value: Math.abs(remainingCalories), color: '#F44336' };
    }
  };

  const status = getCalorieStatus();

  return (
    <div style={{ 
      padding: '20px', 
      border: '1px solid #ddd', 
      borderRadius: '10px', 
      backgroundColor: '#f9f9f9',
      marginBottom: '20px'
    }}>
      <h2 style={{ marginTop: 0, color: '#333' }}>今日のサマリー</h2>
      
      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '20px', marginBottom: '20px' }}>
        <div style={{ textAlign: 'center' }}>
          <div style={{ fontSize: '24px', fontWeight: 'bold', color: '#4CAF50' }}>
            {summary.caloriesIntake.toFixed(0)}
          </div>
          <div style={{ fontSize: '14px', color: '#666' }}>摂取カロリー (kcal)</div>
        </div>
        
        <div style={{ textAlign: 'center' }}>
          <div style={{ fontSize: '24px', fontWeight: 'bold', color: '#FF9800' }}>
            {summary.caloriesBurned.toFixed(0)}
          </div>
          <div style={{ fontSize: '14px', color: '#666' }}>消費カロリー (kcal)</div>
        </div>
      </div>

      <div style={{ 
        textAlign: 'center', 
        padding: '15px', 
        backgroundColor: status.color + '20', 
        borderRadius: '8px',
        border: `2px solid ${status.color}`
      }}>
        <div style={{ fontSize: '20px', fontWeight: 'bold', color: status.color }}>
          {status.text}: {status.value.toFixed(0)} kcal
        </div>
        <div style={{ fontSize: '12px', color: '#666', marginTop: '5px' }}>
          目標: {targetCalories} kcal
        </div>
      </div>

      <div style={{ 
        marginTop: '15px', 
        padding: '10px', 
        backgroundColor: '#fff', 
        borderRadius: '5px',
        fontSize: '14px',
        color: '#666'
      }}>
        純カロリー: {netCalories.toFixed(0)} kcal
      </div>
    </div>
  );
};