---
title: REST/HTTP API
weight: 3
---

While our web frontend might fullfill basic annotation and exploration needs,
a REST/HTTP API allows to interact with the application through
scripts.

We therefore expose a number of HTTP endpoints to handle
requests through json payloads.

Implementation-wise, we leverage a *specs first* approach,
where a `yaml` file defines data models and endpoint URLs,
and we use [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen)
to generate Go code from it.
