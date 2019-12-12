package snmpTree

import (
	"github.com/moffa90/nms-server/constants/snmp"
	"github.com/moffa90/nms-server/db"
	"github.com/moffa90/nms-server/db/models"
	"fmt"
	"os"
	"strconv"
)

type Node struct {
	tag      string
	id       string
	class    string //table, entry, valueInt, valueStr, group, field
	children []*Node
	oid      string
	value    string
}

func NewNode(tag string, id string, class string, children []*Node, oid string, value string) *Node {
	return &Node{tag: tag, id: id, class: class, children: children, oid: oid, value: value}
}

var head Node

func findByTag(root *Node, tag string) *Node {
	queue := make([]*Node, 0)
	queue = append(queue, root)
	for len(queue) > 0 {
		nextUp := queue[0]
		queue = queue[1:]
		if nextUp.tag == tag {
			return nextUp
		}
		if len(nextUp.children) > 0 {
			for _, child := range nextUp.children {
				queue = append(queue, child)
			}
		}
	}
	return nil
}

func RegisterTree(node *Node, closure func(*Node)){
	if node.class == "valueInt" || node.class == "valueStr" {
		closure(node)
	}

	if len(node.children) > 0 {
		for _, child := range node.children {
			RegisterTree(child, closure)
		}
	}
}

