package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/spf13/cobra"
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send message with text or image",
	Long:  "Send a message in a chat as bot with text or an image",
	Args:  validateArgsSend,
	//Link the validation function to the sendTextCmd
	RunE: sendMessage,
	//Link the function with the capabilities of returning an error
}

func init() {
	rootCmd.AddCommand(sendCmd)

	sendCmd.Flags().StringP("token", "t", "", "Token from bot fathers")
	sendCmd.Flags().IntP("chatId", "c", 0, "Your chat ID")
	sendCmd.Flags().StringP("messageText", "m", "", "Message text to send")
	sendCmd.Flags().StringP("filePath", "p", "", "Path of the image/video to send")
	sendCmd.Flags().IntP("fileTimeout", "T", 60, "Timeout in seconds for sending a file")
	sendCmd.Flags().BoolP("pathIsImage", "i", false, "The path is an image to send")
	sendCmd.Flags().BoolP("pathIsVideo", "v", false, "The path is a video to send")
	sendCmd.Flags().BoolP("fileHasSpoiler", "H", false, "The file is send with hidden preview")
	sendCmd.Flags().IntP("replyChatId", "x", 0, "Chat id you want to reply")
	sendCmd.Flags().IntP("replyMessageId", "y", 0, "Message id you want to reply")
	sendCmd.Flags().BoolP("markDownV2", "2", false, "Message text is parsed in markdown v2")
	sendCmd.Flags().BoolP("printMessageId", "M", false, "Print message id of your message")
}

func sendMessage(cmd *cobra.Command, args []string) error {
	token, _ := cmd.Flags().GetString("token")
	chatId, _ := cmd.Flags().GetInt("chatId")
	messageText, _ := cmd.Flags().GetString("messageText")
	filePath, _ := cmd.Flags().GetString("filePath")
	fileTimeout, _ := cmd.Flags().GetInt("fileTimeout")
	pathIsImage, _ := cmd.Flags().GetBool("pathIsImage")
	pathIsVideo, _ := cmd.Flags().GetBool("pathIsVideo")
	fileHasSpoiler, _ := cmd.Flags().GetBool("fileHasSpoiler")
	replyChatId, _ := cmd.Flags().GetInt("replyChatId")
	replyMessageId, _ := cmd.Flags().GetInt("replyMessageId")
	markDownV2, _ := cmd.Flags().GetBool("markDownV2")
	printMessageId, _ := cmd.Flags().GetBool("printMessageId")

	//Create a context
	bgCtx := context.Background()

	//Create the bot
	tgBot, err := bot.New(token)
	if err != nil {
		return err
	}

	//Create the return message structure
	var rtrn *models.Message

	//Create and fill parsing parameters
	var parsing models.ParseMode
	if markDownV2 {
		parsing = models.ParseModeMarkdown
	}

	//If user does no has provided the chat ID use the current one
	if replyChatId == 0 {
		replyChatId = chatId
	}

	//Create and fill reply parameters
	replyParameters := &models.ReplyParameters{}
	if replyMessageId != 0 {
		replyParameters.ChatID = replyChatId
		replyParameters.MessageID = replyMessageId
	}

	//Send a file
	if filePath != "" {
		//Open image
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(fileTimeout)*time.Second)
		defer cancel()

		if pathIsImage {
			//Create image parameters
			parameters := &bot.SendPhotoParams{
				ChatID:          chatId,
				Photo:           &models.InputFileUpload{Filename: filePath, Data: file},
				Caption:         messageText,
				ReplyParameters: replyParameters,
				ParseMode:       parsing,
				HasSpoiler:      fileHasSpoiler,
			}

			//Send image
			rtrn, err = tgBot.SendPhoto(ctx, parameters)

			//Check for errors
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					return fmt.Errorf("send file request exceeded timeout of %d seconds, try a smaller file or increase -T", fileTimeout)
				}
				return err
			}
		} else if pathIsVideo {
			parameters := &bot.SendVideoParams{
				ChatID:          chatId,
				Video:           &models.InputFileUpload{Filename: filePath, Data: file},
				Caption:         messageText,
				ReplyParameters: replyParameters,
				ParseMode:       parsing,
				HasSpoiler:      fileHasSpoiler,
			}

			//Send video
			rtrn, err = tgBot.SendVideo(ctx, parameters)

			//Check for errors
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					return fmt.Errorf("send file request exceeded timeout of %d seconds, try a smaller file or increase -T", fileTimeout)
				}
				return err
			}
		}

	} else { //Send a message
		parameters := &bot.SendMessageParams{
			ChatID:          chatId,
			Text:            messageText,
			ReplyParameters: replyParameters,
			ParseMode:       parsing,
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

func validateArgsSend(cmd *cobra.Command, args []string) error {
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

	message, _ := cmd.Flags().GetString("messageText")
	if message == "" {
		hasMessage = false
	} else {
		hasMessage = true
	}

	//Validate the path
	var hasPath bool

	imagePath, _ := cmd.Flags().GetString("filePath")
	if imagePath == "" {
		hasPath = false
	} else {
		//Check if the path is correct
		if _, err := os.Stat(imagePath); errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("wrong path provided")
		}

		hasPath = true
	}

	if hasPath {
		pathIsImage, _ := cmd.Flags().GetBool("pathIsImage")
		pathIsVideo, _ := cmd.Flags().GetBool("pathIsVideo")

		if pathIsImage && pathIsVideo {
			return fmt.Errorf("images and videos are mutually exclusive")
		} else if !pathIsImage && !pathIsVideo {
			return fmt.Errorf("you need to specify if the path is a video or an image")
		}
	}

	if !(hasMessage || hasPath) {
		return fmt.Errorf("provide at least one parameter between message or file")
	}

	return nil
}
