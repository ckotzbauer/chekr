package ha

var HtmlPage = `
    <h2>High Availability</h2>
    <table class="responsive-table highlight">
        <thead>
            <tr>
                <th>Namespace</th>
                <th>Name</th>
                <th>Type</th>
                <th>Replicas</th>
                <th>Rollout Strategy</th>
                <th>Anti-Affinity</th>
                <th>PVCs</th>
                <th>Rank</th>
            </tr>
        </thead>
        <tbody>
            {{range .Items}}
            <tr>
                <td>{{.Namespace}}</td>
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
    <br/>
    <div style="white-space: pre-wrap;">%s</div>`
