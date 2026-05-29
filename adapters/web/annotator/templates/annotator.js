{{define "annotator"}}

document.addEventListener('alpine:init', () => {

    Alpine.store('labelModal', {
        show: false,
        selectedItem: "",
        open() { this.show = true },
        close() { this.show = false },
        isOpen() { return this.show }
    });

    Alpine.store('annotator', {
        instance: null,
        lastCreatedAnnotation: null,

        setInstance(annotator) {
            this.instance = annotator;
        },

        setLastCreated(annotation) {
            this.lastCreatedAnnotation = annotation;
        }
    });

    const AnnotationAPI = {
        async fetchAll() {
            const res = await fetch(`/ui/annotations?id={{.ImageId}}&collection={{.Collection}}`);
            if (!res.ok) throw new Error('Could not fetch annotations');
            return res.json();
        },
        async setLabel(id, label) {
            const res = await fetch(`/ui/set-label?id=${id}&label=${label}`);
            if (!res.ok) throw new Error('Could not relabel');
        },

        async submit(label, annotation) {
            const res = await fetch("/ui/submit-box", {
                method: "POST",
                headers: { "Content-type": "application/json; charset=UTF-8" },
                body: JSON.stringify({
                    image_id: "{{.ImageId}}",
                    collection: "{{.Collection}}",
                    label,
                    annotation
                })
            });

            if (!res.ok) throw new Error('Could not submit annotation');
        },

        async remove(id) {
            const res = await fetch(`/ui/remove-annotation?id=${id}`);
            if (!res.ok) throw new Error('Could not remove annotation');
        },

        async update_box(annotation) {
            const res = await fetch("/ui/update-box", {
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

        registerEvents(annotator) {
            annotator.on('createAnnotation', (annotation) => {
                Alpine.store("annotator").setLastCreated(annotation);
                Alpine.store("labelModal").open();
            });

            annotator.on('updateAnnotation', (updated) => {
                AnnotationAPI.update_box(updated)
            });

            annotator.on('selectionChanged', (annotations) => {
                // console.log("Selected annotations", annotations);
            });
            annotator.on('mouseEnterAnnotation', (annotation) => {
            console.log('Mouse entered: ' + annotation.id);
            });
            annotator.on('mouseLeaveAnnotation', (annotation) => {
            console.log('Mouse left: ' + annotation.id);
            });
        },

        async draw() {
            try {
                const data = await AnnotationAPI.fetchAll();
                const annotator = Alpine.store("annotator").instance;
                annotator.setAnnotations(data, true);
            } catch (err) {
                console.error(err);
            }
        },

        async submit(label) {
            try {
                const store = Alpine.store("annotator");
                await AnnotationAPI.submit(label, store.lastCreatedAnnotation);
                Alpine.store("labelModal").close();
                await this.refreshUI();

            } catch (err) {
                console.error(err);
                alert(err.message);
            }
        },

        async relabel(id, label) {
            try {
                await AnnotationAPI.setLabel(id, label)
            } catch(err) {
                console.log(err);
                alert(err.message)
            }
        },

        async remove(id) {
            try {
                await AnnotationAPI.remove(id);
                await this.refreshUI();
            } catch (err) {
                console.error(err);
                alert(err.message);
            }
        },

        async refreshUI() {
            this.refreshList();
            await this.draw();
        },

        refreshList() {
            htmx.ajax(
                'GET',
                `/ui/annotation-panel?id={{.ImageId}}&collection={{.Collection}}`,
                '#annotation-list'
            );
        },

        abort() {
            Alpine.store("labelModal").close();
            this.draw();
        }
    };

    window.AnnotatorModule = AnnotatorModule;

});

window.addEventListener('load', () => {
    window.AnnotatorModule.init();
});

{{end}}
