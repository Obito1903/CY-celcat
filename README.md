# CYtech Celcat to ICS

## build

### dépendance

- golang
  - github.com/arran4/golang-ical
  - golang.org/x/net

``` sh
go get github.com/arran4/golang-ical
go golang.org/x/net
```

### build le projet

``` sh
go build ./celcat.go
```

## Utilisation

Premierement il faut mettre vos identifiant celcat dans le fichier `config.json`

une fois que c'est fait il vous suffit d'executer le programme pour obtenir un fichier `data.ics` qui pourra etre importé dans n'importe quel gestionnaire d'agenda.

### windows

```sh
.\celcat.exe
```

### linux

```sh
./celcat
```

### Parametres

#### Fichier de config different

Il est possible de spécifier un fichier de configuration (doit etre dans le meme dossier, flemme de me faire chier a parse les espaces dans les argument)

```sh
./celcat -c example.config.json
```

#### Récuperer l'agenda sur une période différente

Vous pouvez aussi spécifier la période sur laquel récuperer l'agenda

```sh
./celcat -d 2021-03-30 2021-04-30
```