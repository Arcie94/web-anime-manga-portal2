# STAGE 1: Build
FROM node:20-slim AS builder
WORKDIR /app

# Copy package info and install ALL dependencies (including devDeps for building)
COPY package.json ./
# Do NOT copy package-lock from Windows host to avoid platform mismatch
RUN npm install

# Copy source code and build
COPY . .
RUN npm run build

# STAGE 2: Runtime
FROM node:20-slim AS runtime
WORKDIR /app

# Install only production deps for runtime
COPY package.json package-lock.json ./
RUN npm ci --omit=dev

# Copy built assets from builder stage
COPY --from=builder /app/dist ./dist

ENV HOST=0.0.0.0
ENV PORT=4321
EXPOSE 4321

CMD ["node", "./dist/server/entry.mjs"]
