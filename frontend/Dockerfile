# Multi-stage Dockerfile
FROM node:20-alpine AS base

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy source code
COPY . .

# Development stage
FROM base AS development
EXPOSE 3000
ENV CHOKIDAR_USEPOLLING=true
ENV WATCHPACK_POLLING=true
CMD ["npm", "run", "dev", "--", "--host"]

# Testing stage
FROM base AS testing
# Install additional testing dependencies if needed
RUN npm install --save-dev @testing-library/jest-dom @testing-library/react @testing-library/user-event
EXPOSE 3000
CMD ["npm", "test"]

# Build stage
FROM base AS builder
RUN npm run build

# Production stage
FROM nginx:alpine AS production
# Copy built assets from builder stage
COPY --from=builder /app/dist /usr/share/nginx/html
# Copy nginx configuration
COPY nginx.conf /etc/nginx/nginx.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]