// frontend/src/hooks/useAuthActions.ts

import { useState } from 'react';
import { 
  signUp, 
  confirmSignUp, 
  signIn, 
  resetPassword, 
  confirmResetPassword,
  resendSignUpCode,
} from 'aws-amplify/auth';
import type {
  ResendSignUpCodeOutput,
  SignUpOutput,
  SignInOutput,
  ResetPasswordOutput,
} from 'aws-amplify/auth';

// サインアップ用の型定義
interface SignUpInput {
  username: string;
  password: string;
  email: string;
  name?: string;
}

// サインイン用の型定義
interface SignInInput {
  username: string;
  password: string;
}

// パスワードリセット用の型定義
interface ResetPasswordInput {
  username: string;
}

interface ConfirmResetPasswordInput {
  username: string;
  confirmationCode: string;
  newPassword: string;
}

// 開発環境用のテストユーザー
const DEV_TEST_USERS = [
  {
    username: 'test@example.com',
    password: 'password123',
    name: 'テストユーザー',
  },
  {
    username: 'admin@example.com',
    password: 'admin123',
    name: '管理者ユーザー',
  },
  {
    username: 'user@example.com',
    password: 'user123',
    name: '一般ユーザー',
  },
];

// 認証アクション用のフック
export const useAuthActions = () => {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // エラーをクリア
  const clearError = () => setError(null);

  // サインアップ
  const handleSignUp = async (input: SignUpInput): Promise<SignUpOutput | null> => {
    try {
      setIsLoading(true);
      setError(null);

      // 開発環境でのサインアップ処理
      if (import.meta.env.DEV) {
        await new Promise(resolve => setTimeout(resolve, 1000));
        
        // 既存ユーザーのチェック
        if (DEV_TEST_USERS.some(user => user.username === input.username)) {
          setError('このメールアドレスは既に登録されています。');
          return null;
        }
        
        // パスワード強度チェック
        if (input.password.length < 8) {
          setError('パスワードは8文字以上で入力してください。');
          return null;
        }
        
        // 開発環境では即座にサインアップ完了として扱う
        const devUser = {
          userId: `dev-user-${input.username.replace('@', '-').replace('.', '-')}`,
          username: input.username,
          name: input.name || '新規ユーザー',
          signInDetails: {
            loginId: input.username,
          },
        };
        localStorage.setItem('dev-user', JSON.stringify(devUser));
        
        // 認証状態変更イベントを発火
        window.dispatchEvent(new CustomEvent('authStateChange'));
        
        return {
          isSignUpComplete: true,
          userId: devUser.userId,
          nextStep: {
            signUpStep: 'DONE',
          },
        };
      }

      const result = await signUp({
        username: input.username,
        password: input.password,
        options: {
          userAttributes: {
            email: input.email,
            name: input.name || '',
          },
        },
      });

      return result;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Sign up failed';
      setError(errorMessage);
      return null;
    } finally {
      setIsLoading(false);
    }
  };

  // サインアップ確認
  const handleConfirmSignUp = async (
    username: string, 
    confirmationCode: string
  ): Promise<boolean> => {
    try {
      setIsLoading(true);
      setError(null);

      await confirmSignUp({
        username,
        confirmationCode,
      });

      return true;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Confirmation failed';
      setError(errorMessage);
      return false;
    } finally {
      setIsLoading(false);
    }
  };

  // 開発環境での認証チェック
  const validateDevCredentials = (username: string, password: string): boolean => {
    return DEV_TEST_USERS.some(user => 
      user.username === username && user.password === password
    );
  };

  // 開発環境でのユーザー情報取得
  const getDevUserInfo = (username: string) => {
    return DEV_TEST_USERS.find(user => user.username === username);
  };

  // サインイン
  const handleSignIn = async (input: SignInInput): Promise<SignInOutput | null> => {
    try {
      setIsLoading(true);
      setError(null);

      // 開発環境での認証チェック
      if (import.meta.env.DEV) {
        console.log('E2E Test: Attempting login with:', input.username, input.password);
        await new Promise(resolve => setTimeout(resolve, 1000));
        
        // 認証情報をチェック
        const isValid = validateDevCredentials(input.username, input.password);
        console.log('E2E Test: Credentials valid:', isValid);
        
        if (!isValid) {
          setError('メールアドレスまたはパスワードが正しくありません。');
          return null;
        }
        
        // 認証成功時は、localStorageにユーザー情報を保存
        const userInfo = getDevUserInfo(input.username);
        const devUser = {
          userId: `dev-user-${input.username.replace('@', '-').replace('.', '-')}`,
          username: input.username,
          name: userInfo?.name || 'テストユーザー',
          signInDetails: {
            loginId: input.username,
          },
        };
        localStorage.setItem('dev-user', JSON.stringify(devUser));
        console.log('E2E Test: Login successful, user stored:', devUser);
        
        // 認証状態変更イベントを発火
        window.dispatchEvent(new CustomEvent('authStateChange'));
        
        return {
          isSignedIn: true,
          nextStep: {
            signInStep: 'DONE',
          },
        };
      }

      const result = await signIn({
        username: input.username,
        password: input.password,
      });

      return result;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Sign in failed';
      setError(errorMessage);
      return null;
    } finally {
      setIsLoading(false);
    }
  };

  // パスワードリセット開始
  const handleResetPassword = async (input: ResetPasswordInput): Promise<ResetPasswordOutput | null> => {
    try {
      setIsLoading(true);
      setError(null);

      const result = await resetPassword({
        username: input.username,
      });

      return result;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Password reset failed';
      setError(errorMessage);
      return null;
    } finally {
      setIsLoading(false);
    }
  };

  // パスワードリセット確認
  const handleConfirmResetPassword = async (
    input: ConfirmResetPasswordInput
  ): Promise<void> => {
    try {
      setIsLoading(true);
      setError(null);

      await confirmResetPassword({
        username: input.username,
        confirmationCode: input.confirmationCode,
        newPassword: input.newPassword,
      });
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Password reset confirmation failed';
      setError(errorMessage);
      throw err;
    } finally {
      setIsLoading(false);
    }
  };

  // 確認コード再送信
  const handleResendSignUpCode = async (username: string): Promise<ResendSignUpCodeOutput | null> => {
    try {
      setIsLoading(true);
      setError(null);

      const result = await resendSignUpCode({
        username,
      });

      return result;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Resend code failed';
      setError(errorMessage);
      return null;
    } finally {
      setIsLoading(false);
    }
  };

  return {
    isLoading,
    error,
    clearError,
    signUp: handleSignUp,
    confirmSignUp: handleConfirmSignUp,
    signIn: handleSignIn,
    resetPassword: handleResetPassword,
    confirmResetPassword: handleConfirmResetPassword,
    resendSignUpCode: handleResendSignUpCode,
  };
};
