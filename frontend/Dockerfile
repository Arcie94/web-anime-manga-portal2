# SERVE STAGE (Requires local 'npm run build' first)
FROM node:20-slim

WORKDIR /app

# 1. Install Production Dependencies Only (Cached layer)
COPY package.json package-lock.json ./
RUN npm ci --omit=dev

# 2. Copy the pre-built application from Host
COPY dist ./dist

ENV HOST=0.0.0.0
ENV PORT=4321
EXPOSE 4321

CMD ["node", "./dist/server/entry.mjs"]
