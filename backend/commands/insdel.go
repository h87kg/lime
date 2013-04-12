package commands

import (
	"fmt"
	. "lime/backend"
)

type (
	InsertCommand struct {
		DefaultCommand
	}

	LeftDeleteCommand struct {
		DefaultCommand
	}

	RightDeleteCommand struct {
		DefaultCommand
	}
)

func (c *InsertCommand) Run(v *View, e *Edit, args Args) error {
	sel := v.Sel()
	chars, ok := args["characters"].(string)
	if !ok {
		return fmt.Errorf("insert: Missing or invalid characters argument: %v", args)
	}
	for i := 0; i < sel.Len(); i++ {
		r := sel.Get(i)
		if r.Size() == 0 {
			v.Insert(e, r.B, chars)
		} else {
			v.Replace(e, r, chars)
		}
	}
	return nil
}

func (c *LeftDeleteCommand) Run(v *View, e *Edit, args Args) error {
	trim_space := false
	tab_size := 4
	if t, ok := v.Settings().Get("translate_tabs_to_spaces", false).(bool); ok && t {
		if t, ok := v.Settings().Get("use_tab_stops", true).(bool); ok && t {
			trim_space = true
			tab_size, ok = v.Settings().Get("tab_size", 4).(int)
			if !ok {
				tab_size = 4
			}
		}
	}

	sel := v.Sel()
	hasNonEmpty := false
	for _, r := range sel.Regions() {
		if !r.Empty() {
			hasNonEmpty = true
			break
		}
	}
	i := 0
	for {
		l := sel.Len()
		if i >= l {
			break
		}
		r := sel.Get(i)
		if r.A == r.B && !hasNonEmpty {
			d := v.Buffer().Runes()
			if trim_space {
				_, col := v.Buffer().RowCol(r.A)
				prev_col := r.A - (col - (col-tab_size+(tab_size-1))&^(tab_size-1))
				if prev_col < 0 {
					prev_col = 0
				}
				for r.A > prev_col && d[r.A-1] == ' ' {
					r.A--
				}
			}
			if r.A == r.B {
				r.A--
			}
		}
		v.Erase(e, r)
		if sel.Len() != l {
			continue
		}
		i++
	}
	return nil
}

func (c *RightDeleteCommand) Run(v *View, e *Edit, args Args) error {
	sel := v.Sel()
	hasNonEmpty := false
	for _, r := range sel.Regions() {
		if !r.Empty() {
			hasNonEmpty = true
			break
		}
	}
	i := 0
	for {
		l := sel.Len()
		if i >= l {
			break
		}
		r := sel.Get(i)
		if r.A == r.B && !hasNonEmpty {
			r.B++
		}
		v.Erase(e, r)
		if sel.Len() != l {
			continue
		}
		i++
	}
	return nil
}

func init() {
	register([]cmd{
		{"insert", &InsertCommand{}},
		{"left_delete", &LeftDeleteCommand{}},
		{"right_delete", &RightDeleteCommand{}},
	})
}