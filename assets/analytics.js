/* eslint-env browser */

(function(window) {
    var holmesId = window.localStorage.getItem('_holmesId'),
        enrichers = [],
        loadedEvent;

    if(holmesId === null) {
        window.localStorage.setItem('_holmesId', '__HOLMES_ID__');
    }

    function track(trackingObject) {
        trackingObject['holmesId'] = window.localStorage.getItem('_holmesId');

        for(var i=0; i < enrichers.length; i++) {
            enrichers[i](trackingObject);
        }

        var url = '__HOLMES_BASE_URL__'
            + '/track?u='
            + new Date().getTime()
            + '&t=' + encodeURIComponent(JSON.stringify(trackingObject));
        var xhr = new XMLHttpRequest();
        xhr.open('GET', url);
        xhr.send();
    }

    window.Holmes = {
        pageView: function(trackingObject) {
            trackingObject['type'] = 'PAGE_VIEW';
            track(trackingObject);
        },

        addTrackingEnricher: function(trackingEnricher) {
            enrichers.push(trackingEnricher);
        },

        track: track
    };

    loadedEvent = document.createEvent('Event');
    loadedEvent.initEvent('holmesloaded', false, false);
    window.dispatchEvent(loadedEvent);
}(window));
