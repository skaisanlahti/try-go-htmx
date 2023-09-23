# Go + HTMX + Typescript + SASS

HTMX is having a bit of a renaissance right now, so I wanted to give it a try and play around with Go and PostgreSQL at the same time. This repository covers a few ideas:

- Using HTMX to create a frontend
- Using Go standard library to:
  - Run a web server
  - Use HTML templates
  - Write PostgreSQL queries
- Writing PostgreSQL migrations manually
- Building Typescript and SASS assets with Parcel
- Encoding passwords with Bcrypt and Argon2id
- Managing user sessions in a way that is compatible with HTMX

### Technology

- [Go](https://go.dev/#)
- [Task](https://taskfile.dev/)
- [Docker](https://www.docker.com/)
- [PostgreSQL](https://www.postgresql.org/)
- [pnpm](https://pnpm.io/installation)
- [Typescript](https://www.typescriptlang.org/)
- [HTMX](https://htmx.org/)
- [SASS](https://sass-lang.com/)
- [PicoCSS](https://picocss.com/)
- [Parcel](https://parceljs.org/)

## Setup

- Install [Go](https://go.dev/#) to run and build application
- Install [pnpm](https://pnpm.io/installation) to build and bundle web assets
- Install [Task](https://taskfile.dev/) to run project tasks
- Install [Docker](https://www.docker.com/) to run database container
- Run `task dev` in project root to build database container and apply migrations, build web assets, and run application

## What is HTMX?

> HTMX gives you access to AJAX, CSS Transitions, WebSockets and Server Sent Events directly in HTML, using attributes, so you can build modern user interfaces with the simplicity and power of hypertext. HTMX is small (~14k min.gz’d), dependency-free, extendable, IE11 compatible & has reduced code base sizes by 67% when compared with React.

[Official HTMX documentation](https://htmx.org/)

## Why HTMX?

- Simplifies frontend architecture
- Single source of state, no state duplication on the frontend
- Much smaller Javascript bundle resulting in lightweight sites (good for mobile)
- Rendering HTML partials without page refreshes (SPA-like experience)
- Server-side rendering for better Search engine optimization
- Can be used with any server side framework that supports rendering HTML templates as responses

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

## Using HTMX in HTML

Example in [todos/templates/todo_page.HTML](https://github.com/skaisanlahti/try-go-htmx/blob/dev/todos/templates/todo_page.HTML)

- Split a page into partials using "block" calls which combine "define" and "template" to define a render are in-place
- Embedded HTMX properties in HTML elements let us make backend calls without Javascript
- Transitions and swapping strategy can also be defined with HTML properties

## Using HX headers to trigger requests

Example in [todos/handlers/add_todo.go](https://github.com/skaisanlahti/try-go-htmx/blob/6de383d17423c15507fcac301403606dcbc441a7/todos/handlers/add_todo.go#L47)

- Trigger additional HTMX requests using headers `response.Header().Add("HX-Trigger", "GetTodoList")`

## Using HX-Boost and HX-Location to skip page loads

Example in [users/sessions/require_session.go](https://github.com/skaisanlahti/try-go-htmx/blob/48397ea0fa1c3f49ab3df1e1eaa4e0624665b148/users/sessions/require_session.go#L30) and [users/templates/login_page.HTML](https://github.com/skaisanlahti/try-go-htmx/blob/48397ea0fa1c3f49ab3df1e1eaa4e0624665b148/users/templates/login_page.HTML#L15)

- Add HX-Boost to HTML to use HTMX to load new pages directly to the body tag when clicking links on the page
- Add HX-Location header to a server response to use HTMX to load the target page directly to the body tag instead of a full page load
