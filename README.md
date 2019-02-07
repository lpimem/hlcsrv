# hlcsrv

HLC Server


## Benchmarks

### 04/20/2017

```
BenchmarkSavePagenote-8    	   30000	    525521 ns/op	   14174 B/op	     113 allocs/op
BenchmarkSavePagenoteP-8   	   20000	   1002649 ns/op	   14176 B/op	     113 allocs/op
BenchmarkGetPagenote/10-8  	  100000	    198909 ns/op	   16058 B/op	     562 allocs/op
BenchmarkGetPagenote/100-8 	   20000	    776848 ns/op	   72858 B/op	    3540 allocs/op
BenchmarkGetPagenote/1000-8     2000	   6432966 ns/op	  690328 B/op	   33258 allocs/op
BenchmarkGetPagenoteP-8        50000	    293017 ns/op	   73156 B/op	    3543 allocs/op
```