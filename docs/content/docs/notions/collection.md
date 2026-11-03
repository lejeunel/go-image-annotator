---
title: Collections
description: >
  Where we explain how collections can be used to group images,
  configure annotation tasks, and keep track of your work.
weight: 2
linkTitle: Collections
---

At a first level, a collection is akin to
a directory in a file-system tree, in it allows
to group images into a named unit that correspond to an intent.

For instance, you might name a collection `flower-classification`,
in which you would add images of flowers, and pursue with
associating to each image a label denoting its name.

We now elaborate on more advanced features supported by collections,
namely group-ownership, cloning, and profile assignment.

## Group ownership

Also analogous to a folder in a file-system, 
collections can be assigned a group-ownership, as further developped in
the [authorization]({{< ref "authorize.md">}}) section.

This allows to protect your data from non-intended modifications.
For instance, you could restrict modifications of your `flower-classification`
collection to group `botanist`.

## Annotation profile

As a collection is aimed at a specific downstream application, it 
makes sense to configure it with a set of pre-defined labels.

We call this set a `profile`, which should be ideally named
after your downstream use-case. Following the above scenario,
we could create a `flowers` profile which would contain 
labels `common-poppy`, `tulip`, `snow-drop`, etc..

Collections are not required to have an associated profile,
in which case all existing labels will be made available
for annotation.

## Using images across contexts

An important use-case that we wish to cover is when one
wants to use the same image accross different contexts.
For instance, say that we would want to label an image
in a multi-class classification context. We have named its containing
collection `my-multi-class-collection`.

Later on, we wish to use that same image in a multi-label context (several possible
labels per-image). Of course, for reproducibility, we want to keep
intact the first set of annotations.

We achieve this by assigning that image to a new collection `my-multi-label-collection`.
Modifying the latter will leave untouched the annotations applied in `my-multi-class-collection`.

## Versioning collections by cloning

Another common use-case is when we wish to iterate on a dataset, for instance modifying
annotations to take into account defects, or adding new images.

Again, this could be achieved simply by modifying the initial collection in-place. 
For the sake or tracability and reproducibility, we might want to keep the initial collection untouched
and work on a copy.

To do so, the same principle that we saw [above]({{< ref "collection.md#using-images-across-contexts" >}}) applies:
You would clone your initial collection into a new one.
Importantly, you will want to transfer not only images, but also its annotations and meta-data.
To do so, make sure to active the **deep** flag.
