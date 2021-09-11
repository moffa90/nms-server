package usb

import (
	"bufio"
	"github.com/moffa90/triadNMS/constants"
	"github.com/moffa90/triadNMS/db"
	"github.com/moffa90/triadNMS/db/models"
	log "github.com/sirupsen/logrus"
	"github.com/google/gousb"
	"github.com/google/gousb/usbid"
	"github.com/joho/godotenv"
	serial2 "github.com/tarm/serial"
	"strconv"
	"strings"
	"sync"
)

type USBDevice struct {
	VendorID  int
	ProductID int
	HumanDesc string
	Bus       int
	Address   int
	Serial    string
}

var mutexes map[string]*sync.Mutex
var Shared *gousb.Context

func Init() {
	Shared = gousb.NewContext()
	mutexes = make(map[string]*sync.Mutex)
}

func lock(serial string) {
	if m, ok := mutexes[serial]; ok {
		m.Lock()
	} else {
		mutexes[serial] = &sync.Mutex{}
		mutexes[serial].Lock()
	}
}

func unlock(serial string) {
	if m, ok := mutexes[serial]; ok {
		m.Unlock()
	}
}

func GetDevices() []USBDevice {
	// Only one context should be needed for an application.  It should always be closed.
	usbInfo := make([]USBDevice, 0, 1)
	// OpenDevices is used to find the devices to open.
	devs, err := Shared.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		// The usbid package can be used to print out human readable information.
		if desc.Vendor == constants.CYPRESS_VENDOR_ID {
			return true
		}
		return false
	})

	for _, dev := range devs {
		serial, _ := dev.SerialNumber()
		if d, _ := models.GetHardwareBySerial(db.Shared, serial); d == nil {
			usbInfo = append(usbInfo, USBDevice{
				int(dev.Desc.Vendor),
				int(dev.Desc.Product),
				usbid.Describe(dev.Desc),
				dev.Desc.Bus,
				dev.Desc.Address,
				serial,
			})
		}

		dev.Close()
	}
	// OpenDevices can occaionally fail, so be sure to check its return value.
	if err != nil {
		log.Printf("list: %s\n\r", err.Error())
	}

	return usbInfo
}

func openDevice(s string) *serial2.Port {

	s = "/dev/ttyUSB_" + s
	c := &serial2.Config{Name: s, Baud: 115200, ReadTimeout: 500}
	port, err := serial2.OpenPort(c)

	if err != nil {
		log.Errorf("Error opening port %s: %s", s, err.Error())
	}

	return port
}

func strToMapValues(s string) map[string]string {
	i := strings.Index(s, "\r\n")
	if parsed, _ := godotenv.Unmarshal(s[i+1:]); parsed != nil && len(parsed) > 0 {
		return parsed
	} else {
		return nil
	}
}

func localSendCommand(command string, port *serial2.Port) string{
	port.Flush()

	bytes := []byte(command)
	for _, b := range bytes {
		_, err := port.Write([]byte{b})
		if err != nil {
			log.Error(err)
		}
	}

	if command == constants.RESET_COMMAND{
		return ""
	}

	buff := ""

	scanner := bufio.NewScanner(port)
	for scanner.Scan() {
		buff += scanner.Text() + "\r\n"
	}
	if scanner.Err() != nil {
		log.Fatal(scanner.Err())
	}

	if strings.HasPrefix(buff, strings.TrimSuffix(command, "\n")){
		return buff
	}

	return ""
}

func SendCommand(serial, command string) string {
	lock(serial)
	defer unlock(serial)
	port := openDevice(serial)
	// Make sure to close it later.
	if port == nil {
		return ""
	}
	defer port.Close()

	return localSendCommand(command, port)
}

func ExecCommand(s string, command string) (string, map[string]string) {
	rawResponse := SendCommand(s, command)
	if parsed := strToMapValues(rawResponse); parsed != nil {
		return rawResponse, parsed
	}

	return rawResponse, nil
}

