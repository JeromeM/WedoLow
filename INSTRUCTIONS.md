# Test technique - Golang

## Problème

Créer un service stateless sous forme d'un module go muni d'un serveur REST acceptant les appels suivants :

- `POST /user` crée un ou plusieurs utilisateurs aléatoirement en fonction des paramètres de requête suivants :
  - `count` : Nombre d'utilisateurs à ajouter (défaut : 1)
  - `gender` : Genre du ou des utilisateurs à ajouter (`any`, `male`, `female`) (défaut : `any`)
- `GET /user/${id}` récupère un user par ID, l'ID étant un UUID v4
- `GET /users` récupère un plusieurs utilisateurs, avec les paramètres de requête suivants :
  - `limit`: nombre maximum d'utilisateurs à récupérer, 0 pour pas de limite (défaut : 0)
  - `name`: filtre sous forme de regexp sur le prénom ou le nom de manière insensible à la casse (défaut : pas de filtrage)

Quelques exemples :

- `POST /user?count=10` crée 10 utilisateurs de tous genres
- `POST /user?gender=male` crée 1 utilisateur masculin
- `POST /user?count=10&gender=female` crée 10 utilisateurs féminins
- `GET /user/b4912564-4f0b-47d7-822c-bced06317834` récupère l'utilisateur `b4912564-4f0b-47d7-822c-bced06317834`
- `GET /users?limit=10` récupère les 10 premiers utilisateurs
- `GET /users?name=%24a.*` récupère les utilisateurs dont le nom ou le prénom commence par la lettre a ou A (note: `%24` == `$`)
- `GET /users?limit=10&name=%24a.*` récupère les dix premiers utilisateurs dont le nom ou le prénom commence par la lettre a ou A

Les utilisateurs devront être stockés et récupérés via une DB externe de votre choix.

Pour générer les utilisateurs, vous utiliserez [RandomUser API](https://randomuser.me/documentation), soit directement en requêtant l'API publique (le plus simple) soit en déployant leur code dans un container (dans ce cas, les instructions pour le lancer doivent être fournies)

A noter qu'il n'est absolument pas nécessaire de requêter, parser et stocker tous les champs disponibles. Récupérez simplement les champs nécessaires (genre, nom, prénom) et quelques autres champs de votre choix pour enrichir un peu la donnée. Attention à ne pas ajouter de seed dans la requête pour ne pas avoir les mêmes users à chaque fois.

Deux récupérations successives d'utilisateurs identiques devront renvoyer le même résultat

## Livraison

- Vous livrerez le code sur la branche `main` d'un repository que vous pourrez nous transmettre dans une archive compressée. Intégrez cet énoncé à la racine du projet sous le nom `INSTRUCTIONS.md`
- Vous fournirez une documentation autosuffisante pour lancer le service et la DB (pour l'installation de services tiers, une liste de requirements pointant vers la documentation officielle suffit)
- Vous fournirez une liste de scripts ou commandes (e.g, commandes `curl`) permettant de lancer des exemples de requêtes sur le service
- Ecrire un paragraphe pour motiver vos choix d'architectures, de libraries, et sur ce que vous ajouteriez pour améliorer les performances de votre service

## Bonus (facultatifs)

Voici quelques bonus possibles (dans l'ordre, du plus appréciable au plus accessoire) :

- Ecrivez des tests unitaires
- Ecrivez un Dockerfile et un script docker-compose permettant de lancer le service et sa DB en une unique commande et rendant le disque de la DB persistant
- Ajoutez un système de pagination (et documentez les ajouts de paramètres des requêtes)
- Ajoutez un système de tracing (e.g, Jaeger avec Open Telemetry)

## Notes diverses

- Voici ce qui sera jugé en priorité dans ce test :
  - Le respect du cahier des charges
  - L'organisation du code et sa clarté
    - Bonne séparation des handlers HTTP, de la business logic et de la logique de stockage
    - Utilisation d'interfaces aux endroits appropriés
    - Nommage des variables clair
    - Commentaires aux endroits qui peuvent en nécessiter
    - etc.
  - Performances du code
  - Une démonstration de la bonne utilisation de `git` est appréciée (1 branche == 1 feature, pertinence des messages de commit, etc)

En cas de problème d'accès ou pour toute question, n'hésitez pas à envoyer un mail à bbarrois@wedolow.com.
