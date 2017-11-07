install:
	@echo "* Installing Go dependencies"
	@cd backend && dep ensure
	@echo "* Pre-building Go code"
	@cd backend && go install ./...
	@echo "* Installing JavaScript dependencies"
	@cd frontend && npm install &> /dev/null

build: assets
	@echo "* Building binary"
	@go build -tags=binassets -o backend/build/cmdui backend/main.go
	@echo "* Build complete!"
	@echo "* Binary is located at backend/build/cmdui"

build_linux: assets build_docker

build_docker:
	@echo "* Preparing Docker image"
	@cd backend && docker build . -t cmdui:build &>/dev/null
	@echo "* Building Linux binary"
	@cd backend && docker run -i -v ${PWD}/backend/build:/build cmdui:build \
			go build -tags=binassets -o /build/cmdui_linux main.go \
			&>/dev/null
	@echo "* Creating GZIP archive"
	@gzip -kf backend/build/cmdui_linux
	@echo "* Build complete!"
	@echo "* Binary: backend/build/cmdui_linux"
	@echo "* GZIP:   backend/build/cmdui_linux.gz"

assets:
	@echo "* Compiling asset files"
	@cd frontend && npm run build &>/dev/null
	@echo "* Building embedded assets file"
	@go-bindata-assetfs \
		-o=backend/api/assets/bindata_assetfs.go \
		-pkg=assets \
		-prefix=frontend/build \
		frontend/build/... \
		&>/dev/null

clean:
	@echo "* Removing build artefacts"
	@rm -rf frontend/build
	@rm -rf backend/build
	@rm -f backend/api/assets/bindata_assetfs.go

create_db:
	@echo "* Creating a blank database"
	@sqlite3 backend/data/cmdui.db < backend/schema_sqlite.sql
	@echo "* Here it is: backend/data/cmdui.db"

cloc:
	@cloc --exclude-dir=vendor,build,node_modules .
