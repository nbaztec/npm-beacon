package command

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

type OutdatedPackage struct {
	Name    string
	Current string
	Wanted  string
	Latest  string
}

func GetOutdatePackages(dir string) ([]OutdatedPackage, error) {
	cmd := exec.Command("npm", "outdated", "--json")
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	// code := getErrorCode(err)

	// fmt.Printf("code: %v\n", code)
	// fmt.Printf("err: %v\n", string(stderr.Bytes()))
	// fmt.Printf("output: %v \n", string(stdout.Bytes()))

	// out := stdout.Bytes()
	// fmt.Println(">>>>" + string(out))
	var response map[string]npmOutdatedPackage
	if err = json.Unmarshal(stdout.Bytes(), &response); err != nil {
		// fmt.Println(err)
		return nil, err
	}

	// var outdatedPackages []OutdatedPackage
	// // detect outdated packages
	// for k, v := range response {
	// 	vWanted, _ := version.NewVersion(v.Wanted)
	// 	vLatest, _ := version.NewVersion(v.Latest)

	// 	if vWanted.Equal(vLatest) || vWanted.GreaterThan(vLatest) {
	// 		continue
	// 	}

	// 	pkg := OutdatedPackage{
	// 		Name:    k,
	// 		Current: v.Current,
	// 		Latest:  v.Latest,
	// 		Wanted:  v.Wanted,
	// 	}

	// 	outdatedPackages = append(outdatedPackages, pkg)

	// }

	// fmt.Println(len(outdatedPackages))
	// fmt.Printf("%v\n", outdatedPackages)

	// count packages that need updating
	packageLen := 0
	for _, v := range response {
		if v.Wanted != v.Latest {
			packageLen++
		}
	}

	outdatedPackages := make([]OutdatedPackage, packageLen)
	i := 0
	for k, v := range response {
		if v.Wanted == v.Latest {
			continue
		}

		outdatedPackages[i] = OutdatedPackage{
			Name:    k,
			Current: v.Current,
			Latest:  v.Latest,
			Wanted:  v.Wanted,
		}
		i++
	}

	fmt.Printf("'%s' found> %d outdated packages\n", dir, i)
	// fmt.Printf("%v\n", outdatedPackages)

	return outdatedPackages, nil
}

type npmOutdatedPackage struct {
	Current  string
	Wanted   string
	Latest   string
	Location string
}
