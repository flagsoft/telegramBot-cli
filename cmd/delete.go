package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-telegram/bot"
	"github.com/spf13/cobra"
)

var deleteTextCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete message",
	Long:  "Delete a message",
	//Link the validation function to the validateArgsDelete
	Args: validateArgsDelete,
	//Link the function with the capabilities of returning an error
	RunE: deleteMessage,
}

func init() {
	rootCmd.AddCommand(deleteTextCmd)

	deleteTextCmd.Flags().StringP("token", "t", "", "Token from bot fathers")
	deleteTextCmd.Flags().IntP("chatId", "c", 0, "ID of the chat, leave blank or set 0 if you want to listen all chats")
	deleteTextCmd.Flags().IntP("messageId", "i", 0, "ID of the message you wan't to delete")
}

func deleteMessage(cmd *cobra.Command, args []string) error {
	token, _ := cmd.Flags().GetString("token")
	chatId, _ := cmd.Flags().GetInt("chatId")
	messageId, _ := cmd.Flags().GetInt("messageId")

	//Create a context
	bgCtx := context.Background()

	//Create the bot
	tgBot, err := bot.New(token)
	if err != nil {
		return err
	}

	//Populate parameters
	parameters := &bot.DeleteMessageParams{
		ChatID:    chatId,
		MessageID: messageId,
	}

	//Delete message
	res, err := tgBot.DeleteMessage(bgCtx, parameters)
	if !res && err != nil {
		return err
	}

	//Close context
	bgCtx.Done()

	return nil
}

func validateArgsDelete(cmd *cobra.Command, args []string) error {
	//Validate the token
	token, _ := cmd.Flags().GetString("token")
	if token == "" {
		return fmt.Errorf("no token provided")
	}

	//Validate the chat ID
	chatId, _ := cmd.Flags().GetInt("chatId")
	if chatId != 0 && len(strconv.Itoa(chatId)) != 9 {
		return fmt.Errorf("wrong chat ID provided")
	}

	//No need to validate the message ID

	return nil
}
