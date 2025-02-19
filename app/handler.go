package app

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"os"
	"unicode/utf8"

	"github.com/antchfx/htmlquery"
	"github.com/gilliek/go-opml/opml"
	"github.com/microcosm-cc/bluemonday"
	"github.com/mmcdole/gofeed"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/html/charset"

	"github.com/yseto/gion-go/db/db"
	"github.com/yseto/gion-go/internal/client"
	"github.com/yseto/gion-go/internal/pin"
)

type ApiServer struct {
	jwtSignedKeyBin []byte
}

func New(jwtSignedKeyBin []byte) *ApiServer {
	return &ApiServer{jwtSignedKeyBin: jwtSignedKeyBin}
}

var _ StrictServerInterface = (*ApiServer)(nil)

func (*ApiServer) PinnedItems(ctx context.Context, request PinnedItemsRequestObject) (PinnedItemsResponseObject, error) {
	pins, err := DBUserFromContext(ctx).PinnedItems()
	if err != nil {
		return nil, err
	}

	items := []PinnedItem{}
	for _, i := range pins {
		items = append(items, PinnedItem{
			FeedId:   i.EntryFeedID,
			Serial:   i.EntrySerial,
			Title:    i.Title,
			UpdateAt: pin.UpdateTime(i.EntryUpdateAt),
			Url:      i.URL,
		})
	}

	return PinnedItems200JSONResponse(items), nil
}

func (*ApiServer) Profile(ctx context.Context, request ProfileRequestObject) (ProfileResponseObject, error) {
	pin, err := DBUserFromContext(ctx).Profile()
	if err != nil {
		return nil, err
	}
	return Profile200JSONResponse(Profile{
		Autoseen:           pin.AutoSeen,
		OnLoginSkipPinList: pin.OnLoginSkipPinList,
		EntryCount:         pin.EntryCount,
		SubstringLength:    pin.SubstringLength,
	}), nil
}

func (*ApiServer) Categories(ctx context.Context, request CategoriesRequestObject) (CategoriesResponseObject, error) {
	cats, err := DBUserFromContext(ctx).Category()
	if err != nil {
		return nil, err
	}

	items := []Category{}
	for _, i := range cats {
		items = append(items, Category{
			ID:   i.ID,
			Name: i.Name,
		})

	}
	return Categories200JSONResponse(items), nil
}

func (*ApiServer) CategoryAndUnreadEntryCount(ctx context.Context, request CategoryAndUnreadEntryCountRequestObject) (CategoryAndUnreadEntryCountResponseObject, error) {
	cat, err := DBUserFromContext(ctx).CategoryAndUnreadEntryCount()
	if err != nil {
		return nil, err
	}

	items := []CategoryAndUnreadEntryCount{}
	for _, i := range cat {
		items = append(items, CategoryAndUnreadEntryCount{
			ID:    i.ID,
			Name:  i.Name,
			Count: i.Count,
		})
	}

	return CategoryAndUnreadEntryCount200JSONResponse(items), nil
}

func (*ApiServer) UnreadEntry(ctx context.Context, request UnreadEntryRequestObject) (UnreadEntryResponseObject, error) {
	db := DBUserFromContext(ctx)
	u, err := db.Profile()
	if err != nil {
		return nil, err
	}

	cat, err := db.UnreadEntryByCategory(request.Id)
	if err != nil {
		return nil, err
	}

	if u.EntryCount > 0 && len(cat) > int(u.EntryCount) {
		cat = cat[:u.EntryCount]
	}

	items := []UnreadEntry{}

	p := bluemonday.NewPolicy()
	for _, i := range cat {
		d := p.Sanitize(i.Description)
		if u.SubstringLength > 0 && uint64(utf8.RuneCountInString(d)) > u.SubstringLength {
			d = string([]rune(d)[:u.SubstringLength])
		}

		items = append(items, UnreadEntry{
			DateEpoch:      uint64(i.PubDate.Unix()),
			Description:    d,
			FeedId:         i.EntryFeedID,
			Readflag:       i.ReadFlag,
			Serial:         i.EntrySerial,
			SiteTitle:      i.SiteTitle,
			SubscriptionId: i.SubscriptionID,
			Title:          i.Title,
			Url:            i.URL,
		})
	}
	return UnreadEntry200JSONResponse(items), nil
}

