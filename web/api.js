var routes = new Vue({
    el: '#routes',
    data: {
      showCfg:false,
      lastClick:0,
      routeObjects: []
    },
    methods: {
        clickRoute: function(e, inp){
          setFocus(e.explicitOriginalTarget);
          routeClicked(e, inp)
        },
        clickRouteCfg: function(e){
          setFocus(e.explicitOriginalTarget);
          this.showCfg=!this.showCfg
        },
        cancelConf: function(e){
          this.showCfg=false
        }
    }
})

var log = new Vue({
    el: '#log',
    data: {
      logObjects: [
      ]
    }
})

function setFocus(el) {
  if(this.lastClick) {
    this.lastClick.classList.add("btn-info")
    this.lastClick.classList.remove("btn-warning")
  }
  let e = el;
  if (el.tagName.toUpperCase() != "BUTTON") {
    // Clicked image inside button.
    e = el.parentElement
  }
  this.lastClick = e
  e.classList.remove("btn-info")
  e.classList.add("btn-warning")
}

function routeClicked(e, inp) {
  console.log("Get RouteDef for " + inp)
  get_log(inp)
}

function get_routeDefs() {
    fetch('./api/routes')
    .then(function(response) {
      return response.json();
    })
    .then(function(myJson) {
      routes.routeObjects = [{in:""}]
      let ll = myJson.length;
      while ( ll-- ) {
        var ob = myJson[ll]
        routes.routeObjects.push(ob)
      }
    });
}

function get_log(routeDef) {
  let path = '/api/log'
  if (routeDef) {
    path += '?in='+routeDef
  }
  fetch(path)
  .then(function(response) {
    return response.json();
  })
  .then(function(myJson) {
    log.logObjects = []
    let ll = myJson.length;
    while ( ll-- ) {
      let ob = myJson[ll]
      if (ob.OutProtocol && ob.OutProtocol.length  > 0) {
          ob.OutProtocol += ":"
      }
      log.logObjects.push(ob)
    }
  });

}
