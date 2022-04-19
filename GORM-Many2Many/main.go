package main

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Roles struct {
	gorm.Model
	Name      string     `gorm:"size:61;not null"`
	Processes []*Process `gorm:"many2many:process_roles"`
}

type Process struct {
	gorm.Model
	Name string   `gorm:"size:50;not null;"`
	Role []*Roles `gorm:"many2many:process_roles"`
}

func Migrate(db *gorm.DB) {
	rolesPrototype := &Roles{}
	processPrototype := &Process{}
	db.AutoMigrate(rolesPrototype, processPrototype)

}

func Seed(db *gorm.DB) {
	create := &Process{
		Name: "Create",
	}
	read := &Process{
		Name: "Read",
	}
	update := &Process{
		Name: "Update",
	}
	delete := &Process{
		Name: "Delete",
	}
	db.Save(create)
	db.Save(read)
	db.Save(update)
	db.Save(delete)

	Admin := &Roles{
		Name: "Admin",
		Processes: []*Process{
			create, read, update, delete,
		},
	}
	Editor := &Roles{
		Name: "Editor",
		Processes: []*Process{
			read, update,
		},
	}
	Viewer := &Roles{
		Name: "Viewer",
		Processes: []*Process{
			read,
		},
	}

	db.Save(&Admin)
	db.Save(&Editor)
	db.Save(&Viewer)
	fmt.Printf("Roles  created:\n%v\n", *Admin)
	fmt.Printf("Roles  created:\n%v\n", *Editor)
	fmt.Printf("Roles  created:\n%v\n", *Viewer)
	fmt.Printf("Processes created:\n%v\n", []*Process{create, read, update, delete})

}
func ListEverything(db *gorm.DB) {
	roles := []Roles{}
	db.Preload("Process").Preload("Process.Roles").Find(&roles)

	for _, role := range roles {
		fmt.Printf("Role data: %v\n", &role)
		for _, role := range role.Processes {

			fmt.Printf("Role-Process data: %v\n", &role)
		}
	}
}

func ClearEverything(db *gorm.DB) {
	err1 := db.Delete(&Roles{}).Error
	err2 := db.Delete(&Process{}).Error
	fmt.Printf("Deleting the records:\n%v\n", err1)
	fmt.Printf("Deleting the records:\n%v\n", err2)
}
func FindAssociation(db *gorm.DB) {
	var roles []map[string]interface{}
	db.Table("roles").Find(&roles)
	var processes []map[string]interface{}
	db.Table("processes").Find(&processes)

	fmt.Printf("Querying the Role records:\n%v\n", roles)
	fmt.Printf("Querying the Process records:\n%v\n", processes)
}

func main() {
	dsn := "host=18.185.93.196 user=postgres password=postgres dbname=testErtugrul port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	Migrate(db)
	Seed(db)
	ListEverything(db)
	//result := map[string]interface{}{}
	FindAssociation(db)
	//ClearEverything(db)

}
