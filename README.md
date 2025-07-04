# Web Crawler Dashboard

A full-stack web application that crawls websites and displays analytics through a responsive dashboard interface.

## Tech Stack

- **Frontend:** React 18 + TypeScript + Tailwind CSS (Vite)
- **Backend:** Go 1.22 + Gin Framework + GORM
- **Database:** MySQL 8.0
- **Infrastructure:** Docker Multi-Stage Builds + Docker Compose

## Features

- URL management (add, analyze, delete)
- Website crawling and data extraction
- Real-time crawl status tracking
- Responsive data dashboard with sorting, filtering, and pagination
- Detailed analytics with charts
- Bulk operations support

## Quick Start

### Prerequisites

- Docker and Docker Compose installed
- Git

### Environment Commands

#### Development (with hot reloading)
```bash
git clone <repository-url>
cd web-crawler
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --build
```
- Frontend: http://localhost:3000 (hot reloading enabled)
- Backend: http://localhost:8080 (hot reloading enabled)
- Database: localhost:3306

#### Testing (automated tests)
```bash
docker-compose -f docker-compose.yml -f docker-compose.test.yml up --build
```
- Runs frontend and backend tests
- Uses isolated test database
- Generates test coverage reports

#### Production (optimized builds)
```bash
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up --build
```
- Frontend: http://localhost:80 (nginx served)
- Backend: http://localhost:8080 (compiled binary)
- Optimized for performance and security

### Development Workflow

The multi-stage Docker setup provides:

- **Hot Reloading:** Changes automatically refresh in development
- **Automated Testing:** Tests run in isolated environment
- **Production Builds:** Optimized multi-stage builds

### Stopping Applications

```bash
# Stop current environment
docker-compose down

# Remove all data (including database)
docker-compose down -v
```

## Project Structure

```
web-crawler/
├── docker-compose.yml          # Base configuration
├── docker-compose.dev.yml      # Development overrides
├── docker-compose.test.yml     # Testing overrides
├── docker-compose.prod.yml     # Production overrides
├── frontend/                   # React TypeScript application (Vite)
│   ├── Dockerfile             # Multi-stage frontend build
│   ├── package.json
│   ├── tailwind.config.js
│   ├── vite.config.ts         # Vite configuration with Vitest
│   ├── nginx.conf             # Production nginx config
│   └── src/
│       ├── App.test.tsx       # Vitest automated tests
│       └── ...
├── backend/                    # Go 1.22 API server
│   ├── Dockerfile             # Multi-stage backend build
│   ├── .air.toml              # Hot reloading config
│   ├── go.mod
│   └── ...
├── database/                   # Database initialization
│   └── init.sql
└── README.md
```

## Docker Architecture

### Multi-Stage Builds
- **Development:** Hot reloading with full toolchain
- **Testing:** Automated test execution
- **Production:** Minimal optimized images

### Environment Separation
- **Development:** Volume mounting for live code changes
- **Testing:** Isolated test database and CI mode
- **Production:** Security-hardened, minimal attack surface

## Development Status

This project follows modern Docker best practices with multi-environment support:

- [x] Multi-stage Docker builds
- [x] Environment-specific configurations
- [x] Hot reloading for development (Air + Vite)
- [x] Automated testing setup (Vitest + Go tests)
- [x] Production-optimized builds
- [ ] Database schema and migrations
- [ ] Backend API endpoints
- [ ] Frontend components and routing
- [ ] Web crawling functionality
- [ ] Dashboard with data visualization
- [ ] Real-time updates and bulk operations

## Testing

Frontend tests use Vitest + React Testing Library for happy-path scenarios:

```bash
# Run tests in Docker
docker-compose -f docker-compose.yml -f docker-compose.test.yml up --build

# Run tests locally (requires Node.js)
cd frontend && npm run test:run
```

Backend tests use Go's built-in testing framework:

```bash
# Tests run automatically in Docker testing environment
# Individual test execution coming in next commits
```

## API Endpoints

_Coming in next_

## Contributing

This project demonstrates modern full-stack development practices:
- Multi-stage Docker builds following official best practices
- Environment separation (dev/test/prod)
- Automated testing integration
- Professional deployment patterns
- Hot reloading development workflow