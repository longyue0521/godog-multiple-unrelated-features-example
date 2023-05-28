package godogs

import (
	"context"

	"github.com/cucumber/godog"
)

func Eat() *godog.TestSuite {
	return &godog.TestSuite{
		Name: "features/godogs/godogs.feature",
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

			ctx.Step(`^I eat (\d+)$`, iEat)
			ctx.Step(`^there are (\d+) godogs$`, thereAreGodogs)
			ctx.Step(`^there should be none remaining$`, thereShouldBeNoneRemaining)
			ctx.Step(`^there should be (\d+) remaining$`, thereShouldBeRemaining)

			ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
				// add teardown steps after every Scenario stopped
				return ctx, err
			})
		},
	}
}

func iEat(arg1 int) error {
	return nil
}

func thereAreGodogs(arg1 int) error {
	return nil
}

func thereShouldBeNoneRemaining() error {
	return nil
}

func thereShouldBeRemaining(arg1 int) error {
	return nil
}
