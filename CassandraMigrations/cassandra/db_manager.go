package cassandra

import (
	"github.com/gocql/gocql"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/cassandra"
	_ "github.com/mattes/migrate/source/file"
	"log"
	"strconv"
	"sync"
)

type DbManager struct {
	Session *gocql.Session
}

const (
	keyspaceName   = "bank"
	createkeyspace = "CREATE KEYSPACE IF NOT EXISTS bank WITH replication = {'class':'SimpleStrategy', 'replication_factor': 1}"
	port           = 9042
	//use local migration file address for accessing cql files
	migrationFiles = "file://C:/Go/bin/src/Alex_workshop/CassandraMigrations/migration/cql"
)

var (
	instance *DbManager
	once     sync.Once
)

func createInstance() *DbManager {
	instance = new(DbManager)
	host := "localhost"
	instance.handleMigrations(host)

	//create cassandra session
	cluster := gocql.NewCluster(host)
	cluster.Port = port
	cluster.Keyspace = keyspaceName
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Unable to connect to cassandra session: %v", err)
	}
	instance.Session = session
	log.Println("Cassandra session Initialized")

	return instance

}

func (db *DbManager) handleMigrations(host string) {
	//create cassandra keyspace
	createCassandraKeyspace(host)

	cassandra := cassandra.Cassandra{}

	url := "cassandra://" + host + ":" + strconv.Itoa(port) + "/" + keyspaceName
	driver, err := cassandra.Open(url)
	if err != nil {
		log.Fatalf("Unable to connect to cassandra:%v", err)
	}
	mig, err := migrate.NewWithDatabaseInstance(migrationFiles, "cassandra", driver)
	if err != nil {
		log.Fatalf("Failed to validate migration:%v", err)
	}
	if err := mig.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to apply migration %v", err)
	}
	//Close the connection and handle error
	if srcErr, dbErr := mig.Close(); srcErr != nil || dbErr != nil {
		log.Fatalf("Error while executing the migrations %v %v", srcErr, dbErr)
	}
}

func GetInstance() *DbManager {
	once.Do(func() {
		instance = createInstance()
	})
	return instance

}

func createCassandraKeyspace(host string) {
	//Create cassandra session
	cluster := gocql.NewCluster(host)
	cluster.Port = port
	dbSession, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Not able to connect to cassandra %v", err)
	}
	defer dbSession.Close()
	//Create keyspace if not exists
	if err := dbSession.Query(createkeyspace).Exec(); err != nil {
		//Close the cassandra session
		dbSession.Close()
		log.Fatalf("Failed to create keyspace %v", err)
	}
}

func (db *DbManager) close() {
	db.Session.Close()
	log.Println("Closing the cassandra session")
}
