import { useState } from 'react';
import { useQuery } from '@apollo/client';
import { FoodLogModal } from '../records/FoodLogModal';
import { ExerciseLogModal } from '../records/ExerciseLogModal';
import { WeightLogModal } from '../records/WeightLogModal';
import { DailySummaryNumbers } from '../../components/DailySummaryNumbers';
import { PFCProgressBars } from '../../components/PFCProgressBars';
import { DateNavigator } from '../../components/DateNavigator';
import { FloatingActionButton } from '../../components/FloatingActionButton';
import { LogList } from '../../components/LogList';
import { GET_DAILY_SUMMARY_QUERY } from '../../graphql/queries';
import { useAuth } from '../../contexts/AuthContext'; 

export const DashboardPage = () => {
  const [isFoodModalOpen, setFoodModalOpen] = useState(false);
  const [isExerciseModalOpen, setExerciseModalOpen] = useState(false);
  const [isWeightModalOpen, setWeightModalOpen] = useState(false);
  const [selectedDate, setSelectedDate] = useState(new Date().toISOString().split('T')[0]);

  const { user, signOut } = useAuth();
  const { data, loading, error } = useQuery(GET_DAILY_SUMMARY_QUERY, {
    variables: { date: selectedDate }
  });

  const handleLogout = async () => {
    try {
      await signOut();
    } catch (error) {
      console.error('Logout error:', error);
    }
  };

  if (loading) return <p>読み込み中...</p>;
  if (error) return <p>エラーが発生しました: {error.message}</p>;

  return (
    <div>
      <div style={{ 
        display: 'flex', 
        justifyContent: 'space-between', 
        alignItems: 'center',
        marginBottom: '20px',
        padding: '10px 0',
        borderBottom: '1px solid #eee'
      }}>
        <h1>ダッシュボード</h1>
        <div style={{ display: 'flex', alignItems: 'center', gap: '15px' }}>
          {user && (
            <span style={{ color: '#666' }}>
              こんにちは、{user.signInDetails?.loginId || 'ユーザー'}さん
            </span>
          )}
          <button 
            onClick={handleLogout}
            style={{
              padding: '8px 16px',
              backgroundColor: '#f44336',
              color: 'white',
              border: 'none',
              borderRadius: '5px',
              cursor: 'pointer'
            }}
          >
            ログアウト
          </button>
        </div>
      </div>

      {/* 日付ナビゲーター */}
      <DateNavigator 
        selectedDate={selectedDate} 
        onDateChange={setSelectedDate} 
      />

      {/* サマリー表示エリア */}
      <DailySummaryNumbers summary={data.dailySummary} />
      <PFCProgressBars summary={data.dailySummary} />

      {/* ログリスト */}
      <LogList date={selectedDate} />

      {/* フローティングアクションボタン */}
      <FloatingActionButton
        actions={[
          {
            id: 'food',
            label: '食事を記録',
            icon: '🍽️',
            onClick: () => setFoodModalOpen(true),
            color: '#FF9800'
          },
          {
            id: 'exercise',
            label: '運動を記録',
            icon: '🏃',
            onClick: () => setExerciseModalOpen(true),
            color: '#2196F3'
          },
          {
            id: 'weight',
            label: '体重を記録',
            icon: '📏',
            onClick: () => setWeightModalOpen(true),
            color: '#4CAF50'
          }
        ]}
        mainIcon="+"
        mainColor="#4CAF50"
      />

      {/* モーダル */}
      <FoodLogModal isOpen={isFoodModalOpen} onClose={() => setFoodModalOpen(false)} logDate={selectedDate} />
      <ExerciseLogModal isOpen={isExerciseModalOpen} onClose={() => setExerciseModalOpen(false)} logDate={selectedDate} />
      <WeightLogModal isOpen={isWeightModalOpen} onClose={() => setWeightModalOpen(false)} logDate={selectedDate} />
    </div>
  );
};