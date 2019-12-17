var routes = new Vue({
    el: '#routes',
    data: {
      showCfg:false,
      thisIn:"hhhh",
      routeObjects: []
    },
    methods: {
        clickRoute: function(e){
          console.log("Click rout")
          routeClicked(e)
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

function routeClicked(e) {
  this.thisIn="Abba"

  route64=btoa(e.currentTarget.innerText)
  console.log("Get RouteDef for " + e.currentTarget.innerText + " route64:" + route64)
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

function get_log() {
    fetch('./api/log')
    .then(function(response) {
      return response.json();
    })
    .then(function(myJson) {
      // console.dir(myJson)
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
