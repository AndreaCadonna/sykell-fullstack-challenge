# Deployment Documentation

## Overview

The Web Crawler application uses Docker multi-stage builds with environment-specific configurations for development, testing, and production deployments.

## Docker Architecture

### Multi-Stage Strategy

Each service (frontend, backend) uses multi-stage Dockerfiles to optimize for different environments:

- **Development**: Full toolchain with hot reload
- **Testing**: Automated test execution  
- **Production**: Minimal optimized images

### Container Images

| Service | Base Image | Purpose | Size |
|---------|------------|---------|------|
| Frontend Dev | node:20-alpine | Development with hot reload | ~500MB |
| Frontend Prod | nginx:alpine | Static file serving | ~50MB |
| Backend Dev | golang:1.22-alpine | Development with Air | ~400MB |
| Backend Prod | alpine:latest | Compiled binary only | ~20MB |
| Database | mysql:8.0 | Data persistence | ~200MB |

## Environment Configurations

### Development Environment

**File**: `docker-compose.dev.yml`

**Features**:
- Hot reload for both frontend and backend
- Volume mounting for live code changes
- Database exposed for debugging
- Full logging enabled

**Usage**:
```bash
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --build
```

**Services**:
- Frontend: http://localhost:3000 (Vite dev server)
- Backend: http://localhost:8080 (Air hot reload)
- Database: localhost:3306 (exposed port)

### Testing Environment  

**File**: `docker-compose.test.yml`

**Features**:
- Isolated test database
- CI-friendly test execution
- Coverage report generation
- No exposed ports

**Usage**:
```bash
docker-compose -f docker-compose.yml -f docker-compose.test.yml up --build
```

**Test Results**:
- Frontend tests: Vitest with React Testing Library
- Backend tests: Go testing framework
- Integration tests: End-to-end API testing

### Production Environment

**File**: `docker-compose.prod.yml`

**Features**:
- Optimized builds with minimal attack surface
- Security hardening
- Restart policies
- No development tools

**Usage**:
```bash
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up --build
```

**Services**:
- Frontend: http://localhost:80 (Nginx)
- Backend: http://localhost:8080 (compiled binary)
- Database: Internal network only

## Service Configuration

### Frontend Service

**Development Dockerfile**:
```dockerfile
FROM node:20-alpine AS development
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
ENV CHOKIDAR_USEPOLLING=true
ENV WATCHPACK_POLLING=true
EXPOSE 3000
CMD ["npm", "run", "dev", "--", "--host"]
```

**Production Dockerfile**:
```dockerfile
FROM nginx:alpine AS production
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

**Environment Variables**:
- `CHOKIDAR_USEPOLLING=true`: Docker file watching
- `WATCHPACK_POLLING=true`: Vite hot reload
- `REACT_APP_API_URL`: Backend API endpoint

### Backend Service

**Development Configuration**:
```dockerfile
FROM golang:1.22-alpine AS development
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go install github.com/cosmtrek/air@v1.49.0
EXPOSE 8080
CMD ["air"]
```

**Production Configuration**:
```dockerfile
FROM alpine:latest AS production
WORKDIR /app
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
```

**Environment Variables**:
- `DB_HOST=database`: Database hostname
- `DB_PORT=3306`: Database port
- `DB_USER=crawler_user`: Database username
- `DB_PASSWORD=crawler_password`: Database password
- `DB_NAME=crawler_db`: Database name
- `JWT_SECRET`: Token signing secret
- `ENV`: Environment mode (development/production)

### Database Service

**Configuration**:
```yaml
database:
  image: mysql:8.0
  environment:
    - MYSQL_ROOT_PASSWORD=root_password
    - MYSQL_DATABASE=crawler_db
    - MYSQL_USER=crawler_user
    - MYSQL_PASSWORD=crawler_password
  volumes:
    - mysql_data:/var/lib/mysql
    - ./database/init.sql:/docker-entrypoint-initdb.d/init.sql
  healthcheck:
    test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
    timeout: 20s
    retries: 10
    interval: 10s
    start_period: 40s
```

## Networking

### Service Communication

```yaml
# Internal Docker network
networks:
  default:
    driver: bridge

# Service dependencies
services:
  backend:
    depends_on:
      database:
        condition: service_healthy
  
  frontend:
    depends_on:
      - backend
