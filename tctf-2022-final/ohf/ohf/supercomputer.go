import "strings"

func plugin() string {
	var ret []string
	ret = append(ret, "Hacking SkyNet...")
	ret = append(ret, "Success! Getting shell...")
	ret = append(ret, "Linux SkyNet 4.19.0-20-amd64 #1 SMP Debian 4.19.235-1 (2022-03-17) x86_64")
	ret = append(ret, "")
	ret = append(ret, "I'll be back.")
	ret = append(ret, "")
	ret = append(ret, `root@SkyNet:~# id`)
	ret = append(ret, `uid=0(root) gid=0(root) groups=0(root)`)
	return strings.Join(ret, "\n")
}
