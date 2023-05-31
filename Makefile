.PHONY: e2e_all
e2e_all:
	@go test -race ./e2e

.PHONY:e2e_with_progress_and_cucumber
e2e_with_progress_and_cucumber:
	@go test -race ./e2e --godog.format=progress,cucumber:report.json

.PHONY:e2e_with_pretty_and_cucumber
e2e_with_pretty_and_cucumber:
	@go test -race ./e2e --godog.format=cucumber:report.json,pretty

.PHONY: e2e_with_feature
e2e_with_feature:
	@if [ -z "$(f)" ]; then \
		echo "no features ..."; \
	else \
		test_result=0; \
		IFS=',' read -ra features_arr <<< "$${f}"; \
		for feature_file in $${features_arr[@]}; do \
			if [ ! -f "e2e/features/$${feature_file}.feature" ]; then \
				printf "Error: e2e/features/%s.feature does not exist\n" "$${feature_file}"; \
				test_result=1; \
			fi; \
		done; \
		if [ $$test_result -eq 0 ]; then \
			go test -race ./e2e --godog.format=pretty,cucumber:report.json --feature $$(printf 'features/%s.feature ' $${features_arr[@]}); \
		fi; \
	fi

