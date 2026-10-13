---
title: Access Restriction
description: >
  Where we explain how to restrict access to your resources
weight: 4
linkTitle: Authorization
menu: { main: { weight: 20 } }
---

# Groups

Each collection must be assigned to a **group**, which serves
to restrict access to specific users (or groups of users).

The specific of the authorization logic must be customized by 
implementing an interface, which receives a [context](https://pkg.go.dev/context)
that must be parsed to extract the identity of the current request,
and `group` that gives the group assigned to the current annotation.

```
type Auth interface {
	Annotate(ctx context.Context, group string) error
}
```


