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
	"artificer/pkg/cronex"
	"artificer/pkg/health"
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

	"github.com/heptiolabs/healthcheck"
	"github.com/spf13/cobra"

	"sync"

	echo "github.com/labstack/echo/v4"
	middleware "github.com/labstack/echo/v4/middleware"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
)

var (
	ProcessDirectory  string
	serverConfig      *ServerConfig
	healthCheckRecord HealthCheckRecord
	checksMutex       sync.RWMutex
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

func serveHealthCheck() {
	health.CheckIn(health.HealthRecord{
		Name:            "keyvault-api",
		Healthy:         false,
		UnhealthyReason: "Initial stat is always false",
	})
	health.CheckIn(health.HealthRecord{
		Name:            "client-config",
		Healthy:         false,
		UnhealthyReason: "Initial stat is always false",
	})

	healthCheckHandler := healthcheck.NewHandler()
	// Our app is not happy if we've got more than 100 goroutines running.
	healthCheckHandler.AddLivenessCheck("goroutine-threshold", healthcheck.GoroutineCountCheck(100))

	healthCheckHandler.AddReadinessCheck("client-config", health.CreateHealthCheck("client-config"))
	healthCheckHandler.AddReadinessCheck("keyvault-api", health.CreateHealthCheck("keyvault-api"))
	go http.ListenAndServe("0.0.0.0:"+serverConfig.HealthCheckPort, healthCheckHandler)
}

func executeKeyVaultFetch() {
	fmt.Println("CRON Enter ... DoKeyvaultBackground")
	err := keyvault.DoKeyVaultBackground()
	if err != nil {
		health.CheckIn(health.HealthRecord{
			Name:            "keyvault-api",
			Healthy:         false,
			UnhealthyReason: err.Error(),
		})
	} else {
		health.CheckIn(health.HealthRecord{
			Name:            "keyvault-api",
			Healthy:         true,
			UnhealthyReason: "",
		})
	}
	fmt.Println("CRON Complete ... DoKeyvaultBackground")
}
func executeLoadClientConfig() {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel() // Cancel ctx as soon as handleSearch returns.

	fmt.Println("CRON Enter ... LoadClientConfig")
	err := loaders.LoadClientConfig(ctx)
	if err != nil {
		health.CheckIn(health.HealthRecord{
			Name:            "client-config",
			Healthy:         false,
			UnhealthyReason: err.Error(),
		})
	} else {
		health.CheckIn(health.HealthRecord{
			Name:            "client-config",
			Healthy:         true,
			UnhealthyReason: "",
		})
	}
	fmt.Println("CRON Complete ... LoadClientConfig")
}
func serveArtificer() {

	var err error
	serveHealthCheck()
	keyVaultDone := make(chan bool, 1)
	clientConfigDone := make(chan bool, 1)

	c := cron.New()
	cronSpec := viper.GetString("keyVault.cronSpec") // i.e. "@every 10s"
	_, err = cronex.AddFunc(c, true, keyVaultDone, cronSpec, executeKeyVaultFetch)
	if err != nil {
		panic(err.Error())
	}
	cronSpec = viper.GetString("clientConfig.cronSpec") // i.e. "@every 5min"

	_, err = cronex.AddFunc(c, true, clientConfigDone, cronSpec, executeLoadClientConfig)
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

type DNSResolverRecord struct {
	Name string
	DNS  string
}
type HealthCheckRecord struct {
	DnsResolver []DNSResolverRecord
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
		fmt.Println(splash)
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

		err = viper.UnmarshalKey("healthCheck", &healthCheckRecord)
		for _, e := range healthCheckRecord.DnsResolver {
			fmt.Printf("%s:%s", e.Name, e.DNS)
		}
		serveArtificer()
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

	cliFlags = sc.CliFlags_HealthCheckPort()
	addFlag(serveCmd, cliFlags)

}

func validateVehicleRequest(cmd *cobra.Command) (sc *ServerConfig, err error) {
	var (
		val string
	)
	sc = &ServerConfig{}

	val, err = cmd.Flags().GetString("key-vault-client-id")
	if err != nil || len(val) == 0 {
		val = viper.GetString("AF-key-vault-client-id")
	}
	sc.KeyVaultClientId = val

	val, err = cmd.Flags().GetString("key-vault-client-secret")
	if err != nil || len(val) == 0 {
		val = viper.GetString("AF-key-vault-client-secret")
	}
	sc.KeyVaultClientSecret = val

	val, err = cmd.Flags().GetString("az-group-name")
	if err != nil || len(val) == 0 {
		val = viper.GetString("AF-az-group-name")
	}
	sc.AzureGroupName = val

	val, err = cmd.Flags().GetString("az-subscription-id")
	if err != nil || len(val) == 0 {
		val = viper.GetString("AF-az-subscription-id")
	}
	sc.AzureSubscriptionId = val

	val, err = cmd.Flags().GetString("az-tenant-id")
	if err != nil || len(val) == 0 {
		val = viper.GetString("AF-az-tenant-id")
	}
	sc.AzureTenantId = val

	val, err = cmd.Flags().GetString("port")
	if err != nil || len(val) == 0 {
		panic("port is not optional")
	}
	sc.Port = val

	val, err = cmd.Flags().GetString("healthcheck-port")
	if err != nil || len(val) == 0 {
		panic("healthcheck-port is not optional")
	}
	sc.HealthCheckPort = val

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
	_HealthCheckPortStructField      reflect.StructField
)

type ServerConfig struct {
	KeyVaultClientId     string `cli-required:"true" cli-long:"key-vault-client-id" cli-short:"" cli-default:"" cli-description:"Azure KeyVault Client Id" validate:"gt=1  & format=alnum_unicode"`
	KeyVaultClientSecret string `cli-required:"true" cli-long:"key-vault-client-secret" cli-short:"" cli-default:"" cli-description:"Azure KeyVault Client Secret" validate:"gt=1  & format=alnum_unicode"`
	AzureGroupName       string `cli-required:"true" cli-long:"az-group-name" cli-short:"" cli-default:"" cli-description:"Azure Group Name" validate:"gt=1  & format=alnum_unicode"`
	AzureSubscriptionId  string `cli-required:"true" cli-long:"az-subscription-id" cli-short:"" cli-default:"" cli-description:"Azure Subscription Id" validate:"gt=1  & format=alnum_unicode"`
	AzureTenantId        string `cli-required:"true" cli-long:"az-tenant-id" cli-short:"" cli-default:"" cli-description:"Azure Tenant Id" validate:"gt=1  & format=alnum_unicode"`
	Port                 string `cli-required:"false" cli-long:"port" cli-short:"p" cli-default:"" cli-description:"Artifice Server Port" validate:"gt=1  & format=alnum_unicode"`
	HealthCheckPort      string `cli-required:"false" cli-long:"healthcheck-port" cli-short:"" cli-default:"" cli-description:"Artifice Server Port" validate:"gt=1  & format=alnum_unicode"`
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
		_HealthCheckPortStructField, _ = _ServerConfigType.FieldByName("HealthCheckPort")
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
func (m *ServerConfig) CliFlags_HealthCheckPort() *cliFlags {
	return m.getWellknownCliFlags(&_HealthCheckPortStructField)
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
