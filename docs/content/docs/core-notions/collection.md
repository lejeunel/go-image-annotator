---
title: Images & Collections
description: >
  Where we explain how to organize images
weight: 2
linkTitle: Collections
menu: { main: { weight: 20 } }
---

We describe here the major concepts that allow to use this application,
namely how to organize image-data into collections, define labels,
annotate images, and versioning annotations.
Also, we elaborate on the notion of 

# Data Organization

## Images

This combines raw image-data as well as related meta-data. Importantly,
all images must be contained in *at least* one collection (see below).

## Collections

This entity serves to group images together.
Importantly, the same image can appear in several collections,
and in each one, carry different annotations.
In other words, annotations are a property of a `(image,collection)` couple.
