# Monitoring Dashboard Frontend

React, TypeScript, Vite and Tailwind frontend for the Cloud DevOps Monitoring Dashboard.

## Local setup

1. Copy `.env.example` to `.env`.
2. Ensure the Go backend is available at the configured `VITE_API_BASE_URL`.
3. Install dependencies with `npm install`.
4. Start the development server with `npm run dev`.

## Scripts

- `npm run dev` — start the local development server.
- `npm run build` — type-check and create a production build.
- `npm run lint` — run Oxc lint checks.
- `npm run preview` — preview the production build locally.

## Main routes

- `/login`
- `/register`
- `/dashboard`
- `/services/new`
- `/services/:id`
- `/services/:id/edit`

## Production hosting

This is a single-page application. The production web server must rewrite unknown routes to `index.html` so direct visits and refreshes on routes such as `/services/:id` continue to work.
