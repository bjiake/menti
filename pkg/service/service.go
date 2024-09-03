package service

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"menti/pkg/db"
	"menti/pkg/domain/note"
	accountI "menti/pkg/repo/account/interface"
	noteI "menti/pkg/repo/note/interface"
	"net/http"
	"net/url"

	"context"
	interfaces "menti/pkg/service/interface"
	"strconv"
)

type service struct {
	rAccount accountI.AccountRepository
	rNote    noteI.NoteRepository
}

func NewService(
	accountRepository accountI.AccountRepository,
	noteRepository noteI.NoteRepository,
) interfaces.ServiceUseCase {
	return &service{
		rAccount: accountRepository,
		rNote:    noteRepository,
	}
}

func (s *service) Migrate(ctx context.Context) error {
	if err := s.rAccount.Migrate(ctx); err != nil {
		return err
	}
	if err := s.rNote.Migrate(ctx); err != nil {
		return err
	}
	log.Info("all migrate success")
	return nil
}

func (s *service) checkLogin(ctx context.Context, userId string) (int64, error) {
	id, err := strconv.ParseInt(userId, 10, 64)
	if err != nil || id <= 0 {
		log.Error(err)
		return 0, db.ErrAuthorize
	}

	err = s.rAccount.CheckAccount(ctx, id)
	if err != nil {
		log.Error(err)
		return 0, db.ErrAuthorize
	}
	log.Info("check login success", id)
	return id, nil
}

func (s *service) checkText(note note.Note) error {
	urlStr := "https://speller.yandex.net/services/spellservice.json/checkTexts"
	params := url.Values{}
	params.Add("text", note.Name)
	params.Add("text", note.Content)

	u, err := url.Parse(urlStr)
	if err != nil {
		log.Error(err)
		return err
	}

	u.RawQuery = params.Encode()
	client := &http.Client{}

	resp, err := client.Get(u.String())
	if err != nil {
		log.Error(err)
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return err
	}

	var response [][]struct {
		Code        int      `json:"code"`
		Pos         int      `json:"pos"`
		Row         int      `json:"row"`
		Col         int      `json:"col"`
		Len         int      `json:"len"`
		Word        string   `json:"word"`
		Suggestions []string `json:"s"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Error(err)
		return err
	}

	// Check if the response is `[[], []]`
	if len(response) == 2 && len(response[0]) == 0 && len(response[1]) == 0 {
		return nil
	}

	return db.ErrYandexSpeller
}
