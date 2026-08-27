FROM ubuntu:22.04

# Общее
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
      ca-certificates locales tzdata \
      curl wget git \
      mc locate file \
      htop psmisc \
      vim nano \
      build-essential software-properties-common p7zip-full \
    && rm -rf /var/lib/apt/lists/*


# Версия Go
ENV GOLANG_VERSION=1.26.0

# Ставим Go под текущую архитектуру (arm64 на M1)
# На amd64 автоматически подтянется amd64-архив, на arm64 — arm64.
RUN set -eux; \
    arch="$(dpkg --print-architecture)"; \
    case "$arch" in \
      amd64) goarch="amd64" ;; \
      arm64) goarch="arm64" ;; \
      armhf) goarch="armv6l" ;; \
      *) echo "Unsupported arch: $arch" && exit 1 ;; \
    esac; \
    wget -O /tmp/go.tgz "https://go.dev/dl/go${GOLANG_VERSION}.linux-${goarch}.tar.gz"; \
    tar -C /usr/local -xzf /tmp/go.tgz; \
    rm /tmp/go.tgz

# PATH и GOPATH
ENV GOPATH=/go
ENV PATH="/usr/local/go/bin:${GOPATH}/bin:${PATH}"

# Установка утилит Go (соберутся нативно под arm64)
RUN go install github.com/air-verse/air@latest && \
    go install github.com/swaggo/swag/cmd/swag@latest
