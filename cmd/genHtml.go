package main

func genIndexHtml(user string, repo string) string {

	htmlFeed := `<html>
	<head>
		<title>Code Hosting Service Client</title>
	</head>
	<body>
	<h1>Github Client</h1>`
	if user == "" {
		htmlFeed += `<a href="/login">Log in</a>`
	} else {
		htmlFeed += `Logged in to Github as user: ` + user +
			`<br>Creating branch in repo: ` + repo
	}

	htmlFeed += `</body></html>`

	return htmlFeed
}
