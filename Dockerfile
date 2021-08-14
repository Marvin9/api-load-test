FROM golang:alpine AS root

WORKDIR /app/
COPY . .
RUN export CGO_ENABLED=0 && go build

FROM scratch
WORKDIR /app
COPY --from=root /app/ /app/

ENTRYPOINT [ "./api-load-test" ]