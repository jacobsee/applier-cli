package yamlresources

// Template represents an OpenShift template resource
type Template struct {
	APIVersion string                   `yaml:"apiVersion"`
	Kind       string                   `yaml:"kind"`
	Metadata   TemplateMetadata         `yaml:"metadata"`
	Objects    []map[string]interface{} `yaml:"objects"`
	Parameters []TemplateParameter      `yaml:"parameters"`
}

// TemplateMetadata represents the metadata within an OpenShift template resource
type TemplateMetadata struct {
	Name string `yaml:"name"`
}

// TemplateParameter represents a single parameter within an OpenShift template resource
type TemplateParameter struct {
	Description string `yaml:"description"`
	Name        string `yaml:"name"`
}

/*
apiVersion: v1
kind: Template
metadata:
  name: redis-template
  annotations:
    description: "Description"
    iconClass: "icon-redis"
    tags: "database,nosql"
objects:
- apiVersion: v1
  kind: Pod
  metadata:
    name: redis-master
  spec:
    containers:
    - env:
      - name: REDIS_PASSWORD
        value: ${REDIS_PASSWORD}
      image: dockerfile/redis
      name: master
      ports:
      - containerPort: 6379
        protocol: TCP
parameters:
- description: Password used for Redis authentication
  from: '[A-Z0-9]{8}'
  generate: expression
  name: REDIS_PASSWORD
labels:
  redis: master
*/
