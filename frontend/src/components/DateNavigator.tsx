import React from 'react';

type Props = {
  selectedDate: string;
  onDateChange: (date: string) => void;
};

export const DateNavigator = ({ selectedDate, onDateChange }: Props) => {
  const today = new Date();
  const selected = new Date(selectedDate);
  
  // 日付の差分を計算
  const diffDays = Math.floor((selected.getTime() - today.getTime()) / (1000 * 60 * 60 * 24));
  
  // 日付の表示名を取得
  const getDateLabel = (date: Date) => {
    const today = new Date();
    const yesterday = new Date(today);
    yesterday.setDate(yesterday.getDate() - 1);
    const tomorrow = new Date(today);
    tomorrow.setDate(tomorrow.getDate() + 1);
    
    const dateStr = date.toDateString();
    const todayStr = today.toDateString();
    const yesterdayStr = yesterday.toDateString();
    const tomorrowStr = tomorrow.toDateString();
    
    if (dateStr === todayStr) return '今日';
    if (dateStr === yesterdayStr) return '昨日';
    if (dateStr === tomorrowStr) return '明日';
    
    return date.toLocaleDateString('ja-JP', { 
      month: 'short', 
      day: 'numeric',
      weekday: 'short'
    });
  };
  
  // 前の日付に移動
  const goToPreviousDay = () => {
    const newDate = new Date(selected);
    newDate.setDate(newDate.getDate() - 1);
    onDateChange(newDate.toISOString().split('T')[0]);
  };
  
  // 次の日付に移動
  const goToNextDay = () => {
    const newDate = new Date(selected);
    newDate.setDate(newDate.getDate() + 1);
    onDateChange(newDate.toISOString().split('T')[0]);
  };
  
  // 今日に戻る
  const goToToday = () => {
    onDateChange(today.toISOString().split('T')[0]);
  };
  
  return (
    <div style={{
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'space-between',
      padding: '15px 20px',
      backgroundColor: '#f8f9fa',
      borderBottom: '1px solid #e9ecef',
      borderRadius: '8px 8px 0 0'
    }}>
      {/* 前の日付ボタン */}
      <button
        onClick={goToPreviousDay}
        style={{
          padding: '8px 12px',
          backgroundColor: '#6c757d',
          color: 'white',
          border: 'none',
          borderRadius: '6px',
          cursor: 'pointer',
          fontSize: '14px',
          fontWeight: '500',
          transition: 'background-color 0.2s ease'
        }}
        onMouseOver={(e) => e.currentTarget.style.backgroundColor = '#5a6268'}
        onMouseOut={(e) => e.currentTarget.style.backgroundColor = '#6c757d'}
      >
        ←
      </button>
      
      {/* 中央の日付表示 */}
      <div style={{
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        cursor: 'pointer'
      }} onClick={goToToday}>
        <div style={{
          fontSize: '18px',
          fontWeight: 'bold',
          color: '#495057',
          marginBottom: '2px'
        }}>
          {getDateLabel(selected)}
        </div>
        <div style={{
          fontSize: '12px',
          color: '#6c757d'
        }}>
          {selected.toLocaleDateString('ja-JP', { 
            year: 'numeric',
            month: 'long',
            day: 'numeric'
          })}
        </div>
      </div>
      
      {/* 次の日付ボタン */}
      <button
        onClick={goToNextDay}
        style={{
          padding: '8px 12px',
          backgroundColor: '#6c757d',
          color: 'white',
          border: 'none',
          borderRadius: '6px',
          cursor: 'pointer',
          fontSize: '14px',
          fontWeight: '500',
          transition: 'background-color 0.2s ease'
        }}
        onMouseOver={(e) => e.currentTarget.style.backgroundColor = '#5a6268'}
        onMouseOut={(e) => e.currentTarget.style.backgroundColor = '#6c757d'}
      >
        →
      </button>
    </div>
  );
};
