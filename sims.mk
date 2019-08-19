#!/usr/bin/make -f

########################################
### Simulations

SIMAPP = github.com/zar-network/zar-network/app

sim-zar-nondeterminism:
	@echo "Running non-determinism test..."
	@go test -mod=readonly $(SIMAPP) -run TestAppStateDeterminism -Enabled=true \
		-NumBlocks=100 -BlockSize=200 -Commit=true -v -timeout 24h

sim-zar-custom-genesis-fast:
	@echo "Running custom genesis simulation..."
	@echo "By default, ${HOME}/.zard/config/genesis.json will be used."
	@go test -mod=readonly $(SIMAPP) -run TestFullZarSimulation -Genesis=${HOME}/.zard/config/genesis.json \
		-Enabled=true -NumBlocks=100 -BlockSize=200 -Commit=true -Seed=99 -Period=5 -v -timeout 24h

sim-zar-fast:
	@echo "Running quick Zar simulation. This may take several minutes..."
	@go test -mod=readonly $(SIMAPP) -run TestFullZarSimulation -Enabled=true -NumBlocks=100 -BlockSize=200 -Commit=true -Seed=99 -Period=5 -v -timeout 24h

sim-zar-import-export: runsim
	@echo "Running Zar import/export simulation. This may take several minutes..."
	$(GOPATH)/bin/runsim $(SIMAPP) 25 5 TestZarImportExport

sim-zar-simulation-after-import: runsim
	@echo "Running Zar simulation-after-import. This may take several minutes..."
	$(GOPATH)/bin/runsim $(SIMAPP) 25 5 TestZarSimulationAfterImport

sim-zar-custom-genesis-multi-seed: runsim
	@echo "Running multi-seed custom genesis simulation..."
	@echo "By default, ${HOME}/.zard/config/genesis.json will be used."
	$(GOPATH)/bin/runsim -g ${HOME}/.zard/config/genesis.json 400 5 TestFullZarSimulation

sim-zar-multi-seed: runsim
	@echo "Running multi-seed Zar simulation. This may take awhile!"
	$(GOPATH)/bin/runsim $(SIMAPP) 400 5 TestFullZarSimulation

sim-benchmark-invariants:
	@echo "Running simulation invariant benchmarks..."
	@go test -mod=readonly $(SIMAPP) -benchmem -bench=BenchmarkInvariants -run=^$ \
	-Enabled=true -NumBlocks=1000 -BlockSize=200 \
	-Commit=true -Seed=57 -v -timeout 24h

SIM_NUM_BLOCKS ?= 500
SIM_BLOCK_SIZE ?= 200
SIM_COMMIT ?= true

sim-zar-benchmark:
	@echo "Running Zar benchmark for numBlocks=$(SIM_NUM_BLOCKS), blockSize=$(SIM_BLOCK_SIZE). This may take awhile!"
	@go test -mod=readonly -benchmem -run=^$$ $(SIMAPP) -bench ^BenchmarkFullZarSimulation$$  \
		-Enabled=true -NumBlocks=$(SIM_NUM_BLOCKS) -BlockSize=$(SIM_BLOCK_SIZE) -Commit=$(SIM_COMMIT) -timeout 24h

sim-zar-profile:
	@echo "Running Zar benchmark for numBlocks=$(SIM_NUM_BLOCKS), blockSize=$(SIM_BLOCK_SIZE). This may take awhile!"
	@go test -mod=readonly -benchmem -run=^$$ $(SIMAPP) -bench ^BenchmarkFullZarSimulation$$ \
		-Enabled=true -NumBlocks=$(SIM_NUM_BLOCKS) -BlockSize=$(SIM_BLOCK_SIZE) -Commit=$(SIM_COMMIT) -timeout 24h -cpuprofile cpu.out -memprofile mem.out


.PHONY: runsim sim-zar-nondeterminism sim-zar-custom-genesis-fast sim-zar-fast sim-zar-import-export \
	sim-zar-simulation-after-import sim-zar-custom-genesis-multi-seed sim-zar-multi-seed \
	sim-benchmark-invariants sim-zar-benchmark sim-zar-profile
