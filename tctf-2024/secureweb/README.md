# Writeup

Rust is secure, but we can still pwn it :)

## Step1: SQL injection
`api.rs:178` we can escape the quote here. However, we need to bypass checkers in `check_sql`.

- `let blacklist = vec![";", "INSERT", "SELECT", "UNION", "OR", "\\", " "];`: easy bypass through the lowercase and `/**/`.
- `r"\w+(\(.*\)).*"`: `.` does not include `\n`, add a `\n` after `(` can bypass it.
- `r"\([\s\S]*,[\s\S]*"`: well we cannot bypass (as far as I know).

So the problem is reduced to this: Get admin password from `users`
- In one `INSERT` statement column prefixed with `'`
- In 10 request (thus blind injection not possible)
- No error returned
- Can call any function with only one parameter

Unfortunately, due to the `'` in the `INSERT` column, we cannot select out the password directly.
Fortunately, the admin password is a UUID, meaning the characters are within `0-9a-f` (and `-` at fixed index).
The hex value of `0-9a-f` is `0x`+`30-39,61-66`, no `a-f` here! SQL also support string number operation:
```
MariaDB [ctf]> select '0'+'123';
+-----------+
| '0'+'123' |
+-----------+
|       123 |
+-----------+
```
And `from .. for` to avoid `,` in `substr`.
```
MariaDB [ctf]> select substr('abcdef' from 2 for 3);
+-------------------------------+
| substr('abcdef' from 2 for 3) |
+-------------------------------+
| bcd                           |
+-------------------------------+
```

So we can get the password!
```
"0'+(select hex(\nsubstr(\npassword from 1 for 4)) from users where id=1)) # --".replace(' ', '/**/')
```
Note: To avoid scientific notation, we select 4 characters each time.

## Step2: XFF bypass
The proxy splits the input by lines. Then append `X-Forwarded-For` after the first line, together with other lines using `\n`. Finally, replace all `\r` or `\n` to `\r\n`.

`request.rs:86`: `req.headers().get("x-forwarded-for")` from the comment we can know `headers().get` only takes the first header. Thus, our aim is to insert another XFF header before the first `line`.

`str::Lines` comment says it reserve `\r` if `\n` is not immediately followed. So our solution is to add `\rx-forwarded-for:127.0.0.1` between `HTTP/1.1` and `\r\n`.

Note: smuggling is also possible, but I do not want to change the code anymore...

## Step3: UAF
Using `rhai` script, we can define the folder and files to read. However, if the flag is in the output, it only shows the path.

There is a pinned issue: https://github.com/rhaiscript/rhai/issues/894

The problem can be reduced as:
- `malloc` any bytes any times
- `free` all chunks after script executed
- Only get reference to the **last** freed string
- `\r\n\0 \t` is filtered for last string, and used as the file path
- `malloc` `sizeof(file content)` any times (`api.rs:214`)

Our aim is to use `CONFIG_FOLDER` as the return value of `malloc` when we read the flag, so we can leak the flag through the returned `CONFIG_FOLDER`.

GLIBC 2.31 has tcache, the first 16 bytes are used as pointers (which is not valid path characters). However, when the tcache is full, the first fastbin chunk has 8 bytes `\0` that can be filtered.

Flag `0ctf{<uuid>}` is 42 bytes (0x40 chunk). So our solution is:
- `malloc` and `free` x chunks (0x40) (x-1 goes to tcache and 1 goes to fastbin).
- Read flag for x times (x-1 using tcache and 1 use fastbin).
- Get flag from response!

Note: There are many other methods that can consume tcache or arrange the heap.

Note: You can debug the server to determine how many times you need. I got 6 but is not stable, try multiple times if possible.

Note: You need to add more characters to the file list, otherwise the constructed path also need extra chunks.