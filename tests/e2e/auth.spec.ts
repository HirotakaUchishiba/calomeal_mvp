import { test, expect } from '@playwright/test';

test.describe('認証機能のテスト', () => {
  test('ログイン機能', async ({ page }) => {
    await page.goto('/login');
    
    // ログイン画面が表示されることを確認
    await expect(page.locator('h1')).toContainText('ログイン');
    
    // 開発環境用テストアカウント情報が表示されることを確認
    await expect(page.locator('text=開発環境用テストアカウント')).toBeVisible();
    
    // 無効な認証情報でログインを試行
    await page.fill('input[type="email"]', 'invalid@example.com');
    await page.fill('input[type="password"]', 'wrongpassword');
    await page.click('button[type="submit"]');
    
    // エラーメッセージが表示されることを確認
    await expect(page.locator('text=メールアドレスまたはパスワードが正しくありません')).toBeVisible();
    
    // 有効な認証情報でログインを試行
    await page.fill('input[type="email"]', 'test@example.com');
    await page.fill('input[type="password"]', 'password123');
    
    // コンソールログを監視
    page.on('console', msg => console.log('PAGE LOG:', msg.text()));
    
    await page.click('button[type="submit"]');
    
    // 認証状態の更新を待機
    await page.waitForFunction(() => {
      const authState = window.localStorage.getItem('dev-user');
      return authState !== null;
    }, { timeout: 5000 });
    
    // ログイン処理の完了を待機
    await page.waitForTimeout(1000);
    
    // 現在のURLを確認
    const currentUrl = page.url();
    console.log('Current URL after login:', currentUrl);
    
    // ダッシュボードにリダイレクトされることを確認
    await page.waitForURL('/dashboard', { timeout: 10000 });
    await expect(page.locator('h1')).toContainText('ダッシュボード', { timeout: 5000 });
  });

  test('サインアップ機能', async ({ page }) => {
    await page.goto('/login');
    
    // サインアップ画面に遷移
    await page.click('text=サインアップ');
    await expect(page.locator('h1')).toContainText('サインアップ');
    
    // 既存ユーザーでサインアップを試行
    await page.fill('input[type="text"]', 'テストユーザー');
    await page.fill('input[type="email"]', 'test@example.com');
    await page.fill('input[type="password"]', 'password123');
    await page.fill('input[id="confirmPassword"]', 'password123');
    await page.click('button[type="submit"]');
    
    // 既存ユーザーエラーが表示されることを確認
    await expect(page.locator('text=このメールアドレスは既に登録されています')).toBeVisible();
    
    // 新しいユーザーでサインアップを試行
    await page.fill('input[type="text"]', '新規ユーザー');
    await page.fill('input[type="email"]', 'newuser@example.com');
    await page.fill('input[type="password"]', 'newpassword123');
    await page.fill('input[id="confirmPassword"]', 'newpassword123');
    
    // コンソールログを監視
    page.on('console', msg => console.log('PAGE LOG:', msg.text()));
    
    await page.click('button[type="submit"]');
    
    // サインアップ処理の完了を待機
    await page.waitForTimeout(2000);
    
    // 認証状態の更新を待機
    await page.waitForFunction(() => {
      const authState = window.localStorage.getItem('dev-user');
      return authState !== null;
    }, { timeout: 5000 });
    
    // オンボーディングにリダイレクトされることを確認
    await page.waitForURL('/onboarding', { timeout: 10000 });
    await expect(page.locator('h1')).toContainText('ようこそ！目標を設定しましょう');
  });

  test('パスワードリセット機能', async ({ page }) => {
    await page.goto('/');
    
    // パスワードリセットリンクをクリック
    await page.click('text=パスワードを忘れた場合');
    
    // パスワードリセット画面が表示されることを確認
    await expect(page.locator('h2')).toContainText('パスワードリセット');
    
    // メールアドレスを入力
    await page.fill('input[type="email"]', 'test@example.com');
    await page.click('button[type="submit"]');
    
    // 送信完了メッセージが表示されることを確認（ページ内のメッセージ）
    await expect(page.getByText('パスワードリセット用のメールを送信しました。メールをご確認ください。')).toBeVisible();
    
    // 3秒後に自動的にログイン画面に戻ることを確認
    await page.waitForTimeout(4000);
    await expect(page.locator('h1')).toContainText('ログイン');
  });
});
