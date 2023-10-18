# [WIP]upl

upload zip file script.

ブラウザでのファイルアップロードが完全に制限されているクソ環境でファイル転送させるために作った。Tiny File ManagerのエンドポイントにcURLで直にリクエストし、クッキー取得とアップロードを一度にやってくれる。

## install

```
$ go install github.com/kijimaD/upl@main
```

## docker run

```
$ docker run -v "$PWD/":/work -w /work --rm -it ghcr.io/kijimad/upl:latest
```
