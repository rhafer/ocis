// Package channels provides different communication channels to notify users.
package channels

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	groups "github.com/cs3org/go-cs3apis/cs3/identity/group/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/notifications/pkg/config"
	"github.com/pkg/errors"
	mail "github.com/xhit/go-simple-mail/v2"
)

// Channel defines the methods of a communication channel.
type Channel interface {
	// SendMessage sends a message to users.
	SendMessage(ctx context.Context, userIDs []string, msg, subject, senderDisplayName string) error
	// SendMessageToGroup sends a message to a group.
	SendMessageToGroup(ctx context.Context, groupdID *groups.GroupId, msg, subject, senderDisplayName string) error
}

// NewMailChannel instantiates a new mail communication channel.
func NewMailChannel(cfg config.Config, logger log.Logger) (Channel, error) {
	tm, err := pool.StringToTLSMode(cfg.Notifications.GRPCClientTLS.Mode)
	if err != nil {
		logger.Error().Err(err).Msg("could not get gateway client tls mode")
		return nil, err
	}
	gc, err := pool.GetGatewayServiceClient(cfg.Notifications.RevaGateway,
		pool.WithTLSCACert(cfg.Notifications.GRPCClientTLS.CACert),
		pool.WithTLSMode(tm),
	)
	if err != nil {
		logger.Error().Err(err).Msg("could not get gateway client")
		return nil, err
	}

	return Mail{
		gatewayClient: gc,
		conf:          cfg,
		logger:        logger,
	}, nil
}

// Mail is the communication channel for email.
type Mail struct {
	gatewayClient gateway.GatewayAPIClient
	conf          config.Config
	logger        log.Logger
}

func (m Mail) getMailClient() (*mail.SMTPClient, error) {
	server := mail.NewSMTPClient()
	server.Host = m.conf.Notifications.SMTP.Host
	server.Port = m.conf.Notifications.SMTP.Port
	server.Username = m.conf.Notifications.SMTP.Username
	if server.Username == "" {
		// compatibility fallback
		server.Username = m.conf.Notifications.SMTP.Sender
	}
	server.Password = m.conf.Notifications.SMTP.Password
	if server.TLSConfig == nil {
		server.TLSConfig = &tls.Config{}
	}
	server.TLSConfig.InsecureSkipVerify = m.conf.Notifications.SMTP.Insecure

	switch strings.ToLower(m.conf.Notifications.SMTP.Authentication) {
	case "login":
		server.Authentication = mail.AuthLogin
	case "plain":
		server.Authentication = mail.AuthPlain
	case "crammd5":
		server.Authentication = mail.AuthCRAMMD5
	case "none":
		server.Authentication = mail.AuthNone
	default:
		return nil, errors.New("unknown mail authentication method")
	}

	switch strings.ToLower(m.conf.Notifications.SMTP.Encryption) {
	case "tls":
		server.Encryption = mail.EncryptionTLS
		server.TLSConfig.ServerName = m.conf.Notifications.SMTP.Host
	case "starttls":
		server.Encryption = mail.EncryptionSTARTTLS
		server.TLSConfig.ServerName = m.conf.Notifications.SMTP.Host
	case "ssl":
		server.Encryption = mail.EncryptionSSL
	case "ssltls":
		server.Encryption = mail.EncryptionSSLTLS
	case "none":
		server.Encryption = mail.EncryptionNone
	default:
		return nil, errors.New("unknown mail encryption method")
	}

	smtpClient, err := server.Connect()
	if err != nil {
		return nil, err
	}

	return smtpClient, nil
}

// SendMessage sends a message to all given users.
func (m Mail) SendMessage(ctx context.Context, userIDs []string, msg, subject, senderDisplayName string) error {
	if m.conf.Notifications.SMTP.Host == "" {
		return nil
	}

	to, err := m.getReceiverAddresses(ctx, userIDs)
	if err != nil {
		return err
	}

	smtpClient, err := m.getMailClient()
	if err != nil {
		return err
	}

	email := mail.NewMSG()
	if senderDisplayName != "" {
		email.SetFrom(fmt.Sprintf("%s via %s", senderDisplayName, m.conf.Notifications.SMTP.Sender)).AddTo(to...)
	} else {
		email.SetFrom(m.conf.Notifications.SMTP.Sender).AddTo(to...)
	}
	email.SetBody(mail.TextPlain, msg)
	email.SetSubject(subject)

	return email.Send(smtpClient)
}

// SendMessageToGroup sends a message to all members of the given group.
func (m Mail) SendMessageToGroup(ctx context.Context, groupID *groups.GroupId, msg, subject, senderDisplayName string) error {
	res, err := m.gatewayClient.GetGroup(ctx, &groups.GetGroupRequest{GroupId: groupID})
	if err != nil {
		return err
	}
	if res.Status.Code != rpc.Code_CODE_OK {
		return errors.New("could not get group")
	}

	members := make([]string, 0, len(res.Group.Members))
	for _, id := range res.Group.Members {
		members = append(members, id.OpaqueId)
	}

	return m.SendMessage(ctx, members, msg, subject, senderDisplayName)
}

func (m Mail) getReceiverAddresses(ctx context.Context, receivers []string) ([]string, error) {
	addresses := make([]string, 0, len(receivers))
	for _, id := range receivers {
		// Authenticate is too costly but at the moment our only option to get the user.
		// We don't have an authenticated context so calling `GetUser` doesn't work.
		res, err := m.gatewayClient.Authenticate(ctx, &gateway.AuthenticateRequest{
			Type:         "machine",
			ClientId:     "userid:" + id,
			ClientSecret: m.conf.Notifications.MachineAuthAPIKey,
		})
		if err != nil {
			return nil, err
		}
		if res.Status.Code != rpc.Code_CODE_OK {
			m.logger.Error().
				Interface("status", res.Status).
				Str("receiver_id", id).
				Msg("could not get user")
			continue
		}
		addresses = append(addresses, res.User.Mail)
	}

	return addresses, nil
}
