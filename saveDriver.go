
























// MIGRATED TO SAVECASEDRIVER.GO

























package main

import (
	//"html/template"
	"net/http"
	"time"
	
	"strconv"
	
    "golang.org/x/net/context"
    
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	//"appengine/user"
	
	//"appengine/blobstore"
	"google.golang.org/appengine/blobstore" // https://cloud.google.com/appengine/docs/go/blobstore/reference
)

// ========== ========== ========== ========== ========== ========== ========== ========== ========== ==========
func saveDriver(r *http.Request, ctx context.Context) (string) { // context.Context vs appengine.Context
	output := ""
	
	// ========== ========== ========== ========== ==========
	blobkey := saveDriverBlobstore(r)
	output += "<h1>blobkey = ["+blobkey+"]</h1>"
	// ========== ========== ========== ========== ==========
	
	// ========== ========== ========== ========== ==========
	// Pull the POST form fields into a Golang struct
	/*
	type Driver struct {
		Name				string
		FrequencyResponse	string
		Width				int
		Price				int
		
		BlobKey				string
		
		DateAdded			time.Time
	}
	*/
	driverFrequencyLow, _	:= strconv.Atoi(r.FormValue("driverfrequencylow"))
	driverFrequencyHigh, _	:= strconv.Atoi(r.FormValue("driverfrequencyhigh"))
	
	driverdiameter, _		:= strconv.Atoi(r.FormValue("driverdiameter"))
	driverprice, _			:= strconv.Atoi(r.FormValue("driverprice"))
	
	driverData := Driver {
		Name:				r.FormValue("drivername"),
		
		//FrequencyResponse:	r.FormValue("driverfrequencyresponse"),
		FrequencyLow:		int32(driverFrequencyLow), // int32
		FrequencyHigh:		int32(driverFrequencyHigh), // int32
		
		Diameter:			int16(driverdiameter), // int
		Price:				int32(driverprice), // int
		
		BlobKey:			blobkey,
		
		DateAdded:			time.Now(),
	}
	// ========== ========== ========== ========== ==========
	
	
	// ========== ========== ========== ========== ==========
	output += "<h1>r.FormValue(\"drivername\") = ["+r.FormValue("drivername")+"]</h1>"
	output += "<h1>driverData.Name = ["+driverData.Name+"]</h1>"
	if driverData.Name != "" {
		output += saveDriverDatastore(r, ctx, driverData)
	} else if driverData.BlobKey != "" {
		// ========== ========== ========== ========== ==========
		// Delete blobstore entry
		//deleteBlobKey := "EUG76sbkgL8CDsNUokKcRQ=="
		output += "<div>Prepare to delete from blobstore: "+driverData.BlobKey+"</div>"
		blobstore.Delete(ctx, appengine.BlobKey(driverData.BlobKey)) // https://cloud.google.com/appengine/docs/go/blobstore/reference#Delete
		output += "<div>Finished deleting from blobstore</div>"
		output += `<div><a href="/">Return Home</a></div>`
		// ========== ========== ========== ========== ==========
	} else {
		output += "<h1>Form was submitted blank</h1>"
	}
	// ========== ========== ========== ========== ==========
	
	
    return output
}
// ========== ========== ========== ========== ========== ========== ========== ========== ========== ==========

// ========== ========== ========== ========== ========== ========== ========== ========== ========== ==========
func saveDriverBlobstore(r *http.Request) (string) {
	output := ""
	
	// ========== ========== ========== ========== ==========
	// Store the image in the blobstore
	blobs, _, err := blobstore.ParseUpload(r)
	if err != nil {
		//output += "<h1>ERROR: "+err.Error()+"</h1>"
	}
	file := blobs["file"]
	
	if len(file) == 0 {
		//output += "<h1>WARNING: No image file uploaded to blobstore</h1>"
		output = ""
	} else {
		output = string(file[0].BlobKey)
	}
	// ========== ========== ========== ========== ==========
	
    return output
}
// ========== ========== ========== ========== ========== ========== ========== ========== ========== ==========

// ========== ========== ========== ========== ========== ========== ========== ========== ========== ==========
func saveDriverDatastore(r *http.Request, ctx context.Context, driverData Driver) (string) { // context.Context vs appengine.Context
	output := ""
	
	// ========== ========== ========== ========== ==========
	// Store Golang struct in the datastore
	key := datastore.NewIncompleteKey(ctx, "Driver", driverKey(ctx))
	_, err := datastore.Put(ctx, key, &driverData)
	if err != nil {
		output += "<h1>ERROR: datastore failed</h1>"
	} else {
		output += "<h1>SUCCESS: Created datastore entry for new driver</h1>"
		if driverData.BlobKey != "" {
			output += "<img src=\"/serve/?blobKey="+driverData.BlobKey+"\" />"
		} else {
			output += "<h1>No image was uploaded</h1>"
		}
	}
	// ========== ========== ========== ========== ==========
	
    return output
}
// ========== ========== ========== ========== ========== ========== ========== ========== ========== ==========





