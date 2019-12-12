package snmp

const (
	EnterpriseTree = "1.3.6.1.4.1"
	CellgainID     = ".47182"
	CellgainTree   = EnterpriseTree + CellgainID

	HostNameOID = ".1.3.6.1.2.1.1.5.0"
	NMSTree = CellgainTree + ".1"

	CgNumberBackplane = NMSTree + ".1.0"

	//Tables
	CgBackplaneTable = NMSTree + ".2"
	CgSlotTable      = NMSTree + ".3"
	CgDeviceTable    = NMSTree + ".4"
	CgSensorsTable   = NMSTree + ".5"

	//Entries
	CgBackplaneTableEntry = CgBackplaneTable + ".1"
	CgSlotTableEntry      = CgSlotTable + ".1"
	CgDeviceTableEntry    = CgDeviceTable + ".1"
	CgSensorsTableEntry   = CgSensorsTable + ".1"

	//Backplane table columns
	CgBackplaneTableEntryIndex = CgBackplaneTableEntry + ".1"
	CgBackplaneTableEntryDescription = CgBackplaneTableEntry + ".2"
	CgBackplaneTableEntryNumberSlot = CgBackplaneTableEntry + ".3"

	//Slot table columns
	CgSlotTableEntryIndex = CgSlotTableEntry + ".1"
	CgSlotTableEntryDescription = CgSlotTableEntry + ".2"
	CgSlotTableEntryNumberDevice = CgSlotTableEntry + ".3"

	//Device table columns
	CgDeviceTableEntryIndex = CgDeviceTableEntry + ".1"
	CgDeviceTableEntryDescription = CgDeviceTableEntry + ".2"
	CgDeviceTableEntrySerial = CgDeviceTableEntry + ".3"
	CgDeviceTableEntryFwVersion = CgDeviceTableEntry + ".4"
	CgDeviceTableEntryProductID = CgDeviceTableEntry + ".5"
	CgDeviceTableEntryAlarms = CgDeviceTableEntry + ".6"

	//Control table columns
	CgSensorsTableEntryLaserBias    = CgSensorsTableEntry + ".1"
	CgSensorsTableEntryOpticalPower = CgSensorsTableEntry + ".2"
	CgSensorsTableEntryRFPower      = CgSensorsTableEntry + ".3"
	CgSensorsTableEntryTemperature  = CgSensorsTableEntry + ".4"
	CgSensorsTableEntryRFPowerIn    = CgSensorsTableEntry + ".5"
	CgSensorsTableEntryRFPowerOut   = CgSensorsTableEntry + ".6"
)
