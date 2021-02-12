$( document ).ready(function() {

	console.log("Welcome to ggg");
  //LoadFile("/web/menu.html");
  LoadFile("/devices","lstDevices");
  //lstDevices

  $("a").click(function(event) {
      var href = $(this).attr('href');
      console.log(href);
      ChangePage(href);
      event.preventDefault();
  });

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
  
  $("#SelectedDevice").html(dev);
}

function ChangePage(pages){

}