```

### Port Mapping

| Environment | Frontend | Backend | Database |
|-------------|----------|---------|----------|
| Development | 3000:3000 | 8080:8080 | 3306:3306 |
| Testing | No ports | No ports | No ports |
| Production | 80:80 | 8080:8080 | Internal only |

### CORS Configuration

```go
cors.Config{
    AllowOrigins: []string{
        "http://localhost:3000",  // Development
        "http://localhost:80",    // Production
    },
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    AllowCredentials: true,
    MaxAge:           12 * time.Hour,
}
```

## Data Persistence

### Volume Management

**Development Volumes**:
```yaml
volumes:
  - ./frontend:/app                    # Live code editing
  - /app/node_modules                  # Preserve dependencies
  - ./backend:/app                     # Live code editing
  - /app/tmp                          # Air build cache
  - mysql_data:/var/lib/mysql         # Database persistence
```

**Production Volumes**:
```yaml
volumes:
  - mysql_data:/var/lib/mysql         # Database only
```

### Backup Strategy

**Development**:
```bash
# Backup database
docker exec <container> mysqldump -u crawler_user -p crawler_db > backup.sql

# Restore database
docker exec -i <container> mysql -u crawler_user -p crawler_db < backup.sql
```

**Production**:
```bash
# Automated backup script
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
docker exec mysql-container mysqldump \
  -u crawler_user -pcrawler_password crawler_db \
  > "backup_${DATE}.sql"
```

## Health Checks

### Application Health

**Backend Health Check**:
```http
GET /health
Response: {
  "status": "healthy",
  "timestamp": "2025-07-04T13:00:00Z",
  "service": "web-crawler-api"
}
```

**Database Health Check**:
```yaml
healthcheck:
  test: ["CMD", "mysqladmin", "ping", "-h", "localhost", "-u", "crawler_user", "-pcrawler_password"]
  timeout: 20s
  retries: 10
  interval: 10s
  start_period: 40s
```

### Monitoring

**Container Health**:
```bash
# Check container status
docker-compose ps

# View container logs
docker-compose logs <service>

