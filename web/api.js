var routes = new Vue({
    el: '#routes',
    data: {
      showCfg:false,
      thisIn:"hhhh",
      routeObjects: []
    },
    methods: {
        clickRoute: function(e, inp){
          console.log("Click route!", inp)
          routeClicked(e, inp)
        },
        clickRouteCfg: function(e){
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

function routeClicked(e, inp) {
  this.thisIn=inp

  // route64=btoa(e.currentTarget.innerText)
  console.log("Get RouteDef for " + inp)
  get_log(inp)
    // for (let ix in routes.routeObjects) {
    //   let ro = routes.routeObjects[ix]
    //   console.log("RO: " , ro.in , ro.out)
    //   console.dir(ro)

    // }

}

function get_routeDefs() {
    fetch('./api/routes')
    .then(function(response) {
      return response.json();
    })
    .then(function(myJson) {
      routes.routeObjects = [{in:"All"}]
      let ll = myJson.length;
      while ( ll-- ) {
        var ob = myJson[ll]
        routes.routeObjects.push(ob)
      }
    });
}

function get_log(routeDef) {
  fetch('./api/log?in='+routeDef)
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
