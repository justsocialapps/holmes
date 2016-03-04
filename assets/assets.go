package assets 

const (
Analyticsjs = `!function(e){function t(t){t.holmesId=e.localStorage.getItem("_holmesId");var n="__HOLMES_BASE_URL__/track?u="+(new Date).getTime()+"&t="+encodeURIComponent(JSON.stringify(t));o.open("GET",n),o.send()}var o=new XMLHttpRequest,n=e.localStorage.getItem("_holmesId");null===n&&e.localStorage.setItem("_holmesId","__HOLMES_ID__"),e.Holmes={pageView:function(e){e.type="PAGE_VIEW",t(e)}}}(window);`
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
