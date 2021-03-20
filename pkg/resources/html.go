package resources

var HtmlPage = `<html>
    <body>
        <table>
            <thead>
                <tr>
                    <td>Pod</td>
                    <td>Memory Requests</td>
                    <td>Memory Limits</td>
                    <td>CPU Requests</td>
                    <td>CPU Limits</td>
                </tr>
            </thead>
            <tbody>
				{{range .Items}}
                <tr>
                    <td>{{.Pod}}</td>
                    <td>{{.MemoryRequests.Min.FormatMemory}}<br>{{.MemoryRequests.Avg.FormatMemory}}<br>{{.MemoryRequests.Max.FormatMemory}}</td>
                    <td>{{.MemoryLimits.Min.FormatMemory}}<br>{{.MemoryLimits.Avg.FormatMemory}}<br>{{.MemoryLimits.Max.FormatMemory}}</td>
                    <td>{{.CPURequests.Min.FormatCPU}}<br>{{.CPURequests.Avg.FormatCPU}}<br>{{.CPURequests.Max.FormatCPU}}</td>
                    <td>{{.CPULimits.Min.FormatCPU}}<br>{{.CPULimits.Avg.FormatCPU}}<br>{{.CPULimits.Max.FormatCPU}}</td>
                </tr>
				{{end}}
            </tbody>
        </table>
    </body>
</html>`
