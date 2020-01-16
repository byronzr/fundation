package gbar

import (
	"fmt"
	"runtime"
	"time"
)

var (
	ChProgressBar = make(chan ProgressBarData, 1024)
)

type ProgressBarData struct {
	Name  string
	Step  int64
	Msg   string
	Probe bool
	Time  bool
}

func init() {
	go func() {
		max := int64(1)
		mline := make(map[string]int64, 0)
		probe := make(map[string]int64, 0)
		stimes := make(map[string]*time.Time)
		fmt.Printf("\033[1J\033[0;0H\033[0m\n      // STATUS PANEL INFO //\n\n")
		for d := range ChProgressBar {

			var cbar []byte
			var info, color string
			step := int64(float64(d.Step) / 100 * 50)
			t, ok := stimes[d.Name]
			if !ok {
				nt := time.Now()
				stimes[d.Name] = &nt
				t = &nt
			}

			if d.Probe { // 探针模式
				cbar = []byte("..................................................")
				color = "\033[33m"
				if d.Msg == "" {
					info = fmt.Sprintf("%s run\033[0m", color)
				} else {
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
			} else { // 进度模式
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
				fmt.Printf(" %s", info)
				max++
			} else {
				x = max - line
				fmt.Printf("\033[%dF %s", x, info)
			}

			// add bar
			fmt.Printf(" [ %s ]", cbar)

			// add name
			fmt.Printf(" %s%-20s", color, d.Name)

			// add msg
			if d.Msg != "" {
				fmt.Printf(" %s%-20s", color, d.Msg)
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
				fmt.Printf(" %.0f%s\033[K", c, u[i])
			}

			// end close ctrl
			fmt.Printf("\033[0m\033[%dE", x)

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
	ChProgressBar <- d
}
func Progress(name string, step int64) bool {
	d := ProgressBarData{
		Name: name,
		Step: step,
		Time: true,
	}
	ChProgressBar <- d
	if step >= 100 {
		return true
	}
	return false
}
func Info(name, message string) {
	d := ProgressBarData{
		Name:  name,
		Probe: true,
		Msg:   message,
	}
	ChProgressBar <- d
}

// func main() {
// 	names := []string{
// 		"MySQL Write",
// 		"MySQL Read",
// 		"Redis Write",
// 		"Redis Read",
// 		"Channal Push",
// 		"File I/O",
// 		"Logger Write",
// 		"HTTP GET",
// 		"Database Align",
// 		"Analysis Status",
// 	}
// 	for i := int64(1); i <= 9; i++ {
// 		go func(i int64) {
// 			step := int64(0)
// 			for {
// 				step++
// 				r := rand.New(rand.NewSource(time.Now().UnixNano()))
// 				time.Sleep(time.Duration(r.Int63n(int64(time.Second))))
// 				if i%2 == 1 {
// 					//Status(names[i])
// 					Info(names[i], ".....")
// 				} else {
// 					if Progress(names[i], step) {
// 						return
// 					}
// 				}
// 			}
// 		}(i)
// 	}
// 	s := make(chan int)
// 	<-s
// }
