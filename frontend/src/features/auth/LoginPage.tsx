import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../../contexts/AuthContext';
import { useAuthActions } from '../../hooks/useAuthActions';

export const LoginPage = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showPasswordReset, setShowPasswordReset] = useState(false);
  const [resetEmail, setResetEmail] = useState('');
  
  const { isAuthenticated, user, updateAuthState } = useAuth();
  const { signIn, resetPassword, isLoading, error, clearError } = useAuthActions();
  const navigate = useNavigate();

  // 認証済みの場合はダッシュボードにリダイレクト
  useEffect(() => {
    if (isAuthenticated && user) {
      navigate('/dashboard');
    }
  }, [isAuthenticated, user, navigate]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    clearError();

    const result = await signIn({
      username: email,
      password: password,
    });

    if (result) {
      // 開発環境では認証状態を即座に更新
      if (import.meta.env.DEV) {
        const devUser = {
          userId: 'dev-user-id',
          username: email,
          signInDetails: {
            loginId: email,
          },
        };
        updateAuthState(devUser);
      }
      // サインイン成功時はAuthContextが自動的に状態を更新
      navigate('/dashboard');
    }
  };

  const handlePasswordReset = async (e: React.FormEvent) => {
    e.preventDefault();
    clearError();

    const result = await resetPassword({
      username: resetEmail,
    });

    if (result) {
      alert('パスワードリセット用のメールを送信しました。メールをご確認ください。');
      setShowPasswordReset(false);
      setResetEmail('');
    }
  };

  if (isAuthenticated) {
    return (
      <div>
        <h1>ログイン済み</h1>
        <p>ダッシュボードにリダイレクトしています...</p>
      </div>
    );
  }

  return (
    <div style={{ maxWidth: '400px', margin: '0 auto', padding: '20px' }}>
      <h1>ログイン</h1>
      
      {error && (
        <div style={{ 
          color: 'red', 
          backgroundColor: '#ffe6e6', 
          padding: '10px', 
          borderRadius: '5px',
          marginBottom: '20px'
        }}>
          {error}
        </div>
      )}

      {!showPasswordReset ? (
        <form onSubmit={handleSubmit}>
          <div style={{ marginBottom: '15px' }}>
            <label htmlFor="email" style={{ display: 'block', marginBottom: '5px' }}>
              メールアドレス:
            </label>
            <input
              id="email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              style={{ 
                width: '100%', 
                padding: '10px', 
                border: '1px solid #ddd', 
                borderRadius: '5px' 
              }}
            />
          </div>
          
          <div style={{ marginBottom: '15px' }}>
            <label htmlFor="password" style={{ display: 'block', marginBottom: '5px' }}>
              パスワード:
            </label>
            <input
              id="password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              style={{ 
                width: '100%', 
                padding: '10px', 
                border: '1px solid #ddd', 
                borderRadius: '5px' 
              }}
            />
          </div>
          
          <button 
            type="submit" 
            disabled={isLoading}
            style={{ 
              width: '100%', 
              padding: '12px', 
              backgroundColor: '#4CAF50', 
              color: 'white', 
              border: 'none', 
              borderRadius: '5px',
              cursor: isLoading ? 'not-allowed' : 'pointer',
              opacity: isLoading ? 0.7 : 1
            }}
          >
            {isLoading ? 'ログイン中...' : 'ログイン'}
          </button>
        </form>
      ) : (
        <form onSubmit={handlePasswordReset}>
          <h2>パスワードリセット</h2>
          <div style={{ marginBottom: '15px' }}>
            <label htmlFor="resetEmail" style={{ display: 'block', marginBottom: '5px' }}>
              メールアドレス:
            </label>
            <input
              id="resetEmail"
              type="email"
              value={resetEmail}
              onChange={(e) => setResetEmail(e.target.value)}
              required
              style={{ 
                width: '100%', 
                padding: '10px', 
                border: '1px solid #ddd', 
                borderRadius: '5px' 
              }}
            />
          </div>
          
          <button 
            type="submit" 
            disabled={isLoading}
            style={{ 
              width: '100%', 
              padding: '12px', 
              backgroundColor: '#2196F3', 
              color: 'white', 
              border: 'none', 
              borderRadius: '5px',
              cursor: isLoading ? 'not-allowed' : 'pointer',
              opacity: isLoading ? 0.7 : 1
            }}
          >
            {isLoading ? '送信中...' : 'リセットメールを送信'}
          </button>
        </form>
      )}

      <div style={{ marginTop: '20px', textAlign: 'center' }}>
        {!showPasswordReset ? (
          <>
            <p>
              アカウントをお持ちでないですか？ <Link to="/signup">サインアップ</Link>
            </p>
            <p>
              <button 
                type="button"
                onClick={() => setShowPasswordReset(true)}
                style={{ 
                  background: 'none', 
                  border: 'none', 
                  color: '#2196F3', 
                  cursor: 'pointer',
                  textDecoration: 'underline'
                }}
              >
                パスワードを忘れた場合
              </button>
            </p>
          </>
        ) : (
          <p>
            <button 
              type="button"
              onClick={() => setShowPasswordReset(false)}
              style={{ 
                background: 'none', 
                border: 'none', 
                color: '#2196F3', 
                cursor: 'pointer',
                textDecoration: 'underline'
              }}
            >
              ログインに戻る
            </button>
          </p>
        )}
      </div>
    </div>
  );
};