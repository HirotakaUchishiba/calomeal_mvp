import { test, expect } from '@playwright/test';

test.describe('記録機能のテスト', () => {
  test.beforeEach(async ({ page }) => {
    // ログインページにアクセス
    await page.goto('/login');
    
    // テスト用の認証情報でログイン
    await page.fill('input[type="email"]', 'test@example.com');
    await page.fill('input[type="password"]', 'password123');
    await page.click('button[type="submit"]');
    
    // ダッシュボードにリダイレクトされるまで待機
    await page.waitForURL('/dashboard', { timeout: 10000 });
    
    // ダッシュボードが表示されるまで待機
    await expect(page.locator('h1')).toContainText('ダッシュボード', { timeout: 10000 });
  });

  test('食事記録機能', async ({ page }) => {
    // フローティングアクションボタンをクリック
    await page.click('button[aria-label="メニューを開く"]');
    
    // 食事記録ボタンをクリック
    await page.click('button:has-text("食事を記録")');
    
    // 食事記録モーダルが開くことを確認
    await expect(page.locator('h2')).toContainText('食事記録');
    
    // 食品検索
    await page.fill('input[placeholder*="検索"]', 'りんご');
    await page.waitForTimeout(1000);
    
    // 検索結果から選択
    await page.click('text=りんご');
    
    // 数量を入力
    await page.fill('input[placeholder*="量"]', '200');
    await page.selectOption('select', 'g');
    
    // 記録ボタンをクリック
    await page.click('button:has-text("記録する")');
    
    // モーダルが閉じることを確認
    await expect(page.locator('h2')).not.toBeVisible();
    
    // ログリストに記録が表示されることを確認
    await expect(page.locator('text=りんご')).toBeVisible();
    await expect(page.locator('text=200g')).toBeVisible();
  });

  test('運動記録機能', async ({ page }) => {
    // フローティングアクションボタンをクリック
    await page.click('button[aria-label="メニューを開く"]');
    
    // 運動記録ボタンをクリック
    await page.click('button:has-text("運動を記録")');
    
    // 運動記録モーダルが開くことを確認
    await expect(page.locator('h2')).toContainText('運動記録');
    
    // 運動情報を入力
    await page.fill('input[placeholder*="運動名"]', 'ランニング');
    await page.fill('input[placeholder*="時間"]', '30');
    await page.fill('input[placeholder*="カロリー"]', '300');
    
    // 記録ボタンをクリック
    await page.click('button:has-text("記録する")');
    
    // モーダルが閉じることを確認
    await expect(page.locator('h2')).not.toBeVisible();
    
    // ログリストに記録が表示されることを確認
    await expect(page.locator('text=ランニング')).toBeVisible();
    await expect(page.locator('text=30分')).toBeVisible();
  });

  test('体重記録機能', async ({ page }) => {
    // フローティングアクションボタンをクリック
    await page.click('button[aria-label="メニューを開く"]');
    
    // 体重記録ボタンをクリック
    await page.click('button:has-text("体重を記録")');
    
    // 体重記録モーダルが開くことを確認
    await expect(page.locator('h2')).toContainText('体重記録');
    
    // 体重を入力
    await page.fill('input[type="number"]', '65.5');
    
    // 記録ボタンをクリック
    await page.click('button:has-text("記録する")');
    
    // モーダルが閉じることを確認
    await expect(page.locator('h2')).not.toBeVisible({ timeout: 10000 });
    
    // ページをリロードして体重記録が表示されることを確認
    await page.reload();
    await page.waitForTimeout(1000);
    
    // LogListに体重記録が表示されることを確認
    await expect(page.locator('div:has-text("体重")').first()).toBeVisible();
    await expect(page.locator('div:has-text("65.5kg")').first()).toBeVisible();
  });

  test('日付ナビゲーター機能', async ({ page }) => {
    // 日付ナビゲーターが表示されることを確認（より具体的なセレクターを使用）
    await expect(page.locator('div:has-text("今日")').first()).toBeVisible();
    
    // 前の日ボタンをクリック
    await page.click('button:has-text("←")');
    
    // 日付が変更されることを確認
    await expect(page.locator('div:has-text("昨日")').first()).toBeVisible();
    
    // 次の日ボタンをクリック
    await page.click('button:has-text("→")');
    
    // 今日に戻ることを確認
    await expect(page.locator('div:has-text("今日")').first()).toBeVisible();
  });
});
