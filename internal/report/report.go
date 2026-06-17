package report

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"

	"goupd/internal/resolver"
)

type Report struct {
	Module    string            `json:"module"`
	GoVersion string            `json:"goVersion"`
	Updates   []resolver.Update `json:"updates"`
}

func WriteJSON(w io.Writer, r Report) error {
	if r.Updates == nil {
		r.Updates = []resolver.Update{}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}

func WriteText(w io.Writer, r Report) error {
	if _, err := fmt.Fprintf(w, "Module:  %s\n", r.Module); err != nil {
		return err
	}
	goVer := r.GoVersion
	if goVer == "" {
		goVer = "(unspecified)"
	}
	if _, err := fmt.Fprintf(w, "Go:      %s\n\n", goVer); err != nil {
		return err
	}

	if len(r.Updates) == 0 {
		_, err := fmt.Fprintln(w, "All dependencies are up to date.")
		return err
	}

	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "DEPENDENCY\tCURRENT\tLATEST\tTYPE\t")
	for _, u := range r.Updates {
		typeCol := string(u.Type)
		if u.Type == resolver.Major && u.LatestPath != "" && u.LatestPath != u.Path {
			typeCol = fmt.Sprintf("major (-> %s)", u.LatestPath)
		}
		dep := u.Path
		if u.Indirect {
			dep += " (indirect)"
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t\n", dep, u.Current, u.Latest, typeCol)
	}
	return tw.Flush()
}
