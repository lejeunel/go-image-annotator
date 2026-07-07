{{define "annotator"}}

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
        async fetchAll() {
            const res = await fetch(`/ui/annotate/annotations?id={{.ImageId}}&collection={{.Collection}}`);
            if (!res.ok) throw new Error('Could not fetch annotations');
            return res.json();
        },
        async setLabel(id, label) {
            const res = await fetch(`/ui/annotate/set-label?id=${id}&label=${label}`);
            if (!res.ok) throw new Error('Could not relabel');
        },
        async submit_label(label) {
            const res = await fetch(`/ui/annotate/submit-label?image_id={{.ImageId}}&collection={{.Collection}}&label=${label}`)
            if (!res.ok) throw new Error('Could not submit annotation');
        },
        async submit_box(label, annotation) {
            const res = await fetch("/ui/annotate/submit-box", {
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
        async submit_polygon(label, annotation) {
            const res = await fetch("/ui/annotate/submit-polygon", {
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
            const res = await fetch(`/ui/annotate/remove-annotation?id=${id}`);
            if (!res.ok) throw new Error('Could not remove annotation');
        },

        async update_box(annotation) {
            const res = await fetch("/ui/annotate/update-box", {
                method: "POST",
                headers: { "Content-type": "application/json; charset=UTF-8" },
                body: JSON.stringify(annotation),
            });
            if (!res.ok) throw new Error('Could not update annotation');
        },
        async update_polygon(annotation) {
            const res = await fetch("/ui/annotate/update-polygon", {
                method: "POST",
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
                    AnnotationAPI.update_box(updated);
                    break;
                case "POLYGON":
                    AnnotationAPI.update_polygon(updated);
                    break
                default:
                    alert("selector type " + updated.target.selector.type + " not recognized! Should be RECTANGLE or POLYGON")
                }

                try {
                    this.refreshUI();

                } catch (err) {
                    alert(err.message);
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
                const data = await AnnotationAPI.fetchAll();
                const annotator = Alpine.store("annotator").instance;
                annotator.setAnnotations(data, true);
            } catch (err) {
                alert(err.message);
            }
        },

        async submit_label(label) {
            try {
                await AnnotationAPI.submit_label(label);
                Alpine.store("imageLabelModal").close();
                await this.refreshUI();

            } catch (err) {
                alert(err.message);
            }
        },

        async submit_region(label) {
            try {
                const store = Alpine.store("annotator");
                if (store.currentDrawingShape === "rectangle"){
                    await AnnotationAPI.submit_box(label, store.lastCreatedAnnotation);
                } else {
                    await AnnotationAPI.submit_polygon(label, store.lastCreatedAnnotation);
                }
                Alpine.store("regionLabelModal").close();
                await this.refreshUI();

            } catch (err) {
                alert(err.message);
            }
        },

        async relabel(id, label) {
            try {
                await AnnotationAPI.setLabel(id, label)
                await this.refreshList();
            } catch(err) {
                alert(err.message)
            }
        },

        async remove(id) {
            try {
                await AnnotationAPI.remove(id);
                await this.refreshUI();
            } catch (err) {
                alert(err.message);
            }
        },

        async refreshUI() {
            this.refreshList();
            await this.draw();
        },

        async refreshList() {
            htmx.ajax(
                'GET',
                `/ui/annotate/annotation-panel?id={{.ImageId}}&collection={{.Collection}}`,
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
