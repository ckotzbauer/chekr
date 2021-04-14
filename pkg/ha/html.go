package ha

var HtmlPage = `<html>
    <body>
        <table>
            <thead>
                <tr>
                    <td>Name</td>
                    <td>Type</td>
                    <td>Replicas</td>
                    <td>Rollout Strategy</td>
                    <td>Anti-Affinity</td>
                    <td>PVCs</td>
                    <td>Rank</td>
                </tr>
            </thead>
            <tbody>
				{{range .Items}}
                <tr>
                    <td>{{.Name}}</td>
                    <td>{{.Type}}</td>
                    <td>{{.Replicas}}</td>
                    <td>{{.RolloutStrategy}}</td>
                    <td>{{.PodAntiAffinity}}</td>
                    <td>{{.PVC}}</td>
                    <td>{{.Rank}}</td>
                </tr>
				{{end}}
            </tbody>
        </table>
    </body>
</html>`
