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

function notify(variant, message, extra = {}) {
    window.dispatchEvent(new CustomEvent("notify", {
        detail: {
            variant,
            message,
            ...extra,
        },
    }));
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
            url = newURLFromString(endpoints.fetchAnnotations)
            url.searchParams.set("id", "{{.ImageId}}")
            url.searchParams.set("collection", "{{.Collection}}")
            const res = await fetch(url.toString());
            if (!res.ok) throw new Error('Could not fetch annotations');
            return res.json();
        },
        async setLabelToAnnotation(id, label) {
            url = newURLFromString(endpoints.setLabel)
            url.searchParams.set("id", id)
            url.searchParams.set("label", label)
            const res = await fetch(url.toString(), {method: "POST"});
            if (!res.ok) throw new Error('Could not relabel');
        },
        async addImageLabel(label) {
            url = newURLFromString(endpoints.submitImageLabel)
            url.searchParams.set("label", label)
            url.searchParams.set("image_id", "{{.ImageId}}")
            url.searchParams.set("collection", "{{.Collection}}")
            const res = await fetch(url.toString(), {method: "POST"})
            if (!res.ok) throw new Error('Could not submit annotation: ' + res.message);
        },
        async submitBox(label, annotation) {
            const res = await fetch(endpoints.submitBox.toString(), {
                method: "POST",
                headers: { "Content-type": "application/json; charset=UTF-8" },
                body: JSON.stringify({
                    image_id: "{{.ImageId}}",
                    collection: "{{.Collection}}",
                    label,
                    annotation
                })
            });
            if (!res.ok) throw new Error('Could not submit bounding-box');
        },
        async submitPolygon(label, annotation) {
            const res = await fetch(endpoints.submitPolygon.toString(), {
                method: "POST",
                headers: { "Content-type": "application/json; charset=UTF-8" },
                body: JSON.stringify({
                    image_id: "{{.ImageId}}",
                    collection: "{{.Collection}}",
                    label,
                    annotation
                })
            });
            if (!res.ok) throw new Error('Could not submit polygon');
        },

        async remove(id) {
            url = newURLFromString(endpoints.removeAnnotation)
            url.searchParams.set("id", id)
            const res = await fetch(url.toString(), {method: 'DELETE'});
            if (!res.ok) throw new Error('Could not remove annotation');
        },

        async updateBox(annotation) {
            const res = await fetch(endpoints.updateBox.toString(), {
                method: "PUT",
                headers: { "Content-type": "application/json; charset=UTF-8" },
                body: JSON.stringify(annotation),
            });
            if (!res.ok) throw new Error('Could not update annotation');
        },
        async updatePolygon(annotation) {
            const res = await fetch(endpoints.updatePolygon.toString(), {
                method: "PUT",
                headers: { "Content-type": "application/json; charset=UTF-8" },
                body: JSON.stringify(annotation),
            });
            if (!res.ok) throw new Error('Could not update annotation');
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

            annotator.on('updateAnnotation', (updated) => {
                switch(updated.target.selector.type){
                case "RECTANGLE":
                    AnnotationAPI.updateBox(updated);
                    break;
                case "POLYGON":
                    AnnotationAPI.updatePolygon(updated);
                    break
                default:
                    notify("error", "selector type " + updated.target.selector.type + " not recognized! Should be RECTANGLE or POLYGON")
                }
                try {
                    this.refreshUI();
                } catch (err) {
                    notify("error", err.message);
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
                notify("error", err.message);
            }
        },

        async submitImageLabel(label) {
            try {
                await AnnotationAPI.addImageLabel(label);
                Alpine.store("imageLabelModal").close();
                await this.refreshUI();

            } catch (err) {
                notify("error", err.message);
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
                notify("error", err.message);
            }
        },

        async relabel(id, label) {
            try {
                await AnnotationAPI.setLabelToAnnotation(id, label)
                await this.refreshList();
            } catch(err) {
                notify("error", err.message)
            }
        },

        async remove(id) {
            try {
                await AnnotationAPI.remove(id);
                await this.refreshUI();
            } catch (err) {
                notify("error", err.message)
            }
        },

        async refreshUI() {
            this.refreshList();
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
