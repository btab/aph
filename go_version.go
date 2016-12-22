package go-project-harness

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"runtime"
	"sort"
)

func checkLatestGoVersionNumber() {
	latestGVN, err := latestGoVersionNumber()
	if err != nil {
		log.Printf("WARNING: unable to ascertain latest go version number: %s", err)
	} else if runtime.Version() != latestGVN {
		log.Printf("WARNING: current go version number: %s, latest: %s", runtime.Version(), latestGVN)
	} else {
		log.Printf("go version is latest: %s", latestGVN)
	}
}

func latestGoVersionNumber() (string, error) {
	url := "https://golang.org/doc/devel/release.html"
	versionMatcher := "go\\d+\\.\\d+\\.\\d+"

	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	b := bytes.Buffer{}
	_, err = b.ReadFrom(response.Body)
	if err != nil {
		return "", err
	}

	matches := regexp.MustCompile(versionMatcher).FindAllString(b.String(), -1)

	if len(matches) == 0 {
		return "", fmt.Errorf("no strings matching '%s' found in response from '%s'", versionMatcher, url)
	}

	// this is mostly OK as (currently) all releases use single digit major, minor and patch numbers
	sort.Strings(matches)

	return matches[len(matches)-1], nil
}
