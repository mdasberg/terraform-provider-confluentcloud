// Copyright 2021 Confluent Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package provider

import (
	"context"
	"fmt"
	cmk "github.com/confluentinc/ccloud-sdk-go-v2/cmk/v2"
	iam "github.com/confluentinc/ccloud-sdk-go-v2/iam/v2"
	kafkarestv3 "github.com/confluentinc/ccloud-sdk-go-v2/kafkarest/v3"
	mds "github.com/confluentinc/ccloud-sdk-go-v2/mds/v2"
	org "github.com/confluentinc/ccloud-sdk-go-v2/org/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"os"
	"time"
)

func (c *Client) cmkApiContext(ctx context.Context) context.Context {
	if c.apiKey != "" && c.apiSecret != "" {
		return context.WithValue(context.Background(), cmk.ContextBasicAuth, cmk.BasicAuth{
			UserName: c.apiKey,
			Password: c.apiSecret,
		})
	}
	log.Printf("[WARN] Could not find credentials for Confluent Cloud")
	return ctx
}

func (c *Client) iamApiContext(ctx context.Context) context.Context {
	if c.apiKey != "" && c.apiSecret != "" {
		return context.WithValue(context.Background(), iam.ContextBasicAuth, iam.BasicAuth{
			UserName: c.apiKey,
			Password: c.apiSecret,
		})
	}
	log.Printf("[WARN] Could not find credentials for Confluent Cloud")
	return ctx
}

func (c *Client) mdsApiContext(ctx context.Context) context.Context {
	if c.apiKey != "" && c.apiSecret != "" {
		return context.WithValue(context.Background(), mds.ContextBasicAuth, mds.BasicAuth{
			UserName: c.apiKey,
			Password: c.apiSecret,
		})
	}
	log.Printf("[WARN] Could not find credentials for Confluent Cloud")
	return ctx
}

func (c *Client) orgApiContext(ctx context.Context) context.Context {
	if c.apiKey != "" && c.apiSecret != "" {
		return context.WithValue(context.Background(), org.ContextBasicAuth, org.BasicAuth{
			UserName: c.apiKey,
			Password: c.apiSecret,
		})
	}
	log.Printf("[WARN] Could not find credentials for Confluent Cloud")
	return ctx
}

func getTimeoutFor(clusterType string) time.Duration {
	if clusterType == kafkaClusterTypeDedicated {
		return 24 * time.Hour
	} else {
		return 1 * time.Hour
	}
}

func stringToAclResourceType(aclResourceType string) (kafkarestv3.AclResourceType, error) {
	switch aclResourceType {
	case "UNKNOWN":
		return kafkarestv3.ACLRESOURCETYPE_UNKNOWN, nil
	case "ANY":
		return kafkarestv3.ACLRESOURCETYPE_ANY, nil
	case "TOPIC":
		return kafkarestv3.ACLRESOURCETYPE_TOPIC, nil
	case "GROUP":
		return kafkarestv3.ACLRESOURCETYPE_GROUP, nil
	case "CLUSTER":
		return kafkarestv3.ACLRESOURCETYPE_CLUSTER, nil
	case "TRANSACTIONAL_ID":
		return kafkarestv3.ACLRESOURCETYPE_TRANSACTIONAL_ID, nil
	case "DELEGATION_TOKEN":
		return kafkarestv3.ACLRESOURCETYPE_DELEGATION_TOKEN, nil
	}
	return "", fmt.Errorf("unknown ACL resource type was found: %s", aclResourceType)
}

func stringToAclPatternType(aclPatternType string) (kafkarestv3.AclPatternType, error) {
	switch aclPatternType {
	case "UNKNOWN":
		return kafkarestv3.ACLPATTERNTYPE_UNKNOWN, nil
	case "ANY":
		return kafkarestv3.ACLPATTERNTYPE_ANY, nil
	case "MATCH":
		return kafkarestv3.ACLPATTERNTYPE_MATCH, nil
	case "LITERAL":
		return kafkarestv3.ACLPATTERNTYPE_LITERAL, nil
	case "PREFIXED":
		return kafkarestv3.ACLPATTERNTYPE_PREFIXED, nil
	}
	return "", fmt.Errorf("unknown ACL pattern type was found: %s", aclPatternType)
}

