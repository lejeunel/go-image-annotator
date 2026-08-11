{{define "annotator"}}
document.addEventListener('alpine:init', () => {
    Alpine.store('modal', {
        active: false,
        open() { this.active = true },
        close() { this.active = false },
        isOpen() { return this.active }
    });
});

const LabelPicker = (() => {
    return {
        close()         { Alpine.store('modal').close() },
        open() { Alpine.store('modal').open() },
        isOpen() { return Alpine.store('modal').isOpen() },
    };
})();

const Annotator = (() => {
    let instance = null;
    let lastCreatedAnnotation = null;

    const POLYGON_MODE = 'polygon';
    const BOX_MODE = 'rectangle';
    const IMAGE_MODE = 'image';
    const EDIT_LABEL_MODE = 'edit';

    let currentMode = BOX_MODE;



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
        return new URL(urlString, window.location.origin);
    }

    async function apiFetch(url, options, defaultError) {
        const res = await fetch(url, options);
        if (!res.ok) {
            const message = (await res.text()).trim();
            throw new Error(message || defaultError);
        }
        return res;
    }

    const AnnotationAPI = {
        async fetchAllAnnotations() {
            const url = newURLFromString(endpoints.fetchAnnotations);
            url.searchParams.set("id", "{{.ImageId}}");
            url.searchParams.set("collection", "{{.Collection}}");
            const res = await fetch(url.toString());
            if (!res.ok) {
                const message = await res.text();
                throw new Error(message || "Could not fetch annotations");
            }
            return res.json();
        },
        async setLabelToAnnotation(id, label) {
            const url = newURLFromString(endpoints.setLabel);
            url.searchParams.set("id", id);
            url.searchParams.set("label", label);
            await apiFetch(url.toString(), {method: "POST"}, "Could not update image label");
        },
        async addImageLabel(label) {
            const url = newURLFromString(endpoints.submitImageLabel);
            url.searchParams.set("label", label);
            url.searchParams.set("image_id", "{{.ImageId}}");
            url.searchParams.set("collection", "{{.Collection}}");
            await apiFetch(url.toString(), {method: "POST"}, "Could not submit image label");
        },
        async submitBox(label, annotation) {
            await apiFetch(endpoints.submitBox, {
                method: "POST",
                headers: { "Content-type": "application/json; charset=UTF-8" },
                body: JSON.stringify({ image_id: "{{.ImageId}}", collection: "{{.Collection}}", label, annotation })
            }, "Could not submit bounding-box");
        },
        async submitPolygon(label, annotation) {
            await apiFetch(endpoints.submitPolygon, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ image_id: "{{.ImageId}}", collection: "{{.Collection}}", label, annotation })
            }, "Could not submit polygon");
        },
        async remove(id) {
            const url = newURLFromString(endpoints.removeAnnotation);
            url.searchParams.set("id", id);
            await apiFetch(url.toString(), {method: 'DELETE'}, "Could not remove annotation");
        },
        async updateBox(annotation) {
            await apiFetch(endpoints.updateBox, {
                method: "PUT",
                headers: { "Content-type": "application/json; charset=UTF-8" },
                body: JSON.stringify(annotation),
            }, "Could not update bounding-box");
        },
        async updatePolygon(annotation) {
            await apiFetch(endpoints.updatePolygon, {
                method: "PUT",
                headers: { "Content-type": "application/json; charset=UTF-8" },
                body: JSON.stringify(annotation),
            }, "Could not update polygon");
        }
    };

    function styler(annotation) {
        const color = annotation?.properties?.color;
        if (!color) return;
        return { fill: '#ffff', fillOpacity: 0.1, stroke: color, strokeOpacity: 1, strokeWidth: 2 };
    }

    return {
        init() {
            instance = Annotorious.createImageAnnotator('image', {
                userSelectAction: 'EDIT',
                drawingEnabled: {{if .EnableAnnotation}} true {{else}} false {{end}}
            });
            instance.setStyle(styler);
            this.registerEvents(instance);
            this.draw();
            return instance;
        },
        editLabelMode(){
            currentMode = EDIT_LABEL_MODE;
        },
        imageLabelMode(){
            currentMode = IMAGE_MODE;
        },
        polygonMode() {
            instance.setDrawingTool('polygon');
            currentMode = POLYGON_MODE;
        },
        rectangleMode() {
            instance.setDrawingTool('rectangle');
            currentMode = BOX_MODE;
        },

        registerEvents(instance) {
            instance.on('createAnnotation', (annotation) => {
                lastCreatedAnnotation = annotation;
                switch (annotation.target.selector.type) {
                    case "RECTANGLE":
                        this.rectangleMode();
                        break;
                    case "POLYGON":
                        this.polygonMode();
                        break;
                    default:
                        notify("danger", "creating annotation",
                            `annotation type ${annotation.target.selector.type} not recognized`);
                }
                LabelPicker.open();
            });

            instance.on('updateAnnotation', async (updated) => {
                switch (updated.target.selector.type) {
                    case "RECTANGLE":
                        try { await AnnotationAPI.updateBox(updated); await this.refreshUI(); }
                        catch (err) { notify("danger", "updating bounding-box", err.message); }
                        break;
                    case "POLYGON":
                        try { await AnnotationAPI.updatePolygon(updated); await this.refreshUI(); }
                        catch (err) { notify("danger", "updating polygon", err.message); }
                        break;
                    default:
                        notify("danger", "updating annotation",
                            `selector type ${updated.target.selector.type} not recognized`);
                }
            });
        },

        async draw() {
            try {
                const data = await AnnotationAPI.fetchAllAnnotations();
                instance.setAnnotations(data, true);
            } catch (err) {
                notify("danger", "drawing annotator", err.message);
            }
        },

        setAnnotationId(id) {
            currentAnnotationId = id;
        },

        async relabel(id, label) {
            try {
                await AnnotationAPI.setLabelToAnnotation(id, label);
                await this.refreshList();
            } catch (err) {
                notify("danger", "modifying label", err.message);
            }
        },

        async submit(label) {
            try {
                switch (currentMode){
                case BOX_MODE:
                    await AnnotationAPI.submitBox(label, lastCreatedAnnotation);
                    break;
                case POLYGON_MODE:
                    await AnnotationAPI.submitPolygon(label, lastCreatedAnnotation);
                    break;
                case IMAGE_MODE:
                    await AnnotationAPI.addImageLabel(label);
                    break;
                case EDIT_LABEL_MODE:
                    await this.relabel(currentAnnotationId, label);
                    break;
                default:
                        notify("danger", "submitting annotation",
                            `selected mode ${currentMode} not recognized`);
                }
                LabelPicker.close();
                await this.refreshUI();
            } catch (err) {
                notify("danger", "submitting annotation", err.message);
            }
        },

        async remove(id) {
            try {
                await AnnotationAPI.remove(id);
                await this.refreshUI();
            } catch (err) {
                notify("danger", "deleting annotation", err.message);
            }
        },

        async refreshUI() {
            await this.refreshList();
            await this.draw();
        },

        async refreshList() {
            htmx.ajax('GET',
                `{{.URLs.AnnotationPanel}}?id={{.ImageId}}&collection={{.Collection}}`,
                '#annotation-list');
        },

        abort() {
            LabelPicker.close();
            this.draw();
        }
    };
})();

window.Annotation = {
    Annotator,
    LabelPicker,
};

window.addEventListener('load', () => window.Annotation.Annotator.init());

{{end}}
