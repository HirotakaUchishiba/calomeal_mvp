// frontend/src/components/ProtectedRoute.tsx

import React from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';

interface ProtectedRouteProps {
  children: React.ReactNode;
  requireAuth?: boolean;
  redirectTo?: string;
}

export const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ 
  children, 
  requireAuth = true,
  redirectTo 
}) => {
  const { isAuthenticated, isLoading, user } = useAuth();
  const location = useLocation();

  // オンボーディング完了状態をチェック
  const isOnboardingCompleted = user && (user as any).isOnboardingCompleted;

  // ローディング中は表示しない
  if (isLoading) {
    return (
      <div style={{ 
        display: 'flex', 
        justifyContent: 'center', 
        alignItems: 'center', 
        height: '100vh',
        flexDirection: 'column'
      }}>
        <div style={{ 
          width: '40px', 
          height: '40px', 
          border: '4px solid #f3f3f3',
          borderTop: '4px solid #3498db',
          borderRadius: '50%',
          animation: 'spin 1s linear infinite'
        }}></div>
        <p style={{ marginTop: '20px', color: '#666' }}>認証状態を確認中...</p>
        <style>{`
          @keyframes spin {
            0% { transform: rotate(0deg); }
            100% { transform: rotate(360deg); }
          }
        `}</style>
      </div>
    );
  }

  // 認証が必要な場合
  if (requireAuth && !isAuthenticated) {
    const redirectPath = redirectTo || '/login';
    return <Navigate to={redirectPath} state={{ from: location }} replace />;
  }

  // 認証済みユーザーが認証不要なページにアクセスした場合
  if (!requireAuth && isAuthenticated) {
    // オンボーディング完了状態に基づいてリダイレクト先を決定
    const redirectPath = redirectTo || (isOnboardingCompleted ? '/dashboard' : '/onboarding');
    return <Navigate to={redirectPath} replace />;
  }

  // 認証が必要なページで、オンボーディング完了状態に基づいてリダイレクト
  if (requireAuth && isAuthenticated) {
    const currentPath = location.pathname;
    
    // オンボーディング未完了のユーザーがダッシュボードにアクセスした場合
    if (currentPath === '/dashboard' && !isOnboardingCompleted) {
      return <Navigate to="/onboarding" replace />;
    }
    
    // オンボーディング完了済みのユーザーがオンボーディングページにアクセスした場合
    if (currentPath === '/onboarding' && isOnboardingCompleted) {
      return <Navigate to="/dashboard" replace />;
    }
  }

  return <>{children}</>;
};

// 認証が必要なページ用のラッパー
export const RequireAuth: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  return (
    <ProtectedRoute requireAuth={true}>
      {children}
    </ProtectedRoute>
  );
};

// 認証不要なページ用のラッパー（ログイン済みユーザーはリダイレクト）
export const RequireGuest: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  return (
    <ProtectedRoute requireAuth={false}>
      {children}
    </ProtectedRoute>
  );
};

export default ProtectedRoute;
