
# Todo Application - Docker Setup Guide

## Prerequisites
- Docker installed on your system
- Docker Compose installed

## Quick Start

### 1. Build Backend API Image
```bash
➤ cd backend
➤ docker build -t todo-api:v1.0.0 .
```
### 2. Build Frontend Image
```bash
➤ cd frontend
➤ docker build -t todo-web:v1.0.0 .
```
### 3. Run Docker Compose
```bash
➤ docker compose up -d
```
## Application URLs

-   **Frontend**: [http://localhost:80](http://localhost/)
    
-   **Backend API**: [http://localhost:5000](http://localhost:5000/)

## Environment Variables

The application uses the following environment variables:
### Frontend (.env)
```bash
VITE_API_URL=http://localhost:5000/api
```
### Backend

-   `DATABASE_URL`: Database connection string
    
-   `PORT`: Server port (default: 5000)

## Project Structure
```bash
todo-app/
├── frontend/          # React Vite application
├── backend/           # Golang API server
├── docker-compose.yml
└── README.md
```
