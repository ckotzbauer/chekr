package deprecation

var HtmlPage = `
    <h2>API deprecations</h2>
    <table class="responsive-table highlight">
        <thead>
            <tr>
                <th>Namespace</th>
                <th>Name</th>
                <th>Deprecated GV</th>
                <th>Deprecated Kind</th>
                <th>Replacement GV</th>
                <th>Replacement Kind</th>
                <th>Deprecation Version</th>
                <th>Removal Version</th>
            </tr>
        </thead>
        <tbody>
            {{range .Items}}
            <tr>
                <td>{{.Namespace}}</td>
                <td>{{.Name}}</td>
                <td>{{.DeprecatedGroupVersion}}</td>
                <td>{{.DeprecatedKind}}</td>
                <td>{{.ReplacementGroupVersion}}</td>
                <td>{{.ReplacementKind}}</td>
                <td>{{.DeprecationVersion}}</td>
                <td>{{.RemovalVersion}}</td>
            </tr>
            {{end}}
        </tbody>
    </table>
    <br/>`
