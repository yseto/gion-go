package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"unicode/utf8"

	"github.com/antchfx/htmlquery"
	"github.com/hakobe/paranoidhttp"
	"github.com/microcosm-cc/bluemonday"
	"github.com/mmcdole/gofeed"
	"golang.org/x/net/html/charset"

	"github.com/yseto/gion-go/db/db"
	"github.com/yseto/gion-go/internal/pin"
)

type ApiServer struct{}

func NewApiServer() *ApiServer {
	return &ApiServer{}
}

var _ StrictServerInterface = (*ApiServer)(nil)

func (*ApiServer) PinnedItems(ctx context.Context, request PinnedItemsRequestObject) (PinnedItemsResponseObject, error) {
	pins, err := DBUserFromContext(ctx).PinnedItems()
	if err != nil {
		return PinnedItems400Response{}, nil
	}

	items := []PinnedItem{}
	for _, i := range pins {
		items = append(items, PinnedItem{
			FeedId:   i.EntryFeedID,
			Serial:   i.EntrySerial,
			Title:    i.Title,
			UpdateAt: i.EntryUpdateAt,
			Url:      i.URL,
		})
	}

	return PinnedItems200JSONResponse(items), nil
}

func (*ApiServer) Profile(ctx context.Context, request ProfileRequestObject) (ProfileResponseObject, error) {
	pin, err := DBUserFromContext(ctx).Profile()
	if err != nil {
		return Profile400Response{}, nil
	}
	return Profile200JSONResponse(Profile{
		Autoseen:  pin.AutoSeen,
		Nopinlist: pin.NoPinList,
		Numentry:  pin.EntryCount,
		Numsubstr: pin.SubstringLength,
	}), nil
}

