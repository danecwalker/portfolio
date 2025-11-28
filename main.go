package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/danecwalker/portfolio/frontend"
	"github.com/danecwalker/portfolio/internal/notion"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	notionKey := os.Getenv("NOTION_API_KEY")
	if notionKey == "" {
		panic("NOTION_API_KEY not set")
	}

	profilePageId := os.Getenv("PROFILE_PAGE_ID")
	if profilePageId == "" {
		panic("PROFILE_PAGE_ID not set")
	}

	linksDatasourceId := os.Getenv("LINKS_DATASOURCE_ID")
	if linksDatasourceId == "" {
		panic("LINKS_DATASOURCE_ID not set")
	}

	experienceDatasourceId := os.Getenv("EXPERIENCE_DATASOURCE_ID")
	if experienceDatasourceId == "" {
		panic("EXPERIENCE_DATASOURCE_ID not set")
	}

	projectsDatasourceId := os.Getenv("PROJECTS_DATASOURCE_ID")
	if projectsDatasourceId == "" {
		panic("PROJECTS_DATASOURCE_ID not set")
	}

	affiliationsDatasourceId := os.Getenv("AFFILIATIONS_DATASOURCE_ID")
	if affiliationsDatasourceId == "" {
		panic("AFFILIATIONS_DATASOURCE_ID not set")
	}

	client := notion.NewClient(notionKey)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/profile", func(w http.ResponseWriter, r *http.Request) {
		pageReq := notion.Page(profilePageId)
		pageResp, err := pageReq.Fetch(client)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"profileImageUrl": pageResp.Cover.GetURL(),
		})
	})

	mux.HandleFunc("/api/v1/links", func(w http.ResponseWriter, r *http.Request) {
		linksReq := notion.Datasource(linksDatasourceId).Query(map[string]any{
			"property": "Hidden",
			"checkbox": map[string]any{
				"equals": false,
			},
		}, []map[string]any{
			{
				"property":  "Display Order",
				"direction": "ascending",
			},
		})
		linksResp, err := linksReq.Fetch(client)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var links []map[string]any
		for _, row := range linksResp.Results {
			title_row := row.Properties["Name"]
			var title string
			if len(title_row.Title) > 0 {
				title = title_row.Title[0].PlainText
			} else {
				continue
			}
			url := row.Properties["Website URL"].Url
			files := row.Properties["File"].Files

			if len(files) > 0 {
				file := files[0].GetURL()
				links = append(links, map[string]any{
					"title": title,
					"type":  "file",
					"url":   file,
				})
			} else {
				links = append(links, map[string]any{
					"title": title,
					"type":  "link",
					"url":   *url,
				})
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"links": links,
		})
	})

	mux.HandleFunc("/api/v1/content/{pageId}", func(w http.ResponseWriter, r *http.Request) {
		pageId := r.PathValue("pageId")
		blocksReq := notion.Blocks(pageId).Query()
		blocksResp, err := blocksReq.Fetch(client)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var content strings.Builder
		for _, row := range blocksResp.Results {
			switch row.Type {
			case "paragraph":
				// check if we should close a bulleted list
				if strings.HasSuffix(content.String(), "</li>\n") {
					content.WriteString("</ul>\n")
				}
				var paragraphText []string
				for _, rt := range row.Paragraph.RichText {
					var p string = rt.PlainText
					applyBold(&p, rt.Annotations.Bold)
					applyItalic(&p, rt.Annotations.Italic)
					applyStrikethrough(&p, rt.Annotations.Strikethrough)
					applyUnderline(&p, rt.Annotations.Underline)
					applyCode(&p, rt.Annotations.Code)
					paragraphText = append(paragraphText, p)
				}
				content.WriteString("<p>" + strings.Join(paragraphText, " ") + "</p>\n")
			case "heading_1":
				if strings.HasSuffix(content.String(), "</li>\n") {
					content.WriteString("</ul>\n")
				}
				var headingText []string
				for _, rt := range row.Heading1.RichText {
					headingText = append(headingText, rt.PlainText)
				}
				content.WriteString("<h1>" + strings.Join(headingText, " ") + "</h1>\n")
			case "heading_2":
				if strings.HasSuffix(content.String(), "</li>\n") {
					content.WriteString("</ul>\n")
				}
				var headingText []string
				for _, rt := range row.Heading2.RichText {
					headingText = append(headingText, rt.PlainText)
				}
				content.WriteString("<h2>" + strings.Join(headingText, " ") + "</h2>\n")
			case "heading_3":
				if strings.HasSuffix(content.String(), "</li>\n") {
					content.WriteString("</ul>\n")
				}
				var headingText []string
				for _, rt := range row.Heading3.RichText {
					headingText = append(headingText, rt.PlainText)
				}
				content.WriteString("<h3>" + strings.Join(headingText, " ") + "</h3>\n")
			case "bulleted_list_item":
				var listItemText []string
				for _, rt := range row.BulletedListItem.RichText {
					var p string = rt.PlainText
					applyBold(&p, rt.Annotations.Bold)
					applyItalic(&p, rt.Annotations.Italic)
					applyStrikethrough(&p, rt.Annotations.Strikethrough)
					applyUnderline(&p, rt.Annotations.Underline)
					applyCode(&p, rt.Annotations.Code)
					listItemText = append(listItemText, p)
				}
				// if previous block was also a bulleted list item, continue the list
				if strings.HasSuffix(content.String(), "</li>\n") {
					content.WriteString("<li>" + strings.Join(listItemText, " ") + "</li>\n")
				} else {
					content.WriteString("<ul>\n<li>" + strings.Join(listItemText, " ") + "</li>\n")
				}
			default:
				// Handle other block types as needed
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"content": content.String(),
		})
	})

	mux.HandleFunc("/api/v1/experience", func(w http.ResponseWriter, r *http.Request) {
		experienceReq := notion.Datasource(experienceDatasourceId).Query(map[string]any{
			"property": "Hidden",
			"checkbox": map[string]any{
				"equals": false,
			},
		}, []map[string]any{
			{
				"property":  "Date",
				"direction": "descending",
			},
		})
		experienceResp, err := experienceReq.Fetch(client)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var experience []map[string]any
		for _, row := range experienceResp.Results {
			logo := row.Cover.GetURL()
			title_row := row.Properties["Company"]
			var title string
			if len(title_row.Title) > 0 {
				title = title_row.Title[0].PlainText
			} else {
				continue
			}
			position := row.Properties["Position"].RichText[0].PlainText
			date := row.Properties["Date"].Date

			experience = append(experience, map[string]any{
				"company":  title,
				"position": position,
				"logoURL":  logo,
				"start":    date.Start,
				"end":      date.End,
				"current":  date.End == nil,
				"pageId":   row.ID,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"experience": experience,
		})
	})

	mux.HandleFunc("/api/v1/projects", func(w http.ResponseWriter, r *http.Request) {
		projectsReq := notion.Datasource(projectsDatasourceId).Query(map[string]any{
			"property": "Hidden",
			"checkbox": map[string]any{
				"equals": false,
			},
		}, []map[string]any{
			{
				"property":  "Date",
				"direction": "descending",
			},
		})
		projectsResp, err := projectsReq.Fetch(client)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var projects []map[string]any
		for _, row := range projectsResp.Results {
			imageUrl := row.Cover.GetURL()
			title_row := row.Properties["Name"]
			var title string
			if len(title_row.Title) > 0 {
				title = title_row.Title[0].PlainText
			} else {
				continue
			}
			projectUrl := row.Properties["Project URL"].Url
			date := row.Properties["Date"].Date

			projects = append(projects, map[string]any{
				"title":           title,
				"projectImageURL": imageUrl,
				"projectURL":      projectUrl,
				"date":            date.Start,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"projects": projects,
		})
	})

	mux.HandleFunc("/api/v1/affiliations", func(w http.ResponseWriter, r *http.Request) {
		affiliationsReq := notion.Datasource(affiliationsDatasourceId).Query(map[string]any{
			"property": "Hidden",
			"checkbox": map[string]any{
				"equals": false,
			},
		}, []map[string]any{
			{
				"property":  "Display Order",
				"direction": "ascending",
			},
		})
		affiliationsResp, err := affiliationsReq.Fetch(client)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var affiliations []map[string]any
		for _, row := range affiliationsResp.Results {
			logo := row.Cover.GetURL()
			title_row := row.Properties["Name"]
			var title string
			if len(title_row.Title) > 0 {
				title = title_row.Title[0].PlainText
			} else {
				continue
			}
			affiliations = append(affiliations, map[string]any{
				"title":   title,
				"logoURL": logo,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"affiliations": affiliations,
		})
	})

	mux.Handle("/", frontend.SvelteKitHandler())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: CorsMiddleware(mux),
	}

	fmt.Println("Starting server on :8080")
	server.ListenAndServe()
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func applyBold(text *string, bold bool) {
	if bold {
		*text = "<strong>" + *text + "</strong>"
	}
}

func applyItalic(text *string, italic bool) {
	if italic {
		*text = "<em>" + *text + "</em>"
	}
}

func applyUnderline(text *string, underline bool) {
	if underline {
		*text = "<u>" + *text + "</u>"
	}
}

func applyStrikethrough(text *string, strikethrough bool) {
	if strikethrough {
		*text = "<s>" + *text + "</s>"
	}
}

func applyCode(text *string, code bool) {
	if code {
		*text = "<code>" + *text + "</code>"
	}
}
