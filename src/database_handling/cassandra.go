package database_handling

import (
	"context"
	"io"
	"os"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
    "github.com/gocql/gocql"
)


func InitiateCassandra(m session) {

    
// create keyspaces if doesn't exist
err = session.Query("CREATE KEYSPACE IF NOT EXISTS private_cloud WITH REPLICATION = {'class' : 'NetworkTopologyStrategy', 'AWS_VPC_US_WEST_2' : 3};").Exec()
if err != nil {
    log.Println(err)
    return
}

// Create table to store client information (yet to be implemented)      if doesn't exist
err = session.Query("CREATE TABLE IF NOT EXISTS private_cloud.user_details (username text,ownedContainers text, PRIMARY KEY username);").Exec()
if err != nil {
    log.Println(err)
    return
}
// Create table to store system information (yet to be implemented)      if doesn't exist
err = session.Query("CREATE TABLE IF NOT EXISTS private_cloud.system_details (docker_images text, PRIMARY KEY username);").Exec()
if err != nil {
    log.Println(err)
    return
}

// insert some practice data
//err = session.Query("INSERT INTO sleep_centre.sleep_study (name, study_date, sleep_time_hours) VALUES ('James', '2018-01-07', 8.2);").Exec()
//err = session.Query("INSERT INTO sleep_centre.sleep_study (name, study_date, sleep_time_hours) VALUES ('James', '2018-01-08', 6.4);").Exec()
//err = session.Query("INSERT INTO sleep_centre.sleep_study (name, study_date, sleep_time_hours) VALUES ('James', '2018-01-09', 7.5);").Exec()
//err = session.Query("INSERT INTO sleep_centre.sleep_study (name, study_date, sleep_time_hours) VALUES ('Bob', '2018-01-07', 6.6);").Exec()
//err = session.Query("INSERT INTO sleep_centre.sleep_study (name, study_date, sleep_time_hours) VALUES ('Bob', '2018-01-08', 6.3);").Exec()
//err = session.Query("INSERT INTO sleep_centre.sleep_study (name, study_date, sleep_time_hours) VALUES ('Bob', '2018-01-09', 6.7);").Exec()
//err = session.Query("INSERT INTO sleep_centre.sleep_study (name, study_date, sleep_time_hours) VALUES ('Emily', '2018-01-07', 7.2);").Exec()
//err = session.Query("INSERT INTO sleep_centre.sleep_study (name, study_date, sleep_time_hours) VALUES ('Emily', '2018-01-09', 7.5);").Exec()
//if err != nil {
//    log.Println(err)
//    return
//}

// Return average sleep time for James
//var sleep_time_hours float32
//
//sleep_time_output := session.Query("SELECT avg(sleep_time_hours) FROM sleep_centre.sleep_study WHERE name = 'James';").Iter()
//sleep_time_output.Scan(&sleep_time_hours)
//fmt.Println("Average sleep time for James was: ", sleep_time_hours, "h")
//
//// return average sleep time for group
//sleep_time_output = session.Query("SELECT avg(sleep_time_hours) FROM sleep_centre.sleep_study;").Iter()
//sleep_time_output.Scan(&sleep_time_hours)
//fmt.Println("Average sleep time for the group was: ", sleep_time_hours, "h")
 
}


// To be called when authorized user requests storage of data
func StoreData(){




}