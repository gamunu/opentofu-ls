// Copyright (c) Gamunu Balagalla.
// SPDX-License-Identifier: MPL-2.0

package registry

var moduleVersionsMockResponse = `{
  "id": "v0.8.1",
  "published": "2024-08-05T17:16:42+01:00",
  "readme": true,
  "edit_link": "https://github.com/azure/terraform-azurerm-alz/blob/v0.8.1/README.md",
  "variables": {
    "architecture_name": {
      "type": "string",
      "default": null,
      "description": "The name of the architecture to create. This needs to be*.alz_architecture_definition.[json|yaml|yml] files.\n",
      "sensitive": false,
      "required": true
    }, 
    "location": {
      "type": "string",
      "default": null,
      "description": "The default location for resources in this management group. Used for policy managed identities.\n",
      "sensitive": false,
      "required": true
    }
  },
  "outputs": {
    "management_group_resource_ids": {
      "sensitive": false,
      "description": "A map of management group names to their resource ids."
    },
    "policy_assignment_resource_ids": {
      "sensitive": false,
      "description": "A map of policy assignment names to their resource ids."
    }
  },
  "schema_error": "",
  "providers": [],
  "dependencies": [
    {
      "name": "policy_assignment",
      "version_constraint": "",
      "source": "./modules/azapi_helper"
    },
    {
      "name": "policy_definitions",
      "version_constraint": "",
      "source": "./modules/azapi_helper"
    }  
  ],
  "resources": [
    {
      "address": "modtm_telemetry.telemetry",
      "type": "modtm_telemetry",
      "name": "telemetry"
    },
    {
      "address": "random_uuid.telemetry",
      "type": "random_uuid",
      "name": "telemetry"
    }  
  ],
  "link": "https://github.com/azure/terraform-azurerm-alz/tree/v0.8.1",
  "vcs_repository": "",
  "licenses": [
    {
      "spdx": "MIT",
      "confidence": 1,
      "is_compatible": true,
      "file": "LICENSE",
      "link": "https://github.com/azure/terraform-azurerm-alz/blob/v0.8.1/LICENSE"
    }
  ],
  "incompatible_license": false,
  "examples": {
    "default": {
      "readme": true,
      "edit_link": "https://github.com/azure/terraform-azurerm-alz/blob/v0.8.1/examples/default/README.md",
      "variables": {},
      "outputs": {},
      "schema_error": ""
    },
    "policy-assignment-modification-with-custom-lib": {
      "readme": true,
      "edit_link": "https://github.com/azure/terraform-azurerm-alz/blob/v0.8.1/examples/policy-assignment-modification-with-custom-lib/README.md",
      "variables": {},
      "outputs": {},
      "schema_error": ""
    }
  },
  "submodules": {
    "azapi_helper": {
      "readme": true,
      "edit_link": "https://github.com/azure/terraform-azurerm-alz/blob/v0.8.1/modules/azapi_helper/README.md",
      "variables": {
        "body": {
          "type": "dynamic",
          "default": null,
          "description": "The body object of the resource.",
          "sensitive": false,
          "required": true
        },
        "name": {
          "type": "string",
          "default": null,
          "description": "The name of resource.",
          "sensitive": false,
          "required": true
        }
      },
      "outputs": {
        "identity": {
          "sensitive": false,
          "description": "The identity configuration of the resource."
        },
        "name": {
          "sensitive": false,
          "description": "The name of the resource."
        }
      },
      "schema_error": "",
      "providers": [],
      "dependencies": [],
      "resources": [
        {
          "address": "azapi_resource.this",
          "type": "azapi_resource",
          "name": "this"
        },
        {
          "address": "terraform_data.replace_trigger",
          "type": "terraform_data",
          "name": "replace_trigger"
        }
      ]
    }
  }
}`

var moduleDataMockResponse = `{
  "addr": {
    "display": "azure/alz/azurerm",
    "namespace": "azure",
    "name": "alz",
    "target": "azurerm"
  },
  "description": "Terraform module to deploy Azure Landing Zones",
  "versions": [
    {
      "id": "v0.8.1",
      "published": "2024-08-05T17:16:42+01:00"
    },
    {
      "id": "v0.8.0",
      "published": "2024-07-10T16:41:49+01:00"
    },
    {
      "id": "v0.7.0",
      "published": "2024-07-08T17:55:56+01:00"
    },
    {
      "id": "v0.6.0",
      "published": "2024-03-13T12:31:59Z"
    },
    {
      "id": "v0.5.0",
      "published": "2024-03-08T10:29:31Z"
    },
    {
      "id": "v0.4.1",
      "published": "2023-11-06T21:35:38Z"
    },
    {
      "id": "v0.4.0",
      "published": "2023-11-06T21:26:06Z"
    },
    {
      "id": "v0.3.3",
      "published": "2023-11-02T17:22:57Z"
    },
    {
      "id": "v0.3.2",
      "published": "2023-11-02T10:44:07Z"
    },
    {
      "id": "v0.3.1",
      "published": "2023-10-26T14:08:19+01:00"
    },
    {
      "id": "v0.3.0",
      "published": "2023-10-09T20:19:28+01:00"
    },
    {
      "id": "v0.2.0",
      "published": "2023-10-06T16:56:28+01:00"
    },
    {
      "id": "v0.1.1",
      "published": "2023-08-08T17:26:02+01:00"
    },
    {
      "id": "v0.1.0",
      "published": "2023-08-08T15:43:59+01:00"
    }
  ],
  "is_blocked": false
}`
