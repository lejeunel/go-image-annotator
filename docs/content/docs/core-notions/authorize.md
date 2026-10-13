---
title: Access Restriction
description: >
  Where we explain how to restrict access to your resources
weight: 4
linkTitle: Authorization
menu: { main: { weight: 20 } }
---

To safely collaborate whithin a team, it is often
necessary to restrict access according to user privileges.

Privileges are managed using
both **groups** and **roles**, which are defined below.

## Groups

### Definition

Groups are assigned to resources and users. When 
the user is member of group **g**, and the resource is
assigned to group **g**, we say that the user has the 
necessary group-membership to work with said resource.

### Resources

We allow to set group-ownership to `collection` and `profile`
resources.

Regarding `profiles`, this protects the modification
of the `profile` itself, i.e. its name and description, as well
as the modification of the labels assigned to said profile.

As for `collections`, we protect the modification of the 
`collection` itself, i.e. its name and description, and also
all operations on the images it contains. This means that 
if collection `c` is owned by group `g` and contains image `i`, 
modifying annotations of image `i` would require to be member of `g`.

## Roles

### Definition

Roles are assigned to users and are effectively a 
serie of allowed methods, e.g. `CreateCollection`, `Annotate`, ...
We say that a user has the necessary role to perform method **m**
when they are assigned a role that contains method **m**.

### Defaults

On a fresh instance, a number of pre-defined roles are defined that
cover basic uses-cases, namely:

- `annotator`: This role allows to perform create, read, update, and delete on 
  annotations, be it image-wise or region-wise. Also, we allow to create/update and
  delete meta-data.
- `image-contrib`: This role allows to create collections and ingest new images.
  Also, it allows to clone collections.
- `admin`: This role grants every access to every resource. It bypasses the specified
  group-ownership restrictions, meaning that an `admin` does not need to be a member of
  a group to work the corresponding resource.
  
You may remove and/or modify the above at will, and also create new roles as you see fit.
In particular, you can select pre-defined methods for a given role.
