COMMON_FLAGS=-count=1 -race $(if $(v), -v) ./e2e --godog.format=pretty,cucumber:report.json $(if $(t),--godog.tags=$(t))
FEATURES_DIR=e2e/features

.PHONY: e2e
e2e:
	@if [ -z "$(f)" ]; then \
		echo "No specified feature files. Running all e2e tests..."; \
		go test $(COMMON_FLAGS); \
	else \
		IFS=',' read -ra features_arr <<< "$(f)"; \
		feature_files=""; \
		for feature_file in $${features_arr[@]}; do \
			if [ ! -f "$(FEATURES_DIR)/$${feature_file}.feature" ]; then \
				printf "Error: $(FEATURES_DIR)/%s.feature does not exist\n" "$${feature_file}"; \
				exit 1; \
			else \
				feature_files+="features/$${feature_file}.feature,"; \
			fi; \
		done; \
		if [ -n "$${feature_files}" ]; then \
			feature_files="$${feature_files%?}"; \
			ALL_FLAGS="$(COMMON_FLAGS) --feature $$feature_files"; \
			echo "Running e2e tests for feature files: $$feature_files"; \
			go test $$ALL_FLAGS; \
		fi; \
	fi
