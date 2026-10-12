---
title: Configuration & Usage
linkTitle: Usage
menu: { main: { weight: 2 } }
weight: 1
---

## Configuration

### Storage

We distinguish two kinds of storage spaces:
1. The database storage space
2. The raw image storage space

We now explain how to configure them:

#### Database Storage

The database stores every data aside from the raw image data, e.g.
annotations, meta-data, logged events, ...
We implement a [SQLite](https://www.sqlite.org/) backend,
which requires a standard file-system path.

You must therefore provide a path where the database file will be
stored, e.g.

``` sh
GOIA_DB_PATH=/home/user/.cache/go-image-annotator/db.sqlite
```

Next, define a path where images will be stored, e.g.

``` sh
GOIA_ARTEFACT_DIR=/home/user/.cache/go-image-annotator/artefacts
```

## Run

To instantiate the server and listen (e.g. on port `8001`), run
``` sh
./go-image-annotator serve -p 8001
```
