# TCTF2020-QUAL (no)easyphp

Sorry for unintended solution.

Block `cdef` and similar function, how can you call user function through FFI extension? We can learn that `FFI::load` will load the entire file into memory, and FFI has many funtions for memory operation, we can leak some information through this.

Use `$p = FFI::new("char");` to get a pointer and `FFI::addr($p)-$i` to dump the memory and you can find the function defination.