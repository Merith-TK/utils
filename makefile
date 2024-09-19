install:
	for folder in $(wildcard cmd/*); do \
		echo "Installing $$folder"; \
		args=""; \
		if [ -f "$$folder/.buildargs" ]; then \
			args=$$(cat "$$folder/.buildargs"); \
			echo "Args: $$args"; \
		fi; \
		go install $$args ./$$folder; \
	done