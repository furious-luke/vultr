package cmd

import (
	"fmt"
	"log"

	vultr "github.com/JamesClonk/vultr/lib"
	"github.com/jawher/mow.cli"
)

func serversCreate(cmd *cli.Cmd) {
	cmd.Spec = "-n -r -p -o [OPTIONS]"

	name := cmd.StringOpt("n name", "", "Name of new virtual machine")
	regionID := cmd.IntOpt("r region", 0, "Region (DCID)")
	planID := cmd.IntOpt("p plan", 0, "Plan (VPSPLANID)")
	osID := cmd.IntOpt("o os", 0, "Operating system (OSID)")

	// options
	ipxe := cmd.StringOpt("ipxe", "", "Chainload the specified URL on bootup, via iPXE, for custom OS")
	iso := cmd.IntOpt("iso", 0, "ISOID of a specific ISO to mount during the deployment, for custom OS")
	script := cmd.IntOpt("s script", 0, "SCRIPTID of a startup script to execute on boot (see <scripts>)")
	snapshot := cmd.StringOpt("snapshot", "", "SNAPSHOTID (see <snapshots>) to restore for the initial installation")
	sshkey := cmd.StringOpt("k sshkey", "", "SSHKEYID (see <sshkeys>) of SSH key to apply to this server on install")
	ipv6 := cmd.BoolOpt("ipv6", false, "Assign an IPv6 subnet to this virtual machine (where available)")
	privateNetworking := cmd.BoolOpt("private-networking", false, "Add private networking support for this virtual machine")
	autoBackups := cmd.BoolOpt("autobackups", false, "Enable automatic backups for this virtual machine")

	cmd.Action = func() {
		options := &vultr.ServerOptions{
			IPXEChainURL:      *ipxe,
			ISO:               *iso,
			Script:            *script,
			Snapshot:          *snapshot,
			SSHKey:            *sshkey,
			IPV6:              *ipv6,
			PrivateNetworking: *privateNetworking,
			AutoBackups:       *autoBackups,
		}

		server, err := GetClient().CreateServer(*name, *regionID, *planID, *osID, options)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Virtual machine created\n")
		lengths := []int{12, 32, 8, 12, 8}
		printTabbedLine(Columns{"SUBID", "NAME", "DCID", "VPSPLANID", "OSID"}, lengths)
		printTabbedLine(Columns{server.ID, server.Name, server.RegionID, server.PlanID, *osID}, lengths)
		tabsFlush()
	}
}

func serversRename(cmd *cli.Cmd) {
	cmd.Spec = "SUBID -n"
	id := cmd.StringArg("SUBID", "", "SUBID of virtual machine (see <servers>)")
	name := cmd.StringOpt("n name", "", "new name of virtual machine")
	cmd.Action = func() {
		if err := GetClient().RenameServer(*id, *name); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Virtual machine renamed to: %v\n", *name)
	}
}

func serversStart(cmd *cli.Cmd) {
	id := cmd.StringArg("SUBID", "", "SUBID of virtual machine (see <servers>)")
	cmd.Action = func() {
		if err := GetClient().StartServer(*id); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Virtual machine (re)started")
	}
}

func serversHalt(cmd *cli.Cmd) {
	id := cmd.StringArg("SUBID", "", "SUBID of virtual machine (see <servers>)")
	cmd.Action = func() {
		if err := GetClient().HaltServer(*id); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Virtual machine halted")
	}
}

func serversReboot(cmd *cli.Cmd) {
	id := cmd.StringArg("SUBID", "", "SUBID of virtual machine (see <servers>)")
	cmd.Action = func() {
		if err := GetClient().RebootServer(*id); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Virtual machine rebooted")
	}
}

func serversReinstall(cmd *cli.Cmd) {
	id := cmd.StringArg("SUBID", "", "SUBID of virtual machine (see <servers>)")
	cmd.Action = func() {
		if err := GetClient().ReinstallServer(*id); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Virtual machine reinstalled")
	}
}

func serversChangeOS(cmd *cli.Cmd) {
	cmd.Spec = "SUBID -o"
	id := cmd.StringArg("SUBID", "", "SUBID of virtual machine (see <servers>)")
	osID := cmd.IntOpt("o os", 0, "Operating system (OSID)")
	cmd.Action = func() {
		if err := GetClient().ChangeOSofServer(*id, *osID); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Virtual machine operating system changed to: %v\n", *osID)
	}
}

func serversListOS(cmd *cli.Cmd) {
	id := cmd.StringArg("SUBID", "", "SUBID of virtual machine (see <servers>)")
	cmd.Action = func() {
		os, err := GetClient().ListOSforServer(*id)
		if err != nil {
			log.Fatal(err)
		}

		if len(os) == 0 {
			fmt.Println()
			return
		}

		lengths := []int{8, 32, 8, 16, 8, 12}
		printTabbedLine(Columns{"OSID", "NAME", "ARCH", "FAMILY", "WINDOWS", "SURCHARGE"}, lengths)
		for _, os := range os {
			printTabbedLine(Columns{os.ID, os.Name, os.Arch, os.Family, os.Windows, os.Surcharge}, lengths)
		}
		tabsFlush()
	}
}

