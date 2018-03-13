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
	<link rel="stylesheet" href="//fonts.googleapis.com/css?family=Roboto:300,300italic,700,700italic">
	<link rel="stylesheet" href="//cdn.rawgit.com/necolas/normalize.css/master/normalize.css">	
	<link rel="stylesheet" href="//cdn.rawgit.com/milligram/milligram/master/dist/milligram.min.css">
</head> 
<body>
	<div class="main-content">
		<h1 class="blog-title">{{.TITLE}}</h1>
		<ul class="blog-issues">
		{{range .ISSUES}}
			<li class="blog-entry"><a href={{.Link}}>{{.Title}}</a></li>
		{{end}}
		</ul> 
	</div>
</body>
</html>
`

	IssueTmlp = `
	<!doctype html> 
	<head> 
		<meta charset="utf-8"> 
		<meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1"> 
		<title>{{.TITLE}}</title>
		<meta name="description" content=""> 
		<meta name="viewport" content="width=device-width, initial-scale=1"> 
		<link rel="stylesheet" href="//fonts.googleapis.com/css?family=Roboto:300,300italic,700,700italic">
		<link rel="stylesheet" href="//cdn.rawgit.com/necolas/normalize.css/master/normalize.css">	
		<link rel="stylesheet" href="//cdn.rawgit.com/milligram/milligram/master/dist/milligram.min.css">
	</head> 
	<body>
		<div class="main-content">
			<h1 class="blog-title">{{.TITLE}}</h1>
			<h2 class="issue-title">{{.ISSUE_TITLE}}</h2>
			<p class="issue-content">{{.ISSUE_CONTENT}}</p>		
		</div>
	</body>
	</html>
`
)
