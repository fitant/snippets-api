## Snippets API
---
Try out Snippets hosted by Fitant here: https://snippets.fitant.cloud

Or, if you prefer, use the API directly: https://snippets.fitant.cloud/api

```
> echo "Hello, World!" > hello-world.txt 

> curl -F "snippet=@hello-world.txt" https://snippets.fitant.cloud/api
{"URL":"https://snippets.fitant.cloud/api/r/MorallyStimulate"}

> curl https://snippets.fitant.cloud/api/r/MorallyStimulate
Hello, World!

> # Or, just use curl --upload-file

> curl --upload-file hello-world.txt https://snippets.fitant.cloud/api 
{"URL":"https://snippets.fitant.cloud/api/r/EmbossChemicals"}

> curl https://snippets.fitant.cloud/api/r/EmbossChemicals
Hello, World!

> echo "Snippets is Awesome!"
Snippets is Awesome!
```
---

### Features:
- Store and fetch encrypted snippets
- Snippets get saved against an easy to remember, id like `HedgingSmitten` 
- Snippets are compressed using [brotli](https://github.com/google/brotli)
- Use MongoDB as storage backend
- Send snippet as formatted JSON, raw body or multi-part form
- Use `POST` and `PUT` at endpoint `/snippets` to save snippet
- Get snippet using `GET` at endpoint `/snippets/<ID>` or `/snippets/r/<ID>`
- Snippets are by default ephemeral and stored in a [capped collection](https://docs.mongodb.com/manual/core/capped-collections/) 

---

### Quick Start:
***Prerequisites: Git, Docker, Docker Compose and curl***
- `git clone --depth=1 https://github.com/fitant/snippets-api`
- `cd snippets-api`
- `docker compose up -d`
- Upload a snippet
    - `curl --upload-file dev.env http://localhost:8080/snippets`
    - Copy the URL field returned in JSON
- Fetch snippet
    - `curl <URL you copied>`
- `docker compose down`

### Setup for Development:
- `git clone --depth=1 https://github.com/fitant/snippets-api`
- `cd snippets-api`
- Start database server
  - `docker compose -f docker-compose.dev.yaml up -d`
- `go mod download`
- Start application server
  - `env $(cat dev.env | xargs -L 1) go run src/main.go`

### Deployment / Self-Hosting:
The easiest way to self-host is to run it in a container (image: `realsidsun/snippets-api`), reverse-proxy it and connect it to a [MongoDB Atlas](https://www.mongodb.com/atlas/database) instance

***NOTE:*** Look at the Config section down below before you deploy

---

### Cryptographic Specification
- ***Refer [b60e2da](https://github.com/fitant/snippets-api/commit/b60e2dadfc89a8307dc2811273415b5e1e158c0d) to understand cryptographic strategy***

---

### Config
Configuration is done through environment variables

#### General Config:

| Name      | Type / Options                     | Description                         | Required | Default |
|-----------|------------------------------------|-------------------------------------|----------|---------|
| ENV       | string                             | Application Environment             | no       | dev     |
| LOG_LEVEL | debug / info / warn / error        | Log Level to print                  | no       | debug   |
| OVERRIDES | comma and colon seperated mappings | override certain IDs for About, etc | no       |         |

***Example Overrides:*** About:BackwashLicorice,PrivacyPolicy:TranceUnsterile

#### Cryptographic Config:

| Name            | Type / Options | Description                         | Required      | Default |
|-----------------|----------------|-------------------------------------|---------------|---------|
| SALT            | string         | SALT used for ARGON2 Key Derivation | yes           |         |
| CIPHER          | AES / SEAT     | Block Cipher used for encryption    | no            | AES     |
| CIPHER_UNTESTED | boolean        | Enable Untested Ciphers             | If using SEAT | false   |
| ARGON2_MEM      | number         | ARGON2 Memory / space param in MB   | no            | 32      |
| ARGON2_ROUNDS   | number         | ARGON2 rounds / iterations param    | no            | 8       |
| ARGON2_ID_ROUNDS| number         | ARGON2 rounds for ID generation     | no            | 8       |

#### AWS S3 Config:

| Name           | Type   | Description                        | Required |
|----------------|--------|------------------------------------|----------|
| AWS_ACCESS_KEY | string | AWS Programmatic Access Key / ID   | yes      |
| AWS_SECRET_KEY | string | Associated Programmatic Secret Key | yes      |
| AWS_REGION     | string | AWS Hosting Region                 | yes      |
| AWS_S3_BUCKET  | string | S3 Bucket Name                     | yes      |

#### HTTP Server Config:

| Name               | Type / Options           | Description                                   | Required | Default               |
|--------------------|--------------------------|-----------------------------------------------|----------|-----------------------|
| HTTP_LISTEN_HOST   | string                   | HTTP Server listen host                       | no       | 127.0.0.1             |
| HTTP_LISTEN_PORT   | number                   | Replica Set name if using replicaset instance | no       | 8080                  |
| HTTP_CORS_LIST     | comma seperated strings  | Allowed HTTP cross origins list               | no       | http://localhost:*    |
| HTTP_BASE_URL      | string                   | HTTP/S frontend URL to use for formatting     | no       | http://localhost:8080 |
| HTTP_API_ENDPOINT  | string                   | API mount Endpoint from base                  | no       | /snippets             |
| HTTP_RETURN_FORMAT | json / raw               | Default URI for URL to created snippet        | no       | raw                   |
