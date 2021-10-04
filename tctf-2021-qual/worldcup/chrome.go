package main

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
)

func xssrun(id int, level int) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if level == 1 {
		chrome(ctx, fmt.Sprintf("http://localhost:8080/wTf_Pa1h_5604m/%v/%v", id, level), "set cookie: level1 "+level1)
	} else {
		chrome(ctx, fmt.Sprintf("http://localhost:8080/wTf_Pa1h_5604m/%v/%v", id, level), "set cookie: level2 "+level2)
	}
}

func chrome(ctx context.Context, url string, token string) {
	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	ch := chromedp.WaitNewTarget(ctx, func(info *target.Info) bool {
		return info.URL != ""
	})
	chromedp.Run(ctx, emulation.SetUserAgentOverride(token), chromedp.Navigate(url))
	newCtx, cancel := chromedp.NewContext(ctx, chromedp.WithTargetID(<-ch))
	defer cancel()
	chromedp.Run(newCtx, emulation.SetUserAgentOverride(token))
}
