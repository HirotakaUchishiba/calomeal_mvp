import { test, expect } from '@playwright/test';

test.describe('E2E-HP-001: ハッピーパス', () => {
  test('新規ユーザー登録から食事記録まで', async ({ page }) => {
    // 1. Navigate: アプリケーションのURLにアクセスする
    await page.goto('/');
    
    // ログイン画面が表示されることを確認
    await expect(page).toHaveTitle(/CaloMeal/);
    await expect(page.locator('h1')).toContainText('ログイン');

    // 2. Action: 「サインアップ」し、新しいユーザーを登録する
    await page.click('text=アカウントをお持ちでない方はこちら');
    
    // サインアップ画面に遷移することを確認
    await expect(page.locator('h1')).toContainText('アカウント作成');
    
    // テスト用のユーザー情報を入力
    const testEmail = `test-${Date.now()}@example.com`;
    const testPassword = 'TestPassword123!';
    
    await page.fill('input[type="email"]', testEmail);
    await page.fill('input[type="password"]', testPassword);
    await page.fill('input[placeholder*="確認"]', testPassword);
    
    // サインアップボタンをクリック
    await page.click('button[type="submit"]');
    
    // 3. Action: オンボーディング画面でプロフィールと目標を設定する
    // オンボーディング画面が表示されることを確認
    await expect(page.locator('h1')).toContainText('プロフィール設定');
    
    // プロフィール情報を入力
    await page.fill('input[placeholder*="身長"]', '170');
    await page.fill('input[placeholder*="体重"]', '65');
    await page.selectOption('select', 'moderate'); // 活動レベル
    
    // 目標設定
    await page.fill('input[placeholder*="目標体重"]', '60');
    await page.fill('input[type="date"]', '2025-12-31');
    
    // 完了ボタンをクリック
    await page.click('button:has-text("完了")');
    
    // 4. Action: ダッシュボード画面で食事記録モーダルを開く
    // ダッシュボードが表示されることを確認
    await expect(page.locator('h1')).toContainText('ダッシュボード');
    
    // フローティングアクションボタンをクリック
    await page.click('button[aria-label="メニューを開く"]');
    
    // 食事記録ボタンをクリック
    await page.click('button:has-text("食事を記録")');
    
    // 5. Action: 「ごはん」を「150g」で記録する
    // 食事記録モーダルが開くことを確認
    await expect(page.locator('h2')).toContainText('食事記録');
    
    // 食品検索
    await page.fill('input[placeholder*="検索"]', 'ごはん');
    await page.waitForTimeout(1000); // 検索結果の読み込み待ち
    
    // 検索結果から「ごはん」を選択
    await page.click('text=ごはん');
    
    // 数量を入力
    await page.fill('input[placeholder*="量"]', '150');
    await page.selectOption('select', 'g');
    
    // 記録ボタンをクリック
    await page.click('button:has-text("記録する")');
    
    // モーダルが閉じることを確認
    await expect(page.locator('h2')).not.toBeVisible();
    
    // 6. Assert: ダッシュボードの「摂取カロリー」表示が、ごはん150g分 (252 kcal) 増加していることを検証する
    // 摂取カロリーが252kcalになっていることを確認
    await expect(page.locator('text=摂取')).toContainText('252');
    
    // ログリストに記録が表示されることを確認
    await expect(page.locator('text=ごはん')).toBeVisible();
    await expect(page.locator('text=150g')).toBeVisible();
    
    // 7. Action: ログアウト操作を行う
    // ユーザーメニューを開く（ヘッダーのユーザー名をクリック）
    await page.click('button:has-text("ログアウト")');
    
    // 8. Assert: ログイン画面にリダイレクトされていることを検証する
    await expect(page.locator('h1')).toContainText('ログイン');
    await expect(page.url()).toContain('/login');
  });
});
