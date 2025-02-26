package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-telegram/bot"
	"github.com/spf13/cobra"
)

var editTextCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit message",
	Long:  "Edit a text message",
	//Link the validation function to the validateArgsEdit
	Args: validateArgsEdit,
	//Link the function with the capabilities of returning an error
	RunE: editMessage,
}

func init() {
	rootCmd.AddCommand(editTextCmd)

	editTextCmd.Flags().StringP("token", "t", "", "Token from bot fathers")
	editTextCmd.Flags().IntP("chatId", "c", 0, "ID of the chat, leave blank or set 0 if you want to listen all chats")
	editTextCmd.Flags().IntP("messageId", "i", 0, "ID of the message you wan't to edit")
	editTextCmd.Flags().StringP("messageText", "m", "", "Text of the new message")
}

func editMessage(cmd *cobra.Command, args []string) error {
	token, _ := cmd.Flags().GetString("token")
	chatId, _ := cmd.Flags().GetInt("chatId")
	messageId, _ := cmd.Flags().GetInt("messageId")
	message, _ := cmd.Flags().GetString("messageText")

	//Create a context
	bgCtx := context.Background()

	//Create the bot
	tgBot, err := bot.New(token)
	if err != nil {
		return err
	}

	//Populate parameters
	parameters := &bot.EditMessageTextParams{
		ChatID:    chatId,
		MessageID: messageId,
		Text:      message,
	}

	//Delete message
	_, err = tgBot.EditMessageText(bgCtx, parameters)
	if err != nil {
		return err
	}

	//Close context
	bgCtx.Done()

	return nil
}

func validateArgsEdit(cmd *cobra.Command, args []string) error {
	//Validate the token
	token, _ := cmd.Flags().GetString("token")
	if token == "" {
		return fmt.Errorf("no token provided")
	}

	//Validate the chat ID
	chatId, _ := cmd.Flags().GetInt("chatId")
	if len(strconv.Itoa(chatId)) != 9 {
		return fmt.Errorf("wrong chat ID provided")
	}

	//Validate the text
	text, _ := cmd.Flags().GetString("messageText")
	if text == "" {
		return fmt.Errorf("no message provided")
	}

	//No need to validate the message ID

	return nil
}
