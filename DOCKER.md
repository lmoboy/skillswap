# Running with Docker

You can run the Skillswap frontend using Docker for a consistent and isolated development or preview environment. The provided Docker setup uses Node.js version 22.13.1 (slim) and exposes the SvelteKit preview server on port 5173.

## Docker Requirements

- Docker and Docker Compose installed on your system
- No additional environment variables are required by default (uncomment the `env_file` line in `docker-compose.yml` if you add a `.env` file to `frontend/`)

## Build and Run Instructions

1. **Build and start the frontend service:**

    ```bash
    docker compose up --build
    ```

    This will build the frontend image from the `./frontend` directory and start the SvelteKit preview server.

2. **Access the app:**

    Open your browser and navigate to [http://localhost:5173](http://localhost:5173) to view the frontend.

## Service Details

- **Service name:** `typescript-frontend`
- **Port exposed:** `5173` (mapped to host)
- **Network:** `appnet` (Docker bridge network)
- **Node.js version:** 22.13.1 (as specified in the Dockerfile)

No external databases or persistent volumes are required for the frontend service. If you need to provide environment variables, add a `.env` file to the `frontend/` directory and uncomment the `env_file` line in the `docker-compose.yml`.

---

*The above instructions are specific to running the frontend with Docker. For local development or backend setup, refer to the sections above.*
