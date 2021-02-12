package main

import(
	"net"
)
//device objs

type ObjAddress struct{
	ip net.IP 
	subnet net.IPMask
}

type ObjDevice struct{
	name string
	description string
	devices_addresses string
	addresses [] ObjAddress
	text_log [] string
}