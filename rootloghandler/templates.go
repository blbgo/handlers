package rootloghandler

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
{{define "log"}}
{{template "header"}}
<h2>Log</h2>
<nav class="f5 b mb3">
	{{- if ne .From "" -}}
		<a href="/log/view/{{.ID}}/f" class="dim link near-white mr3">Start</a>
		<a href="/log/view/{{.ID}}/r{{.From}}" class="dim link near-white mr3">Previous</a>
	{{- else -}}
		<span class="white-30 mr3">Start</span>
		<span class="white-30 mr3">Previous</span>
	{{- end -}}
	{{- if ne .To "" -}}
		<a href="/log/view/{{.ID}}/f{{.To}}" class="dim link near-white mr3">Next</a>
		<a href="/log/view/{{.ID}}/r" class="dim link near-white mr3">End</a>
	{{- else -}}
			<span class="white-30 mr3">Next</span>
			<span class="white-30 mr3">End</span>
	{{- end -}}
	{{- /*
		<form action="/log/prune" method="post" class="di">
			<input type="hidden" name="id" value="{{.ID}}" />
			<input type="submit" value="Prune" class="dim link bg-transparent bn pa0 near-white mr3" />
		</form>
	*/ -}}
</nav>
<table class="collapse ba br2 b--white-10 pv2 ph3">
	<thead>
		<tr class="bg-black-30">
			<th class="pv2 ph3 tl f6 fw6">Time</th>
			<th class="pv2 ph3 tl f6 fw6">Message</th>
		</tr>
	</thead>
	<tbody>
		{{range .Records -}}
			<tr class="stripe-dark">
				<td class="pv2 ph3">{{.When}}</td>
				<td class="pv2 ph3">{{.Message}}</td>
			</tr>
		{{- end}}
	</tbody>
</table>
{{template "footer"}}
{{end}}

{{define "logs"}}
{{template "header"}}
<h2>Logs</h2>
<table class="collapse ba br2 b--white-10 pv2 ph3">
	<thead>
		<tr class="bg-black-30">
			<th class="pv2 ph3 tl f6 fw6">Created</th>
			<th class="pv2 ph3 tl f6 fw6">Name</th>
			<th class="pv2 ph3 tl f6 fw6">Action</th>
		</tr>
	</thead>
	<tbody>
		{{range . -}}
			<tr class="stripe-dark">
				<td class="pv2 ph3">{{.Created}}</td>
				<td class="pv2 ph3">
					<a href="/log/view/{{.ID}}/r" class="dim link near-white b">{{.Name}}</a>
				</td>
				<td class="pv2 ph3">
					<form action="/log/del" method="post" class="di">
						<input type="hidden" name="id" value="{{.ID}}">
						<input type="submit" value="Delete" class="link bg-transparent bn pa0 dim white b">
					</form>
				</td>
			</tr>
		{{- end}}
	</tbody>
</table>
{{template "footer"}}
{{end}}
`

func (r templates) Template() string {
	return templateString
}
