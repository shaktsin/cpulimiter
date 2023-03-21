# GO CPU Limiter 
This is a proof-of-concept - process cpu usage limiter. It is not production ready, it is only meant for understadning the concepts. This only works on Linux. 

# Compilation
```
make build
```

# Usage 

```
./cpulimit -p=<target_process_pid> -l=<target avg cpu usgae across all cores (float valueb/w 0-1, 0.05 = 5%)>
```

# Test






