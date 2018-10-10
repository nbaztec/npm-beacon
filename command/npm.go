package command

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"time"
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

func GetPackageReleaseDate(name string, version string) (time.Time, error) {
	cmd := exec.Command("npm", "view", name+"@"+version, "time", "--json")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	fmt.Printf("get-version-date> %s@%s : ", name, version)

	code := getErrorCode(err)
	if code != 0 {
		fmt.Printf("[ERROR] %s\n", string(stderr.Bytes()))
		return time.Now(), err
	}
	// fmt.Printf("code : %d\n", code)
	// fmt.Printf("code : %d\n", string(stdout.Bytes()))

	var response map[string]string
	if err = json.Unmarshal(stdout.Bytes(), &response); err != nil {
		fmt.Printf("[ERROR] %s\n", err)
		return time.Now(), err
	}

	v, ok := response[version]
	if ok == false {
		fmt.Printf("[ERROR] version %s found in response\n", version)
		return time.Now(), errors.New("version '" + version + "' not found in npm response")
	}

	t, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return time.Now(), errors.New("time '" + v + "' cannot be parsed")
	}

	fmt.Println("done")

	return t, nil
}

type npmOutdatedPackage struct {
	Current  string
	Wanted   string
	Latest   string
	Location string
}
