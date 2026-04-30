FROM golang:1.25.9-alpine AS builder

RUN apk add --no-cache --update git build-base
ENV CGO_ENABLED=0

WORKDIR /app

COPY . /app

WORKDIR /app/src
RUN go mod download
RUN go build -o /app/ezpay .

FROM alpine:3.22 AS runner
ENV TZ=Asia/Shanghai
RUN apk --no-cache add ca-certificates tzdata \
    && addgroup -S ezpay \
    && adduser -S -D -H -G ezpay ezpay
ARG API_RATE_URL=""

WORKDIR /app
COPY --from=builder --chown=ezpay:ezpay /app/src/www /app/www
COPY --from=builder --chown=ezpay:ezpay /app/src/.env.example /app/.env
RUN if [ -n "$API_RATE_URL" ]; then \
      sed -i "s|^api_rate_url=.*$|api_rate_url=${API_RATE_URL}|" /app/.env; \
    fi \
    && mkdir -p /app/conf /app/data/runtime /app/runtime /app/runtime/logs /app/static \
    && chown -R ezpay:ezpay /app
COPY --from=builder --chown=ezpay:ezpay /app/ezpay .

VOLUME /app/conf
USER ezpay
ENTRYPOINT ["./ezpay", "http", "start"]
