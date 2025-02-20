FROM debian:bookworm-slim AS base

# 環境変数の設定
ENV GOLANG_VERSION=1.22.1
ENV PATH=/usr/local/go/bin:$PATH
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH

# vscodeユーザーの作成
RUN useradd -l -u 1000 -G sudo -md /home/vscode -s /bin/bash -p vscode vscode \
    && groupadd docker \
    && usermod -aG docker vscode

# 基本パッケージのインストール
RUN apt-get update && apt-get install -y \
    git \
    curl \
    wget \
    gcc \
    g++ \
    make \
    direnv \
    sudo \
    unzip \
    tree \
    && rm -rf /var/lib/apt/lists/*

# sudoers設定
RUN echo "vscode ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers.d/vscode \
    && chmod 0440 /etc/sudoers.d/vscode

# Goのインストール
RUN curl -O https://dl.google.com/go/go${GOLANG_VERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go${GOLANG_VERSION}.linux-amd64.tar.gz && \
    rm go${GOLANG_VERSION}.linux-amd64.tar.gz

# CompileDaemon（ホットリロード用）のインストール
RUN go install github.com/githubnemo/CompileDaemon@latest

# direnvの設定
RUN echo 'eval "$(direnv hook bash)"' >> /home/vscode/.bashrc

# 作業ディレクトリの設定
WORKDIR /workspace

# Go モジュールのキャッシュディレクトリ作成とパーミッション設定
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && \
    chown -R vscode:vscode "$GOPATH" && \
    chown -R vscode:vscode /workspace

# .envrcのセットアップスクリプトを追加
COPY setup-envrc.sh /usr/local/bin/
COPY setup-adminlte.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/setup-envrc.sh && \
    chmod +x /usr/local/bin/setup-adminlte.sh && \
    chown vscode:vscode /usr/local/bin/setup-envrc.sh && \
    chown vscode:vscode /usr/local/bin/setup-adminlte.sh

# デフォルトユーザーをvscodeに設定
USER vscode

# コンテナ起動時の初期化
CMD ["/usr/local/bin/setup-envrc.sh"]

