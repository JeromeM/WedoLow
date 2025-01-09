# User Service

Ce service permet de créer et gérer des utilisateurs aléatoires en utilisant l'API RandomUser.

## Prérequis

- Go 1.23+
- Docker et Docker Compose
- PostgreSQL (si exécution locale)
- Jaeger (si exécution locale)

## Installation et Lancement

### Avec Docker Compose (recommandé)

1. Clonez le repository
2. Lancez le service avec :
```bash
docker compose up --build
```

Le service sera accessible sur `http://localhost:8080`

### Installation Locale

1. Installez PostgreSQL
2. Créez une base de données :
```sql
CREATE DATABASE userdb;
```

3. Installez Jaeger

3. Installez les dépendances :
```bash
go mod download
```

4. Créez un fichier `.env` :
```
PORT=8080
DATABASE_URL=postgres://postgres:postgres@localhost:5432/userdb?sslmode=disable
RANDOM_USER_API=https://randomuser.me/api/
JAEGER_ENDPOINT=http://jaeger:14268/api/traces
```

5. Lancez l'application :
```bash
go run main.go
```

## Exemples d'Utilisation

> **Note**: 
> 
> Pour faciliter la lecture des réponses JSON, il est recommandé d'installer `jq`, un utilitaire en ligne de commande pour traiter les données JSON. Vous pouvez l'installer en utilisant votre gestionnaire de paquets préféré, par exemple :
> ```bash
> sudo apt install jq       # pour Debian/Ubuntu
> brew install jq           # pour macOS
> ```
>
> Vous pouvez aussi tout simplement utiliser un outil comme Postman ou Thunder Client directement dans VSCode.

### Création d'utilisateurs

Créer 10 utilisateurs de tous genres :
```bash
curl -X POST "http://localhost:8080/user?count=10"
```

Créer un utilisateur masculin :
```bash
curl -X POST "http://localhost:8080/user?gender=male"
```

Créer 10 utilisateurs féminins :
```bash
curl -X POST "http://localhost:8080/user?count=10&gender=female"
```

### Récupération d'utilisateurs

Récupérer un utilisateur par ID :
```bash
curl "http://localhost:8080/user/b4912564-4f0b-47d7-822c-bced06317834"
```

Récupérer les 10 premiers utilisateurs :
```bash
curl "http://localhost:8080/users?limit=10"
```

Récupérer les utilisateurs dont le nom commence par 'a' :
```bash
curl "http://localhost:8080/users?name=^a.*"
```

## Architecture et Choix Techniques

### Architecture

Le projet utilise une architecture hexagonale avec les couches suivantes :
- **API** (handlers) : gestion des requêtes HTTP
- **Service** : logique métier
- **Database** : accès aux données
- **Model** : structures de données

### Bibliothèques Utilisées

| Bibliothèque    | Description                                                                   |
| :-------------: | :---------------------------------------------------------------------------- |
| *gin-gonic/gin* | Framework web performant et simple d'utilisation                              |
| *gorm*          | ORM pour la gestion de la base de données. Utilisé ici avec une base Postgres |
| *uuid*          | Génération d'identifiants uniques                                             |
| *godotenv*      | Gestion des variables d'environnement                                         |
| *otel*          | Télémétrie et communication avec Jaeger                                       |

### Performances et Améliorations Possibles

1. **Cache** : Implémenter un cache Redis pour les requêtes fréquentes
2. **Monitoring** : Intégrer Prometheus et Grafana
3. **Rate Limiting** : Ajouter une limitation de requêtes
4. **Indexes** : Optimiser les index PostgreSQL pour les recherches par nom
5. **Connection Pooling** : Configurer un pool de connexions optimal

## Tests

Lancer les tests unitaires :
```bash
go test ./... -v
```

## Bonus Implémentés

- [x] Tests unitaires (*partiels*)
- [x] Docker et Docker Compose
- [x] Pagination
- [x] Tracing (*partiel*)