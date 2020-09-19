package main

import (
	"os"
	"fmt"
	"log"
	"flag"
	"regexp"
	"net/http"
	// "encoding/json"
	"path/filepath"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
)

type Mapping struct {
	gorm.Model
	Loadbalancer 	string `json:"loadbalancer"`
	Bigip 			string `json:"bigip"`
}

type Program struct {
	DB *gorm.DB
	// DBPath string
	// Port string
	UUIDRegex *regexp.Regexp
}

var prog Program

func main() {
	// fmt.Println("golang + gin + gorm")
	progPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	dbDefault := filepath.Join(progPath, "mapping.db")
	var port string
	var dbpath string
	flag.StringVar(&port, "port", "8080", "The port to listen.")
	flag.StringVar(&dbpath, "dbpath", dbDefault, "The database path.")
	flag.Parse()
	
	db, err := gorm.Open(sqlite.Open(dbpath), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&Mapping{})

	reg, _ := regexp.Compile("[a-f0-9]{8}-([a-f0-9]{4}-){3}[a-f0-9]{12}")

	prog = Program{
		// DBPath: dbpath,
		DB: db,
		// Port: port,
		UUIDRegex: reg,
	}

	r := gin.Default()

	r.GET("/", Hello)
	r.GET("/mapping", GET)
	r.POST("/mapping", POST)
	r.DELETE("/mapping", DELETE)

	r.Run(fmt.Sprintf(":%s", port))
}

func Hello(c *gin.Context) { 
	c.String(
		200, 
		"Here, Find the Mapping Between LoadBalancer and BIG-IP.",
	)
}

func GET(c *gin.Context) {
	lb := c.DefaultQuery("loadbalancer", "")
	bigip := c.DefaultQuery("bigip", "")

	var m []Mapping
	var err error

	if lb == "" && bigip == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing query 'loadbalancer' or 'bigip'",
			"result": []map[string]interface{}{},
		})
		return
	} else if lb != "" && bigip == "" {
		m, err = GetMappingOfLoadbalancer(lb)
	} else if lb == "" && bigip != "" {
		m, err = GetMappingOfBigip(bigip)
	} else {
		m, err = GetMapping(lb, bigip)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
			"result": []map[string]interface{}{},
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"error": "",
			"result": m,
		})
		return
	}
}

func POST(c *gin.Context) {
	var m Mapping
	err := c.BindJSON(&m)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	} else {
		err = PostMapping(m.Loadbalancer, m.Bigip)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(http.StatusCreated, gin.H{
				"error": "",
			})
		}
		return
	}
}

func DELETE(c *gin.Context) {
	lb := c.DefaultQuery("loadbalancer", "")
	bigip := c.DefaultQuery("bigip", "")

	if lb == "" || bigip == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Query 'loadbalancer' and 'bigip' cannot be none.",
		})
		return
	}

	err := DeleteMapping(lb, bigip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusAccepted, gin.H{
			"error": "",
		})
	}
}

func GetMappingOfLoadbalancer(lb string) ([]Mapping, error) {
	if !ValidateUUID(lb) {
		return []Mapping{}, fmt.Errorf("Invalid UUID: %s", lb)
	}
	m := []Mapping{}
	if err := prog.DB.Where("loadbalancer = ?", lb).Find(&m).Error; err != nil {
		return []Mapping{}, err
	} else {
		return m, nil
	}
}

func GetMappingOfBigip(bigip string) ([]Mapping, error) {
	if !ValidateUUID(bigip) {
		return []Mapping{}, fmt.Errorf("Invalid UUID: %s", bigip)
	}
	m := []Mapping{}
	if err := prog.DB.Where("bigip = ?", bigip).Find(&m).Error; err != nil {
		return []Mapping{}, err
	} else {
		return m, nil
	}
}

func GetMapping(lb string, bigip string) ([]Mapping, error) {
	if !ValidateUUID(bigip) {
		return []Mapping{}, fmt.Errorf("Invalid UUID: %s", bigip)
	}
	if !ValidateUUID(lb) {
		return []Mapping{}, fmt.Errorf("Invalid UUID: %s", lb)
	}

	m := []Mapping{}
	if err := prog.DB.Where(
			"loadbalancer = ? AND bigip = ?", lb, bigip,
		).Find(&m).Error; err != nil {
		return []Mapping{}, err
	} else {
		return m, nil
	}
}

func PostMapping(lb string, bigip string) error {
	if !ValidateUUID(bigip) {
		return fmt.Errorf("Invalid UUID: %s", bigip)
	}
	if !ValidateUUID(lb) {
		return fmt.Errorf("Invalid UUID: %s", lb)
	}

	m := Mapping{Loadbalancer: lb, Bigip: bigip}
	rlt := prog.DB.Create(&m)
	log.Printf("Saving %v, affected rows: %d\n", m, rlt.RowsAffected)
	return rlt.Error
}

func DeleteMapping(lb string, bigip string) error {
	if !ValidateUUID(bigip) {
		return fmt.Errorf("Invalid UUID: %s", bigip)
	}
	if !ValidateUUID(lb) {
		return fmt.Errorf("Invalid UUID: %s", lb)
	}

	m := Mapping{Loadbalancer: lb, Bigip: bigip}
	rlt := prog.DB.Where("loadbalancer = ? AND bigip = ?", lb, bigip).Delete(&m)
	log.Printf("Deleting %v, affected rows: %d\n", m , rlt.RowsAffected)
	return rlt.Error
}

func ValidateUUID(uuid string) bool {
	return prog.UUIDRegex.Match([]byte(uuid))
}
