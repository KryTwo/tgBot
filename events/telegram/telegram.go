package telegram

import (
	"errors"
	"fmt"
	"main/clients/telegram"
	"main/events"
	"main/lib/e"
	"main/storage"
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID   int
	Username string
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func New(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

// Fetch получает события
func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	// получим все апдейты
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap("can't get events", err)
	}
	// если апдейтов нет - возвращаем ничего
	if len(updates) == 0 {
		return nil, nil
	}
	// аллоцируем память под результат
	res := make([]events.Event, 0, len(updates))

	// перебираем все значения и преобразуем их в тип event
	for _, u := range updates {
		res = append(res, event(u))
	}

	// обновляем параметр offset
	p.offset = updates[len(updates)-1].ID + 1

	fmt.Println(res)
	return res, nil
}

// Process обрабатывает события
func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return e.Wrap("can't process message", ErrUnknownEventType)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	// поскольку мы здесь работаем не с update а с event , нам придется первым делом получить мету
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("can't process message", err)
	}

	if err := p.doCmd(event.Text, meta.ChatID, meta.Username); err != nil {
		return e.Wrap("can't process message", err)
	}

	return nil
}

func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)
	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}
	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}
	return res
}

func meta(event events.Event) (Meta, error) {
	// для поля мета мы сделаем т.н. TypeAssertion и если здесь будет что то другое то вторым параметром вернется false
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("can't get meta", ErrUnknownMetaType)
	}

	return res, nil
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}
	return upd.Message.Text
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}

	return events.Message
}
