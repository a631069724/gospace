package catchpack

import (
	"fmt"

	"github.com/axgle/mahonia"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func CatchPack() {

	fmt.Println("packet start...")
	ifs, err := pcap.FindAllDevs()
	if err != nil {
		fmt.Println("find all devs err:", err)
	}
	for _, infa := range ifs {
		fmt.Println(infa.Name, infa.Flags, infa.Addresses)
	}
	deviceName := "{24EB13F6-A71F-41ED-901C-72AAE6342FD9}"
	snapLen := int32(65535)
	port := uint16(443)
	filter := getFilter(52140, port)
	fmt.Printf("device:%v, snapLen:%v, port:%v\n", deviceName, snapLen, port)
	fmt.Println("filter:", filter)

	handle, err := pcap.OpenLive(deviceName, snapLen, true, pcap.BlockForever)
	if err != nil {
		fmt.Printf("pcap open live failed:%s\n", mahonia.NewDecoder("gbk").ConvertString(err.Error()))

		return
	}

	if err := handle.SetBPFFilter(filter); err != nil {
		fmt.Printf("set bpf filter failed: %v", err)
		return
	}
	defer handle.Close()
	/*
		var tcp layers.TCP
		parser := gopacket.NewDecodingLayerParser(layers.LayerTypeTCP, &tcp)
		decodedLayers := []gopacket.LayerType{}
	*/
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packetSource.NoCopy = true
	for packet := range packetSource.Packets() {
		if packet.NetworkLayer() == nil || packet.TransportLayer() == nil || packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
			fmt.Println("unexpected packet")
			continue
		}
		tcpLayer := packet.Layer(layers.LayerTypeTCP)
		if tcpLayer != nil {
			tcp, _ := tcpLayer.(*layers.TCP)
			tcp.
				fmt.Println(string(tcp.Payload))

		}
	}
}

func getFilter(lport uint16, port uint16) string {
	filter := fmt.Sprintf("tcp and ((src port %v) or (dst port %v))", lport, port)
	return filter
}
