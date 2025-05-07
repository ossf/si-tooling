package si

import (
	"fmt"
	"log"
)

func Example_read() {
	// read Insights data directly from a public repository using Read
	si, err := Read("ossf", "si-tooling", ".github/security-insights.yml")
	if err != nil {
		log.Fatalf("error reading Security Insights data from ossf/si-tooling/.github/security-insights.yml: %v", err)
	}
	fmt.Println(si.Repository.License.Url)

	// output:
	// https://github.com/ossf/si-tooling?tab=Apache-2.0-1-ov-file#readme
}

func Example_load() {
	// load Insights data from a byte slice using Load
	si, err := Load(minimalTestData())
	if err != nil {
		log.Fatalf("error loading Security Insights data: %v", err)
	}
	fmt.Println(si.Repository.License.Url)

	// output:
	// https://example.com/LICENSE
}
