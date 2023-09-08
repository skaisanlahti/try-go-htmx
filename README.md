# Go + HTMX + Typescript + SASS

HTMX is having a bit of a renaissance right now, so I wanted to give it a try and play around with Go and PostgreSQL at the same time. This repository covers a few ideas:

- Creating an interactive frontend experience with HTMX
- Using Go standard library to:
  - Run a web server
  - Use html templates
  - Write PostgreSQL queries
- Writing PostgreSQL migrations manually
- Building Typescript and SASS assets with Parcel

### Technology

- [Go](https://go.dev/#)
- [HTMX](https://htmx.org/)
- [PostgreSQL](https://www.postgresql.org/)
- [Typescript](https://www.typescriptlang.org/)
- [SASS](https://sass-lang.com/)
- [PicoCSS](https://picocss.com/)
- [Parcel](https://parceljs.org/)

## Setup

```bash
# start database container
docker-compose up -d

# run migrations
cd migrations
cat "0_create_migrations.up.sql" | docker exec -i try-go-htmx-db psql -U dbuser -d try-go-htmx-db
cat "1_create_todos.up.sql" | docker exec -i try-go-htmx-db psql -U dbuser -d try-go-htmx-db
cd ..

# build frontend assets
pnpm install
pnpm build

# run application
go run .
```

## What is HTMX?

> HTMX gives you access to AJAX, CSS Transitions, WebSockets and Server Sent Events directly in HTML, using attributes, so you can build modern user interfaces with the simplicity and power of hypertext. HTMX is small (~14k min.gz’d), dependency-free, extendable, IE11 compatible & has reduced code base sizes by 67% when compared with React.

[Official HTMX documentation](https://htmx.org/)

## Why HTMX?

- Simplifies frontend architecture
- Single source of state, no state duplication on the frontend
- Much smaller Javascript bundle resulting in lightweight sites (good for mobile)
- Rendering HTML partials without page refreshes (SPA-like experience)
- Server-side rendering for better Search engine optimization
- Can be used with any server side framework that supports rendering html templates as responses

## Why Go?

- Simple syntax, easy to learn
- Great performance
- Easy and flexible concurrency
- Static type safety out of the box along with a good formatter
- Good standard library
- Fast compilation times

## Possible use cases

- Rendering large data sets
- Public facing interfaces requiring SEO
- Simple applications with minimal user interaction
- Applications aimed for mobile and low-power devices
- Giving backend focused teams easy to use frontend capabilities

# Points of interest

## Using HTMX in html

Example in [todos/templates/todo_page.html](https://github.com/skaisanlahti/try-go-htmx/blob/dev/todos/templates/todo_page.html)

- Split a page into partials using "block" calls which combine "define" and "template" to define a render are in-place
- Embedded HTMX properties in html elements let us make backend calls without JS
- Transitions and swapping strategy can also be defined with html properties

## Using HX headers to trigger requests

Example in [todos/htmx/add_todo.go](https://github.com/skaisanlahti/try-go-htmx/blob/52e40d35a723b6ddf5018fea5312ce82f0d3f785/todos/htmx/add_todo.go#L47)

- Trigger additional HTMX requests using headers `response.Header().Add("HX-Trigger", "GetTodoList")`

## Feature focused code organization

Vertical slices and feature folders are cool. It's a bit difficult to find the balance between code reuse and feature isolation. In this project I decided to keep each http handler as it's own feature in the htmx folder, but separate repository on it's own, because at least 4 handlers reuse the same database operations. Domain types and functions are also split into their own folder because at least the type is shared by all handlers.
