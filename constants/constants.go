package constants



const (

	//Main path
	SYSTEM_DIRECTORY = "/etc/nms-triad/"
	SYSTEM_NETWORK_CONF_DIRECTORY = "/etc/systemd/network/"

	LOCAL_DIRECTORY = "/data/nms-triad/"
	LOCAL_NETWORK_CONF_DIRECTORY = LOCAL_DIRECTORY + "network/"
	LOCAL_SNMP_CONF_DIRECTORY = LOCAL_DIRECTORY + "snmp/"

	ENV_FILE_PATH = LOCAL_DIRECTORY + ".env"
	TEMPLATE_ENV_PATH = "goAssets/.env"

	//HTML templates paths
	TEMPLATE_PAGE_HEADER_PATH = "goAssets/html/header.html"
	TEMPLATE_PAGE_FOOTER_PATH = "goAssets/html/footer.html"
	TEMPLATE_PAGE_MENU_PATH = "goAssets/html/menu.html"
	TEMPLATE_PAGE_LOGIN_PATH = "goAssets/html/login.html"
	TEMPLATE_PAGE_HOME_PATH = "goAssets/html/index.html"
	TEMPLATE_PAGE_INFO_PATH = "goAssets/html/information.html"
	TEMPLATE_PAGE_USERS_PATH = "goAssets/html/users.html"
	TEMPLATE_PAGE_EDIT_USERS_PATH = "goAssets/html/editUser.html"
	TEMPLATE_PAGE_ETHERNET_PATH = "goAssets/html/ethernet.html"
	TEMPLATE_PAGE_WIFI_PATH = "goAssets/html/wifi.html"
	TEMPLATE_PAGE_ERROR_PATH = "goAssets/html/error.html"
	TEMPLATE_PAGE_SUCCESS_PATH = "goAssets/html/success.html"
	TEMPLATE_PAGE_HARDWARE_PATH = "goAssets/html/hardware.html"
	TEMPLATE_PAGE_HARDWARE_INFO_PATH = "goAssets/html/hardwareInfo.html"
	TEMPLATE_PARTIAL_CELL_INFO_INFO_PATH = "goAssets/html/devInfoTableHome.html"
	TEMPLATE_PARTIAL_CELL_REMOTE_INFO_PATH = "goAssets/html/remoteCellTemplate.html"
	TEMPLATE_PAGE_SNMP_PATH = "goAssets/html/snmp.html"
	TEMPLATE_PAGE_CONFIG_PATH = "goAssets/html/configuration.html"
	TEMPLATE_PAGE_REMOTES_PATH = "goAssets/html/remotes.html"

	//Network configuration templates paths
	TEMPLATE_ETH_STATIC_PATH = "goAssets/scripts/eth-static.txt"
	TEMPLATE_ETH_DHCP_PATH = "goAssets/scripts/eth-dhcp.txt"
	TEMPLATE_AP_PATH = "goAssets/scripts/start-ap.txt"
	TEMPLATE_SNMPD_CONF_PATH = "goAssets/scripts/snmpd.txt"

	//Scripts paths
	ETH_SCRIPT_PATH = SYSTEM_DIRECTORY + "eth-conf.sh"
	START_AP_SCRIPT_PATH = SYSTEM_DIRECTORY + "start-ap.sh"
	STOP_AP_SCRIPT_PATH = SYSTEM_DIRECTORY + "stop-ap.sh"

	//Configuration paths
	LOCAL_ETH_CONF_PATH = LOCAL_NETWORK_CONF_DIRECTORY + "eth-conf.sh"
	LOCAL_AP_CONF_PATH = LOCAL_NETWORK_CONF_DIRECTORY + "start-ap.sh"
	LOCAL_SNMPD_CONF_PATH = LOCAL_SNMP_CONF_DIRECTORY + "snmpd.conf"

	//USB vendor ID
	CYPRESS_VENDOR_ID = 0x04B4

	//Serial Commands
	INFO_COMMAND                  = "info\n"
	READ_SENSORS_COMMAND          = "readSensors\n"
	READ_CONTROLS_COMMAND         = "readControls\n"
	READ_ALARMS_COMMAND           = "readAlarms\n"
	READ_ADC_COMMAND              = "readAdc\n"
	READ_ALARMS_THRESHOLD_COMMAND = "readAlarmThresholds\n"
	READ_CAL_RF_DET_COMMAND       = "readCalRFDet\n"
	READ_CAL_RF_DET2_COMMAND      = "readCalRFDet2\n"
	READ_RF_AGC_COMMAND           = "readCalRFAgc\n"
	READ_LASER_APC_COMMAND        = "readLaserAPC\n"
	READ_CAL_OPT_PWR_COMMAND      = "readCalOptPwr\n"

	SET_RF_ATTENUATOR_1_COMMAND = "setRFattenuator1 "
	SET_RF_ATTENUATOR_2_COMMAND = "setRFattenuator2 "
	SET_RF_ATTENUATOR_3_COMMAND = "setRFattenuator3 "
	SET_AGC_MODE_COMMAND        = "setRFAGC "
	SET_RF_SWITCH_COMMAND       = "setRFswitch "
	SET_RF_SQUELCH_COMMAND      = "setRFsquelch "
	SET_LASER_APC_COMMAND       = "setLaserAPC "

	WRITE_THRESHOLDS_COMMAND     = "writeAlarmThresholds "
	WRITE_CAL_RF_DET_COMMAND     = "writeCalRFDet "
	WRITE_CAL_RF_DET2_COMMAND    = "writeCalRFDet2 "
	WRITE_CAL_RF_AGC_COMMAND     = "writeCalRFAgc "
	WRITE_DAC_LASER_BIAS_COMMAND = "writeDacLaserBias "
	WRITE_OPT_DET_OFFSET_COMMAND = "writeOptDetOffset "
	WRITE_OPT_DET_TARGET_COMMAND = "writeOptDetTarget "
    WRITE_CAL_OPT_PWR_COMMAND    = "writeCalOptPwr "

	SAVE_COMMAND         = "save\n"
	RESTORE_COMMAND      = "restore\n"
	CLEAR_ALARMS_COMMAND = "clearAlarms\n"
	RESET_COMMAND        = "reset\n"

	//Hardware types
	MOD_100_REV1 = "CGW-MOD-100-REV1"
	MOD_200_REV1 = "CGW-MOD-200-REV1"
	FOTX_200_REV1 = "CGW-FOTX-200-REV1"
	FORX_200_REV1 = "CGW-FORX-200-REV1"


	//Work request actions

	SerialWorkRequestUpdateDevice = iota
	SerialWorkRequestUpdateSensorsDevice
	SerialWorkRequestUpdateInfoDevice
	SerialWorkRequestExecCommand
	SerialWorkRequestCustom

	SNMPWorkRequestGetHostname
	SNMPWorkRequestGetAlarms
)

