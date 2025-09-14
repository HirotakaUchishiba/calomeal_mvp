# services/foods/Dockerfile
FROM golang:1.25-alpine AS builder

# 作業ディレクトリを設定
WORKDIR /app

# 必要なパッケージをインストール
RUN apk add --no-cache git

# go.modとgo.sumをコピー
COPY go.mod go.sum ./

# ソースコードをコピー
COPY . .

# 依存関係をダウンロード（protoはボリュームマウントで提供）
RUN go mod download

# プロトコルバッファの生成（必要に応じて）
# RUN protoc --go_out=. --go_opt=paths=source_relative \
#     --go-grpc_out=. --go-grpc_opt=paths=source_relative \
#     proto/foods/v1/foods.proto

# アプリケーションをビルド
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 実行用の軽量なイメージ
FROM alpine:latest

# 必要なパッケージをインストール
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# ビルドしたバイナリをコピー
COPY --from=builder /app/main .

# ポートを公開
EXPOSE 50051

# アプリケーションを実行
CMD ["./main"]
