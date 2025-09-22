# syntax=docker/dockerfile:1

ARG NODE_VERSION=24.7.0
FROM node:${NODE_VERSION}-slim AS builder
WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY . .
RUN npm run build

FROM node:${NODE_VERSION}-slim AS final
WORKDIR /app

RUN addgroup --system appuser && adduser --system --ingroup appuser appuser

COPY --from=builder /app/ ./
# COPY --from=builder /app/node_modules ./node_modules
# COPY --from=builder /app/.svelte-kit ./.svelte-kit

ENV NODE_ENV=production
USER appuser

EXPOSE 3000
CMD ["node", "build"]
