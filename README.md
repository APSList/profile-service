# Profile Service
Mikrostoritev za upravljanje uporabniških profilov, nastavitev organizacij in identitet v sistemu Hostflow.
## Odgovornosti
profile-service je odgovoren za hrambo in upravljanje podatkov o uporabnikih in njihovih organizacijah.

Glavne odgovornosti zajemajo:

Upravljanje profilov: CRUD operacije nad osebnimi podatki uporabnikov.

Upravljanje organizacij: Hramba podatkov o organizacijah.

Logična izolacija (Multi-tenancy): Zagotavljanje, da uporabniki dostopajo le do podatkov svoje organizacije na podlagi organization_id.

IAM sinhronizacija: Integracija s Supabase Auth za posodabljanje uporabniških metapodatkov.

## Tehnološki sklad
- **Go (Golang)**
- **Gin Web Framework (HTTP API)**
- **PostgreSQL** (Supabase)
- **UberFx**
- **Swagger**
- **HealthChecks** (liveness/readiness)
- **Prometheus** (prometheus-net)


## API
### Swagger
- Swagger UI: `/swagger`
- OpenAPI JSON: `/swagger/v1/swagger.json`

OpenAPI specifikacija se generira z ukazom: swag init -g main.go

## Avtorizacija
Servis zahteva veljaven Supabase JWT žeton v glavi Authorization. V Swaggerju uporabite gumb Authorize in vnesite žeton v formatu: Bearer <token>.

## Model napak
Servis vrača standardne JSON odgovore v obliki:

```
{
"error": "Opis napake",
"message": "Podrobnejše sporočilo"
}
```

## Konfiguracija
Servis uporablja okoljske spremenljivke (.env), ki se nalagajo ob zagonu prek Uber Fx modula.

### Nastavitve (env var)
```
DATABASE_URL=Povezovalni niz za povezavo s PostgreSQL/Supabase bazo
APP_HOST=localhost
APP_PORT=8080
```

## Lokalno testiranje

## CI/CD in pravila razvoja
GitHub Actions workflowi
#### PR validacija (`pr.yaml`)
- **Trigger**: PR → `main`
- **Koraki**: restore → build → test
- **Pravila**: naslov PR mora slediti “conventional” prefiksom:
    - `feat:`, `fix:`, `chore:`, `docs:`, `style:`, `refactor:`, `perf:`, `test:`, `ci:`

#### DEV CI/CD (`dev.yaml`)
- **Trigger**: `push` → `dev`
- **Koraki**:
    1) restore/build/test
    2) build Docker image
    3) push image v registry z tagom **kratkega SHA** (`${GITHUB_SHA::7}`)
    4) checkout deployment repota (`APSList/Hostflow`, veja `dev`)
    5) `helm upgrade --install` za **DEV** okolje (nastavi `image.tag` na kratek SHA)

#### Release PR (`release-please.yaml`)
- **Trigger**: `push` → `main`
- **Namen**: `release-please` pripravi/posodobi **release PR** (changelog + bump verzije) na podlagi conventional sprememb.

#### PROD release (`release.yaml`)
- **Trigger**: `git tag vX.Y.Z` (npr. `v1.2.3`)
- **Koraki**:
    1) restore/build/test
    2) build + push Docker image z tagom **verzije** (`vX.Y.Z`)
    3) checkout deployment repota (`APSList/Hostflow`, privzeta veja)
    4) `helm upgrade --install` za **PROD** okolje (nastavi `image.tag` na `vX.Y.Z`)

---

### Deploy model (booking-service repo → deployment repo)

1. **booking-service repo** zgradi artefakt:
    - Docker image se zgradi iz trenutnega commita.
    - Image se pushne v registry (DockerHub/registry).

2. **Deployment repo** definira, *kako* in *kam* se deploya:
    - Helm chart + `values.yaml` (in pogosto `values-dev.yaml`/`values-prod.yaml`) so v deployment repotu.
    - Deployment repo je “source of truth” za:
        - namespace, ingress, replicas, resources
        - env var/secret reference (DB, Stripe, Kafka, itd.)
        - health probes, autoscaling, service/ports

3. **Helm deploy**:
    - Pipeline naredi `helm upgrade --install` in ob tem nastavi vsaj:
        - `image.repository`
        - `image.tag` (DEV = kratek SHA, PROD = verzija)

---