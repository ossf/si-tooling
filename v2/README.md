# v2/si-tooling [![Go Reference](https://pkg.go.dev/badge/github.com/ossf/si-tooling.svg)](https://pkg.go.dev/github.com/ossf/si-tooling)

This is a go module for marshaling/unmarshaling `security-insights.yml` data.

## Usage

Unmarshal the `security-insights.yml` data in `ossf/security-insights-spec`

`github.com/ossf/si-tooling/v2`

```go
import (
    "fmt"

    "github.com/ossf/si-tooling/v2/si"
)

func main() {
    insights, err := si.Read("ossf", "security-insights-spec", ".github/security-insights.yml")
    message = fmt.Sprintf("Repository license is: %s", insights.Repository.License.Expression)
}
```