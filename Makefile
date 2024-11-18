windows:
	air -c .air.win.toml

mac:
	air -c .air.mac.toml

build:
	bunx tailwindcss -i ./base.css -o ./static/style.css
	templ generate
	go generate ./ent
	go build