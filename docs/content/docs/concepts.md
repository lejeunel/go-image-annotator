---
title: Core Concepts
description: >
  Where we explain how data are organized and how to restrict access to it
weight: 1
---

{{% pageinfo %}}

We describe here the major concepts that allow to use this application,
namely how to organize image-data into collections, define labels,
annotate images, and versioning annotations.
Also, we elaborate on the notion of 

{{% /pageinfo %}}

## Data Organization

### Images

This combines raw image-data as well as related meta-data. Importantly,
all images must be contained in *at least* one collection (see below).

### Collections

This entity serves to group images together.
Importantly, the same image can appear in several collections,
and in each one, carry different annotations.
In other words, annotations are a property of a `(image,collection)` couple.

### Labels

Prior to annotate an image, we need to define
labels that define semantic information on what the image contains.

### Annotations

An annotation is a user-defined semantic information on an image.
Currently, they can be of two kinds:

- **Image-wise**: where the annotation concerns the image as a whole
- **Region-wise**: where the annotation concerns a sub-region of the image

## Group-based Authorization

Each collection must be assigned to a **group**, which serves
to restrict access to specific users (or groups of users).

The specific of the authorization logic must be customized by 
implementing an interface, which receives a [context](https://pkg.go.dev/context)
that must be parsed to extract the identity of the current request,
and `group` that gives the group assigned to the current annotation.

```
type Auth interface {
	AnnotateGroup(ctx context.Context, group string) error
}
```


