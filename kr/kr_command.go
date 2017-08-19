package main

import (
	"encoding/base64"
	"os"
	"strconv"

	"github.com/kryptco/kr"
	"github.com/kryptco/kr/krdclient"
	"github.com/urfave/cli"
)

func createTeamCommand(c *cli.Context) (err error) {
	request, err := kr.NewRequest()
	if err != nil {
		PrintFatal(os.Stderr, err.Error())
		return
	}
	name := c.String("name")
	if name == "" {
		PrintFatal(os.Stderr, "Team name requied.")
	}
	request.CreateTeamRequest = &kr.CreateTeamRequest{
		Name: name,
	}
	response, err := krdclient.Request(request)
	if err != nil {
		PrintFatal(os.Stderr, err.Error())
	}
	createTeamResponse := response.CreateTeamResponse
	if createTeamResponse == nil {
		PrintFatal(os.Stderr, "Invalid response from Kryptonite app.")
	}
	if createTeamResponse.Error != nil {
		PrintFatal(os.Stderr, "Error creating team: "+*createTeamResponse.Error)
	}
	privateKeySeed := createTeamResponse.PrivateKeySeed
	if privateKeySeed == nil {
		PrintFatal(os.Stderr, "No team admin private key returned from Kryptonite app.")
	}
	kr.SaveAdminKeypair(*privateKeySeed)
	return
}

func createInviteCommand(c *cli.Context) (err error) {
	kr.CreateInvite()
	return
}

func setPolicyCommand(c *cli.Context) (err error) {
	var window *int64
	if c.String("window") != "" {
		windowInt, err := strconv.ParseInt(c.String("window"), 10, 64)
		if err != nil {
			PrintFatal(os.Stderr, "Could not parse window as integer seconds: "+err.Error())
		}
		window = &windowInt
	}
	kr.SetApprovalWindow(window)
	return
}

func getMembersCommand(c *cli.Context) (err error) {
	var query *string
	if c.String("query") != "" {
		queryStr := c.String("query")
		query = &queryStr
	}
	kr.GetMembers(query, c.Bool("ssh"), c.Bool("pgp"))
	return
}

func pinHostKeyCommand(c *cli.Context) (err error) {
	if c.String("public-key") == "" {
		kr.PinKnownHostKeys(c.String("host"), c.Bool("update-from-server"))
		return
	}
	pk, err := base64.StdEncoding.DecodeString(c.String("public-key"))
	if err != nil {
		PrintFatal(os.Stderr, "error decoding public-key, make sure it is base64 encoded without the key type prefix (i.e. no 'ssh-rsa' or 'ssh-ed25519') "+err.Error())
	}
	kr.PinHostKey(c.String("host"), pk)
	return
}

func unpinHostKeyCommand(c *cli.Context) (err error) {
	pk, err := base64.StdEncoding.DecodeString(c.String("public-key"))
	if err != nil {
		PrintFatal(os.Stderr, "error decoding public-key, make sure it is base64 encoded without the key type prefix (i.e. no 'ssh-rsa' or 'ssh-ed25519') "+err.Error())
	}
	kr.UnpinHostKey(c.String("host"), pk)
	return
}

func listPinnedKeysCommand(c *cli.Context) (err error) {
	host := c.String("host")
	if host == "" {
		kr.GetAllPinnedHostKeys()
		return
	}
	kr.GetPinnedHostKeys(host, c.Bool("search"))
	return
}
