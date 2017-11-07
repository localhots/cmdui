install:
	cd backend && dep ensure
	cd backend && go install ./...
	cd frontend && npm install

build: assets
	go build -tags=binassets -o backend/build/cmdui backend/main.go
	@echo Build complete!
	@echo Binary is located at backend/build/cmdui

build_linux: assets build_docker

build_docker:
	cd backend && docker build . -t cmdui:build
	cd backend && docker run -i -v ${PWD}/backend/build:/build cmdui:build \
		go build -tags=binassets -o /build/cmdui_linux main.go
	@echo Build complete!
	@echo Binary is located at backend/build/cmdui_linux

assets:
	cd frontend && npm run build
	go-bindata-assetfs \
		-o=backend/api/assets/bindata_assetfs.go \
		-pkg=assets \
		-prefix=frontend/build \
		frontend/build/...

clean:
	rm -rf frontend/build
	rm -rf backend/build
	rm -f backend/api/assets/bindata_assetfs.go

create_db:
	sqlite3 backend/data/cmdui.db < backend/schema_sqlite.sql

cloc:
	cloc --exclude-dir=vendor,build,node_modules .
