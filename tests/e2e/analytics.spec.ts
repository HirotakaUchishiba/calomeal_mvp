import { test, expect } from '@playwright/test';

test.describe('Analytics機能のテスト', () => {
  test.beforeEach(async ({ page }) => {
    // ログイン処理
    await page.goto('/login');
    await page.fill('input[type="email"]', 'test@example.com');
    await page.fill('input[type="password"]', 'password123');
    await page.click('button[type="submit"]');
    
    // ダッシュボードにリダイレクトされるまで待機
    await page.waitForURL('/dashboard', { timeout: 10000 });
  });

  test('アナリティクスページへの遷移', async ({ page }) => {
    // アナリティクスボタンをクリック
    await page.click('a[href="/analytics"]');
    
    // アナリティクスページに遷移することを確認
    await page.waitForURL('/analytics', { timeout: 5000 });
    await expect(page.locator('h1')).toContainText('アナリティクス');
  });

  test('栄養サマリータブの表示', async ({ page }) => {
    await page.goto('/analytics');
    
    // 栄養サマリータブが選択されていることを確認
    await expect(page.locator('button:has-text("📊 日次サマリー")')).toBeVisible();
    
    // 日付選択フィールドが表示されることを確認
    await expect(page.locator('input[type="date"]')).toBeVisible();
    
    // 栄養サマリーカードが表示されることを確認（ローディング状態またはデータ）
    await expect(page.locator('text=栄養サマリー')).toBeVisible();
  });

  test('栄養トレンドタブの表示', async ({ page }) => {
    await page.goto('/analytics');
    
    // 栄養トレンドタブをクリック
    await page.click('button:has-text("📈 栄養トレンド")');
    
    // 開始日と終了日の選択フィールドが表示されることを確認
    await expect(page.locator('input[type="date"]').first()).toBeVisible();
    await expect(page.locator('input[type="date"]').nth(1)).toBeVisible();
    
    // 栄養トレンドチャートが表示されることを確認
    await expect(page.locator('text=栄養トレンド')).toBeVisible();
  });

  test('体重進捗タブの表示', async ({ page }) => {
    await page.goto('/analytics');
    
    // 体重進捗タブをクリック
    await page.click('button:has-text("⚖️ 体重進捗")');
    
    // 開始日と終了日の選択フィールドが表示されることを確認
    await expect(page.locator('input[type="date"]').first()).toBeVisible();
    await expect(page.locator('input[type="date"]').nth(1)).toBeVisible();
    
    // 体重進捗チャートが表示されることを確認
    await expect(page.locator('text=体重進捗')).toBeVisible();
  });

  test('カロリーバランスタブの表示', async ({ page }) => {
    await page.goto('/analytics');
    
    // カロリーバランスタブをクリック
    await page.click('button:has-text("⚖️ カロリーバランス")');
    
    // 開始日と終了日の選択フィールドが表示されることを確認
    await expect(page.locator('input[type="date"]').first()).toBeVisible();
    await expect(page.locator('input[type="date"]').nth(1)).toBeVisible();
    
    // カロリーバランスチャートが表示されることを確認
    await expect(page.locator('text=カロリーバランス')).toBeVisible();
  });

  test('日付選択機能', async ({ page }) => {
    await page.goto('/analytics');
    
    // 今日の日付を取得
    const today = new Date().toISOString().split('T')[0];
    
    // 日付フィールドに今日の日付を設定
    await page.fill('input[type="date"]', today);
    
    // 日付が正しく設定されることを確認
    const dateValue = await page.inputValue('input[type="date"]');
    expect(dateValue).toBe(today);
  });

  test('期間選択機能', async ({ page }) => {
    await page.goto('/analytics');
    
    // 栄養トレンドタブに移動
    await page.click('button:has-text("📈 栄養トレンド")');
    
    // 今日の日付を取得
    const today = new Date().toISOString().split('T')[0];
    const weekAgo = new Date(Date.now() - 7 * 24 * 60 * 60 * 1000).toISOString().split('T')[0];
    
    // 開始日と終了日を設定
    await page.fill('input[type="date"]', weekAgo);
    await page.fill('input[type="date"]', today);
    
    // 日付が正しく設定されることを確認
    const startDateValue = await page.inputValue('input[type="date"]');
    const endDateValue = await page.inputValue('input[type="date"]');
    expect(startDateValue).toBe(weekAgo);
    expect(endDateValue).toBe(today);
  });

  test('エラーハンドリング', async ({ page }) => {
    await page.goto('/analytics');
    
    // ネットワークエラーをシミュレート
    await page.route('**/query', route => route.abort());
    
    // エラーメッセージが表示されることを確認
    await expect(page.locator('text=エラー:')).toBeVisible({ timeout: 10000 });
  });

  test('ローディング状態の表示', async ({ page }) => {
    await page.goto('/analytics');
    
    // ローディング状態が表示されることを確認（短時間）
    await expect(page.locator('text=読み込み中...')).toBeVisible({ timeout: 1000 });
  });

  test('レスポンシブデザイン', async ({ page }) => {
    await page.goto('/analytics');
    
    // デスクトップ表示を確認
    await expect(page.locator('h1')).toBeVisible();
    
    // モバイル表示に変更
    await page.setViewportSize({ width: 375, height: 667 });
    
    // モバイル表示でもコンテンツが表示されることを確認
    await expect(page.locator('h1')).toBeVisible();
    await expect(page.locator('button:has-text("📊 日次サマリー")')).toBeVisible();
  });

  test('タブナビゲーション', async ({ page }) => {
    await page.goto('/analytics');
    
    // 各タブを順番にクリックして内容が切り替わることを確認
    const tabs = [
      { name: '📊 日次サマリー', content: '栄養サマリー' },
      { name: '📈 栄養トレンド', content: '栄養トレンド' },
      { name: '⚖️ 体重進捗', content: '体重進捗' },
      { name: '⚖️ カロリーバランス', content: 'カロリーバランス' }
    ];
    
    for (const tab of tabs) {
      await page.click(`button:has-text("${tab.name}")`);
      await expect(page.locator(`text=${tab.content}`)).toBeVisible();
    }
  });
});
