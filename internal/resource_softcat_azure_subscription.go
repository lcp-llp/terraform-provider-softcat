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
	FriendlyName  string
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
	Status                string                    `json:"status"`
	Currency              string                    `json:"currency"`
	PONumber              string                    `json:"poNumber"`
	OrderName             string                    `json:"orderName"`
	OrderID               string                    `json:"orderId"`
	Creator               *AzureSubscriptionCreator `json:"creator"`
}

type createAzureSubscriptionResponse struct {
	CreateAndOrderAzureSubscription []AzureSubscriptionOrder `json:"createAndOrderAzureSubscription"`
}

type updateAzureSubscriptionResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type updateAzureSubscriptionResponse struct {
	UpdateAzureSubscription updateAzureSubscriptionResult `json:"updateAzureSubscription"`
}

type cancelAzureSubscriptionResponse struct {
	CancelAzureSubscription updateAzureSubscriptionResult `json:"cancelAzureSubscription"`
}

type AzureSubscription struct {
	SubscriptionID string `json:"subscriptionId"`
	PlanID         string `json:"planId"`
	FriendlyName   string `json:"friendlyName"`
	Status         string `json:"status"`
	OrderID        string `json:"orderId"`
}

type getAzureSubscriptionsResponse struct {
	GetAzureSubscriptions []AzureSubscription `json:"getAzureSubscriptions"`
}

func ResourceAzureSubscription() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAzureSubscriptionCreate,
		ReadContext:   resourceAzureSubscriptionRead,
		UpdateContext: resourceAzureSubscriptionUpdate,
		DeleteContext: resourceAzureSubscriptionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceAzureSubscriptionImport,
		},
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
				Description:  "Microsoft tenant ID required by the Softcat API.",
			},
			"azure_budget": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "Budget value passed to the Softcat Azure subscription order mutation.",
			},
			"azure_contact": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "Primary contact email for the Azure subscription order.",
			},
			"friendly_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "Friendly name for the Azure subscription.",
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
			"subscription_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Azure subscription identifier returned by the follow-up subscription lookup.",
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

	if len(response.CreateAndOrderAzureSubscription) == 0 || response.CreateAndOrderAzureSubscription[0].OrderID == "" {
		return diag.FromErr(fmt.Errorf("create azure subscription: response did not include orderId"))
	}

	order := response.CreateAndOrderAzureSubscription[0]

	d.SetId(order.OrderID)

	subscription, err := lookupAzureSubscription(ctx, client, request.MsID, order.OrderID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("create azure subscription: %w", err))
	}

	if err := flattenAzureSubscriptionOrder(d, order); err != nil {
		return diag.FromErr(err)
	}

	if err := flattenAzureSubscription(d, subscription); err != nil {
		return diag.FromErr(err)
	}

	return resourceAzureSubscriptionRead(ctx, d, meta)
}


func resourceAzureSubscriptionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.Id() == "" {
		return nil
	}

	client := meta.(*Client)

	subscription, err := lookupAzureSubscription(ctx, client, d.Get("msid").(string), d.Id())
	if err != nil {
		return diag.FromErr(fmt.Errorf("read azure subscription: %w", err))
	}

	if subscription.OrderID != "" {
		if err := d.Set("order_id", subscription.OrderID); err != nil {
			return diag.FromErr(fmt.Errorf("set order_id: %w", err))
		}
	}

	if subscription.Status != "" {
		if err := d.Set("status", subscription.Status); err != nil {
			return diag.FromErr(fmt.Errorf("set status: %w", err))
		}
	}

	if err := flattenAzureSubscription(d, subscription); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceAzureSubscriptionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client)

	subscriptionID := d.Get("subscription_id").(string)
	if subscriptionID == "" {
		subscription, err := lookupAzureSubscription(ctx, client, d.Get("msid").(string), d.Id())
		if err != nil {
			return diag.FromErr(fmt.Errorf("update azure subscription: %w", err))
		}
		subscriptionID = subscription.SubscriptionID
	}

	var response updateAzureSubscriptionResponse
	if err := client.DoGraphQL(ctx, buildUpdateAzureSubscriptionMutation(
		d.Get("msid").(string),
		subscriptionID,
		d.Get("friendly_name").(string),
		d.Get("azure_budget").(string),
		d.Get("azure_contact").(string),
	), &response); err != nil {
		return diag.FromErr(fmt.Errorf("update azure subscription: %w", err))
	}

	if !response.UpdateAzureSubscription.Success {
		if response.UpdateAzureSubscription.Message == "" {
			response.UpdateAzureSubscription.Message = "update failed"
		}
		return diag.FromErr(fmt.Errorf("update azure subscription: %s", response.UpdateAzureSubscription.Message))
	}

	return resourceAzureSubscriptionRead(ctx, d, meta)
}

func resourceAzureSubscriptionImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	msID, orderID, err := parseAzureSubscriptionImportID(d.Id())
	if err != nil {
		return nil, err
	}

	if err := d.Set("msid", msID); err != nil {
		return nil, fmt.Errorf("set msid: %w", err)
	}

	if err := d.Set("order_id", orderID); err != nil {
		return nil, fmt.Errorf("set order_id: %w", err)
	}

	d.SetId(orderID)

	return []*schema.ResourceData{d}, nil
}

func resourceAzureSubscriptionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client)

	subscriptionID := d.Get("subscription_id").(string)
	if subscriptionID == "" {
		subscription, err := lookupAzureSubscription(ctx, client, d.Get("msid").(string), d.Id())
		if err != nil {
			return diag.FromErr(fmt.Errorf("cancel azure subscription: %w", err))
		}
		subscriptionID = subscription.SubscriptionID
	}

	var response cancelAzureSubscriptionResponse
	if err := client.DoGraphQL(ctx, buildCancelAzureSubscriptionMutation(
		d.Get("msid").(string),
		subscriptionID,
	), &response); err != nil {
		return diag.FromErr(fmt.Errorf("cancel azure subscription: %w", err))
	}

	if !response.CancelAzureSubscription.Success {
		if response.CancelAzureSubscription.Message == "" {
			response.CancelAzureSubscription.Message = "cancel failed"
		}
		return diag.FromErr(fmt.Errorf("cancel azure subscription: %s", response.CancelAzureSubscription.Message))
	}

	d.SetId("")
	return nil
}

func expandAzureSubscriptionRequest(d *schema.ResourceData) (AzureSubscriptionRequest, error) {
	request := AzureSubscriptionRequest{
		MsID:          d.Get("msid").(string),
		AzureBudget:   d.Get("azure_budget").(string),
		AzureContact:  d.Get("azure_contact").(string),
		FriendlyName: d.Get("friendly_name").(string),
		Quantity:     d.Get("quantity").(int),
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
		fmt.Sprintf("azureNickname: %s", graphQLString(request.FriendlyName)),
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
    status
    currency
    poNumber
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

func buildUpdateAzureSubscriptionMutation(msID string, subscriptionID string, friendlyName string, budget string, budgetContact string) string {
	return fmt.Sprintf(`mutation UpdateAzureSubscription {
	updateAzureSubscription(
		msid: %s
		subscriptionId: %s
		friendlyName: %s
		budget: %s
		budgetContact: %s
	) {
		success
		message
	}
}`,
		graphQLString(msID),
		graphQLString(subscriptionID),
		graphQLString(friendlyName),
		graphQLString(budget),
		graphQLString(budgetContact),
	)
}

func buildCancelAzureSubscriptionMutation(msID string, subscriptionID string) string {
	return fmt.Sprintf(`mutation CancelAzureSubscription {
	cancelAzureSubscription(msid: %s, subscriptionId: %s) {
		success
		message
	}
}`,
		graphQLString(msID),
		graphQLString(subscriptionID),
	)
}


func lookupAzureSubscription(ctx context.Context, client *Client, msID string, orderID string) (AzureSubscription, error) {
	if msID == "" {
		return AzureSubscription{}, fmt.Errorf("msid must not be empty")
	}

	if orderID == "" {
		return AzureSubscription{}, fmt.Errorf("orderId must not be empty")
	}

	var response getAzureSubscriptionsResponse
	if err := client.DoGraphQL(ctx, buildGetAzureSubscriptionsQuery(msID, orderID), &response); err != nil {
		return AzureSubscription{}, fmt.Errorf("get azure subscriptions: %w", err)
	}

	for _, subscription := range response.GetAzureSubscriptions {
		if subscription.SubscriptionID != "" {
			return subscription, nil
		}
	}

	return AzureSubscription{}, fmt.Errorf("get azure subscriptions: response did not include subscriptionId for orderId %q", orderID)
}

func buildGetAzureSubscriptionsQuery(msID string, orderID string) string {
	return fmt.Sprintf(`query GetAzureSubscriptions {
	getAzureSubscriptions(msid: %s, orderId: %s) {
    subscriptionId
    planId
    friendlyName
    status
    orderId
    creationDate
    effectiveStartDate
    billingCycle
    commitmentEndDate
    autoRenewEnabled
    budget
    budgetContact
  }
}`,
		graphQLString(msID),
		graphQLString(orderID))
}

func parseAzureSubscriptionImportID(value string) (string, string, error) {
	parts := strings.SplitN(value, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("unexpected import ID format %q, expected msid/order_id", value)
	}

	return parts[0], parts[1], nil
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

func flattenAzureSubscription(d *schema.ResourceData, subscription AzureSubscription) error {
	if err := d.Set("friendly_name", subscription.FriendlyName); err != nil {
		return fmt.Errorf("set friendly_name: %w", err)
	}

	if err := d.Set("subscription_id", subscription.SubscriptionID); err != nil {
		return fmt.Errorf("set subscription_id: %w", err)
	}

	return nil
}

func graphQLString(value string) string {
	return fmt.Sprintf("%q", value)
}
