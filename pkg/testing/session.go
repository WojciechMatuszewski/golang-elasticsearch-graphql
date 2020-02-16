package testing

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

// LocalSession returns new local aws session
func LocalSession() *session.Session {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("local"),
		Credentials: credentials.NewStaticCredentials("local", "local", "local"),
	})
	if err != nil {
		panic(err.Error())
	}

	return sess
}
