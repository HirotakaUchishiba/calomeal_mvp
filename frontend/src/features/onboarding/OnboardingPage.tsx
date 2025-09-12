import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { gql } from '@apollo/client';
import { useMutation } from '@apollo/client/react';
import { useAuthActions } from '../../hooks/useAuthActions';

// バックエンドのschema.graphqlで定義したミューテーションを記述します
const COMPLETE_ONBOARDING_MUTATION = gql`
  mutation CompleteOnboarding($profile: UserProfileInput!, $goal: UserGoalInput!) {
    completeOnboarding(profile: $profile, goal: $goal) {
      id
      email
    }
  }
`;

export const OnboardingPage = () => {
  const navigate = useNavigate();
  const { completeOnboarding: markOnboardingComplete } = useAuthActions();
  
  // プロフィール情報の状態
  const [height, setHeight] = useState('');
  const [weight, setWeight] = useState('');
  const [activityLevel, setActivityLevel] = useState('normal'); // 'low', 'normal', 'high'

  // 目標情報の状態
  const [targetWeight, setTargetWeight] = useState('');
  const [targetDate, setTargetDate] = useState('');

  // useMutationフックを呼び出し、API通信用の関数や状態を取得します
  const [completeOnboarding, { loading, error }] = useMutation(COMPLETE_ONBOARDING_MUTATION);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    const profileInput = {
      height: parseFloat(height),
      weight: parseFloat(weight),
      activityLevel,
    };
    const goalInput = {
      targetWeight: parseFloat(targetWeight),
      targetDate,
    };

    // フォームの入力値をバックエンドAPIに送信します
    completeOnboarding({
      variables: {
        profile: profileInput,
        goal: goalInput,
      }
    }).then((response: any) => {
      console.log('Onboarding Succeeded:', response.data);
      // オンボーディング完了状態を更新
      markOnboardingComplete();
      // 成功したらダッシュボードへリダイレクト
      navigate('/dashboard');
    }).catch((err: any) => {
      console.error('Onboarding Failed:', err);
    });
  };

  return (
    <div>
      {/* ヘッダー */}
      <div style={{ 
        display: 'flex', 
        justifyContent: 'space-between', 
        alignItems: 'center',
        marginBottom: '20px',
        padding: '10px 0',
        borderBottom: '1px solid #eee'
      }}>
        <h1>ようこそ！目標を設定しましょう</h1>
        <div style={{ color: '#666', fontSize: '14px' }}>
          あなたの健康管理を始めましょう
        </div>
      </div>

      <form onSubmit={handleSubmit} style={{ maxWidth: '600px', margin: '0 auto' }}>
        {/* プロフィールセクション */}
        <div style={{
          background: '#f9f9f9',
          borderRadius: '8px',
          padding: '20px',
          marginBottom: '20px',
          border: '1px solid #eee'
        }}>
          <h2 style={{
            fontSize: '18px',
            fontWeight: '600',
            color: '#333',
            margin: '0 0 20px 0',
            display: 'flex',
            alignItems: 'center',
            gap: '8px'
          }}>
            <span style={{ fontSize: '20px' }}>👤</span>
            プロフィール
          </h2>
          
          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '20px' }}>
            <div>
              <label htmlFor="height" style={{
                display: 'block',
                fontSize: '14px',
                fontWeight: '500',
                color: '#666',
                marginBottom: '5px'
              }}>
                身長 (cm)
              </label>
              <input
                id="height"
                type="number"
                value={height}
                onChange={(e) => setHeight(e.target.value)}
                required
                style={{
                  width: '100%',
                  padding: '10px',
                  border: '1px solid #ddd',
                  borderRadius: '5px',
                  fontSize: '16px',
                  boxSizing: 'border-box'
                }}
              />
            </div>
            
            <div>
              <label htmlFor="weight" style={{
                display: 'block',
                fontSize: '14px',
                fontWeight: '500',
                color: '#666',
                marginBottom: '5px'
              }}>
                現在の体重 (kg)
              </label>
              <input
                id="weight"
                type="number"
                value={weight}
                onChange={(e) => setWeight(e.target.value)}
                required
                style={{
                  width: '100%',
                  padding: '10px',
                  border: '1px solid #ddd',
                  borderRadius: '5px',
                  fontSize: '16px',
                  boxSizing: 'border-box'
                }}
              />
            </div>
          </div>
          
          <div style={{ marginTop: '20px' }}>
            <label htmlFor="activityLevel" style={{
              display: 'block',
              fontSize: '14px',
              fontWeight: '500',
              color: '#666',
              marginBottom: '5px'
            }}>
              活動レベル
            </label>
            <select
              id="activityLevel"
              value={activityLevel}
              onChange={(e) => setActivityLevel(e.target.value)}
              style={{
                width: '100%',
                padding: '10px',
                border: '1px solid #ddd',
                borderRadius: '5px',
                fontSize: '16px',
                backgroundColor: 'white',
                boxSizing: 'border-box'
              }}
            >
              <option value="low">低い（デスクワーク中心）</option>
              <option value="normal">普通（軽い運動を含む）</option>
              <option value="high">高い（定期的な運動）</option>
            </select>
          </div>
        </div>

        {/* 目標セクション */}
        <div style={{
          background: '#f9f9f9',
          borderRadius: '8px',
          padding: '20px',
          marginBottom: '20px',
          border: '1px solid #eee'
        }}>
          <h2 style={{
            fontSize: '18px',
            fontWeight: '600',
            color: '#333',
            margin: '0 0 20px 0',
            display: 'flex',
            alignItems: 'center',
            gap: '8px'
          }}>
            <span style={{ fontSize: '20px' }}>🎯</span>
            目標
          </h2>
          
          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '20px' }}>
            <div>
              <label htmlFor="targetWeight" style={{
                display: 'block',
                fontSize: '14px',
                fontWeight: '500',
                color: '#666',
                marginBottom: '5px'
              }}>
                目標体重 (kg)
              </label>
              <input
                id="targetWeight"
                type="number"
                value={targetWeight}
                onChange={(e) => setTargetWeight(e.target.value)}
                required
                style={{
                  width: '100%',
                  padding: '10px',
                  border: '1px solid #ddd',
                  borderRadius: '5px',
                  fontSize: '16px',
                  boxSizing: 'border-box'
                }}
              />
            </div>
            
            <div>
              <label htmlFor="targetDate" style={{
                display: 'block',
                fontSize: '14px',
                fontWeight: '500',
                color: '#666',
                marginBottom: '5px'
              }}>
                目標期日
              </label>
              <input
                id="targetDate"
                type="date"
                value={targetDate}
                onChange={(e) => setTargetDate(e.target.value)}
                required
                style={{
                  width: '100%',
                  padding: '10px',
                  border: '1px solid #ddd',
                  borderRadius: '5px',
                  fontSize: '16px',
                  boxSizing: 'border-box'
                }}
              />
            </div>
          </div>
        </div>

        {/* エラーメッセージ */}
        {error && (
          <div style={{
            background: '#ffebee',
            border: '1px solid #f44336',
            borderRadius: '5px',
            padding: '10px',
            color: '#c62828',
            fontSize: '14px',
            marginBottom: '20px'
          }}>
            エラーが発生しました: {error.message}
          </div>
        )}

        {/* 送信ボタン */}
        <div style={{ textAlign: 'center' }}>
          <button 
            type="submit" 
            disabled={loading}
            style={{
              background: loading ? '#ccc' : '#4CAF50',
              color: 'white',
              border: 'none',
              borderRadius: '5px',
              padding: '12px 24px',
              fontSize: '16px',
              fontWeight: '500',
              cursor: loading ? 'not-allowed' : 'pointer',
              minWidth: '120px'
            }}
          >
            {loading ? '送信中...' : 'はじめる'}
          </button>
        </div>
      </form>
    </div>
  );
};