package discord

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"

	"github.com/mtzvd/ironroll/core/roll"
)

// Command defines the /ironroll slash command.
var Command = &discordgo.ApplicationCommand{
	Name:        "ironroll",
	Description: "Perform an Ironsworn roll",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "modifier",
			Description: "Optional action modifier (Z)",
			Required:    false,
		},
	},
}

// HandleInteraction handles the /ironroll command interaction.
//
// This handler is stateless and performs a single roll per invocation.
func HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.ApplicationCommandData().Name != "ironroll" {
		return
	}

	modifier := 0
	for _, opt := range i.ApplicationCommandData().Options {
		if opt.Name == "modifier" {
			modifier = int(opt.IntValue())
		}
	}

	result := roll.Roll(modifier)
	content := formatResult(result)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
	if err != nil {
		slog.Error("discord interaction response failed", "error", err)
	}
}
