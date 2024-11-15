# Datenschutz Training

inspired by [https://github.com/nwrenger/htmx-gorm-gin/blob/main/main.go](https://github.com/nwrenger/htmx-gorm-gin/blob/main/main.go)

## Features

- **RESTful API:** Built on the principles of Representational State Transfer (REST) for efficient and standardized communication.
- **SQLite Database:** Utilizes a SQLite database for data storage, providing a simple and self-contained solution.

## Usage

- **Install Perquisites:** You have to have [go](https://go.dev/), [air](https://github.com/cosmtrek/air) and [bun](https://bun.sh/) installed.
- **Install Dependencies:** Install dependencies of go (in [root](/)) and bun (in [content](static/content/)).
- **Run Dev:** Finally, You have to use the `air` command, it's pre-configured in the [air-toml](.air.toml).
- **Build:** To build the project you have to run the following command, **make sure to include in your export the static files**:

```sh
cd ./static/content && bunx tailwindcss -i ./../../web/input.css -o ./dist/output.css && cd ./../../ && templ generate && go build .
```
