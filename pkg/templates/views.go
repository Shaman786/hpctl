package templates

// ---------------------------------------------------------
// USER / IDENTITY VIEWS
// ---------------------------------------------------------

const WhoAmI = `
👤  USER PROFILE
------------------------------------------------
ID:		{{.ID}}
Name:		{{.FirstName}} {{.LastName}}
Email:		{{.Email}}
Phone:		{{.Phone}}
Location:	{{.City}}, {{.Country}}
------------------------------------------------
`

// ---------------------------------------------------------
// INFRASTRUCTURE / SERVICES VIEWS
// ---------------------------------------------------------

// ServicesList loops through the .Services array automatically
const ServicesList = `
YOUR INFRASTRUCTURE
--------------------------------------------------------------------------------
ID	NAME	DOMAIN	PRICE	STATUS	NEXT DUE
{{range .Services}}{{.ID}}	{{.Name}}	{{.Domain}}	{{.Price}}	{{.Status}}	{{.NextDue}}
{{end}}--------------------------------------------------------------------------------
`

const ServiceDetail = `
📋 SERVICE DETAILS
------------------------------------------------
Name:        {{ .Name }}
Domain:      {{ .Domain }}
Status:      {{ .Status }}
IP Address:  {{ .DedicatedIP }}
Price:       {{ .Amount }} {{ .BillingCycle }}
Next Due:    {{ .NextDueDate }}
Payment:     {{ .PaymentMethod }}
------------------------------------------------
`
