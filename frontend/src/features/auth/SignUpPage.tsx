import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../../contexts/AuthContext';
import { useAuthActions } from '../../hooks/useAuthActions';

export const SignUpPage = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [name, setName] = useState('');
  const [showConfirmation, setShowConfirmation] = useState(false);
  const [confirmationCode, setConfirmationCode] = useState('');
  const [, setSignUpResult] = useState<any>(null);
  
  const { isAuthenticated, user } = useAuth();
  const { signUp, confirmSignUp, resendSignUpCode, isLoading, error, clearError } = useAuthActions();
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

    // パスワード確認
    if (password !== confirmPassword) {
      alert('パスワードが一致しません。');
      return;
    }

    // パスワード強度チェック
    if (password.length < 8) {
      alert('パスワードは8文字以上で入力してください。');
      return;
    }

    const result = await signUp({
      username: email,
      password: password,
      email: email,
      name: name,
    });

    if (result) {
      setSignUpResult(result);
      setShowConfirmation(true);
    }
  };

  const handleConfirmation = async (e: React.FormEvent) => {
    e.preventDefault();
    clearError();

    const success = await confirmSignUp(email, confirmationCode);

    if (success) {
      alert('アカウントが確認されました。ログインしてください。');
      navigate('/login');
    }
  };

  const handleResendCode = async () => {
    clearError();
    const result = await resendSignUpCode(email);
    if (result) {
      alert('確認コードを再送信しました。');
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
      <h1>サインアップ</h1>
      
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

      {!showConfirmation ? (
        <form onSubmit={handleSubmit}>
          <div style={{ marginBottom: '15px' }}>
            <label htmlFor="name" style={{ display: 'block', marginBottom: '5px' }}>
              お名前:
            </label>
            <input
              id="name"
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              style={{ 
                width: '100%', 
                padding: '10px', 
                border: '1px solid #ddd', 
                borderRadius: '5px' 
              }}
            />
          </div>

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
              minLength={8}
              style={{ 
                width: '100%', 
                padding: '10px', 
                border: '1px solid #ddd', 
                borderRadius: '5px' 
              }}
            />
            <small style={{ color: '#666' }}>
              8文字以上で入力してください
            </small>
          </div>

          <div style={{ marginBottom: '15px' }}>
            <label htmlFor="confirmPassword" style={{ display: 'block', marginBottom: '5px' }}>
              パスワード確認:
            </label>
            <input
              id="confirmPassword"
              type="password"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
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
            {isLoading ? 'サインアップ中...' : 'サインアップ'}
          </button>
        </form>
      ) : (
        <form onSubmit={handleConfirmation}>
          <h2>メール確認</h2>
          <p>
            {email} に確認コードを送信しました。<br />
            メールに記載された6桁のコードを入力してください。
          </p>
          
          <div style={{ marginBottom: '15px' }}>
            <label htmlFor="confirmationCode" style={{ display: 'block', marginBottom: '5px' }}>
              確認コード:
            </label>
            <input
              id="confirmationCode"
              type="text"
              value={confirmationCode}
              onChange={(e) => setConfirmationCode(e.target.value)}
              required
              maxLength={6}
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
              opacity: isLoading ? 0.7 : 1,
              marginBottom: '10px'
            }}
          >
            {isLoading ? '確認中...' : '確認'}
          </button>

          <button 
            type="button"
            onClick={handleResendCode}
            disabled={isLoading}
            style={{ 
              width: '100%', 
              padding: '10px', 
              backgroundColor: '#2196F3', 
              color: 'white', 
              border: 'none', 
              borderRadius: '5px',
              cursor: isLoading ? 'not-allowed' : 'pointer',
              opacity: isLoading ? 0.7 : 1
            }}
          >
            コードを再送信
          </button>
        </form>
      )}

      <div style={{ marginTop: '20px', textAlign: 'center' }}>
        {!showConfirmation ? (
          <p>
            すでにアカウントをお持ちですか？ <Link to="/login">ログイン</Link>
          </p>
        ) : (
          <p>
            <button 
              type="button"
              onClick={() => setShowConfirmation(false)}
              style={{ 
                background: 'none', 
                border: 'none', 
                color: '#2196F3', 
                cursor: 'pointer',
                textDecoration: 'underline'
              }}
            >
              サインアップに戻る
            </button>
          </p>
        )}
      </div>
    </div>
  );
};