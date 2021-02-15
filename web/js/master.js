var lookingat = ""
var start = 0;
var finish = -1;
var current = -1;

$( document ).ready(function() {

	console.log("Welcome to paws");
  //LoadFile("/web/menu.html");
  LoadFile("/devices","lstDevices");
  //lstDevices

  $("a").click(function(event) {
      var href = $(this).attr('href');
      console.log(href);
      ChangePage(href);
      event.preventDefault();
  });
  StartLoop();
});


function LoadFile(file,id){
  $.get(file, function(data){
    //console.log(data);
      //$(this).children("div:first").html(data);
      $("#"+id).html(data);
  });
}

function Select_Device(dev){
  
  
  $("#SelectedDevice").html(dev);
  var res = dev.split("/"); 
  console.log("Getting log: " + res[2]);
  lookingat = res[2];
  SetupFinish();
}

function ChangePage(pages){
  if(pages == "#Log"){
    $("#CodeShow").css({"display": "block"});
  }else if(pages == "#Home"){
    $("#CodeShow").css({"display": "none"});
  }

}
//---------------------------------------------------log management
function StartLoop(){
  var tid = setInterval(LoopDeLoop, 2000);
  

}

function LoopDeLoop() {
  console.log("CHECK");

  if(lookingat != "" && finish != -1){
    //console.log("about to req data");
    console.log("start: " + start + " finish: " + finish);
    console.log(" current: " + current);
    if(current != finish){
      console.log("req data");
      var req = "/dev?dev=" + lookingat + "&start=" + start + "&finish=" + finish
      $.get(req, function(data){
        //finish = data;
        $("#CodeShow").append(data);
        console.log("append data: " + req);
        
      });
     // current = finish;
      start = finish;
    }
  }

  if(start == finish){
    console.log("update between " + finish + " and " + current);
    var req = "/dev?dev=" + lookingat + "&start=" + finish + "&finish=" + current
    $.get(req, function(data){
      //finish = data;
      $("#CodeShow").append(data);
      //console.log("append data: " + req);
      finish = current;
      start = current;
    });
  }

  if (lookingat != ""){
    $.get("/dev?dev=" + lookingat, function(data){
      current = data;
    });
  }
}
function StopLoop() { // to be called when you want to stop the timer
  clearInterval(tid);
}

//-- requests

function SetupFinish(){
  $.get("/dev?dev=" + lookingat, function(data){
    finish = data;
  });
}