import { test, expect } from '@playwright/test';

test.describe('記録機能のテスト', () => {
  test.beforeEach(async ({ page }) => {
    // ログインページにアクセス
    await page.goto('/login');
    
    // テスト用の認証情報でログイン
    await page.fill('input[type="email"]', 'test@example.com');
    await page.fill('input[type="password"]', 'password123');
    await page.click('button[type="submit"]');
    
    // 認証状態の更新を待機
    await page.waitForFunction(() => {
      const authState = window.localStorage.getItem('dev-user');
      return authState !== null;
    }, { timeout: 5000 });
    
    // ログイン処理の完了を待機
    await page.waitForTimeout(1000);
    
    // オンボーディングページに遷移することを確認
    await page.waitForURL('/onboarding', { timeout: 10000 });
    await expect(page.locator('h1')).toContainText('ようこそ！目標を設定しましょう');
    
    // オンボーディング情報を入力
    await page.fill('input[id="height"]', '170');
    await page.fill('input[id="weight"]', '70');
    await page.selectOption('select[id="activityLevel"]', 'normal');
    await page.fill('input[id="targetWeight"]', '65');
    await page.fill('input[id="targetDate"]', '2025-12-31');
    
    // オンボーディングを完了
    await page.click('button[type="submit"]');
    
    // ダッシュボードに遷移するまで待機
    await page.waitForURL('/dashboard', { timeout: 10000 });
    
    // ダッシュボードが表示されるまで待機
    await expect(page.locator('h1')).toContainText('ダッシュボード', { timeout: 10000 });
  });

  test('食事記録機能', async ({ page }) => {
    // フローティングアクションボタンをクリック
    await page.click('button:has-text("+")');
    
    // メニューが開くまで待機
    await page.waitForSelector('text=食事を記録', { timeout: 5000 });
    
    // 食事記録ボタンをクリック（アイコンボタンをクリック）
    await page.click('button:has-text("️")');
    
    // 食事記録モーダルが開くことを確認
    await expect(page.locator('h2:has-text("食事を記録")')).toBeVisible();
    
    // 食品検索
    await page.fill('input[placeholder*="例: 鶏むね肉"]', 'ごはん');
    await page.click('button:has-text("検索")');
    await page.waitForTimeout(1000);
    
    // 検索結果から選択
    await page.click('li:has-text("ごはん (白米)")');
    
    // 数量を入力
    await page.fill('input', '200');
    await page.selectOption('select', 'g');
    
    // 記録ボタンをクリック
    await page.click('button:has-text("この食事を記録する")');
    
    // モーダルが閉じることを確認
    await expect(page.locator('h2:has-text("食事を記録")')).not.toBeVisible();
    
    // ページをリフレッシュして記録を表示
    await page.reload();
    
    // ログリストに記録が表示されることを確認（修正）
    await expect(page.locator('div:has-text("選択された食品")').first()).toBeVisible();
    await expect(page.locator('div:has-text("200g")').first()).toBeVisible();
  });

  test('運動記録機能', async ({ page }) => {
    // フローティングアクションボタンをクリック
    await page.click('button:has-text("+")');
    
    // メニューが開くまで待機
    await page.waitForSelector('text=運動を記録', { timeout: 5000 });
    
    // 運動記録ボタンをクリック（アイコンボタンをクリック）
    await page.click('button:has-text("🏃")');
    
    // 運動記録モーダルが開くことを確認
    await expect(page.locator('h2:has-text("運動を記録")')).toBeVisible();
    
    // 運動情報を入力
    await page.fill('#exerciseName', 'ランニング');
    await page.fill('#duration', '30');
    await page.fill('#calories', '300');
    
    // 記録ボタンをクリック
    await page.click('button:has-text("この運動を記録する")');
    
    // モーダルが閉じることを確認
    await expect(page.locator('h2:has-text("運動を記録")')).not.toBeVisible();
    
    // ページをリフレッシュして記録を表示
    await page.reload();
    await page.waitForTimeout(1000);
    
    // ログリストに記録が表示されることを確認（修正）
    await expect(page.locator('div:has-text("ランニング")').first()).toBeVisible();
    await expect(page.locator('div:has-text("30分")').first()).toBeVisible();
  });

// 修正後のコード
test('体重記録機能', async ({ page }) => {
  // フローティングアクションボタンをクリック
  await page.click('button:has-text("+")');
  
  // メニューが開くまで待機
  await page.waitForSelector('text=体重を記録', { timeout: 5000 });
  
  // 体重記録ボタンをクリック（アイコンボタンをクリック）
  await page.click('button:has-text("📏")');
  
  // 体重記録モーダルが開くことを確認
  await expect(page.locator('h2:has-text("体重記録")')).toBeVisible();
  
  // 体重を入力
  await page.fill('input[type="number"]', '65.5');
  
  // 記録ボタンをクリック
  await page.click('button:has-text("記録する")');
  
  // モーダルが閉じるまで待機（修正）
  await page.waitForSelector('h2:has-text("体重記録")', { state: 'hidden', timeout: 10000 });
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
