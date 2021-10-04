# TCTF2021-QUAL worldcup

Golang SSTI and xxs

## First layer

Inject point is in `nickname`, you can escape simply use a `'`, use meta refresh to lead the page to your site.

## Second layer

Because of some problem some unintended solution could work... Anyway, the Golang render engine will not parse ` correctly, so we can escape with it.

## Third layer

From [here](https://pkg.go.dev/text/template) you can see that `{{.}}` could print all the variables, however it's filtered. Instead you can use `{{$}}` to print all the variables(thus you can see the first and second number).

Tips is in the result test, the backend does just as the text says: calculate results, render the page and update database. Thus we need to find a way to break render when we make a wrong guess(so we could avoid punishment). You can use `{{if}}` keyword to check your guess and use `call xxx` that invoke an unexisting function or something similar to force render engine exit when you miss the guess. Thus you can always win and get the flag :)