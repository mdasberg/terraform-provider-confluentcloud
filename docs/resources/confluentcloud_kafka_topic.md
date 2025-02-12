---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "confluentcloud_kafka_topic Resource - terraform-provider-confluentcloud"
subcategory: ""
description: |-
  
---

# confluentcloud_kafka_topic Resource

`confluentcloud_kafka_topic` provides a Kafka Topic resource that enables creating and deleting Kafka Topics on a Kafka cluster on Confluent Cloud.

## Example Usage

```terraform
resource "confluentcloud_environment" "test-env" {
  display_name = "Development"
}

resource "confluentcloud_kafka_cluster" "basic-cluster" {
  display_name = "basic_kafka_cluster"
  availability = "SINGLE_ZONE"
  cloud        = "GCP"
  region       = "us-central1"
  basic {}

  environment {
    id = confluentcloud_environment.test-env.id
  }
}

resource "confluentcloud_kafka_topic" "orders" {
  kafka_cluster      = confluentcloud_kafka_cluster.basic-cluster.id
  topic_name         = "orders"
  partitions_count   = 4
  http_endpoint      = confluentcloud_kafka_cluster.basic-cluster.http_endpoint
  config = {
    "cleanup.policy"    = "compact"
    "max.message.bytes" = "12345"
    "retention.ms"      = "67890"
  }
  credentials {
    key    = "<Kafka API Key for confluentcloud_kafka_cluster.basic-cluster>"
    secret = "<Kafka API Secret for confluentcloud_kafka_cluster.basic-cluster>"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Argument Reference

The following arguments are supported:

- `kafka_cluster` - (Required String) The ID of the Kafka cluster, for example, `lkc-abc123`.
- `topic_name` - (Required String) The name of the topic, for example, `orders-1`. The topic name can be up to 255 characters in length and can contain only alphanumeric characters, hyphens, and underscores.
- `http_endpoint` - (Required String) The REST endpoint of the Kafka cluster, for example, `https://pkc-00000.us-central1.gcp.confluent.cloud:443`).
- `credentials` (Required Configuration Block) supports the following:
    - `key` - (Required String) The Kafka API Key.
    - `secret` - (Required String) The Kafka API Secret.

-> **Note:** A Kafka API key consists of a key and a secret. Kafka API keys are required to interact with Kafka clusters in Confluent Cloud. Each Kafka API key is valid for one specific Kafka cluster.

-> **Note:** To rotate a Kafka API key, create a new Kafka API key, update `credentials` block in all configuration files to use the new Kafka API key, run `terraform apply`, and remove the old Kafka API key.

- `partitions_count` - (Optional Number) The number of partitions to create in the topic. Defaults to `6`.
- `config` - (Optional String Map) The custom topic configurations to set:
    - `name` - (Required String) The configuration name, for example, `cleanup.policy`.
    - `value` - (Required String) The configuration value, for example, `compact`.

!> **Warning:** Terraform doesn't encrypt the sensitive `credentials` value of the `confluentcloud_kafka_topic` resource, so you must keep your state file secure to avoid exposing it. Refer to the [Terraform documentation](https://www.terraform.io/docs/language/state/sensitive-data.html) to learn more about securing your state file.

## Attributes Reference

In addition to the preceding arguments, the following attributes are exported:

- `id` - (String) The ID of the Kafka topic, in the format `<Kafka cluster ID>/<Kafka Topic name>`, for example, `lkc-abc123/orders-1`.

## Import

-> **Note:** `KAFKA_API_KEY` (`credentials.key`), `KAFKA_API_SECRET` (`credentials.secret`), and `KAFKA_HTTP_ENDPOINT` (`http_endpoint`) environment variables must be set before importing a Kafka topic.

Import Kafka topics by using the Kafka cluster ID and Kafka topic name in the format `<Kafka cluster ID>/<Kafka topic name>`, for example:

```shell
$ export KAFKA_API_KEY="<kafka_api_key>"
$ export KAFKA_API_SECRET="<kafka_api_secret>"
$ export KAFKA_HTTP_ENDPOINT="<kafka_http_endpoint>"
$ terraform import confluentcloud_kafka_topic.my_topic lkc-abc123/orders-123
```
