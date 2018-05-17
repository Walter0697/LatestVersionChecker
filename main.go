package main

import (
	"context"
	"fmt"
	"strings"
	"os"

	"github.com/coreos/go-semver/semver"
	"github.com/google/go-github/github"
)

// Latest Versions retruns a sorted slice with the highest version as its first element and the highest version of the smaller minor versins in a descending order
func LatestVersions(releases []*semver.Version, minVersion *semver.Version) []*semver.Version {
	var versionSlice []*semver.Version

	//sorting the list using the library function
	semver.Sort(releases)

	//checking index, since sorting will sort in incending order, and we would like to find the highest version first
	//therefore we would go from the highest value to check
	i := len(releases) - 1

	//set the current_version to impossible to reach
	current_version := semver.New("999.999.999")

	//i represents the checking index, so negative number will end the loop
	for i >= 0 {

		if releases[i].LessThan(*minVersion) {
			return versionSlice
		}

		checking_version := findVersion(releases[i])

		//append if it is another version
		//again, since it is sorted, the first one must be the highest one
		if !checking_version.Equal(*current_version) {
			versionSlice = append(versionSlice, releases[i])
			current_version = checking_version
		}
		i--
	}
	

	return versionSlice
}

// Find out what higher version is this
// it will force the third number into zero
// such as 1.4.3 -> 1.4.0
func findVersion(release *semver.Version) *semver.Version {
	//convert the version into String first
	versionString := release.String()
	resultString := ""

	//find out the first "."
	i := strings.Index(versionString, ".")
	resultString += versionString[:i+1]
	versionString = versionString[i+1:]

	
	//find out the second "."
	i = strings.Index(versionString, ".")
	resultString += versionString[:i] + ".0"

	return semver.New(resultString)
}


func main() {
	//getting the information from command line argument
	var github_repo, github_name string
	var minVersion *semver.Version
	
	//can have two ways to use this program
	//example 1 : go run main.go kubernetes/kubernetes 1.8.2
	//example 2 : go run main.go kubernetes kubernetes 1.8.2
	if len(os.Args) == 3 {
		//find out the position of '/'
		pos := strings.Index(os.Args[1], "/")
		github_repo = os.Args[1][:pos]
		github_name = os.Args[1][pos+1:]
		minVersion = semver.New(os.Args[2])
	} else if len(os.Args) >= 4 {
		//ignore the extra argument when there are more than 4 arguments
		github_repo = os.Args[1]
		github_name = os.Args[2]
		minVersion = semver.New(os.Args[3])
	} else {
		fmt.Println("invalid numbers of command line argument")
		return
	}

	// Github
	client := github.NewClient(nil)
	ctx := context.Background()
	opt := &github.ListOptions{PerPage: 10}
	releases, _, err := client.Repositories.ListReleases(ctx, github_repo, github_name, opt)

	//command line error is already checked above
	if err != nil {
		fmt.Println("github error")
		//print out the error 
		fmt.Println(err.Error())
	}

	allReleases := make([]*semver.Version, len(releases))
	for i, release := range releases {
		versionString := *release.TagName
		if versionString[0] == 'v' {
			versionString = versionString[1:]
		}
		fmt.Println(versionString)
		allReleases[i] = semver.New(versionString)
	}

	versionSlice := LatestVersions(allReleases, minVersion)

	fmt.Printf("latest versions of %s/%s: %s\n", github_repo, github_name, versionSlice)
}