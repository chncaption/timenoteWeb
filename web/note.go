package web

import (
	"github.com/gin-gonic/gin"
	"strconv"
	. "timenoteWeb/config"
	"timenoteWeb/loader"
)

func NoteListPage(c *gin.Context) {
	data := loader.LoadLastDataFile()
	notes := make([]simpleNote, len(data.Notes))
	for i, note := range data.Notes {
		notes[i] = simpleNote{
			ID:           strconv.FormatInt(note.ID, 10),
			Title:        note.Title,
			Date:         note.GetDateStr(),
			Weather:      note.GetWeatherStr(),
			WeatherEmoji: note.GetWeatherEmoji(),
			Mood:         note.GetMoodStr(),
			MoodEmoji:    note.GetMoodEmoji(),
			Location:     note.Location,
			CategoryID:   strconv.FormatInt(note.CategoryID, 10),
			CategoryName: func() string {
				c, s := data.FindCategory(note)
				if !s {
					return ""
				} else {
					return c.CategoryName
				}
			}(),
		}
	}
	var pData noteListData
	pData.Notes = notes
	pData.Title = "日记列表"
	pData.Nickname = AppConfig.Web.Nickname
	pData.Source = data.Source
	c.HTML(200, "notes.html", pData)
}

func NotePage(c *gin.Context) {
	id := c.Param("id")
	data := loader.LoadLastDataFile()
	var nData note
	for _, n := range data.Notes {
		if strconv.FormatInt(n.ID, 10) == id {
			nData.Title = n.Title
			nData.Content = n.GetContentHTML()
			nData.ID = strconv.FormatInt(n.ID, 10)
			nData.Date = n.GetDateStr()
			nData.Weather = n.GetWeatherStr()
			nData.WeatherEmoji = n.GetWeatherEmoji()
			nData.Mood = n.GetMoodStr()
			nData.MoodEmoji = n.GetMoodEmoji()
			nData.CategoryID = strconv.FormatInt(n.CategoryID, 10)
			nData.CategoryName = func() string {
				c, s := data.FindCategory(n)
				if !s {
					return ""
				} else {
					return c.CategoryName
				}
			}()
			nData.Location = n.Location
		}
	}
	var pData notePageData
	pData.Note = nData
	pData.Title = "日记《" + nData.Title + "》"
	pData.Nickname = AppConfig.Web.Nickname
	pData.Source = data.Source
	c.HTML(200, "note.html", pData)
}