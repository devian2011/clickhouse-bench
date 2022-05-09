package froze

import (
	"time"
)

type Froze struct {
	begin time.Time
	end   time.Time
}

func (f *Froze) Start() {
	f.begin = time.Now()
}

func (f *Froze) Stop() {
	f.end = time.Now()
}

func (f *Froze) GetDiff() time.Duration {
	return f.end.Sub(f.begin)
}
