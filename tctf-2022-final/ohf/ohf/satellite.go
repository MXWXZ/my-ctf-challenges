import "strings"

func plugin() string {
	var ret []string
	ret = append(ret, "Hacking satellite...")
	ret = append(ret, "Success! Getting shell...")
	ret = append(ret, "Microsoft Windows [Version 10.0.22000.918]")
	ret = append(ret, "(c) Microsoft Corporation. All rights reserved.")
	ret = append(ret, "")
	ret = append(ret, `C:\Windows\system32>whoami`)
	ret = append(ret, `nt authority\system`)
	return strings.Join(ret, "\n")
}
