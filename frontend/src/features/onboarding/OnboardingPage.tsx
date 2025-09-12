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

        <button type="submit" disabled={loading}>
          {loading? '送信中...' : 'はじめる'}
        </button>

        {error && <p style={{ color: 'red' }}>エラーが発生しました: {error.message}</p>}
      </form>
    </div>
  );
};