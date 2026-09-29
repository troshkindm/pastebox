FROM golang:1.27.1-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/pastebox ./cmd/pastebox

FROM alpine:3.22
COPY --from=build /out/pastebox /usr/local/bin/pastebox
USER nobody
EXPOSE 8080
CMD ["pastebox"]