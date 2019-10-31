/*
Copyright © 2019 Artificer Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"artificer/pkg/api/handlers"
	"artificer/pkg/client/loaders"
	"artificer/pkg/config"
	"artificer/pkg/keyvault"
	"artificer/pkg/util"
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	echo "github.com/labstack/echo/v4"
	middleware "github.com/labstack/echo/v4/middleware"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
)

var (
	ProcessDirectory string
	serverConfig     *ServerConfig
)

func init() {
	var err error
	ProcessDirectory, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ProcessDirectory)

	viper.SetConfigFile("config/appsettings.json")
	err = viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	useKeyVault := viper.GetBool("clientConfig.useKeyVault") // true|false

	loaders.InitializeClientConfig(loaders.ClientConfigOptions{
		RootFolder:  ProcessDirectory,
		UseKeyVault: useKeyVault,
	})
}

func Alive() {
	now := time.Now().UTC()
	fmt.Println(fmt.Sprintf("Alive:%s", now))
}

func executeEcho() {

	var err error

	keyVaultDone := make(chan bool, 1)
	clientConfigDone := make(chan bool, 1)
	go func() {
		fmt.Println("Startup Enter ... DoKeyvaultBackground")
		keyvault.DoKeyvaultBackground()
		fmt.Println("Startup Complete ... DoKeyvaultBackground")
		keyVaultDone <- true
	}()
	go func() {
		var (
			ctx    context.Context
			cancel context.CancelFunc
		)
		ctx, cancel = context.WithCancel(context.Background())
		defer cancel() // Cancel ctx as soon as handleSearch returns.

		fmt.Println("Startup Enter ... LoadClientConfig")
		loaders.LoadClientConfig(ctx)
		fmt.Println("Startup Complete ... LoadClientConfig")
		clientConfigDone <- true
	}()

	c := cron.New()
	cronSpec := viper.GetString("keyVault.cronSpec") // i.e. "@every 10s"
	_, err = c.AddFunc(cronSpec, func() {
		fmt.Println("CRON Enter ... DoKeyvaultBackground")
		keyvault.DoKeyvaultBackground()
		fmt.Println("CRON Complete ... DoKeyvaultBackground")
	})
	if err != nil {
		panic(err.Error())
	}
	cronSpec = viper.GetString("clientConfig.cronSpec") // i.e. "@every 5min"

	_, err = c.AddFunc(cronSpec, func() {
		var (
			ctx    context.Context
			cancel context.CancelFunc
		)
		ctx, cancel = context.WithCancel(context.Background())
		defer cancel() // Cancel ctx as soon as handleSearch returns.

		fmt.Println("CRON Enter ... LoadClientConfig")
		loaders.LoadClientConfig(ctx)
		fmt.Println("CRON Complete ... LoadClientConfig")
	})
	if err != nil {
		panic(err.Error())
	}
	c.Start()

	// Creating a new Echo instance.
	e := echo.New()
	// Configure Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost},
	}))

	//e.File("/", "static/index.html")

	// in order to serve static assets
	//e.Static("/static", "static")

	// Route / to handler function
	e.GET("/", handlers.Index)
	e.GET("/health", handlers.HealthCheck)
	e.GET("/.well-known/openid-configuration", handlers.WellKnownOpenidConfiguration)
	e.GET("/.well-known/openid-configuration/jwks", handlers.WellKnownOpenidConfigurationJwks)
	e.GET("/mint-test-token", handlers.MintTestToken)
	e.POST("/connect/token", handlers.TokenEndpoint)
	e.GET("/get-test-secret", handlers.GetTestSecret)
	// V1 Routes
	// v1 := e.Group("/v1")
	// v1Tokens := v1.Group("/tokens")
	// v1Tokens.GET("/tokens", handlers)
	fmt.Println()
	fmt.Println("------ Waiting for initial go routines to complete ------")
	fmt.Println()

	allDoneChannel := util.WaitOnAllChannels(keyVaultDone, clientConfigDone)
	<-allDoneChannel

	// If you start me up, I'll never stop
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", serverConfig.Port)))
}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serves artificer oauth2 server",
	Long: `serves artificer oauth2 server.
	Environment Variables Win:
		AF-key-vault-client-id
		AF-key-vault-client-secret
		AF-az-group-name
		AF-az-subscription-id
		AF-az-tenant-id
		AF-port`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("serve called")
		var err error
		serverConfig, err = validateVehicleRequest(cmd)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		config.SetGroupName(serverConfig.AzureGroupName)
		config.SetTenantID(serverConfig.AzureTenantId)
		config.SetSubscriptionID(serverConfig.AzureSubscriptionId)
		config.SetClientID(serverConfig.KeyVaultClientId)
		config.SetClientSecret(serverConfig.KeyVaultClientSecret)

		executeEcho()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	var (
		cliFlags *cliFlags
	)
	sc := ServerConfig{}

	cliFlags = sc.CliFlags_KeyVaultClientId()
	addFlag(serveCmd, cliFlags)

	cliFlags = sc.CliFlags_KeyVaultClientSecret()
	addFlag(serveCmd, cliFlags)

	cliFlags = sc.CliFlags_AzureGroupName()
	addFlag(serveCmd, cliFlags)

	cliFlags = sc.CliFlags_AzureSubscriptionId()
	addFlag(serveCmd, cliFlags)

	cliFlags = sc.CliFlags_AzureTenantId()
	addFlag(serveCmd, cliFlags)

	cliFlags = sc.CliFlags_Port()
	addFlag(serveCmd, cliFlags)
}

func validateVehicleRequest(cmd *cobra.Command) (sc *ServerConfig, err error) {
	var (
		val string
	)
	sc = &ServerConfig{}

	val = viper.GetString("AF-key-vault-client-id")
	if len(val) == 0 {
		val, err = cmd.Flags().GetString("key-vault-client-id")
		if err != nil {
			return
		}
		sc.KeyVaultClientId = val
	}

	val = viper.GetString("AF-key-vault-client-secret")
	if len(val) == 0 {
		val, err = cmd.Flags().GetString("key-vault-client-secret")
		if err != nil {
			return
		}
		sc.KeyVaultClientSecret = val
	}

	val = viper.GetString("AF-az-group-name")
	if len(val) == 0 {
		val, err = cmd.Flags().GetString("az-group-name")
		if err != nil {
			return
		}
		sc.AzureGroupName = val
	}

	val = viper.GetString("AF-az-subscription-id")
	if len(val) == 0 {
		val, err = cmd.Flags().GetString("az-subscription-id")
		if err != nil {
			return
		}
		sc.AzureSubscriptionId = val
	}

	val = viper.GetString("AF-az-tenant-id")
	if len(val) == 0 {
		val, err = cmd.Flags().GetString("az-tenant-id")
		if err != nil {
			return
		}
		sc.AzureTenantId = val
	}

	val = viper.GetString("AF-port")
	if len(val) == 0 {
		val, err = cmd.Flags().GetString("port")
		if err != nil {
			return
		}
		sc.Port = val
	}

	return
}

var (
	_ServerConfigType                reflect.Type
	_KeyVaultClientIdStructField     reflect.StructField
	_KeyVaultClientSecretStructField reflect.StructField
	_AzureGroupNameStructField       reflect.StructField
	_AzureSubscriptionIdStructField  reflect.StructField
	_AzureTenantIdStructField        reflect.StructField
	_PortStructField                 reflect.StructField
)

type ServerConfig struct {
	KeyVaultClientId     string `cli-required:"true" cli-long:"key-vault-client-id" cli-short:"" cli-default:"" cli-description:"Azure KeyVault Client Id" validate:"gt=1  & format=alnum_unicode"`
	KeyVaultClientSecret string `cli-required:"true" cli-long:"key-vault-client-secret" cli-short:"" cli-default:"" cli-description:"Azure KeyVault Client Secret" validate:"gt=1  & format=alnum_unicode"`
	AzureGroupName       string `cli-required:"true" cli-long:"az-group-name" cli-short:"" cli-default:"" cli-description:"Azure Group Name" validate:"gt=1  & format=alnum_unicode"`
	AzureSubscriptionId  string `cli-required:"true" cli-long:"az-subscription-id" cli-short:"" cli-default:"" cli-description:"Azure Subscription Id" validate:"gt=1  & format=alnum_unicode"`
	AzureTenantId        string `cli-required:"true" cli-long:"az-tenant-id" cli-short:"" cli-default:"" cli-description:"Azure Tenant Id" validate:"gt=1  & format=alnum_unicode"`
	Port                 string `cli-required:"false" cli-long:"port" cli-short:"p" cli-default:"9000" cli-description:"Artifice Server Port" validate:"gt=1  & format=alnum_unicode"`
}

func buildServerConfigReflectionData() {
	if _ServerConfigType == nil {
		_ServerConfigType = reflect.TypeOf(ServerConfig{})
		_KeyVaultClientIdStructField, _ = _ServerConfigType.FieldByName("KeyVaultClientId")
		_KeyVaultClientSecretStructField, _ = _ServerConfigType.FieldByName("KeyVaultClientSecret")
		_AzureGroupNameStructField, _ = _ServerConfigType.FieldByName("AzureGroupName")
		_AzureSubscriptionIdStructField, _ = _ServerConfigType.FieldByName("AzureSubscriptionId")
		_AzureTenantIdStructField, _ = _ServerConfigType.FieldByName("AzureTenantId")
		_PortStructField, _ = _ServerConfigType.FieldByName("Port")
	}
}
func (m *ServerConfig) CliFlags_KeyVaultClientId() *cliFlags {
	return m.getWellknownCliFlags(&_KeyVaultClientIdStructField)
}
func (m *ServerConfig) CliFlags_KeyVaultClientSecret() *cliFlags {
	return m.getWellknownCliFlags(&_KeyVaultClientSecretStructField)
}
func (m *ServerConfig) CliFlags_AzureGroupName() *cliFlags {
	return m.getWellknownCliFlags(&_AzureGroupNameStructField)
}
func (m *ServerConfig) CliFlags_AzureSubscriptionId() *cliFlags {
	return m.getWellknownCliFlags(&_AzureSubscriptionIdStructField)
}
func (m *ServerConfig) CliFlags_AzureTenantId() *cliFlags {
	return m.getWellknownCliFlags(&_AzureTenantIdStructField)
}
func (m *ServerConfig) CliFlags_Port() *cliFlags {
	return m.getWellknownCliFlags(&_PortStructField)
}
func (m *ServerConfig) getWellknownCliFlags(sf *reflect.StructField) *cliFlags {
	buildServerConfigReflectionData()
	cliLong, _ := sf.Tag.Lookup("cli-long")
	cliShort, _ := sf.Tag.Lookup("cli-short")
	cliDescription, _ := sf.Tag.Lookup("cli-description")
	cliDefault, _ := sf.Tag.Lookup("cli-default")
	cliRequired, _ := sf.Tag.Lookup("cli-required")
	required, err := strconv.ParseBool(cliRequired)
	if err != nil {
		panic(err)
	}

	return &cliFlags{
		LongFlag:    cliLong,
		ShortFlag:   cliShort,
		Description: cliDescription,
		Default:     cliDefault,
		Required:    required,
	}
}
func addFlag(cmd *cobra.Command, cliFlags *cliFlags) {
	if len(cliFlags.ShortFlag) == 1 {
		cmd.Flags().StringP(cliFlags.LongFlag, cliFlags.ShortFlag, cliFlags.Default, cliFlags.Description)
	} else {
		cmd.Flags().String(cliFlags.LongFlag, cliFlags.Default, cliFlags.Description)

	}
	if cliFlags.Required {
		cobra.MarkFlagRequired(cmd.Flags(), cliFlags.LongFlag)
	}
}
