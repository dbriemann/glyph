package tmpl

const (
	IndexTmpl = `
<!doctype html> 
<head> 
	<meta charset="utf-8"> 
	<meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1"> 
	<title>{{.TITLE}}</title>
	<meta name="description" content=""> 
	<meta name="viewport" content="width=device-width, initial-scale=1"> 
	<link rel="stylesheet" href="css/normalize.min.css">
	<link rel="stylesheet" href="css/milligram.min.css">
</head> 
<body>
	<div id="main-content">
		<h1>{{.TITLE}}</h1>
		<ul>
		{{range .ISSUES}}
			<li><a href={{.Link}}>{{.Title}}</a></li>
		{{end}}
		</ul> 
	</div>
</body>
</html>
`
)
