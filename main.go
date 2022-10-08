package main

import (
    "html/template"
    "net/http"
    "fmt"
    "log"
    "database/sql"
    //"encoding/csv"
    //"io"
    //"os"
    //"encoding/json"

	_ "github.com/go-sql-driver/mysql"
)

type GenderData struct {
    PageTitle   string    
    Genders     []Gender 
}
type Gender struct {
    Code        string 
    Name        string 
}

type PositiveData struct {
    PageTitle   string
    Positives   []Positive
}

type Positive struct {
    Count       string
    Lat         string
    Lng         string
}

func main() {

    fs := http.FileServer(http.Dir("templates/assets/"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    fmt.Printf("Connecting to db\n")
    db, err := sql.Open("mysql", "root:root@tcp(golang-mysql:3306)/hackujstat")
    defer db.Close()
    if err != nil {
        log.Fatal(err)
    }

    tmplCovid := template.Must(template.ParseFiles("templates/main.html"))
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    
        res, err := db.Query("SELECT * FROM pohlavi")
        defer res.Close()
    
        if err != nil {
            log.Fatal(err)
        }
    
        gd := GenderData{
            PageTitle: "Genders",
        }
    
        for res.Next() {
            var g Gender
            if err := res.Scan(&g.Code, &g.Name); err != nil {  }
            gd.Genders = append(gd.Genders, Gender{Code: g.Code, Name: g.Name})
        }        
    
       tmplCovid.Execute(w, gd)
    }) 
    
    tmplHospital := template.Must(template.ParseFiles("templates/hospitalizace.html"))
    http.HandleFunc("/hospitalizace", func(w http.ResponseWriter, r *http.Request) {
    
        res, err := db.Query("SELECT * FROM pohlavi")
        defer res.Close()
    
        if err != nil {
            log.Fatal(err)
        }
    
        gd := GenderData{
            PageTitle: "Hospitalizace",
        }
    
        for res.Next() {
            var g Gender
            if err := res.Scan(&g.Code, &g.Name); err != nil {  }
            gd.Genders = append(gd.Genders, Gender{Code: g.Code, Name: g.Name})
        }        
    
       tmplHospital.Execute(w, gd)
    }) 


    tmplJip := template.Must(template.ParseFiles("templates/jip.html"))
    http.HandleFunc("/jip", func(w http.ResponseWriter, r *http.Request) {
    
        res, err := db.Query("SELECT * FROM pohlavi")
        defer res.Close()
    
        if err != nil {
            log.Fatal(err)
        }
    
        gd := GenderData{
            PageTitle: "JIP",
        }
    
        for res.Next() {
            var g Gender
            if err := res.Scan(&g.Code, &g.Name); err != nil {  }
            gd.Genders = append(gd.Genders, Gender{Code: g.Code, Name: g.Name})
        }        
    
       tmplJip.Execute(w, gd)
    }) 


    http.HandleFunc("/getCovid", func(w http.ResponseWriter, r *http.Request) {
    
        res, err := db.Query("select count(okres_lau_kod) as amount, okresy.lat, okresy.lng from pozitivnipripady left join okresy on pozitivnipripady.okres_lau_kod = okresy.code where okresy.lat != '' group by okresy.lat,okresy.lng")
        defer res.Close()
    
        if err != nil {
            log.Fatal(err)
        }

        w.Header().Add("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)

        for res.Next() {
            var p Positive
            if err := res.Scan(&p.Count, &p.Lat, &p.Lng); err != nil {  }
		    w.Write([]byte(`{"lat": `+ p.Lat + `, "lng": `+ p.Lng + `, "count": `+ p.Count + `},`))

        }
    }) 


    http.HandleFunc("/getHospital", func(w http.ResponseWriter, r *http.Request) {
    
        res, err := db.Query("select count(okres_lau_kod) as amount, okresy.lat, okresy.lng from pozitivnipripady left join okresy on pozitivnipripady.okres_lau_kod = okresy.code where hospitalizace = '1' and okresy.lat != '' group by okresy.lat,okresy.lng")
        defer res.Close()
    
        if err != nil {
            log.Fatal(err)
        }

        w.Header().Add("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)

        for res.Next() {
            var p Positive
            if err := res.Scan(&p.Count, &p.Lat, &p.Lng); err != nil {  }
		    w.Write([]byte(`{"lat": `+ p.Lat + `, "lng": `+ p.Lng + `, "count": `+ p.Count + `},`))

        }
    }) 

    http.HandleFunc("/getJip", func(w http.ResponseWriter, r *http.Request) {
    
        res, err := db.Query("select count(okres_lau_kod) as amount, okresy.lat, okresy.lng from pozitivnipripady left join okresy on pozitivnipripady.okres_lau_kod = okresy.code where jip = '1' and okresy.lat != '' group by okresy.lat,okresy.lng;")
        defer res.Close()
    
        if err != nil {
            log.Fatal(err)
        }

        w.Header().Add("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)

        for res.Next() {
            var p Positive
            if err := res.Scan(&p.Count, &p.Lat, &p.Lng); err != nil {  }
		    w.Write([]byte(`{"lat": `+ p.Lat + `, "lng": `+ p.Lng + `, "count": `+ p.Count + `},`))

        }
    }) 

    fmt.Printf("Server running (port=8080), route: http://localhost:8080/\n")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }

}

/*

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    
        res, err := db.Query("select count(okres_lau_kod), okres_lau_kod, okresy.name, okresy.lat, okresy.lng from pozitivnipripady left join okresy on pozitivnipripady.okres_lau_kod = okresy.code group by pozitivnipripady.okres_lau_kod,okresy.name,okresy.lat,okresy.lng order by pozitivnipripady.okres_lau_kod")
        defer res.Close()
    
        if err != nil {
            log.Fatal(err)
        }

        pd := PositiveData{
            PageTitle: "Positive",
        }

        for res.Next() {
            var p Positive
            if err := res.Scan(&p.City, &p.Lat, &p.Lng, &p.Count); err != nil {  }
            fmt.Printf("%s", p.Lat)
            pd.Positives = append(pd.Positives, Positive{City: p.City, Lat: p.Lat, Lng: p.Lng, Count: p.Count})
        }        

       tmpl.Execute(w, pd)
    }) 



http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    
    res, err := db.Query("SELECT * FROM pohlavi")
    defer res.Close()

    if err != nil {
        log.Fatal(err)
    }

    gd := GenderData{
        PageTitle: "Genders",
    }

    for res.Next() {
        var g Gender
        if err := res.Scan(&g.Code, &g.Name); err != nil {  }
        gd.Genders = append(gd.Genders, Gender{Code: g.Code, Name: g.Name})
    }        

   tmpl.Execute(w, gd)
}) 
*/