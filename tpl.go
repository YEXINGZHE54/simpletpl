package simpletpl

import (
	"fmt"
	"strconv"
	"strings"
)

type (
	Template struct {
		ID      string
		Tpl     string
		Vars    map[string][]int
		Content []string
	}
	Data map[string]interface{}
)

const (
	PREFIX  = "{{"
	POSTFIX = "}}"
)

var (
	EMPTY_TPL     = fmt.Errorf("empty tpl")
	BADKEY_ERROR  = fmt.Errorf("bad key in tpl data")
	KEY_NONEXISTS = fmt.Errorf("key does not exist in data")
)

func NewTemplate(pattern string, id string) (t *Template, err error) {
	t = new(Template)
	t.ID = id
	t.Tpl = pattern
	err = t.Compile()
	return
}

func (t *Template) TrueID() string {
	return t.ID
}

func (t *Template) Compile() (err error) {
	if t.Tpl == "" {
		return EMPTY_TPL
	}
	t.Vars = make(map[string][]int)
	idx := 0
	for idx < len(t.Tpl) {
		post := strings.Index(t.Tpl[idx:], POSTFIX)
		if post < 0 {
			// not found break
			break
		}
		// update idx
		newidx := idx + post + len(POSTFIX)
		prev := strings.LastIndex(t.Tpl[idx:idx+post], PREFIX)
		if prev < 0 {
			// no prev found,yet
			// so we store the literal first, then switch to next search loop
			t.Content = append(t.Content, t.Tpl[idx:newidx])
			idx = newidx
			continue
		}
		//if we found one, first step is to store literal before prev
		literal := t.Tpl[idx : idx+prev]
		t.Content = append(t.Content, literal)
		tplvar := strings.TrimSpace(t.Tpl[idx+prev+len(PREFIX) : idx+post])
		if tplvar == "" {
			idx = newidx
			continue
		}
		vlist, ok := t.Vars[tplvar]
		if !ok {
			vlist = make([]int, 0)
		}
		vlist = append(vlist, len(t.Content))
		t.Vars[tplvar] = vlist
		t.Content = append(t.Content, "")
		idx = newidx
		continue
	}
	if idx < len(t.Tpl) {
		// it means: we did not find post and break
		t.Content = append(t.Content, t.Tpl[idx:])
	}
	return nil
}

func (t *Template) Render(data map[string]interface{}) (msg string, err error) {
	var content []string
	for _, s := range t.Content {
		content = append(content, s)
	}
	for k, vlist := range t.Vars {
		vstr, err := Force(data, k)
		if err != nil {
			return "", err
		}
		for _, idx := range vlist {
			content[idx] = vstr
		}
	}
	msg = strings.Join(content, "")
	return
}

func Force(data interface{}, key string) (value string, err error) {
	vlen := 0
	leninfo := strings.Split(key, ":")
	key = leninfo[0]
	if len(leninfo) > 1 {
		vlen, err = strconv.Atoi(leninfo[1])
		if err != nil {
			return
		}
	}

	fields := strings.Split(key, ".")
	var ok bool
	for _, k := range fields {
		switch data.(type) {
		case map[string]interface{}:
			m, _ := data.(map[string]interface{})
			data, ok = m[k]
			if !ok {
				err = fmt.Errorf("part key %s does not exist in key %s, cur map: %v", k, key, m)
				return
			}
		default:
			err = fmt.Errorf("bad key, expecting map[string]interface{}, got: %T", data)
			return
		}
	}
	switch v := data.(type) {
	case map[string]interface{}:
		err = BADKEY_ERROR
	case string:
		value = v
		if vlen > 0 && len(value) > vlen {
			value = value[:vlen]
		}
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		var format string
		if vlen > 0 {
			format = fmt.Sprintf("%%0%dd", vlen)
		} else {
			format = "%d"
		}
		value = fmt.Sprintf(format, v)
	case []byte:
		value = string(v)
		if vlen > 0 && len(value) > vlen {
			value = value[:vlen]
		}
	}
	return
}
