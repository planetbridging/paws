$( document ).ready(function() {

	console.log("Welcome to ggg");
  //LoadFile("/web/menu.html");
  console.log(LoadFile("/devices","lstDevices"));
  //lstDevices
});


function LoadFile(file,id){
  $.get(file, function(data){
    //console.log(data);
      //$(this).children("div:first").html(data);
      $("#"+id).html(data);
  });
}

function Select_Device(dev){
  console.log(dev);
}