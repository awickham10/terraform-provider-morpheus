package morpheus

import (
	"context"

	"log"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCheckboxOptionType() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus checkbox option type resource",
		CreateContext: resourceCheckboxOptionTypeCreate,
		ReadContext:   resourceCheckboxOptionTypeRead,
		UpdateContext: resourceCheckboxOptionTypeUpdate,
		DeleteContext: resourceCheckboxOptionTypeDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the checkbox option type",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the checkbox option type",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the checkbox option type",
				Optional:    true,
				Computed:    true,
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "The organization labels associated with the option type (Only supported on Morpheus 5.5.3 or higher)",
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"field_name": {
				Type:        schema.TypeString,
				Description: "The field name of the checkbox option type",
				Required:    true,
			},
			"export_meta": {
				Type:        schema.TypeBool,
				Description: "Whether to export the checkbox option type as a tag",
				Optional:    true,
				Default:     false,
			},
			"dependent_field": {
				Type:        schema.TypeString,
				Description: "The field or code used to trigger the reloading of the field",
				Optional:    true,
				Computed:    true,
			},
			"visibility_field": {
				Type:        schema.TypeString,
				Description: "The field or code used to trigger the visibility of the field",
				Optional:    true,
				Computed:    true,
			},
			"require_field": {
				Type:        schema.TypeString,
				Description: "The field or code used to trigger the requirement of this field",
				Optional:    true,
				Computed:    true,
			},
			"show_on_edit": {
				Type:        schema.TypeBool,
				Description: "Whether the option type will display in the edit section of the provisioned resource",
				Optional:    true,
				Computed:    true,
			},
			"editable": {
				Type:        schema.TypeBool,
				Description: "Whether the value of the option type can be edited after the initial request",
				Optional:    true,
				Computed:    true,
			},
			"display_value_on_details": {
				Type:        schema.TypeBool,
				Description: "Display the selected value of the checkbox option type on the associated resource's details page",
				Optional:    true,
				Default:     false,
			},
			"field_label": {
				Type:        schema.TypeString,
				Description: "The label associated with the field in the UI",
				Required:    true,
			},
			"default_checked": {
				Type:        schema.TypeBool,
				Description: "Whether the checkbox option type is checked by default",
				Optional:    true,
				Computed:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceCheckboxOptionTypeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		for _, s := range attr.(*schema.Set).List() {
			labelsPayload = append(labelsPayload, s.(string))
		}
	}
	var defaultCheck string
	if d.Get("default_checked").(bool) {
		defaultCheck = "on"
	} else {
		defaultCheck = "off"
	}
	req := &morpheus.Request{
		Body: map[string]interface{}{
			"optionType": map[string]interface{}{
				"name":                  name,
				"description":           description,
				"labels":                labelsPayload,
				"fieldName":             d.Get("field_name").(string),
				"exportMeta":            d.Get("export_meta"),
				"dependsOnCode":         d.Get("dependent_field").(string),
				"visibleOnCode":         d.Get("visibility_field"),
				"requireOnCode":         d.Get("require_field").(string),
				"showOnEdit":            d.Get("show_on_edit").(bool),
				"editable":              d.Get("editable").(bool),
				"displayValueOnDetails": d.Get("display_value_on_details"),
				"type":                  "checkbox",
				"fieldLabel":            d.Get("field_label").(string),
				"defaultValue":          defaultCheck,
			},
		},
	}
	resp, err := client.CreateOptionType(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.CreateOptionTypeResult)
	environment := result.OptionType
	// Successfully created resource, now set id
	d.SetId(int64ToString(environment.ID))

	resourceCheckboxOptionTypeRead(ctx, d, meta)
	return diags
}

func resourceCheckboxOptionTypeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	name := d.Get("name").(string)

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindOptionTypeByName(name)
	} else if id != "" {
		resp, err = client.GetOptionType(toInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("OptionType cannot be read without name or id")
	}

	if err != nil {
		// 404 is ok?
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
	result := resp.Result.(*morpheus.GetOptionTypeResult)
	optionType := result.OptionType
	if optionType != nil {
		d.SetId(int64ToString(optionType.ID))
		d.Set("name", optionType.Name)
		d.Set("description", optionType.Description)
		d.Set("labels", optionType.Labels)
		d.Set("field_name", optionType.FieldName)
		d.Set("export_meta", optionType.ExportMeta)
		d.Set("dependent_field", optionType.DependsOnCode)
		d.Set("visibility_field", optionType.VisibleOnCode)
		d.Set("require_field", optionType.RequireOnCode)
		d.Set("show_on_edit", optionType.ShowOnEdit)
		d.Set("editable", optionType.Editable)
		d.Set("display_value_on_details", optionType.DisplayValueOnDetails)
		d.Set("field_label", optionType.FieldLabel)
		if optionType.DefaultValue == "on" {
			d.Set("default_checked", true)
		} else {
			d.Set("default_checked", true)
		}
	} else {
		return diag.Errorf("read operation: option type not found in response data") // should not happen
	}

	return diags
}

func resourceCheckboxOptionTypeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)
	id := d.Id()
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		for _, s := range attr.(*schema.Set).List() {
			labelsPayload = append(labelsPayload, s.(string))
		}
	}
	var defaultCheck string
	if d.Get("default_checked").(bool) {
		defaultCheck = "on"
	} else {
		defaultCheck = "off"
	}
	req := &morpheus.Request{
		Body: map[string]interface{}{
			"optionType": map[string]interface{}{
				"name":                  name,
				"description":           description,
				"labels":                labelsPayload,
				"fieldName":             d.Get("field_name").(string),
				"exportMeta":            d.Get("export_meta"),
				"dependsOnCode":         d.Get("dependent_field").(string),
				"visibleOnCode":         d.Get("visibility_field"),
				"requireOnCode":         d.Get("require_field").(string),
				"showOnEdit":            d.Get("show_on_edit").(bool),
				"editable":              d.Get("editable").(bool),
				"displayValueOnDetails": d.Get("display_value_on_details"),
				"type":                  "checkbox",
				"fieldLabel":            d.Get("field_label").(string),
				"defaultValue":          defaultCheck,
			},
		},
	}
	resp, err := client.UpdateOptionType(toInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.UpdateOptionTypeResult)
	optionType := result.OptionType
	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(int64ToString(optionType.ID))
	return resourceCheckboxOptionTypeRead(ctx, d, meta)
}

func resourceCheckboxOptionTypeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*morpheus.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteOptionType(toInt64(id), req)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)
			return nil
		} else {
			log.Printf("API FAILURE: %s - %s", resp, err)
			return diag.FromErr(err)
		}
	}
	log.Printf("API RESPONSE: %s", resp)
	d.SetId("")
	return diags
}
