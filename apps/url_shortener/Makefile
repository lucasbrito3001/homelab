.PHONY: all test mock clean build

# Variáveis
MOCKERY=mockery
PORTS_DIR=internal/app/ports
MOCKS_DIR=$(PORTS_DIR)/mocks

all: build

mock:
	@echo "Gerando mocks com Mockery..."
	@rm -rf $(MOCKS_DIR)
	@$(MOCKERY) --all --dir=$(PORTS_DIR) --output=$(MOCKS_DIR) --case=underscore --outpkg=mocks
	@echo "Mocks gerados com sucesso em $(MOCKS_DIR)"

test:
	@echo "Executando testes..."
	@go test -v ./...

test-cov:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

build:
	@echo "Compilando..."
	@go build -o bin/api cmd/api/main.go

clean:
	@rm -rf bin/
	@rm -f coverage.out
	@echo "Limpeza concluída."