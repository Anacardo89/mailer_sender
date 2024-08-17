package handlers

var (
	registerSubject = `Welcome to the Blog {{.User}}. Please verify the account.`

	registerBody = `
	We're glad you could join us {{.User}}. Please click the link below to verify your account:
	{{.Link}}
	`

	recoverPasswordSubject = `Password Recovery for TPSI25 Blog`

	recoverPasswordBody = `
	Here's your recovery link {{.User}}.
	
	Please click the link below:
	{{.Link}}
	`
)
