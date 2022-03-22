/**
 * Created by GoLand.
 * User: nankai
 * Date: 2019-07-25
 * Time: 17:27
 */

package helper

import (
	"net"
)

// get server ip
// func GetServerIp() (ip string) {

// addrs, err := net.InterfaceAddrs()

// if err != nil {
// return ""
// }

// for _, address := range addrs {
// // Check the ip address to determine whether it is a loopback address
// if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
// if ipNet.IP.To4() != nil {
// ip = ipNet.IP.String()
// }
// }
// }

// return
// }

/**** Problem: I am running a distributed scenario on a local multi-NIC machine. The ip returned by this function is incorrect, causing the rpc connection to fail. Then google results are as follows:
  *** 1. https://www.jianshu.com/p/301aabc06972
  *** 2. https://www.cnblogs.com/chaselogs/p/11301940.html
****/
func GetServerIp() string {
	ip, err := externalIP()
	if err != nil {
		return ""
	}
	return ip.String()
}

func externalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, err
}

func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}
