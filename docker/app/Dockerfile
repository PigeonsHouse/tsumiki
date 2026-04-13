# Build the frontend
FROM node:24.14 AS build
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable

WORKDIR /app/frontend

COPY ./frontend/package.json ./frontend/pnpm-lock.yaml /app/frontend/

RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile

COPY ./frontend/ /app/frontend/
RUN pnpm build

# Build the backend
FROM golang:1.26 AS backend
WORKDIR /app/backend
COPY ./backend/go.mod ./backend/go.sum /app/backend/
RUN go mod download
COPY ./backend/ /app/backend/
RUN CGO_ENABLED=0 GOOS=linux go build -o api-server ./cmd/server.go

# Build the final image
FROM alpine:latest
WORKDIR /app
COPY --from=backend /app/backend/api-server /app
COPY --from=build /app/frontend/dist /app/view

CMD ["/app/api-server"]
