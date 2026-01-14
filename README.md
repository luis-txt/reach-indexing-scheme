# Fast Reachability Using Indexed Transitivity (FRUIT)

## Description
This repository contains an implementation of the reachability indexing scheme from the paper "Fast Reachability Using DAG Decomposition" by Giorgos Kritikakis and Ioannis G. Tollis from the University of Crete.

Giorgos Kritikakis and Ioannis G. Tollis. Fast Reachability Using DAG Decomposition. In 21st International Symposium on Experimental Algorithms (SEA 2023). Leibniz International Proceedings in Informatics (LIPIcs), Volume 265, pp. 2:1-2:17, Schloss Dagstuhl – Leibniz-Zentrum für Informatik (2023) [https://doi.org/10.4230/LIPIcs.SEA.2023.2](https://doi.org/10.4230/LIPIcs.SEA.2023.2)


The appraoch includes the following steps: 
- Transforming a given directed graph to a directed acyclic graph (DAG) using Tarjan's strongly connected components algorithm
- Creating a chain decomposition of the calculated DAG
- Reducing the number of transitive edges of the DAG using the described heuristic from the paper
- Creating the indexing scheme as described in the paper
- Carrying out reachability queries as described in the paper

This reposetory also includes a script to visualize the graphs and the process (see Visuals). Additionally we provide a more basic approach (DFS based) to answer reachability queries and verify the results from the approach of the paper.

## Installation
This project requires Go (version 1.23) and the optional scripts require Python (version 3.12.5).

To install the project, follow these steps:
1. Clone the repository.
2. Build the project using the `go build` command.
3. Install required packages for the provided scripts (optional):
```
pip install numpy matplotlib pandas networkx
```

## Usage
### Indexing Scheme implementation
#### Running the Project
To run the project, you can either:
- Execute the code directly using: `go run fruit [flags] [input-graph-file-path]`
- Or, after building the project, execute the binary with: `./fruit [flags] [input-graph-file-path]`

#### Flags
You can customize the behavior of the implementation by using the following flags:
- -v: Enables verbose mode for detailed algorithm output.
- -m: Outputs the transitive closure matrix.
- -b: Measures and returns various performance metrics.
- -no: Uses the Node-Order (NO) heuristic for chain decomposition.
- -co: Uses the Chain-Order (CO) heuristic.
- -noc: Uses the NO heuristic followed by the concatenation (CONC) heuristic.
- -coc: Uses the CO heuristic followed by the concatenation (CONC) heuristic.
If no chain decomposition flag is set, the default heuristic is H3-Concat.

#### Unit Tests
Execute the unit-tests using: `go test fruit`. Use the *-v* flag to get detailed information.

### Additional Scripts
For benchmarking, plotting and creating tables, the following scripts are provided:
- benchmark_fruit.sh: This bash script handles the entire benchmark process.
- generate_graph.py: Generates random digraphs for testing.
- create_test_data.py: Creates a set of test graphs.
- run_benches.py: Runs benchmarks and measures performance.
- create_plots.py: Generates plots from the benchmark data.
- create_table.py: Produces tables from the benchmarking results.

#### Visuals
Benchmarking results can be visualized using the create_plots.py script. The following types of plots are supported:
- Comparison of decomposition methods.
- Visualization of the effects of different heuristics as the number of vertices and edges increases in the $G_{n,m}$ model.

To generate plots, run: `python3 scripts/create_plots.py`.
Tables can also be generated using: 
```
python3 scripts/create_table.py [sub-directory] [field]
```
where *[sub-directory]* refers to a sub-directory of the *benches* directory while *field* is a field of the generated benchmarking data.
