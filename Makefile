.PHONY: dev build clean test lint fmt e2e-install e2e e2e-ui e2e-update e2e-report

dev:
	@echo "Starting dev server..."
	@cd frontend && npm run dev & cd backend && air -- cmd/server/main.go

build:
	@cd frontend && npm run build -- --emptyOutDir
	@mkdir -p dist
	@cd backend && go build -o ../dist/warmisle .

clean:
	@rm -rf frontend/dist backend/frontend/dist dist

test:
	cd backend && go test ./...
	cd frontend && npm test

lint:
	cd backend && golangci-lint run ./...
	cd frontend && npx eslint src/

fmt:
	cd backend && gofmt -w .
	cd frontend && npx eslint --fix src/

# E2E 测试
e2e-install:
	cd e2e && npx playwright install chromium

e2e: build
	cd e2e && npx playwright test

e2e-ui:
	cd e2e && npx playwright test --ui

e2e-update: build
	cd e2e && npx playwright test --update-snapshots

e2e-report:
	cd e2e && npx playwright show-report
