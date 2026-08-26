# Cloud Deployment

StatusWatch is deployed on an Oracle Cloud Infrastructure (OCI) Ubuntu VM and is publicly available at:

`https://statuswatch.duckdns.org`

## Production architecture

```text
Internet
   |
   v
DuckDNS hostname
statuswatch.duckdns.org
   |
   v
OCI public IP / VCN security rules
   |
   v
Ubuntu host firewall (iptables)
   |
   v
Caddy reverse proxy
   |-- /      -> frontend container
   `-- /api/* -> backend container

Docker Compose network
   |-- frontend
   |-- backend
   |     `-- PostgreSQL
   |-- Prometheus
   `-- Grafana
```

Caddy is installed directly on the VM and runs as a `systemd` service. The application stack runs with Docker Compose.

## Cloud VM

- Provider: Oracle Cloud Infrastructure
- Region: UK South (London)
- OS: Ubuntu 24.04 LTS
- Architecture: ARM64 / Ampere A1
- VM shape: `VM.Standard.A1.Flex`
- Resources used for this deployment: 2 OCPUs and 12 GB RAM
- Application directory: `~/apps/Cloud-Service-Monitoring-Dashboard`

## Docker services

The Compose stack contains:

- `postgres` - PostgreSQL application database
- `backend` - Go/Gin API and background health checker
- `frontend` - React/Vite frontend served by nginx
- `prometheus` - platform metrics collection
- `grafana` - operator-facing metrics dashboards

Production networking is deliberately restricted:

- PostgreSQL is not published to the VM's public network.
- Prometheus is not published to the VM's public network.
- Grafana binds only to `127.0.0.1:3001` and is accessed through an SSH tunnel.
- Frontend and backend container ports are not exposed through OCI security rules; public application traffic enters through Caddy on ports 80/443.

## Reverse proxy and HTTPS

Caddy is configured in `/etc/caddy/Caddyfile` on the VM.

```caddy
statuswatch.duckdns.org {
    handle /api/* {
        reverse_proxy localhost:8080
    }

    handle {
        reverse_proxy localhost:3000
    }
}
```

Caddy provides a single public entry point and automatically manages TLS certificates for `statuswatch.duckdns.org`.

Traffic routing:

- `https://statuswatch.duckdns.org/` -> React frontend
- `https://statuswatch.duckdns.org/api/*` -> Go backend

HTTP on port 80 remains enabled so Caddy can redirect clients to HTTPS and perform certificate-related HTTP validation when required.

## Network security

The deployment has two inbound filtering layers:

1. OCI VCN security rules.
2. The Ubuntu VM's `iptables` firewall.

Public inbound ports used by the deployment:

- `22/tcp` - SSH administration
- `80/tcp` - HTTP / redirect to HTTPS
- `443/tcp` - HTTPS application traffic

Application and infrastructure ports such as `3000`, `8080`, `3001`, `9090`, and `5432` are not publicly reachable.

During deployment, port 80 initially remained unreachable even after it was opened in the OCI security list. The VM's `iptables` chain contained a general REJECT rule before the HTTP ACCEPT rule, so traffic was rejected before the allow rule could be evaluated. Moving the port 80 rule above the REJECT rule resolved the issue. The final firewall rules are persisted with `netfilter-persistent`.

## Grafana administration

Grafana is intentionally not exposed to the public internet. It can be accessed from an administrator's machine through an SSH tunnel:

```bash
ssh -i <private-key> -L 3001:localhost:3001 ubuntu@<vm-public-ip>
```

With the tunnel active, Grafana is available locally at:

`http://localhost:3001`

Prometheus remains private inside the Docker network and Grafana queries it internally.

## Production environment variables

Production secrets and environment-specific values are stored in the VM's root project `.env` file and are not committed to Git.

Important deployment values include:

```env
FRONTEND_URLS=https://statuswatch.duckdns.org
VITE_API_BASE_URL=https://statuswatch.duckdns.org
APP_BASE_URL=https://statuswatch.duckdns.org
```

The file also contains the PostgreSQL password, JWT secret, Grafana administrator password, and Brevo SMTP credentials. These values must never be committed to the repository.

The `.env` file is restricted on the VM with:

```bash
chmod 600 .env
```

`VITE_API_BASE_URL` is a Vite build-time variable, so changes to it require rebuilding the frontend image.

## Deployment / update workflow

From the VM:

```bash
cd ~/apps/Cloud-Service-Monitoring-Dashboard
git checkout main
git pull origin main
docker compose up -d --build
```

Check container state:

```bash
docker compose ps
```

Useful internal health checks:

```bash
curl http://localhost:8080/health
curl http://localhost:3000/healthz
docker exec monitoring-backend wget -qO- http://prometheus:9090/-/healthy
curl -I http://localhost:3001
```

Expected behaviour:

- Backend health endpoint responds successfully.
- Frontend health endpoint responds successfully.
- Prometheus reports healthy.
- Grafana responds locally and redirects to `/login` when unauthenticated.

## Production validation

The live deployment has been smoke-tested for:

- HTTPS access through Caddy
- frontend-to-backend API requests through `/api`
- registration and login
- email verification using the production HTTPS hostname
- service monitoring and health checks
- downtime alerts and email delivery
- reports
- private PostgreSQL and Prometheus networking
- Grafana access through SSH tunnel only
- blocked public access to ports 3000, 8080, 3001, 9090, and 5432
- healthy Docker containers after testing

## DNS and email

The current public hostname is provided by DuckDNS:

`statuswatch.duckdns.org`

Brevo SMTP is used for verification and downtime alert emails. The current sender remains an independently verified Brevo sender because a DuckDNS subdomain does not provide the same DNS control as a privately owned domain for full sender-domain authentication (DKIM/DMARC). A future custom domain can be authenticated with Brevo for branded sender addresses.
