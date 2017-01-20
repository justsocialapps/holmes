package assets 

const (
Analyticsjs = `!function(e){function t(t){t.holmesId=e.localStorage.getItem("_holmesId");var o="__HOLMES_BASE_URL__/track?u="+(new Date).getTime()+"&t="+encodeURIComponent(JSON.stringify(t)),n=new XMLHttpRequest;n.open("GET",o),n.send()}var o=e.localStorage.getItem("_holmesId");null===o&&e.localStorage.setItem("_holmesId","__HOLMES_ID__"),e.Holmes={pageView:function(e){e.type="PAGE_VIEW",t(e)},track:t}}(window);`
Bannertxt = `
       ,_
     ,'  `+"`"+`\,_       ██╗  ██╗ ██████╗ ██╗     ███╗   ███╗███████╗███████╗
     |_,-'_)        ██║  ██║██╔═══██╗██║     ████╗ ████║██╔════╝██╔════╝
     /##c '\  (     ███████║██║   ██║██║     ██╔████╔██║█████╗  ███████╗
    ' |'  -{.  )    ██╔══██║██║   ██║██║     ██║╚██╔╝██║██╔══╝  ╚════██║
      /\__-' \[]    ██║  ██║╚██████╔╝███████╗██║ ╚═╝ ██║███████╗███████║
     /`+"`"+`-_`+"`"+`\         ╚═╝  ╚═╝ ╚═════╝ ╚══════╝╚═╝     ╚═╝╚══════╝╚══════╝
     '     \
`
)
