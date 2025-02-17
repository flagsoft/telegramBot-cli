package cmd

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/spf13/cobra"
)

var receiveTextCmd = &cobra.Command{
	Use:   "receive-text",
	Short: "Receive message text",
	Long:  "Receive a message in a chat id as bot with the pattern below\nDATA | CHAT_ID | MESSAGE_ID | MESSAGE",
	//Link the validation function to the receiveTextCmd
	Args: validateArgsReceiveText,
	//Link the function with the capabilities of returning an error
	RunE: receiveMessage,
}

func init() {
	rootCmd.AddCommand(receiveTextCmd)

	receiveTextCmd.Flags().StringP("token", "t", "", "Token from bot fathers")
	receiveTextCmd.Flags().IntP("chatId", "c", 0, "Your chatId, leave blank or set 0 if you want to listen all chats")
	receiveTextCmd.Flags().IntP("messageCounter", "n", 0, "Numer of messages, leave blank or set 0 for continuous receiving")
	receiveTextCmd.Flags().BoolP("wantSync", "s", false, "Sync old messages sended while the bot was not running")
	receiveTextCmd.Flags().BoolP("wantChatId", "C", false, "Print the chat ID")
	receiveTextCmd.Flags().BoolP("wantTimestamp", "U", false, "Print the UNIX datetime")
	receiveTextCmd.Flags().BoolP("wantTimestampHuman", "H", false, "Print the datetime human readable")
	receiveTextCmd.Flags().BoolP("wantMessageId", "M", false, "Print the message ID of each message")
}

func receiveMessage(cmd *cobra.Command, args []string) error {
	token, _ := cmd.Flags().GetString("token")
	chatId, _ := cmd.Flags().GetInt("chatId")
	maxMessages, _ := cmd.Flags().GetInt("messageCounter")
	wantChatId, _ := cmd.Flags().GetBool("wantChatId")
	sync, _ := cmd.Flags().GetBool("wantSync")
	wantTimestamp, _ := cmd.Flags().GetBool("wantTimestamp")
	wantTimestampHuman, _ := cmd.Flags().GetBool("wantTimestampHuman")
	wantMessageId, _ := cmd.Flags().GetBool("wantMessageId")
	counter := 0

	//Create a context
	bgCtx, cancel := context.WithCancel(context.Background())

	//Create the handler
	defaultHandler := func(ctx context.Context, tgBot *bot.Bot, update *models.Update, cancelFunc context.CancelFunc) {
		//Handle only messages
		if update.Message != nil {

			if int64(update.Message.Date) < time.Now().Unix() && !sync {
				return
			}

			//Listen only for the specified chat ID
			if update.Message.Chat.ID == int64(chatId) || chatId == 0 {

				outputMessage := ""

				//Append the Date and Time
				if wantTimestampHuman {
					outputMessage += fmt.Sprintf("DATE:%s|", time.Unix(int64(update.Message.Date), 0))
				} else if wantTimestamp {
					outputMessage += fmt.Sprintf("DATE:%d|", update.Message.Date)
				}

				//Append Chat ID
				if wantChatId {
					outputMessage += fmt.Sprintf("CHAT_ID:%d|", update.Message.Chat.ID)
				}

				//Append Message ID
				if wantMessageId {
					outputMessage += fmt.Sprintf("MESSAGE_ID:%d|", update.Message.ID)
				}

				//Append message
				outputMessage += update.Message.Text

				//Print out complete message
				fmt.Println(outputMessage)

				//Increase the counter only if user want a cuntdown
				if maxMessages != 0 {
					counter++

				}

				//Check if counter has reach the user value
				if counter > maxMessages {
					//Close the bot
					tgBot.Close(ctx)

					//Cancel the Context
					cancelFunc()
				}

			}
		}
	}

	opts := []bot.Option{
		//Link the handler to the bot
		bot.WithDefaultHandler(func(ctx context.Context, b *bot.Bot, update *models.Update) {
			//Pass the param from the default handler + the context cancellation function
			defaultHandler(ctx, b, update, cancel)
		}),
	}

	//Create the bot
	tgBot, err := bot.New(token, opts...)
	if nil != err {
		return (err)
	}

	//Start the bot
	tgBot.Start(bgCtx)

	//Close context
	bgCtx.Done()

	return nil
}

func validateArgsReceiveText(cmd *cobra.Command, args []string) error {
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

	//No need to validate the messageCounter

	return nil
}
