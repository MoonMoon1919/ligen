package ligen

import "text/template"

// Template for notice file used for most licenses
const SimpleNoticeTemplateBody = `{{.ProjectName}}
Copyright {{.StartYear}}{{if (gt .EndYear 0) }}-{{.EndYear}}{{end}} {{.Holder}}`

var SimpleNoticeTemplate = template.Must(template.New("SimpleNotice").Parse(SimpleNoticeTemplateBody))
