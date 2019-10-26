#!/usr/bin/make -f

########################################
### Simulations

SIMAPP = github.com/barkisnet/barkis/app

sim-barkis-nondeterminism:
	@echo "Running nondeterminism test..."
	@go test -mod=readonly $(SIMAPP) -run TestAppStateDeterminism -Enabled=true \
		-NumBlocks=100 -BlockSize=200 -Commit=true -v -timeout 24h

sim-barkis-custom-genesis-fast:
	@echo "Running custom genesis simulation..."
	@echo "By default, ${HOME}/.barkisd/config/genesis.json will be used."
	@go test -mod=readonly $(SIMAPP) -run TestFullBarkisSimulation -Genesis=${HOME}/.barkisd/config/genesis.json \
		-Enabled=true -NumBlocks=100 -BlockSize=200 -Commit=true -Seed=99 -Period=5 -v -timeout 24h

sim-barkis-fast:
	@echo "Running quick Barkis simulation. This may take several minutes..."
	@go test -mod=readonly $(SIMAPP) -run TestFullBarkisSimulation -Enabled=true -NumBlocks=100 -BlockSize=200 -Commit=true -Seed=99 -Period=5 -v -timeout 24h

sim-barkis-import-export: runsim
	@echo "Running Barkis import/export simulation. This may take several minutes..."
	$(GOPATH)/bin/runsim $(SIMAPP) 25 5 TestBarkisImportExport

sim-barkis-simulation-after-import: runsim
	@echo "Running Barkis simulation-after-import. This may take several minutes..."
	$(GOPATH)/bin/runsim $(SIMAPP) 25 5 TestBarkisSimulationAfterImport

sim-barkis-custom-genesis-multi-seed: runsim
	@echo "Running multi-seed custom genesis simulation..."
	@echo "By default, ${HOME}/.barkisd/config/genesis.json will be used."
	$(GOPATH)/bin/runsim $(SIMAPP) -g ${HOME}/.barkisd/config/genesis.json 400 5 TestFullBarkisSimulation

sim-barkis-multi-seed: runsim
	@echo "Running multi-seed Barkis simulation. This may take awhile!"
	$(GOPATH)/bin/runsim $(SIMAPP) 400 5 TestFullBarkisSimulation

sim-benchmark-invariants:
	@echo "Running simulation invariant benchmarks..."
	@go test -mod=readonly $(SIMAPP) -benchmem -bench=BenchmarkInvariants -run=^$ \
	-Enabled=true -NumBlocks=1000 -BlockSize=200 \
	-Commit=true -Seed=57 -v -timeout 24h

SIM_NUM_BLOCKS ?= 500
SIM_BLOCK_SIZE ?= 200
SIM_COMMIT ?= true
sim-barkis-benchmark:
	@echo "Running Barkis benchmark for numBlocks=$(SIM_NUM_BLOCKS), blockSize=$(SIM_BLOCK_SIZE). This may take awhile!"
	@go test -mod=readonly -benchmem -run=^$$ $(SIMAPP) -bench ^BenchmarkFullBarkisSimulation$$  \
		-Enabled=true -NumBlocks=$(SIM_NUM_BLOCKS) -BlockSize=$(SIM_BLOCK_SIZE) -Commit=$(SIM_COMMIT) -timeout 24h

sim-barkis-profile:
	@echo "Running Barkis benchmark for numBlocks=$(SIM_NUM_BLOCKS), blockSize=$(SIM_BLOCK_SIZE). This may take awhile!"
	@go test -mod=readonly -benchmem -run=^$$ $(SIMAPP) -bench ^BenchmarkFullBarkisSimulation$$ \
		-Enabled=true -NumBlocks=$(SIM_NUM_BLOCKS) -BlockSize=$(SIM_BLOCK_SIZE) -Commit=$(SIM_COMMIT) -timeout 24h -cpuprofile cpu.out -memprofile mem.out


.PHONY: runsim sim-barkis-nondeterminism sim-barkis-custom-genesis-fast sim-barkis-fast sim-barkis-import-export \
	sim-barkis-simulation-after-import sim-barkis-custom-genesis-multi-seed sim-barkis-multi-seed \
	sim-benchmark-invariants sim-barkis-benchmark sim-barkis-profile
