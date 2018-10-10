package command

import (
	"bytes"
	"fmt"
	"os/exec"
	"syscall"
)

func CloneRepo(repo string, targetDir string) error {
	cmd := exec.Command("git", "clone", repo, ".")
	cmd.Dir = targetDir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	fmt.Printf("'%s' clone> '%s' : ", targetDir, repo)

	if err := cmd.Run(); err != nil {
		fmt.Printf("[ERROR] %s\n", string(stderr.Bytes()))
		return err
	}

	fmt.Println("done")
	return nil
}

func CheckRemoteExists(repo string, name string) bool {
	cmd := exec.Command("git", "show-branch", "remotes/origin/"+name)
	cmd.Dir = repo
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()

	code := getErrorCode(err)

	return code == 0
}

func CheckoutMaster(repo string) error {
	cmd := exec.Command("git", "checkout", "master")
	cmd.Dir = repo
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()

	code := getErrorCode(err)

	fmt.Printf("'%s' checkout-master> : ", repo)
	if code != 0 {
		fmt.Printf("[ERROR] %s\n", string(stderr.Bytes()))
		return err
	}

	fmt.Println("done")

	return nil
}

func ResetHeadHard(repo string) error {
	cmd := exec.Command("git", "reset", "--hard")
	cmd.Dir = repo
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()

	code := getErrorCode(err)

	fmt.Printf("'%s' reset-head-hard> : ", repo)
	if code != 0 {
		fmt.Printf("[ERROR] %s\n", string(stderr.Bytes()))
		return err
	}

	fmt.Println("done")

	return nil
}

func CreateBranch(repo string, name string) error {
	cmd := exec.Command("git", "checkout", "-b", name)
	cmd.Dir = repo
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()

	code := getErrorCode(err)

	fmt.Printf("'%s' branch> '%s' : ", repo, name)
	if code != 0 {
		fmt.Printf("[ERROR] %s\n", string(stderr.Bytes()))
		return err
	}

	fmt.Println("done")

	return nil
}

func StageChanges(repo string) error {
	cmd := exec.Command("git", "add", "-u")
	cmd.Dir = repo
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()

	code := getErrorCode(err)

	fmt.Printf("'%s' stage> : ", repo)
	if code != 0 {
		fmt.Printf("[ERROR] %s\n", string(stderr.Bytes()))
		return err
	}

	fmt.Println("done")
	return nil
}

func CommitChanges(repo string, message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Dir = repo
	var stderr, stdout bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	err := cmd.Run()

	code := getErrorCode(err)

	fmt.Printf("'%s' commit> '%s' : ", repo, message)
	if code != 0 {
		fmt.Printf("[ERROR] %s\n", string(stderr.Bytes()))
		fmt.Println(string(stdout.Bytes()))
		return err
	}

	fmt.Println("done")
	return nil
}

func PushBranch(repo string, name string) error {
	cmd := exec.Command("git", "push", "origin", name)
	cmd.Dir = repo
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()

	code := getErrorCode(err)

	fmt.Printf("'%s' push> '%s' : ", repo, name)

	if code != 0 {
		fmt.Printf("[ERROR] %s\n", string(stderr.Bytes()))
		return err
	}

	fmt.Println("done")
	return nil
}

func getErrorCode(err error) int {
	code := 0
	if msg, ok := err.(*exec.ExitError); ok { // there is error code
		code = msg.Sys().(syscall.WaitStatus).ExitStatus()
	}

	return code
}
