package handler

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/nbaztec/npm-beacon/command"
)

func Process2(repositories []string, githubToken string) bool {
	success := true
	for _, repo := range repositories {
		if err := handleOneRepository(repo, githubToken); err != nil {
			fmt.Printf("Failed processing: '%s' > %v\n", repo, err)
			success = false
		}
	}

	return success
}

func Process(repositories []string, githubToken string) bool {
	result := make(chan bool)
	done := make(chan bool)
	totalTasks := len(repositories)

	go resultAggregator(result, done, totalTasks)

	for _, repo := range repositories {
		go process(repo, githubToken, result)
	}

	return <-done
}

func resultAggregator(result <-chan bool, done chan<- bool, totalTasks int) {
	current := true
	counter := 0
	for r := range result {
		current = current && r
		counter++

		if counter == totalTasks {
			break
		}
	}

	done <- current
}

func process(repo string, githubToken string, result chan<- bool) {
	if err := handleOneRepository(repo, githubToken); err != nil {
		fmt.Printf("Failed processing: '%s' > %v\n", repo, err)
		result <- false
	}

	result <- true
}

func handleOneRepository(repository string, githubToken string) error {
	tempDir, err := ioutil.TempDir(".beacon", "repo-")
	defer os.RemoveAll(tempDir)

	fmt.Printf("Executing in '%s' fo '%s'\n", tempDir, repository)
	if err != nil {
		return err
	}

	if err = command.CloneRepo(repository, tempDir); err != nil {
		fmt.Println("Failed to clone repository")
		return err
	}

	packages, err := getOutdatePackages(tempDir)
	if err != nil {
		fmt.Println("Failed to retrieve outdated packages")
		return err
	}

	for _, pkg := range packages {
		err = createBranchAndUpdatePackage(repository, tempDir, pkg, githubToken)
		if err != nil {
			fmt.Printf("Error handling '%s' > '%v' : %v\n", repository, pkg, err)
		}
	}

	return nil
}

func getOutdatePackages(dir string) ([]command.OutdatedPackage, error) {
	packages, err := command.GetOutdatePackages(dir)
	if err != nil {
		return nil, err
	}

	return packages, nil
}

func generateBranchName(pkg command.OutdatedPackage) string {
	return fmt.Sprintf("update-%s-%s", strings.ToLower(pkg.Name), pkg.Latest)
}
func createBranchAndUpdatePackage(repository string, dir string, pkg command.OutdatedPackage, githubToken string) error {
	branch := generateBranchName(pkg)

	// check if a PR is already created, if yes then do nothing
	if command.CheckRemoteExists(dir, branch) {
		fmt.Printf("remote already exists: '%s'\n", branch)
		return nil
	}

	if err := command.CheckoutMaster(dir); err != nil {
		fmt.Println("error creating branch")
		return err
	}

	if err := command.CreateBranch(dir, branch); err != nil {
		fmt.Println("error creating branch")
		return err
	}

	if err := updatePackageVersion(dir, pkg); err != nil {
		fmt.Println("error updating package version")
		return err
	}

	if err := command.StageChanges(dir); err != nil {
		fmt.Println("error staging changes")
		return err
	}

	commitMessage := fmt.Sprintf("updates package to %s", pkg.Latest)
	if err := command.CommitChanges(dir, commitMessage); err != nil {
		fmt.Println("error commiting changes")
		return err
	}

	if err := command.PushBranch(dir, branch); err != nil {
		return err
	}

	prTitle := fmt.Sprintf("updates %s to %s", pkg.Name, pkg.Latest)
	prBody := fmt.Sprintf("This is an automated pull request to update the package **`%s`** from `%s` to `%s`.\n\n"+
		"If the tests are green the PR can be safely merged.", pkg.Name, pkg.Wanted, pkg.Latest)
	if err := command.OpenPullRequest(githubToken, repository, branch, prTitle, prBody); err != nil {
		fmt.Println("error opening pull request")
		return err
	}

	return nil
}

func updatePackageVersion(dir string, pkg command.OutdatedPackage) error {
	filename := path.Join(dir, "package.json")
	input, err := ioutil.ReadFile(filename)

	fmt.Printf("'%s' update-package> '%s', '%s' => '%s' : ", dir, pkg.Name, pkg.Wanted, pkg.Latest)

	if err != nil {
		fmt.Println("[ERROR] Error reading package.json")
		return err
	}

	output := strings.Replace(string(input), fmt.Sprintf(`"%s": "%s"`, pkg.Name, pkg.Wanted), fmt.Sprintf(`"%s": "%s"`, pkg.Name, pkg.Latest), -1)
	err = ioutil.WriteFile(filename, []byte(output), 0644)
	if err != nil {
		fmt.Println("[ERROR] Error writing package.json")
		return err
	}

	fmt.Println("done")

	return nil
}
