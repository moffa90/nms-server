package home

import (
	"github.com/moffa90/triadNMS/constants"
	"github.com/moffa90/triadNMS/db"
	"github.com/moffa90/triadNMS/db/models"
	"github.com/moffa90/triadNMS/utils"
	"github.com/moffa90/triadNMS/utils/security"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/shirou/gopsutil/host"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func Handler(w http.ResponseWriter, req *http.Request) {
	HostInfo, _ := host.Info()
	mode := os.Getenv("mode")

	type Slot struct {
		Number int
		Hardware []models.Hardware
	}
	type Backplane struct {
		Number	int
		Slots	[]Slot
	}

	//TODO: make backplane number configurable, default 2

	var backplaneQty, slotQty, hardwareQty int

	switch mode {
	case "hec-1":
		backplaneQty = 1
		slotQty = 5
		hardwareQty = 2

		break

	case "hec-2":
		backplaneQty = 2
		slotQty = 5
		hardwareQty = 2

		break

	case "remote":
		backplaneQty = 1
		slotQty = 4
		hardwareQty = 1
		break
	}

	backplanes := make([]Backplane, backplaneQty)
	for i, b := range backplanes{
		b.Slots = make([]Slot, slotQty)
		b.Number = i

		for j, s := range b.Slots {
			s.Number = j
			s.Hardware = make([]models.Hardware, hardwareQty)
			b.Slots[j] = s
		}
		backplanes[i] = b
	}

	hardwares, _ := models.GetHardware(db.Shared)

	for _, h := range hardwares{
		//backplanes[h.Backplane].Slots[0].Hardware = append([]models.Hardware{h}, backplanes[h.Backplane].Slots[0].Hardware...)
		if h.Backplane < backplaneQty && h.Address <= slotQty  && h.Address > 0{
			backplanes[h.Backplane].Slots[h.Address -1].Hardware = append([]models.Hardware{h}, backplanes[h.Backplane].Slots[h.Address - 1].Hardware...)
		}else{
			//invalid address
			log.Errorf("Invalid address for: %s", h.ToString())
		}
	}

	infoStruct := struct{
		Host host.InfoStat
		Mender string
		Artifact string
		Active string
		Info map[string]string
		Backplanes []Backplane
		Remotes [][]models.Remote
		Mode string
	}{
		*HostInfo,
		getMenderVersion(),
		getMenderArtifact(),
		"index",
		security.CookieGetInfo(w, req),
		backplanes,
		nil,
		mode,
	}

	remotes := models.GetRemotes(db.Shared)
	infoStruct.Remotes = make([][]models.Remote, 4)

	for i, _ := range infoStruct.Remotes{
		infoStruct.Remotes[i] = make([]models.Remote, 5)
	}

	for _, r := range remotes{
		infoStruct.Remotes[r.Remote][r.Group] = r
	}

	utils.RenderPage(w, constants.TEMPLATE_PAGE_HOME_PATH, infoStruct)
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