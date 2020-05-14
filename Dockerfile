FROM golang:1.13 AS builder

ENV GO111MODULE on
ENV GOPRIVATE github.com/AlpacaLabs

ARG GITHUB_USER
ARG GITHUB_PASS

COPY go.mod go.sum /go/app/
WORKDIR /go/app

# 1. add credentials on build
# 2. make sure your domain is accepted
# 3. configure git to use ssh instead of https
# 4. download dependencies
ARG SSH_PRIVATE_KEY
RUN mkdir /root/.ssh/ && \
    echo "${SSH_PRIVATE_KEY}" > /root/.ssh/id_rsa && \
    chmod 600 /root/.ssh/id_rsa && \
    touch /root/.ssh/known_hosts && \
    ssh-keyscan github.com >> /root/.ssh/known_hosts && \
    git config --global url.git@github.com:.insteadOf https://github.com/ && \
    go mod download

COPY . /go/app

CMD ["go", "run", "main.go"]