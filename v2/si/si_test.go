package si

import (
	"fmt"
	"log"
)

func ExampleRead() {
	si, err := Read("ossf", "si-tooling", ".github/security-insights.yml")
	if err != nil {
		log.Fatalf("error reading Security Insights data from ossf/si-tooling/.github/security-insights.yml: %v", err)
	}
	fmt.Println(string(si.Repository.License.Url))

	// output:
	// https://github.com/ossf/security-insights-spec?tab=Apache-2.0-1-ov-file#readme
}

func ExampleLoad() {
	si, err := Load(minimalTestData())
	if err != nil {
		log.Fatalf("error loading Security Insights data: %v", err)
	}
	fmt.Println(string(si.Repository.License.Url))

	// output:
	// https://example.com/LICENSE
}
