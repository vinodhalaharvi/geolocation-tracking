let drawnPolygons = []; // Assume this is filled with your drawn polygons
let carMarkers = new Map(); // Stores markers with car IDs as keys
let map; // Declare `map` globally

function initWebSocket() {
    console.log("initWebSocket");
    const socket = new WebSocket('ws://localhost:8080/ws');

    socket.onmessage = function (event) {
        const data = JSON.parse(event.data);
        console.log(data); // Assuming data is an array of car objects

        data.forEach(car => {
            const carLocation = new google.maps.LatLng(car.location.latitude, car.location.longitude);
            let isInsidePolygon = false;

            // Assuming `drawnPolygons` is an array of polygons with their `id` properties set
            const assignedPolygon = drawnPolygons.find(polygon => polygon.get('id') === car.polygonId);
            if (assignedPolygon && google.maps.geometry.poly.containsLocation(carLocation, assignedPolygon)) {
                isInsidePolygon = true;
            }

            const markerIcon = {
                path: google.maps.SymbolPath.CIRCLE,
                fillColor: isInsidePolygon ? '#0000FF' : '#FF0000', // Blue if inside, red if outside
                fillOpacity: 0.6,
                strokeColor: '#FFFFFF',
                strokeWeight: 2,
                scale: 6, // Size of the dot
            };

            if (carMarkers.has(car.id)) {
                // Move the existing marker and change color if needed
                const marker = carMarkers.get(car.id);
                marker.setPosition(carLocation);
                marker.setIcon(markerIcon);
            } else {
                // Create a new marker
                const carMarker = new google.maps.Marker({
                    position: carLocation,
                    map: map,
                    icon: markerIcon,
                });
                carMarkers.set(car.id, carMarker);
            }
        });
    };

    socket.onclose = function () {
        console.log("WebSocket connection closed. Attempting to reconnect...");
        // setTimeout(initWebSocket, 1000); // Adjust delay as appropriate
    };

    socket.onerror = function (err) {
        console.error("WebSocket error observed:", err);
    };
}

function initMap() {

    map = new google.maps.Map(document.getElementById("map"), {
        zoom: 13, center: {lat: 38.9072, lng: -77.0369}, // Washington DC coordinates
        mapTypeId: "terrain",
    });

    const drawingManager = new google.maps.drawing.DrawingManager({
        drawingMode: google.maps.drawing.OverlayType.POLYGON, drawingControl: true, drawingControlOptions: {
            position: google.maps.ControlPosition.TOP_CENTER, drawingModes: ['polygon']
        }, polygonOptions: {
            strokeColor: "#FF0000",
            strokeOpacity: 0.8,
            strokeWeight: 2,
            fillColor: "#FF0000",
            fillOpacity: 0.35,
            editable: true,
            draggable: true,
        },
    });

    google.maps.event.addListener(drawingManager, 'drawingmode_changed', function () {
        if (drawingManager.getDrawingMode()) {
            document.getElementById("map").classList.add("custom-cursor");
        } else {
            document.getElementById("map").classList.remove("custom-cursor");
        }
    });
    drawingManager.setMap(map);

    // Listen for the Escape key to exit drawing mode
    document.addEventListener('keyup', function (event) {
        if (event.key === "Escape") {
            drawingManager.setDrawingMode(null);
        }
    });


    // Inside initMap function, after initializing drawingManager
    google.maps.event.addListener(drawingManager, 'overlaycomplete', function (event) {
        console.log("overlaycomplete")
        if (event.type === google.maps.drawing.OverlayType.POLYGON) {
            // Get the coordinates of the polygon's path
            const vertices = event.overlay.getPath().getArray().map(vertex => ({
                latitude: vertex.lat(), longitude: vertex.lng()
            }));

            // Add the polygon to the global array
            drawnPolygons.push(event.overlay);

            // Send the coordinates to the backend
            fetch('/polygons', {
                method: 'POST', headers: {
                    'Content-Type': 'application/json',
                }, // add random polygon id
                body: JSON.stringify({points: vertices}),
            })
                .then(response => {
                    return response.json();
                })
                .then(data => {
                    console.log('Polygon saved with ID:', data.id);
                    event.overlay.set('id', data.id); // Associate the returned ID with the polygon
                    // Attach a click listener to the polygon
                    google.maps.event.addListener(event.overlay, 'click', function (mapsMouseEvent) {
                        // Reuse the logic you have for handling clicks on the map
                        let clickedLat = mapsMouseEvent.latLng.lat();
                        let clickedLng = mapsMouseEvent.latLng.lng();

                        // Logic to handle the click inside the polygon
                        let polygonId = event.overlay.get('id'); // Assuming your polygons have an 'id' property

                        fetch('/addAsset', {
                            method: 'POST', headers: {
                                'Content-Type': 'application/json',
                            }, body: JSON.stringify({
                                polygonId: polygonId, location: {latitude: clickedLat, longitude: clickedLng}
                            }),
                        })
                            .then(response => response.json())
                            .then(data => {
                                console.log('Asset added:', data);
                                const carId = data.id; // Assuming your backend returns a unique ID for the car
                                const carLocation = {lat: data.location.latitude, lng: data.location.longitude};

                                if (carMarkers.has(carId)) {
                                    // Correctly use Map's get method to access existing marker
                                    const marker = carMarkers.get(carId);
                                    marker.setPosition(carLocation);
                                } else {
                                    // Create a new marker for the car
                                    const carMarker = new google.maps.Marker({
                                        position: carLocation,
                                        map: map, // Ensure the marker is added to the same map
                                        icon: {
                                            path: google.maps.SymbolPath.CIRCLE,
                                            fillColor: 'transparent',
                                            fillOpacity: 0,
                                            strokeColor: '#0955ce',
                                            strokeWeight: 4,
                                            scale: 3, // Adjust size of the halo slightly bigger than the dot
                                        },
                                        animation: google.maps.Animation.DROP, // Optional: Adds animation when marker is added to map
                                    });
                                    // Correctly use Map's set method to update the carMarkers map
                                    carMarkers.set(carId, carMarker);
                                }
                            })
                            .catch(error => console.error('Error adding car:', error));
                    });
                    drawnPolygons.push(event.overlay); // Add the polygon to the global array
                })
                .catch(error => console.error('Error saving polygon:', error));

            // Optional: Remove the polygon after saving, or handle as needed
            // event.overlay.setMap(null);
        }
    });

    document.getElementById('simulateBtn').addEventListener('click', () => {
        clearAssetMarkers(); // Clear existing car markers before starting a new simulation

        initWebSocket();

        // Send a simulate request to the backend
        fetch('/simulate', {
            method: 'POST', headers: {
                'Content-Type': 'application/json',
            }, body: JSON.stringify({action: "start"}),
        })
            .then(response => response.json())
            .then(data => {
                console.log('Simulation started:', data);
            })
            .catch(error => console.error('Error starting simulation:', error));
    });

}

window.initMap = initMap;

function clearAssetMarkers() {
    carMarkers.forEach((marker) => {
        marker.setMap(null); // Removes the marker from the map
    });
    carMarkers.clear(); // Clears the Map, removing all key-value pairs
}