func (*ApiServer) Subscriptions(ctx context.Context, request SubscriptionsRequestObject) (SubscriptionsResponseObject, error) {
	dbClient := DBUserFromContext(ctx)
	subs, err := dbClient.Subscriptions()
	if err != nil {
		return nil, err
	}
	cat, err := dbClient.Category()
	if err != nil {
		return nil, err
	}

	var resp []Subscription
	for i := range cat {
		var subsOnCategory []CategorySubscription
		for j := range subs {
			if cat[i].ID == subs[j].CategoryID {
				subsOnCategory = append(subsOnCategory, CategorySubscription{
					CategoryId: subs[j].CategoryID,
					HttpStatus: subs[j].HTTPStatus,
					FeedID:     subs[j].FeedID,
					Siteurl:    subs[j].SiteURL,
					Title:      subs[j].FeedTitle,
				})
			}
		}

		resp = append(resp, Subscription{
			CategoryID:   cat[i].ID,
			Name:         cat[i].Name,
			Subscription: subsOnCategory,
		})
	}

	return Subscriptions200JSONResponse(resp), nil
}

func (*ApiServer) SetAsRead(ctx context.Context, request SetAsReadRequestObject) (SetAsReadResponseObject, error) {
	if request.Body == nil {
		return SetAsRead400Response{}, nil
	}

	// return SetAsRead201Response{}, nil // FOR DEBUG

	db := DBUserFromContext(ctx)
	for _, i := range *request.Body {
		if err := db.UpdateEntrySeen(i.FeedID, i.Serial); err != nil {
			return nil, err
		}
	}
	return SetAsRead201Response{}, nil
}

func (*ApiServer) SetPin(ctx context.Context, request SetPinRequestObject) (SetPinResponseObject, error) {
	var readflag pin.ReadFlag
	if request.Body.Readflag == pin.Setpin {
		readflag = pin.Seen
	} else {
		readflag = pin.Setpin
	}

	fmt.Printf("PIN feed_id:%d\tserial:%d\treadflag:%s\n", request.Body.FeedId, request.Body.Serial, readflag)

	tx := DBUserFromContext(ctx).MustBegin()
	if err := tx.UpdateEntry(request.Body.FeedId, request.Body.Serial, readflag); err != nil {
		tx.Rollback()
		return nil, err
	}

	return SetPin200JSONResponse{readflag}, tx.Commit()
}

func (*ApiServer) RegisterCategory(ctx context.Context, request RegisterCategoryRequestObject) (RegisterCategoryResponseObject, error) {
	categoryName := request.Body.Name
	if categoryName == "" {
		return RegisterCategory400Response{}, nil
	}

	tx := DBUserFromContext(ctx).MustBegin()

	cat, err := tx.CategoryByName(categoryName)
	if err != nil && err != sql.ErrNoRows {
		tx.Rollback()
		return nil, err
	}
	if cat != nil {
		return RegisterCategory409Response{}, tx.Commit()
	}

	if err = tx.InsertCategory(categoryName); err != nil {
		tx.Rollback()
		return nil, err
	}

	return RegisterCategory201Response{}, tx.Commit()
}

func insertFeed(ctx context.Context, rssUrl, siteUrl, title string) (*db.Feed, error) {
	tx := DBUserFromContext(ctx).MustBegin()
	feed, err := tx.FeedByUrl(rssUrl, siteUrl)
	if err != nil && err != sql.ErrNoRows {
		tx.Rollback()
		return nil, err
	}
	if feed != nil {
		return feed, tx.Commit()
	}

	err = tx.InsertFeed(rssUrl, siteUrl, title)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	feed, err = tx.FeedByUrl(rssUrl, siteUrl)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return feed, tx.Commit()
}

