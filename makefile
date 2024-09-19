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

build: 
	for folder in $(wildcard cmd/*); do \
		echo "Building $$folder"; \
		args=""; \
		if [ -f "$$folder/.buildargs" ]; then \
			args=$$(cat "$$folder/.buildargs"); \
			echo "Args: $$args"; \
		fi; \
		go build -o .build/ $$args ./$$folder; \
	done