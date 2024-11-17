# Datenschutz Training

Using the all new SHEET Stack.

- **_S_**essions (Gorilla Sessions)
- **_H_**otwire / Turbo
- **_E_**cho Webframework
- **_E_**nt ORM
- **_T_**empl and TailwindCSS

## Features

- [x] [Echo Webserver](https://echo.labstack.com)
- [x] [Templ](https://templ.guide)
- [x] [Gorilla Sessions](https://gorilla.github.io)
- [x] [Ent ORM](https://entgo.io/docs/getting-started)
- [x] [HotWire](https://hotwired.dev)
- [x] [TailwindCSS](https://tailwindcss.com)

Lesen:
https://www.akmittal.dev/posts/hotwire-go/
https://adrianhesketh.com/2021/06/04/hotwired-go-with-templ/

## Database

### Create Schema

```bash
go run -mod=mod entgo.io/ent/cmd/ent new User
```

### Create Edge (Relation)

```bash
go run -mod=mod entgo.io/ent/cmd/ent new Car Group
```

### Generate Schema

```bash
go generate ./ent
```

### Inspect The Ent Schema

```bash
atlas schema inspect \
  -u "ent://ent/schema" \
  --dev-url "sqlite://file?mode=memory&_fk=1" \
  -w
```

### Generating migrations

```bash
atlas migrate diff migration_name \
  --dir "file://ent/migrate/migrations" \
  --to "ent://ent/schema" \
  --dev-url "sqlite://file?mode=memory&_fk=1"
```

### Applying migrations

```bash
atlas migrate apply \
  --dir "file://ent/migrate/migrations" \
  --url "sqlite://file.db?_fk=1"
```
