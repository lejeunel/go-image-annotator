---
title: Frontend
weight: 1
---

We try to keep this project lightweight, by avoiding
heavy frontend framework. Instead, we make the choice to
generate most content on the server-side, and let
the client be a *thin viewer*.

In particular,
we rely on [HTMX](https://htmx.org/) to generate views on the server and update
the client, and [Alpine.js](https://alpinejs.dev/) for client-side interactivity.

For the annotation client-side widget, we use [Annotorious.js](https://annotorious.dev/),
which allows to draw and update different kinds of shapes.

Last, we style our pages using [tailwindcss](https://tailwindcss.com/).
