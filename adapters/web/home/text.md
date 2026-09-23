## Quickstart

To get started annotating images, ensure that
you have the permissions needed to create a collection,
ingest images, and annotate. These are displayed in your [dashboard](/dashboard).
Ask your admin for help if this is not the case.

### Creating the recipient collection

Prior to contributing your first images, you must first create
a collection in which they will be contained.

To do so, head to the [collection](/collections) page
and create one with a name of your choosing, e.g. **my-first-collection**.
You may omit the other fields for now.

### Ingesting images

From your local machine, create a zip archive that contains
your images. Next, head to the [listing of your new collection](/images?collection=my-first-collection),
and open the **Ingest** dialog. You may drop here your zip-archive.

A background process will start once the archive is finished uploading. You may inspect
its progression in your [logs](/dashboard/logs).

### Adding labels

Annotations can be seen as a relation between an image (or a region of an image) and a label.
Usually, these must be named according to your downstream use-cases, e.g. if you are
classifying animals, your label could be **dog**, **cat**, **pigeon**, ...

Go ahead and add a couple labels that you wish to use on your images.

### Annotating

Now we are all set to start annotating. Head back to the [image listing](/images?collection=my-first-collection),
and click on the first entry. This will open the annotator view.
You may choose to assign the label to the whole image or to a sub-region (polygon or bounding-box).

For the former, simply click on the **+** button on the **Labels** frame.

For the latter, you must first draw a region by clicking on the image. A pop-up will appear to let you
choose the corresponding label.



