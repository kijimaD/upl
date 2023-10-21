# upl

upload zip file script.

ブラウザでのファイルアップロードが完全に制限されている環境でファイル転送させるために作った。Tiny File ManagerのエンドポイントにcURLで直にリクエストし、クッキー取得とアップロードを一度にやってくれる。

## install

```
$ go install github.com/kijimaD/upl@main
```

## how to use

このリポジトリでの設定で動かした Tiny File Manager にアップロードするのを前提としている。

```
$ docker-compose up -d
```

実行。カレントディレクトリにある `upload.zip` を、指定パスにある Tiny File Manager の指定ディレクトリにアップロードする。`.`の場合は、Tiny File Manager のトップに設定されているディレクトリに`upload.zip`をアップロードする。

```
$ upl localhost:7777 .
upl [ベースパス] [アップロード先ディレクトリ]
```
( http://localhost:7777 で Tiny File Managerにアクセスできているものとする。)

## docker run

```
$ docker run -v "$PWD/":/work -w /work --rm -it ghcr.io/kijimad/upl:latest localhost:7777 .
```
