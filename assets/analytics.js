/* eslint-env browser */

(function(window) {
    var holmesId = window.localStorage.getItem('_holmesId'),
        enrichers = [],
        loadedEvent;

    function getUUID() {
        var lut = []; for (var i=0; i<256; i++) { lut[i] = (i<16?'0':'')+(i).toString(16); }
        var d0 = Math.random()*0x100000000>>>0;
        var d1 = Math.random()*0x100000000>>>0;
        var d2 = Math.random()*0x100000000>>>0;
        var d3 = Math.random()*0x100000000>>>0;
        return lut[d0&0xff]+lut[d0>>8&0xff]+lut[d0>>16&0xff]+lut[d0>>24&0xff]+'-'+
            lut[d1&0xff]+lut[d1>>8&0xff]+'-'+lut[d1>>16&0x0f|0x40]+lut[d1>>24&0xff]+'-'+
            lut[d2&0x3f|0x80]+lut[d2>>8&0xff]+'-'+lut[d2>>16&0xff]+lut[d2>>24&0xff]+
            lut[d3&0xff]+lut[d3>>8&0xff]+lut[d3>>16&0xff]+lut[d3>>24&0xff];
    }

    if(holmesId === null) {
        window.localStorage.setItem('_holmesId', getUUID());
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
