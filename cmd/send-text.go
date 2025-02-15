package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/spf13/cobra"
)

var sendTextCmd = &cobra.Command{
	Use:   "send-text",
	Short: "Send message text",
	Long:  "Send a message in a chat id as bot",
	Args:  validateArgsText,
	//Link the validation function to the sendTextCmd
	RunE: sendMessage,
	//Link the function with the capabilities of returning an error
}

func init() {
	rootCmd.AddCommand(sendTextCmd)

	sendTextCmd.Flags().StringP("token", "t", "", "Token from bot fathers")
	sendTextCmd.Flags().IntP("chatId", "c", 0, "Your chatId")
	sendTextCmd.Flags().StringP("message", "m", "", "Message text to send")
	sendTextCmd.Flags().IntP("replyMessageId", "M", 0, "Message id you want to reply")
	sendTextCmd.Flags().IntP("replyChatId", "C", 0, "Chat id you want to reply")
}

func sendMessage(cmd *cobra.Command, args []string) error {
	token, _ := cmd.Flags().GetString("token")
	chatId, _ := cmd.Flags().GetInt("chatId")
	message, _ := cmd.Flags().GetString("message")
	printMessageId, _ := cmd.Flags().GetBool("printMessageId")
	replyMessageId, _ := cmd.Flags().GetInt("replyMessageId")
	replyChatId, _ := cmd.Flags().GetInt("replyChatId")

	//Create a context
	bgCtx := context.Background()

	//Create the bot
	tgBot, err := bot.New(token)
	if err != nil {
		return err
	}

	tgMessage := &bot.SendMessageParams{
		ChatID: chatId,
		Text:   message,
	}

	//If user does no has provided the chat ID use the current one
	if replyChatId == 0 {
		replyChatId = chatId
	}

	//Add reply parameters
	if replyMessageId != 0 {
		tgMessage.ReplyParameters = &models.ReplyParameters{}

		tgMessage.ReplyParameters.ChatID = replyChatId
		tgMessage.ReplyParameters.MessageID = replyMessageId
	}

	//Send the message
	messageRtrn, err := tgBot.SendMessage(bgCtx, tgMessage)

	//If requested print messsage ID
	if printMessageId {
		fmt.Println(messageRtrn.ID)
	}

	//Check for errors
	if err != nil {
		return err
	}

	//Close context
	bgCtx.Done()

	return nil
}

func validateArgsText(cmd *cobra.Command, args []string) error {
	//Validate the token
	token, _ := cmd.Flags().GetString("token")
	if token == "" {
		return fmt.Errorf("no token provided")
	}

	//Validate the chat ID
	chatId, _ := cmd.Flags().GetInt("chatId")
	if chatId == 0 || len(strconv.Itoa(chatId)) != 9 {
		return fmt.Errorf("wrong chat ID provided")
	}

	//Validate the message
	message, _ := cmd.Flags().GetString("message")
	if message == "" {
		return fmt.Errorf("empty message provided")
	}

	return nil
}
