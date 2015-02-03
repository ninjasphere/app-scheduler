build-deps:
	npm install
	tar -C app -czf app-bower-components.tgz bower_components

unpack-build-deps:
	rm -rf app/bower_components
	tar -C app -xcf app-bower-components.tgz