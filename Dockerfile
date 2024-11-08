
FROM golang:1.22.5-alpine3.19 AS builder


WORKDIR /app


COPY go.mod go.sum ./


COPY . .


ENV GOCACHE=/root/.cache/go-build


RUN --mount=type=cache,target="/root/.cache/go-build" go build -o /pdf_generator


FROM alpine:3.14 AS runner


WORKDIR /data


COPY --from=builder /pdf_generator /bin/pdf_generator

COPY --from=builder /app/internal /app/internal

COPY --from=builder /app/qrcode.png /data/qrcode.png


EXPOSE 2222


CMD ["/bin/pdf_generator"]
