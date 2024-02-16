FROM golang:1.22 as BUILD
WORKDIR /src
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux make build

FROM alpine:latest
COPY --from=BUILD /src/favsapi /usr/bin
ENTRYPOINT ["favsapi"]