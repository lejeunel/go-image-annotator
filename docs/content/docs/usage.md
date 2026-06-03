---
title: Configuration & Usage
linkTitle: Usage
menu: { main: { weight: 2 } }
weight: 1
---

## Configuration

This application stores meta and raw-data on a local mount.
In particular, define a path for the SQLite database, e.g.

``` sh
GOIA_DBPATH=/home/user/.cache/go-image-annotator/db.sqlite
```

Next, define a path where images will be stored, e.g.

``` sh
GOIA_ARTEFACTDIR=/home/user/.cache/go-image-annotator/artefacts
```

## Build and Run

Build the application with:

``` sh
go build .
```

To run the service, run
``` sh
./go-image-annotator serve -p 8001
```
