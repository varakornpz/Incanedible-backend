package models


import (
    "time"
	"gorm.io/gorm"
    "encoding/json"
    "github.com/google/uuid"
    "database/sql/driver"
    "errors"
)

type RegisteredCane struct {
    Name    string  `json:"name"`
    CaneID      string  `json:"cane_id"`
}

type RegisteredCanes []RegisteredCane


func (a RegisteredCanes) Value() (driver.Value, error) {
    if len(a) == 0 {
        return "[]" , nil
    }
    return json.Marshal(a)
}

func (a *RegisteredCanes) Scan(value interface{}) error {
	if value == nil {
		*a = make(RegisteredCanes, 0)
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		str, okStr := value.(string)
		if !okStr {
			return errors.New("type assertion to []byte or string failed")
		}
		b = []byte(str)
	}

	return json.Unmarshal(b, a)
}

type User struct{
    CreatedAt      time.Time
    UpdatedAt      time.Time
    DeletedAt      gorm.DeletedAt `gorm:"index"`
    UUID            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"` 
    Email          string
    Name            string
    ProfilePic     string   `gorm:"type:text"`
    RegisteredCanes RegisteredCanes `gorm:"type:jsonb;index:idx_attributes,type:gin"`
}