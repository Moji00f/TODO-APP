
# Todo Application - Docker Setup Guide

## Prerequisites
- Docker installed on your system
- Docker Compose installed

## Quick Start

### 1. Build Backend API Image
```bash
➤ cd backend
➤ docker build -t 192.168.240.141:5000/todo-api:v1.1.0 .
```
### 2. Build Frontend Image
```bash
➤ cd frontend
➤ docker build -t 192.168.240.141:5000/todo-web:v1.1.0 .

➤ dockeer push 192.168.240.141:5000/todo-api:v1.1.0
➤ docker push 192.168.240.141:5000/todo-web:v1.1.0
```

### 3. Run Docker Stack
```bash
➤ docker stack deploy -c docker-compose.yaml GoReact
```
## Application URLs

-   **Frontend**: [http://goreact.local](http://goreact.local/)
    
-   **Backend API**: [http://api.goreact.local](http://api.goreact.local/)

## Environment Variables

The application uses the following environment variables:
### Frontend (.env)
```bash
VITE_API_URL=http://api.goreact.local
```
### Backend

-   `DATABASE_URL`: Database connection string
    
-   `PORT`: Server port (default: 4000)

## Project Structure
```bash
todo-app/
├── frontend/          # React Vite application
├── backend/           # Golang API server
├── nginx         
├── docker-compose.yml
└── README.md
```
