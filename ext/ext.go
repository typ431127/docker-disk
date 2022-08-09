package ext

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PodInfo struct {
	Name       string
	Kubernetes bool
	Deployment string
	Podname    string
	Namespace  string
	Restart    string
	Size       float64
	ID         string
	FormatUnit string
	PausePOD   bool
	Status     string
}

func ParseName(name string) PodInfo {
	if strings.HasPrefix(name, "k8s") {
		parse := false
		namelist := strings.Split(name, "_")
		if namelist[1] == "POD" {
			parse = true
		}
		return PodInfo{
			Name:       name,
			Kubernetes: true,
			Deployment: namelist[1],
			Podname:    namelist[2],
			Namespace:  namelist[3],
			Restart:    namelist[5],
			PausePOD:   parse,
		}
	} else {
		return PodInfo{
			Name:       name,
			Kubernetes: false,
			PausePOD:   false,
		}
	}

}

func UnitCheck(size string) bool {
	unit := size[len(size)-1:]
	for _, un := range []string{"k", "m", "g", "t"} {
		if unit == un {
			return true
		}
	}
	return false
}

func UnitParse(size string) float64 {
	size = strings.ToLower(size)
	unit := size[len(size)-1:]
	var sizefloat64 float64
	if UnitCheck(size) == false {
		fmt.Println("-size单位错误,只能选择 k m g t")
		os.Exit(2)
	}
	sizefloat64, _ = strconv.ParseFloat(strings.Split(size, unit)[0], 64)
	switch {
	case unit == "k":
		sizefloat64 = sizefloat64 * float64(1000)
	case unit == "m":
		sizefloat64 = sizefloat64 * float64(1000*1000)
	case unit == "g":
		sizefloat64 = sizefloat64 * float64(1000*1000*1000)
	case unit == "t":
		sizefloat64 = sizefloat64 * float64(1000*1000*1000*1000)
	}
	return sizefloat64

}
func UnitConvert(fileSize int64) (size string) {
	if fileSize < 1000 {
		return fmt.Sprintf("%.2fB", float64(fileSize)/float64(1))
	} else if fileSize < (1000 * 1000) {
		return fmt.Sprintf("%.2fKB", float64(fileSize)/float64(1000))
	} else if fileSize < (1000 * 1000 * 1000) {
		return fmt.Sprintf("%.2fMB", float64(fileSize)/float64(1000*1000))
	} else if fileSize < (1000 * 1000 * 1000 * 1000) {
		return fmt.Sprintf("%.2fGB", float64(fileSize)/float64(1000*1000*1000))
	} else if fileSize < (1000 * 1000 * 1000 * 1000 * 1000) {
		return fmt.Sprintf("%.2fTB", float64(fileSize)/float64(1000*1000*1000*1000))
	} else {
		return fmt.Sprintf("%.2fEB", float64(fileSize)/float64(1000*1000*1000*1000*1000))
	}
}
