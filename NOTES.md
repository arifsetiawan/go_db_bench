
# Run Postgres

```
docker run --name postgres-bench -p 5432:5432 -e POSTGRES_USER=one -e POSTGRES_PASSWORD=two -e POSTGRES_DB=bench -d postgres:10.1
```

# benchmark

```
goos: darwin
goarch: amd64
pkg: gitlab.com/github.com/arifsetiawan/go_db_bench
BenchmarkPgxNativeSelectSingleShortString-2       	    5000	    366872 ns/op	     292 B/op	       6 allocs/op
BenchmarkPgxStdlibSelectSingleShortString-2       	    3000	    388180 ns/op	     932 B/op	      24 allocs/op
BenchmarkPgSelectSingleShortString-2              	    3000	    338188 ns/op	     120 B/op	       8 allocs/op
BenchmarkPqSelectSingleShortString-2              	    3000	    351577 ns/op	     736 B/op	      24 allocs/op
BenchmarkRawSelectSingleShortValue-2              	    5000	    322354 ns/op	     121 B/op	       4 allocs/op
BenchmarkPgxNativeSelectSingleShortBytes-2        	    3000	    338826 ns/op	     313 B/op	       7 allocs/op
BenchmarkPgxStdlibSelectSingleShortBytes-2        	    3000	    367871 ns/op	     950 B/op	      25 allocs/op
BenchmarkPqSelectSingleShortBytes-2               	    3000	    509873 ns/op	     760 B/op	      25 allocs/op
BenchmarkPgxNativeSelectSingleRow-2               	    3000	    389751 ns/op	     479 B/op	       8 allocs/op
BenchmarkPgxStdlibSelectSingleRow-2               	    3000	    499946 ns/op	    1724 B/op	      41 allocs/op
BenchmarkPgSelectSingleRow-2                      	    3000	    461667 ns/op	     380 B/op	      16 allocs/op
BenchmarkPqSelectSingleRow-2                      	    3000	    463889 ns/op	    1135 B/op	      36 allocs/op
BenchmarkRawSelectSingleRow-2                     	    5000	    340203 ns/op	     121 B/op	       4 allocs/op
BenchmarkPgxNativeSelectMultipleRows-2            	    3000	    502369 ns/op	    2942 B/op	      80 allocs/op
BenchmarkPgxStdlibSelectMultipleRows-2            	    2000	    601745 ns/op	    7604 B/op	     346 allocs/op
BenchmarkPgSelectMultipleRowsCollect-2            	    2000	    649171 ns/op	   13101 B/op	     284 allocs/op
BenchmarkPgSelectMultipleRowsAndDiscard-2         	    2000	    520421 ns/op	      16 B/op	       3 allocs/op
BenchmarkPqSelectMultipleRows-2                   	    2000	    661722 ns/op	    6171 B/op	     389 allocs/op
BenchmarkRawSelectMultipleRows-2                  	    2000	    686293 ns/op	     302 B/op	      10 allocs/op
BenchmarkPgxNativeSelectMultipleRowsBytes-2       	    3000	    491083 ns/op	    3470 B/op	     155 allocs/op
BenchmarkPgxStdlibSelectMultipleRowsBytes-2       	    2000	    597577 ns/op	    8384 B/op	     421 allocs/op
BenchmarkPqSelectMultipleRowsBytes-2              	    2000	    659583 ns/op	    7050 B/op	     464 allocs/op
BenchmarkPgxNativeSelectBatch3Query-2             	    2000	    833879 ns/op	    1744 B/op	      34 allocs/op
BenchmarkPgxNativeSelectNoBatch3Query-2           	    1000	   1014867 ns/op	     800 B/op	      13 allocs/op
BenchmarkPgxStdlibSelectNoBatch3Query-2           	    2000	   1043029 ns/op	    2668 B/op	      65 allocs/op
BenchmarkPqSelectNoBatch3Query-2                  	    2000	   1020652 ns/op	    2122 B/op	      65 allocs/op
BenchmarkPgxNativeSelectLargeTextString1KB-2      	    5000	    338517 ns/op	    2637 B/op	       6 allocs/op
BenchmarkPgxNativeSelectLargeTextString8KB-2      	    2000	    693436 ns/op	   22006 B/op	       8 allocs/op
BenchmarkPgxNativeSelectLargeTextString64KB-2     	     500	   2150516 ns/op	  143588 B/op	       8 allocs/op
BenchmarkPgxNativeSelectLargeTextString512KB-2    	     100	  20417637 ns/op	 1061166 B/op	       8 allocs/op
BenchmarkPgxNativeSelectLargeTextString4096KB-2   	      10	 137771594 ns/op	 8400955 B/op	       8 allocs/op
BenchmarkPgxStdlibSelectLargeTextString1KB-2      	    5000	    343882 ns/op	    3269 B/op	      24 allocs/op
BenchmarkPgxStdlibSelectLargeTextString8KB-2      	    2000	    703094 ns/op	   22638 B/op	      26 allocs/op
BenchmarkPgxStdlibSelectLargeTextString64KB-2     	     500	   2483891 ns/op	  144220 B/op	      26 allocs/op
BenchmarkPgxStdlibSelectLargeTextString512KB-2    	      50	  23181002 ns/op	 1061798 B/op	      26 allocs/op
BenchmarkPgxStdlibSelectLargeTextString4096KB-2   	      10	 149805175 ns/op	 8401588 B/op	      26 allocs/op
BenchmarkPgSelectLargeTextString1KB-2             	    5000	    330015 ns/op	    1144 B/op	       8 allocs/op
BenchmarkPgSelectLargeTextString8KB-2             	    2000	    686795 ns/op	   16504 B/op	       9 allocs/op
BenchmarkPgSelectLargeTextString64KB-2            	    1000	   2676390 ns/op	  131192 B/op	       9 allocs/op
BenchmarkPgSelectLargeTextString512KB-2           	      50	  23879377 ns/op	 1048696 B/op	       9 allocs/op
BenchmarkPgSelectLargeTextString4096KB-2          	      10	 147470792 ns/op	 8388729 B/op	       9 allocs/op
BenchmarkPqSelectLargeTextString1KB-2             	    5000	    341067 ns/op	    2904 B/op	      25 allocs/op
BenchmarkPqSelectLargeTextString8KB-2             	    2000	    661195 ns/op	   18392 B/op	      25 allocs/op
BenchmarkPqSelectLargeTextString64KB-2            	     500	   2286978 ns/op	  140024 B/op	      25 allocs/op
BenchmarkPqSelectLargeTextString512KB-2           	      50	  24569181 ns/op	 1057574 B/op	      25 allocs/op
BenchmarkPqSelectLargeTextString4096KB-2          	      10	 150756376 ns/op	 8397635 B/op	      26 allocs/op
BenchmarkPgxNativeSelectLargeTextBytes1KB-2       	    5000	    332326 ns/op	    3677 B/op	       7 allocs/op
BenchmarkPgxNativeSelectLargeTextBytes8KB-2       	    2000	    691018 ns/op	   30220 B/op	       9 allocs/op
BenchmarkPgxNativeSelectLargeTextBytes64KB-2      	     500	   2221073 ns/op	  209164 B/op	       9 allocs/op
BenchmarkPgxNativeSelectLargeTextBytes512KB-2     	     100	  14347953 ns/op	 1585347 B/op	       9 allocs/op
BenchmarkPgxNativeSelectLargeTextBytes4096KB-2    	      10	 136302146 ns/op	12596507 B/op	       9 allocs/op
BenchmarkPgxStdlibSelectLargeTextBytes1KB-2       	    3000	    359484 ns/op	    4309 B/op	      25 allocs/op
BenchmarkPgxStdlibSelectLargeTextBytes8KB-2       	    2000	    735925 ns/op	   30852 B/op	      27 allocs/op
BenchmarkPgxStdlibSelectLargeTextBytes64KB-2      	     500	   2257734 ns/op	  209796 B/op	      27 allocs/op
BenchmarkPgxStdlibSelectLargeTextBytes512KB-2     	     100	  14360314 ns/op	 1586102 B/op	      27 allocs/op
BenchmarkPgxStdlibSelectLargeTextBytes4096KB-2    	      10	 151544734 ns/op	12595908 B/op	      27 allocs/op
BenchmarkPqSelectLargeTextBytes1KB-2              	    5000	    345738 ns/op	    3944 B/op	      26 allocs/op
BenchmarkPqSelectLargeTextBytes8KB-2              	    2000	    679366 ns/op	   26603 B/op	      26 allocs/op
BenchmarkPqSelectLargeTextBytes64KB-2             	     500	   2254728 ns/op	  205591 B/op	      26 allocs/op
BenchmarkPqSelectLargeTextBytes512KB-2            	     100	  13964848 ns/op	 1581880 B/op	      26 allocs/op
BenchmarkPqSelectLargeTextBytes4096KB-2           	      10	 138223624 ns/op	12591958 B/op	      27 allocs/op
PASS
ok  	gitlab.com/github.com/arifsetiawan/go_db_bench	107.077s

```

