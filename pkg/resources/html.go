package resources

var HtmlPage = `
    <h2>Resource usage</h2>
    <table class="responsive-table highlight">
        <thead>
            <tr>
                <th>Namespace</th>
                <th>Pod</th>
                <th>Memory Requests</th>
                <th>Memory Limits</th>
                <th>CPU Requests</th>
                <th>CPU Limits</th>
            </tr>
        </thead>
        <tbody>
            {{range .Items}}
            <tr>
                <td>{{.Namespace}}</td>
                <td>{{.Pod}}</td>
                <td>{{.MemoryRequests.Min.FormatMemory}}<br>{{.MemoryRequests.Avg.FormatMemory}}<br>{{.MemoryRequests.Max.FormatMemory}}</td>
                <td>{{.MemoryLimits.Min.FormatMemory}}<br>{{.MemoryLimits.Avg.FormatMemory}}<br>{{.MemoryLimits.Max.FormatMemory}}</td>
                <td>{{.CPURequests.Min.FormatCPU}}<br>{{.CPURequests.Avg.FormatCPU}}<br>{{.CPURequests.Max.FormatCPU}}</td>
                <td>{{.CPULimits.Min.FormatCPU}}<br>{{.CPULimits.Avg.FormatCPU}}<br>{{.CPULimits.Max.FormatCPU}}</td>
            </tr>
            {{end}}
        </tbody>
    </table>`
