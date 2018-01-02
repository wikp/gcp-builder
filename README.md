# Usage on travis

## Add GITHUB_OAUTH_TOKEN secret

Token must be generated using full repository access scope

## Install tool
```
source <(curl -H 'Accept: application/vnd.github.v3.raw' \
     -s \
     https://api.github.com/repos/wendigo/travis-builder/contents/installer.sh)
```

## Create project descriptor file (service.yml)
```
project:
  name: test
  domain: test
  context: test
  versionPrefix: test-

variables:
  - name: replicas
    value: 1
  - name: containerPort
    value: 9000

environments:
  - name: test
    gcloud:
      registry: eu.gcr.io/project-test
      project: project-test
    kubernetes:
      cluster: container-test
      zone: europe-west1-b
      template: deployment.yml
      variables:
         - name: replicas
           value: 2
images:
  - build: dockerfiles/1
    name: container1
  - build: dockerfiles/2
    name: container2

```

## Create project deployment template for kubernetes (deployment.yml) - this is example only :)

```
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: {{ .Config.Project.FullName }}
  labels:
    domain: {{ .Config.Project.Domain }}
    context: {{ .Config.Project.Context }}
    name: {{ .Config.Project.Name }}
    environment: {{ .EnvironmentName }}
spec:
  replicas: {{ .Variable "replicas" }}
  selector:
    matchLabels:
      domain: {{ .Config.Project.Domain }}
      context: {{ .Config.Project.Context }}
      name: {{ .Config.Project.Name }}
      environment: {{ .EnvironmentName }}
  template:
    metadata:
      labels:
        domain: {{ .Config.Project.Domain }}
        context: {{ .Config.Project.Context }}
        name: {{ .Config.Project.Name }}
        environment: {{ .EnvironmentName }}
    spec:
      volumes:
        - name: chat-brain-service-sql-credentials
          secret:
            secretName: chat-brain-service-sql-credentials
        - name: cloudsql
          emptyDir:
        - name: ssl-certs
          hostPath:
            path: /etc/ssl/certs
      containers:
      - name: brain
        image: {{ .Container "container1" }}
        imagePullPolicy: Always
        ports:
        - containerPort: {{ .Variable "containerPort" }}
---
apiVersion: v1
kind: Service
metadata:
  name: {{ .Config.Project.FullName }}
  labels:
    domain: {{ .Config.Project.Domain }}
    context: {{ .Config.Project.Context }}
    name: {{ .Config.Project.Name }}
    environment: {{ .EnvironmentName }}
spec:
  type: NodePort
  ports:
  - port: {{ .Variable "containerPort" }}
    targetPort: {{ .Variable "containerPort" }}
    protocol: TCP
    name: http
  selector:
    domain: {{ .Config.Project.Domain }}
    context: {{ .Config.Project.Context }}
    name: {{ .Config.Project.Name }}
    environment: {{ .EnvironmentName }}
```