func (*ApiServer) RegisterSubscription(ctx context.Context, request RegisterSubscriptionRequestObject) (RegisterSubscriptionResponseObject, error) {
	rssUrl, rErr := url.Parse(request.Body.Rss)
	siteUrl, sErr := url.Parse(request.Body.Url)
	title := request.Body.Title
	category := request.Body.Category
	if rErr != nil || sErr != nil || title == "" {
		return RegisterSubscription400Response{}, nil
	}

	feed, err := insertFeed(ctx, rssUrl.String(), siteUrl.String(), title)
	if err != nil {
		return nil, err
	}

	db := DBUserFromContext(ctx)
	tx := db.MustBegin()
	sub, err := tx.SubscriptionByFeedID(feed.ID)
	if err != nil && err != sql.ErrNoRows {
		tx.Rollback()
		return nil, err
	}
	if sub != nil {
		tx.Rollback()
		return RegisterSubscription409Response{}, nil
	}

	cat, err := db.CategoryByID(category)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if cat == nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.InsertSubscription(feed.ID, cat.ID); err != nil {
		tx.Rollback()
		return nil, err
	}

	return RegisterSubscription201Response{}, tx.Commit()
}

func (*ApiServer) DeleteSubscription(ctx context.Context, request DeleteSubscriptionRequestObject) (DeleteSubscriptionResponseObject, error) {
	db := DBUserFromContext(ctx)
	if err := db.DeleteSubscription(request.Id); err != nil {
		return nil, err
	}
	return DeleteSubscription204Response{}, nil
}

func (*ApiServer) DeleteCategory(ctx context.Context, request DeleteCategoryRequestObject) (DeleteCategoryResponseObject, error) {
	db := DBUserFromContext(ctx)
	if err := db.DeleteCategory(request.Id); err != nil {
		return nil, err
	}
	return DeleteCategory204Response{}, nil
}

func (*ApiServer) ChangeSubscription(ctx context.Context, request ChangeSubscriptionRequestObject) (ChangeSubscriptionResponseObject, error) {
	if err := DBUserFromContext(ctx).UpdateSubscription(request.Id, request.Body.Category); err != nil {
		return nil, err
	}
	return ChangeSubscription204Response{}, nil
}

func (*ApiServer) UpdateProfile(ctx context.Context, request UpdateProfileRequestObject) (UpdateProfileResponseObject, error) {
	if request.Body == nil {
		return UpdateProfile400Response{}, nil
	}

	err := DBUserFromContext(ctx).UpdateProfile(db.UserProfile{
		AutoSeen:           request.Body.Autoseen,
		EntryCount:         request.Body.EntryCount,
		OnLoginSkipPinList: request.Body.OnLoginSkipPinList,
		SubstringLength:    request.Body.SubstringLength,
	})

	if err != nil {
		return UpdateProfile400Response{}, nil
	}
	return UpdateProfile204Response{}, nil
}

func (*ApiServer) RemoveAllPin(ctx context.Context, request RemoveAllPinRequestObject) (RemoveAllPinResponseObject, error) {
	if err := DBUserFromContext(ctx).RemovePinnedItem(); err != nil {
		return nil, err
	}
	return RemoveAllPin204Response{}, nil
}

var (
	ErrFeedIsMissing = errors.New("FEED IS MISSING")
)

func examineSubscriptionGetContent(rawUrl string) (*ExamineSubscription, error) {
	urlParam, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}

	resp, err := client.Get(urlParam.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}

	doc, err := htmlquery.Parse(r)
	if err != nil {
		return nil, err
	}

	titleElement := htmlquery.FindOne(doc, "//title")
	if titleElement == nil {
		return nil, ErrFeedIsMissing
	}

	var title = htmlquery.InnerText(titleElement)

	// most blog service is /html/head/link....
	// but youtube /hrml/body/link....
	// each inclusive query.
	rssUrlNode := htmlquery.FindOne(doc, `//link[@type="application/rss+xml"][1]/@href`)
	if rssUrlNode == nil {
		rssUrlNode = htmlquery.FindOne(doc, `//link[@type="application/atom+xml"][1]/@href`)
	}
	if rssUrlNode == nil {
		return nil, ErrFeedIsMissing
	}
	var feedUrl string
	if r := htmlquery.InnerText(rssUrlNode); r != "" {
		u, err := url.Parse(r)
		if err != nil {
			return nil, err
		}
		feedUrl = urlParam.ResolveReference(u).String()
	}
	if len(feedUrl) == 0 {
		return nil, ErrFeedIsMissing
	}

	return &ExamineSubscription{
		URL:   feedUrl,
		Title: title,
	}, nil
}