func GetDeviceInfo(s string) map[string]string {

	if parsed := strToMapValues(SendCommand(s, constants.INFO_COMMAND)); parsed != nil {
		if val, ok := parsed["AddressID"]; ok {
			address, _ := strconv.ParseInt(strings.Replace(val, "x", "", -1), 16, 32)
			parsed["Backplane"] = strconv.Itoa(int(address >> 3))
			parsed["Slot"] = strconv.Itoa(int(address & 0x07))
			parsed["Serial"] = s
		}
		return parsed
	}

	return nil
}

func GetFullData(s string) map[string]map[string]string{
	lock(s)

	port := openDevice(s)

	defer unlock(s)
	// Make sure to close it later.
	if port == nil {
		return nil
	}
	defer port.Close()

	data := make(map[string]map[string]string)

	if data["info"] = strToMapValues(localSendCommand(constants.INFO_COMMAND, port)); data["info"] != nil {
		if val, ok := data["info"]["AddressID"]; ok {
			address, _ := strconv.ParseInt(strings.Replace(val, "x", "", -1), 16, 32)
			data["info"]["Backplane"] = strconv.Itoa(int(address >> 3))
			data["info"]["Slot"] = strconv.Itoa(int(address & 0x07))
		}
	}

	data["sensors"] = strToMapValues(localSendCommand(constants.READ_SENSORS_COMMAND, port))

	data["alarms"] = strToMapValues(localSendCommand(constants.READ_ALARMS_COMMAND, port))

	data["controls"] = strToMapValues(localSendCommand(constants.READ_CONTROLS_COMMAND, port))

	data["thresholds"] = strToMapValues(localSendCommand(constants.READ_ALARMS_THRESHOLD_COMMAND, port))

	data["RFCalDet"] = strToMapValues(localSendCommand(constants.READ_CAL_RF_DET_COMMAND, port))

	data["ADC"] = strToMapValues(localSendCommand(constants.READ_ADC_COMMAND, port))

	switch data["info"]["ProductID"] {
	case constants.FOTX_200_REV1:
		data["extra1"] = strToMapValues(localSendCommand(constants.READ_LASER_APC_COMMAND, port))
		data["extra2"] = strToMapValues(localSendCommand(constants.READ_CAL_OPT_PWR_COMMAND, port))
		break
	case constants.MOD_100_REV1:
		data["extra1"] = strToMapValues(localSendCommand(constants.READ_CAL_RF_DET2_COMMAND, port))
		data["extra2"] = strToMapValues(localSendCommand(constants.READ_RF_AGC_COMMAND, port))
		break
	case constants.FORX_200_REV1:
		data["extra1"] = strToMapValues(localSendCommand(constants.READ_CAL_OPT_PWR_COMMAND, port))
		data["extra1"]["NameExtra"] = "Optical Power"
		break
	case constants.MOD_200_REV1:

		break
	}

	return data
}

func GetSensorsData(s string) map[string]map[string]string{
	lock(s)

	port := openDevice(s)

	defer unlock(s)
	// Make sure to close it later.
	if port == nil {
		return nil
	}
	defer port.Close()

	data := make(map[string]map[string]string)

	data["sensors"] = strToMapValues(localSendCommand(constants.READ_SENSORS_COMMAND, port))

	data["alarms"] = strToMapValues(localSendCommand(constants.READ_ALARMS_COMMAND, port))

	data["ADC"] = strToMapValues(localSendCommand(constants.READ_ADC_COMMAND, port))

	return data
}

func GetInfoData(s string) map[string]map[string]string{
	lock(s)

	port := openDevice(s)

	defer unlock(s)
	// Make sure to close it later.
	if port == nil {
		return nil
	}
	defer port.Close()

	data := make(map[string]map[string]string)

	if data["info"] = strToMapValues(localSendCommand(constants.INFO_COMMAND, port)); data["info"] != nil {
		if val, ok := data["info"]["AddressID"]; ok {
			address, _ := strconv.ParseInt(strings.Replace(val, "x", "", -1), 16, 32)
			data["info"]["Backplane"] = strconv.Itoa(int(address >> 3))
			data["info"]["Slot"] = strconv.Itoa(int(address & 0x07))
		}
	}

	return data
}