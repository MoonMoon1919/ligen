package ligen

import "text/template"

// NOTICES
const SimpleNoticeTemplateBody = `{{.ProjectName}}
Copyright {{.StartYear}}{{if (gt .EndYear 0) }}-{{.EndYear}}{{end}} {{.Holder}}`

var SimpleNoticeTemplate = template.Must(template.New("SimpleNotice").Parse(SimpleNoticeTemplateBody))
