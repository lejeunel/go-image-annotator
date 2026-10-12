---
title: Home
linkTitle: Home
menu: { main: { weight: 1 } }
---

A lightweight web application to manage, store, and annotate images
built in [Go](https://go.dev/).

This application is well suited for small/medium sized teams
working on computer vision projects, where annotating images 
allow to train and evaluate Machine Learning models.

Its primary aim is to be a lightweight and self-contained alternative
to the more complex solution. Concretely, this app holds in a single binary
and works without any moving parts by leveraging [SQLite](www.sqlite.org)
as its backend database, thereby only needing a local file-system
for persistence.

## Features

- Annotate images with image-wise labels, bounding-boxes or polygons.
- Clone image and re-define annotations.
- Assign custom meta-data and apply filtering/querying using a custom query language.
- Automate tasks with API calls.
- Create annotation profiles to restrict pre-defined subsets of labels.
- Leverage group and role-based access control for safer collaboration.
- Authenticate users with simple password or use an OAuth provider.
