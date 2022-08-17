package app

import (
	"context"
	"docker-disk/config"
	"docker-disk/ext"
	"flag"
	"fmt"
	"github.com/docker/docker/client"
	"github.com/olekukonko/tablewriter"
	"os"
	"sort"
)

func init() {
	flag.StringVar(&config.Size, "size", "0k", "指定容器存储大小")
	flag.StringVar(&config.WithVersion, "withversion", "1.39", "指定客户端API版本，不指定则自动协商API版本")
	flag.BoolVar(&config.K8s, "k8s", false, "解析k8s名称格式")
	flag.BoolVar(&config.NoPause, "nopause", false, "不显示pause容器")
	flag.BoolVar(&config.PrintDelete, "delete-commond", false, "生成k8s pod删除语句")
}

func Show() {
	flag.Parse()
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	sizefloat64 := ext.UnitParse(config.Size)
	var Containers []ext.PodInfo
	var cli *client.Client
	var err error
	if config.WithVersion == "1.39" {
		cli, err = client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	} else {
		cli, err = client.NewClientWithOpts(client.WithVersion(config.WithVersion))
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	defer cli.Close()
	data, err := cli.DiskUsage(context.Background())
	if err != nil {
		panic(err)
	}
	for _, con := range data.Containers {
		if con.SizeRootFs >= int64(sizefloat64) {
			podinfo := ext.ParseName(con.Names[0][1:])
			podinfo.Size = float64(con.SizeRootFs)
			podinfo.Image = con.Image
			podinfo.ID = con.ID[:12]
			podinfo.Status = con.State
			podinfo.FormatUnit = ext.UnitConvert(con.SizeRootFs)
			Containers = append(Containers, podinfo)
		}
	}
	sort.SliceStable(Containers, func(i, j int) bool {
		return Containers[i].Size > Containers[j].Size
	})
	for _, container := range Containers {
		var table []string
		if config.K8s {
			if container.Kubernetes {
				if config.NoPause {
					if container.PausePOD == false {
						table = append(table, container.Podname, container.Namespace, container.ID, container.FormatUnit, container.Restart, container.Status)
						config.Tabledata = append(config.Tabledata, table)
					}
				} else {
					table = append(table, container.Podname, container.Namespace, container.ID, container.FormatUnit, container.Restart, container.Status)
					config.Tabledata = append(config.Tabledata, table)
				}
			}

		} else {
			table = append(table, container.Name, container.ID, container.FormatUnit, container.Image, container.Status)
			config.Tabledata = append(config.Tabledata, table)
		}
	}
	var deletetext string
	if config.PrintDelete && config.K8s {
		for _, container := range config.Tabledata {
			deletetext = deletetext + fmt.Sprintf("kubectl delete pod %s -n %s\n", container[0], container[1])
		}
		fmt.Printf(deletetext)
		os.Exit(0)
	}
	table := tablewriter.NewWriter(os.Stdout)
	if config.K8s {
		table.SetHeader([]string{"Name", "NAMESPACE", "ID", "SIZE", "RESTART", "STATUS"})
	} else {
		table.SetHeader([]string{"Name", "ID", "SIZE", "IMAGE", "STATUS"})
	}
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)
	table.AppendBulk(config.Tabledata)
	table.Render()
}
