package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"text/template"

	"github.com/jllopis/try5/account"
	"github.com/mgutz/logxi/v1"
	"github.com/spf13/cobra"
)

const (
	VERSION = "v0.0.1"
	try5uri = "https://b2d:9000"
	//try5uri = "https://srv18.acb.info:9000"
)

var (
	logger   log.Logger
	pAccTmpl *template.Template
)

func init() {
	logger = log.New("try5cli")
	var err error
	pAccTmpl, err = template.New("pAccount").Parse(`{"email": "{{ .Email }}", "name": "{{ .Name }}", "password": "{{ .Password }}"}`)
	if err != nil {
		logger.Fatal("init", "error", err.Error())
	}
}

func main() {
	var (
		batchFile string
		newAcc    string
	)

	var cmdVersion = &cobra.Command{
		Use:   "version",
		Short: "Muestra la versión del programa",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("try5cli", VERSION)
		},
	}

	var cmdAccount = &cobra.Command{
		Use:   "account [account name]",
		Short: "Operaciones con las cuentas",
		Long:  `account permite operar con las cuentas, a través de los subcomandos, para crear, modificar, eliminar y obtener las cuentas.`,
	}

	var cmdListAccounts = &cobra.Command{
		Use:   "list",
		Short: "List accounts",
		Run: func(cmd *cobra.Command, args []string) {
			var msg []map[string]interface{}
			tr := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
			client := &http.Client{Transport: tr}
			resp, err := client.Get(try5uri + "/api/v1/accounts")
			if err != nil {
				logger.Warn("accounts list", "error", err.Error())
				return
			}
			defer resp.Body.Close()
			err = json.NewDecoder(resp.Body).Decode(&msg)
			if resp.StatusCode != 200 {
				logger.Warn("accounts list", "error", resp.Status)
				return
			}
			output, _ := json.Marshal(msg)
			fmt.Printf("%v\n", string(output))
		},
	}

	var cmdNewAccount = &cobra.Command{
		Use:   "new",
		Short: "create a new account",
		Run: func(cmd *cobra.Command, args []string) {
			var wg sync.WaitGroup
			done := make(chan struct{})
			comms := make(chan string)
			var accs []account.Account
			if newAcc == "" && batchFile == "" {
				logger.Warn("account new", "error", "se necesita la cuenta a crear en formato JSON")
				return
			}
			if newAcc != "" {
				err := json.Unmarshal([]byte(newAcc), &accs)
				if err != nil {
					logger.Warn("account new", "error", err.Error())
					return
				}
			}
			if batchFile != "" {
				tmp, err := ioutil.ReadFile(batchFile)
				if err != nil {
					logger.Warn("account new", "error", err.Error())
					if newAcc == "" {
						return
					}
				}
				var tmpAccs []account.Account
				err = json.Unmarshal(tmp, &tmpAccs)
				if err != nil {
					logger.Warn("account new", "error", err.Error())
					if newAcc == "" {
						return
					}
				}
				for _, v := range tmpAccs {
					accs = append(accs, v)
				}

			}
			for _, a := range accs {
				if a.Email == nil {
					logger.Warn("accounts new", "error", "email no puede ser nulo")
					continue
				}
				if a.Name == nil {
					a.Name = a.Email
				}
				x, err := json.Marshal(a)
				if err != nil {
					logger.Warn("account new", "error", err.Error())
				}
				wg.Add(1)
				go func(a account.Account) {
					defer wg.Done()
					comms <- fmt.Sprintf("\033[33m[INF]\033[0m Creando cuenta: %s (%s)", *a.Name, *a.Email)
					resp, err := httpClient().Post(try5uri+"/api/v1/accounts", "application/json; charset=utf-8", bytes.NewBuffer(x))
					if err != nil {
						//logger.Warn("accounts new", "error", err.Error())
						comms <- fmt.Sprintf("\033[31m[ERR]\033[0m accounts new: %v", err)
						return
					}
					if resp.StatusCode != 201 {
						//logger.Warn("accounts new", "error", resp.Status)
						comms <- fmt.Sprintf("\033[31m[ERR]\033[0m accounts new: Respuesta incorrecta del servidor: %s", resp.Status)
						return
					}
					defer resp.Body.Close()
					comms <- fmt.Sprintf("\033[32m[OK]\033[0m Creada cuenta: %s (%s)", *a.Name, *a.Email)
				}(a)
			}
			go func() {
				for {
					select {
					case msg := <-comms:
						fmt.Println(msg)
					case <-done:
						return
					}
					logger.Info("accounts new", "message", "quitting listener goroutine")
				}
			}()
			wg.Wait()
		},
	}

	cmdNewAccount.Flags().StringVarP(&batchFile, "file", "f", "", "files with a bunch of json formatted lines corresponding accounts")
	cmdNewAccount.Flags().StringVarP(&newAcc, "account", "a", "", "accepts a JSON list of Account objects that will be created")

	var rootCmd = &cobra.Command{Use: "try5cli"}
	rootCmd.AddCommand(cmdVersion, cmdAccount)
	cmdAccount.AddCommand(cmdListAccounts, cmdNewAccount)
	rootCmd.Execute()
}

func httpClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	return client
}
