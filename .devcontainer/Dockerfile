FROM ubuntu:22.04

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

ARG TARGETARCH
ARG GO_VERSION=1.20
ARG GOPLS_VERSION=v0.11.0
ARG KUBECTL_VERSION=v1.27.3
ARG KIND_VERSION=v
ARG SHFMT_VERSION=v3.7.0
ARG YQ_VERSION=v4.30.6

ENV LANG='en_US.UTF-8' LANGUAGE='en_US:en' LC_ALL='en_US.UTF-8'
ENV PATH=/usr/local/go/bin:$GOPATH/bin:$PATH
ENV GOPATH=$HOME/go
# Grafana Tempo のライブラリが v2 対応していないため、+incompatible で回避
# https://christina04.hatenablog.com/entry/go-modules-major-version
ENV GOPRIVATE="github.com/grafana/tempo"
ENV EDITOR=vim

RUN apt-get update && apt-get install -y \
	bash-completion \
	build-essential \
	ca-certificates \
	cmake \
	curl \
	gnupg \
	htop \
	jq \
	less \
	locales \
	lsb-release \
	lsof \
	man-db \
	neovim \
	shellcheck \
	ssl-cert \
	sudo \
	tmux \
	unzip \
	xz-utils \
	zip \
	&& apt-get autoremove -y \
	&& apt-get clean \
	&& rm -rf /var/cache/* /var/lib/apt/lists/* \
	# install golang and toolchains
	&& echo curl -fsSL https://dl.google.com/go/go$GO_VERSION.linux-$TARGETARCH.tar.gz \
	&& curl -fsSL https://dl.google.com/go/go$GO_VERSION.linux-$TARGETARCH.tar.gz | tar zx -C /usr/local \
	&& echo 'export PATH=/usr/local/go/bin:$GOPATH/bin:$PATH' >> $HOME/.bashrc \
	&& go install -v golang.org/x/tools/gopls@$GOPLS_VERSION \
	&& go install -v mvdan.cc/sh/v3/cmd/shfmt@$SHFMT_VERSION \
	&& go install -v github.com/mikefarah/yq/v4@$YQ_VERSION \
	# install kubectl and toolchains
	&& curl -LO "https://dl.k8s.io/release/$KUBECTL_VERSION/bin/linux/$TARGETARCH/kubectl" \
	# configure locales to UTF8
	&& locale-gen en_US.UTF-8 \
	# install Starship
	&& curl -sS https://starship.rs/install.sh | sh -s -- -y \
	&& echo 'eval "$(starship init bash)"' >> $HOME/.bashrc
