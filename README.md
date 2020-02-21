# mapslicecomp
This is just a collection of functions and benchmarks I've been using to compare the performance of a map and a slice for common use-cases.
I frequently heard at go meetups and related conferences that it is more efficient to just iterate through a slice than to retrieve the value from a map when dealing with a relatively small number of values. 

Currently the main purpose of this repo is to demonstrate the usage of my [gobenchplot](https://github.com/ShawnROGrady/gobenchplot) tool, which I developed to visualize the results of the benchmarks.

I want to come up with some more variations and use cases before claiming any results, but at least on my computer with the current benchmarks the map and slice implementations have practically equivalent performance for a small number of values (think `num_elems<10`) but after that the map implementation is far better. 

Below are the results of the benchmark for deduping a slice of strings which takes in to account both insertion and traversal. This was generated as:
```
go test . -run ! -bench BenchmarkDedupe -benchmem -json -timeout 0 -benchtime 10000x -count 3 | tee tmp.txt | gobenchplot --bench='BenchmarkDedupe' --x='num_elems' --group-by='finder'
```
![bench_dedupe](https://github.com/ShawnROGrady/mapslicecomp/blob/master/assets/benchmark_dedupe_time-v-num_elems.png)

Focusing on just the portion of interest:
```
gobenchplot --bench='BenchmarkDedupe' --x='num_elems' --group-by='finder' --filter-by='num_elems<=10' tmp.txt
```
![focused_bench_dedupe](https://github.com/ShawnROGrady/mapslicecomp/blob/master/assets/focused_benchmark_dedupe_time-v-num_elems.png)
