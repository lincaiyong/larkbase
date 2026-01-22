package tos

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"github.com/volcengine/ve-tos-golang-sdk/v2/tos"
	"io"
	"os"
)

func md5Of(b []byte) string {
	sum := md5.Sum(b)
	return hex.EncodeToString(sum[:])
}

var ak, sk, region, endpoint, bucket string

func init() {
	ak = os.Getenv("VOLC_TOS_AK")
	sk = os.Getenv("VOLC_TOS_SK")
	region = os.Getenv("VOLC_TOS_REGION")
	endpoint = os.Getenv("VOLC_TOS_ENDPOINT")
	bucket = os.Getenv("VOLC_TOS_BUCKET")
}

func Put(ctx context.Context, bs []byte) (string, error) {
	hash := md5Of(bs)
	return hash, put(ctx, hash, bytes.NewReader(bs))
}

func put(ctx context.Context, key string, content io.Reader) error {
	client, err := tos.NewClientV2(endpoint, tos.WithRegion(region), tos.WithCredentials(tos.NewStaticCredentials(ak, sk)))
	if err != nil {
		return err
	}
	output, err := client.ListObjectsType2(ctx, &tos.ListObjectsType2Input{
		Bucket:  bucket,
		MaxKeys: 1,
		Prefix:  key,
	})
	if err != nil {
		return err
	}
	if len(output.Contents) > 0 {
		return nil
	}
	_, err = client.PutObjectV2(ctx, &tos.PutObjectV2Input{
		PutObjectBasicInput: tos.PutObjectBasicInput{
			Bucket: bucket,
			Key:    key,
		},
		Content: content,
	})
	return err
}

func Get(ctx context.Context, key string) ([]byte, error) {
	client, err := tos.NewClientV2(endpoint, tos.WithRegion(region), tos.WithCredentials(tos.NewStaticCredentials(ak, sk)))
	if err != nil {
		return nil, err
	}
	output, err := client.GetObjectV2(ctx, &tos.GetObjectV2Input{
		Bucket: bucket,
		Key:    key,
	})
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = output.Content.Close()
	}()
	data, err := io.ReadAll(output.Content)
	if err != nil {
		return nil, err
	}
	return data, nil
}
