# backend/Dockerfile

# 1. ベースとなるGo言語のイメージを選択
FROM golang:1.21-alpine

# 2. コンテナ内の作業ディレクトリを指定
WORKDIR /app

# 3. このコンテナが起動したときに実行されるデフォルトのコマンド
# (コンテナを起動し続けるためのおまじない。後でGoアプリの起動コマンドに置き換えます)
CMD [ "tail", "-f", "/dev/null" ]