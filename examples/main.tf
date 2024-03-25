resource "keep_provider" "prometheus" {
  name = "prometheus-dev"
  type = "prometheus"
  auth_config = {
	url = "http://localhost:9090"
    /*
    from keep cli you can easily get the which config params are needed for the provider you want to connect
    ~ keep provider connect prometheus --help
    +------------+--------------+----------+-----------------------+
    |  Provider  | Config Param | Required |      Description      |
    +------------+--------------+----------+-----------------------+
    | prometheus |    --url     |   True   | Prometheus server URL |
    |            |  --username  |  False   |  Prometheus username  |
    |            |  --password  |  False   |  Prometheus password  |
    +------------+--------------+----------+-----------------------+
    */
  }
}

resource "keep_workflow" "example_workflow" {
  workflow_file_path = "path/to/workflow.yml"
}

resource "keep_mapping" "example_mapping" {
  name = "example_mapping"
  mapping_file_path = "path/to/mapping.yml"
  matchers = [
    "your unique matcher",
  ]
}

output "keep_provider_id" {
  value = keep_provider.prometheus.id
}