package main

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/mail"
	"net/textproto"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/kothawoc/go-nntp"
	"github.com/kothawoc/kothawoc/pkg/messages"
	serr "github.com/kothawoc/kothawoc/pkg/serror"
)

const createArticleIndex string = `
CREATE TABLE IF NOT EXISTS articles (
	timestamp INT NOT NULL DEFAULT 0,
	messageid TEXT NOT NULL UNIQUE,
	refs TEXT NOT NULL DEFAULT ""
	);
`

type ArticleDB struct {
	*sql.DB
}

func NewArticleDB(path string) (*ArticleDB, error) {

	db, err := sql.Open("sqlite3", path+"/gui-article.db")
	if err != nil {
		return nil, serr.New(err)
	}

	if _, err := db.Exec(createArticleIndex); err != nil {
		slog.Info("FAILED Create DB database query", "error", err, "path", path, "query", createArticleIndex)
		return nil, serr.New(err)
	}
	if _, err := db.Exec("PRAGMA journal_mode = WAL;pragma synchronous = normal;pragma temp_store = memory;pragma mmap_size = 30000000000;pragma page_size = 32768;pragma auto_vacuum = incremental;pragma incremental_vacuum;"); err != nil {
		slog.Info("FAILED PRAGMA journal_mode = WAL;pragma synchronous = normal;pragma temp_store = memory;pragma mmap_size = 30000000000;pragma page_size = 32768;pragma auto_vacuum = incremental;pragma incremental_vacuum;", "error", err, "path", path, "query", createArticleIndex)
		return nil, serr.New(err)
	}

	a := &ArticleDB{db}

	go a.run()
	return a, nil
}

func (a *ArticleDB) Insert(time int, messageid, references string) error {
	_, err := a.Exec("INSERT INTO ARTICLES (timestamp,messageid,refs) VALUES(?,?,?);", time, messageid, references)

	//slog.Info("Inserting", "error", err, "time", time, "messageid", messageid, "references", references)
	return serr.New(err)
}

func (a *ArticleDB) GetLength() (int64, error) {
	length := int64(0)
	row := a.QueryRow("SELECT COUNT(timestamp) FROM articles WHERE refs=\"\";")
	err := row.Scan(&length)
	return length, serr.New(err)
}

func (a *ArticleDB) GetItem(num int) (*messages.MessageTool, error) {
	messageid := ""
	row := a.QueryRow("SELECT messageid FROM articles ORDER BY timestamp DESC LIMIT 1 OFFSET ?;", num)
	err := row.Scan(&messageid)
	if err != nil {
		return nil, serr.Errorf("cant find message number [%q]", err)
	}
	_, _, rawArticle, err := kc.NNTPclient.Article(messageid)

	if err != nil {
		return nil, serr.New(err)
	}
	mailMsg, err := mail.ReadMessage(rawArticle)
	if err != nil {
		return nil, serr.New(err)
	}

	article := &nntp.Article{
		Header: textproto.MIMEHeader(mailMsg.Header),
		Body:   mailMsg.Body,
		//	Bytes:  len([]byte(body)),
		//	Lines:  strings.Count(body, "\n"),
	}

	msg := messages.NewMessageToolFromArticle(article)

	return msg, serr.New(err)
}

func (a *ArticleDB) run() {
	for {
		groups, err := kc.NNTPclient.List("")
		if err != nil {
			slog.Warn("Can't list groups")
		}
		for _, group := range groups {

			a, err := kc.NNTPclient.Group(group.Name)

			if err != nil {
				log.Printf("Error in listing newsgroup conftent: [%v]", err)
				return
			}
			//tmp, err := kc.NNTPclient.ListOverviewFmt()
			//log.Printf("listoverviewformat [%#v][%v]", tmp, err)
			res, err := kc.NNTPclient.Over(int(a.Low), int(a.High))
			for _, art := range res {
				//date := time.art.Date

				//s := "Tue Sep 16 21:58:58 +0000 2014"
				// Date: Tue, 10 Sep 2024 15:11:19 +0000
				const rfc2822 = "Mon, 02 Jan 2006 15:04:05 -0700 "
				//  MessageId:"Sat, 31 Aug 2024 19:59:07 +0000", References:"<1jd6tgb-snjkpy2l
				slog.Info(fmt.Sprintf("[%#v][%v]\n", art, art))
				// actually the data TODO FIX THIS BUG
				t, err := time.Parse(rfc2822, art.Date)
				if err != nil {
					slog.Info("Error", "error", err)
					return
				}

				ADB.Insert(int(t.UTC().Unix()), art.MessageId, art.References)
				//fmt.Println(u)
				//f := t.Format(time.UnixDate)
				//fmt.Println(f)
			}
		}

		slog.Info("looping")
		time.Sleep(time.Second * 5)
	}
}

func displayMainFeed() {

	var List *widget.List

	list := widget.NewList(
		func() int {
			length, err := ADB.GetLength()
			if err != nil {
				slog.Error("Article db, can't get length", "error", err)
			}
			return int(length)
		},
		func() fyne.CanvasObject {
			vbox := container.NewVBox()
			tl := widget.NewRichTextFromMarkdown("From")
			tc := widget.NewRichTextFromMarkdown("Subject")
			tr := widget.NewRichTextFromMarkdown("Date")
			bbox := container.NewBorder(nil, nil, tl, tr, tc)
			vbox.Add(bbox)
			body := widget.NewRichTextFromMarkdown("Content")
			body.Wrapping = fyne.TextWrapWord
			vbox.Add(body)
			return vbox
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {

			msg, _ := ADB.GetItem(i)
			shortFrom := msg.Article.Header.Get("From")
			shortFrom = shortFrom[:3] + "..." + shortFrom[len(shortFrom)-4:]
			body := fmt.Sprintf("%s\r\n%s", msg.Preamble, msg.Parts)
			splitDate := strings.Split(msg.Article.Header.Get("Date"), " ")
			date := strings.Join(splitDate[1:len(splitDate)-1], " ")

			o.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.RichText).ParseMarkdown("**" + msg.Article.Header.Get("Subject") + "**")
			o.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.RichText).ParseMarkdown(shortFrom)
			o.(*fyne.Container).Objects[0].(*fyne.Container).Objects[2].(*widget.RichText).ParseMarkdown(date)
			o.(*fyne.Container).Objects[1].(*widget.RichText).ParseMarkdown(body)

			height := o.(*fyne.Container).Objects[0].Size().Height + o.(*fyne.Container).Objects[1].Size().Height

			List.SetItemHeight(i, height)
		})

	List = list

	parent.Trailing = list

	parent.Trailing.Refresh()
	parent.Refresh()

}
