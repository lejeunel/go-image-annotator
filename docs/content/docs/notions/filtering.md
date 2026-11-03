---
title: Slices
description: >
  Describe how you can generate views using filtering and sorting queries.
weight: 3
linkTitle: Slices
---

We have already seen how collections allow to group images together.
In practice, this solution might be limited in that
collections are not easy to split and combine. For instance,
you might need to quicly scroll through images possibly belonging
to several collections.

To do so, you will resort to what we call a *slice*, i.e. a view
on all images contained in the database. Slices can be generated on-the-fly
using a custom query language. They allow you to use pre-built 
filtering criteria, as well as using your own meta-data. Last,
you may provide ordering criteria to sort your images.

## Filtering

A filter is a string of text that contains one or several
selection criteria combined using logical operators.

As an example, consider:

```
collection=flower-dataset and (meta.size-cm>=5 and meta.color=red)
```

The first term selects all images in collection `flower-dataset`. The second
term further restricts the result by picking flowers whose size is above a given
value and are red.

## Sorting

Results can be ordered by concatenating field names along with 
keywords `asc` and `desc`, which stand for *ascending* and *descending*,
respectively.

Note that the ordering of terms matter here.

As an example, take:

```
meta.size-cm:asc ingested_at:desc 
```

The above orders images by their size, with the smallest coming first.
When several images have the same size, they are ordered such that the
last ingested image appears first.


## Scrolling

Now that you have created the slice of your liking, you will
be presented a table of results. Importantly,
clicking on an image link will bring you to the image viewer,
where your filtering and ordering criteria will be set.

This means that your filtering and ordering criteria will
affect the behaviour of the `next` and `previous` navigation buttons.
