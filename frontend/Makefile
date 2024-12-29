.PHONY: build clean

DIST_DIR = dist
TARGET_DIR = ../embedded/admin/static

build: node_modules
	npm run build
	@mkdir -p $(TARGET_DIR)/js $(TARGET_DIR)/css
	@cp $(DIST_DIR)/assets/*.js $(TARGET_DIR)/js/app.js
	@cp $(DIST_DIR)/assets/*.css $(TARGET_DIR)/css/app.css

node_modules: package.json
	npm install

clean:
	rm -rf $(DIST_DIR)
	rm -rf node_modules
