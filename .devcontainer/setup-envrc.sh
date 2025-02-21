#!/bin/bash

# 権限を設定
chown -R vscode:vscode /workspace

# テンプレートファイルの存在確認
if [ -f "/workspace/.devcontainer/.envrc.template" ]; then
    # .envrcが存在しない場合、テンプレートをコピー
    if [ ! -f "/workspace/.envrc" ]; then
        cp /workspace/.devcontainer/.envrc.template /workspace/.envrc
        direnv allow /workspace/.envrc
    fi
    direnv allow
fi

echo $?
