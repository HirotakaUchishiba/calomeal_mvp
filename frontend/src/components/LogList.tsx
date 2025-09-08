import React from 'react';
import { useQuery } from '@apollo/client';
import { GET_FOOD_LOGS_QUERY, GET_EXERCISE_LOGS_QUERY } from '../graphql/queries';

// ログアイテムの型定義
type LogItem = {
  id: string;
  type: 'food' | 'exercise';
  name: string;
  details: string;
  calories: number;
  loggedAt: string;
  icon: string;
  color: string;
};

type Props = {
  date: string;
};

export const LogList = ({ date }: Props) => {
  // 食事記録を取得
  const { data: foodData, loading: foodLoading, error: foodError } = useQuery(GET_FOOD_LOGS_QUERY, {
    variables: { date }
  });

  // 運動記録を取得
  const { data: exerciseData, loading: exerciseLoading, error: exerciseError } = useQuery(GET_EXERCISE_LOGS_QUERY, {
    variables: { date }
  });

  // データを統合してLogItem形式に変換
  const logs: LogItem[] = React.useMemo(() => {
    const foodLogs: LogItem[] = (foodData?.foodLogs || []).map((log: any) => ({
      id: log.id,
      type: 'food' as const,
      name: log.foodName,
      details: `${log.quantity}${log.unit}`,
      calories: log.calories,
      loggedAt: log.loggedAt,
      icon: '🍽️',
      color: '#FF9800'
    }));

    const exerciseLogs: LogItem[] = (exerciseData?.exerciseLogs || []).map((log: any) => ({
      id: log.id,
      type: 'exercise' as const,
      name: log.exerciseName,
      details: `${log.durationMinutes}分`,
      calories: log.caloriesBurned,
      loggedAt: log.loggedAt,
      icon: '🏃',
      color: '#2196F3'
    }));

    // 食事記録と運動記録を統合して時間順でソート
    return [...foodLogs, ...exerciseLogs].sort((a, b) => 
      new Date(b.loggedAt).getTime() - new Date(a.loggedAt).getTime()
    );
  }, [foodData, exerciseData]);

  // ローディング状態
  if (foodLoading || exerciseLoading) {
    return (
      <div style={{
        padding: '20px',
        textAlign: 'center',
        color: '#666',
        backgroundColor: '#f8f9fa',
        borderRadius: '8px',
        marginTop: '20px'
      }}>
        <div style={{ fontSize: '24px', marginBottom: '10px' }}>⏳</div>
        <p>記録を読み込み中...</p>
      </div>
    );
  }

  // エラー状態
  if (foodError || exerciseError) {
    return (
      <div style={{
        padding: '20px',
        textAlign: 'center',
        color: '#d32f2f',
        backgroundColor: '#ffebee',
        borderRadius: '8px',
        marginTop: '20px'
      }}>
        <div style={{ fontSize: '24px', marginBottom: '10px' }}>⚠️</div>
        <p>記録の読み込みに失敗しました</p>
        <p style={{ fontSize: '14px', marginTop: '5px' }}>
          {foodError?.message || exerciseError?.message}
        </p>
      </div>
    );
  }

  // ログがない場合の表示
  if (logs.length === 0) {
    return (
      <div style={{
        padding: '20px',
        textAlign: 'center',
        color: '#666',
        backgroundColor: '#f8f9fa',
        borderRadius: '8px',
        marginTop: '20px'
      }}>
        <div style={{ fontSize: '48px', marginBottom: '10px' }}>📝</div>
        <p>まだ記録がありません</p>
        <p style={{ fontSize: '14px', marginTop: '5px' }}>
          右下のボタンから記録を開始しましょう
        </p>
      </div>
    );
  }

  // 時間でソート（新しい順）
  const sortedLogs = [...logs].sort((a, b) => 
    new Date(b.loggedAt).getTime() - new Date(a.loggedAt).getTime()
  );

  // 時間をフォーマット
  const formatTime = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleTimeString('ja-JP', { 
      hour: '2-digit', 
      minute: '2-digit' 
    });
  };

  return (
    <div style={{
      marginTop: '20px',
      backgroundColor: '#fff',
      borderRadius: '8px',
      boxShadow: '0 2px 4px rgba(0, 0, 0, 0.1)'
    }}>
      <div style={{
        padding: '15px 20px',
        borderBottom: '1px solid #e9ecef',
        backgroundColor: '#f8f9fa',
        borderRadius: '8px 8px 0 0'
      }}>
        <h3 style={{ 
          margin: 0, 
          fontSize: '16px', 
          fontWeight: '600',
          color: '#495057'
        }}>
          今日の記録
        </h3>
      </div>

      <div style={{ padding: '0' }}>
        {sortedLogs.map((log, index) => (
          <div
            key={log.id}
            style={{
              display: 'flex',
              alignItems: 'center',
              padding: '15px 20px',
              borderBottom: index < sortedLogs.length - 1 ? '1px solid #f1f3f4' : 'none',
              transition: 'background-color 0.2s ease'
            }}
            onMouseOver={(e) => {
              e.currentTarget.style.backgroundColor = '#f8f9fa';
            }}
            onMouseOut={(e) => {
              e.currentTarget.style.backgroundColor = 'transparent';
            }}
          >
            {/* アイコン */}
            <div style={{
              width: '40px',
              height: '40px',
              borderRadius: '50%',
              backgroundColor: log.color + '20',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              fontSize: '18px',
              marginRight: '15px',
              flexShrink: 0
            }}>
              {log.icon}
            </div>

            {/* ログ情報 */}
            <div style={{ flex: 1, minWidth: 0 }}>
              <div style={{
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'flex-start',
                marginBottom: '4px'
              }}>
                <div>
                  <div style={{
                    fontSize: '16px',
                    fontWeight: '600',
                    color: '#212529',
                    marginBottom: '2px'
                  }}>
                    {log.name}
                  </div>
                  <div style={{
                    fontSize: '14px',
                    color: '#6c757d'
                  }}>
                    {log.details}
                  </div>
                </div>
                <div style={{
                  textAlign: 'right',
                  flexShrink: 0,
                  marginLeft: '10px'
                }}>
                  <div style={{
                    fontSize: '14px',
                    fontWeight: '600',
                    color: log.type === 'food' ? '#FF9800' : '#2196F3'
                  }}>
                    {log.calories} kcal
                  </div>
                  <div style={{
                    fontSize: '12px',
                    color: '#6c757d',
                    marginTop: '2px'
                  }}>
                    {formatTime(log.loggedAt)}
                  </div>
                </div>
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* 合計表示 */}
      <div style={{
        padding: '15px 20px',
        backgroundColor: '#f8f9fa',
        borderTop: '1px solid #e9ecef',
        borderRadius: '0 0 8px 8px'
      }}>
        <div style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center'
        }}>
          <span style={{
            fontSize: '14px',
            fontWeight: '600',
            color: '#495057'
          }}>
            合計
          </span>
          <div style={{
            display: 'flex',
            gap: '20px'
          }}>
            <span style={{
              fontSize: '14px',
              color: '#FF9800',
              fontWeight: '600'
            }}>
              摂取: {logs.filter(log => log.type === 'food').reduce((sum, log) => sum + log.calories, 0)} kcal
            </span>
            <span style={{
              fontSize: '14px',
              color: '#2196F3',
              fontWeight: '600'
            }}>
              消費: {logs.filter(log => log.type === 'exercise').reduce((sum, log) => sum + log.calories, 0)} kcal
            </span>
          </div>
        </div>
      </div>
    </div>
  );
};
