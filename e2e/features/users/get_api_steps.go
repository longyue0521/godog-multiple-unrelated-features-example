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
