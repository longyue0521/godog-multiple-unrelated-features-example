# godog-multiple-unrelated-features-example
Show how to manage multiple unrelated feature files in a web project and use Godog to build end-to-end acceptance tests


# install godog and gocure

install godog binary used for generating steps definitions of features files.

```shell
go install github.com/cucumber/godog/cmd/godog@v0.12.6
```

import godog package used for writing acceptance tests or end-to-end tests.
```shell
go get -u github.com/cucumber/godog@v0.12.6
```

import gocure used for generating HTML report from above acceptance tests' json reports.
```shell
go get -u gitlab.com/rodrigoodhin/gocure
```

# Step 1: Write Feature file

# Step 2: Generate Steps Definition

1. create step definitions file `e2e/features/users/get_api_steps.go`
2. add a function like following

```go
package users

import (
	"context"

	"github.com/cucumber/godog"
)

// a description name
func GetAPI() *godog.TestSuite {
	return &godog.TestSuite{
		// path to your feature file
		Name: "features/users/get_api.feature",
		TestSuiteInitializer: func(ctx *godog.TestSuiteContext) {
			ctx.BeforeSuite(func() {
				// add setup steps before the test suite to be running
			})
			ctx.AfterSuite(func() {
				// add teardown steps after the test suite stopped
			})
		},
		ScenarioInitializer: func(ctx *godog.ScenarioContext) {
			ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
				// add setup steps before every Scenario to run
				return ctx, nil
			})

			InitializeScenario(ctx)

			ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
				// add teardown steps after every Scenario stopped
				return ctx, err
			})
		},
	}
}
```

```shell
go test -v e2e/e2e_test.go --tag=e2e
```

copy generated step definitions from terminal to `e2e/features/users/get_api_steps.go`

