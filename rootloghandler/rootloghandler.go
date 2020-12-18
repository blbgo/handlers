package rootloghandler

import (
	"encoding/base64"
	"errors"
	"net/http"
	"time"

	"github.com/blbgo/httpserver"
	"github.com/blbgo/record/rootlog"
)

type log struct {
	httpserver.RouteParams
	rootlog.RootLog
	httpserver.Renderer
}

type logModel struct {
	ID      string
	Created string
	Name    string
}

type messagesModel struct {
	ID      string
	From    string
	To      string
	Records []*messageModel
}

const logTimeFormat = "2006-01-02 15:04:05"

type messageModel struct {
	When    string
	Message string
}

// Setup setups all the routes for log viewing
func Setup(
	router httpserver.Router,
	rp httpserver.RouteParams,
	rootLog rootlog.RootLog,
	rend httpserver.Renderer,
) {
	r := &log{RouteParams: rp, RootLog: rootLog, Renderer: rend}
	router.Handler("GET", "/log", http.HandlerFunc(r.log))
	router.Handler("GET", "/log/view/:id/:offset", http.HandlerFunc(r.view))
	router.Handler("POST", "/log/del", http.HandlerFunc(r.del))
	router.Handler("POST", "/log/prune", http.HandlerFunc(r.prune))
}

func (r *log) log(rw http.ResponseWriter, req *http.Request) {
	var logList []*logModel
	var innerError error
	err := r.Range(rootlog.MinTime(), false, func(created time.Time, name string) bool {
		var createdString string
		createdString, innerError = encodeTime(created)
		if innerError != nil {
			return false
		}
		logList = append(
			logList,
			&logModel{
				ID:      createdString,
				Created: created.Format(logTimeFormat),
				Name:    name,
			},
		)
		return true
	})
	if innerError != nil {
		r.Error(rw, "error", innerError)
		return
	}
	if err != nil {
		r.Error(rw, "error", err)
		return
	}
	r.OK(rw, "logs", logList)
}

func (r *log) view(rw http.ResponseWriter, req *http.Request) {
	params := r.Get(req)
	if len(params) != 2 {
		r.Error(rw, "error", errors.New("invalid request"))
		return
	}
	created, err := decodeTime(params[0])
	if err != nil {
		r.Error(rw, "error", err)
		return
	}
	offset := params[1]
	if len(offset) < 1 {
		r.Error(rw, "error", errors.New("invalid request"))
		return
	}
	reverse := offset[0] != 'f'
	offset = offset[1:]
	var start time.Time
	if len(offset) == 0 {
		if reverse {
			start = rootlog.MaxTime()
		} else {
			start = rootlog.MinTime()
		}
	} else {
		start, err = decodeTime(offset)
		if err != nil {
			r.Error(rw, "error", err)
			return
		}
		if reverse {
			start = start.Add(time.Nanosecond * -1)
		} else {
			start = start.Add(time.Nanosecond)
		}
	}

	messages := &messagesModel{ID: params[0]}
	var innerError error
	err = r.RangeLog(created, start, reverse, func(created time.Time, message string) bool {
		messages.Records = append(
			messages.Records,
			&messageModel{
				When:    created.Format(logTimeFormat),
				Message: message,
			},
		)
		records := len(messages.Records)
		if records == 1 {
			messages.From, innerError = encodeTime(created)
			if innerError != nil {
				return false
			}
		}
		if records == 20 {
			messages.To, innerError = encodeTime(created)
			if innerError != nil {
				return false
			}
		}
		return records != 21
	})
	if innerError != nil {
		r.Error(rw, "error", innerError)
		return
	}
	if err != nil {
		r.Error(rw, "error", err)
		return
	}
	records := len(messages.Records)
	if reverse {
		messages.To, messages.From = messages.From, messages.To

		if len(offset) == 0 {
			messages.To = ""
		}
		if records != 21 {
			messages.From = ""
		} else {
			messages.Records = messages.Records[:records-1]
			records--
		}
		recordsM1 := records - 1
		recordsHalf := records / 2
		for i := 0; i < recordsHalf; i++ {
			messages.Records[i], messages.Records[recordsM1-i] =
				messages.Records[recordsM1-i], messages.Records[i]
		}
	} else {
		if len(offset) == 0 {
			messages.From = ""
		}
		if records != 21 {
			messages.To = ""
		} else {
			messages.Records = messages.Records[:records-1]
		}
	}
	r.OK(rw, "log", messages)
}

func (r *log) del(rw http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		r.Error(rw, "error", err)
		return
	}
	created, err := decodeTime(req.PostForm.Get("id"))
	if err != nil {
		r.Error(rw, "error", err)
		return
	}
	r.Delete(created)
	if err != nil {
		r.Error(rw, "error", err)
		return
	}
	http.Redirect(rw, req, "/log", http.StatusSeeOther)
}

func (r *log) prune(rw http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		r.Error(rw, "error", err)
		return
	}
	//id, err := strconv.ParseInt(req.PostForm.Get("id"), 10, 64)
	//if err != nil {
	//	r.Error(rw, "error", err)
	//	return
	//}
	r.Error(rw, "error", errors.New("not implemented"))
	//err = r.DeleteLogAllButX(id, 5)
	//if err != nil {
	//	r.Error(rw, "error", err)
	//	return
	//}
	//http.Redirect(rw, req, fmt.Sprint("/log/view/", id, "/0"), http.StatusSeeOther)
}

func decodeTime(encoded string) (time.Time, error) {
	createdBytes, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return time.Time{}, err
	}
	var created time.Time
	err = created.UnmarshalBinary(createdBytes)
	if err != nil {
		return time.Time{}, err
	}
	return created, nil
}

func encodeTime(time time.Time) (string, error) {
	timeBytes, err := time.MarshalBinary()
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(timeBytes), nil
}
