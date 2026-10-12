---
title: Configuration
linkTitle: Configuration
menu: { main: { weight: 2 } }
weight: 1
---

# Introduction

This page elaborates on your configuration options
to setup this application. In particular, we explain
how to setup:

1. The database storage space where meta-data will reside.
2. The raw image storage space.
3. The SMTP service which will send e-mails to users.

## Database Storage

The database stores every data aside from the raw image data, e.g.
annotations, meta-data, logged events, ...
We implement a [SQLite](https://www.sqlite.org/) backend,
which requires to setup a file-system path, e.g.

``` sh
GOIA_DB_PATH=/home/user/.cache/go-image-annotator/db.sqlite
```


## Raw-image Storage

To store raw-images, we provide 2 options:
1. The local file-system
2. A S3-compatible blob-store.

Note that the file-system-based solution will get you started quickly,
but might not scale well for your use-case.

### Local File-System

To store images on the file-system, simply provide a file URI, e.g.

``` sh
GOIA_IMAGES_URI=file:///home/user/.cache/go-image-annotator/images
```

### S3 blob-store

To setup an S3 bucket hosted on amazon, provide the s3 URI, e.g.

``` sh
GOIA_IMAGES_URI=s3://my-bucket/a-prefix/images
```

The authentication will by default follow the [credential provider chain](https://docs.aws.amazon.com/sdkref/latest/guide/standardized-credentials.html),
i.e. it will:
1. Look for `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` env. variables.
2. If they don't exist, use the login credentials, usually stored in `$HOME/.aws`
3. Else, assume an IAM role.

#### Usage outside of AWS

You may also connect to other S3-compatible object-stores, such
as [minio](https://github.com/minio/minio), or [backblaze](https://www.backblaze.com/).

To do so, you must set both `GOIA_S3_ACCESS_KEY` and `GOIA_S3_SECRET`.

Also, MinIO and most self-hosted services
will need to have `GOIA_S3_USE_PATH_STYLE=true`.

## SMTP server

