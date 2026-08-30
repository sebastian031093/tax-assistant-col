# Tax Assistant CO

Tax Assistant CO es un asistente financiero y tributario personal para Colombia. El proyecto busca ayudar al usuario a organizar su información financiera, identificar datos faltantes o inconsistentes, clasificar conceptos tributarios y obtener estimaciones explicables sobre su declaración de renta.

> [!IMPORTANT]
> La aplicación es una herramienta de apoyo. No representa una determinación oficial de la DIAN ni sustituye la revisión de un contador o asesor tributario.

El repositorio también funciona como laboratorio de ingeniería de software: cada incremento debe fortalecer el producto, las prácticas de desarrollo y el conocimiento del dominio tributario.

## Estado actual

El proyecto se encuentra en su etapa de fundamentos. Actualmente incluye:

- una aplicación web con React, TypeScript, Vite y Chakra UI;
- una API HTTP en Go;
- conexión de la API a PostgreSQL mediante `pgx`;
- una base de datos PostgreSQL local administrada con Docker Compose;
- un endpoint de salud en `GET /health`;
- apagado controlado de la API ante señales del sistema.

## Principios del producto

- Los cálculos tributarios deben ser determinísticos, versionados y explicables.
- Todo resultado importante debe conservar trazabilidad hacia sus datos de origen, la regla aplicada y su versión.
- La inteligencia artificial puede proponer o explicar; el usuario debe confirmar.
- Durante la etapa de fundamentos solo se deben utilizar datos sintéticos.
- La arquitectura comienza como un monolito modular y evolucionará de acuerdo con necesidades observables.

## Estructura del repositorio

```text
tax-assistant-col/
├── apps/
│   ├── api/                 # API en Go
│   │   ├── cmd/api/         # Punto de entrada
│   │   └── internal/        # Configuración, base de datos y servidor HTTP
│   └── web/                 # Frontend React + TypeScript + Vite
├── deployments/
│   └── docker/              # Docker Compose para PostgreSQL
└── docs/                    # Documentación técnica y de aprendizaje
```

## Requisitos

Antes de comenzar, instala:

- [Go 1.25 o superior](https://go.dev/doc/install)
- [Node.js](https://nodejs.org/) y [Yarn 1.x](https://classic.yarnpkg.com/lang/en/docs/install/)
- [Docker Desktop](https://www.docker.com/products/docker-desktop/) con Docker Compose
- Git

Los comandos siguientes se ejecutan desde la raíz del repositorio, salvo que se indique lo contrario.

## Inicio rápido

La aplicación se levanta en tres procesos: PostgreSQL, API y frontend. Inícialos en ese orden y utiliza una terminal independiente para cada uno.

### 1. Iniciar PostgreSQL

Crea el archivo `deployments/docker/.env` con credenciales exclusivamente locales:

```dotenv
POSTGRES_USER=tax_assistant_app
POSTGRES_PASSWORD=local_development_password
POSTGRES_DB=tax_assistant_co
```

Este archivo está ignorado por Git y no debe contener credenciales de producción.

Levanta la base de datos:

```bash
docker compose --env-file deployments/docker/.env -f deployments/docker/compose.yaml up -d
```

Verifica que el contenedor esté en ejecución:

```bash
docker compose --env-file deployments/docker/.env -f deployments/docker/compose.yaml ps
```

PostgreSQL quedará disponible en `localhost:5432`. Para detenerlo sin eliminar sus datos:

```bash
docker compose --env-file deployments/docker/.env -f deployments/docker/compose.yaml down
```

### 2. Iniciar la API

La API necesita la variable `DATABASE_URL`. Usa los mismos valores configurados para PostgreSQL.

En PowerShell:

```powershell
$env:DATABASE_URL = "postgres://tax_assistant_app:local_development_password@localhost:5432/tax_assistant_co?sslmode=disable"
Set-Location apps/api
go mod download
go run ./cmd/api
```

En Bash, Zsh o Git Bash:

```bash
export DATABASE_URL="postgres://tax_assistant_app:local_development_password@localhost:5432/tax_assistant_co?sslmode=disable"
cd apps/api
go mod download
go run ./cmd/api
```

La API se inicia por defecto en `http://localhost:8080`. Comprueba su estado desde otra terminal:

```bash
curl http://localhost:8080/health
```

Respuesta esperada:

```json
{
  "status": "ok",
  "service": "tax-assistant-api"
}
```

Puedes personalizar el servidor con estas variables de entorno:

| Variable | Valor predeterminado | Descripción |
| --- | ---: | --- |
| `PORT` | `8080` | Puerto HTTP de la API |
| `READ_TIMEOUT_SECS` | `5` | Tiempo máximo de lectura, en segundos |
| `WRITE_TIMEOUT_SECS` | `10` | Tiempo máximo de escritura, en segundos |
| `SHUTDOWN_TIMEOUT_SECS` | `15` | Tiempo máximo para el apagado controlado |
| `DATABASE_URL` | — | URL de conexión a PostgreSQL; es obligatoria |

Detén la API con `Ctrl+C`.

### 3. Iniciar el frontend

En una terminal nueva:

```bash
cd apps/web
yarn install
yarn dev
```

Vite mostrará en la terminal la URL local del frontend, normalmente `http://localhost:5173`.

## Verificación local

### API

```bash
cd apps/api
go test ./...
```

### Frontend

```bash
cd apps/web
yarn lint
yarn build
```

## Roadmap resumido

1. Fundamentos: monorepo, frontend, API, PostgreSQL, Docker y CI.
2. Perfil y año gravable: cuentas, activos, pasivos y primer flujo CRUD.
3. Importación desde Excel con validación, vista previa e idempotencia.
4. Motor tributario con reglas versionadas y estimaciones explicables.
5. Importación y reconciliación de información exógena de la DIAN.
6. Asistencia con IA para extracción, clasificación y detección de inconsistencias, siempre con confirmación humana.

Fuera del alcance se encuentran la presentación automática ante la DIAN, el almacenamiento de credenciales MUISCA, la firma electrónica y el reemplazo de la asesoría profesional.

## Referencias del proyecto

Este README sintetiza la configuración actual del repositorio y la documentación interna de Notion:

- [Tax Assistant CO — Product & Engineering](https://app.notion.com/p/3bfe37f34bf381e0a836d80f7086cc1d)
- [Product & Tax Domain Brief](https://app.notion.com/p/3bfe37f34bf381e1a40fd5a420ceacb4)
- [Engineering Handbook](https://app.notion.com/p/3bfe37f34bf381d4b93af36d40b13b66)
- [Sprint 0 — Fundamentos](https://app.notion.com/p/3bfe37f34bf381b79bc0e9b31913cd68)
