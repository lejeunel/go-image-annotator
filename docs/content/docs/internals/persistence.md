---
title: Data Persistence
weight: 2
---

## Meta-data

To store relevant meta-data, our use-cases define *Repository* interfaces.

We provide adapters that leverage [SQLite](https://sqlite.org/),
which is a very capable option for small/medium scale projects.

## Raw-data

To store the raw image data, we resort to a local file system,
where a local directory contains files that are named after
image IDs.