# Monitor resource usage
docker stats
```

## Security

### Container Security

**Non-Root Users**:
```dockerfile
# Backend production
FROM alpine:latest
RUN addgroup -g 1001 -S nodejs
RUN adduser -S nextjs -u 1001
USER nextjs
```

**Minimal Attack Surface**:
- Alpine Linux base images
- Only necessary packages installed
- No development tools in production
- Read-only file systems where possible

### Network Security

**Internal Communication**:
- Services communicate via internal Docker network
- Database not exposed in production
- API Gateway for external access (future)

**HTTPS Termination**:
```nginx
# Production nginx configuration
server {
    listen 443 ssl http2;
    server_name your-domain.com;
    
    ssl_certificate /etc/ssl/certs/cert.pem;
    ssl_certificate_key /etc/ssl/private/key.pem;
    
    location / {
        root /usr/share/nginx/html;
        try_files $uri $uri/ /index.html;
    }
    
    location /api/ {
        proxy_pass http://backend:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### Secrets Management

**Environment Variables**:
```yaml
# Production secrets
environment:
  - DB_PASSWORD=${DB_PASSWORD}
  - JWT_SECRET=${JWT_SECRET}
```

**Docker Secrets** (Future):
```yaml
secrets:
  db_password:
    external: true
  jwt_secret:
    external: true
```

## CI/CD Pipeline

### GitHub Actions Example

```yaml
name: Deploy
on:
  push:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Run Tests
        run: |
          docker-compose -f docker-compose.yml -f docker-compose.test.yml up --build --abort-on-container-exit

  deploy:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Deploy to Production
        run: |
          docker-compose -f docker-compose.yml -f docker-compose.prod.yml up --build -d
```

### Build Optimization

**Docker Build Cache**:
```bash
# Use BuildKit for faster builds
export DOCKER_BUILDKIT=1
docker-compose build --parallel
```

**Multi-Platform Builds**:
```bash
# Build for different architectures
docker buildx build --platform linux/amd64,linux/arm64 .
```

## Scaling

### Horizontal Scaling

**Load Balancer Configuration**:
```yaml
# nginx-lb.conf
upstream backend {
    server backend1:8080;
    server backend2:8080;
    server backend3:8080;
}

server {
    location /api/ {
        proxy_pass http://backend;
    }
}
```

**Database Replication**:
```yaml
# Master-slave MySQL setup
mysql-master:
  image: mysql:8.0
  environment:
    - MYSQL_REPLICATION_MODE=master

mysql-slave:
  image: mysql:8.0
  environment:
    - MYSQL_REPLICATION_MODE=slave
    - MYSQL_MASTER_HOST=mysql-master
```

### Vertical Scaling

**Resource Limits**:
```yaml
services:
  backend:
    deploy:
      resources:
        limits:
          cpus: '2.0'
          memory: 1G
        reservations:
          cpus: '1.0'
          memory: 512M
```

## Production Deployment

### Cloud Deployment Options

**AWS ECS**:
```yaml
# ecs-task-definition.json
{
  "family": "web-crawler",
  "containerDefinitions": [
    {
      "name": "frontend",
      "image": "your-repo/web-crawler-frontend:latest",
      "memory": 512,
      "portMappings": [{"containerPort": 80}]
    }
  ]
}
```

**Kubernetes**:
```yaml
# k8s-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web-crawler-backend
spec:
  replicas: 3
  selector:
    matchLabels:
      app: web-crawler-backend
  template:
    spec:
      containers:
      - name: backend
        image: web-crawler-backend:latest
        ports:
        - containerPort: 8080
```

**Digital Ocean App Platform**:
```yaml
# .do/app.yaml
name: web-crawler
services:
- name: frontend
  source_dir: /frontend
  github:
    repo: your-username/web-crawler
    branch: main
  run_command: npm start
```

### SSL/TLS Configuration

**Let's Encrypt with Nginx**:
```bash
# Install certbot
docker run --rm -it \
  -v letsencrypt:/etc/letsencrypt \
  -v $(pwd)/nginx:/var/www/html \
  certbot/certbot certonly --webroot \
  -w /var/www/html \
  -d your-domain.com
```

**Nginx SSL Configuration**:
```nginx
server {
    listen 443 ssl http2;
    ssl_certificate /etc/letsencrypt/live/your-domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/your-domain.com/privkey.pem;
    
    # Security headers
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Frame-Options DENY always;
    add_header X-Content-Type-Options nosniff always;
}
```

## Monitoring and Logging

### Application Logging

**Structured Logging**:
```json
{
  "timestamp": "2025-07-04T13:00:00Z",
  "level": "info",
  "service": "web-crawler-api",
  "message": "Request processed",
  "request_id": "req-123",
  "duration_ms": 245
}
```

**Log Aggregation**:
```yaml
# ELK Stack example
logging:
  driver: "fluentd"
  options:
    fluentd-address: localhost:24224
    tag: web-crawler.{{.Name}}
```

### Metrics Collection

**Prometheus Integration**:
```go
// Future backend metrics
import "github.com/prometheus/client_golang/prometheus"

var (
    requestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
)
```

**Grafana Dashboard**:
- HTTP request rates
- Response time percentiles
- Error rates by endpoint
- Database connection pool usage

## Troubleshooting

### Common Issues

**Container Won't Start**:
```bash
# Check container logs
docker-compose logs <service>

# Inspect container configuration
docker inspect <container-id>

# Check resource usage
docker stats
```

**Database Connection Issues**:
```bash
# Test database connectivity
docker exec -it backend-container sh
nc -zv database 3306

# Check database logs
docker-compose logs database
```

**Permission Issues**:
```bash
# Fix volume permissions
sudo chown -R $USER:$USER ./data
chmod -R 755 ./data
```

### Debug Commands

**Container Debugging**:
```bash
# Enter running container
docker exec -it <container> sh

# Run container with shell override
docker run -it --entrypoint sh <image>

# Check environment variables
docker exec <container> env
```

**Network Debugging**:
```bash
# List networks
docker network ls

# Inspect network
docker network inspect <network-name>

# Test connectivity between containers
docker exec <container1> ping <container2>
```

## Maintenance

### Updates and Upgrades

**Rolling Updates**:
```bash
# Update with zero downtime
docker-compose pull
docker-compose up -d --no-deps <service>
```

**Database Migrations**:
```bash
# Backup before migration
./scripts/backup-db.sh

# Run migrations
docker exec backend-container ./migrate up

# Verify migration
docker exec database-container mysql -e "SHOW TABLES"
```

**Security Updates**:
```bash
# Update base images
docker pull node:20-alpine
docker pull golang:1.22-alpine
docker pull mysql:8.0

# Rebuild with updated images
docker-compose build --no-cache
```

### Backup Procedures

**Automated Backup Script**:
```bash
#!/bin/bash
# backup.sh
set -e

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/backups"

# Database backup
docker exec mysql-container mysqldump \
  -u crawler_user -pcrawler_password \
  --single-transaction \
  --routines \
  --triggers \
  crawler_db > "${BACKUP_DIR}/db_${DATE}.sql"

# Compress backup
gzip "${BACKUP_DIR}/db_${DATE}.sql"

# Upload to S3 (optional)
aws s3 cp "${BACKUP_DIR}/db_${DATE}.sql.gz" \
  s3://your-backup-bucket/

# Cleanup old backups (keep last 30 days)
find ${BACKUP_DIR} -name "db_*.sql.gz" -mtime +30 -delete
```

**Disaster Recovery**:
```bash
# Full system restore
./scripts/restore-backup.sh backup_20250704_130000.sql.gz

# Test restoration in staging
docker-compose -f docker-compose.test.yml up -d
./scripts/restore-backup.sh backup.sql.gz
./scripts/run-tests.sh
```