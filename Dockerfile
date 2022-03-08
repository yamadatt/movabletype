FROM golang:1.17.7-alpine3.15
#FROM golang:latest

# アップデートとgitのインストール！！
RUN apk add --update &&  apk add git

# appディレクトリの作成
RUN mkdir /go/src/app

# ワーキングディレクトリの設定
WORKDIR /go/src/app

#ADD . /go/src/app
#RUN go mod init github.com/yamadatt/movabletype

VOLUME $(pwd):/go/src/app

RUN go get -u github.com/yamadatt/movabletype




# ホストのファイルをコンテナの作業ディレクトリに移行
#ADD . /go/src/app
