package taskshandler

import (
	"github.com/blbgo/httpserver"
)

type templates struct{}

// NewTemplates provides the templates needed by the log handler
func NewTemplates() httpserver.TemplateProvider {
	return templates{}
}

// these templates assume "header" and "footer" templates are defined elsewhere
var templateString = `
{{define "tasks"}}
{{template "header"}}
<h2>Tasks</h2>
<ul class="list pl3">
	{{range $i, $e := . -}}
		<li class="pb2"><a href="/task/input/{{$i}}" class="dim link near-white b">{{.Name}}</a></li>
	{{- end}}
</ul>
{{template "footer"}}
{{end}}

{{define "taskInputs"}}
{{template "header"}}
<h2>Run task "{{.Name}}"</h2>
<form action="/task/run" method="post">
	<input type="hidden" name="id" value="{{.ID}}">
	{{range $i, $e := .Inputs -}}
		<input class="input-reset ba br3 b--mid-gray pa2 mb3 db w-100 outline-0" placeholder="{{.}}" type="text" name="I{{$i}}">
	{{- end}}
	<input type="submit" value="Run">
</form>
{{template "footer"}}
{{end}}

{{define "task"}}
{{template "header"}}
<h2>Task</h2>
<h3>Task "{{.Name}}" has been started.</h3>
{{template "footer"}}
{{end}}
`

func (r templates) Template() string {
	return templateString
}
