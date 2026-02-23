
{{define "annotator"}}

document.addEventListener('alpine:init', () => {
    Alpine.store('labelModal', {
        show:false,
        selectedItem: "",
        isOpen() {
            return this.show;
        },
        open(){
            this.show=true;
        },
        close(){
            this.show =false;
        },
    })
    Alpine.store('annotator', {annotations: {{.Annotations}}})
})

function myStyler(annotation, state) {
    if (annotation.hasOwnProperty('properties')) {
        if(annotation.properties.hasOwnProperty('color')) {
            var style = {
                fill: '#ffff',
                fillOpacity: 0.1,
                stroke: annotation.properties.color,
                strokeOpacity: 1,
                strokeWidth: 2
            }
            return style

        }
    }
}

function newAnnotator (){
    annotator = Annotorious.createImageAnnotator('image',
                    {userSelectAction: 'SELECT',
                    drawingEnabled: {{if .EnableAnnotation}} true {{else}} false {{end}}});

    annotator.on('selectionChanged', (annotations) => {
        if (annotations.length == 0) {
            modifyAnnotation(Alpine.store("annotator").lastSelectedId)
        }
    });
    Alpine.store("annotator").annotator = annotator

    annotator.setStyle(myStyler);
    setAnnotations()

    registerEventsOnAnnotator();

    return annotator;
}

function registerEventsOnAnnotator() {
    annotator = Alpine.store("annotator").annotator
    annotator.on('createAnnotation', (annotation) => {
        openLabelModal();
    });

}

function setAnnotations() {
    annotations = Alpine.store("annotator").annotations
    annotator = Alpine.store("annotator").annotator
    if (annotator.getAnnotations().length > 0) {
        annotator.clearAnnotations();
    }

    if(annotations !== null){
        annotator.setAnnotations(annotations, replace=true);
    }
}

function editAnnotation(annotationId){
    annotator = Alpine.store("annotator").annotator
    Alpine.store("annotator").lastSelectedId = annotationId
    annotator.setSelected(annotationId, true);

}

window.onload = function() {
    var annotator = newAnnotator();
    Alpine.store("annotator").annotator = annotator
}

function closeLabelModal(){
    Alpine.store("labelModal").close();
    setAnnotations()
}

function openLabelModal(){
    Alpine.store("labelModal").open();
}

function upsertAnnotation(label, annotation) {
    var body = {"image_id": "{{.ImageId}}",
                "collection_id": "{{.CollectionId}}",
                "label": label,
                "annotation": annotation};
    var headers = {
            "Content-type": "application/json; charset=UTF-8"
    };

    fetch("/ui/submit-annotation?image_id={{.ImageId}}&collection_id={{.CollectionId}}&origin_entity={{.OriginType}}&origin_id={{.OriginId}}&ordering={{.Ordering}}&descending={{.Descending}}", {
        method: "POST",
        body: JSON.stringify(body),
        headers: {
            "Content-type": "application/json; charset=UTF-8"
        }
    }).then(response => {
            if (response.ok) {
                return response.json()
            }
            throw new Error('Could not upsert annotation')
        })
    .then(data =>{
        Alpine.store("annotator").annotations = data
        setAnnotations();
    })
    .then(() => {
        closeLabelModal();
    })
    .then(() => {
        redrawAnnotationList();
    })
    .catch((error) => {
        console.log(error)

    });

}

function submitNewAnnotation(label){
    const annotator = Alpine.store("annotator").annotator
    const currentAnnotations = annotator.getAnnotations()
    const lastDrawnAnnotation = currentAnnotations[currentAnnotations.length - 1]

    upsertAnnotation(label, lastDrawnAnnotation.target);
}

function modifyAnnotation(annotationId) {
    const annotator = Alpine.store("annotator").annotator
    const annotations = annotator.getAnnotations()
    annotations.forEach((a) => {
        if (a.id == annotationId) {
            upsertAnnotation(a.properties.label, a.target)
        }
    })

}

function redrawAnnotationList(){
    htmx.ajax('GET',
            '/ui/annotation-panel?image_id={{.ImageId}}&collection_id={{.CollectionId}}&origin_entity={{.OriginType}}&origin_id={{.OriginId}}&ordering={{.Ordering}}&descending={{.Descending}}',
            '#annotation-panel')
    location.hash = "#image-annotation-panel";

}

function deleteAnnotation(annotationId){
    var url = "/ui/delete-annotation?image_id={{.ImageId}}&collection_id={{.CollectionId}}&origin_entity={{.OriginType}}&origin_id={{.OriginId}}&ordering={{.Ordering}}&descending={{.Descending}}&annotation_id=" + annotationId
    const annotator = Alpine.store("annotator").annotator

    fetch(url,
          {method: "DELETE"})
    .then(response => {
            if (response.ok) {
                return response.json()
            }
            throw new Error('Could not delete annotation')
    })
    .then(data =>{
        annotator.clearAnnotations();
        Alpine.store("annotator").annotations = [];
        if (data !== null) {
            annotator.setAnnotations(data);
        }
    })
    .then(() => {
        redrawAnnotationList();
    })
    .catch((error) => {
        console.log(error)
    });
}

{{end}}
