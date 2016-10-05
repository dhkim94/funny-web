package cktime

import (
	"time"
	"strings"
	"fmt"
)

type Cktime struct {
	tm		time.Time
	format		string
	_cformat 	string
}

type _fme struct {
	from	string
	to	string
}

var _fmlist = []_fme {
	_fme{"YYYY",	"2006"},
	_fme{"YY",	"06"},
	_fme{"MMMM",	"January"},
	_fme{"MMM",	"Jan"},
	_fme{"MM",	"01"},
	_fme{"M",	"1"},
	_fme{"DDDD",	"Monday"},
	_fme{"DDD",	"Mon"},
	_fme{"DD",	"02"},
	_fme{"D",	"2"},
	_fme{"hh",	"15"},
	_fme{"h",	"03"},
	_fme{"mm",	"04"},
	_fme{"SSS",	"00000"},
	_fme{"ss",	"05"},
}

// Replace time format
func replace(_in string) (out string)  {
	out = _in
	for _, val := range _fmlist {
		out = strings.Replace(out, val.from, val.to, 1)
	}

	return
}

// Create new Cktime
func NewCktime(_tm time.Time, _format string) *Cktime {
	if len(_format) <= 0 {
		_format = "YYYY-MM-DD hh:mm:ss.SSS"
	}

	return &Cktime{
		tm: _tm,
		format: _format,
	}
}

// Set time format of Cktime
func (tm *Cktime) SetFormat(_format string) {
	if len(_format) <= 0 {
		_format = "YYYY-MM-DD hh:mm:ss.SSS"
	}

	tm._cformat = ""
	tm.format = _format
}

// Set time of Cktime
func (tm *Cktime) SetTime(_tm time.Time) {
	tm.tm = _tm
}

// Convert from time to string
func (tm *Cktime) ToString() (out string)  {
	if len(tm._cformat) <= 0 {
		tm._cformat = replace(tm.format)
	}

	out = tm.tm.Format(tm._cformat)

	return
}

// Convert from string to time
func (tm *Cktime) ToTime(str ...string) (otm time.Time, oerr error) {
	if len(str) <= 0 {
		otm = time.Now()
		oerr = fmt.Errorf("not found time format for tm [%s]", otm)
		return
	}


	fmt.Println(len(str))

	//if len(tm._cformat) <= 0 {
	//	tm._cformat = replace(tm.format)
	//}
	//
	//otm, oerr = time.Parse(tm._cformat, _stm)

	return
}