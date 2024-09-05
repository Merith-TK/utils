install:
	for folder in $(wildcard cmd/*); do \
		echo "Installing $$folder"; \
		go install ./$$folder; \
	done