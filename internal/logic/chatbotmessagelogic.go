package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"dropshipbe/common/utils"
	"dropshipbe/dropshipbe"
	"dropshipbe/internal/svc"

	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
	"github.com/zeromicro/go-zero/core/logx"
)

type ChatbotMessageLogic struct {
	ctx          context.Context
	svcCtx       *svc.ServiceContext
	openAIClient *openai.Client
	logx.Logger
}

func NewChatbotMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatbotMessageLogic {
	return &ChatbotMessageLogic{
		ctx:          ctx,
		svcCtx:       svcCtx,
		openAIClient: openai.NewClient(svcCtx.Config.OpenAI.APIKey),
		Logger:       logx.WithContext(ctx),
	}
}

func defineAITools() []openai.Tool {
	return []openai.Tool{
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "search_products",
				Description: "Search for products when customers ask general questions regarding features or recommendations.",
				Parameters: jsonschema.Definition{
					Type: jsonschema.Object,
					Properties: map[string]jsonschema.Definition{
						"semantic_query": {Type: jsonschema.String, Description: "The semantic keyword or query for searching"},
					},
					Required: []string{"semantic_query"},
				},
			},
		},
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "check_variant_stock",
				Description: "Check the stock levels and pricing for a specific product variant.",
				Parameters: jsonschema.Definition{
					Type: jsonschema.Object,
					Properties: map[string]jsonschema.Definition{
						"product_name": {Type: jsonschema.String, Description: "The name of the product"},
					},
					Required: []string{"product_name"},
				},
			},
		},
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "get_order_details",
				Description: "Retrieve order details including invoice information, total amount, and payment method.",
				Parameters: jsonschema.Definition{
					Type: jsonschema.Object,
					Properties: map[string]jsonschema.Definition{
						"order_number": {Type: jsonschema.String, Description: "The order ID/number, e.g., ORD123"},
					},
					Required: []string{"order_number"},
				},
			},
		},
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "track_shipment_status",
				Description: "Check the shipment status, tracking number, and estimated delivery date.",
				Parameters: jsonschema.Definition{
					Type: jsonschema.Object,
					Properties: map[string]jsonschema.Definition{
						"identifier": {Type: jsonschema.String, Description: "The order number or tracking ID"},
					},
					Required: []string{"identifier"},
				},
			},
		},
	}
}

func (l *ChatbotMessageLogic) dbSearchProductsSemantic(client *openai.Client, query string) (string, error) {
	req := openai.EmbeddingRequest{
		Input: []string{query},
		Model: openai.SmallEmbedding3,
	}
	resp, err := client.CreateEmbeddings(l.ctx, req)
	if err != nil {
		return "Error creating embedding.", err
	}
	vec := resp.Data[0].Embedding

	return l.svcCtx.EcommerceRepo.SearchProductsSemantic(l.ctx, vec, query)
}

func (l *ChatbotMessageLogic) dbCheckVariantStock(productName string) string {

	product, err := l.svcCtx.EcommerceRepo.GetProductByName(l.ctx, productName)
	if err != nil {
		return fmt.Sprintf("No data found for product: %s", productName)
	}

	if len(product.Variants) == 0 {
		return fmt.Sprintf("Product %s currently has no variants.", product.Name)
	}

	result := fmt.Sprintf("Stock status for %s:\n", product.Name)
	for _, v := range product.Variants {
		count, err := utils.GetInventory(l.ctx, l.svcCtx.Redis, v.ID)
		if err != nil {
			return fmt.Sprintf("Error fetching inventory for variant: %v", err)
		}
		result += fmt.Sprintf("- SKU: %s | Price: %.2f | Stock: %d units\n", v.Sku, v.Price, count)
	}
	return result
}

