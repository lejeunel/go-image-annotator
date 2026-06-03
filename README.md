# Go Image Annotation Tool

A simple web application
to manage, store, and annotate images using bounding boxes
built in [Go](https://go.dev/).

## Goals

- Store and organize images in collections
- Annotate images through a simple web interface

## Features

- Simple web frontend to add bounding boxes
- HTTP/REST API to (among other things):
  - Add new annotation labels
  - Ingest images into collections
  - Add bounding box annotations
    
## Usage

### Configuration

This application stores meta and raw-data on a local mount.
In particular, define a path for the SQLite database, e.g.

``` sh
GOIA_DBPATH=/home/user/.cache/go-image-annotator/db.sqlite
```

Next, define a path where images will be stored, e.g.

``` sh
GOIA_ARTEFACTDIR=/home/user/.cache/go-image-annotator/artefacts
```

`

### Building

To build the binary, run

``` sh
go build .
```

### Ingesting images from a local directory

Aside from the HTTP/REST endpoints, we
provide basic CLI functions to get started quickly.

First, create a new collection that will hold images:

``` sh
./go-image-annotator create-collection my-new-collection
```

### Run web server

You may then launch the web server on port `8001` with:

``` sh
./go-image-annotator serve -p 8001
```

## Development dependencies

To (re)-generate the HTTP endpoints, you will need
`oapi-codegen` version `>=2.5.1`.

Last, CSS must be generated using `tailwindcss` version `4`.

All other dependencies are bundled in this repository
as javascript files.

Optionally, we provide a [Nix Flake](https://nixos.wiki/wiki/Flakes)
that provides a development environment (dev-shell)
with all requirements included.

