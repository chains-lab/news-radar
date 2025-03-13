Below is a simple README file that you can include in your repository. It explains how developers can run the necessary 
containers for the re-news application using the provided commands, with the default ports:

---

# re-news Development Environment Setup

This document explains how to start the required containers for developing the re-news application. For local 
development, you can use the following Docker commands to run PostgreSQL, Redis, and MongoDB. These commands launch the 
containers with default port mappings:

- **PostgreSQL**: Internal port 5432 mapped to host port **7000**
- **MongoDB**: Internal port 27017 mapped to host port **7100**
- **Redis**: Internal port 6379 mapped to host port **7200**

> **Note:** These settings are intended for development only. For production deployments, consider using managed services or scaling the containers separately.

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/) must be installed on your system.

## Running the Containers

Open your terminal and execute the following commands:

### PostgreSQL

```bash
docker network create re-news-net
```

```bash
docker compose --project-name sso-oauth up -d
docker compose --project-name sso-oauth down -v
```

```bash
docker network inspect re-news-net
docker network connect re-news-net sso-oauth
```

## Starting the re-news Application

After the containers are running, you can start the re-news application normally (e.g., via your IDE or another Docker 
container that runs the app). The application will connect to the dependencies on the default ports defined above.

Happy coding!

---