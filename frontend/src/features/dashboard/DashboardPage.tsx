import { useState } from 'react';
import { useQuery } from '@apollo/client';
import { FoodLogModal } from '../records/FoodLogModal';
import { ExerciseLogModal } from '../records/ExerciseLogModal';
import { DailySummaryNumbers } from '../../components/DailySummaryNumbers';
import { PFCProgressBars } from '../../components/PFCProgressBars';
import { GET_DAILY_SUMMARY_QUERY } from '../../graphql/queries';
import { useAuth } from '../../contexts/AuthContext'; 

export const DashboardPage = () => {
  const [isFoodModalOpen, setFoodModalOpen] = useState(false);
  const [isExerciseModalOpen, setExerciseModalOpen] = useState(false);
  const today = new Date().toISOString().split('T')[0];

  const { user, signOut } = useAuth();
  const { data, loading, error } = useQuery(GET_DAILY_SUMMARY_QUERY, {
    variables: { date: today }
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