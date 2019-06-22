// GORMはGO言語用のORMフレームワークである。
// ORM(オブジェクト関係マッピング)とはデータベースとオブジェクト指向プログラミング言語の間の非互換なデータを変換するプログラミング技法である。


// importでGORMと使用するDBのドライバをインポート。
import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"time"
)

func main() {
	db, err := gorm.Open("postgres", "user=postgres password=postgres dbname=gorm sslmode=disable") //DBによって構文が違う
	if err != nil {
		panic(err.Error())
	}
	defer db.Close() // 実行完了後DB接続を閉じる

	seedDB(db)
	//*************************Retrieving Single Record******************************
	// 1レコードを取り出す
	u := User{}
	db.Debug().First(&u)										// usersテーブルから１のIDであるレコードを取り出す
	db.Debug().FirstOrInit(&u, &User{UserName:"fprefect"})		// user_name="fpreect"のレコードを取り出す
	db.Debug().FirstOrCreate(&u, &User{UserName:"lprosser"})	// usersテーブルにuser_nameが"lprosser"であるレコードがあれば取り出す。ない場合はDBに追加する
	db.Debug().Last(&u)

	//*************************Retrieving Record Sets********************************
	// 複数のデータを取り出す
	users := []User{}
	db.Debug().Find(&users)													// 全てのレコードを取り出す
	db.Debug().Find(&users, &User{UserName: "fprefect"})
	db.Debug().Find(&users, map[string]interface{}{"user_name":"fprefect"})
	db.Debug().Find(&users, "user_name = ?", "fprefect")
	for _, u := range users{
		fmt.Printf("\n%v\n", u)
	}

	//**************************Where Clauses****************************************
	// WHERE句とは、SQL文で検索条件を指定するためのものである
	// 条件によりレコードを取り出す
	users := []User {}
	db.Debug().Where("user_name = ?", "adent").Find(&users)
	db.Debug().Where(&User{UserName: "adent"}).Find(&users)
	db.Debug().Where(map[string]interface{}{"user_name":"adent"}).Find(&users)
	db.Debug().Where("user_name in (?)", []string{"adent", "tmacmillan"}).Find(&users)
	db.Debug().Where("user_name like ?", "%mac%" ).Find(&users)
	db.Debug().Where("user_name like ? and first_name = ? ", "%mac%", "Tricia" ).Find(&users)
	db.Debug().Where("created_at < ?", time.Now()).Find(&users)
	db.Debug().Where("created_at BETWEEN ? and ?", time.Now().Add(-30*24*time.Hour), time.Now()).Find(&users)
	db.Debug().Not("user_name = ?", "adent").Find(&users)
	db.Debug().Where("user_name = ?", "adent").Or("user_name = ?" , "fprefect").Find(&users)
	for _, u := range users {
		fmt.Printf("\n%v\n", u)
	}

	//***************************Preloading Child Objects******************************
	users := []User{}
	db.Debug().Preload("Calendar.Appointments").Find(&users)
	for _, u := range users {
		fmt.Printf("\n%v\n", u.Calendar)
	}

	//*****************************Limits Offsets Ordering******************************
	//LIMIT句を指定した場合は先頭のデータから指定した数のデータを取得しますが、先頭からではなく指定した位置からデータを取得することもできます。データの取得を行う最初の位置を指定するにはOFFSET句を使用します。
	users := []User{}
	db.Debug().Limit(2).Offset(2).Order("first_name ").Find(&users)
	for _, u := range users {
		fmt.Printf("\n%v\n", u)
	}

	//****************************Selecting data subsets*******************************
	users := []User{}
	db.Debug().Select([]string{"first_name", "last_name"}).Find(&users)
	db.Debug().Model(&User{}).Pluck("first_name", &usernames)
	userVMs := []UserViewModel{}
	db.Debug().Model(&User{}).Select([]string{"first_name", "last_name"}).Scan(&userVMs)
	for _, u := range userVMs {
		fmt.Printf("\n%v\n", u)
	}
	var count int
	db.Debug().Model(&User{}).Count(&count)
	fmt.Println(count)

	//********************Using Attrs and assign to provide default value *************
	u := User{}
	db.Debug().Where("user_name = ?", "adent").Attrs(&User{FirstName: "Eddie"}).FirstOrInit(&u)
	db.Debug().Where("user_name = ?", "adent").Assign(&User{FirstName: "Eddie"}).FirstOrInit(&u)
	fmt.Printf("\n%v\n", u)

	//***********************Creating Projections with Joins****************************
	usersVMS := []UserViewModel{}
	db.Debug().Model(&User{}).Joins("inner join calendars on calendars.user_id = users.id").Select("users.first_name, users.last_name, calendars.name").Scan(&usersVMS)
	
	for _, u := range usersVMS{
		fmt.Printf("\n%v\n", u)
	}

	//*********************Creating Aggregations with Group and Having*******************
	rows, _ := db.Debug().Model(&Appointment{}).Select("calendar_id, sum(length) ").Group("calendar_id").Rows()
	rows, _ := db.Debug().Model(&Appointment{}).Select("calendar_id, sum(length) ").Group("calendar_id").Having("calendar_id = ?", 1).Rows()
	for rows.Next() {
		var id, length int
		rows.Scan(&id, &length)
		fmt.Println(id, length)
	}

	//************************************Using Raw SQL**********************************
	users := []User{}
	//db.Find(&users)
	db.Debug().Exec("SELECT * FROM users").Find(&users)
	for _, u := range users{
		fmt.Printf("\n%v\n", u)
	}

	//************************************Scopes*****************************************
	appts := []Appointment{}
	db.Scopes(LongMeetings).Find(&appts)
	for _, appt := range appts {
		fmt.Printf("\n%v\n", appt)
	}

}

