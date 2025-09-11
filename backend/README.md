## Running the Project with Docker

This project is containerized using Docker and Docker Compose for easy setup and deployment. Below are the specific instructions and requirements for running this project:

### Project-Specific Docker Requirements
- **Go Version:** The build uses Go `1.25.0` (as specified in the Dockerfile).
- **Dependencies:** All Go dependencies are managed via `go.mod` and `go.sum`.
- **Database:** The project uses PostgreSQL, with initialization SQL in `bac.sql`.

### Environment Variables
- The PostgreSQL service uses the following environment variables (set in `docker-compose.yml`):
  - `POSTGRES_USER=postgres`
  - `POSTGRES_PASSWORD=postgres`
  - `POSTGRES_DB=bac`
- No additional environment variables are required for the Go app by default. If you need to add any, uncomment and use the `env_file` section in the compose file.

### Build and Run Instructions
1. **Build and Start Services:**
   ```sh
   docker compose up --build
   ```
   This will build the Go application and start both the app and PostgreSQL services.

2. **Database Initialization:**
   - The `bac.sql` file is included in the app container. To pre-load the database, consider mounting `bac.sql` to `/docker-entrypoint-initdb.d/` in the PostgreSQL service or use a custom init script if needed.

### Ports Exposed
- **Go App:**
  - Exposes port `8080` (mapped to host `8080`)
- **PostgreSQL:**
  - Default PostgreSQL port (`5432`) is used internally; not exposed to host by default.

### Special Configuration
- The Go app runs as a non-root user for security.
- Persistent storage for PostgreSQL is provided via the `pgdata` volume.
- Both services are connected via the `appnet` Docker network.

---
*Ensure Docker and Docker Compose are installed on your system before proceeding. For any custom environment variables or database initialization, adjust the `docker-compose.yml` as needed for your workflow.*
