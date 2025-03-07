# Blog - ffss.dev

This is the repository containing my personal blog, you can check it out at https://ffss.dev.

This blog is written in Go and uses [sqlite](https://www.sqlite.org/) for enabling
full-text search of article content and the excellent [goldmark](https://github.com/yuin/goldmark) library for parsing Markdown.

## Requirements

- [Go 1.24](https://go.dev)
- [Node 22](https://nodejs.org/en)
- [Make](https://www.gnu.org/software/make/)
- [Goose](https://github.com/pressly/goose) - Database migration tool
- [Reflex](https://github.com/cespare/reflex) - Reruns commands when files change

## Configuration

Server configuration is done using command line flags. All Available flags are documented below:

- `addr` - Sets the listen address of the server (default `:4000`)
- `dev` - Sets the application in development mode (default: `true`)
- `articles` - Sets the articles dir path (default: `articles`)
- `static` - Sets the static assets dir path (default: `web/static`)
- `views`- Sets the HTML templates dir path (default: `web/views`)
- `db-path` - Sets the SQLite database path (default: `blog.db`)

## Development

To start development, first install the Node dependencies using the command below:

```bash
npm ci
```

The application uses [tailwind](https://tailwindcss.com/) as a CSS framework.
The input file is located at `web/input.css`.

To start compiling tailwind classes, run the command below:

```bash
npm run tw:watch
```

Make sure all migrations are applied by running:

```bash
goose -dir migrations sqlite blog.db up
```

To start the server in develpment mode, run the Make script below:

```bash
make watch
```

## Production

First, you'll need to compile the server with the command below:

```bash
make build
```

After that, compile all tailwind classes with the npm script below:

```bash
npm run tw:build
```

Make sure all migrations are applied by running:

```bash
goose -dir migrations sqlite blog.db up
```

Then, start the server by running the command below. Make sure to specify
that dev mode is off, since it is on by default:

```bash
bin/server --dev=false
```
