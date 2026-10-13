---
title: Notions
description: >
    Where we develop on the core notions to organize, and version your work. 
weight: 2
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

The following sub-sections provides technical details for use by developers.
