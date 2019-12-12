package snmpClient

import (
	"github.com/moffa90/triadNMS/constants/snmp"
	"errors"
	g "github.com/soniah/gosnmp"
	"net"
	"strconv"
)

func GetHostnameRemoteSNMP(ip string,port string, community string) (string, error){
	//check if IP is valid
	addressIP := net.ParseIP(ip)
	p, _ := strconv.ParseUint(port, 10, 16)
	if addressIP.To4() != nil {
		g.Default.Target = ip
		g.Default.Community = community
		g.Default.Port = uint16(p)
		g.Default.ExponentialTimeout = false
		err := g.Default.Connect()
		if err != nil {
			return "", err
		}
		defer g.Default.Conn.Close()
		oids := []string{snmp.HostNameOID}
		result, err2 := g.Default.Get(oids) // Get() accepts up to g.MAX_OIDS
		if err2 != nil {
			return  "", err2
		}

		for _, variable := range result.Variables {
			switch variable.Type {
			case g.OctetString:
				return string(variable.Value.([]byte)), nil
				break
			default:
				return "", errors.New("no string")
			}
		}

		return "", errors.New("no results")
	} else {
		return "", errors.New("invalid IP")
	}
}

