var routes = new Vue({
    el: '#routes',
    data: {
      routeObjects: [
        {in:"hej", out:"abc"},
        {in:"hopp", out:"abc2"},
        {in:"knasen3", out:"abc2"}
      ]
    },
    methods: {
        clickRoute: function(e){
            routeClicked(e)
        }
    }
  })

function routeClicked(e) {
    console.log("Clicked " + e.currentTarget.innerText)
}

function get_routes(el) {
    let routes_list="<h4>kalas</h4>";
    //document.getElementById(el).innerHTML = routes_list;
    fetch('./api/routes')
    .then(function(response) {
      //document.getElementById(el).innerHTML = "No routes";
      return response.json();
    })
    .then(function(myJson) {
      console.log("myJson")
      routes.routeObjects = [{in:"All"}]
      var ll = myJson.length;
      while ( ll-- ) {
        var ob = myJson[ll]
        routes.routeObjects.push(ob)
        
      }
    });
}

function get_log(el) {
    fetch('./api/log')
    .then(function(response) {
      document.getElementById(el).innerHTML = "RRR!";
      console.log("response log A " , response.status);

      return response.json();
    })
    .then(function(myJson) {
      var ll = myJson.length;
      var list = '<ul class="list-group">';
      while ( ll-- ) {
        var ob = myJson[ll]
        list = list + "<li class='list-group-item'><div class='d-inline mr-2'><b>" + myJson[ll].Time + ":</b></div>";
        list = list + "<div class='d-inline mr-2'>"+ myJson[ll].Input + "</div><div class=d-inline>";
        if (myJson[ll].Output.length > 0) {
          if (myJson[ll].OutProtocol.length  > 0) {
            list = list + myJson[ll].OutProtocol + ":" + myJson[ll].Output;
          } else {
            list = myJson[ll].Output;
          }
        }
        list = list + "</div></li>"
      }
      list = list + "</ol>";

      document.getElementById(el).innerHTML = list;

    });

}
