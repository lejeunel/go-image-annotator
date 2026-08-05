{{define "annotator"}}

const endpoints = {
    fetchAnnotations : "{{.URLs.FetchAnnotations}}",
    setLabel : "{{.URLs.SetLabel}}",
    submitImageLabel : "{{.URLs.SubmitImageLabel}}",
    submitBox : "{{.URLs.SubmitBox}}",
    submitPolygon : "{{.URLs.SubmitPolygon}}",
    removeAnnotation : "{{.URLs.RemoveAnnotation}}",
    updateBox : "{{.URLs.UpdateBox}}",
    updatePolygon : "{{.URLs.UpdatePolygon}}",
};

function newURLFromString(urlString) {
    return new URL(urlString, window.location.origin)
}

async function apiFetch(url, options, defaultError) {
    const res = await fetch(url, options);

    if (!res.ok) {
        const message = (await res.text()).trim();
        throw new Error(message || defaultError);
    }

    return res;
}

document.addEventListener('alpine:init', () => {
    Alpine.store('imageLabelModal', {
        show: false,
        selectedItem: "",
        open() { this.show = true },
        close() { this.show = false },
        isOpen() { return this.show }
    });

    Alpine.store('regionLabelModal', {
        show: false,
        selectedItem: "",
        open() { this.show = true },
        close() { this.show = false },
        isOpen() { return this.show }
    });

    Alpine.store('annotator', {
        instance: null,
        lastCreatedAnnotation: null,
        currentDrawingShape: "rectangle",

        setInstance(annotator) {
            this.instance = annotator;
        },
        DrawWithPolygon() {
            this.currentDrawingShape = "polygon";
        },
        DrawWithRectangle() {
            this.currentDrawingShape = "rectangle";
        },

        setLastCreated(annotation) {
            this.lastCreatedAnnotation = annotation;
        }
    });

    const AnnotationAPI = {
        async fetchAllAnnotations() {
            const url = newURLFromString(endpoints.fetchAnnotations)
            url.searchParams.set("id", "{{.ImageId}}")
            url.searchParams.set("collection", "{{.Collection}}")
            const res = await fetch(url.toString());
            if (!res.ok) {
                const message = await res.text();
                throw new Error(message || "Could not fetch annotations");
            }
            return res.json();
        },
        async setLabelToAnnotation(id, label) {
            const url = newURLFromString(endpoints.setLabel)
            url.searchParams.set("id", id)
            url.searchParams.set("label", label)
            await apiFetch(url.toString(), {method: "POST"}, "Could not update image label");
        },
        async addImageLabel(label) {
            const url = newURLFromString(endpoints.submitImageLabel)
            url.searchParams.set("label", label)
            url.searchParams.set("image_id", "{{.ImageId}}")
            url.searchParams.set("collection", "{{.Collection}}")
            await apiFetch(url.toString(), {method: "POST"},
                          "Could not submit image label")
        },
        async submitBox(label, annotation) {
            await apiFetch(endpoints.submitBox.toString(),
                {
                    method: "POST",
                    headers: { "Content-type": "application/json; charset=UTF-8" },
                    body: JSON.stringify({
                        image_id: "{{.ImageId}}",
                        collection: "{{.Collection}}",
                        label,
                        annotation
                    })
                },
                          "Could not submit bounding-box");
        },
        async submitPolygon(label, annotation) {
            await apiFetch(
                endpoints.submitPolygon,
                {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify({
                        image_id: "{{.ImageId}}",
                        collection: "{{.Collection}}",
                        label,
                        annotation
                    })
                },
                "Could not submit polygon"
            );
        },

        async remove(id) {
            const url = newURLFromString(endpoints.removeAnnotation)
            url.searchParams.set("id", id)
            await apiFetch(url.toString(), {method: 'DELETE'}, "Could not remove annotation");
        },

        async updateBox(annotation) {
            await apiFetch(endpoints.updateBox.toString(), {
                method: "PUT",
                headers: { "Content-type": "application/json; charset=UTF-8" },
                body: JSON.stringify(annotation),
            }, "Could not update bounding-box");
        },
        async updatePolygon(annotation) {
            await apiFetch(endpoints.updatePolygon.toString(), {
                method: "PUT",
                headers: { "Content-type": "application/json; charset=UTF-8" },
                body: JSON.stringify(annotation),
            }, "Could not update polygon");
        }

    };

    function styler(annotation) {
        const color = annotation?.properties?.color;
        if (!color) return;

        return {
            fill: '#ffff',
            fillOpacity: 0.1,
            stroke: color,
            strokeOpacity: 1,
            strokeWidth: 2
        };
    }

    const AnnotatorModule = {

        init() {
            const annotator = Annotorious.createImageAnnotator('image', {
                userSelectAction: 'EDIT',
                drawingEnabled: {{if .EnableAnnotation}} true {{else}} false {{end}}
            });

            annotator.setStyle(styler);

            this.registerEvents(annotator);
            Alpine.store("annotator").setInstance(annotator);

            this.draw();

            return annotator;
        },
        drawPolygon(){
            const annotator = Alpine.store("annotator").instance;
            annotator.setDrawingTool('polygon');
            Alpine.store("annotator").DrawWithPolygon();
        },
        drawRectangle(){
            const annotator = Alpine.store("annotator").instance;
            annotator.setDrawingTool('rectangle');
            Alpine.store("annotator").DrawWithRectangle();
        },
        registerEvents(annotator) {
            annotator.on('createAnnotation', (annotation) => {
                Alpine.store("annotator").setLastCreated(annotation);
                Alpine.store("regionLabelModal").open();
            });

            annotator.on('updateAnnotation', async (updated) => {
                switch(updated.target.selector.type){
                case "RECTANGLE":
                    try {
                        await AnnotationAPI.updateBox(updated);
                        await this.refreshUI();
                    } catch (err) {
                        notify("danger", "updating bounding-box", err.message);
                    }
                    break;
                case "POLYGON":
                    try {
                        await AnnotationAPI.updatePolygon(updated);
                        await this.refreshUI();
                    } catch (err) {
                        notify("danger", "updating polygon", err.message);
                    }
                    break;
                default:
                    notify("danger", "updating annotation", "selector type " + updated.target.selector.type + " not recognized! Should be RECTANGLE or POLYGON")
                }
            });

            annotator.on('selectionChanged', (annotations) => {
            });
            annotator.on('mouseEnterAnnotation', (annotation) => {
            });
            annotator.on('mouseLeaveAnnotation', (annotation) => {
            });
        },

        async draw() {
            try {
                const data = await AnnotationAPI.fetchAllAnnotations();
                const annotator = Alpine.store("annotator").instance;
                annotator.setAnnotations(data, true);
            } catch (err) {
                notify("danger", "drawing annotator", err.message);
            }
        },

        async submitImageLabel(label) {
            try {
                await AnnotationAPI.addImageLabel(label);
                Alpine.store("imageLabelModal").close();
                await this.refreshUI();

            } catch (err) {
                notify("danger", "submiting image label",  err.message);
            }
        },

        async submitRegion(label) {
            try {
                const store = Alpine.store("annotator");
                if (store.currentDrawingShape === "rectangle"){
                    await AnnotationAPI.submitBox(label, store.lastCreatedAnnotation);
                } else {
                    await AnnotationAPI.submitPolygon(label, store.lastCreatedAnnotation);
                }
                Alpine.store("regionLabelModal").close();
                await this.refreshUI();

            } catch (err) {
                notify("danger", "submitting region", err.message);
            }
        },

        async relabel(id, label) {
            try {
                await AnnotationAPI.setLabelToAnnotation(id, label)
                await this.refreshList();
            } catch(err) {
                notify("danger", "modifying label", err.message)
            }
        },

        async remove(id) {
            try {
                await AnnotationAPI.remove(id);
                await this.refreshUI();
            } catch (err) {
                notify("danger", "deleting annotation", err.message)
            }
        },

        async refreshUI() {
            await this.refreshList();
            await this.draw();
        },

        async refreshList() {
            htmx.ajax(
                'GET',
                `{{.URLs.AnnotationPanel}}?id={{.ImageId}}&collection={{.Collection}}`,
                '#annotation-list'
            );
        },

        abort() {
            Alpine.store("regionLabelModal").close();
            Alpine.store("imageLabelModal").close();
            this.draw();
        }
    };

    window.AnnotatorModule = AnnotatorModule;

});

window.addEventListener('load', () => {
    window.AnnotatorModule.init();
});

{{end}}
