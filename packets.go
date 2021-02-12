

package main

import (
    "fmt"
    "log"
    "github.com/google/gopacket/pcap"
    "github.com/google/gopacket/layers"
    "github.com/google/gopacket"
    "time"
    "strings"
    "strconv"
)




var lst_devices [] ObjDevice


func Print_Devices(){
    for n := range lst_devices {
        num := strconv.Itoa(n)

        if len(lst_devices[n].addresses) > 1{
            fmt.Println(num + " ::: " + lst_devices[n].addresses[0].ip.String() +" ::: "+lst_devices[n].addresses[1].ip.String())
        }
       
    }
}

func find_devices(){
    devices, err := pcap.FindAllDevs()
    if err != nil {
        log.Fatal(err)
    }

    // Print device information
    fmt.Println("finding devices")
    for _, device := range devices {
        //fmt.Println("\nName: ", device.Name)
        //fmt.Println("Description: ", device.Description)
        //fmt.Println("Devices addresses: ", device.Description)
        newtmpobjdev := ObjDevice{}
        newtmpobjdev.name = device.Name
        newtmpobjdev.description = device.Description
        newtmpobjdev.devices_addresses = device.Description

        var lst_addresses_tmp [] ObjAddress
        for _, address := range device.Addresses {
            newtmpobjaddy := ObjAddress{}
            newtmpobjaddy.ip = address.IP
            newtmpobjaddy.subnet = address.Netmask
            lst_addresses_tmp = append(lst_addresses_tmp,newtmpobjaddy)
            //fmt.Println("- IP address: ", address.IP)
            //fmt.Println("- Subnet mask: ", address.Netmask)
        }
        newtmpobjdev.addresses = lst_addresses_tmp
        lst_devices = append(lst_devices,newtmpobjdev)
    }
}

func start_live_capture(){
    for n := range lst_devices {
        go live_capture(n,lst_devices[n].name)       
    }
}

func live_capture(num int,dev string){
    var (
        device      string = dev
        snapshotLen int32  = 1024
        promiscuous bool   = false
        err         error
        timeout     time.Duration = 5 * time.Second
        handle      *pcap.Handle
    )

    handle, err = pcap.OpenLive(device, snapshotLen, promiscuous, timeout)
    if err != nil {log.Fatal(err) }
    defer handle.Close()

    packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
    for packet := range packetSource.Packets() {
        //printPacketInfo(packet)
        saveTextPacketInfo(num,packet)
    }
}


func saveTextPacketInfo(num int,packet gopacket.Packet) {
    // Let's see if the packet is an ethernet packet
    ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
    if ethernetLayer != nil {
        
        lst_devices[num].text_log = append(lst_devices[num].text_log,"Ethernet layer detected.")
        ethernetPacket, _ := ethernetLayer.(*layers.Ethernet)
        lst_devices[num].text_log = append(lst_devices[num].text_log,"Source MAC: " + ethernetPacket.SrcMAC.String())
        lst_devices[num].text_log = append(lst_devices[num].text_log,"Destination MAC: " + ethernetPacket.DstMAC.String())
        lst_devices[num].text_log = append(lst_devices[num].text_log,"Ethernet type: " + ethernetPacket.EthernetType.String())
    }

    // Let's see if the packet is IP (even though the ether type told us)
    ipLayer := packet.Layer(layers.LayerTypeIPv4)
    if ipLayer != nil {
        lst_devices[num].text_log = append(lst_devices[num].text_log,"IPv4 layer detected.")
        ip, _ := ipLayer.(*layers.IPv4)
        lst_devices[num].text_log = append(lst_devices[num].text_log,"From" + ip.SrcIP.String() + " to " +  ip.DstIP.String())
        lst_devices[num].text_log = append(lst_devices[num].text_log,"Protocol: " + ip.Protocol.String())
    }

    // Let's see if the packet is TCP
    tcpLayer := packet.Layer(layers.LayerTypeTCP)
    if tcpLayer != nil {
        lst_devices[num].text_log = append(lst_devices[num].text_log,"TCP layer detected.")
        tcp, _ := tcpLayer.(*layers.TCP)
        lst_devices[num].text_log = append(lst_devices[num].text_log,"From port " + tcp.SrcPort.String() + " to " + tcp.DstPort.String())
        tcpseq := strconv.Itoa(int(tcp.Seq))
        lst_devices[num].text_log = append(lst_devices[num].text_log,"Sequence number: " + tcpseq)
    }

    // Iterate over all layers, printing out each layer type
    lst_devices[num].text_log = append(lst_devices[num].text_log,"All packet layers:")
    for _, layer := range packet.Layers() {
        lst_devices[num].text_log = append(lst_devices[num].text_log, "- " + layer.LayerType().String())
    }

    // When iterating through packet.Layers() above,
    // if it lists Payload layer then that is the same as
    // this applicationLayer. applicationLayer contains the payload
    applicationLayer := packet.ApplicationLayer()
    if applicationLayer != nil {
        lst_devices[num].text_log = append(lst_devices[num].text_log,"Application layer/Payload found.")
        paydata := string(applicationLayer.Payload())
        lst_devices[num].text_log = append(lst_devices[num].text_log,paydata)
        if strings.Contains(string(applicationLayer.Payload()), "HTTP") {
            //fmt.Println("HTTP found!")
        }
    }

    // Check for errors
    if err := packet.ErrorLayer(); err != nil {
        fmt.Println("Error decoding some part of the packet:", err)
    }
}

