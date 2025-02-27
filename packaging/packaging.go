package packaging

import (
	"errors"
	"homework/order"
	"strings"
)

const (
	noPackaging   = "none"
	bagPackaging  = "bag"
	boxPackaging  = "box"
	wrapPackaging = "wrap"
)

const (
	bagCost  float64 = 5.00
	boxCost  float64 = 20.00
	wrapCost float64 = 1.00
)

const (
	bagMinWeight float64 = 10.00
	boxMinWeight float64 = 30.00
)

var (
	errWrongPackagingType error = errors.New("wrong packaging type")
	errNotEnoughWeight          = errors.New("not enough weight")
	errWrongPackaging           = errors.New("wrong packaging")
	// errWrongPackagingOrder       = errors.New("wrong packaging order")
)

type Packaging interface {
	String() string
	Pack(order *order.Order) error
}

func GetPackaging(packagingType string) (Packaging, error) {
	switch packagingType {
	case bagPackaging:
		return &Bag{}, nil
	case boxPackaging:
		return &Box{}, nil
	case wrapPackaging:
		return &Wrap{}, nil
	default:
		return nil, errWrongPackagingType
	}
}

func FormPackagingString(packagings []Packaging) string {
	if len(packagings) == 0 {
		return noPackaging
	}

	var sb strings.Builder
	for i, packaging := range packagings {
		sb.WriteString(packaging.String())
		if i != len(packagings)-1 {
			sb.WriteString(" in ")
		}
	}

	return sb.String()
}

func CheckPackaging(packagings []Packaging) error {
	if len(packagings) > 1 && (packagings[1].String() != wrapPackaging || packagings[0].String() == wrapPackaging) {
		return errWrongPackaging
	}

	return nil
}

type Bag struct {
}

func (b *Bag) String() string {
	return bagPackaging
}

func (b *Bag) Pack(order *order.Order) error {
	if order.Weight < bagMinWeight {
		return errNotEnoughWeight
	}

	order.Price += bagCost

	return nil
}

type Box struct {
}

func (b *Box) String() string {
	return boxPackaging
}

func (b *Box) Pack(order *order.Order) error {
	if order.Weight < boxMinWeight {
		return errNotEnoughWeight
	}

	order.Price += boxCost

	return nil
}

type Wrap struct {
}

func (w *Wrap) String() string {
	return wrapPackaging
}

func (w *Wrap) Pack(order *order.Order) error {
	order.Price += wrapCost

	return nil
}
