package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/angelfluffyookami/247BVR/modules/common"
	"github.com/bwmarrin/discordgo"
)

const (
	token           = "tokenhere"
	adminserver     = "adminserveridhere"
	adminchannel    = "adminchannelidhere"
	infochannel     = "infochannelidhere"
	warnchannel     = "warnchannelidhere"
	errchannel      = "errorchannelidhere"
	updatechannel   = "updatechannelidhere"
	feedbackchannel = "feedbackchannelidhere"
)

// Function reads json file, returning variable of type config.Data
func ReadJSON(filename string) (common.Data, error) {
	//	Reads file and saves []byte to variable data, then checks if there was an error, if error, then return nil config and error.
	data, err := os.ReadFile(filename)
	if err != nil {
		return config, err
	}

	//	Unmarshals []byte of type json into variable config of type config.Data, if error, then return nil config, and error.
	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	// 	If no error, return config, and nil error.
	return config, nil
}

// Registers commands, and adds command handling functions to discord session after discord session is opened.
func initDiscordHandlers() {
	//	Adds handler of type i *discordgo.InteractionCreate that selects the appropriate handler over a map of command handlers depending on the application
	s.AddHandler(func(i *discordgo.InteractionCreate) {
		if h, ok := allCommandHandlers[i.ApplicationCommandData().Name]; ok {
			fmt.Println(i.ApplicationCommandData().Name)
			h(i)
		}
	})

	log.Println("Adding commands...")
	// registers commands with the discord session.
	registeredCommands := make([]*discordgo.ApplicationCommand, len(allCommands))
	for i, v := range allCommands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, "", &v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
}

// Clean up session by removing commands and closing the session afterwards
func exit() {
	fmt.Println("Removing commands...")

	// Get a list of all the commands registered within.
	commands, err := s.ApplicationCommands(s.State.User.ID, "")
	if err != nil {
		panic(err)
	}

	// iterates through all commands and deletes them.
	for x := 0; x < len(commands); x++ {
		err := s.ApplicationCommandDelete(s.State.User.ID, "", commands[x].ID)
		if err != nil {
			panic(err)
		}
	}
	// closes discord session.
	s.Close()
}

// beautifies a json file
func beautifyJSONFile(filename string) {
	// Open the given file
	jsonFile, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
	}

	// Unmarshal the contents of the file into an interface
	var jsonData interface{}
	json.Unmarshal(jsonFile, &jsonData)

	// Marshal the data into a byte array with indentation
	prettyJSON, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	// Write the byte array back to the file
	err = os.WriteFile(filename, prettyJSON, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

// CreateOrUpdateJSON creates or updates a JSON file with two keys
func CreateOrUpdateJSON(file string) error {
	// Read the existing file
	data := common.Data{}
	bytes, err := os.ReadFile(file)
	if err != nil {
		// File does not exist, create it
		data = common.Data{
			Token:         token,
			AdminServer:   adminserver,
			AdminChannel:  adminchannel,
			InfoChannel:   infochannel,
			WarnChannel:   warnchannel,
			ErrChannel:    errchannel,
			UpdateChannel: updatechannel,
		}
	} else {
		// File exists, parse it
		if err := json.Unmarshal(bytes, &data); err != nil {
			return fmt.Errorf("failed to parse existing file: %v", err)
		}
		// Check if keys are missing.
		if data.Token == "" {
			data.Token = token
		}
		if data.AdminServer == "" {
			data.AdminServer = adminserver
		}
		if data.AdminChannel == "" {
			data.AdminChannel = adminchannel
		}
		if data.InfoChannel == "" {
			data.InfoChannel = infochannel
		}
		if data.WarnChannel == "" {
			data.WarnChannel = warnchannel
		}
		if data.ErrChannel == "" {
			data.ErrChannel = errchannel
		}
		if data.UpdateChannel == "" {
			data.UpdateChannel = updatechannel
		}
	}

	// Marshal the data back to JSON
	bytes, err = json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	// Write the data back to the file
	if err := os.WriteFile(file, bytes, 0644); err != nil {
		return fmt.Errorf("failed to write data to file: %v", err)
	}

	return nil
}

func checkAPIEndpoint() {
	log.Println("Probing API endpoints to ensure responses match expected.")
	URL := common.Config.APIEndpoint
	if !strings.HasSuffix(URL, "/") {
		common.Config.APIEndpoint = common.Config.APIEndpoint + "/"
	}

	var fatal bool
	err = checkOnlineEndpoint()
	if err != nil {
		fatal = true
		log.Println(err)
	}
	err = checkUsersEndpoint()
	if err != nil {
		fatal = true
		log.Println(err)
	}
	err = checkUserID()
	if err != nil {
		fatal = true
		log.Println(err)
	}
	err = checkKillsEndpoint()
	if err != nil {
		fatal = true
		log.Println(err)
	}
	err = checkDeathEndpoint()
	if err != nil {
		fatal = true
		log.Println(err)
	}
	err = checkLogEndpoint()
	if err != nil {
		fatal = true
		log.Println(err)
	}
	if fatal {
		log.Fatal("One or more JSON responses do not match expected.")
	}
}

func checkOnlineEndpoint() error {
	req, err := http.Get(common.Config.APIEndpoint + "online")
	if err != nil {
		return fmt.Errorf("error trying to submit HTTP GET request to /users")
	}
	log.Println("API Endpoint HTTP GET /online successful. Checking if JSON response matches expected.")

	var Online []common.Online
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return fmt.Errorf("error trying to read /online response body")
	}
	err = json.Unmarshal(body, &Online)
	if err != nil {
		return fmt.Errorf("the JSON response does not match expected. /online \n %s", body)
	}
	log.Println("JSON response matches expected. /online")
	return nil
}

