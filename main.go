package main

import (
    "encoding/json"
    "encoding/xml"
    "fmt"
    "html/template"
    "io/ioutil"
    "log"
    "net/http"
    "os/exec"
    "regexp"
    "strconv"
    "strings"
)

type Config struct {
    Command        string `json:"command"`
    MatchType      string `json:"match_type"`      // "exact", "regex", "integer"
    MatchValue     string `json:"match_value"`
    Port           string `json:"port"`
    ConfigFilePath string
}

func (c *Config) Load() error {
    data, err := ioutil.ReadFile(c.ConfigFilePath)
    if err != nil {
        return err
    }
    return json.Unmarshal(data, c)
}

func (c *Config) Save() error {
    data, err := json.MarshalIndent(c, "", "  ")
    if err != nil {
        return err
    }
    return ioutil.WriteFile(c.ConfigFilePath, data, 0644)
}

type Result struct {
    Success    bool   `json:"success"`
    Output     string `json:"output"`
    Command    string `json:"command"`
    MatchType  string `json:"match_type"`
    MatchValue string `json:"match_value"`
}

func checkHandler(w http.ResponseWriter, r *http.Request, config *Config, responseType string) {
    cmd := exec.Command("sh", "-c", config.Command)
    output, err := cmd.Output()
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to run command: %v", err), http.StatusInternalServerError)
        return
    }

    success := checkOutput(string(output), config.MatchType, config.MatchValue)
    result := Result{
        Success:    success,
        Output:     string(output),
        Command:    config.Command,
        MatchType:  config.MatchType,
        MatchValue: config.MatchValue,
    }

    switch responseType {
    case "json":
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(result)
    case "xml":
        w.Header().Set("Content-Type", "application/xml")
        xml.NewEncoder(w).Encode(result)
    case "status":
        w.Header().Set("Content-Type", "text/plain")
        if success {
            fmt.Fprintln(w, "success")
        } else {
            fmt.Fprintln(w, "failed")
        }
    }
}

func checkOutput(output, matchType, matchValue string) bool {
    switch matchType {
    case "exact":
        return strings.TrimSpace(output) == matchValue
    case "regex":
        matched, _ := regexp.MatchString(matchValue, output)
        return matched
    case "integer":
        outputInt, err := strconv.Atoi(strings.TrimSpace(output))
        if err != nil {
            return false
        }
        matchInt, err := strconv.Atoi(matchValue)
        if err != nil {
            return false
        }
        return outputInt == matchInt
    default:
        return false
    }
}

func configureHandler(w http.ResponseWriter, r *http.Request, config *Config) {
    if r.Method == http.MethodPost {
        r.ParseForm()
        config.Command = r.FormValue("command")
        config.MatchType = r.FormValue("match_type")
        config.MatchValue = r.FormValue("match_value")
        config.Port = r.FormValue("port")
        if err := config.Save(); err != nil {
            http.Error(w, fmt.Sprintf("Failed to save config: %v", err), http.StatusInternalServerError)
            return
        }
    }

    tmpl := template.Must(template.New("config").Parse(`
<!DOCTYPE html>
<html>
<head>
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/css/bootstrap.min.css">
    <title>Configure</title>
</head>
<body>
<div class="container">
    <h1>Configuration</h1>
    <form method="POST">
        <div class="form-group">
            <label for="command">Command</label>
            <input type="text" class="form-control" id="command" name="command" value="{{.Command}}">
        </div>
        <div class="form-group">
            <label for="match_type">Match Type</label>
            <select class="form-control" id="match_type" name="match_type">
                <option value="exact" {{if eq .MatchType "exact"}}selected{{end}}>Exact</option>
                <option value="regex" {{if eq .MatchType "regex"}}selected{{end}}>Regex</option>
                <option value="integer" {{if eq .MatchType "integer"}}selected{{end}}>Integer</option>
            </select>
        </div>
        <div class="form-group">
            <label for="match_value">Match Value</label>
            <input type="text" class="form-control" id="match_value" name="match_value" value="{{.MatchValue}}">
        </div>
        <div class="form-group">
            <label for="port">HTTP Listen Port</label>
            <input type="text" class="form-control" id="port" name="port" value="{{.Port}}">
        </div>
        <button type="submit" class="btn btn-primary">Save</button>
    </form>
</div>
</body>
</html>
`))

    w.Header().Set("Content-Type", "text/html")
    tmpl.Execute(w, config)
}

func main() {
    config := &Config{ConfigFilePath: "config.json"}
    if err := config.Load(); err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    http.HandleFunc("/check/json", func(w http.ResponseWriter, r *http.Request) {
        checkHandler(w, r, config, "json")
    })
    http.HandleFunc("/check/xml", func(w http.ResponseWriter, r *http.Request) {
        checkHandler(w, r, config, "xml")
    })
    http.HandleFunc("/check/status", func(w http.ResponseWriter, r *http.Request) {
        checkHandler(w, r, config, "status")
    })
    http.HandleFunc("/configure", func(w http.ResponseWriter, r *http.Request) {
        configureHandler(w, r, config)
    })

    log.Printf("Server started on :%s", config.Port)
    log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}
