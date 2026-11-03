---
title: Annotating
description: >
  Where we describe different ways to annotate images
weight: 3
linkTitle: Annotations
---

An annotation is a user-defined piece of information about an image. Their
actual use depends on the downstream application but typically lie in:

- Building training and evaluation datasets for Machine Learning applications.
- Retrieving specific subsets of images for a review task, e.g. data quality check.

We provide two approaches:

1. By explicitly adding before-hand a set of **labels** and assigning them to images of parts of images.
2. Assign "free-form" meta-data entry to an image.

We now elaborate on the specificity and differences of these two approaches.

## Using Labels

The label-based approach is suited when you need to inform about the existence, or non-existence 
of objects **that are known in advance** on an images, or on parts of an image.

For example, you have pre-defined set of flower genre, e.g. `common-poppy`, `tulip`, `snow-drop`,
and a set of images that contain zero or more instances of them.

This approach therefore requires to create one `label` object for each genre,
and pick for each image, one or more label from the menu.

## Using Free-form Meta-data

This approach is best suited when you want sub-region granularity, e.g. 
when you need to segment parts of images using polygons or bounding-boxes.

Prior to annotate an image, we need to define
labels that define semantic information on what the image contains.
