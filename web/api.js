var routes = new Vue({
    el: '#routes',
    data: {
        routeObjects: [
        ]
    },
    methods: {
        clickRoute: function(e){
            routeClicked(e)
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
    console.log("Clicked " + e.currentTarget.innerText)
}

function get_routes() {
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
