package testing_ondatra

import (
	"fmt"
	"os"
	"testing"

	"github.com/openconfig/gnoigo/system"
	"github.com/openconfig/ondatra"
	"github.com/openconfig/ondatra/gnmi"
	"github.com/openconfig/ondatra/gnmi/oc"
	"github.com/openconfig/ondatra/gnoi"

	kinit "github.com/openconfig/ondatra/knebind/init"
)

func TestMain(m *testing.M) {
	//Init the kinit to read the kne yaml file
	ondatra.RunTests(m, kinit.Init)
}

func Allduts(t *testing.T) []string {
	devs := ondatra.DUTs(t)
	var dutsslice []string
	for _, r := range devs {
		//Loop through the names of the duts and append it to dusslice
		dutsslice = append(dutsslice, r.Name())
	}
	return dutsslice
}

func TestAllConfigs(t *testing.T) {
	//Get the pwd
	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	for _, i := range Allduts(t) {
		dut := ondatra.DUT(t, i)
		config, err := os.ReadFile(path + "/" + i + "-append.cfg")
		if err != nil {
			fmt.Print(err)
		}
		//Either config can be appended or wiped.
		//dut.Config().New().WithAristaText(string(config)).Push(t)
		dut.Config().New().WithAristaText(string(config)).Append(t)
	}
}

func TestGNMISystem(t *testing.T) {
	dut := ondatra.DUT(t, "r1")
	sys := gnmi.Lookup(t, dut, gnmi.OC().System().Hostname().State())
	configHostname, _ := sys.Val()
	if configHostname != "r1" {
		t.Errorf("Expected result to be %s, got %s", configHostname, "r1")
	}

}

func TestGNMILLDPNeighbors(t *testing.T) {
	dut := ondatra.DUT(t, "r1")
	lldp := gnmi.Lookup(t, dut, gnmi.OC().Lldp().State())
	if !lldp.IsPresent() {
		t.Fatalf("No LLDP for %v", dut)
	}
}

func TestInterfaceEth1(t *testing.T) {
	InterfaceEth := gnmi.OC().Interface("Ethernet1").AdminStatus().State()
	dut := ondatra.DUT(t, "r1")
	Ethernet1 := gnmi.Lookup(t, dut, InterfaceEth)
	if !Ethernet1.IsPresent() {
		t.Fatalf("Check Interface %v is present", dut)
	}
	if result, ok := Ethernet1.Val(); ok && result == oc.Interface_AdminStatus_DOWN {
		t.Fatalf("Check Interface %v is up", dut)
	}
}

func TestDestination(t *testing.T) {
	dut := ondatra.DUT(t, "r2")
	dest := gnoi.Execute(t, dut, system.NewPingOperation().Destination("8.8.8.8").Count(5))
	for _, d := range dest {
		if d.Received == 5 {
			if d.Received == 0 {
				t.Fatalf("Cannot ping destination")
			}
		}
	}
}