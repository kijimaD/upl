# upl

upload zip file script.

ブラウザでのファイルアップロードが完全に制限されている環境でファイル転送させるために作った。Tiny File Managerのエンドポイントに直にリクエストし、クッキー取得とアップロードを一度にやってくれる。

このツールはクロスプラットフォームを念頭に置いて作成した。Linux, Windows, MacOS で動作する。実行可能バイナリ以外の依存性はない。

## install

```
$ go install github.com/kijimaD/upl@main
```

あるいは https://github.com/kijimaD/upl/releases からダウンロードする。

## prerequisite

v2.4.3 Tiny File Manager が動いているのを前提とする。Tiny File Manager のバージョンによっては動作しない。開発用としてファイルマネージャーを起動するには以下のコマンドを使う。

```
$ docker-compose up -d
```

## how to use

実行する。カレントディレクトリにある `upload.zip` を、指定パスにある Tiny File Manager の指定ディレクトリにアップロードする。`.`の場合は、Tiny File Manager のトップに設定されているディレクトリに`upload.zip`をアップロードする。

```
$ upl localhost:7777 .
upl [ベースパス] [アップロード先ディレクトリ]
```
( http://localhost:7777 で Tiny File Managerにアクセスできているものとする。)

## docker run

```
$ docker run -v "$PWD/":/work -w /work --rm -it ghcr.io/kijimad/upl:latest localhost:7777 .
```
