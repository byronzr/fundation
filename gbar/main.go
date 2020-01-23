package gbar

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	maxChannell   = 1024
	ChProgressBar = make(chan ProgressBarData, maxChannell)
	WG            = sync.WaitGroup{}

	max    = 1
	mline  = make(map[string]int64, 0)
	probe  = make(map[string]int64, 0)
	stimes = make(map[string]*time.Time)
	counts = make(map[string]int64, 0)
)

type ProgressBarData struct {
	Name  string
	Step  int
	Msg   string
	Count bool
	Probe bool
	Time  bool
}

func init() {

	go func() {
		WG.Add(1)
		max := int64(1)
		mline := make(map[string]int64, 0)
		probe := make(map[string]int64, 0)
		stimes := make(map[string]*time.Time)
		counts := make(map[string]int64, 0)
		fmt.Printf("\033[1J\033[0;0H\033[0m\n      // Fundation Gbar Paneld //\n\n")
		for d := range ChProgressBar {

			name := strings.ToLower(d.Name)

			var cbar, buf []rune
			var info, color string
			var step int
			step = int(float64(d.Step) / 100 * 50)
			t, ok := stimes[name]
			if !ok {
				nt := time.Now()
				stimes[name] = &nt
				t = &nt
			}

			if d.Probe {
				// 探针模式
				// ──────────────────────────────────────────────────┬┴┰┸
				// ┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄
				// ┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈
				// ──────────────────────────────────────────────────
				// ❚❚❚❚❚❚❚❚❚❚❚❚❚❚❚❚❚❚❚❚❚❚❚❚❚❚❚❚❚❚❚❚❚❚❚❚❚❚❚❚❚❚❚❚❚❚❚❚❚❚
				point := '┸'
				cbar = []rune("──────────────────────────────────────────────────")
				// RUN
				if d.Msg == "" {
					color = "\033[33m"
					info = fmt.Sprintf("%s RUN\033[0m", color)

				} else {
					// INF
					color = "\033[35m"
					info = fmt.Sprintf("%s INF\033[0m", color)
				}
				if p, ok := probe[name]; !ok {
					probe[name] = 0
					//cbar[0] = '+'
					cbar[0] = point
				} else {
					n := int64(0)
					if p < 49 {
						n = p + 1
					}
					//cbar[n] = '+'
					cbar[n] = point
					probe[name] = n
				}
			} else if d.Count {
				// 统计模式
				color = "\033[36m"
				info = fmt.Sprintf("%sCALC\033[0m", color)
				counts[name]++
				cbar = []rune(fmt.Sprintf("% 50d", counts[name]))
			} else {
				// 进度模式
				if d.Step == 0 {
					color = "\033[31m"
					info = fmt.Sprintf("%s ERR\033[0m", color)
				} else {
					color = "\033[32m"
					info = fmt.Sprintf("%s%3d%%\033[0m", color, d.Step)
				}

				cbar = []rune("\033[32m")
				for i := 0; i < 50; i++ {
					if i == step+1 {
						cbar = append(cbar, []rune("❚\033[0m")...)
					} else {
						if i <= step {
							cbar = append(cbar, '❚')
						} else {
							cbar = append(cbar, ' ')
						}
					}
				}
				cbar = append(cbar, []rune("\033[0m")...)
			}

			// Add status
			line, ok := mline[name]
			x := int64(0)
			if !ok {
				mline[name] = max
				buf = []rune(fmt.Sprintf(" %s", info))
				max++
			} else {
				x = max - line
				buf = []rune(fmt.Sprintf("\033[%dF %s", x, info))
			}

			// add bar
			buf = append(buf, []rune(fmt.Sprintf(" [ %s ]", string(cbar)))...)

			// add name
			buf = append(buf, []rune(fmt.Sprintf(" %s%-20s", color, name))...)

			// add time
			if d.Time {
				c := time.Since(*t).Seconds()
				u := []string{"s", "m", "h", "d", "month", "year"}
				d := []float64{60, 60, 24, 30, 365}
				i := 0
				for {
					if c > d[i] {
						c = c / d[i]
						i++
					} else {
						break
					}
				}
				buf = append(buf, []rune(fmt.Sprintf(" %3.0f%s\033[K", c, u[i]))...)
			}

			// add msg
			if d.Msg != "" {
				buf = append(buf, []rune(fmt.Sprintf("  : %s%-20s", color, d.Msg))...)
			}

			// end close ctrl
			buf = append(buf, []rune(fmt.Sprintf("\033[%dE\033[0m", x))...)

			os.Stdout.WriteString(string(buf))

			// clear set
			if d.Step >= 100 {
				// delete(mline, name)
				// delete(probe, name)
				// delete(stimes, name)
				// runtime.GC()
			}

		}
		WG.Done()
		os.Stdout.Sync()
	}()
}

func Status(name string) {
	d := ProgressBarData{
		Name:  name,
		Probe: true,
		Time:  true,
	}
	send(d)
}

func Progress(name string, step int, msg string) bool {
	d := ProgressBarData{
		Name: name,
		Step: step,
		Time: true,
		Msg:  msg,
	}
	if step > 100 {
		return true
	}
	send(d)
	return false
}

func Info(name, message string) {
	d := ProgressBarData{
		Name:  name,
		Probe: true,
		Msg:   message,
	}
	send(d)
}

func Count(name string) {
	d := ProgressBarData{
		Name:  name,
		Time:  true,
		Count: true,
	}
	send(d)
}

func Fatal() {
	close(ChProgressBar)
	WG.Wait()
	os.Exit(0)
}

func send(d ProgressBarData) {
	if len(ChProgressBar) > maxChannell {
		panic("maxChannel out range!")
	}
	ChProgressBar <- d
}
