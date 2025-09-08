# frontend/Dockerfile

# 1. ベースとなるNode.jsのイメージを選択
FROM node:18-alpine

# 2. コンテナ内の作業ディレクトリを指定
WORKDIR /app

# 3. package.jsonとpackage-lock.jsonをコピーして依存関係をインストール
COPY frontend/package*.json ./
RUN npm install

# 4. アプリケーションのソースコードをコピー
COPY frontend/ .

# 5. 依存関係を再インストール
RUN npm install

# 6. viteパッケージが正しくインストールされていることを確認
RUN npx vite --version

# 7. 開発サーバーを起動（ホスト0.0.0.0でリッスンしてコンテナ外からアクセス可能にする）
CMD ["sh", "-c", "echo 'Starting Vite development server...' && npx vite --host 0.0.0.0 --port 3000"]