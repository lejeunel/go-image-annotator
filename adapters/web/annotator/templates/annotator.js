
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
    Alpine.store('annotator', {lastCreatedAnnotation: null,
                              })
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

    annotator.on('createAnnotation', (annotation) => {
        Alpine.store("annotator").lastCreatedAnnotation = annotation
        openLabelModal();
    });
    // annotator.on('updateAnnotation', ((updated, previous) => void) => {
    //     Alpine.store("annotator").lastUpdatedAnnotation = updated
    // });

    Alpine.store("annotator").annotator = annotator

    annotator.setStyle(myStyler);
    setAnnotations()

    return annotator;
}

window.onload = function() {
    var annotator = newAnnotator();
    Alpine.store("annotator").annotator = annotator
}

function closeLabelModal(){
    Alpine.store("labelModal").close();
}

function openLabelModal(){
    Alpine.store("labelModal").open();
}

function submitAnnotation(label) {
    lastAnnotation = Alpine.store("annotator").lastCreatedAnnotation
    var body = {"image_id": "{{.ImageId}}",
                "collection": "{{.Collection}}",
                "label": label,
                "annotation": lastAnnotation};
    var headers = {
            "Content-type": "application/json; charset=UTF-8"
    };

    fetch("/ui/submit-box", {
        method: "POST",
        body: JSON.stringify(body),
        headers: {
            "Content-type": "application/json; charset=UTF-8"
        }
    })
    .then(response => {
            closeLabelModal();
            if (!response.ok) {
                throw new Error('Could not submit annotation')
            }
    })
    .then(() => {
        refreshAnnotationList();
    })
    .then(() => {
        setAnnotations();
    })
    .catch(error => {
        console.error(error);
    });
}

function refreshAnnotationList(){
    htmx.ajax('GET',
            '/ui/annotation-panel?id={{.ImageId}}&collection={{.Collection}}',
            '#annotation-list')
}
function setAnnotations(){
    fetch("/ui/annotations?id={{.ImageId}}&collection={{.Collection}}", {
        method: "GET",
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Could not fetch annotations')
        } 
        return response.json()
    })
    .then( data => {
        Alpine.store("annotator").annotator.setAnnotations(data, replace=true)
    })
    .catch(error => {
        console.error(error);
    });
}

function removeAnnotation(id) {
    fetch("/ui/remove-annotation?id="+id, {
        method: "GET",
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Could not remove annotation')
        } 
    })
    .catch(error => {
        console.error(error);
    });
    refreshAnnotationList()
    setAnnotations()
}

{{end}}
