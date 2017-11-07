install:
	cd backend && dep ensure
	cd backend && go install ./...
	cd frontend && npm install

build:
	cd frontend && npm run build
	go-bindata-assetfs \
		-o=backend/api/assets/bindata_assetfs.go \
		-pkg=assets \
		-prefix=frontend/build \
		frontend/build/...
	go build -tags=binassets -o backend/build/cmdui backend/main.go

build_linux:
	cd backend && docker build . -t cmdui:build
	cd backend && docker run -i -v ${PWD}/backend/build:/artefacts cmdui:build \
		go build -o /artefacts/cmdui_linux main.go

create_db:
	sqlite3 backend/data/cmdui.db < backend/schema_sqlite.sql

cloc:
	cloc --exclude-dir=vendor,build,node_modules .
