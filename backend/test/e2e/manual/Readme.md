Tests in this directory are not meant to be run by CICD automation, but rather manually by developers on their machines.
They serve as workflow tests and mimic the sequence of actions that would be performed via the UI. These
tests will typically connect to real data-sources, so any data-source specific data/credential needs to be supplied
externally. The convention we are using is to wrap all such variables in a Config struct, placed in a `models.go` file,
which is loaded in at runtime. Populating these structs and making them available to the tests will be your responsibility.
See the example below.

You can also add your own manual tests by having the test files follow the pattern *_local_test.go to exclude them
from git's tracking.

Example:

In `models.go` define
```go
package azuredevops

    type TestConfig struct {
        Org     string
        Project string
        Token   string
    }
```

and load it into your test function (if you write one) via
```go
    cfg := helper.GetTestConfig[TestConfig]()
```

In `azure_local_test.go` (or any git-ignorable file) you write your setup.
```go
package azuredevops

import "github.com/apache/incubator-devlake/test/helper"

func init() {
	helper.SetTestConfig(TestConfig{
		Org:     "???",
		Project: "???",
		Token:   "??????",
	})
}

// Your custom test cases (optional)
```