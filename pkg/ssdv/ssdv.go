package ssdv

// #cgo CFLAGS: -g -Wall
// #include <stdlib.h>
// #include "ssdv.h"
// #include "rs8.h"
// #include "decode.h"
import "C"
import (
	"fmt"
	"os"
	"path/filepath"
	"unsafe"
)

const (
	SSDV_HEADER_IMAGE      = 5
	SSDV_HEADER_PACKET_MSB = 6
	SSDV_HEADER_PACKET_LSB = 7
	SSDV_HEADER_FLAGS      = 10
)

type SSDVBasicInfo struct {
	Packet     uint16
	LastPacket bool
	Image      uint8
}

func SSDVPacketInfo(packet []byte) SSDVBasicInfo {
	info := SSDVBasicInfo{}
	info.Image = packet[SSDV_HEADER_IMAGE]
	info.Packet = (uint16(packet[SSDV_HEADER_PACKET_MSB]) << 8) + uint16(packet[SSDV_HEADER_PACKET_LSB])
	if (packet[SSDV_HEADER_FLAGS]&0b00000100)>>2 != 0 {
		info.LastPacket = true
	} else {
		info.LastPacket = false
	}
	return info
}

func SSDVDecodePacket(packet []byte, path string) (string, string, error) {
	info := SSDVPacketInfo(packet)

	// files
	input := filepath.Join(path, fmt.Sprintf("ssdv%d.bin", info.Image))
	output := filepath.Join(path, fmt.Sprintf("ssdv%d.jpg", info.Image))

	// save SSDV bin data
	var file *os.File
	var err error
	// if first packet open new file
	if info.Packet == 0 {
		file, err = os.Create(input)
		if err != nil {
			return "", "", err
		}
	} else {
		// append
		file, err = os.OpenFile(input, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return "", "", err
		}
	}

	// add data
	file.Seek(int64(info.Packet)*256, 0)
	file.Write([]byte{0x55}) // add removed sync byte
	file.Write(packet)
	file.Close()

	// decode
	in := C.CString(input)
	out := C.CString(output)
	mission := C.CString("       ") // 7 chars
	defer C.free(unsafe.Pointer(in))
	defer C.free(unsafe.Pointer(out))
	defer C.free(unsafe.Pointer(mission))
	data := C.decode_ssdv_file(in, out, mission)
	fmt.Printf("Paquetes %d, Imagen %d\n", data.decoded_packets, data.image_id)
	return output, C.GoString(mission), nil
}
