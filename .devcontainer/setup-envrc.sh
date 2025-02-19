#!/bin/bash

# .envrcが存在しない場合、テンプレートをコピー
if [ ! -f "/workspace/.envrc" ]; then
    cp /workspace/.devcontainer/.envrc.template /workspace/.envrc
    direnv allow /workspace/.envrc
fi

# コンテナを起動し続ける
exec sleep infinity
