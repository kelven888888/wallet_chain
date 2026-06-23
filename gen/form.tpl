{{define "main"}}
{{.Main}}
	<div class="x-body">
		<form class="layui-form">
	    {{.Items}}


			{{.Btnsub}}
		</form>
	</div>
	{{.Scrip}}

</body>
{{.Ends}}
{{end}}