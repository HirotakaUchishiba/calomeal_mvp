# 1. ベースとなるGo言語のイメージを選択
FROM golang:1.25-alpine

# bashとnetcatをインストール（wait-for-it.shスクリプト用）
RUN apk add --no-cache bash netcat-openbsd

# 2. コンテナ内の作業ディレクトリを指定
WORKDIR /app

# 3. Goモジュールの依存関係を先にコピーしてキャッシュする
COPY backend/go.mod backend/go.sum ./
RUN go mod tidy

# 4. wait-for-itスクリプトは不要（ncコマンドを使用）

# 5. アプリケーションのソースコードをコピーする
COPY backend/ .

RUN go mod tidy

# 6. アプリケーションをビルドする
# CGO_ENABLED=0 は、静的リンクバイナリを生成するためのおまじない
RUN CGO_ENABLED=0 go build -o /server ./cmd/server

# 7. コンテナが起動したときに実行されるコマンド
# データベースの準備完了を待ってからサーバーを起動する
# ポート8080でリッスンしてコンテナ外からアクセス可能にする
EXPOSE 8080
CMD ["sh", "-c", "until nc -z db 5432; do echo 'Waiting for database...'; sleep 2; done; echo 'Database is ready!'; /server"]