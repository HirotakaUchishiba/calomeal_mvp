# 1. ベースとなるGo言語のイメージを選択
FROM golang:1.25-alpine

# 2. コンテナ内の作業ディレクトリを指定
WORKDIR /app

# 3. Goモジュールの依存関係を先にコピーしてキャッシュする
COPY go.mod go.sum ./
RUN go mod tidy

# 4. wait-for-itスクリプトをコピー
COPY wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

# 5. アプリケーションのソースコードをコピーする
COPY . .

RUN go mod tidy

# 6. アプリケーションをビルドする
# CGO_ENABLED=0 は、静的リンクバイナリを生成するためのおまじない
RUN CGO_ENABLED=0 go build -o /server ./cmd/server

# 7. コンテナが起動したときに実行されるコマンド
# データベースの準備完了を待ってからサーバーを起動する
# ポート8080でリッスンしてコンテナ外からアクセス可能にする
EXPOSE 8080
CMD ["/wait-for-it.sh", "db:5432", "--", "/server"]