func checkUsersEndpoint() error {
	req, err := http.Get(common.Config.APIEndpoint + "users")
	if err != nil {
		return fmt.Errorf("error trying to submit HTTP GET request to /users")
	}
	log.Println("API Endpoint HTTP GET /users successful. Checking if JSON response matches expected.")

	var Users []common.User
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return fmt.Errorf("error trying to read /users response body")
	}
	err = json.Unmarshal(body, &Users)
	if err != nil {
		return fmt.Errorf("the JSON response does not match expected. /users \n %s", body)
	}
	log.Println("JSON response matches expected. /users")
	return nil
}

func checkUserID() error {
	req, err := http.Get(common.Config.APIEndpoint + "users/" + common.DefaultID)
	if err != nil {
		return fmt.Errorf("error trying to submit HTTP GET request to /users/<id>")
	}
	log.Println("API Endpoint HTTP GET /users/<id> successful. Checking if JSON response matches expected.")

	var User common.User
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return fmt.Errorf("error trying to read /users/<id> response body")
	}
	err = json.Unmarshal(body, &User)
	if err != nil {
		return fmt.Errorf("the JSON response does not match expected. /users \n %s", body)
	}
	log.Println("JSON response matches expected. /users/<id>")
	return nil
}

func checkKillsEndpoint() error {
	req, err := http.Get(common.Config.APIEndpoint + "kills")
	if err != nil {
		return fmt.Errorf("error trying to submit HTTP GET request to /kills")
	}
	log.Println("API Endpoint HTTP GET /kills successful. Checking if JSON response matches expected.")

	var Kills []common.Kill
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return fmt.Errorf("error trying to read /kills response body")
	}
	err = json.Unmarshal(body, &Kills)
	if err != nil {
		return fmt.Errorf("the JSON response does not match expected. /kills \n %s", body)
	}
	log.Println("JSON response matches expected. /kills")
	return nil
}

func checkDeathEndpoint() error {
	req, err := http.Get(common.Config.APIEndpoint + "/deaths")
	if err != nil {
		return fmt.Errorf("error trying to submit HTTP GET request to /deaths")
	}
	log.Println("API Endpoint HTTP GET /deaths successful. Checking if JSON response matches expected.")

	var Kills []common.Kill
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return fmt.Errorf("error trying to read /deaths response body")
	}
	err = json.Unmarshal(body, &Kills)
	if err != nil {
		return fmt.Errorf("the JSON response does not match expected. /deaths \n %s", body)
	}
	log.Println("JSON response matches expected. /deaths")
	return nil
}

func checkLogEndpoint() error {
	req, err := http.Get(common.Config.APIEndpoint + "log/" + common.DefaultID)
	if err != nil {
		return fmt.Errorf("error trying to submit HTTP GET request to /log")
	}
	log.Println("API Endpoint HTTP GET /kills successful. Checking if JSON response matches expected.")

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return fmt.Errorf("error trying to read /kills response body")
	}

	log.Printf("TXT response received. /kills \n%s", body)
	return nil
}
