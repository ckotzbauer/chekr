package deprecation

var KyvernoTemplate = `
apiVersion: kyverno.io/v1
kind: ClusterPolicy
metadata:
  name: chekr-check-deprecated-apis
  labels:
    creator: chekr
  annotations:
    policies.kyverno.io/title: Check deprecated APIs
    policies.kyverno.io/category: {{.Category}}
    policies.kyverno.io/subject: {{.Subject}}
    policies.kyverno.io/description: >-
      Kubernetes APIs are sometimes deprecated and removed after a few releases.
      As a best practice, older API versions should be replaced with newer versions.
      This policy validates for APIs that are deprecated or scheduled for removal.
      Note that checking for some of these resources may require modifying the Kyverno
      ConfigMap to remove filters.
spec:
  validationFailureAction: {{.ValidationFailureAction}}
  background: {{.Background}}
  rules:
  {{- range .Versions}}
  - name: {{.Name}}
    match:
      resources:
        kinds:
        {{- range .Kinds}}
        - {{.}}
        {{- end}}
    validate:
      message: >-
        {{.Message}}
        See: https://kubernetes.io/docs/reference/using-api/deprecation-guide/
      deny: {}
  {{- end}}
`
