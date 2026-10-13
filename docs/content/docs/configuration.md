---
title: Configuration
linkTitle: Configuration
menu: { main: { weight: 2 } }
weight: 1
---

We explain here how to setup:

1. The database location, i.e. where pretty much all the application state will reside.
2. The raw image storage space.
3. The SMTP service which will send e-mails to users.

# Database Storage

The database stores every data aside from the raw image data, e.g.
annotations, meta-data, logged events, ...
Our technology of choice is [SQLite](https://www.sqlite.org/),
which manages an RDBMS in a single random access file.

Specify a path on the local file-system to the SQLite file, e.g.

``` sh
GOIA_DB_PATH=/home/user/.cache/go-image-annotator/db.sqlite
```


# Raw-image Storage

To store raw-images, we provide two options, each having their pros and cons.
1. Local file-system: 
  - Easy to setup
  - Straight-forward to parse and backup.
  - Difficult to scale for large amounts of data and high resolution images.
  
2. A S3-compatible blob-store: More scalable
  - More involved setup
  - Usually slower than a local drive
  - Scales better

## Local File-System

To store images on the file-system, simply provide a file URI, e.g.

``` sh
GOIA_IMAGES_URI=file:///home/user/.cache/go-image-annotator/images
```

## S3 blob-store

To setup an S3 bucket hosted on amazon, provide the s3 URI, e.g.

``` sh
GOIA_IMAGES_URI=s3://my-bucket/a-prefix/images
```

The authentication will by default follow the [credential provider chain](https://docs.aws.amazon.com/sdkref/latest/guide/standardized-credentials.html),
i.e. it will:
1. Look for `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` env. variables.
2. If they don't exist, use the login credentials, usually stored in `$HOME/.aws`
3. Else, assume an IAM role.

### Usage outside of AWS

You may also connect to other S3-compatible object-stores, such
as [minio](https://github.com/minio/minio), or [backblaze](https://www.backblaze.com/).

To do so, you must set both `GOIA_S3_ACCESS_KEY` and `GOIA_S3_SECRET`.

Also, MinIO and most self-hosted services
will need to have `GOIA_S3_USE_PATH_STYLE=true`.

## SMTP server

The email module is responsible to send to users
temporary credentials when:
- They are logging-in to a newly created account using the email/password flow.
- They have forgotten their password and need to reset it.

In the general case, the app would still work as expected in the case where
your users use the OAuth flow. Concretely, they will use
credentials given by a third-party (Google, Github), and therefore
will not need to reset their password.

In the case where you choose the latter approach, you will still
find a warning message on startup, which you may safely ignore.

To setup the SMTP server, you will need to set the following:

``` sh
GOIA_URL=https://www.example.com # URL where users will be directed to reset their password
GOIA_SMTP_HOST=smtp.example.com
GOIA_SMTP_PORT=587
GOIA_SMTP_USERNAME=admin@example.com
```