func serversDelete(cmd *cli.Cmd) {
	id := cmd.StringArg("SUBID", "", "SUBID of virtual machine (see <servers>)")
	cmd.Action = func() {
		if err := GetClient().DeleteServer(*id); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Virtual machine deleted")
	}
}

func serversBandwidth(cmd *cli.Cmd) {
	id := cmd.StringArg("SUBID", "", "SUBID of virtual machine (see <servers>)")
	cmd.Action = func() {
		bandwidth, err := GetClient().BandwidthOfServer(*id)
		if err != nil {
			log.Fatal(err)
		}

		if len(bandwidth) == 0 {
			fmt.Println()
			return
		}

		lengths := []int{24, 24, 24}
		printTabbedLine(Columns{"DATE", "INCOMING", "OUTGOING"}, lengths)
		for _, b := range bandwidth {
			printTabbedLine(Columns{b["date"], b["incoming"], b["outgoing"]}, lengths)
		}
		tabsFlush()
	}
}

func serversList(cmd *cli.Cmd) {
	cmd.Action = func() {
		servers, err := GetClient().GetServers()
		if err != nil {
			log.Fatal(err)
		}

		if len(servers) == 0 {
			fmt.Println()
			return
		}

		lengths := []int{12, 16, 24, 32, 32, 32, 8, 8, 24, 12, 8}
		printTabbedLine(Columns{
			"SUBID",
			"STATUS",
			"IP",
			"NAME",
			"OS",
			"LOCATION",
			"VCPU",
			"RAM",
			"DISK",
			"BANDWIDTH",
			"COST"}, lengths)
		for _, server := range servers {
			printTabbedLine(Columns{
				server.ID,
				server.Status,
				server.MainIP,
				server.Name,
				server.OS,
				server.Location,
				server.VCpus,
				server.RAM,
				server.Disk,
				server.AllowedBandwidth,
				server.Cost,
			}, lengths)
		}
		tabsFlush()
	}
}

func serversShow(cmd *cli.Cmd) {
	cmd.Spec = "SUBID [-f | --full]"

	id := cmd.StringArg("SUBID", "", "SUBID of virtual machine (see <servers>)")
	full := cmd.BoolOpt("f full", false, "Display full length of KVM URL")

	cmd.Action = func() {
		server, err := GetClient().GetServer(*id)
		if err != nil {
			log.Fatal(err)
		}

		if server.ID == "" {
			fmt.Printf("No virtual machine with SUBID %v found!\n", *id)
			return
		}

		keyLength := 64
		if *full {
			keyLength = 1024
		}
		lengths := []int{24, keyLength}

		printTabbedLine(Columns{"Id (SUBID):", server.ID}, lengths)
		printTabbedLine(Columns{"Name:", server.Name}, lengths)
		printTabbedLine(Columns{"Operating system:", server.OS}, lengths)
		printTabbedLine(Columns{"Status:", server.Status}, lengths)
		printTabbedLine(Columns{"Power status:", server.PowerStatus}, lengths)
		printTabbedLine(Columns{"Location:", server.Location}, lengths)
		printTabbedLine(Columns{"Region (DCID):", server.RegionID}, lengths)
		printTabbedLine(Columns{"VCPU count:", server.VCpus}, lengths)
		printTabbedLine(Columns{"RAM:", server.RAM}, lengths)
		printTabbedLine(Columns{"Disk:", server.Disk}, lengths)
		printTabbedLine(Columns{"Allowed bandwidth:", server.AllowedBandwidth}, lengths)
		printTabbedLine(Columns{"Current bandwidth:", server.CurrentBandwidth}, lengths)
		printTabbedLine(Columns{"Cost per month:", server.Cost}, lengths)
		printTabbedLine(Columns{"Pending charges:", server.PendingCharges}, lengths)
		printTabbedLine(Columns{"Plan (VPSPLANID):", server.PlanID}, lengths)
		printTabbedLine(Columns{"IP:", server.MainIP}, lengths)
		printTabbedLine(Columns{"Netmask:", server.NetmaskV4}, lengths)
		printTabbedLine(Columns{"Gateway:", server.GatewayV4}, lengths)
		printTabbedLine(Columns{"Internal IP:", server.InternalIP}, lengths)
		printTabbedLine(Columns{"IPv6 IP:", server.MainIPV6}, lengths)
		printTabbedLine(Columns{"IPv6 Network:", server.NetworkV6}, lengths)
		printTabbedLine(Columns{"IPv6 Network Size:", server.NetworkSizeV6}, lengths)
		printTabbedLine(Columns{"Created date:", server.Created}, lengths)
		printTabbedLine(Columns{"Default password:", server.DefaultPassword}, lengths)
		printTabbedLine(Columns{"Auto backups:", server.AutoBackups}, lengths)
		printTabbedLine(Columns{"KVM URL:", server.KVMUrl}, lengths)
		tabsFlush()
	}
}