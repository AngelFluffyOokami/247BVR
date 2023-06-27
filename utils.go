package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/angelfluffyookami/247BVR/modules/common"
	"github.com/bwmarrin/discordgo"
	"github.com/goccy/go-json"
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

/*
* Probes API endpoint to ensure API responses match expected.
* Ran at init to ensure URL is not incorrect and/or changes were made to API.
* Concurrent for ease of debugging purposes.
* (I'm not gonna spend 3 minutes or so waiting for a sanity check every time I compile and debug this)
 */
func sanityCheck() {
	log.Println("Probing API endpoints to ensure responses match expected.")

	// Get URL from global variable.
	URL := common.Config.APIEndpoint

	// Checks if URL has slash at the end, if not, add it and save it to global variable.
	if !strings.HasSuffix(URL, "/") {

		common.Config.APIEndpoint = common.Config.APIEndpoint + "/"
	}

	// A massive fuck you to whichever end receives these HTTP GETs depending on JSON dataset size.
	go checkOnlineEndpoint(common.Config.Debugging)
	go checkUsersEndpoint(common.Config.Debugging)
	go checkUserID(common.Config.Debugging)
	go checkKillsEndpoint(common.Config.Debugging)
	go checkDeathEndpoint(common.Config.Debugging)
	go checkLogEndpoint(common.Config.Debugging)

}

// Sanity Check /online
func checkOnlineEndpoint(debug bool) {

	getBench := time.Now()

	req, err := http.Get(common.Config.APIEndpoint + "online")

	getExecTime := time.Since(getBench)

	if err != nil {
		log.Panicf("Sanity check: %s", err)
	}

	if debug {
		log.Println("API Endpoint HTTP GET /online successful." + `
	` + "\nHTTP GET Time: " + getExecTime.String() + `
	` + "Checking if JSON response matches expected.")
	}

	var Online []common.Online

	unmarshalBench := time.Now()

	err = json.NewDecoder(req.Body).Decode(&Online)

	unmarshalExecTime := time.Since(unmarshalBench)

	if err != nil {
		log.Panicf("Sanity check: %s", err)
	}

	if debug {
		log.Println("JSON response matches expected. /online" + `
	` + "JSON Unmarshal Exec Time: " + unmarshalExecTime.String())
	}

}

// Sanity Check /users
func checkUsersEndpoint(debug bool) {

	getBench := time.Now()

	req, err := http.Get(common.Config.APIEndpoint + "users")

	getExecTime := time.Since(getBench)

	if err != nil {
		log.Panicf("Sanity check: %s", err)
	}

	if debug {
		log.Println("API Endpoint HTTP GET /users successful." + `
	` + "\nHTTP GET Time: " + getExecTime.String() + `
	` + "Checking if JSON response matches expected.")
	}

	var Users []common.User

	unmarshalBench := time.Now()

	err = json.NewDecoder(req.Body).Decode(&Users)

	unmarshalExecTime := time.Since(unmarshalBench)

	if err != nil {
		log.Panicf("Sanity check: %s", err)
	}

	if debug {
		log.Println("JSON response matches expected. /users" + `
		` + "JSON Unmarshal Exec Time: " + unmarshalExecTime.String())
	}

}

// Sanity Check /user/<id>
func checkUserID(debug bool) {
	getBench := time.Now()

	req, err := http.Get(common.Config.APIEndpoint + "users/" + common.DefaultID)

	getExecTime := time.Since(getBench)

	if err != nil {
		log.Panicf("Sanity check: %s", err)
	}

	if debug {
		log.Println("API Endpoint HTTP GET users/<id> successful." + `
	` + "\nHTTP GET Time: " + getExecTime.String() + `
	` + "Checking if JSON response matches expected.")
	}

	var User common.User

	unmarshalBench := time.Now()

	err = json.NewDecoder(req.Body).Decode(&User)

	unmarshalExecTime := time.Since(unmarshalBench)

	if err != nil {
		log.Panicf("Sanity check: %s", err)
	}

	if debug {
		log.Println("JSON response matches expected. /users/<id>" + `
		` + "JSON Unmarshal Exec Time: " + unmarshalExecTime.String())
	}

}

// Sanity Check /kills
func checkKillsEndpoint(debug bool) {

	getBench := time.Now()

	req, err := http.Get(common.Config.APIEndpoint + "kills")

	getExecTime := time.Since(getBench)

	if err != nil {
		log.Panicf("Sanity check: %s", err)
	}

	if debug {
		log.Println("API Endpoint HTTP GET /kills successful." + `
	` + "\nHTTP GET Time: " + getExecTime.String() + `
	` + "Checking if JSON response matches expected.")
	}

	var Kills []common.Kill

	unmarshalBench := time.Now()

	err = json.NewDecoder(req.Body).Decode(&Kills)

	unmarshalExecTime := time.Since(unmarshalBench)

	if err != nil {
		log.Panicf("Sanity check: %s", err)
	}

	if debug {
		log.Println("JSON response matches expected. /kills" + `
	` + "JSON Unmarshal Exec Time: " + unmarshalExecTime.String())
	}
}

// Sanity Check /deaths
func checkDeathEndpoint(debug bool) {

	getBench := time.Now()

	req, err := http.Get(common.Config.APIEndpoint + "/deaths")

	getExecTime := time.Since(getBench)

	if err != nil {
		log.Panicf("Sanity check: %s", err)
	}

	if debug {
		log.Println("API Endpoint HTTP GET /deaths successful." + `
	` + "\nHTTP GET Time: " + getExecTime.String() + `
	` + "Checking if JSON response matches expected.")
	}

	var Kills []common.Kill
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Panicf("Sanity check: %s", err)
	}

	err = json.Unmarshal(body, &Kills)
	if err != nil {
		log.Panicf("Sanity check: %s", err)
	}

	if debug {
		log.Println("JSON response matches expected. /deaths")
	}

}

// Sanity Check /log/<id>
func checkLogEndpoint(debug bool) {
	getBench := time.Now()

	req, err := http.Get(common.Config.APIEndpoint + "log/" + common.DefaultID)

	getExecTime := time.Since(getBench)

	if err != nil {
		log.Panicf("Sanity check: %s", err)
	}

	if debug {
		log.Println("API Endpoint HTTP GET /log successful." + `
		` + "HTTP GET Time: " + getExecTime.String() + `
		` + "Checking if TXT response matches expected.")
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Panicf("Sanity check: %s", err)
	}

	if debug {
		log.Printf("TXT response received. /kills \n%s", body)
	}
}
