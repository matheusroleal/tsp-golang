# Parallel Travelling Salesman Problem
Implementation of a Parallel Travelling Salesman Problem (TSP) by Branch and Bound Algorithm using Golang

## What it is ?

The travelling salesman problem (TSP) asks the following question: "Given a list of cities and the distances between each pair of cities, what is the shortest possible route that visits each city exactly once and returns to the origin city?" It is an NP-hard problem in combinatorial optimization, important in theoretical computer science and operations research.

In the theory of computational complexity, the decision version of the TSP (where given a length L, the task is to decide whether the graph has a tour of at most L) belongs to the class of NP-complete problems. Thus, it is possible that the worst-case running time for any algorithm for the TSP increases superpolynomially (but no more than exponentially) with the number of cities.

## Implementation
The program uses an affiliate and linked approach to resolve the issue. This means that each instance contains a path, a count of edges already passed, and a minimal cost of a round trip with that path at the beginning. This information is managed in a Priority Queue which is currently implemented as an array-based binary heap. We use Goroutines for parallel execution of threads. The idea is that we will simultaneously be able to reach the ideal solution more quickly. An instance is extended, generating all possible possibilities from the given path. We then add all available edges, recalculate the limits, and finally not be found back into the heap. To prevent race conditions from happening to a system queue, we use sync Mutex. This is a method used as a blocking mechanism to ensure that only one Goroutine accesses the critical section of code at a time.

One way we approach this problem is by using lower bounds. We calculate a value that represents a definite minimum length for each possible path that has a certain subpath. To do this, we start by calculating the lowest weight of an output edge for each vertex. Now we can add all of these together. We can do the same for the input edges. Since these two values ​​are definite lower bounds for the shortest path length, we only need to consider the higher value.
Now our first subpath. It contains exactly one vertex and no edges. This path can be extended by all neighbors of the last visited node (the only one so far). We now create at most n-1 new candidates. For all of them, we calculate the ultimate lower bound again. The only difference is that, at this point, we have a node where we know the output weight and a node where we know the input weight for sure.

As this decision can change the lower limit, we can now choose the candidate with the lowest minimum cost. If there are two candidates with the same threshold but different lengths, we choose the shorter one, since more effective nodes than the threshold is more likely to be the actual final value.

At some point, we'll end up with a “candidate” that has exactly n edges and the cost of which is at most equal to another candidate's lower bound. So there is no other shortest way.

## Usage

There is no actual CLI-Parser. You need to build the code before runing a scenario:
```
$ make setup
$ make build
```
After that, you just need to run the binary with the necessary arguments:
```
$ tsp <path> <matrix-size> <go-routines>
```