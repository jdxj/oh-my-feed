package bot

import (
	"context"
	"errors"
	"fmt"
	"strings"

	tbi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"

	"github.com/jdxj/oh-my-feed/log"
	"github.com/jdxj/oh-my-feed/model"
)

var (
	commands = []*command{
		newHelloCmd(),
		newTestInlineKeyboardCmd(),
		newAddFeedCmd(),
	}

	cmdSlice = func() []tbi.BotCommand {
		var bcs []tbi.BotCommand
		for _, v := range commands {
			bcs = append(bcs, tbi.BotCommand{
				Command:     v.name,
				Description: v.description,
			})
		}
		return bcs
	}()

	cmdMap = func() map[string]*command {
		m := make(map[string]*command)
		for _, v := range commands {
			_, ok := m[v.name]
			if ok {
				log.Fatalf("duplicated: %s", v.name)
			} else {
				m[v.name] = v
			}
		}
		return m
	}()
)

type handler func(args []string, update tbi.Update) tbi.Chattable

type command struct {
	name        string
	description string
	h           handler
}

func registerCmd() {
	deleteCmdReq := tbi.NewDeleteMyCommands()
	deleteCmdRsp, err := client.Request(deleteCmdReq)
	if err != nil {
		log.Warnf("request delete cmd err: %s", err)
	} else {
		if !deleteCmdRsp.Ok {
			log.Desugar().Warn(
				"delete-cmd",
				zap.String("description", deleteCmdRsp.Description),
			)
		}
	}

	setCmdReq := tbi.NewSetMyCommands(cmdSlice...)
	setCmdRsp, err := client.Request(setCmdReq)
	if err != nil {
		log.Fatalf("request set commands err: %s", err)
	}
	if !setCmdRsp.Ok {
		log.Desugar().Fatal(
			"set-cmd",
			zap.String("description", setCmdRsp.Description),
		)
	}
}

// todo: 优雅退出
func handlers(updates tbi.UpdatesChannel) {
	go func() {
		// todo: update并发测试
		for update := range updates {
			if update.Message != nil {
				txt := update.Message.Text
				cli, err := parseCmdLine(txt)
				if err != nil {
					switch {
					case errors.Is(err, ErrNotCmd):
						log.Infof("receive msg: %s", txt)
					case errors.Is(err, ErrCmdNotFound):
						log.Warnf("not register cmd: %s", txt)
					default:
						log.Errorf("parse cmdline err: %s", err)
					}
					continue
				}

				msg := cli.cmd.h(cli.args, update)
				_, err = client.Send(msg)
				if err != nil {
					log.Warnf("send msg err: %s", err)
				}
			}
		}

		log.Infof("quit range updates")
	}()
}

type cmdLine struct {
	cmd  *command
	args []string
}

var (
	ErrNotCmd      = errors.New("not cmd")
	ErrCmdNotFound = errors.New("cmd not found")
)

func parseCmdLine(str string) (*cmdLine, error) {
	if !strings.HasPrefix(str, "/") {
		return nil, ErrNotCmd
	}
	str = strings.TrimPrefix(str, "/")

	args := strings.Split(str, " ")
	var noSpaceArgs []string
	for _, v := range args {
		str = strings.TrimSpace(v)
		if str == "" {
			continue
		}
		noSpaceArgs = append(noSpaceArgs, str)
	}
	args = noSpaceArgs

	if len(args) == 0 {
		return nil, ErrCmdNotFound
	}
	cmd, ok := cmdMap[args[0]]
	if !ok {
		return nil, ErrCmdNotFound
	}

	return &cmdLine{
		cmd:  cmd,
		args: args[1:],
	}, nil
}

func newHelloCmd() *command {
	return &command{
		name:        "hello",
		description: "say hello",
		h: func(args []string, update tbi.Update) tbi.Chattable {
			txt := "world"
			if len(args) > 0 {
				txt = fmt.Sprintf("%s的天呐! 竟然有%s!", txt, args[0])
			}
			return tbi.NewMessage(update.Message.Chat.ID, txt)
		},
	}
}

func newTestInlineKeyboardCmd() *command {
	return &command{
		name:        "test_inline_keyboard",
		description: "测试",
		h:           testInlineKeyboard,
	}
}

func testInlineKeyboard(args []string, update tbi.Update) tbi.Chattable {
	row1 := tbi.NewInlineKeyboardRow(
		tbi.NewInlineKeyboardButtonData("abc", "123"),
		tbi.NewInlineKeyboardButtonData("def", "456"),
	)
	row2 := tbi.NewInlineKeyboardRow(
		tbi.NewInlineKeyboardButtonData("cba", "321"),
		tbi.NewInlineKeyboardButtonData("fed", "654"),
	)
	inlineKeyboard := tbi.NewInlineKeyboardMarkup(row1, row2)

	msg := tbi.NewMessage(update.Message.Chat.ID, "hhh")
	msg.ReplyMarkup = inlineKeyboard
	return msg
}

func newAddFeedCmd() *command {
	return &command{
		name:        "addfeed",
		description: "添加订阅",
		h: func(args []string, update tbi.Update) tbi.Chattable {
			chatID := update.Message.Chat.ID
			msg := tbi.NewMessage(chatID, "")

			if len(args) == 0 {
				msg.Text = "需要指定一个订阅地址"
				return msg
			}

			err := model.AddUserFeed(context.TODO(), chatID, args[0])
			if err != nil {
				log.Desugar().Error(
					"add-feed",
					zap.String("feed", args[0]),
					zap.String("err", err.Error()),
				)
				msg.Text = "添加订阅地址失败"
			} else {
				msg.Text = "添加订阅地址成功"
			}
			msg.ReplyToMessageID = update.Message.MessageID
			return msg
		},
	}
}
