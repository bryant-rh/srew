FROM docker.io/library/golang:1.18-buster AS build-env

FROM build-env AS builder

WORKDIR /go/src
COPY ./ ./

ARG Github_UserName
ARG Github_Token
# build
RUN go env -w GOPRIVATE="github.com/kunlun-qilian" && \ 
    git config --global url."https://${Github_UserName}:${Github_Token}@github.com".insteadOf  "https://github.com" && \
    make build.srew.server WORKSPACE=server 

# runtime
FROM alpine
COPY --from=builder /go/src/cmd/server/server /go/bin/srew-server
#COPY --from=builder /go/src/cmd/server/openapi.json /go/bin/openapi.json
EXPOSE 80

ARG PROJECT_NAME
ARG PROJECT_VERSION
ENV GOENV=DEV PROJECT_NAME=${PROJECT_NAME} PROJECT_VERSION=${PROJECT_VERSION}

WORKDIR /go/bin
ENTRYPOINT ["/go/bin/srew-server"]