# wrk

wrk -t20 -c100 -d30s http://localhost:8080/people/pgx-native

wrk -t20 -c100 -d30s http://localhost:8080/people/pgx-stdlib

wrk -t20 -c100 -d30s http://localhost:8080/people/pq

wrk -t20 -c100 -d30s http://localhost:8080/people/pg

```
InfiniteSquares:Documents arifsetiawan$ wrk -t20 -c100 -d30s http://localhost:8080/people/pg
Running 30s test @ http://localhost:8080/people/pg
  20 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   147.66ms   94.27ms 927.03ms   89.95%
    Req/Sec    40.39     15.35   101.00     72.61%
  21700 requests in 30.10s, 90.62MB read
Requests/sec:    720.90
Transfer/sec:      3.01MB

InfiniteSquares:Documents arifsetiawan$ wrk -t20 -c100 -d30s http://localhost:8080/people/pq
Running 30s test @ http://localhost:8080/people/pq
  20 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   208.35ms  291.22ms   2.00s    83.80%
    Req/Sec    31.77     30.42   213.00     84.88%
  12584 requests in 30.10s, 52.52MB read
  Socket errors: connect 0, read 0, write 0, timeout 115
  Non-2xx or 3xx responses: 6
Requests/sec:    418.02
Transfer/sec:      1.74MB

InfiniteSquares:Documents arifsetiawan$ wrk -t20 -c100 -d30s http://localhost:8080/people/pgx-native
Running 30s test @ http://localhost:8080/people/pgx-native
  20 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   112.88ms   66.80ms 810.99ms   84.77%
    Req/Sec    48.29     17.84   158.00     66.06%
  28117 requests in 30.10s, 117.40MB read
Requests/sec:    934.05
Transfer/sec:      3.90MB

InfiniteSquares:Documents arifsetiawan$ wrk -t20 -c100 -d30s http://localhost:8080/people/pgx-stdlib
Running 30s test @ http://localhost:8080/people/pgx-stdlib
  20 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   166.53ms  254.77ms   1.96s    84.11%
    Req/Sec    37.85     39.77   420.00     84.25%
  13307 requests in 30.09s, 55.53MB read
  Socket errors: connect 0, read 0, write 0, timeout 108
  Non-2xx or 3xx responses: 2
Requests/sec:    442.28
Transfer/sec:      1.85MB
```
