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

var imageCmd = &cobra.Command{
	Use:   "send-image",
	Short: "Send image",
	Long:  "Send an image from local path in a chat id as bot",
	//Link the validation function to the imageCmd
	Args: validateArgsImage,
	//Link the function with the capabilities of returning an error
	RunE: sendImage,
}

func init() {
	rootCmd.AddCommand(imageCmd)

	imageCmd.Flags().StringP("token", "t", "", "Token from bot fathers")
	imageCmd.Flags().IntP("chatId", "c", 0, "Your chatId")
	imageCmd.Flags().StringP("message", "m", "", "Message text to send")
	imageCmd.Flags().StringP("imagePath", "p", "", "Path of the image to send")
}

func sendImage(cmd *cobra.Command, args []string) error {
	token, _ := cmd.Flags().GetString("token")
	chatId, _ := cmd.Flags().GetInt("chatId")
	message, _ := cmd.Flags().GetString("message")
	imagePath, _ := cmd.Flags().GetString("imagePath")

	//Create a context
	bgCtx := context.Background()

	//Create the bot
	tgBot, err := bot.New(token)
	if err != nil {
		return err
	}

	//Open image
	image, err := os.Open(imagePath)
	if err != nil {
		return err
	}

	//Send image
	_, err = tgBot.SendPhoto(bgCtx, &bot.SendPhotoParams{
		ChatID:  chatId,
		Photo:   &models.InputFileUpload{Filename: imagePath, Data: image},
		Caption: message,
	})

	//Check for errors
	if err != nil {
		return err
	}

	//Close context
	bgCtx.Done()

	return nil
}

func validateArgsImage(cmd *cobra.Command, args []string) error {
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

	//Validate the path
	imagePath, _ := cmd.Flags().GetString("imagePath")
	if imagePath == "" {
		return fmt.Errorf("no image path provided")
	}

	if _, err := os.Stat(imagePath); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("wrong path provided")
	}

	return nil
}