func BuildNMSTree() *Node{

	// Callback function returning an int32 value
	mode := os.Getenv("mode")
	var backplaneNumber int
	var slotQty int
	var hardwareQty int
	switch mode {
	case "hec-1":
		backplaneNumber = 1
		slotQty = 5
		hardwareQty = 2
		break

	case "hec-2":
		backplaneNumber = 1
		slotQty = 5
		hardwareQty = 2
		break

	case "remote":
		backplaneNumber = 1
		slotQty = 4
		hardwareQty = 1
		break

	default:
		backplaneNumber =  0
		slotQty = 0
		hardwareQty = 0
		break
	}

	head = *NewNode("Cellgain", "47182", "group", []*Node{}, snmp.CellgainTree, "")

	head.children = append(head.children, NewNode("NMS", "1", "group", []*Node{}, snmp.NMSTree, ""))

	if n := findByTag(&head, "NMS"); n != nil {
		n.children = append(n.children, NewNode("cgNumberBackplanes", "1", "valueInt", nil, snmp.CgNumberBackplane, strconv.Itoa(backplaneNumber)))
		n.children = append(n.children, NewNode("cgBackplaneTable", "2", "table", []*Node{}, snmp.CgBackplaneTable, ""))
		n.children = append(n.children, NewNode("cgSlotTable", "3", "table", []*Node{}, snmp.CgSlotTable, ""))
		n.children = append(n.children, NewNode("cgDeviceTable", "4", "table", []*Node{}, snmp.CgDeviceTable, ""))
		n.children = append(n.children, NewNode("cgSensorsTable", "5", "table", []*Node{}, snmp.CgSensorsTable, ""))
	}

	//Backplane table
	//TODO: Modularize
	if n := findByTag(&head, "cgBackplaneTable"); n != nil {
		n.children = append(n.children, NewNode("cgBackplaneEntry", "1", "entry", []*Node{}, snmp.CgBackplaneTableEntry, ""))
	}

	if n := findByTag(&head, "cgBackplaneEntry"); n != nil {
		n.children = append(n.children, NewNode("cgBackplaneEntryIndex", "1", "field", []*Node{}, snmp.CgBackplaneTableEntryIndex, ""))
		n.children = append(n.children, NewNode("cgBackplaneEntryDescription", "2", "field", []*Node{}, snmp.CgBackplaneTableEntryDescription, ""))
		n.children = append(n.children, NewNode("cgBackplaneEntryNumberSlot", "3", "field", []*Node{}, snmp.CgBackplaneTableEntryNumberSlot, ""))
	}

	if n := findByTag(&head, "cgBackplaneEntryIndex"); n != nil {
		for i := 0; i < backplaneNumber; i++ {
			n.children = append(n.children, NewNode("cgBackplaneEntryIndexValue", strconv.Itoa(i + 1), "valueInt", nil, snmp.CgBackplaneTableEntryIndex + "." + strconv.Itoa(i+1), strconv.Itoa(i+1)))
		}

	}

	if n := findByTag(&head, "cgBackplaneEntryDescription"); n != nil {
		for i := 0; i < backplaneNumber; i++ {
			n.children = append(n.children, NewNode("cgBackplaneEntryDescriptionValue", strconv.Itoa(i + 1), "valueStr", nil, snmp.CgBackplaneTableEntryDescription + "." + strconv.Itoa(i+1), "Backplane #" + strconv.Itoa(i)))
		}

	}

	if n := findByTag(&head, "cgBackplaneEntryNumberSlot"); n != nil {
		for i := 0; i < backplaneNumber; i++ {
			n.children = append(n.children, NewNode("cgBackplaneEntryNumberSlotValue", strconv.Itoa(i + 1), "valueInt", nil, snmp.CgBackplaneTableEntryNumberSlot + "." + strconv.Itoa(i + 1), strconv.Itoa(slotQty)))
		}

	}

	//Slot table
	//TODO: Modularize
	if n := findByTag(&head, "cgSlotTable"); n != nil {
		n.children = append(n.children, NewNode("cgSlotEntry", "1", "entry", []*Node{}, snmp.CgSlotTableEntry, ""))
	}


	if n := findByTag(&head, "cgBackplaneEntry"); n != nil {
		n.children = append(n.children, NewNode("cgSlotEntryIndex", "1", "field", []*Node{}, snmp.CgSlotTableEntryIndex, ""))
		n.children = append(n.children, NewNode("cgSlotEntryDescription", "2", "field", []*Node{}, snmp.CgSlotTableEntryDescription, ""))
		n.children = append(n.children, NewNode("cgSlotEntryNumberDevice", "3", "field", []*Node{}, snmp.CgSlotTableEntryNumberDevice, ""))
	}

	if n := findByTag(&head, "cgSlotEntryIndex"); n != nil {
		for i := 0; i < backplaneNumber; i++ {
			for j := 0; j < slotQty; j++ {
				n.children = append(n.children,
					NewNode(
						"cgSlotEntryIndexValue",
						strconv.Itoa(j + 1),
						"valueInt",
						nil,
						snmp.CgSlotTableEntryIndex + "." + strconv.Itoa(i + 1) + "." + strconv.Itoa(j + 1),
						strconv.Itoa(j + 1),
					),
				)
			}
		}
	}

	if n := findByTag(&head, "cgSlotEntryDescription"); n != nil {
		for i := 0; i < backplaneNumber; i++ {
			for j := 0; j < slotQty; j++ {
				n.children = append(n.children,
					NewNode(
						"cgSlotEntryDescriptionValue",
						strconv.Itoa(j + 1),
						"valueStr",
						nil,
						snmp.CgSlotTableEntryDescription + "." + strconv.Itoa(i + 1) + "." + strconv.Itoa(j + 1),
						"Backplane #" + strconv.Itoa(i) +" Slot #" + strconv.Itoa(j + 1),
					),
				)
			}
		}
	}

	if n := findByTag(&head, "cgSlotEntryNumberDevice"); n != nil {
		for i := 0; i < backplaneNumber; i++ {
			for j := 0; j < slotQty; j++ {
				n.children = append(n.children, NewNode("cgSlotEntryNumberDeviceValue", strconv.Itoa(j + 1), "valueInt", nil, snmp.CgSlotTableEntryNumberDevice + "." + strconv.Itoa(i + 1) + "." + strconv.Itoa(j + 1), strconv.Itoa(hardwareQty)))
			}
		}
	}

	//Device table
	//TODO: Modularize

	if n := findByTag(&head, "cgDeviceTable"); n != nil {
		n.children = append(n.children, NewNode("cgDeviceEntry", "1", "entry", []*Node{}, snmp.CgDeviceTableEntry, ""))
	}


	if n := findByTag(&head, "cgDeviceEntry"); n != nil {
		n.children = append(n.children, NewNode("cgDeviceEntryIndex", "1", "field", []*Node{}, snmp.CgDeviceTableEntryIndex, ""))
		n.children = append(n.children, NewNode("cgDeviceEntryDescription", "2", "field", []*Node{}, snmp.CgDeviceTableEntryDescription, ""))
		n.children = append(n.children, NewNode("cgDeviceEntrySerial", "3", "field", []*Node{}, snmp.CgDeviceTableEntrySerial, ""))
		n.children = append(n.children, NewNode("cgDeviceEntryFwVersion", "4", "field", []*Node{}, snmp.CgDeviceTableEntryFwVersion, ""))
		n.children = append(n.children, NewNode("cgDeviceEntryProductID", "5", "field", []*Node{}, snmp.CgDeviceTableEntryProductID, ""))
		n.children = append(n.children, NewNode("cgDeviceEntryAlarms", "6", "field", []*Node{}, snmp.CgDeviceTableEntryAlarms, ""))
	}

	devices, _ := models.GetHardware(db.Shared)

	if len(devices) > 0 {
		for _, d:= range devices{
			if n := findByTag(&head, "cgDeviceEntryIndex"); n != nil {
				n.children = append(n.children,
					NewNode(
						"cgDeviceEntryIndexValue",
						"1",
						"valueInt",
						nil,
						snmp.CgDeviceTableEntryIndex + "." + strconv.Itoa(d.Backplane + 1) + "." + strconv.Itoa(d.Address) + ".1",
						"1",
					),
				)
			}

			if n := findByTag(&head, "cgDeviceEntryDescription"); n != nil {
				n.children = append(n.children,
					NewNode(
						"cgDeviceEntryDescriptionValue",
						"1",
						"valueStr",
						nil,
						snmp.CgDeviceTableEntryDescription + "." + strconv.Itoa(d.Backplane + 1) + "." + strconv.Itoa(d.Address) + ".1",
						"Backplane #" + strconv.Itoa(d.Backplane) + ", Slot #" + strconv.Itoa(d.Address) + ", Device #1",
					),
				)
			}

			if n := findByTag(&head, "cgDeviceEntrySerial"); n != nil {
				n.children = append(n.children,
					NewNode(
						"cgDeviceEntrySerialValue",
						"1",
						"valueStr",
						nil,
						snmp.CgDeviceTableEntrySerial + "." + strconv.Itoa(d.Backplane + 1) + "." + strconv.Itoa(d.Address) + ".1",
						d.Serial,
					),
				)
			}

			if n := findByTag(&head, "cgDeviceEntryFwVersion"); n != nil {
				n.children = append(n.children,
					NewNode(
						"cgDeviceEntryFwVersionValue",
						"1",
						"valueStr",
						nil,
						snmp.CgDeviceTableEntryFwVersion + "." + strconv.Itoa(d.Backplane + 1) + "." + strconv.Itoa(d.Address) + ".1",
						d.FWversion,
					),
				)
			}

			if n := findByTag(&head, "cgDeviceEntryProductID"); n != nil {
				n.children = append(n.children,
					NewNode(
						"cgDeviceEntryProductIDValue",
						"1",
						"valueStr",
						nil,
						snmp.CgDeviceTableEntryProductID + "." + strconv.Itoa(d.Backplane + 1) + "." + strconv.Itoa(d.Address) + ".1",
						d.ProductId,
					),
				)
			}

			if n := findByTag(&head, "cgDeviceEntryAlarms"); n != nil {
				n.children = append(n.children,
					NewNode(
						"cgDeviceEntryAlarmsValue",
						"1",
						"valueStr",
						nil,
						snmp.CgDeviceTableEntryAlarms + "." + strconv.Itoa(d.Backplane + 1) + "." + strconv.Itoa(d.Address) + ".1",
							fmt.Sprintf("0x%08X", d.Alarms),
					),
				)
			}
		}

	}


	//Sensors table
	//TODO: Modularize
	if n := findByTag(&head, "cgSensorsTable"); n != nil {
		n.children = append(n.children, NewNode("cgSensorsEntry", "1", "entry", []*Node{}, snmp.CgSensorsTableEntry, ""))
	}

	if n := findByTag(&head, "cgSensorsEntry"); n != nil {
		n.children = append(n.children, NewNode("cgSensorsEntryLaserBias", "1", "field", []*Node{}, snmp.CgSensorsTableEntryLaserBias, ""))
		n.children = append(n.children, NewNode("cgSensorsEntryOpticalPower", "2", "field", []*Node{}, snmp.CgSensorsTableEntryOpticalPower, ""))
		n.children = append(n.children, NewNode("cgSensorsEntryRFPower", "3", "field", []*Node{}, snmp.CgSensorsTableEntryRFPower, ""))
		n.children = append(n.children, NewNode("cgSensorsEntryTemperature", "4", "field", []*Node{}, snmp.CgSensorsTableEntryTemperature, ""))
		n.children = append(n.children, NewNode("cgSensorsEntryRFPowerIn", "5", "field", []*Node{}, snmp.CgSensorsTableEntryRFPowerIn, ""))
		n.children = append(n.children, NewNode("cgSensorsEntryRFPowerOut", "6", "field", []*Node{}, snmp.CgSensorsTableEntryRFPowerOut, ""))
	}

	if len(devices) > 0 {
		for _, d:= range devices{
			characteristics, _ := models.GetHardwareCharacteristicsByCat(d.Serial, "sensors", db.Shared)
			for _, c:= range characteristics{

				switch c.Key {
				case "LaserBias":
					if n := findByTag(&head, "cgSensorsEntryLaserBias"); n != nil {
						n.children = append(n.children,
							NewNode(
								"cgSensorsEntryLaserBiasValue",
								"1",
								"valueStr",
								nil,
								snmp.CgSensorsTableEntryLaserBias+ "." + strconv.Itoa(d.Backplane + 1) + "." + strconv.Itoa(d.Address) + ".1",
								c.Key + ": " + c.Value,
							),
						)
					}
				case "RFPowerInp":
					if n := findByTag(&head, "cgSensorsEntryRFPowerIn"); n != nil {
						n.children = append(n.children,
							NewNode(
								"cgSensorsEntryRFPowerInValue",
								"1",
								"valueStr",
								nil,
								snmp.CgSensorsTableEntryRFPowerIn+ "." + strconv.Itoa(d.Backplane + 1) + "." + strconv.Itoa(d.Address) + ".1",
								c.Key + ": " + c.Value,
							),
						)
					}
					break
					break
				}

			}
		}
	}

	return &head
}

func (n Node) Value() string {
	return n.value
}

func (n Node) Oid() string {
	return n.oid
}

func (n Node) Children() []*Node {
	return n.children
}

func (n Node) Class() string {
	return n.class
}

func (n Node) Id() string {
	return n.id
}

func (n Node) Tag() string {
	return n.tag
}