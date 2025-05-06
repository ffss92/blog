# Stage 1: Bundle assets
FROM node:22-slim AS assets
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm install --frozen-lockfile
COPY . .
RUN npm run tw:build

# Stage 2: Build server
FROM golang:1.24-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/server cmd/server/*.go

# Stage 3: Deploy
FROM alpine:latest
COPY --from=build /out/server /server
COPY --from=assets /app/web /web
COPY --from=assets /app/articles /articles
EXPOSE 4000
ENTRYPOINT [ "/server" ]
