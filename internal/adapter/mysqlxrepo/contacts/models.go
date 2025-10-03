package contacts

import (
	"database/sql"
	"fmt"
	"labra/internal/entity"
	"strconv"
	"strings"
	"time"
)

type LinkedEntity struct {
	ContactID  int          `db:"contact_id"`
	EntityID   string       `db:"entity_id"`
	EntityType string       `db:"entity_type"`
	VerifiedAt sql.NullTime `db:"verified_at"`
	CreatedAt  time.Time    `db:"created_at"`
	UpdatedAt  time.Time    `db:"updated_at"`
}

type LinkedEntities []LinkedEntity

type Contacts []Contact

type Contact struct {
	ID               int64          `db:"id"`
	Type             string         `db:"type"`
	Value            string         `db:"value"`
	CreatedAt        time.Time      `db:"created_at"`
	UpdatedAt        time.Time      `db:"updated_at"`
	DeletedAt        sql.NullTime   `db:"deleted_at"`
	LinkedEntityID   sql.NullInt64  `db:"linked_entity_id"`
	LinkedEntityType sql.NullString `db:"linked_entity_type"`
}

func NewFromEntity(contact entity.Contact) Contact {
	return Contact{
		ID:    int64(contact.ID),
		Type:  string(contact.Type),
		Value: contact.Value,
	}
}

func (l LinkedEntity) buildEntity() entity.LinkedEntity {
	entityID, _ := strconv.Atoi(l.EntityID)

	return entity.LinkedEntity{
		ID:         entityID,
		ContactID:  l.ContactID,
		Type:       l.EntityType,
		IsVerified: l.VerifiedAt.Valid,
	}
}

func (l LinkedEntities) buildEntity() []entity.LinkedEntity {
	res := make([]entity.LinkedEntity, 0, len(l))

	for _, link := range l {
		res = append(res, link.buildEntity())
	}

	return res
}

func (c Contact) buildEntity() entity.Contact {
	return entity.Contact{
		ID:    int(c.ID),
		Type:  entity.ContactType(c.Type),
		Value: c.Value,
	}
}

func (c Contacts) buildEntity() entity.Contacts {
	res := make(entity.Contacts, 0, len(c))

	linkedEntMap := c.MapByEntity()

	for entKey, cont := range linkedEntMap {
		eID, eType := parseKey(entKey)
		eUID, _ := strconv.Atoi(eID)
		entCont := cont.buildEntity()
		entCont.LinkedEntities = append(entCont.LinkedEntities, entity.LinkedEntity{
			ID:        eUID,
			Type:      eType,
			ContactID: entCont.ID,
		})

		res = append(res, entCont)
	}

	return res
}

func (c Contacts) MapByEntity() map[string]Contact {
	m := map[string]Contact{}
	for _, cont := range c {
		if cont.LinkedEntityType.Valid {
			m[cont.EntityKey()] = cont
		}
	}

	return m
}

func (c Contact) EntityKey() string {
	return fmt.Sprintf("%d_%s", c.LinkedEntityID.Int64, c.LinkedEntityType.String)
}

func parseKey(entKey string) (entID string, entType string) {
	parts := strings.Split(entKey, "_")
	if len(parts) != 2 {
		return "", ""
	}

	return parts[0], parts[1]
}
