package morpheus

import (
	"context"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceApiOptionList() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus api option list resource.",
		CreateContext: resourceApiOptionListCreate,
		ReadContext:   resourceApiOptionListRead,
		UpdateContext: resourceApiOptionListUpdate,
		DeleteContext: resourceApiOptionListDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the api option list",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the option list",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the option list",
				Optional:    true,
			},
			"visibility": {
				Type:         schema.TypeString,
				Description:  "Whether the option list is visible in sub-tenants or not",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"private", "public", ""}, false),
				Default:      "private",
			},
			"option_list": {
				Type:         schema.TypeString,
				Description:  "The Morpheus object option list",
				ValidateFunc: validation.StringInSlice([]string{"clouds", "instanceTypeClouds", "environments", "groups", "instances", "instance-wiki", "networks", "instanceNetworks", "servicePlans", "resourcePools", "securityGroups", "servers", "server-wiki"}, false),
				Optional:     true,
				Computed:     true,
			},
			"translation_script": {
				Type:        schema.TypeString,
				Description: "A js script to translate the result data object into an Array containing objects with properties 'name’ and 'value’.",
				Optional:    true,
				Computed:    true,
			},
			"request_script": {
				Type:        schema.TypeString,
				Description: "A js script to manipulate the request payload.",
				Optional:    true,
				Computed:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceApiOptionListCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	req := &morpheus.Request{
		Body: map[string]interface{}{
			"optionTypeList": map[string]interface{}{
				"name":              name,
				"description":       description,
				"type":              "api",
				"apiType":           d.Get("option_list").(string),
				"visibility":        d.Get("visibility").(string),
				"translationScript": d.Get("translation_script").(string),
				"requestScript":     d.Get("request_script").(string),
			},
		},
	}
	resp, err := client.CreateOptionList(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.CreateOptionListResult)
	optionList := result.OptionList
	// Successfully created resource, now set id
	d.SetId(int64ToString(optionList.ID))

	resourceApiOptionListRead(ctx, d, meta)
	return diags
}

func resourceApiOptionListRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindOptionListByName(name)
	} else if id != "" {
		resp, err = client.GetOptionList(toInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Option list cannot be read without name or id")
	}

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)
			return diag.FromErr(err)
		} else {
			log.Printf("API FAILURE: %s - %s", resp, err)
			return diag.FromErr(err)
		}
	}
	log.Printf("API RESPONSE: %s", resp)

	// store resource data
	result := resp.Result.(*morpheus.GetOptionListResult)
	optionList := result.OptionList
	if optionList != nil {
		d.SetId(int64ToString(optionList.ID))
		d.Set("name", optionList.Name)
		d.Set("description", optionList.Description)
		d.Set("visibility", optionList.Visibility)
		d.Set("option_list", optionList.APIType)
		d.Set("translation_script", optionList.TranslationScript)
		d.Set("request_script", optionList.RequestScript)
	} else {
		log.Println(optionList)
		return diag.Errorf("read operation: option list not found in response data") // should not happen
	}

	return diags
}

func resourceApiOptionListUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()
	name := d.Get("name").(string)
	description := d.Get("description").(string)

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"optionTypeList": map[string]interface{}{
				"name":              name,
				"description":       description,
				"type":              "api",
				"apiType":           d.Get("option_list").(string),
				"visibility":        d.Get("visibility").(string),
				"translationScript": d.Get("translation_script").(string),
				"requestScript":     d.Get("request_script").(string),
			},
		},
	}
	resp, err := client.UpdateOptionList(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.UpdateOptionListResult)
	optionList := result.OptionList
	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(int64ToString(optionList.ID))
	return resourceApiOptionListRead(ctx, d, meta)
}

func resourceApiOptionListDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteOptionList(toInt64(id), req)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)
			return diag.FromErr(err)
		} else {
			log.Printf("API FAILURE: %s - %s", resp, err)
			return diag.FromErr(err)
		}
	}
	log.Printf("API RESPONSE: %s", resp)
	d.SetId("")
	return diags
}
