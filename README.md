# Database Hosting Services (DBHS) API

[![Go](https://img.shields.io/badge/Go-1.24.4-blue.svg)](https://golang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-blue.svg)](https://postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-7+-red.svg)](https://redis.io/)
[![Docker](https://img.shields.io/badge/Docker-Ready-blue.svg)](https://docker.com/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

A comprehensive, production-ready REST API for managing database hosting services. DBHS provides a complete solution for database-as-a-service operations, including user management, project creation, schema management, analytics, and AI-powered assistance.

## 🚀 Features

### Core Functionality
- **🔐 User Authentication & Authorization**: JWT-based authentication with role-based access control
- **📊 Project Management**: Create, manage, and organize database projects
- **🗄️ Database Schema Management**: Dynamic schema creation, modification, and inspection
- **📋 Table Operations**: Full CRUD operations for database tables with real-time schema synchronization
- **🔍 SQL Editor**: Execute custom SQL queries with syntax validation and error handling
- **📈 Analytics & Monitoring**: Real-time database usage statistics, storage monitoring, and cost analysis
- **🤖 AI Assistant**: Intelligent database assistant powered by advanced AI for query optimization and suggestions
- **⚡ Caching**: Redis-based caching for improved performance
- **🔒 Security**: Rate limiting, input validation, and secure database connections

### Advanced Features
- **📊 Real-time Analytics**: Track database usage, query performance, and storage metrics
- **🔄 Background Workers**: Automated data collection and maintenance tasks
- **📖 API Documentation**: Auto-generated Swagger/OpenAPI documentation
- **🐳 Containerized Deployment**: Docker and Fly.io ready
- **📧 Email Integration**: User verification and notification system
- **🛡️ Middleware Stack**: Comprehensive security and logging middleware

## 🏗️ Architecture

![Database Design](./public/database-design.png)

The DBHS API follows a modular microservices-inspired architecture with the following components:

### System Architecture
- **API Gateway Layer**: Gorilla Mux router with comprehensive middleware
- **Service Layer**: Business logic organized by domain (accounts, projects, tables, etc.)
- **Data Access Layer**: PostgreSQL with connection pooling and Redis caching
- **Background Processing**: Cron-based workers for analytics and maintenance
- **External Integrations**: AI services, email notifications, and monitoring

### Database Strategy
- **Metadata Database**: Stores user accounts, project configurations, and system metadata
- **Dynamic User Databases**: Isolated databases created per user project for data isolation
- **Connection Pooling**: Efficient database connection management with pgxpool
- **Multi-tenancy**: Secure isolation between user projects and data

## 📁 Project Structure

```
├── main/                    # Application entry point and routing
├── config/                  # Configuration management and database connections
├── accounts/                # User authentication and profile management
├── projects/                # Project creation and management
├── tables/                  # Database table operations and schema management
├── schemas/                 # Database schema inspection and metadata
├── SqlEditor/               # SQL query execution and validation
├── AI/                      # AI assistant and intelligent features
├── analytics/               # Usage analytics and monitoring
├── indexes/                 # Database indexing operations
├── middleware/              # HTTP middleware (auth, rate limiting, CORS)
├── response/                # Standardized API response handling
├── utils/                   # Shared utilities and helpers
├── workers/                 # Background processing and cron jobs
├── caching/                 # Redis caching implementation
├── test/                    # Comprehensive test suites
├── docs/                    # API documentation (Swagger/OpenAPI)
├── public/                  # Static assets and diagrams
└── templates/               # Email and notification templates
```

## 🛠️ Technology Stack

### Backend Technologies
- **Language**: Go 1.24.4
- **Web Framework**: Gorilla Mux
- **Database**: PostgreSQL 15+ with pgx driver
- **Caching**: Redis 7+
- **Authentication**: JWT with golang-jwt/jwt
- **Documentation**: Swagger/OpenAPI with go-swaggo
- **Background Jobs**: Robfig Cron
- **Email**: GoMail v2

### AI & Analytics
- **AI Integration**: Custom AI agent for database assistance
- **Monitoring**: Real-time analytics and usage tracking
- **Performance**: Query optimization and caching strategies

## 🚀 Quick Start

### Prerequisites
- Go 1.24.4 or higher
- PostgreSQL 15+
- Redis 7+
- Docker (optional)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/your-org/database-hosting-services-api.git
   cd database-hosting-services-api
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Required Environment Variables**
   ```env
   # Database Configuration
   DATABASE_URL=postgresql://user:password@localhost:5432/dbhs_main
   DATABASE_ADMIN_URL=postgresql://admin:password@localhost:5432/postgres
   TEST_DATABASE_URL=postgresql://user:password@localhost:5432/dbhs_test
   
   # Redis Configuration
   REDIS_URL=redis://localhost:6379
   
   # JWT Configuration
   JWT_SECRET=your-super-secret-jwt-key
   
   # Email Configuration
   SMTP_HOST=smtp.gmail.com
   SMTP_PORT=587
   SMTP_USERNAME=your-email@gmail.com
   SMTP_PASSWORD=your-app-password
   
   # AI Configuration
   AI_API_KEY=your-ai-api-key
   
   # Application Configuration
   API_PORT=8000
   ENV=development
   ```

5. **Run database migrations**
   ```bash
   # Set up main database schema
   psql -d $DATABASE_URL -f scripts/migrations/001_initial_schema.sql
   ```

6. **Build and run the application**
   ```bash
   make build
   make run
   ```

## 📖 API Documentation

### Interactive Documentation
- **Scalar API Reference**: https://orbix.fly.dev/reference
- **Swagger UI**: https://orbix.fly.dev/swagger/index.html
- **ReDoc**: https://orbix.fly.dev/redoc

### Generate Documentation
```bash
make generate-docs
```

## 🔗 API Endpoints

### Authentication
- `POST /api/user/sign-up` - User registration
- `POST /api/user/sign-in` - User login  
- `POST /api/user/verify` - Email verification
- `POST /api/user/resend-code` - Resend verification code
- `POST /api/user/forget-password` - Password reset request
- `POST /api/user/forget-password/verify` - Password reset verification

### User Management
- `GET /api/users/me` - Get current user profile
- `POST /api/users/update-password` - Update user password
- `PATCH /api/users/{id}` - Update user profile

### Projects
- `GET /api/projects` - List user projects
- `POST /api/projects` - Create new project
- `GET /api/projects/{project_id}` - Get project details
- `PATCH /api/projects/{project_id}` - Update project
- `DELETE /api/projects/{project_id}` - Delete project

### Tables & Schema
- `GET /api/projects/{project_id}/tables` - List project tables
- `POST /api/projects/{project_id}/tables` - Create new table
- `GET /api/projects/{project_id}/tables/{table_id}` - Get table data (with pagination)
- `POST /api/projects/{project_id}/tables/{table_id}` - Insert row into table
- `PUT /api/projects/{project_id}/tables/{table_id}` - Update table schema
- `DELETE /api/projects/{project_id}/tables/{table_id}` - Delete table
- `GET /api/projects/{project_id}/tables/{table_id}/schema` - Get table schema

### Database Schema
- `GET /api/projects/{project_id}/schema/tables` - Get database schema
- `GET /api/projects/{project_id}/schema/tables/{table_id}` - Get specific table schema

### Indexes
- `GET /api/projects/{project_id}/indexes` - List project indexes
- `POST /api/projects/{project_id}/indexes` - Create new index
- `GET /api/projects/{project_id}/indexes/{index_oid}` - Get specific index
- `PUT /api/projects/{project_id}/indexes/{index_oid}` - Update index name
- `DELETE /api/projects/{project_id}/indexes/{index_oid}` - Delete index

### SQL Editor
- `POST /api/projects/{project_id}/sqlEditor/run-query` - Execute SQL query

### Analytics
- `GET /api/projects/{project_id}/analytics/storage` - Database storage analytics
- `GET /api/projects/{project_id}/analytics/execution-time` - Query execution time statistics
- `GET /api/projects/{project_id}/analytics/usage` - Database usage statistics

### AI Assistant
- `GET /api/projects/{project_id}/ai/report` - Get AI-generated database report
- `POST /api/projects/{project_id}/ai/chatbot/ask` - Chat with AI assistant
- `POST /api/projects/{project_id}/ai/agent` - AI agent for database operations
- `POST /api/projects/{project_id}/ai/agent/accept` - Accept AI agent suggestions
- `POST /api/projects/{project_id}/ai/agent/cancel` - Cancel AI agent operations


## 🔧 Development

## 🧪 Testing

### Run Tests
```bash
# Run all tests
make test

# Run specific test suite
go test ./test/accounts_test.go -v
go test ./test/projects_test.go -v
go test ./test/tables/ -v

# Run tests with coverage
go test -cover ./...
```

### Test Structure
- **Unit Tests**: Individual component testing
- **Integration Tests**: End-to-end API testing
- **Service Tests**: Business logic validation
- **Database Tests**: Data access layer testing

## 📊 Monitoring & Analytics

### Built-in Analytics
- **Database Usage**: Query counts, execution time, resource usage
- **Storage Metrics**: Database size, growth trends, optimization suggestions
- **Performance Monitoring**: Query performance, connection pooling metrics
- **User Analytics**: Usage patterns, feature adoption

### External Integrations
- **Axiom**: Advanced logging and analytics
- **Health Checks**: Application and database health monitoring
- **Alerting**: Automated alerts for critical issues

## 🔒 Security Features

### Authentication & Authorization
- **JWT Tokens**: Secure, stateless authentication
- **Role-Based Access**: Fine-grained permission control
- **Project Isolation**: Secure multi-tenancy

### Security Measures
- **Rate Limiting**: API abuse prevention
- **Input Validation**: SQL injection and XSS protection
- **CORS Configuration**: Cross-origin request security
- **Database Security**: Connection encryption and access control