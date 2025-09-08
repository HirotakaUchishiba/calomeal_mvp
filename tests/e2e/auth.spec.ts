import { test, expect } from '@playwright/test';

test.describe('認証機能のテスト', () => {
  test('ログイン機能', async ({ page }) => {
    await page.goto('/');
    
    // ログイン画面が表示されることを確認
    await expect(page.locator('h1')).toContainText('ログイン');
    
    // 無効な認証情報でログインを試行
    await page.fill('input[type="email"]', 'invalid@example.com');
    await page.fill('input[type="password"]', 'wrongpassword');
    await page.click('button[type="submit"]');
    
    // エラーメッセージが表示されることを確認
    await expect(page.locator('text=ログインに失敗しました')).toBeVisible();
  });

  test('サインアップ機能', async ({ page }) => {
    await page.goto('/');
    
    // サインアップ画面に遷移
    await page.click('text=アカウントをお持ちでない方はこちら');
    await expect(page.locator('h1')).toContainText('アカウント作成');
    
    // 無効なパスワードでサインアップを試行
    await page.fill('input[type="email"]', 'test@example.com');
    await page.fill('input[type="password"]', '123'); // 短すぎるパスワード
    await page.fill('input[placeholder*="確認"]', '123');
    
    // バリデーションエラーが表示されることを確認
    await expect(page.locator('text=パスワードは8文字以上')).toBeVisible();
  });

  test('パスワードリセット機能', async ({ page }) => {
    await page.goto('/');
    
    // パスワードリセットリンクをクリック
    await page.click('text=パスワードを忘れた方');
    
    // パスワードリセット画面が表示されることを確認
    await expect(page.locator('h1')).toContainText('パスワードリセット');
    
    // メールアドレスを入力
    await page.fill('input[type="email"]', 'test@example.com');
    await page.click('button[type="submit"]');
    
    // 送信完了メッセージが表示されることを確認
    await expect(page.locator('text=リセットメールを送信しました')).toBeVisible();
  });
});
