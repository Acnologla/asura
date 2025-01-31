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

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "config",
		Description: "Config asura",
		Run:         runConfig,
		Cooldown:    5,
		Options: utils.GenerateOptions(&disgord.ApplicationCommandOption{
			Type:        disgord.OptionTypeSubCommand,
			Name:        "disable",
			Description: "Desabilita as batalhas de chefe no seu servidor",
		},
			&disgord.ApplicationCommandOption{
				Type:        disgord.OptionTypeSubCommand,
				Name:        "enable",
				Description: "Habilita as batalhas de chefe no seu servidor",
			},
			&disgord.ApplicationCommandOption{
				Type:        disgord.OptionTypeSubCommand,
				Name:        "channel",
				Description: "Seta um canal para as batalhas de chefe",
				Options: utils.GenerateOptions(&disgord.ApplicationCommandOption{
					Type:        disgord.OptionTypeChannel,
					Required:    true,
					Name:        "canal",
					Description: "Canal para batalhas de chefe",
				}),
			},
			&disgord.ApplicationCommandOption{
				Type:        disgord.OptionTypeSubCommand,
				Name:        "allchannels",
				Description: "Libera o comando para todos os canais",
			},
		),
	})
}

func runConfig(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	perm, err := itc.Member.GetPermissions(ctx, handler.Client)
	guildDiscord, errGuild := handler.Client.Guild(itc.GuildID).Get()
	if errGuild != nil {
		return entities.CreateMsg().
			Content("Erro ao pegar informações do servidor").
			Res()
	}
	hasPerm := perm.Contains(disgord.PermissionManageServer) || perm.Contains(disgord.PermissionAdministrator) || guildDiscord.OwnerID == itc.Member.User.ID
	if !hasPerm || err != nil {
		return entities.CreateMsg().
			Content("Voce não tem permissão para usar esse comando").
			Res()
	}
	command := itc.Data.Options[0].Name
	guild := database.Guild.GetGuild(ctx, itc.GuildID)
	switch command {
	case "disable":
		guild.DisableLootbox = true
		database.Guild.UpdateGuild(ctx, guild)
		return entities.CreateMsg().
			Content("Lootbox foi desabilitada").
			Res()
	case "enable":
		guild.DisableLootbox = false
		database.Guild.UpdateGuild(ctx, guild)
		return entities.CreateMsg().
			Content("Lootbox foi habilitada").
			Res()
	case "channel":
		channel := itc.Data.Options[0].Options[0].Value.(string)
		guild.LootBoxChannel = disgord.ParseSnowflakeString(channel)
		database.Guild.UpdateGuild(ctx, guild)
		return entities.CreateMsg().
			Content(fmt.Sprintf("Canal <#%s> de lootbox foi setado", channel)).
			Res()
	case "allchannels":
		guild.LootBoxChannel = 0
		database.Guild.UpdateGuild(ctx, guild)
		return entities.CreateMsg().
			Content("Todos os canais estão liberados para lootbox").
			Res()
	}
	return nil
}
