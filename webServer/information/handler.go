package information

import (
	"github.com/moffa90/nms-server/constants"
	"github.com/moffa90/nms-server/utils"
	"github.com/moffa90/nms-server/utils/security"
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"net/http"
	"os/exec"
	"strings"
)

func Handler(w http.ResponseWriter, req *http.Request) {
	HostInfo, _ := host.Info()
	Cores, _ := cpu.Counts(true)
	Mem, _ := mem.VirtualMemory()
	Net, _ := net.Interfaces()
	usagePercent, _ := cpu.Percent(0, false)

	type CPUStruct struct{
		Cores int
		Usage float64
	}

	CPUInfo := CPUStruct{
		Cores,
		usagePercent[0],
	}

	infoStruct := struct{
		Host host.InfoStat
		CPU CPUStruct
		Mem mem.VirtualMemoryStat
		Net []net.InterfaceStat
		Mender string
		Artifact string
		Active string
		Info map[string]string
	}{
		*HostInfo,
		CPUInfo,
		*Mem,
		Net,
		getMenderVersion(),
		getMenderArtifact(),
		"info",
		security.CookieGetInfo(w, req),
	}

	utils.RenderPage(w, constants.TEMPLATE_PAGE_INFO_PATH, infoStruct)


}

func getMenderVersion() string{
	cmd := exec.Command("mender", "-version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("cmd.Run() failed with %s\n", err)
	}

	s := string(out)
	subStr:= "module=main"
	fields := strings.LastIndex(s, subStr)
	return fmt.Sprintf("%s", s[fields+len(subStr)+1:])
}

func getMenderArtifact() string{
	cmd := exec.Command("mender", "-show-artifact")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("cmd.Run() failed with %s\n", err)
	}

	s := string(out)
	subStr:= "module=main"
	fields := strings.LastIndex(s, subStr)
	return fmt.Sprintf("%s", s[fields+len(subStr)+1:])
}