func stringToAclOperation(aclOperation string) (kafkarestv3.AclOperation, error) {
	switch aclOperation {
	case "UNKNOWN":
		return kafkarestv3.ACLOPERATION_UNKNOWN, nil
	case "ANY":
		return kafkarestv3.ACLOPERATION_ANY, nil
	case "ALL":
		return kafkarestv3.ACLOPERATION_ALL, nil
	case "READ":
		return kafkarestv3.ACLOPERATION_READ, nil
	case "WRITE":
		return kafkarestv3.ACLOPERATION_WRITE, nil
	case "CREATE":
		return kafkarestv3.ACLOPERATION_CREATE, nil
	case "DELETE":
		return kafkarestv3.ACLOPERATION_DELETE, nil
	case "ALTER":
		return kafkarestv3.ACLOPERATION_ALTER, nil
	case "DESCRIBE":
		return kafkarestv3.ACLOPERATION_DESCRIBE, nil
	case "CLUSTER_ACTION":
		return kafkarestv3.ACLOPERATION_CLUSTER_ACTION, nil
	case "DESCRIBE_CONFIGS":
		return kafkarestv3.ACLOPERATION_DESCRIBE_CONFIGS, nil
	case "ALTER_CONFIGS":
		return kafkarestv3.ACLOPERATION_ALTER_CONFIGS, nil
	case "IDEMPOTENT_WRITE":
		return kafkarestv3.ACLOPERATION_IDEMPOTENT_WRITE, nil
	}
	return "", fmt.Errorf("unknown ACL operation was found: %s", aclOperation)
}

func stringToAclPermission(aclPermission string) (kafkarestv3.AclPermission, error) {
	switch aclPermission {
	case "UNKNOWN":
		return kafkarestv3.ACLPERMISSION_UNKNOWN, nil
	case "ANY":
		return kafkarestv3.ACLPERMISSION_ANY, nil
	case "DENY":
		return kafkarestv3.ACLPERMISSION_DENY, nil
	case "ALLOW":
		return kafkarestv3.ACLPERMISSION_ALLOW, nil
	}
	return "", fmt.Errorf("unknown ACL permission was found: %s", aclPermission)
}

type Acl struct {
	ResourceType kafkarestv3.AclResourceType
	ResourceName string
	PatternType  kafkarestv3.AclPatternType
	Principal    string
	Host         string
	Operation    kafkarestv3.AclOperation
	Permission   kafkarestv3.AclPermission
}

type KafkaImportEnvVars struct {
	kafkaApiKey       string
	kafkaApiSecret    string
	kafkaHttpEndpoint string
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func checkEnvironmentVariablesForKafkaImportAreSet() (KafkaImportEnvVars, error) {
	kafkaApiKey := getEnv("KAFKA_API_KEY", "")
	kafkaApiSecret := getEnv("KAFKA_API_SECRET", "")
	kafkaHttpEndpoint := getEnv("KAFKA_HTTP_ENDPOINT", "")
	canImport := kafkaApiKey != "" && kafkaApiSecret != "" && kafkaHttpEndpoint != ""
	if !canImport {
		return KafkaImportEnvVars{}, fmt.Errorf("KAFKA_API_KEY, KAFKA_API_SECRET, and KAFKA_HTTP_ENDPOINT must be set for kafka topic / ACL import")
	}
	return KafkaImportEnvVars{
		kafkaApiKey:       kafkaApiKey,
		kafkaApiSecret:    kafkaApiSecret,
		kafkaHttpEndpoint: kafkaHttpEndpoint,
	}, nil
}

type KafkaRestClient struct {
	apiClient        *kafkarestv3.APIClient
	clusterId        string
	clusterApiKey    string
	clusterApiSecret string
	httpEndpoint     string
}

func (c *KafkaRestClient) apiContext(ctx context.Context) context.Context {
	if c.clusterApiKey != "" && c.clusterApiSecret != "" {
		return context.WithValue(context.Background(), kafkarestv3.ContextBasicAuth, kafkarestv3.BasicAuth{
			UserName: c.clusterApiKey,
			Password: c.clusterApiSecret,
		})
	}
	log.Printf("[WARN] Could not find cluster credentials for Confluent Cloud for clusterId=%s", c.clusterId)
	return ctx
}

type KafkaRestClientFactory struct {
	userAgent string
}

func (f KafkaRestClientFactory) CreateKafkaRestClient(httpEndpoint, clusterId, clusterApiKey, clusterApiSecret string) *KafkaRestClient {
	config := kafkarestv3.NewConfiguration()
	config.BasePath = httpEndpoint
	config.UserAgent = f.userAgent
	return &KafkaRestClient{
		apiClient:        kafkarestv3.NewAPIClient(config),
		clusterId:        clusterId,
		clusterApiKey:    clusterApiKey,
		clusterApiSecret: clusterApiSecret,
		httpEndpoint:     httpEndpoint,
	}
}

func extractStringAttributeFromListBlockOfSizeOne(d *schema.ResourceData, blockName, attributeName string) (string, error) {
	// d.Get() will return "" if the key is not present
	value := d.Get(fmt.Sprintf("%s.0.%s", blockName, attributeName)).(string)
	if value == "" {
		return "", fmt.Errorf("could not find %s attribute in %s block", attributeName, blockName)
	}
	return value, nil
}
