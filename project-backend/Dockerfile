# Build stage
FROM golang:1.23 AS builder


WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o server

# Final stage
FROM debian:bullseye-slim

WORKDIR /app
COPY --from=builder /app/server /app/
COPY --from=builder /app/templates /app/templates
COPY --from=builder /app/admin_templates /app/admin_templates
COPY --from=builder /app/fonts /app/fonts
COPY --from=builder /app/output /app/output
COPY --from=builder /app/uploads /app/uploads
COPY --from=builder /app/frontend /app/frontend

RUN apt update && apt install -y wkhtmltopdf fontconfig libxrender1 libxext6 xfonts-base && rm -rf /var/lib/apt/lists/*
EXPOSE 8090

CMD ["./server"]
