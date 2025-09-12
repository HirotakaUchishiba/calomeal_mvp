import { test, expect } from '@playwright/test';

test.describe('E2E-HP-001: ハッピーパス', () => {
  test('新規ユーザー登録から食事記録まで', async ({ page }) => {
    // 1. Navigate: アプリケーションのURLにアクセスする
    await page.goto('/');
    
    // ログイン画面が表示されることを確認
    await expect(page).toHaveTitle(/CaloMeal/);
    await expect(page.locator('h1')).toContainText('ログイン');

    // 2. Action: 「サインアップ」し、新しいユーザーを登録する
    await page.click('text=サインアップ');
    
    // サインアップ画面に遷移することを確認
    await expect(page.locator('h1')).toContainText('サインアップ');
    
    // テスト用のユーザー情報を入力
    const testName = `test`;
    const testEmail = `test-${Date.now()}@example.com`;
    const testPassword = 'TestPassword123!';
    
    await page.fill('input[type="text"]', testName);
    await page.fill('input[type="email"]', testEmail);
    await page.fill('input[type="password"]', testPassword);
    await page.fill('input[id="confirmPassword"]', testPassword);
    
    // サインアップボタンをクリック
    await page.click('button[type="submit"]');
    
    // 3. Action: オンボーディング画面でプロフィールと目標を設定する
    // オンボーディング画面が表示されることを確認
    await expect(page.locator('h1')).toContainText('ようこそ！目標を設定しましょう');
    
    // プロフィール情報を入力
    await page.fill('id=height', '170');
    await page.fill('id=weight', '65');
    await page.selectOption('select[id="activityLevel"]', 'normal'); 
    
    // 目標設定
    await page.fill('id=targetWeight', '60');
    await page.fill('id=targetDate', '2025-12-31');
    
    // はじめるボタンをクリック
    await page.click('button:has-text("はじめる")');
    
    // 4. Action: ダッシュボード画面で食事記録モーダルを開く
    // ダッシュボードが表示されることを確認
    await expect(page.locator('h1')).toContainText('ダッシュボード');
    
    // フローティングアクションボタンをクリック
    await page.click('button:has-text("+")');
    
    // メニューが開くまで待機
    await page.waitForSelector('text=食事を記録', { timeout: 5000 });
    
    // 食事記録ボタンをクリック（アイコンボタンをクリック）
    await page.click('button:has-text("️")');
    
    // 5. Action: 「ごはん」を「200g」で記録する
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
    
    // 6. Assert: 記録が正しく保存されていることを検証する
    // ページをリフレッシュして記録を表示
    await page.reload();
    
    // ログリストに記録が表示されることを確認
    await expect(page.locator('div:has-text("選択された食品")').first()).toBeVisible();
    await expect(page.locator('div:has-text("200g")').first()).toBeVisible();
    
    // 7. Action: ログアウト操作を行う
    // ユーザーメニューを開く（ヘッダーのユーザー名をクリック）
    await page.click('button:has-text("ログアウト")');
    
    // 8. Assert: ログイン画面にリダイレクトされていることを検証する
    await expect(page.locator('h1')).toContainText('ログイン');
    await expect(page.url()).toContain('/login');
  });
});
