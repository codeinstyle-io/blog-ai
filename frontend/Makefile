.PHONY: build clean dev format format-check

DIST_DIR = dist
TARGET_DIR = ../embedded/admin/static

build: node_modules
	npm run build
	@mkdir -p $(TARGET_DIR)/js $(TARGET_DIR)/css
	@cp $(DIST_DIR)/app.js $(TARGET_DIR)/js/app.js
	@cp $(DIST_DIR)/index.css $(TARGET_DIR)/css/app.css

node_modules: package.json
	npm install

clean:
	rm -rf $(DIST_DIR)
	rm -rf node_modules

dev: node_modules
	npm run dev

format: node_modules
	npx prettier --write "src/**/*.{js,ts,svelte,css}"

format-check: node_modules
	npx prettier --check "src/**/*.{js,ts,svelte,css}"