func (*ApiServer) ExamineSubscription(ctx context.Context, request ExamineSubscriptionRequestObject) (ExamineSubscriptionResponseObject, error) {
	var empty = ExamineSubscription200JSONResponse{
		Success: false,
	}

	content, err := examineSubscriptionGetContent(request.Body.Url)
	if err != nil {
		return empty, nil
	}

	empty.Title = content.Title

	resp, err := client.Get(content.URL)
	if err != nil {
		return empty, nil
	}
	defer resp.Body.Close()

	feed, err := gofeed.NewParser().Parse(resp.Body)
	if err != nil {
		return empty, nil
	}

	var sampleFeeed = []ExamineFeed{}
	for i := range feed.Items {
		var date string
		if dt := feed.Items[i].PublishedParsed; dt != nil {
			date = dt.Format("01/02 15:04")
		}
		sampleFeeed = append(sampleFeeed, ExamineFeed{
			Title: feed.Items[i].Title,
			Url:   feed.Items[i].Link,
			Date:  date,
		})

		if len(sampleFeeed) == 3 {
			break
		}
	}
	content.PreviewFeed = sampleFeeed
	content.Success = true

	return ExamineSubscription200JSONResponse(*content), nil
}

func (*ApiServer) OpmlExport(ctx context.Context, request OpmlExportRequestObject) (OpmlExportResponseObject, error) {
	db := DBUserFromContext(ctx)
	cats, err := db.Category()
	if err != nil {
		return nil, err
	}

	o := opml.Body{}
	for i := range cats {
		feeds, err := db.FeedsByCategoryID(cats[i].ID)
		if err != nil {
			return nil, err
		}

		b := opml.Outline{Text: cats[i].Name, Title: cats[i].Name}
		for j := range feeds {
			b.Outlines = append(b.Outlines, opml.Outline{
				Type:    "rss",
				Text:    feeds[j].Title,
				Title:   feeds[j].Title,
				HTMLURL: feeds[j].SiteURL,
				XMLURL:  feeds[j].URL})
		}
		o.Outlines = append(o.Outlines, b)
	}

	xml, err := opml.OPML{Version: "1.0", Head: opml.Head{Title: "export data"}, Body: o}.XML()
	if err != nil {
		return nil, err
	}

	return OpmlExport200JSONResponse{Xml: xml}, nil
}

func categoryByName(ctx context.Context, categoryName string) (*db.Category, error) {
	db := DBUserFromContext(ctx)
	tx := db.MustBegin()
	cat, err := tx.CategoryByName(categoryName)
	if err != nil && err != sql.ErrNoRows {
		tx.Rollback()
		return nil, err
	}
	if cat != nil {
		return cat, tx.Commit()
	}

	err = tx.InsertCategory(categoryName)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	cat, err = tx.CategoryByName(categoryName)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return cat, tx.Commit()
}

func saveOutline(ctx context.Context, categoryName string, o opml.Outline) error {
	if o.XMLURL == "" || o.HTMLURL == "" || o.Title == "" {
		fmt.Printf("missing required parameter : %+v\n", o)
		return nil
	}

	category, err := categoryByName(ctx, categoryName)
	if err != nil {
		return err
	}

	feed, err := insertFeed(ctx, o.XMLURL, o.HTMLURL, o.Title)
	if err != nil {
		return err
	}

	tx := DBUserFromContext(ctx).MustBegin()

	sub, err := tx.SubscriptionByFeedID(feed.ID)
	if err != nil && err != sql.ErrNoRows {
		tx.Rollback()
		return err
	}
	if sub != nil {
		fmt.Printf("already registered : %s\n", o.Title)
		return tx.Commit()
	}

	if err = tx.InsertSubscription(feed.ID, category.ID); err != nil {
		tx.Rollback()
		return err
	}
	fmt.Printf("registered : %s\n", o.Title)
	return tx.Commit()
}

