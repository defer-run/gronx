package gronx

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

var literals = strings.NewReplacer(
	"SUN", "0", "MON", "1", "TUE", "2", "WED", "3", "THU", "4", "FRI", "5", "SAT", "6",
	"JAN", "1", "FEB", "2", "MAR", "3", "APR", "4", "MAY", "5", "JUN", "6", "JUL", "7",
	"AUG", "8", "SEP", "9", "OCT", "10", "NOV", "11", "DEC", "12",
)

var expressions = map[string]string{
	"@yearly":    "0 0 1 1 *",
	"@annually":  "0 0 1 1 *",
	"@monthly":   "0 0 1 * *",
	"@weekly":    "0 0 * * 0",
	"@daily":     "0 0 * * *",
	"@hourly":    "0 * * * *",
	"@always":    "* * * * *",
	"@5minutes":  "*/5 * * * *",
	"@10minutes": "*/10 * * * *",
	"@15minutes": "*/15 * * * *",
	"@30minutes": "0,30 * * * *",
}

var spaceRe = regexp.MustCompile(`\s+`)

func normalize(expr string) []string {
	expr = strings.Trim(expr, " \t")
	if e, ok := expressions[strings.ToLower(expr)]; ok {
		expr = e
	}

	expr = spaceRe.ReplaceAllString(expr, " ")
	expr = literals.Replace(strings.ToUpper(expr))

	return strings.Split(strings.ReplaceAll(expr, "  ", " "), " ")
}

// Gronx is the main program.
type Gronx struct {
	C Checker
}

// New initializes Gronx with factory defaults.
func New() Gronx {
	return Gronx{&SegmentChecker{}}
}

// IsDue checks if cron expression is due for given reference time (or now).
// It returns bool or error if any.
func (g *Gronx) IsDue(expr string, ref ...time.Time) (bool, error) {
	segs := normalize(expr)
	if len(segs) < 5 || len(segs) > 6 {
		return false, errors.New("expr should contain 5 to 6 segments separated by space")
	}

	if len(ref) > 0 {
		g.C.SetRef(ref[0])
	} else {
		g.C.SetRef(time.Now())
	}

	for pos, seg := range segs {
		if seg == "*" || seg == "?" {
			continue
		}

		if due, err := g.C.CheckDue(seg, pos); !due {
			return due, err
		}
	}

	return true, nil
}
