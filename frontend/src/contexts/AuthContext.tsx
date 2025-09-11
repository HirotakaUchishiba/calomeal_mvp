// frontend/src/contexts/AuthContext.tsx

import React, { createContext, useContext, useEffect, useState } from 'react';
import type { ReactNode } from 'react';
import { getCurrentUser, signOut } from 'aws-amplify/auth';
import type { AuthUser } from 'aws-amplify/auth';

// 認証状態の型定義
interface AuthState {
  user: AuthUser | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  error: string | null;
}

// 認証コンテキストの型定義
interface AuthContextType extends AuthState {
  signOut: () => Promise<void>;
  clearError: () => void;
  refreshUser: () => Promise<void>;
  updateAuthState: (user: any) => void;
  setE2EAuthState: (user: any) => void;
}

// 認証コンテキストの作成
const AuthContext = createContext<AuthContextType | undefined>(undefined);

// 認証プロバイダーのProps
interface AuthProviderProps {
  children: ReactNode;
}

// 認証プロバイダーコンポーネント
export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [authState, setAuthState] = useState<AuthState>({
    user: null,
    isLoading: true,
    isAuthenticated: false,
    error: null,
  });

  // 現在のユーザー情報を取得
  const fetchCurrentUser = async () => {
    try {
      setAuthState(prev => ({ ...prev, isLoading: true, error: null }));
      
      // 開発環境では認証機能を無効化
      if (import.meta.env.DEV) {
        // E2Eテスト用の認証状態を確認
        const isE2ETest = window.location.search.includes('e2e-test=true');
        if (isE2ETest) {
          // E2Eテスト用の認証状態を強制的に設定
          const e2eUser = {
            userId: 'e2e-user-id',
            username: 'testuser@example.com',
            signInDetails: {
              loginId: 'testuser@example.com',
            },
          };
          setAuthState({
            user: e2eUser as any,
            isLoading: false,
            isAuthenticated: true,
            error: null,
          });
          return;
        }
        
        // 通常の開発環境では、localStorageから認証状態を確認
        const devUser = localStorage.getItem('dev-user');
        if (devUser) {
          try {
            const user = JSON.parse(devUser);
            setAuthState({
              user,
              isLoading: false,
              isAuthenticated: true,
              error: null,
            });
          } catch (error) {
            console.error('Failed to parse dev-user:', error);
            setAuthState({
              user: null,
              isLoading: false,
              isAuthenticated: false,
              error: null,
            });
          }
        } else {
          setAuthState({
            user: null,
            isLoading: false,
            isAuthenticated: false,
            error: null,
          });
        }
        return;
      }
      
      const user = await getCurrentUser();
      setAuthState({
        user,
        isLoading: false,
        isAuthenticated: true,
        error: null,
      });
    } catch (error) {
      console.log('No authenticated user:', error);
      setAuthState({
        user: null,
        isLoading: false,
        isAuthenticated: false,
        error: null,
      });
    }
  };

  // ユーザー情報をリフレッシュ
  const refreshUser = async () => {
    await fetchCurrentUser();
  };

  // 開発環境用の認証状態更新
  const updateAuthState = (user: any) => {
    setAuthState({
      user,
      isLoading: false,
      isAuthenticated: true,
      error: null,
    });
  };

  // E2Eテスト用の認証状態設定
  const setE2EAuthState = (user: any) => {
    if (import.meta.env.DEV) {
      localStorage.setItem('dev-user', JSON.stringify(user));
      setAuthState({
        user,
        isLoading: false,
        isAuthenticated: true,
        error: null,
      });
    }
  };

  // サインアウト
  const handleSignOut = async () => {
    try {
      setAuthState(prev => ({ ...prev, isLoading: true, error: null }));
      
      // 開発環境では認証機能を無効化
      if (import.meta.env.DEV) {
        localStorage.removeItem('dev-user');
        setAuthState({
          user: null,
          isLoading: false,
          isAuthenticated: false,
          error: null,
        });
        return;
      }
      
      await signOut();
      setAuthState({
        user: null,
        isLoading: false,
        isAuthenticated: false,
        error: null,
      });
    } catch (error) {
      console.error('Sign out error:', error);
      setAuthState(prev => ({
        ...prev,
        isLoading: false,
        error: error instanceof Error ? error.message : 'Sign out failed',
      }));
    }
  };

  // エラーをクリア
  const clearError = () => {
    setAuthState(prev => ({ ...prev, error: null }));
  };

  // コンポーネントマウント時にユーザー情報を取得
  useEffect(() => {
    fetchCurrentUser();
  }, []);

  // 認証状態の変更を監視
  useEffect(() => {
    const handleAuthStateChange = () => {
      fetchCurrentUser();
    };

    // 認証状態の変更イベントを監視
    window.addEventListener('authStateChange', handleAuthStateChange);
    
    return () => {
      window.removeEventListener('authStateChange', handleAuthStateChange);
    };
  }, []);

  const contextValue: AuthContextType = {
    ...authState,
    signOut: handleSignOut,
    clearError,
    refreshUser,
    updateAuthState,
    setE2EAuthState,
  };

  return (
    <AuthContext.Provider value={contextValue}>
      {children}
    </AuthContext.Provider>
  );
};

// 認証コンテキストを使用するカスタムフック
export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

// 認証が必要なコンポーネント用のフック
export const useRequireAuth = (): AuthContextType => {
  const auth = useAuth();
  
  useEffect(() => {
    if (!auth.isLoading && !auth.isAuthenticated) {
      // 認証されていない場合はログインページにリダイレクト
      window.location.href = '/login';
    }
  }, [auth.isLoading, auth.isAuthenticated]);

  return auth;
};

export default AuthContext;
