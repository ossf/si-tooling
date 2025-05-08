# v2/si-tooling [![Go Reference](https://pkg.go.dev/badge/github.com/ossf/si-tooling/v2.svg)](https://pkg.go.dev/github.com/ossf/si-tooling/v2)

This is a go module for working with [Security Insights](https://github.com/ossf/security-insights-spec) data in YAML `security-insights.yml` and Go `si.SecurityInsights`.

## Usage

Unmarshal the `security-insights.yml` data in [ossf/security-insights-spec](https://github.com/ossf/security-insights-spec)

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
