package packer

const CSSTemplate = `
.{{.Prefix}} {
  background-image: url({{.ImgURL}}/{{.Name}}.{{.Format}});
  background-repeat: no-repeat;
  display: block;
}

{{range .Images}}
.{{.Name}}{{.Hover}} {
  background-position: {{.X}}px {{.Y}}px;
  width: {{.Width}}px;
  height: {{.Height}}px;
}
{{end}}
`

const HTMLTemplate = `
<html>
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width">
    <title>Sprite Test</title>
    <link href='{{.CSSPath}}'  rel='stylesheet' type='text/css'>
    <style>html {background:#333; color:white;}</style>
  </head>
  <body>
    <table cellspacing="40">
      <tr><th>Sprites</th><th>Image</th></tr>
      <tr>
        <td>
          <table cellpadding="4">
            <th>Name</th><th>Icon</th><th>(X,Y)</th><th>W x H</th></tr>
            {{range .Images}}{{if not .Hover }}<tr><td>{{.Name}}</td><td><div class="sprite {{.Name}}"></div></td><td>({{.X}}, {{.Y}})<td>{{.Width}} x {{.Height}}</td></tr>{{end}}{{end}}
          </table>
        </td>
        <td valign=top>
          <img src="{{.ImgPath}}">
        </td>
      </tr>
    </table>
  </body>
</html>
`

/*
{{#sprite}}
{{class}} {
  background-image: url('{{{escaped_image}}}');
}

{{/sprite}}
{{#retina}}
@media (min--moz-device-pixel-ratio: 1.5), (-o-min-device-pixel-ratio: 3/2), (-webkit-min-device-pixel-ratio: 1.5), (min-device-pixel-ratio: 1.5), (min-resolution: 1.5dppx) {
  {{class}} {
    background-image: url('{{{escaped_image}}}');
    background-size: {{px.total_width}} {{px.total_height}};
  }
}

{{/retina}}
{{#items}}
{{class}} {
  background-position: {{px.offset_x}} {{px.offset_y}};
  width: {{px.width}};
  height: {{px.height}};
}

{{/items}}
*/
