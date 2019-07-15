#!/usr/bin/make -f

########################################
### Simulations

SIMAPP = github.com/zar-network/zar-network/app

sim-zar-nondeterminism:
	@echo "Running nondeterminism test..."
	@go test -mod=readonly $(SIMAPP) -run TestAppStateDeterminism -SimulationEnabled=true -v -timeout 10m

sim-zar-custom-genesis-fast:
	@echo "Running custom genesis simulation..."
	@echo "By default, ${HOME}/.zard/config/genesis.json will be used."
	@go test -mod=readonly github.com/zar-network/zar-network/app -run TestFullZarSimulation -SimulationGenesis=${HOME}/.zard/config/genesis.json \
		-SimulationEnabled=true -SimulationNumBlocks=100 -SimulationBlockSize=200 -SimulationCommit=true -SimulationSeed=99 -SimulationPeriod=5 -v -timeout 24h

sim-zar-fast:
	@echo "Running quick Zar simulation. This may take several minutes..."
	@go test -mod=readonly github.com/zar-network/zar-network/app -run TestFullZarSimulation -SimulationEnabled=true -SimulationNumBlocks=100 -SimulationBlockSize=200 -SimulationCommit=true -SimulationSeed=99 -SimulationPeriod=5 -v -timeout 24h

sim-zar-import-export: runsim
	@echo "Running Zar import/export simulation. This may take several minutes..."
	$(GOPATH)/bin/runsim 25 5 TestZarImportExport

sim-zar-simulation-after-import: runsim
	@echo "Running Zar simulation-after-import. This may take several minutes..."
	$(GOPATH)/bin/runsim 25 5 TestZarSimulationAfterImport

sim-zar-custom-genesis-multi-seed: runsim
	@echo "Running multi-seed custom genesis simulation..."
	@echo "By default, ${HOME}/.zard/config/genesis.json will be used."
	$(GOPATH)/bin/runsim -g ${HOME}/.zard/config/genesis.json 400 5 TestFullZarSimulation

sim-zar-multi-seed: runsim
	@echo "Running multi-seed Zar simulation. This may take awhile!"
	$(GOPATH)/bin/runsim 400 5 TestFullZarSimulation

sim-benchmark-invariants:
	@echo "Running simulation invariant benchmarks..."
	@go test -mod=readonly github.com/zar-network/zar-network/app -benchmem -bench=BenchmarkInvariants -run=^$ \
	-SimulationEnabled=true -SimulationNumBlocks=1000 -SimulationBlockSize=200 \
	-SimulationCommit=true -SimulationSeed=57 -v -timeout 24h

SIM_NUM_BLOCKS ?= 500
SIM_BLOCK_SIZE ?= 200
SIM_COMMIT ?= true
sim-zar-benchmark:
	@echo "Running Zar benchmark for numBlocks=$(SIM_NUM_BLOCKS), blockSize=$(SIM_BLOCK_SIZE). This may take awhile!"
	@go test -mod=readonly -benchmem -run=^$$ github.com/zar-network/zar-network/app -bench ^BenchmarkFullZarSimulation$$  \
		-SimulationEnabled=true -SimulationNumBlocks=$(SIM_NUM_BLOCKS) -SimulationBlockSize=$(SIM_BLOCK_SIZE) -SimulationCommit=$(SIM_COMMIT) -timeout 24h

sim-zar-profile:
	@echo "Running Zar benchmark for numBlocks=$(SIM_NUM_BLOCKS), blockSize=$(SIM_BLOCK_SIZE). This may take awhile!"
	@go test -mod=readonly -benchmem -run=^$$ github.com/zar-network/zar-network/app -bench ^BenchmarkFullZarSimulation$$ \
		-SimulationEnabled=true -SimulationNumBlocks=$(SIM_NUM_BLOCKS) -SimulationBlockSize=$(SIM_BLOCK_SIZE) -SimulationCommit=$(SIM_COMMIT) -timeout 24h -cpuprofile cpu.out -memprofile mem.out


.PHONY: runsim sim-zar-nondeterminism sim-zar-custom-genesis-fast sim-zar-fast sim-zar-import-export \
	sim-zar-simulation-after-import sim-zar-custom-genesis-multi-seed sim-zar-multi-seed \
	sim-benchmark-invariants sim-zar-benchmark sim-zar-profile
