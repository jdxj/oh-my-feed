package bot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	tbi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/jdxj/oh-my-feed/internal/app/model"
	"github.com/jdxj/oh-my-feed/internal/app/task"
	"github.com/jdxj/oh-my-feed/internal/pkg/config"
	"github.com/jdxj/oh-my-feed/internal/pkg/log"
)

type (
	handler func(args []string, update tbi.Update) tbi.Chattable
	command struct {
		bc  tbi.BotCommand
		bcs tbi.BotCommandScope
		h   handler
	}
)

var (
	cmdMap map[string]*command
)

func initCmd() {
	commands := []*command{
		newHelloCmd(),
		newSubscribeCmd(),
		newUnsubscribeCmd(),
		newIntervalCmd(),
	}

	cmdGroup := make(map[tbi.BotCommandScope][]tbi.BotCommand)
	for _, v := range commands {
		cmdGroup[v.bcs] = append(cmdGroup[v.bcs], v.bc)
	}

	cmdMap = make(map[string]*command)
	for _, v := range commands {
		_, ok := cmdMap[v.bc.Command]
		if ok {
			log.Fatalf("duplicated: %s", v.bc.Command)
		} else {
			cmdMap[v.bc.Command] = v
		}
	}

	registerCmd(cmdGroup)
}

func registerCmd(cmdGroup map[tbi.BotCommandScope][]tbi.BotCommand) {
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

	for scope, cmd := range cmdGroup {
		setCmdReq := tbi.NewSetMyCommandsWithScope(scope, cmd...)
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
}

func handlers(updates tbi.UpdatesChannel) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		for update := range updates {
			// todo: 用于调试, 应该删除
			data, err := json.MarshalIndent(update, "", "  ")
			if err != nil {
				log.Errorf("marshal update err: %s", err)
			} else {
				log.Debugf("marshaled update: %s", data)
			}

			select {
			case <-stop:
				log.Infof("stop handle update")
				// todo: 打印剩余的update?
				return
			default:
			}

			update := update
			// todo: 测试 submit 同时close
			_ = gp.Submit(func() {
				// todo: 不合理的判断, 应该结合上下文
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
						return
					}

					msg := cli.cmd.h(cli.args, update)
					_, err = client.Send(msg)
					if err != nil {
						log.Warnf("send msg err: %s", err)
					}
				}
			})
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
		bc: tbi.BotCommand{
			Command:     "hello",
			Description: "say hello",
		},
		bcs: tbi.NewBotCommandScopeDefault(),
		h: func(args []string, update tbi.Update) tbi.Chattable {
			txt := "world"
			if len(args) > 0 {
				txt = fmt.Sprintf("%s的天呐! 竟然有%s!", txt, args[0])
			}
			return tbi.NewMessage(update.Message.Chat.ID, txt)
		},
	}
}

func newSubscribeCmd() *command {
	return &command{
		bc: tbi.BotCommand{
			Command:     "subscribe",
			Description: "订阅",
		},
		bcs: tbi.NewBotCommandScopeDefault(),
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
					zap.Error(err),
				)
				msg.Text = "订阅失败"
			} else {
				msg.Text = "订阅成功"
			}
			msg.ReplyToMessageID = update.Message.MessageID
			return msg
		},
	}
}

func newUnsubscribeCmd() *command {
	return &command{
		bc: tbi.BotCommand{
			Command:     "unsubscribe",
			Description: "退订",
		},
		bcs: tbi.NewBotCommandScopeDefault(),
		h: func(args []string, update tbi.Update) tbi.Chattable {
			chatID := update.Message.Chat.ID
			msg := tbi.NewMessage(chatID, "")

			if len(args) == 0 {
				msg.Text = "需要指定一个订阅地址"
				return msg
			}

			err := model.DelUserFeed(context.TODO(), chatID, args[0])
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					msg.Text = "没有该订阅"
				} else {
					log.Errorf("del user feed err: %s", err)
					msg.Text = "退订失败"
				}
			} else {
				msg.Text = "退订成功"
			}
			msg.ReplyToMessageID = update.Message.MessageID
			return msg
		},
	}
}

func newIntervalCmd() *command {
	return &command{
		bc: tbi.BotCommand{
			Command:     "interval",
			Description: "更新间隔",
		},
		bcs: tbi.NewBotCommandScopeChat(config.Telegram.Owner),
		h: func(args []string, update tbi.Update) tbi.Chattable {
			chatID := update.Message.Chat.ID
			msg := tbi.NewMessage(chatID, "")

			if len(args) == 0 {
				msg.Text = "需要指定一个延时 e.g., 1h, 10m"
				return msg
			}

			err := task.SetInterval(args[0])
			if err != nil {
				if errors.Is(err, task.ErrIntervalTooSmall) {
					msg.Text = "间隔过小"
				} else {
					log.Errorf("set interval err: %s", err)
					msg.Text = "更新间隔失败"
				}
			} else {
				msg.Text = "更新间隔成功"
			}
			msg.ReplyToMessageID = update.Message.MessageID
			return msg
		},
	}
}
