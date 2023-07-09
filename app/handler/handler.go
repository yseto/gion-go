package handler

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
