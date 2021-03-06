package gbar

import (
	"fmt"
	"runtime"
	"time"
)

var (
	maxChannell   = 1024
	ChProgressBar = make(chan ProgressBarData, maxChannell)
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
		fmt.Printf("\033[1J\033[0;0H\033[0m\n      // Fundation Gbar Panel string //\n\n")
		for d := range ChProgressBar {
			var cbar []byte
			var info, color string
			buf := make([]byte, 0)
			step := int64(float64(d.Step) / 100 * 50)
			t, ok := stimes[d.Name]
			if !ok {
				nt := time.Now()
				stimes[d.Name] = &nt
				t = &nt
			}

			if d.Probe {
				// 探针模式
				cbar = []byte("..................................................")
				if d.Msg == "" {
					color = "\033[33m"
					info = fmt.Sprintf("%s run\033[0m", color)
				} else {
					color = "\033[35m"
					info = fmt.Sprintf("%s inf\033[0m", color)
				}
				if p, ok := probe[d.Name]; !ok {
					probe[d.Name] = 0
					cbar[0] = '+'
				} else {
					n := int64(0)
					if p < 49 {
						n = p + 1
					}
					cbar[n] = '+'
					probe[d.Name] = n
				}
			} else if d.Count {
				// 统计模式
				color = "\033[36m"
				info = fmt.Sprintf("%sCalc\033[0m", color)
				counts[d.Name]++
				cbar = []byte(fmt.Sprintf("% 50d", counts[d.Name]))
			} else {
				// 进度模式
				color = "\033[32m"
				info = fmt.Sprintf("%s%3d%%\033[0m", color, d.Step)
				cbar = []byte("                                                  ")
				if step < 50 {
					cbar[step] = '>'
				}
				for i := int64(0); i < step; i++ {
					cbar[i] = '='
				}
			}

			// add status
			line, ok := mline[d.Name]
			x := int64(0)
			if !ok {
				mline[d.Name] = max
				buf = append(buf, []byte(fmt.Sprintf(" %s", info))...)
				max++
			} else {
				x = max - line
				buf = append(buf, []byte(fmt.Sprintf("\033[%dF %s", x, info))...)
			}

			// add bar
			buf = append(buf, []byte(fmt.Sprintf(" [ %s ]", cbar))...)

			// add name
			buf = append(buf, []byte(fmt.Sprintf(" %s%-20s", color, d.Name))...)

			// add msg
			if d.Msg != "" {
				buf = append(buf, []byte(fmt.Sprintf(" %s%-20s", color, d.Msg))...)
			}

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
				buf = append(buf, []byte(fmt.Sprintf(" %.0f%s\033[K", c, u[i]))...)
			}

			// end close ctrl
			buf = append(buf, []byte(fmt.Sprintf("\033[0m\033[%dE", x))...)
			fmt.Print(string(buf))

			// clear set
			if d.Step >= 100 {
				delete(mline, d.Name)
				delete(probe, d.Name)
				delete(stimes, d.Name)
				runtime.GC()
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
func Progress(name string, step int64) bool {
	d := ProgressBarData{
		Name: name,
		Step: step,
		Time: true,
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

func send(d ProgressBarData) {
	if len(ChProgressBar) > maxChannell {
		panic("maxChannel out range!")
	}
	ChProgressBar <- d
}
