package commands

import (
	"asura/src/database"
	"asura/src/entities"
	"asura/src/handler"
	"asura/src/utils"
	"context"
	"fmt"

	"github.com/andersfylling/disgord"
)

const TRANSACTIONS_PER_PAGE = 9

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "transactions",
		Description: "Veja suas transaçoes em dinheiro",
		Run:         runTransactions,
		Cooldown:    5,
		Category:    handler.Profile,
	})
}

func generateTransactionsEmbed(chunkedTransactions [][]*entities.Transaction, page int) *disgord.Embed {
	transactions := chunkedTransactions[page]
	text := ""
	for _, transaction := range transactions {
		transactionAuthor, _ := handler.Client.User(transaction.AuthorID).Get()
		transactionUsername := "Desconhecido"
		if transactionAuthor != nil {
			transactionUsername = transactionAuthor.Username
		}
		transactionDate := utils.FormatDate(transaction.CreatedAt)
		text += fmt.Sprintf("🗓️ Data: %s\nAutor: %s (**%d** 💰)\n\n", transactionDate, transactionUsername, transaction.Amount)
	}
	return &disgord.Embed{
		Color:       65535,
		Title:       fmt.Sprintf("Transaçoes (%d/%d)", page+1, len(chunkedTransactions)),
		Description: text,
	}
}

func runTransactions(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	user := database.User.GetUser(ctx, itc.Member.UserID, "Transactions")
	chunkedTransactions := utils.Chunk(user.Transactions, TRANSACTIONS_PER_PAGE)
	currentPage := 0
	components := []*disgord.MessageComponent{
		{
			Type: disgord.MessageComponentActionRow,
			Components: []*disgord.MessageComponent{
				{
					Type:  disgord.MessageComponentButton,
					Style: disgord.Primary,
					Emoji: &disgord.Emoji{
						Name: "⬅️",
					},
					Label:    "\u200f",
					CustomID: "back",
				},
				{
					Type:  disgord.MessageComponentButton,
					Style: disgord.Primary,
					Emoji: &disgord.Emoji{
						Name: "➡️",
					},
					Label:    "\u200f",
					CustomID: "next",
				},
			},
		},
	}

	handler.Client.SendInteractionResponse(ctx, itc, &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.CreateInteractionResponseData{
			Embeds: []*disgord.Embed{
				generateTransactionsEmbed(chunkedTransactions, currentPage),
			},
			Components: components,
		},
	})
	handler.RegisterHandler(itc.ID, func(interaction *disgord.InteractionCreate) {
		if itc.Member.User.ID == interaction.Member.User.ID {
			if interaction.Data.CustomID == "back" {
				if currentPage == 0 {
					currentPage = len(chunkedTransactions) - 1
				} else {
					currentPage--
				}
			} else {
				if currentPage == len(chunkedTransactions)-1 {
					currentPage = 0
				} else {
					currentPage++
				}
			}

			handler.Client.SendInteractionResponse(ctx, interaction, &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackUpdateMessage,
				Data: &disgord.CreateInteractionResponseData{
					Embeds: []*disgord.Embed{
						generateTransactionsEmbed(chunkedTransactions, currentPage),
					},
					Components: components,
				},
			})
		}
	}, 60*2)
	return nil
}