func (*ApiServer) Categories(ctx context.Context, request CategoriesRequestObject) (CategoriesResponseObject, error) {
	cats, err := DBUserFromContext(ctx).Category()
	if err != nil {
		return Categories400Response{}, nil
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
		return CategoryAndUnreadEntryCount400Response{}, nil
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
		return UnreadEntry400Response{}, nil
	}

	cat, err := db.UnreadEntryByCategory(request.Body.Category)
	if err != nil {
		return UnreadEntry400Response{}, nil
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
			Readflag:       i.ReadFlag.ToPinReadFlag(),
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
		return Subscriptions400Response{}, nil
	}
	cat, err := dbClient.Category()
	if err != nil {
		return Subscriptions400Response{}, nil
	}

	var resp []Subscription
	for i := range cat {
		var subsOnCategory []SubscriptionForUser
		for j := range subs {
			if cat[i].ID == subs[j].CategoryID {
				subsOnCategory = append(subsOnCategory, SubscriptionForUser{
					CategoryId: subs[j].CategoryID,
					HttpStatus: subs[j].HTTPStatus,
					FeedID:     subs[j].FeedID,
					Siteurl:    subs[j].SiteURL,
					Title:      subs[j].FeedTitle,
				})
			}
		}

		resp = append(resp, Subscription{
			ID:           cat[i].ID,
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

	// return SetAsRead200JSONResponse{Result: "OK"}, nil // FOR DEBUG

	db := DBUserFromContext(ctx)
	for _, i := range *request.Body {
		err := db.UpdateEntrySeen(i.FeedID, i.Serial)
		if err != nil {
			return SetAsRead400Response{}, nil
		}
	}
	return SetAsRead200JSONResponse{Result: "OK"}, nil
}

func (*ApiServer) SetPin(ctx context.Context, request SetPinRequestObject) (SetPinResponseObject, error) {
	var readflag db.ReadFlag
	if request.Body.Readflag == pin.Setpin {
		readflag = db.Seen
	} else {
		readflag = db.SetPin
	}

	fmt.Printf("PIN feed_id:%d\tserial:%d\treadflag:%d\n", request.Body.FeedId, request.Body.Serial, readflag)

	tx := DBUserFromContext(ctx).MustBegin()
	if tx.UpdateEntry(request.Body.FeedId, request.Body.Serial, readflag) != nil {
		tx.Rollback()
		return SetPin400Response{}, nil
	}
	tx.Commit()

	return SetPin200JSONResponse{readflag.ToPinReadFlag()}, nil
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
		return RegisterCategory400Response{}, nil
	}
	if cat != nil {
		tx.Commit()
		return RegisterCategory200JSONResponse{Result: "ERROR_ALREADY_REGISTER"}, nil
	}

	if err = tx.InsertCategory(categoryName); err != nil {
		tx.Rollback()
		return RegisterCategory400Response{}, nil
	}
	tx.Commit()

	return RegisterCategory200JSONResponse{Result: "OK"}, nil
}

func insertFeed(ctx context.Context, rssUrl, siteUrl, title string) (*db.Feed, error) {
	tx := DBUserFromContext(ctx).MustBegin()
	feed, err := tx.FeedByUrl(rssUrl, siteUrl)
	if err != nil && err != sql.ErrNoRows {
		tx.Rollback()
		return nil, err
	}
	if feed != nil {
		tx.Commit()
		return feed, nil
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
	tx.Commit()
	return feed, nil
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
		return RegisterSubscription400Response{}, nil
	}

	db := DBUserFromContext(ctx)
	tx := db.MustBegin()
	sub, err := tx.SubscriptionByFeedID(feed.ID)
	if err != nil && err != sql.ErrNoRows {
		tx.Rollback()
		return RegisterSubscription400Response{}, nil
	}
	if sub != nil {
		tx.Rollback()
		return RegisterSubscription200JSONResponse{"ERROR_ALREADY_REGISTER"}, nil
	}

	cat, err := db.CategoryByID(category)
	if err != nil {
		tx.Rollback()
		return RegisterSubscription400Response{}, nil
	}
	if cat == nil {
		tx.Rollback()
		return RegisterSubscription400Response{}, nil
	}
	if tx.InsertSubscription(feed.ID, cat.ID) != nil {
		tx.Rollback()
		return RegisterSubscription400Response{}, nil
	}
	tx.Commit()

	return RegisterSubscription200JSONResponse{"OK"}, nil
}

func (*ApiServer) DeleteSubscription(ctx context.Context, request DeleteSubscriptionRequestObject) (DeleteSubscriptionResponseObject, error) {
	deleteType := request.Body.Subscription
	id := request.Body.Id
	if deleteType == "" {
		return DeleteSubscription400Response{}, nil
	}

	var err error
	db := DBUserFromContext(ctx)
	switch deleteType {
	case DeleteSubscriptionFormdataBodySubscriptionCategory:
		err = db.DeleteCategory(id)
	case DeleteSubscriptionFormdataBodySubscriptionEntry:
		err = db.DeleteSubscription(id)
	default:
		err = fmt.Errorf("invalid type")
	}
	if err != nil {
		return DeleteSubscription400Response{}, nil
	}
	return DeleteSubscription200JSONResponse{Result: "OK"}, nil
}

func (*ApiServer) ChangeSubscription(ctx context.Context, request ChangeSubscriptionRequestObject) (ChangeSubscriptionResponseObject, error) {
	if DBUserFromContext(ctx).UpdateSubscription(request.Body.Id, request.Body.Category) != nil {
		return ChangeSubscription400Response{}, nil
	}
	return ChangeSubscription200JSONResponse{Result: "OK"}, nil
}

func (*ApiServer) UpdateProfile(ctx context.Context, request UpdateProfileRequestObject) (UpdateProfileResponseObject, error) {
	if request.Body == nil {
		return UpdateProfile400Response{}, nil
	}

	err := DBUserFromContext(ctx).UpdateProfile(db.UserProfile{
		AutoSeen:        request.Body.Autoseen,
		EntryCount:      request.Body.Numentry,
		NoPinList:       request.Body.Nopinlist,
		SubstringLength: request.Body.Numsubstr,
	})

	if err != nil {
		return UpdateProfile400Response{}, nil
	}
	return UpdateProfile200JSONResponse{Result: "OK"}, nil
}

func (*ApiServer) RemoveAllPin(ctx context.Context, request RemoveAllPinRequestObject) (RemoveAllPinResponseObject, error) {
	if DBUserFromContext(ctx).RemovePinnedItem() != nil {
		return RemoveAllPin400Response{}, nil
	}
	return RemoveAllPin200JSONResponse{Result: "OK"}, nil
}

var (
	ErrFeedIsMissing = errors.New("FEED IS MISSING")
)

func examineSubscriptionGetContent(rawUrl string) (*ExamineSubscription, error) {
	urlParam, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}

	resp, err := paranoidhttp.DefaultClient.Get(urlParam.String())
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
		return nil, err
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

	resp, err := paranoidhttp.DefaultClient.Get(content.URL)
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
		return OpmlExport400Response{}, nil
	}

	o := opml.Body{}
	for i := range cats {
		feeds, err := db.FeedsByCategoryID(cats[i].ID)
		if err != nil {
			return OpmlExport400Response{}, nil
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
		return OpmlExport400Response{}, nil
	}

	return OpmlExport200JSONResponse{Xml: xml}, nil
}
