# syntax=docker/dockerfile:1

# --- Build Stage ---
ARG NODE_VERSION=22.13.1
FROM node:${NODE_VERSION}-slim AS builder
WORKDIR /app

# Install dependencies (using cache for npm)
COPY --link package.json package-lock.json ./
RUN --mount=type=cache,target=/root/.npm \
    npm ci

# Copy source files (excluding .git, .env, lock files, etc.)
COPY --link . .

# Build the app
RUN npm run build

# --- Production Stage ---
FROM node:${NODE_VERSION}-slim AS final
WORKDIR /app

# Create non-root user
RUN addgroup --system appuser && adduser --system --ingroup appuser appuser

# Copy built app and production deps only
COPY --from=builder /app/package.json ./
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/.svelte-kit ./build
# COPY --from=builder /app/dist ./dist
COPY --from=builder /app/static ./static

ENV NODE_ENV=production
ENV NODE_OPTIONS="--max-old-space-size=4096"
USER appuser

EXPOSE 5173
CMD ["npm", "run", "preview", "--", "--port", "5173", "--host"]
