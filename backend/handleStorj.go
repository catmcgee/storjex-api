package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/rs/xid"
	"storj.io/uplink"
)

const (
    apiKey = "1dfJhEcjZxSiQBvWWqm7n9PeDYDqB3iroEipqznUkiPdcn6fJsgf4bCV2EAn8QPXo2MgkDtUnmcdoKrGPmEJ4wm5T71kiiAazW4qoBeKuMdxb7eu3XaQ"
    satelliteAddress = "12L9ZFwhzVpuEKMUNUqkaTLGzwY9G24tbiigLiXpmZWKwmcNDDs@eu1.storj.io:7777"
)

func UploadData(ctx context.Context,
    accessGrant, bucketName, name string,
    data[] byte, numberOfDownloads int) (string, string, string) {
    // TODO: Create two accesses - admin and user, and get passphrases for both.
    // Call twice
    randomKey:= xid.New()
    objectKey:= randomKey.String()
    // Parse the admin Access Grant.
   project := ConnectToStorjexProject(accessGrant)
    
        // Intitiate the upload of our Object to the specified bucket and key.
    upload, err:= project.UploadObject(ctx, bucketName, objectKey, nil)
    if err != nil {
         fmt.Errorf("could not initiate upload: %v", err)
    }

    // Copy the data to the upload.
    buf:= bytes.NewBuffer(data)
    _, err = io.Copy(upload, buf)
    if err != nil {
        _ = upload.Abort()
         fmt.Errorf("could not upload data: %v", err)
    }

    // Commit the uploaded object.
    err = upload.Commit()
    if err != nil {
         fmt.Errorf("could not commit uploaded object: %v", err)
    }
    // Create passphrases & user access token
    adminAccessToken:= createAccessToken(name, objectKey, 0)
    userAccessToken:= createAccessToken(name, objectKey, 1)

    adminPassphrase, userPassphrase := generatePassphrases(
        context.Background(), 
        adminAccessToken, 
        userAccessToken,
        bucket, 
        name,
        objectKey,
        numberOfDownloads)

    return userPassphrase, adminPassphrase, objectKey

}
func DownloadData(passphrase string)(fileContents []byte, err error) {
        var accessGrant string
        var bucket string
        var key string
        var numberOfDownloads int
        ctx:= context.Background()

        conn := ConnectToDataBase()

        query:= fmt.Sprintf(`SELECT accessGrant, bucket, key, numberOfDownloads FROM passphrases WHERE passphrase = '%s'`,
            passphrase)
        defer conn.Close(ctx)
        err = conn.QueryRow(ctx, query).Scan( &accessGrant, &bucket, &key, &numberOfDownloads)
        if err != nil {
            fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
        }
        if (numberOfDownloads <= 0) {
            return nil, fmt.Errorf("This file has reached its maximum number of downloads")
        } else {

       project := ConnectToStorjexProject(myAccessGrant)
        download, err:= project.DownloadObject(ctx, bucket, key, nil)
        if err != nil {
             return nil, fmt.Errorf("Could not open object: %v", err)
        }
        defer download.Close()
        newNumberOfDownloads := numberOfDownloads - 1

        query= fmt.Sprintf(`UPDATE passphrases SET numberOfDownloads = %b WHERE passphrase = '%s'`,
            newNumberOfDownloads, passphrase)
        defer conn.Close(ctx)
        _, err = conn.Query(ctx, query)
        if err != nil {
           return nil, fmt.Errorf("Could not get file from database: %v", err)
        }

        // Read everything from the download stream
        receivedContents, err:= ioutil.ReadAll(download)
		if err != nil {
            return nil, fmt.Errorf("Could not read file data, may be corrupted: %v", err)
	   }
    
        // Check that the downloaded data is the same as the uploaded data.
        return receivedContents, nil
    }
}

func createAccessToken(name string, objectKey string, level int)(string) {
    // TODO: create two access tokens - one admin (can 
    // delete and update object) and one just for downloading. Call twice
    if (level == 0) {
    // TODO: set variables to put into admin permissions
    } else {}
    // set variables for user permission
    // TODO: THIS DOES NOT WORK
	access, err := uplink.RequestAccessWithPassphrase(context.Background(), 
	satelliteAddress, apiKey, "passphrase")
	if err != nil {
		fmt.Println(err)
	}
	
	// create an access grant for reading bucket "storjex"
	permission := uplink.ReadOnlyPermission()
	shared := uplink.SharePrefix{Bucket: "storjex"}
	restrictedAccess, err := access.Share(permission, shared)
	if err != nil {
		fmt.Println(err)
	}
	
	// serialize the restricted access grant
	serializedAccess, err := restrictedAccess.Serialize()
	if err != nil {
		fmt.Println(err)
	}

    return serializedAccess
}

func HandleDelete(passphrase string) (string) {
    var bucket string 
    var key string
    var accessGrant string
    ctx := context.Background()
    conn := ConnectToDataBase()

    query := fmt.Sprintf(`SELECT adminAccessGrant, bucket, key FROM passphrases WHERE adminPassphrase = '%s'`,
        passphrase)

    err := conn.QueryRow(ctx, query).Scan( &accessGrant, &bucket, &key)
    if err != nil {
        return fmt.Sprintf("Could not find file: %v\n", err)
    }

    project := ConnectToStorjexProject(myAccessGrant)

    _, err = project.DeleteObject(ctx, bucket, key) 
    if err != nil {
        return fmt.Sprintf("Could not delete file: %v", err) 
   }

   query = fmt.Sprintf(`DELETE FROM passphrases WHERE adminPassphrase = '%s'`,
        passphrase)

    _, err = conn.Query(ctx, query)
    if err != nil {
        return fmt.Sprintf("Could not delete from database: %v\n", err)
    }

return "Successfully deleted file"

}

func generatePassphrases(
    cxt context.Context,
    adminAccessToken string, 
    userAccessToken string,
    bucket string, 
    name string, 
    objectKey string,
    numberOfDownloads int)(string, string) {

    randomPass1:= xid.New()
    randomPass2:= xid.New()
    adminPassphrase:= "admin-" + name + "-" + randomPass1.String()
    userPassphrase:= name + "-" + randomPass2.String()
    // put the passes in database with access grant, bucket, key
    conn := ConnectToDataBase()

       
    query:= fmt.Sprintf(`INSERT INTO passphrases (passphrase, adminPassphrase, adminAccessGrant, accessGrant, bucket, key, numberOfDownloads) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%b')`,
            userPassphrase, adminPassphrase, adminAccessToken, userAccessToken, bucket, objectKey, numberOfDownloads)
        if _, err:= conn.Exec(context.Background(),
            query);
        err != nil {
            fmt.Println(err) }
		return adminPassphrase, userPassphrase
    
    }