func printPacketInfo(packet gopacket.Packet) {
    // Let's see if the packet is an ethernet packet
    ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
    if ethernetLayer != nil {
        fmt.Println("Ethernet layer detected.")
        ethernetPacket, _ := ethernetLayer.(*layers.Ethernet)
        fmt.Println("Source MAC: ", ethernetPacket.SrcMAC)
        fmt.Println("Destination MAC: ", ethernetPacket.DstMAC)
        // Ethernet type is typically IPv4 but could be ARP or other
        fmt.Println("Ethernet type: ", ethernetPacket.EthernetType)
        fmt.Println()
    }

    // Let's see if the packet is IP (even though the ether type told us)
    ipLayer := packet.Layer(layers.LayerTypeIPv4)
    if ipLayer != nil {
        fmt.Println("IPv4 layer detected.")
        ip, _ := ipLayer.(*layers.IPv4)

        // IP layer variables:
        // Version (Either 4 or 6)
        // IHL (IP Header Length in 32-bit words)
        // TOS, Length, Id, Flags, FragOffset, TTL, Protocol (TCP?),
        // Checksum, SrcIP, DstIP
        fmt.Printf("From %s to %s\n", ip.SrcIP, ip.DstIP)
        fmt.Println("Protocol: ", ip.Protocol)
        fmt.Println()
    }

    // Let's see if the packet is TCP
    tcpLayer := packet.Layer(layers.LayerTypeTCP)
    if tcpLayer != nil {
        fmt.Println("TCP layer detected.")
        tcp, _ := tcpLayer.(*layers.TCP)

        // TCP layer variables:
        // SrcPort, DstPort, Seq, Ack, DataOffset, Window, Checksum, Urgent
        // Bool flags: FIN, SYN, RST, PSH, ACK, URG, ECE, CWR, NS
        fmt.Printf("From port %d to %d\n", tcp.SrcPort, tcp.DstPort)
        fmt.Println("Sequence number: ", tcp.Seq)
        fmt.Println()
    }

    // Iterate over all layers, printing out each layer type
    fmt.Println("All packet layers:")
    for _, layer := range packet.Layers() {
        fmt.Println("- ", layer.LayerType())
    }

    // When iterating through packet.Layers() above,
    // if it lists Payload layer then that is the same as
    // this applicationLayer. applicationLayer contains the payload
    applicationLayer := packet.ApplicationLayer()
    if applicationLayer != nil {
        fmt.Println("Application layer/Payload found.")
        fmt.Printf("%s\n", applicationLayer.Payload())

        // Search for a string inside the payload
        if strings.Contains(string(applicationLayer.Payload()), "HTTP") {
            fmt.Println("HTTP found!")
        }
    }

    // Check for errors
    if err := packet.ErrorLayer(); err != nil {
        fmt.Println("Error decoding some part of the packet:", err)
    }
}