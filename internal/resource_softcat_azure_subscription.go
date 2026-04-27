package internal

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type AzureSubscriptionRequest struct {
	BasketName    *string
	MsID          string
	AzureBudget   string
	AzureContact  string
	AzureNickname string
	CheckoutData  AzureSubscriptionCheckoutData
	Quantity      int
}

type AzureSubscriptionCheckoutData struct {
	PurchaseOrderNumber   string
	AdditionalInformation string
	CSPTerms              bool
}

type AzureSubscriptionCreator struct {
	UserID  string `json:"userId"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Account string `json:"account"`
}

type AzureSubscriptionOrder struct {
	DateCreated           string                    `json:"dateCreated"`
	DateStored            string                    `json:"dateStored"`
	ItemCount             int                       `json:"itemCount"`
	Status                string                    `json:"status"`
	BasketTotalIncVat     string                    `json:"basketTotalIncVat"`
	BasketTotalExcVat     string                    `json:"basketTotalExcVat"`
	ProductsTotalIncVat   string                    `json:"productsTotalIncVat"`
	ProductsTotalExcVat   string                    `json:"productsTotalExcVat"`
	ShippingIncVat        string                    `json:"shippingIncVat"`
	ShippingExcVat        string                    `json:"shippingExcVat"`
	Currency              string                    `json:"currency"`
	DeliveryContact       string                    `json:"deliveryContact"`
	DeliveryContactPhone  string                    `json:"deliveryContactPhone"`
	DeliveryContactEmail  string                    `json:"deliveryContactEmail"`
	PONumber              string                    `json:"poNumber"`
	PaymentMethod         string                    `json:"paymentMethod"`
	StoreCollection       bool                      `json:"storeCollection"`
	Collected             bool                      `json:"collected"`
	PurchasingInstruction string                    `json:"purchasingInstruction"`
	OrderName             string                    `json:"orderName"`
	OrderID               string                    `json:"orderId"`
	Creator               *AzureSubscriptionCreator `json:"creator"`
}

type createAzureSubscriptionResponse struct {
	CreateAndOrderAzureSubscription AzureSubscriptionOrder `json:"createAndOrderAzureSubscription"`
}

func ResourceAzureSubscription() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAzureSubscriptionCreate,
		ReadContext:   resourceAzureSubscriptionRead,
		DeleteContext: resourceAzureSubscriptionDelete,
		Schema: map[string]*schema.Schema{
			"basket_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "Optional basket name used when placing the Azure subscription order.",
			},
			"msid": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "Microsoft tenant or account identifier required by the Softcat API.",
			},
			"azure_budget": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "Budget value passed to the Softcat Azure subscription order mutation.",
			},
			"azure_contact": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "Primary contact email for the Azure subscription order.",
			},
			"azure_nickname": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "Friendly name assigned to the Azure subscription order.",
			},
			"quantity": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(1),
				Description:  "Quantity passed to the order mutation. Softcat currently expects a positive integer.",
			},
			"checkout_data": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				MinItems:    1,
				MaxItems:    1,
				Description: "Checkout metadata required by the createAndOrderAzureSubscription mutation.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"purchase_order_number": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
							Description:  "Purchase order number to include in the Softcat checkout payload.",
						},
						"additional_information": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Optional free-form information sent with the checkout data.",
						},
						"csp_terms": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Confirms acceptance of CSP terms for the order.",
						},
					},
				},
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Order status returned by the Softcat API.",
			},
			"order_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Softcat order identifier returned by the mutation.",
			},
			"order_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Order name returned by the Softcat API.",
			},
			"currency": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Currency associated with the created order.",
			},
			"po_number": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Purchase order number echoed back by the API.",
			},
			"date_created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp when the order was created.",
			},
			"date_stored": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp when the order was stored by the API.",
			},
			"creator": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Creator metadata returned by the Softcat API.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"email": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"account": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceAzureSubscriptionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client)

	request, err := expandAzureSubscriptionRequest(d)
	if err != nil {
		return diag.FromErr(err)
	}

	var response createAzureSubscriptionResponse
	if err := client.DoGraphQL(ctx, buildCreateAzureSubscriptionMutation(request), &response); err != nil {
		return diag.FromErr(fmt.Errorf("create azure subscription: %w", err))
	}

	if response.CreateAndOrderAzureSubscription.OrderID == "" {
		return diag.FromErr(fmt.Errorf("create azure subscription: response did not include orderId"))
	}

	d.SetId(response.CreateAndOrderAzureSubscription.OrderID)

	if err := flattenAzureSubscriptionOrder(d, response.CreateAndOrderAzureSubscription); err != nil {
		return diag.FromErr(err)
	}

	return resourceAzureSubscriptionRead(ctx, d, meta)
}

func resourceAzureSubscriptionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAzureSubscriptionDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	d.SetId("")

	return diag.Diagnostics{{
		Severity: diag.Warning,
		Summary:  "Resource removed from state only",
		Detail:   "The Softcat API integration does not yet expose a delete operation for Azure subscription orders, so Terraform can only forget this resource locally.",
	}}
}

func expandAzureSubscriptionRequest(d *schema.ResourceData) (AzureSubscriptionRequest, error) {
	request := AzureSubscriptionRequest{
		MsID:          d.Get("msid").(string),
		AzureBudget:   d.Get("azure_budget").(string),
		AzureContact:  d.Get("azure_contact").(string),
		AzureNickname: d.Get("azure_nickname").(string),
		Quantity:      d.Get("quantity").(int),
	}

	if basketName, ok := d.GetOk("basket_name"); ok {
		value := basketName.(string)
		request.BasketName = &value
	}

	checkoutData := d.Get("checkout_data").([]interface{})
	if len(checkoutData) != 1 {
		return AzureSubscriptionRequest{}, fmt.Errorf("checkout_data must contain exactly one item")
	}

	checkoutMap := checkoutData[0].(map[string]interface{})
	request.CheckoutData = AzureSubscriptionCheckoutData{
		PurchaseOrderNumber:   checkoutMap["purchase_order_number"].(string),
		AdditionalInformation: checkoutMap["additional_information"].(string),
		CSPTerms:              checkoutMap["csp_terms"].(bool),
	}

	return request, nil
}

func buildCreateAzureSubscriptionMutation(request AzureSubscriptionRequest) string {
	arguments := []string{
		fmt.Sprintf("checkoutData: { purchaseOrderNumber: %s, additionalInformation: %s, cspTerms: %t }", graphQLString(request.CheckoutData.PurchaseOrderNumber), graphQLString(request.CheckoutData.AdditionalInformation), request.CheckoutData.CSPTerms),
		fmt.Sprintf("msid: %s", graphQLString(request.MsID)),
		fmt.Sprintf("azureBudget: %s", graphQLString(request.AzureBudget)),
		fmt.Sprintf("azureNickname: %s", graphQLString(request.AzureNickname)),
		fmt.Sprintf("azureContact: %s", graphQLString(request.AzureContact)),
		fmt.Sprintf("quantity: %d", request.Quantity),
	}

	if request.BasketName != nil {
		arguments = append(arguments, fmt.Sprintf("basketName: %s", graphQLString(*request.BasketName)))
	}

	return fmt.Sprintf(`mutation CreateAndOrderAzureSubscription {
  createAndOrderAzureSubscription(
    %s
  ) {
    dateCreated
    dateStored
    itemCount
    status
    basketTotalIncVat
    basketTotalExcVat
    productsTotalIncVat
    productsTotalExcVat
    shippingIncVat
    shippingExcVat
    currency
    deliveryContact
    deliveryContactPhone
    deliveryContactEmail
    poNumber
    paymentMethod
    storeCollection
    collected
    purchasingInstruction
    orderName
    orderId
    creator {
      userId
      name
      email
      account
    }
  }
}`,
		strings.Join(arguments, "\n    "))
}

func flattenAzureSubscriptionOrder(d *schema.ResourceData, order AzureSubscriptionOrder) error {
	if err := d.Set("status", order.Status); err != nil {
		return fmt.Errorf("set status: %w", err)
	}
	if err := d.Set("order_id", order.OrderID); err != nil {
		return fmt.Errorf("set order_id: %w", err)
	}
	if err := d.Set("order_name", order.OrderName); err != nil {
		return fmt.Errorf("set order_name: %w", err)
	}
	if err := d.Set("currency", order.Currency); err != nil {
		return fmt.Errorf("set currency: %w", err)
	}
	if err := d.Set("po_number", order.PONumber); err != nil {
		return fmt.Errorf("set po_number: %w", err)
	}
	if err := d.Set("date_created", order.DateCreated); err != nil {
		return fmt.Errorf("set date_created: %w", err)
	}
	if err := d.Set("date_stored", order.DateStored); err != nil {
		return fmt.Errorf("set date_stored: %w", err)
	}

	if order.Creator == nil {
		if err := d.Set("creator", nil); err != nil {
			return fmt.Errorf("set creator: %w", err)
		}
		return nil
	}

	if err := d.Set("creator", []interface{}{map[string]interface{}{
		"user_id": order.Creator.UserID,
		"name":    order.Creator.Name,
		"email":   order.Creator.Email,
		"account": order.Creator.Account,
	}}); err != nil {
		return fmt.Errorf("set creator: %w", err)
	}

	return nil
}

func graphQLString(value string) string {
	return fmt.Sprintf("%q", value)
}
