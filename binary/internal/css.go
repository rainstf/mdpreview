package internal

func Css() string {
	return `<style>
	body {
		margin: 0 auto;
		max-width: 792px;
		font-family: -apple-system, BlinkMacSystemFont, Segoe UI, Helvetica, Arial, sans-serif, Apple Color Emoji, Segoe UI Emoji;
		line-height: 1.6;
		font-size: 16px;
		color: #333;
		background-color: #fff;
		padding: 20px;
	}
	h1, h2, h3, h4, h5, h6 {
		margin-top: 24px;
		margin-bottom: 16px;
		font-weight: 600;
	}
	p {
    margin-top: 0;
    margin-bottom: 16px;
}

a {
    color: #0366d6;
    text-decoration: none;
}

a:hover {
    text-decoration: underline;
}

code {
    background-color: #f6f8fa;
    padding: .2em .4em;
    border-radius: 5px;
}

pre {
    background-color: #f6f8fa;
    border-radius: 6px;
    font-family: SFMono-Regular, Consolas, Liberation Mono, Menlo, monospace;
    padding: 16px;
    font-size: 85%;
    overflow: auto;
}

pre code {
    background-color: transparent;
    padding: 0;
    margin: 0;
    border: 0;
}

ul, ol {
    padding-left: 2em;
}

blockquote {
    color: #57606a;
    border-left: .25em solid #dfe2e5;
    padding: 0 1em;
    margin-left: 0;
    margin-right: 0;
    margin-top: 0;
    margin-bottom: 16px;
}

table {
    border-collapse: collapse;
    margin-top: 0;
    margin-bottom: 16px;
    width: 100%;
}

table, th, td {
    border: 1px solid #dfe2e5;
}

th, td {
    padding: 6px 13px;
}

th {
    font-weight: 600;
}
  </style>`
}
