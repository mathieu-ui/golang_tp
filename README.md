# golang_tp

## Documentation de l'API

Cette API est écrite en Go et utilise le framework Chi pour la gestion des routes HTTP. Elle est conçue pour être exécutée sur le port 3333.

> 'old.main.go' est une version de l'API sans l'utilisation de la base de donnée.

## Routes

### GET /moutonlist

Cette route est destinée à renvoyer une liste de moutons.

```bash
curl --location --request GET 'http://localhost:3333/moutonlist'
```

### Post /mouton

Cette route est destinée à créer un mouton. Exemple avec un mouton nommé "juju" ou on change le poid.

```bash
curl --location 'http://localhost:3333/mouton' \
--header 'Content-Type: application/json' \
--data '{  
    "Name":"juji",
    "Age":10,
    "Weight": 100
}  '
```

### Post /updateMouton

Cette route est destinée à mettre à jour un mouton. Exemple avec un mouton nommé "juju".

```bash
curl --location 'http://localhost:3333/updatemouton' \
--header 'Content-Type: application/json' \
--data '{  
    "Id":9,
    "Name":"juju",   
    "Age":100,
    "Weight": 90
}  '
```

### Post /dellmouton

Cette route est destinée à supprimer un mouton. Exemple avec un mouton nommé "juju".

```bash
curl --location 'http://localhost:3333/dellmouton' \
--header 'Content-Type: application/json' \
--data '{  
    "Id":7,
    "Name":"juju",
    "Age":10,
    "Weight": 100
}  '
```

> Les données sont stockées dans une BBD Postgres. Cette base de donnée est dans un container Docker. Pour lancer le container, il faut utiliser la commande suivante : 'docker compose up -d'
