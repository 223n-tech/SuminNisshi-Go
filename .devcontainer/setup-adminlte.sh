#!/bin/bash

# 必要なディレクトリを作成
mkdir -p /workspace/web/static/adminlte
mkdir -p /workspace/web/template
mkdir -p /workspace/web/template/layouts
mkdir -p /workspace/web/template/partials
mkdir -p /workspace/web/template/pages

# AdminLTE 最新版（v3.2.0）をダウンロード
wget https://github.com/ColorlibHQ/AdminLTE/archive/refs/tags/v3.2.0.zip -O adminlte.zip

# 解凍
unzip adminlte.zip

# 必要なファイルを配置
cp -r AdminLTE-3.2.0/dist/* /workspace/web/static/adminlte/
cp -r AdminLTE-3.2.0/pages/* /workspace/web/template/
cp -r AdminLTE-3.2.0/plugins /workspace/web/static/adminlte/

# 不要なファイルを削除
rm -rf adminlte.zip AdminLTE-3.2.0

# 権限を設定
chown -R vscode:vscode /workspace/web
