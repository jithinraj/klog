package mac_widget

import (
	"fmt"
	"klog"
	"klog/app"
	"klog/lib/caseymrm/menuet"
	"klog/service"
	"time"
)

var blinker = blinkerT{1}

func render(ctx *app.Context, agent *launchAgent) []menuet.MenuItem {
	var items []menuet.MenuItem

	items = append(items, func() []menuet.MenuItem {
		if rs, file, err := ctx.Bookmark(); err == nil {
			return renderRecords(ctx, rs, file)
		}
		return []menuet.MenuItem{{
			Text:       "No input file specified",
			FontWeight: menuet.WeightBold,
		}, {
			Text: "Specify one by running:",
		}, {
			Text: "klog widget --file yourfile.klg",
		}}
	}()...)

	items = append(items, menuet.MenuItem{
		Type: menuet.Separator,
	}, menuet.MenuItem{
		Text: "Settings",
		Children: func() []menuet.MenuItem {
			return []menuet.MenuItem{{
				Text:  "Launch widget on login",
				State: agent.isActive(),
				Clicked: func() {
					var err error
					if agent.isActive() {
						err = agent.deactivate()
					} else {
						err = agent.activate()
					}
					if err != nil {
						fmt.Println(err)
					}
				},
			}}
		},
	})

	return items
}

func renderRecords(ctx *app.Context, records []klog.Record, file app.File) []menuet.MenuItem {
	var items []menuet.MenuItem

	now := time.Now()
	today := service.FindFilter(records, service.Filter{Dates: []klog.Date{klog.NewDateFromTime(now)}})
	if today != nil {
		total, isOngoing := service.HypotheticalTotal(now, today...)
		indicator := ""
		if isOngoing {
			indicator = "  " + blinker.blink()
		}
		items = append(items, menuet.MenuItem{
			Text: "Today: " + total.ToString() + indicator,
		})
	}

	items = append(items, menuet.MenuItem{
		Text: file.Name,
		Children: func() []menuet.MenuItem {
			total := service.Total(records...)
			should := service.ShouldTotalSum(records...)
			diff := klog.NewDuration(0, 0).Minus(should).Plus(total)
			plus := ""
			if diff.InMinutes() > 0 {
				plus = "+"
			}
			items := []menuet.MenuItem{
				{
					Text: "Show in Finder...",
					Clicked: func() {
						_ = ctx.OpenInFileBrowser(file.Location)
					},
				},
				{Type: menuet.Separator},
				{Text: "Total: " + total.ToString()},
				{Text: "Should: " + should.ToString()},
				{Text: "Diff: " + plus + diff.ToString()},
			}
			showMax := 7
			for i, r := range service.Sort(records, true) {
				if i == 0 {
					items = append(items, menuet.MenuItem{Type: menuet.Separator})
				}
				if i == showMax {
					items = append(items, menuet.MenuItem{Text: fmt.Sprintf("(%d more)", len(records)-showMax)})
					break
				}
				items = append(items, menuet.MenuItem{Text: r.Date().ToString() + ": " + service.Total(r).ToString()})
			}
			return items
		},
	})

	return items
}

type blinkerT struct {
	cycle int
}

func (b *blinkerT) blink() string {
	b.cycle++
	switch b.cycle {
	case 1:
		return "◷"
	case 2:
		return "◶"
	case 3:
		return "◵"
	}
	b.cycle = 0
	return "◴"
}
