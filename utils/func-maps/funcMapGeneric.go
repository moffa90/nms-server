package func_maps

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"html/template"
	"strings"
	"time"
)

var Generic template.FuncMap = template.FuncMap{
	"ToUpper": strings.ToUpper,
	"Capitalize": strings.Title,
	"EpochToDate": EpochToDate,
	"SecondsToMinutes": SecondsToMinutes,
	"FormatPercent" : FormatPercent,
	"BytesToHuman" : byteToHuman,
	"TimeToDate": TimeToDate,
	"ToHex" : IntToHex,
	"GetSensorsTitle": GetSensorsTitle,
	"RemoveSpaces": RemoveSpaces,
}

func EpochToDate(e uint64) string{
	return time.Unix(int64(e), 0).Format(time.RFC822Z)
}

func TimeToDate(t time.Time) string{
	if t.IsZero(){
		return  "Never"
	}
	return t.Format("Mon _2 Jan 2006 15:04")
}

func SecondsToMinutes(sec uint64) string{
	var substract int64 = -1 * int64(sec)
	return humanize.Time(time.Now().Add(time.Second * time.Duration(substract)))
}

func FormatPercent(percentage float64) string{
	return fmt.Sprintf("%.2f%s", percentage, "%")
}

func byteToHuman(b uint64) string{
	return humanize.Bytes(b)
}

func IntToHex(i int) string{
	return fmt.Sprintf("%#02x",i)
}

func GetSensorsTitle(s string, values map[string]string) string{
	switch s {
	case "OpticalPower":
		return fmt.Sprintf("%s (min:%s, max:%s)",s, values["Optical Power Alarm Lo"], values["Optical Power Alarm Hi"])

	case "RFPower":
		return fmt.Sprintf("%s (min:%s, max:%s)",s, values["RF Detector Alarm Lo"], values["RF Detector Alarm Hi"])

	case "Temperature":
		return fmt.Sprintf("%s (min:%s, max:%s)",s, values["Temperature Alarm Lo"], values["Temperature Alarm Hi"])
	default:
		return ""
	}
}

func RemoveSpaces(s string) string{
	return strings.Replace(s, " ", "", -1)
}

