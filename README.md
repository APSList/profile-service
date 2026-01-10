# Profile Service

This microservice is used for managing user profiles, 
organization-based authorization, 
and user preferences. 
Built with **Go** and integrated 
with **Supabase Auth**.

---

## Prerequisites

* **Go 1.22+**
* **Supabase project**  instance

---

## Local Development

### 1. Setup
reate a .env file in the service root:
```
DATABASE_URL=Supabase url
APP_HOST=localhost
APP_PORT=8080
```

```bash
go mod tidy
main.go
docker build -t hostflow-profile-service .
docker run -p 50052:50052 --env-file .env hostflow-profile-service