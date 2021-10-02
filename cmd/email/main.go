package main

import (
	"crypto/tls"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/itsabgr/authentify"
	"github.com/itsabgr/authentify/pkg/smtpSender"
	"github.com/itsabgr/go-handy"
	"google.golang.org/grpc"
	"net/smtp"
	"time"
)

var (
	emailServer   = ""
	emailFrom     = ""
	emailIdentity = ""
	emailPass     = ""
	emailUser     = ""
	emailHost     = ""
	badgerPath    = ""
	serverCert    = []byte("")
	serverKey     = []byte("")
	serverAddr    = ":443"
)

type EmailTemplate struct{}

func (_ EmailTemplate) RenderMsg(_ string, code authentify.Code, exp time.Time) ([]byte, error) {
	return []byte(fmt.Sprintf("Code: %s\nExpire at: %s\n", code, exp)), nil
}

var emailTemplate = EmailTemplate{}

func main() {
	smtpAuth := smtp.PlainAuth(emailIdentity, emailUser, emailPass, emailHost)
	emailSender, err := smtpSender.NewEmailSender(emailServer, emailFrom, emailTemplate, smtpAuth)
	handy.Throw(err)
	db, err := badger.Open(badger.DefaultOptions(badgerPath))
	handy.Throw(err)
	defer db.Close()
	repo, err := authentify.BadgerAsRepo(db)
	handy.Throw(err)
	defer repo.Close()
	auth, err := authentify.NewAuthenticator(authentify.Options{
		Repo:        repo,
		PrefixChars: []rune("ABCDEGHKLNOPSTXYZ"),
		CodeChars:   []rune("0123456789"),
		CodeLength:  6,
		SaltLength:  32,
		TTL:         1 * time.Minute,
	})
	handy.Throw(err)
	sendersMap, err := authentify.SendersToMap(emailSender)
	handy.Throw(err)
	authServer, err := authentify.NewGrpcProtoServer(auth, sendersMap)
	handy.Throw(err)
	tlsCert, err := tls.X509KeyPair(serverCert, serverKey)
	handy.Throw(err)
	tlsServer, err := tls.Listen("tcp", serverAddr, &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
	})
	handy.Throw(err)
	defer tlsServer.Close()
	grpcServer := grpc.NewServer()
	defer grpcServer.Stop()
	grpcServer.RegisterService(authentify.GrpcProtoServiceDesc(), authServer)
	handy.Throw(grpcServer.Serve(tlsServer))
}
