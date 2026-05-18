.PHONY: dev build clean

dev:
	@echo "Starting dev server..."
	@cd frontend && npm run dev & cd backend && air -- cmd/server/main.go

build:
	@cd frontend && npm run build
	@cd backend && go build -o warmisle .

clean:
	@rm -rf frontend/dist backend/frontend/dist backend/warmisle backend/warmisle.exe