func (l *ChatbotMessageLogic) dbGetOrderDetails(orderNumber string) string {
	order, err := l.svcCtx.EcommerceRepo.GetOrderDetailByOrderNumber(l.ctx, orderNumber)
	if err != nil {
		return fmt.Sprintf("Order number %s not found.", orderNumber)
	}

	result := fmt.Sprintf("Order %s (Payment Status: %s):\n", order.OrderNumber, order.FinancialStatus)
	result += fmt.Sprintf("Total: %.2f %s\n", order.TotalPrice, order.Currency)

	if order.ShippingAddress != nil {
		result += "Shipping Address:\n"
		result += fmt.Sprintf("- Shipping Address: %s\n", order.ShippingAddress.String())
	} else {
		result += "Shipping Address: Not provided or not found.\n"
	}

	result += "Product Details:\n"
	for _, item := range order.OrderItems {
		result += fmt.Sprintf("- %s (Qty: %d) | Total: %.2f\n", item.ProductName, item.Quantity, item.Total)
	}

	return result
}

func (l *ChatbotMessageLogic) dbTrackShipment(identifier string) (string, error) {
	return l.svcCtx.EcommerceRepo.TrackShipment(l.ctx, identifier)
}

func (l *ChatbotMessageLogic) ChatbotMessage(in *dropshipbe.ChatbotRequest) (*dropshipbe.ChatbotResponse, error) {
	tools := defineAITools()

	systemPrompt := `You are a friendly virtual assistant for a dropshipping store.
MANDATORY RULES:
1. DO NOT fabricate information. Always use Tools to retrieve real-time data.
2. If a customer asks about an order or shipment WITHOUT providing an ID, POLITELY ask them for their order number.
3. Always format currency in GBP (£).`

	messages := []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleSystem, Content: systemPrompt},
		{Role: openai.ChatMessageRoleUser, Content: in.Message},
	}

	req := openai.ChatCompletionRequest{
		Model:    openai.GPT3Dot5Turbo,
		Messages: messages,
		Tools:    tools,
	}

	resp, err := l.openAIClient.CreateChatCompletion(l.ctx, req)
	if err != nil {
		l.Errorf("OpenAI Error: %v", err)
		return nil, err
	}

	msg := resp.Choices[0].Message

	if len(msg.ToolCalls) > 0 {
		messages = append(messages, msg)

		for _, toolCall := range msg.ToolCalls {
			var args map[string]interface{}
			json.Unmarshal([]byte(toolCall.Function.Arguments), &args)

			var dbResult string

			switch toolCall.Function.Name {
			case "search_products":
				dbResult, err = l.dbSearchProductsSemantic(l.openAIClient, args["semantic_query"].(string))
				if err != nil {
					dbResult = "Error during product search."
				}
			case "check_variant_stock":
				dbResult = l.dbCheckVariantStock(args["product_name"].(string))
			case "get_order_details":
				dbResult = l.dbGetOrderDetails(args["order_number"].(string))
			case "track_shipment_status":
				dbResult, err = l.dbTrackShipment(args["identifier"].(string))
				if err != nil {
					dbResult = "Error during shipment tracking."
				}
			}

			messages = append(messages, openai.ChatCompletionMessage{
				Role:       openai.ChatMessageRoleTool,
				Content:    dbResult,
				Name:       toolCall.Function.Name,
				ToolCallID: toolCall.ID,
			})
		}

		req.Messages = messages
		finalResp, err := l.openAIClient.CreateChatCompletion(l.ctx, req)
		if err != nil {
			return nil, err
		}

		// save to DB
		err = l.svcCtx.EcommerceRepo.SaveChatbotInteraction(l.ctx, in.Message, finalResp.Choices[0].Message.Content)
		if err != nil {
			l.Errorf("Error saving chatbot interaction: %v", err)
		}

		return &dropshipbe.ChatbotResponse{Reply: finalResp.Choices[0].Message.Content}, nil
	}

	// save to DB
	err = l.svcCtx.EcommerceRepo.SaveChatbotInteraction(l.ctx, in.Message, msg.Content)
	if err != nil {
		l.Errorf("Error saving chatbot interaction: %v", err)
	}

	return &dropshipbe.ChatbotResponse{Reply: msg.Content}, nil
}
