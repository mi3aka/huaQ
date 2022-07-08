## /proc

`/proc/self`指向当前进程的`/proc/PID/`

`/proc/self/root/`是指向`/`的符号链接

如果可以任意文件读取,可以尝试对`/proc/PID/cmdline`进行读取,可能会读取到敏感信息

