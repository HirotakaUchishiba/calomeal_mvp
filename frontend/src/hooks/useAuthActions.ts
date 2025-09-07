// frontend/src/hooks/useAuthActions.ts

import { useState } from 'react';
import { 
  signUp, 
  confirmSignUp, 
  signIn, 
  resetPassword, 
  confirmResetPassword,
  resendSignUpCode,
  ResendSignUpCodeOutput,
  SignUpOutput,
  SignInOutput,
  ResetPasswordOutput,
  ConfirmResetPasswordOutput,
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

  // サインイン
  const handleSignIn = async (input: SignInInput): Promise<SignInOutput | null> => {
    try {
      setIsLoading(true);
      setError(null);

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
  ): Promise<ConfirmResetPasswordOutput | null> => {
    try {
      setIsLoading(true);
      setError(null);

      const result = await confirmResetPassword({
        username: input.username,
        confirmationCode: input.confirmationCode,
        newPassword: input.newPassword,
      });

      return result;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Password reset confirmation failed';
      setError(errorMessage);
      return null;
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
