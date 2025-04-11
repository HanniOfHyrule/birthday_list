package main

import (
    "encoding/json"
    "html/template"
    "net/http"
    "os"
)

type Birthday struct {
    Name     string
    Birthday string
}

var birthdays = func() []Birthday {
    file, err := os.Open("birthdays.json")
    if err != nil {
        panic(err)
    }
    defer file.Close()

    var loadedBirthdays []Birthday
    if err := json.NewDecoder(file).Decode(&loadedBirthdays); err != nil {
        panic(err)
    }
    return loadedBirthdays
}()

func main() {
    http.HandleFunc("/", listBirthdays)
    http.HandleFunc("/add-birthday-form", addBirthdayForm)
    http.HandleFunc("/add-birthday", addBirthday)

    http.ListenAndServe(":8088", nil)
}

func listBirthdays(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.New("birthdays").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Geburtstagsliste</title>
    <script src="https://unpkg.com/htmx.org"></script>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        table { width: 100%; border-collapse: collapse; margin-bottom: 20px; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        th { background-color: #f4f4f4; }
        button { padding: 10px 15px; background-color: #007BFF; color: white; border: none; cursor: pointer; }
        button:hover { background-color: #0056b3; }
    </style>
</head>
<body>
    <h1>Geburtstagsliste</h1>
    <table>
        <thead>
            <tr>
                <th>Name</th>
                <th>Geburtstag</th>
            </tr>
        </thead>
        <tbody id="birthday-list">
            {{range .Birthdays}}
            <tr>
                <td>{{.Name}}</td>
                <td>{{.Birthday}}</td>
            </tr>
            {{end}}
        </tbody>
    </table>
    <button hx-get="/add-birthday-form" hx-target="#form-container">Neuen Geburtstag hinzufügen</button>
    <div id="form-container"></div>
</body>
</html>
`))

    tmpl.Execute(w, map[string]interface{}{
        "Birthdays": birthdays,
    })
}

func addBirthdayForm(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.New("addBirthdayForm").Parse(`
<form hx-post="/add-birthday" hx-target="#birthday-list" hx-swap="beforeend">
    <label for="name">Name:</label>
    <input type="text" id="name" name="name" required>
    <label for="birthday">Geburtstag:</label>
    <input type="date" id="birthday" name="birthday" required>
    <button type="submit">Hinzufügen</button>
</form>
`))

    tmpl.Execute(w, nil)
}

func addBirthday(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        name := r.FormValue("name")
        birthday := r.FormValue("birthday")
        newBirthday := Birthday{Name: name, Birthday: birthday}
        birthdays = append(birthdays, newBirthday)

        // Aktualisiere die JSON-Datei
        file, err := os.Create("birthdays.json")
        if err != nil {
            panic(err)
        }
        defer file.Close()

        if err := json.NewEncoder(file).Encode(birthdays); err != nil {
            panic(err)
        }

        // Render the new row for HTMX
        tmpl := template.Must(template.New("row").Parse(`
        <tr>
            <td>{{.Name}}</td>
            <td>{{.Birthday}}</td>
        </tr>`))
        tmpl.Execute(w, newBirthday)
    }
}