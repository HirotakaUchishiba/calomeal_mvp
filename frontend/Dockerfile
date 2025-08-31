# frontend/Dockerfile

# 1. ベースとなるNode.jsのイメージを選択
FROM node:18-alpine

# 2. コンテナ内の作業ディレクトリを指定
WORKDIR /app

# 3. このコンテナが起動したときに実行されるデフォルトのコマンド
# (コンテナを起動し続けるためのおまじない。後でReactの起動コマンドに置き換えます)
CMD [ "tail", "-f", "/dev/null" ]