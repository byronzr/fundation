package gbar

import (
	"fmt"
	"os"
	"time"
)

var (
	maxChannell   = 1024
	ChProgressBar = make(chan ProgressBarData, maxChannell)

	max    = int64(1)
	mline  = make(map[string]int64, 0)
	probe  = make(map[string]int64, 0)
	stimes = make(map[string]*time.Time)
	counts = make(map[string]int64, 0)
)

type ProgressBarData struct {
	Name  string
	Step  int64
	Msg   string
	Count bool
	Probe bool
	Time  bool
}

func init() {

	go func() {
		max := int64(1)
		mline := make(map[string]int64, 0)
		probe := make(map[string]int64, 0)
		stimes := make(map[string]*time.Time)
		counts := make(map[string]int64, 0)
		fmt.Printf("\033[1J\033[0;0H\033[0m\n      // Fundation Gbar Paneld //\n\n")
		for d := range ChProgressBar {

			var cbar []rune
			var info, color string
			step := int64(float64(d.Step) / 100 * 50)
			t, ok := stimes[d.Name]
			if !ok {
				nt := time.Now()
				stimes[d.Name] = &nt
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
				//cbar = []rune("┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈")
				//cbar = []byte("..................................................")
				// RUN
				if d.Msg == "" {
					color = "\033[33m"
					info = fmt.Sprintf("%s RUN\033[0m", color)

				} else {

					if d.Name == "ERR" {
						// ERR
						color = "\033[31m"
						info = fmt.Sprintf("%s ERR\033[0m", color)
					} else {
						// INF
						color = "\033[35m"
						info = fmt.Sprintf("%s INF\033[0m", color)
					}
				}
				if p, ok := probe[d.Name]; !ok {
					probe[d.Name] = 0
					//cbar[0] = '+'
					cbar[0] = point
				} else {
					n := int64(0)
					if p < 49 {
						n = p + 1
					}
					//cbar[n] = '+'
					cbar[n] = point
					probe[d.Name] = n
				}
			} else if d.Count {
				// 统计模式
				color = "\033[36m"
				info = fmt.Sprintf("%sCALC\033[0m", color)
				counts[d.Name]++
				cbar = []rune(fmt.Sprintf("% 50d", counts[d.Name]))
			} else {
				// 进度模式
				color = "\033[32m"
				info = fmt.Sprintf("%s%3d%%\033[0m", color, d.Step)
				// cbar = []byte("                                                  ")
				// if step < 50 {
				// 	cbar[step] = '>'
				// }
				// for i := int64(0); i < step; i++ {
				// 	cbar[i] = '='
				// }
				cbar = []rune("\033[32m")
				for i := int64(0); i < 50; i++ {
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
			line, ok := mline[d.Name]
			x := int64(0)
			if !ok {
				mline[d.Name] = max
				fmt.Printf(" %s", info)
				max++
			} else {
				x = max - line
				fmt.Printf("\033[%dF %s", x, info)
			}

			// add bar
			fmt.Printf(" [ %s ]", string(cbar))

			// add name
			fmt.Printf(" %s%-20s", color, d.Name)

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
				fmt.Printf(" %.0f%s\033[K", c, u[i])
			}

			// add msg
			if d.Msg != "" {
				fmt.Printf(" %s%-20s", color, d.Msg)
			}

			// end close ctrl
			fmt.Printf("\033[0m\033[%dE", x)

			// clear set
			if d.Step >= 100 {
				// delete(mline, d.Name)
				// delete(probe, d.Name)
				// delete(stimes, d.Name)
				// runtime.GC()
			}

			// exit
			if d.Name == "ERR" {
				os.Exit(1)
			}
		}
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

func Progress(name string, step int64, msg string) bool {
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

func Err(msg string) {
	d := ProgressBarData{
		Name: "ERR",
		Time: true,
		Msg:  msg,
	}
	send(d)

}

func send(d ProgressBarData) {
	if len(ChProgressBar) > maxChannell {
		panic("maxChannel out range!")
	}
	ChProgressBar <- d
}
