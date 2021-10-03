package resources

var HtmlPage = `
    <h2>Resource usage</h2>
    <table class="responsive-table highlight">
        <thead>
            <tr>
                <th>Namespace</th>
                <th>Pod</th>
                <th>Container</th>
                <th> </th>
                <th>Current value</th>
                <th>Min value</th>
                <th>Average value</th>
                <th>Max value</th>
            </tr>
        </thead>
        <tbody>
            {{range $pod := .Items}}
            {{range $container := $pod.Containers}}
            <tr>
                <td>{{$pod.Namespace}}</td>
                <td>{{$pod.Pod}}</td>
                <td>{{$container.Name}}</td>
                {{if $container.MemoryRequests.HasValue}}
                <td>Memory Requests</td>
                <td>{{$container.MemoryRequests.Current.FormatMemory}}</td>
                {{else}}
                <td>Memory</td>
                <td> </td>
                {{end}}
                <td>{{$container.MemoryRequests.Min.FormatMemory}}</td>
                <td>{{$container.MemoryRequests.Avg.FormatMemory}}</td>
                <td>{{$container.MemoryRequests.Max.FormatMemory}}</td>
            </tr>
            <tr>
                <td>{{$pod.Namespace}}</td>
                <td>{{$pod.Pod}}</td>
                <td>{{$container.Name}}</td>
                {{if $container.MemoryLimits.HasValue}}
                <td>Memory Limits</td>
                <td>{{$container.MemoryLimits.Current.FormatMemory}}</td>
                {{else}}
                <td>Memory</td>
                <td> </td>
                {{end}}
                <td>{{$container.MemoryLimits.Min.FormatMemory}}</td>
                <td>{{$container.MemoryLimits.Avg.FormatMemory}}</td>
                <td>{{$container.MemoryLimits.Max.FormatMemory}}</td>
            </tr>
            <tr>
                <td>{{$pod.Namespace}}</td>
                <td>{{$pod.Pod}}</td>
                <td>{{$container.Name}}</td>
                {{if $container.CPURequests.HasValue}}
                <td>CPU Requests</td>
                <td>{{$container.CPURequests.Current.FormatCPU}}</td>
                {{else}}
                <td>CPUs</td>
                <td> </td>
                {{end}}
                <td>{{$container.CPURequests.Min.FormatCPU}}</td>
                <td>{{$container.CPURequests.Avg.FormatCPU}}</td>
                <td>{{$container.CPURequests.Max.FormatCPU}}</td>
            </tr>
            <tr>
                <td>{{$pod.Namespace}}</td>
                <td>{{$pod.Pod}}</td>
                <td>{{$container.Name}}</td>
                {{if $container.CPULimits.HasValue}}
                <td>CPU Limits</td>
                <td>{{$container.CPULimits.Current.FormatCPU}}</td>
                {{else}}
                <td>CPUs</td>
                <td> </td>
                {{end}}
                <td>{{$container.CPULimits.Min.FormatCPU}}</td>
                <td>{{$container.CPULimits.Avg.FormatCPU}}</td>
                <td>{{$container.CPULimits.Max.FormatCPU}}</td>
            </tr>
            {{end}}
            {{end}}
        </tbody>
    </table>`