func walkOutlines(ctx context.Context, categoryName string, o []opml.Outline) error {
	for i := range o {
		if o[i].Type == "rss" {
			if err := saveOutline(ctx, categoryName, o[i]); err != nil {
				return err
			}
		}
		if err := walkOutlines(ctx, o[i].Title, o[i].Outlines); err != nil {
			return err
		}
	}
	return nil
}

func (*ApiServer) OpmlImport(ctx context.Context, request OpmlImportRequestObject) (OpmlImportResponseObject, error) {
	o, err := opml.NewOPML([]byte(request.Body.Xml))
	if err != nil {
		return nil, err
	}

	if err := walkOutlines(ctx, "Default Category", o.Outlines()); err != nil {
		return nil, err
	}

	return OpmlImport201Response{}, nil
}

// https://echo.labstack.com/middleware/jwt/
// https://echo.labstack.com/cookbook/jwt/
func (a *ApiServer) Login(ctx context.Context, request LoginRequestObject) (LoginResponseObject, error) {
	db := DBCommonFromContext(ctx)

	user, err := db.UserByName(request.Body.Id)
	if err != nil {
		return Login401Response{}, nil
	}

	if check := bcrypt.CompareHashAndPassword([]byte(user.Digest), []byte(request.Body.Password)); check != nil {
		return Login401Response{}, nil
	}

	if err := db.UpdateUserLastLogin(user.ID); err != nil {
		fmt.Println(err)
	}

	signedToken, err := GenerateToken(user.ID, a.jwtSignedKeyBin)
	if err != nil {
		return nil, err
	}

	return Login200JSONResponse{
		Autoseen: user.UserProfile.AutoSeen,
		Token:    signedToken,
	}, nil
}

func (*ApiServer) Logout(ctx context.Context, request LogoutRequestObject) (LogoutResponseObject, error) {
	return Logout204Response{}, nil
}

func (*ApiServer) Index(ctx context.Context, request IndexRequestObject) (IndexResponseObject, error) {
	b, err := os.ReadFile("public/index.html")
	if err != nil {
		return nil, err
	}

	return Index200TexthtmlResponse{
		Body: bytes.NewReader(b),
	}, nil
}

func (*ApiServer) ServeRootFile(ctx context.Context, request ServeRootFileRequestObject) (ServeRootFileResponseObject, error) {
	var filename, contentType string
	switch request.Filename {
	case "index.html":
		filename = "public/index.html"
		contentType = "text/html"
	case "gion.js":
		filename = "public/gion.js"
		contentType = "text/javascript"
	case "favicon.ico":
		filename = "public/favicon.ico"
		contentType = "image/x-icon"
	case "apple-touch-icon-precomposed.png":
		filename = "public/apple-touch-icon-precomposed.png"
		contentType = "image/png"
	default:
		return ServeRootFile404Response{}, nil
	}

	b, err := os.ReadFile(filename)
	if err != nil && os.IsNotExist(err) {
		return ServeRootFile404Response{}, nil
	}
	if err != nil {
		return nil, err
	}

	return ServeRootFile200AsteriskResponse{
		Body:        bytes.NewReader(b),
		ContentType: contentType,
	}, nil
}

func (*ApiServer) UpdatePassword(ctx context.Context, request UpdatePasswordRequestObject) (UpdatePasswordResponseObject, error) {
	passwordOld := request.Body.PasswordOld
	password := request.Body.Password
	passwordCheck := request.Body.Passwordc

	if password != passwordCheck || utf8.RuneCountInString(password) < 8 {
		return UpdatePassword400JSONResponse{"error"}, nil
	}

	db := DBUserFromContext(ctx)
	user, err := db.User()
	if err != nil {
		return nil, err
	}

	if check := bcrypt.CompareHashAndPassword([]byte(user.Digest), []byte(passwordOld)); check != nil {
		return UpdatePassword400JSONResponse{"unmatch current password"}, nil
	}

	newDigest, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		return nil, err
	}

	if err := db.UpdateUserDigest(string(newDigest)); err != nil {
		return nil, err
	}
	return UpdatePassword201Response{}, nil
}
