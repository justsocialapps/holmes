package assets

var (
	Analyticsjs = `!function(e){function t(t){t.holmesId=e.localStorage.getItem("_holmesId");for(var n=0;n<o.length;n++)o[n](t);var a="__HOLMES_BASE_URL__/track?u="+(new Date).getTime()+"&t="+encodeURIComponent(JSON.stringify(t)),l=new XMLHttpRequest;l.open("GET",a),l.send()}var n,o=[];null===e.localStorage.getItem("_holmesId")&&e.localStorage.setItem("_holmesId","__HOLMES_ID__"),e.Holmes={pageView:function(e){e.type="PAGE_VIEW",t(e)},addTrackingEnricher:function(e){o.push(e)},track:t},(n=document.createEvent("Event")).initEvent("holmesloaded",!1,!1),e.dispatchEvent(n)}(window);`
	Bannertxt   = `
       ,_
     ,'  ` + "`" + `\,_       ██╗  ██╗ ██████╗ ██╗     ███╗   ███╗███████╗███████╗
     |_,-'_)        ██║  ██║██╔═══██╗██║     ████╗ ████║██╔════╝██╔════╝
     /##c '\  (     ███████║██║   ██║██║     ██╔████╔██║█████╗  ███████╗
    ' |'  -{.  )    ██╔══██║██║   ██║██║     ██║╚██╔╝██║██╔══╝  ╚════██║
      /\__-' \[]    ██║  ██║╚██████╔╝███████╗██║ ╚═╝ ██║███████╗███████║
     /` + "`" + `-_` + "`" + `\         ╚═╝  ╚═╝ ╚═════╝ ╚══════╝╚═╝     ╚═╝╚══════╝╚══════╝
     '     \
`
)