```go
package users

import (
	"context"

	"github.com/cucumber/godog"
)

func GetAPI() *godog.TestSuite {
	return &godog.TestSuite{
		Name: "features/users/get_api.feature",
		TestSuiteInitializer: func(ctx *godog.TestSuiteContext) {
			ctx.BeforeSuite(func() {
				// add setup steps before the test suite to be running
			})
			ctx.AfterSuite(func() {
				// add teardown steps after the test suite stopped
			})
		},
		ScenarioInitializer: func(ctx *godog.ScenarioContext) {
			ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
				// add setup steps before every Scenario to run
				return ctx, nil
			})

			ctx.Step(`^I send "([^"]*)" request to "([^"]*)"$`, iSendRequestTo)
			ctx.Step(`^the response code should be (\d+)$`, theResponseCodeShouldBe)
			ctx.Step(`^the response should match json:$`, theResponseShouldMatchJson)
			ctx.Step(`^there are users:$`, thereAreUsers)

			ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
				// add teardown steps after every Scenario stopped
				return ctx, err
			})
		},
	}
}

func iSendRequestTo(arg1, arg2 string) error {
	return godog.ErrPending
}

func theResponseCodeShouldBe(arg1 int) error {
	return godog.ErrPending
}

func theResponseShouldMatchJson(arg1 *godog.DocString) error {
	return godog.ErrPending
}

func thereAreUsers(arg1 *godog.Table) error {
	return godog.ErrPending
}
```

run `go test e2e/e2e_test.go --tags=e2e` again, you will get error.

```shell
FFeature: users
  In order to use users api
  As an API user
  I need to be able to manage users

  Scenario: should get users when there is only one # features/users/get_api.feature:37
    Given there are users:                          # get_api_steps.go:68 -> github.com/longyue0521/godog-multiple-unrelated-features-example/e2e/features/users.thereAreUsers
      | username | email           |
      | gopher   | gopher@mail.com |
    after scenario hook failed: step implementation is pending, step error: step implementation is pending
    When I send "GET" request to "/users"           # get_api_steps.go:56 -> github.com/longyue0521/godog-multiple-unrelated-features-example/e2e/features/users.iSendRequestTo
    Then the response code should be 200            # get_api_steps.go:60 -> github.com/longyue0521/godog-multiple-unrelated-features-example/e2e/features/users.theResponseCodeShouldBe
    And the response should match json:             # get_api_steps.go:64 -> github.com/longyue0521/godog-multiple-unrelated-features-example/e2e/features/users.theResponseShouldMatchJson
      """
      {
        "users": [
          {
            "username": "gopher"
          }
        ]
      }
      """

  Scenario: should get users              # features/users/get_api.feature:16
    Given there are users:                # get_api_steps.go:68 -> github.com/longyue0521/godog-multiple-unrelated-features-example/e2e/features/users.thereAreUsers
      | username | email             |
      | john     | john.doe@mail.com |
      | jane     | jane.doe@mail.com |
    after scenario hook failed: step implementation is pending, step error: step implementation is pending
    When I send "GET" request to "/users" # get_api_steps.go:56 -> github.com/longyue0521/godog-multiple-unrelated-features-example/e2e/features/users.iSendRequestTo
    Then the response code should be 200  # get_api_steps.go:60 -> github.com/longyue0521/godog-multiple-unrelated-features-example/e2e/features/users.theResponseCodeShouldBe
    And the response should match json:   # get_api_steps.go:64 -> github.com/longyue0521/godog-multiple-unrelated-features-example/e2e/features/users.theResponseShouldMatchJson
      """
      {
        "users": [
          {
            "username": "john"
          },
          {
            "username": "jane"
          }
        ]
      }
      """

  Scenario: should get empty users        # features/users/get_api.feature:6
    When I send "GET" request to "/users" # get_api_steps.go:56 -> github.com/longyue0521/godog-multiple-unrelated-features-example/e2e/features/users.iSendRequestTo
    after scenario hook failed: step implementation is pending, step error: step implementation is pending
    Then the response code should be 200  # get_api_steps.go:60 -> github.com/longyue0521/godog-multiple-unrelated-features-example/e2e/features/users.theResponseCodeShouldBe
    And the response should match json:   # get_api_steps.go:64 -> github.com/longyue0521/godog-multiple-unrelated-features-example/e2e/features/users.theResponseShouldMatchJson
      """
      {
        "users": []
      }
      """

--- Failed steps:

  Scenario: should get empty users # features/users/get_api.feature:6
    When I send "GET" request to "/users" # features/users/get_api.feature:7
      Error: after scenario hook failed: step implementation is pending, step error: step implementation is pending

  Scenario: should get users # features/users/get_api.feature:16
    Given there are users: # features/users/get_api.feature:17
      Error: after scenario hook failed: step implementation is pending, step error: step implementation is pending

  Scenario: should get users when there is only one # features/users/get_api.feature:37
    Given there are users: # features/users/get_api.feature:38
      Error: after scenario hook failed: step implementation is pending, step error: step implementation is pending


3 scenarios (3 failed)
11 steps (3 failed, 8 skipped)
1.024337ms

Randomized with seed: 1685289426157075000
--- FAIL: TestE2E (0.00s)
    --- FAIL: TestE2E/users/get_api (0.00s)
        --- FAIL: TestE2E/users/get_api/should_get_users_when_there_is_only_one (0.00s)
            suite.go:451: after scenario hook failed: step implementation is pending, step error: step implementation is pending
        --- FAIL: TestE2E/users/get_api/should_get_users (0.00s)
            suite.go:451: after scenario hook failed: step implementation is pending, step error: step implementation is pending
        --- FAIL: TestE2E/users/get_api/should_get_empty_users (0.00s)
            suite.go:451: after scenario hook failed: step implementation is pending, step error: step implementation is pending
        e2e_test.go:128: 
                Error Trace:    e2e_test.go:128
                Error:          Not equal: 
                                expected: 0
                                actual  : 1
                Test:           TestE2E/users/get_api
                Messages:       0 - success
                                1 - failed
                                2 - command line usage error
                                128 - or higher, os signal related error exit codes
FAIL
FAIL    command-line-arguments  0.172s
FAIL
```
change all steps:

```go
func iSendRequestTo(arg1, arg2 string) error {
	return nil
}

func theResponseCodeShouldBe(arg1 int) error {
	return nil
}

func theResponseShouldMatchJson(arg1 *godog.DocString) error {
	return nil
}

func thereAreUsers(arg1 *godog.Table) error {
	return nil
}
```


using `@wip` on feature file, and `go test . -v --godog.tags=@wip` and `go test . -v --godog.tags=~@wip`

@wip - run all scenarios with wip tag
~@wip - exclude all scenarios with wip tag
@wip && ~@new - run wip scenarios, but exclude new
@wip,@undone - run wip or undone scenarios