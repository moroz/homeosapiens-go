# Homeo sapiens (Go edition)

## Setting up

Prerequisites:

* [Mise](https://mise.jdx.dev/),
* [PostgreSQL 18](https://www.postgresql.org/).

Ensure experimental features are enabled in Mise:

```shell
mise settings experimental=true
```

Install dependencies using Mise:

```shell
mise trust
mise install
cd assets && pnpm install
```

(Optional) Set up Google OAuth secrets. This step is required to enable Sign in with Google. You need to obtain Google OAuth credentials for this step.

```shell
# Copy the example mise settings file and edit the values
cp mise.local.toml.example mise.local.toml
```

Create, migrate, and seed the database:

```shell
mise r db:setup
```