func LongMeetings(db *gorm.DB) *gorm.DB {
	return db.Where("length > ?", 60)
}

// DBにデータを追加する関数である。
func seedDB(db *gorm.DB){	

	db.DropTableIfExists(&User{})
	db.CreateTable(&User{})　			//usersテーブルは自動生成
	db.DropTableIfExists(&Calendar{})
	db.CreateTable(&Calendar{})			//calendarsテーブルは自動生成
	db.DropTableIfExists(&Appointment{})
	db.CreateTable(&Appointment{})		//appointmentsテーブルは自動生成

	users := map[string]*User{
		"adent":		&User{UserName:"adent", FirstName: "Arthur", LastName:"Dent"},
		"fprefect":		&User{UserName:"fprefect", FirstName: "Ford", LastName:"Prefect"},
		"tmacmillan":	&User{UserName:"tmacmillan", FirstName: "Tricia", LastName:"Macmillan"},
		"zbeeblebrox":	&User{UserName:"zbeeblebrox", FirstName: "Zaphod", LastName:"Beeblebrox"},
		"mrobot":		&User{UserName:"mrobot", FirstName: "Marvin", LastName:"Robot"},
	}
	for _, user := range users {
		user.Calendar = Calendar{Name: "Calendar"}
	}

	users["adent"].AddAppointment(&Appointment{
		Subject:	"Save House",
		StartTime:	parseTime("1979-07-02 08:00"),
		Length:		60,
	})

	users["fprefect"].AddAppointment(&Appointment{
		Subject:	"Get a Drink at Local Pub",
		StartTime:	parseTime("1979-07-02 10:00"),
		Length:		11,
		Attendees: 	[]*User{users["adent"]},
	})

	users["fprefect"].AddAppointment(&Appointment{
		Subject:	"Hitch a ride",
		StartTime:	parseTime("1979-07-02 10:12"),
		Length:		60,
		Attendees: 	[]*User{users["adent"]},
	})

	users["fprefect"].AddAppointment(&Appointment{
		Subject:	"Attend Poetry Reading",
		StartTime:	parseTime("1979-07-02 11:00"),
		Length:		30,
		Attendees: 	[]*User{users["adent"]},
	})

	users["fprefect"].AddAppointment(&Appointment{
		Subject:	"Get Thrown into Space",
		StartTime:	parseTime("1979-07-02 10:40"),
		Length:		5,
		Attendees: 	[]*User{users["adent"]},
	})

	users["fprefect"].AddAppointment(&Appointment{
		Subject:	"Get saved from Space",
		StartTime:	parseTime("1979-07-02 11:45"),
		Length:		1,
		Attendees: 	[]*User{users["adent"]},
	})

	users["zbeeblebrox"].AddAppointment(&Appointment{
		Subject:	"Explore Planet Builder's HomeWorld",
		StartTime:	parseTime("1979-07-03 11:00"),
		Length:		240,
		Attendees: 	[]*User{users["adent"]},
	})

	for _, user := range users {
		db.Save(&user)	//DBに以上のデータを保存する
	}
}


func parseTime(timeRaw string) time.Time {
	const timeLayout = "2006-01-02 15:04"
	t, _ := time.Parse(timeLayout, timeRaw)
	return t
}

// テーブル定義　
// 構造体でカラムを定義する時、変数名は先頭が大文字でなければ行けない
// usersというテーブルが自動的に生成される
// 七つのカラムが生成される - id(PK), created_at, updated_at, deleted_at, user_name, first_name, last_name
type User struct {
	gorm.Model
	UserName	string
	FirstName	string
	LastName	string
	Calendar	Calendar
}

func(u *User) AddAppointment(appt *Appointment) {
	u.Calendar.Appointments = append(u.Calendar.Appointments, appt)
}

// calendarsというテーブルが自動的に生成される
type Calendar struct {
	gorm.Model
	Name	string
	UserID 	uint
	Appointments []*Appointment
}

// appointmentsというテーブルが自動的に生成される
type Appointment struct {
	gorm.Model
	Subject		string
	Description string
	StartTime	time.Time
	Length		uint
	CalendarID	uint
	Attendees	[]*User `gorm:"many2many:appointment_user"` // many-to-manyの関係になるようなリレーションを示します
}

type UserViewModel struct {
	FirstName	string
	LastName	string
	CalendarName string `gorm:"column:name"`
}
