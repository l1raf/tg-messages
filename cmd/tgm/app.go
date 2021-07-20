package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/gorilla/mux"
	"github.com/gotd/contrib/middleware/floodwait"
	"github.com/gotd/contrib/middleware/ratelimit"
	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/telegram/updates"
	"github.com/gotd/td/tg"
	"github.com/jinzhu/gorm"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"
	"golang.org/x/xerrors"
	"tg-messages/internal/config"
	"tg-messages/internal/database"
	"tg-messages/internal/handlers"
	"tg-messages/internal/models"
	"tg-messages/internal/msg"
	"tg-messages/internal/utils"
)

type App struct {
	db         *gorm.DB
	cfg        config.Config
	msgRepo    *database.MessageRepository
	router     *mux.Router
	msgHandler *handlers.MessagesHandler
	logger     *zap.Logger
}

func NewApp() (*App, error) {
	cfg, err := config.Parse()
	if err != nil {
		return nil, err
	}

	a := App{cfg: *cfg}
	messagesDB, err := database.OpenConnection(cfg.ConnectionString)
	if err != nil {
		return nil, err
	}

	a.db = messagesDB
	err = database.Init()

	if err != nil {
		return nil, err
	}

	a.msgRepo = database.NewMessageStore(a.db)
	a.router = mux.NewRouter()
	a.msgHandler = handlers.NewMessagesHandler(a.msgRepo, cfg.Chats)

	a.logger, _ = zap.NewDevelopment(
		zap.IncreaseLevel(zapcore.DebugLevel),
	)

	return &a, nil
}

func (a *App) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)
	home, err := os.UserHomeDir()
	if err != nil {
		return xerrors.Errorf("get home: %w", err)
	}

	sessionDir := filepath.Join(home, ".td")

	if err := os.MkdirAll(sessionDir, 0700); err != nil {
		return xerrors.Errorf("mkdir: %w", err)
	}

	dispatcher := tg.NewUpdateDispatcher()
	client := telegram.NewClient(a.cfg.AppID, a.cfg.AppHash, telegram.Options{
		SessionStorage: &session.FileStorage{
			Path: filepath.Join(sessionDir, "session.json"),
		},
		UpdateHandler: dispatcher,
		Middlewares: []telegram.Middleware{
			floodwait.NewSimpleWaiter(),
			ratelimit.New(rate.Every(100*time.Millisecond), 5),
		},
		// Logger: a.logger,
		ReconnectionBackoff: func() backoff.BackOff {
			return backoff.WithMaxRetries(backoff.NewConstantBackOff(time.Nanosecond), 2)
		},
	})

	api := tg.NewClient(client)

	server := http.Server{
		Addr:    ":" + strconv.Itoa(a.cfg.Port),
		Handler: a.router,
	}

	group.Go(server.ListenAndServe)

	group.Go(func() error {
		<-ctx.Done()
		shutCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)

		defer cancel()

		if err := server.Shutdown(shutCtx); err != nil {
			return multierr.Append(err, server.Close())
		}

		return nil
	})

	group.Go(func() error {
		return client.Run(ctx, func(ctx context.Context) error {
			var flow auth.Flow

			if a.cfg.Password != "" {
				flow = auth.NewFlow(
					auth.Constant(a.cfg.Phone, a.cfg.Password, auth.CodeAuthenticatorFunc(getCode)),
					auth.SendCodeOptions{})
			} else {
				flow = auth.NewFlow(auth.CodeOnly(a.cfg.Phone, auth.CodeAuthenticatorFunc(getCode)), auth.SendCodeOptions{})
			}

			if err := client.Auth().IfNecessary(ctx, flow); err != nil {
				log.Println(xerrors.Errorf("auth: %w", err))
				return xerrors.Errorf("auth: %w", err)
			}

			tgMessages, err := msg.GetLastNMessages(ctx, api, a.cfg.Chats, a.cfg.N)
			if err != nil {
				log.Println(xerrors.Errorf("get messages: %w", err))
				tgMessages, err = retry(ctx, api, a.cfg.Chats, a.cfg.N)
			}

			if err == nil {
				a.saveMessages(tgMessages)
			}

			dispatcher.OnNewMessage(a.msgHandler.NewMessageHandler)
			dispatcher.OnNewChannelMessage(a.msgHandler.NewChannelMessageHandler)
			dispatcher.OnEditMessage(a.msgHandler.EditMessageHandler)

			var handleGraceful telegram.UpdateHandlerFunc = func(ctx context.Context, u tg.UpdatesClass) error {
				if err := dispatcher.Handle(ctx, u); err != nil {
					log.Println("handle update:", err)
				}

				return nil
			}

			gaps := updates.New(updates.Config{
				RawClient: api,
				Handler:   handlers.NewGapAdapter(ctx, handleGraceful),
				// Logger:    a.logger,
			})

			return gaps.Run(ctx)
		})
	})

	return group.Wait()
}

func retry(ctx context.Context, api *tg.Client, chats []int, n int) ([]tg.Message, error) {
	var err error
	var tgMessages []tg.Message

	start := time.Now()
	elapsed := time.Since(start)

	for elapsed.Seconds() < 15 {
		tgMessages, err = msg.GetLastNMessages(ctx, api, chats, n)

		if err == nil {
			return tgMessages, nil
		}

		elapsed = time.Since(start)
	}

	return tgMessages, err
}

func (a *App) saveMessages(messages []tg.Message) {
	for i := range messages {
		err := a.msgRepo.Create(models.Message{
			MessageID:  messages[i].ID,
			FromID:     utils.GetPeerID(messages[i].FromID),
			PeerID:     utils.GetPeerID(messages[i].PeerID),
			Text:       messages[i].Message,
			Date:       messages[i].Date,
			EditDate:   messages[i].EditDate,
			PostAuthor: messages[i].PostAuthor,
		})
		if err != nil {
			log.Println("save message:", err)
		}
	}
}

func (a *App) Close() error {
	err := a.db.Close()
	return multierr.Append(err, a.logger.Sync())
}

func getCode(ctx context.Context, sentCode *tg.AuthSentCode) (string, error) {
	fmt.Print("Enter code: ")

	code, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return code, err
	}

	code = strings.ReplaceAll(code, "\n", "")

	return code, nil
}
