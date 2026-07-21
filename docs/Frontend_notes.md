# Frontend Notes

This folder is for the frontend application.

The backend is already set up with:

```text
Authentication
Service management
Health check history
Dashboard summary data
CORS support
```

## Recommended Stack

Use one of these:

```text
React + Vite
Angular
```

Recommended if using React:

```text
React + Vite
Tailwind CSS
React Router
Axios or Fetch API
Recharts
```

## Backend URL

The backend runs on:

```text
http://localhost:8080
```

If using React + Vite, create a frontend `.env` file:

```env
VITE_API_BASE_URL=http://localhost:8080
```

## CORS

The backend allows common local frontend URLs through the backend `.env` file.

Backend `.env` contains:

```env
FRONTEND_URLS=http://localhost:5173,http://127.0.0.1:5173,http://localhost:3000,http://127.0.0.1:3000,http://localhost:4200,http://127.0.0.1:4200
```

This means the frontend can run on:

```text
http://localhost:5173
http://localhost:3000
http://localhost:4200
```

If the frontend uses a different port, add that URL to `FRONTEND_URLS` in `backend/.env`, then restart the backend.

## Planned Pages

```text
/login
/register
/dashboard
/add-service
/services/:id
```

## Main Frontend Tasks

```text
Create login page
Create register page
Store JWT token after login
Protect dashboard routes
Create dashboard summary cards
Create service list/cards
Create add service form
Create service detail page
Display health check history
Display response time chart
Show loading and error states
```

## API Documentation

See:

```text
docs/api.md
```

That file explains every backend route, request body and response shape.