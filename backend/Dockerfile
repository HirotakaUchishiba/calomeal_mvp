# 1. ベースとなるGo言語のイメージを選択
FROM golang:1.25-alpine

# 2. コンテナ内の作業ディレクトリを指定
WORKDIR /app

# 3. Goモジュールの依存関係を先にコピーしてキャッシュする
COPY go.mod go.sum ./

# 4. 依存関係を解決（protoはボリュームマウントで提供）
RUN go mod tidy

# 6. アプリケーションのソースコードをコピーする
COPY . .

# 7. アプリケーションをビルドする
# CGO_ENABLED=0 は、静的リンクバイナリを生成するためのおまじない
RUN CGO_ENABLED=0 go build -o /server ./cmd/server

# 8. コンテナが起動したときに実行されるコマンド
# ビルドして生成された /server という実行ファイルを起動する
CMD ["/server"]