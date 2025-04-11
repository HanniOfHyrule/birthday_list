package main

import (
    "encoding/json"
    "html/template"
    "net/http"
    "os"
    "time"    
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
    tmpl := template.New("birthdays").Funcs(template.FuncMap{
        "formatDate": func(date string) string {
            parsedDate, err := time.Parse("2006-01-02", date)
            if err != nil {
                return date // Fallback, falls das Datum nicht geparst werden kann
            }
            return parsedDate.Format("02.01.2006") // Format für Deutschland
        },
    })

    tmpl = template.Must(tmpl.Parse(`
<!DOCTYPE html>
<html lang="de">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Geburtstagsliste</title>
    <script src="https://unpkg.com/htmx.org"></script>
    <style>
        body {
            font-family: 'Arial', sans-serif;
            margin: 20px;
            background-color: #f4f4f9;
            color: #333;
        }
        h1 {
            text-align: center;
            color: #444;
            padding: 1rem 1rem 1rem 2rem 
        }
        table {
            width: 90%;
            margin: 3rem;
            border-collapse: collapse;
            background-color: #fff;
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
            border-radius: 10px;
            overflow: hidden;
        }
        th, td {
            padding: 12px 15px;
            text-align: left;
            border-bottom: 1px solid #ddd;
        }
        th {
            background-color: #6c63ff;
            color: #fff;
            text-transform: uppercase;
            font-size: 14px;
        }
        tr:hover {
            background-color: #f1f1f1;
        }
        button {
            display: block;
            margin: 20px auto;
            padding: 12px 25px;
            background-color: #6c63ff;
            color: #fff;
            border: none;
            border-radius: 25px;
            cursor: pointer;
            font-size: 16px;
            transition: background-color 0.3s ease;
        }
        button:hover {
            background-color: #4b47cc;
        }
        form {
            margin: 20px auto;
            padding: 20px;
            background-color: #fff;
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
            border-radius: 10px;
            max-width: 400px;
        }
        label {
            display: block;
            margin-bottom: 8px;
            font-weight: bold;
        }
        input {
            width: calc(100% - 20px); /* Platz für Padding */
            padding: 10px;
            margin-bottom: 15px;
            margin-right: 10px; /* Abstand nach rechts */
            border: 1px solid #ddd;
            border-radius: 5px;
        }
        input[type="date"] {
            font-family: 'Arial', sans-serif;
        }
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
                <td>{{.Birthday | formatDate}}</td>
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