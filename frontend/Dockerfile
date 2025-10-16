FROM node:20-alpine AS builder

WORKDIR /app

COPY package*.json ./

RUN npm ci

COPY . .

# Build with limited memory to prevent OOM
ENV NODE_OPTIONS="--max-old-space-size=4096"
RUN npm run build

FROM node:20-alpine

WORKDIR /app

# Copy built application
COPY --from=builder /app/build ./build
COPY --from=builder /app/package*.json ./

# Install ONLY production dependencies (for SvelteKit runtime)
RUN npm ci --omit=dev

# Expose port
EXPOSE 3000

# Runtime environment
ENV NODE_ENV=production
ENV PORT=3000
ENV HOST=0.0.0.0

# Start the SvelteKit server
CMD ["node", "build"]
