ARG GO_VERSION="1.20"

FROM golang:${GO_VERSION} AS build
WORKDIR /src
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY ./ ./
RUN go test -timeout 30s ./...
RUN go build \
    -ldflags "-linkmode 'external' -extldflags '-static'" \
    -o /app
RUN touch /db.sqlite3

FROM gcr.io/distroless/static-debian11:nonroot AS final
LABEL maintainer="travis@thisguy.codes"
USER nonroot:nonroot

COPY --from=build --chown=nonroot:nonroot /app /app
COPY --from=build --chown=nonroot:nonroot /db.sqlite3 /db.sqlite3

ENTRYPOINT [ "/app" ]
