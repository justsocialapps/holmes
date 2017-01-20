/* eslint-env browser */

(function(window) {
    var holmesId = window.localStorage.getItem('_holmesId');
    if(holmesId === null) {
        window.localStorage.setItem('_holmesId', '__HOLMES_ID__');
    }
    
    function track(trackingObject) {
        trackingObject['holmesId'] = window.localStorage.getItem('_holmesId');
        var url = '__HOLMES_BASE_URL__'
        + '/track?u='
        +new Date().getTime()
        +'&t=' + encodeURIComponent(JSON.stringify(trackingObject));
        var xhr = new XMLHttpRequest();
        xhr.open('GET', url);
        xhr.send();
    }

    window.Holmes = {
        pageView: function(trackingObject) {
            trackingObject['type'] = 'PAGE_VIEW';
            track(trackingObject);
        },
        track: track
    };
}(window));
