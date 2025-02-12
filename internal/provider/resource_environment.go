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
	org "github.com/confluentinc/ccloud-sdk-go-v2/org/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"net/http"
)

func environmentResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: environmentCreate,
		ReadContext:   environmentRead,
		UpdateContext: environmentUpdate,
		DeleteContext: environmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			paramDisplayName: {
				Type:         schema.TypeString,
				Description:  "A human-readable name for the Environment.",
				ValidateFunc: validation.StringIsNotEmpty,
				Required:     true,
			},
		},
	}
}

func environmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.HasChangeExcept(paramDisplayName) {
		return diag.FromErr(fmt.Errorf(fmt.Sprintf("only %s field can be updated for an environment", paramDisplayName)))
	}

	updatedEnvironment := org.NewOrgV2Environment()
	updatedDisplayName := d.Get(paramDisplayName).(string)
	updatedEnvironment.SetDisplayName(updatedDisplayName)

	c := meta.(*Client)
	_, _, err := c.orgClient.EnvironmentsOrgV2Api.UpdateOrgV2Environment(c.orgApiContext(ctx), d.Id()).OrgV2Environment(*updatedEnvironment).Execute()

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func environmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*Client)

	displayName := d.Get(paramDisplayName).(string)
	env := org.NewOrgV2Environment()
	env.SetDisplayName(displayName)

	createdEnv, resp, err := executeEnvironmentCreate(c.orgApiContext(ctx), c, env)
	if err != nil {
		log.Printf("[ERROR] Environment create failed %v, %v, %s", env, resp, err)
		return diag.FromErr(err)
	}
	d.SetId(createdEnv.GetId())
	log.Printf("[DEBUG] Created environment %s", createdEnv.GetId())

	return nil
}

func executeEnvironmentCreate(ctx context.Context, c *Client, environment *org.OrgV2Environment) (org.OrgV2Environment, *http.Response, error) {
	req := c.orgClient.EnvironmentsOrgV2Api.CreateOrgV2Environment(c.orgApiContext(ctx)).OrgV2Environment(*environment)

	return req.Execute()
}

func environmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*Client)

	req := c.orgClient.EnvironmentsOrgV2Api.DeleteOrgV2Environment(c.orgApiContext(ctx), d.Id())
	_, err := req.Execute()

	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting environment (%s), err: %s", d.Id(), err))
	}

	return nil
}

func executeEnvironmentRead(ctx context.Context, c *Client, environmentId string) (org.OrgV2Environment, *http.Response, error) {
	req := c.orgClient.EnvironmentsOrgV2Api.GetOrgV2Environment(c.orgApiContext(ctx), environmentId)
	return req.Execute()
}

func environmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Environment read for %s", d.Id())
	c := meta.(*Client)
	environment, resp, err := executeEnvironmentRead(c.orgApiContext(ctx), c, d.Id())
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		log.Printf("[ERROR] Environment get failed for id %s, %v, %s", d.Id(), resp, err)
	}
	if err == nil {
		err = d.Set(paramDisplayName, environment.GetDisplayName())
	}
	return diag.FromErr(err)
}
