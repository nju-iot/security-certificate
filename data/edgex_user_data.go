package dal

import (
	"time"

	"github.com/nju-iot/security-certificate/caller"
	"github.com/nju-iot/security-certificate/logs"
	"gorm.io/gorm"
)

// EdgexUser ...
type EdgexUser struct {
	ID           int64     `gorm:"column:id" json:"id"`
	Username     string    `gorm:"column:username" json:"username"`
	Password     string    `gorm:"column:password" json:"password"`
	PhoneNumber  string    `gorm:"column:phone_number" json:"phone_number"`
	Email        string    `gorm:"column:email" json:"email"`
	Deleted      int32     `gorm:"column:deleted" json:"deleted"`
	CreatedTime  time.Time `gorm:"column:created_time" json:"created_time"`
	ModifiedTime time.Time `gorm:"column:modified_time" json:"modified_time"`
}

// AddEdgexUser ...
func AddEdgexUser(db *gorm.DB, user *EdgexUser) (err error) {
	dbRes := db.Debug().Model(&EdgexUser{}).Create(user)
	if dbRes.Error != nil {
		err = dbRes.Error
		logs.Error("[AddEdgexUser] add user failed: user=%v, err=%v", user, err)
		return
	}
	return
}

// GetEdgexUserByName ...
func GetEdgexUserByName(name string) (user *EdgexUser, err error) {
	userList := make([]*EdgexUser, 0)
	dbRes := caller.EdgexDB.Debug().Model(&EdgexUser{}).Where("username = ?", name).Find(&userList)
	if dbRes.Error != nil {
		err = dbRes.Error
		logs.Error("[GetEdgexUserByName] get edgex user failed: name=%v, err=%v", name, err)
		return
	}
	if len(userList) > 0 {
		user = userList[0]
	}
	return
}

// GetEdgexUserByID ...
func GetEdgexUserByID(id int64) (user *EdgexUser, err error) {
	userList := make([]*EdgexUser, 0)
	dbRes := caller.EdgexDB.Debug().Model(&EdgexUser{}).Where("id = ?", id).Find(&userList)
	if dbRes.Error != nil {
		err = dbRes.Error
		logs.Error("[GetEdgexUserByID] get edgex user failed: userID=%v, err=%v", id, err)
		return
	}
	if len(userList) > 0 {
		user = userList[0]
	}
	return
}

// GetEdgexUserByEmail ...
func GetEdgexUserByNameAndEmail(username string, email string) (user *EdgexUser, err error) {
	userList := make([]*EdgexUser, 0)
	dbRes := caller.EdgexDB.Debug().Model(&EdgexUser{}).Where("username = ? or email = ?", username, email).Find(&userList)
	if dbRes.Error != nil {
		err = dbRes.Error
		logs.Error("[GetEdgexUserByEmail] get edgex user failed: username=%v, email=%v, err=%v", username, email, err)
		return
	}
	if len(userList) > 0 {
		user = userList[0]
	}
	return
}

// UpdateEdgexUser ...
func UpdateEdgexUser(userID int64, fieldsMap map[string]interface{}) error {
	dbRes := caller.EdgexDB.Debug().Model(&EdgexUser{}).Where("id = ?", userID).Updates(fieldsMap)
	if dbRes.Error != nil {
		logs.Error("[UpdateEdgexRelatedUser] update EdgexRelatedUser failed: id=%+v, filedsMap=%+v, err=%v", userID, fieldsMap, dbRes.Error)
		return dbRes.Error
	}
	return nil
}

// MGetEdgexUserMapByIDs ...
func MGetEdgexUserMapByIDs(ids []int64) (userMap map[int64]*EdgexUser, err error) {
	userMap = make(map[int64]*EdgexUser)
	if len(ids) == 0 {
		return
	}
	userList := make([]*EdgexUser, 0)
	dbRes := caller.EdgexDB.Debug().Model(&EdgexUser{}).Where("id IN (?)", ids).Find(&userList)
	if dbRes.Error != nil {
		err = dbRes.Error
		logs.Error("[GetEdgexUserByID] get edgex user failed: userIDs=%v, err=%v", ids, err)
		return
	}
	for _, user := range userList {
		userMap[user.ID] = user
	}
	return
}
