---
title: Internals
description: >
    Where we explain how this application is built, with what technologies, and how one might proceed to extend it
weight: 8
---

This project is meant as a starter-kit for small/medium sized
annotation projects. Therefore, we attempt to fullfil the following requirements:

- **Modular architecture** to facilitate its extension and modification. Concretely,
    where some projects would favor low-code patterns such as
    [Active records](https://en.wikipedia.org/wiki/Active_record_pattern), we
    rather favor
    more granular composition and layering with
    the [Clean architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html),
    where we separate our code into use-cases and *plugins*.
- **Few dependencies** on third-party packages. We strive to use Go's standard library,
    and lean to simple alternatives when necessary.

