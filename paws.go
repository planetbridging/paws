package main
//paws	= packets are web stuff
import(

)


func main(){

	find_devices()
	
	Print_Devices()
	go start_live_capture()

	start_hosting()
	//device = lst_devices[3].name
	//live_capture()
}
