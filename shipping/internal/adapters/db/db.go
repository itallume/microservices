package db

import (
    "fmt"

    "github.com/itallume/microservices/shipping/internal/application/core/domain"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

type Shipping struct {
    gorm.Model
    BillId       int64
    Items        []Item
}

type Item struct {
    gorm.Model
    ProductCode string
    UnitPrice   float32
    Quantity    int32
    ShippingID  uint
}

type Adapter struct {
    db *gorm.DB
}

func NewAdapter(dataSourceUrl string) (*Adapter, error) {
    db, openErr := gorm.Open(mysql.Open(dataSourceUrl), &gorm.Config{})
    if openErr != nil {
        return nil, fmt.Errorf("db connection error: %v", openErr)
    }
    err := db.AutoMigrate(&Shipping{}, Item{})
    if err != nil {
        return nil, fmt.Errorf("db migration error: %v", err)
    }
    return &Adapter{db: db}, nil
}

func (a Adapter) Get(id string) (domain.Shipping, error) {
    var shippingEntity Shipping
    res := a.db.First(&shippingEntity, id)
    var items []domain.Item
    for _, item := range shippingEntity.Items {
        items = append(items, domain.Item{
            ProductCode: item.ProductCode,
            UnitPrice:   item.UnitPrice,
            Quantity:    item.Quantity,
        })
    }
    shipping := domain.Shipping{
        ID:           int64(shippingEntity.ID),
        BillId:       shippingEntity.BillId,
        Items:        items,
        CreatedAt:    shippingEntity.CreatedAt.UnixNano(),
    }
    return shipping, res.Error
}

func (a Adapter) Save(shipping *domain.Shipping) error {
    var items []Item
    for _, item := range shipping.Items {
        items = append(items, Item{
            ProductCode: item.ProductCode,
            UnitPrice:   item.UnitPrice,
            Quantity:    item.Quantity,
        })
    }
    shippingModel := Shipping{
        BillId:       shipping.BillId,
        Items:        items,
    }
    res := a.db.Create(&shippingModel)
    if res.Error == nil {
        shipping.ID = int64(shippingModel.ID)
    }
    return res.Error
}