package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/spf13/cobra"
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send message with text or image",
	Long:  "Send a message with text or image in a chat id as bot",
	Args:  validateArgsText,
	//Link the validation function to the sendTextCmd
	RunE: sendMessage,
	//Link the function with the capabilities of returning an error
}

func init() {
	rootCmd.AddCommand(sendCmd)

	sendCmd.Flags().StringP("token", "t", "", "Token from bot fathers")
	sendCmd.Flags().IntP("chatId", "c", 0, "Your chatId")
	sendCmd.Flags().StringP("message", "m", "", "Message text to send")
	sendCmd.Flags().StringP("imagePath", "i", "", "Path of the image to send")
	sendCmd.Flags().BoolP("printMessageId", "I", false, "Print message id of your message")
	sendCmd.Flags().IntP("replyChatId", "C", 0, "Chat id you want to reply")
	sendCmd.Flags().IntP("replyMessageId", "M", 0, "Message id you want to reply")
}

func sendMessage(cmd *cobra.Command, args []string) error {
	token, _ := cmd.Flags().GetString("token")
	chatId, _ := cmd.Flags().GetInt("chatId")
	message, _ := cmd.Flags().GetString("message")
	imagePath, _ := cmd.Flags().GetString("imagePath")
	printMessageId, _ := cmd.Flags().GetBool("printMessageId")
	replyChatId, _ := cmd.Flags().GetInt("replyChatId")
	replyMessageId, _ := cmd.Flags().GetInt("replyMessageId")

	//Create a context
	bgCtx := context.Background()

	//Create the bot
	tgBot, err := bot.New(token)
	if err != nil {
		return err
	}

	var rtrn *models.Message
	replyParameters := &models.ReplyParameters{}

	//If user does no has provided the chat ID use the current one
	if replyChatId == 0 {
		replyChatId = chatId
	}

	//Fill reply parameters
	if replyMessageId != 0 {
		replyParameters.ChatID = replyChatId
		replyParameters.MessageID = replyMessageId
	}

	//Send the image
	if imagePath != "" {
		//Open image
		image, err := os.Open(imagePath)
		if err != nil {
			return err
		}

		//Create image parameters
		parameters := &bot.SendPhotoParams{
			ChatID:          chatId,
			Photo:           &models.InputFileUpload{Filename: imagePath, Data: image},
			Caption:         message,
			ReplyParameters: replyParameters,
		}

		//Send image
		rtrn, err = tgBot.SendPhoto(bgCtx, parameters)

	} else { //Send a message
		parameters := &bot.SendMessageParams{
			ChatID:          chatId,
			Text:            message,
			ReplyParameters: replyParameters,
		}

		//Send the message
		rtrn, err = tgBot.SendMessage(bgCtx, parameters)
	}

	//Check for errors
	if err != nil {
		return err
	}

	//If requested print messsage ID
	if printMessageId {
		fmt.Printf("CHAT_ID:%d\n", rtrn.ID)
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
	var hasMessage bool

	message, _ := cmd.Flags().GetString("message")
	if message == "" {
		hasMessage = false
	} else {
		hasMessage = true
	}

	//Validate the path
	var hasImage bool

	imagePath, _ := cmd.Flags().GetString("imagePath")
	if imagePath == "" {
		hasImage = false
	} else {
		//Check if the path is correct
		if _, err := os.Stat(imagePath); errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("wrong path provided")
		}

		hasImage = true
	}

	if !(hasMessage || hasImage) {
		return fmt.Errorf("provide at least one parameter between message or image")
	}

	return nil
}
