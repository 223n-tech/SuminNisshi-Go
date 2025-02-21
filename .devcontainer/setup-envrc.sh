#!/bin/bash

# 権限を設定
chown -R vscode:vscode /workspace

# direnvが利用可能かチェック
if ! command -v direnv &> /dev/null; then
    echo "direnv command not found"
    exec sleep infinity
    exit 1
fi

# テンプレートファイルの存在確認
if [ ! -f "/workspace/.devcontainer/.envrc.template" ]; then
    echo ".envrc.template not found"
    exec sleep infinity
    exit 1
fi

# .envrcが存在しない場合、テンプレートをコピー
if [ ! -f "/workspace/.envrc" ]; then
    cp /workspace/.devcontainer/.envrc.template /workspace/.envrc
    direnv allow /workspace/.envrc
fi

direnv allow

# コンテナを起動し続ける
exec sleep infinity
