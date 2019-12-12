package subagentSNMP

import (
	"github.com/moffa90/nms-server/assets"
	"github.com/moffa90/nms-server/constants"
	"github.com/moffa90/nms-server/constants/snmp"
	"github.com/moffa90/nms-server/snmpTree"
	"github.com/moffa90/nms-server/utils"
	log "github.com/Sirupsen/logrus"
	"gitlab.com/martinclaro/go-agentx"
	"gitlab.com/martinclaro/go-agentx/pdu"
	"gitlab.com/martinclaro/go-agentx/value"
	"gopkg.in/errgo.v1"
	"os"
	"os/exec"
	"strconv"
	"text/template"
	"time"
)

func StartSubagent() {

	if !utils.ExistsFile(constants.LOCAL_SNMP_CONF_DIRECTORY){
		if e := os.MkdirAll(constants.LOCAL_SNMP_CONF_DIRECTORY, os.ModePerm); e != nil {
			log.Info(e.Error())
		}
	}

	if !utils.ExistsFile(constants.LOCAL_SNMPD_CONF_PATH) {
		snmpd, _ := assets.Asset(constants.TEMPLATE_SNMPD_CONF_PATH)

		f, errFile := utils.CreateFile(constants.LOCAL_SNMPD_CONF_PATH)

		if errFile != nil {
			log.Info("Error creating File:" + errFile.Error())
		}

		tpl := template.New("temp")
		tpl.New("snmpd").Parse(string(snmpd))

		snmpdInfo := struct {
			WriteString string
			ReadString string
			Port uint64
			SysContact string
			SysLocation string
		}{
			"private",
			"public",
			161,
			"TBD",
			"TBD",
		}

		if ErrTpl := tpl.ExecuteTemplate(f, "snmpd", snmpdInfo); ErrTpl != nil {
			log.Info(ErrTpl.Error())
		}

		f.Close()

		cmd := exec.Command("systemctl", "restart", "snmpd")
		out, e := cmd.CombinedOutput()
		if e !=nil {
			log.Info(out)
			log.Info(e.Error())
		}

		time.Sleep(time.Second)
	}

	client := &agentx.Client{
		Net:               "tcp",
		Address:           "localhost:705",
		Timeout:           1 * time.Minute,
		ReconnectInterval: 1 * time.Second,
	}

	if err := client.Open(); err != nil {
		log.Fatalf(errgo.Details(err))
	}

	session, err := client.Session()
	if err != nil {
		log.Fatalf(errgo.Details(err))
	}

	listHandler := &agentx.ListHandler{}

	tree := snmpTree.BuildNMSTree()
	snmpTree.RegisterTree(tree, func(node *snmpTree.Node) {

		log.Info(node.Oid())
		item := listHandler.Add(node.Oid())
		if node.Class() == "valueInt"{
			item.Type = pdu.VariableTypeInteger
			item.Value = func() int32 {
				n, _ := strconv.ParseInt(node.Value(), 10, 4)
				return int32(n)
			}
		}else if node.Class() == "valueStr"{
			item.Type = pdu.VariableTypeOctetString
			item.Value = func() string {
				return  node.Value()
			}
		}

	})

	session.Handler = listHandler

	if err := session.Register(127, value.MustParseOID(snmp.NMSTree)); err != nil {
		log.Fatalf(errgo.Details(err))
	}

	for {
		time.Sleep(100 * time.Millisecond)
	}
}

