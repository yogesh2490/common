package couchbase

import (
	"github.com/couchbase/gocb"
	"utils/common"
	"utils/models"
)

var Cluster = &gocb.Cluster{}

type ICouch interface {
	Init(couchbase *models.CouchbaseDetails) (*gocb.Bucket, error, common.IResult)
}

type Couch struct {
	config common.IConfigGetter
}

//GetCouch returns an instance of an ICouch interface implementation for use.
func GetCouch(
	configs common.IConfigGetter,
) ICouch {
	return &Couch{
		config: configs,
	}
}

/*
	Init makes a connection to a Couchbase cluster, returns a bucket connection,
	and stores the connection in memory for future use.
*/
func (this Couch) Init(CouchbaseDetails *models.CouchbaseDetails) (Couchbase *gocb.Bucket, err error, result common.IResult) {
	result = common.MakeAPIResult(this.config)

	for _, host := range CouchbaseDetails.Nodes {
		connection_string := "couchbase://" + host
		result.Infof("Trying connection to couchbase at: %s\n", connection_string)

		Cluster, err = gocb.Connect(connection_string)
		if err != nil {
			result.Errorf("Not able to connect to Couchbase at %s because %s. Trying another node...\n", host, err.Error())
			continue
		} else {
			result.Infof("Opening couchbase bucket: %s with password: %s", CouchbaseDetails.Bucket, CouchbaseDetails.BucketPassword)
			Cluster.Authenticate(gocb.PasswordAuthenticator{
				this.config.MustGetConfigVar("COUCHBASE_USERNAME"),
				this.config.MustGetConfigVar("COUCHBASE_PASSWORD")})
			Couchbase, err = Cluster.OpenBucket(this.config.MustGetConfigVar("COUCHBASE_BUCKET"), this.config.MustGetConfigVar("BUCKET_PASSWORD"))
			if err != nil {
				result.Errorf("Not able to open Couchbase bucket at %s because: %s. Trying another node...\n", host, err.Error())
			} else {
				break
			}
		}
	}

	return
}
