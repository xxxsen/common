package s3

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/xxxsen/common/errs"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var Client *S3Client

type S3Client struct {
	c      *config
	sess   *session.Session
	client *s3.S3
}

func toBase64MD5CheckSum(val string) *string {
	raw, err := hex.DecodeString(val)
	if err != nil {
		log.Printf("invalid md5 checksum:%s, err:%v", val, err)
		return aws.String("invalid")
	}
	return aws.String(base64.StdEncoding.EncodeToString(raw))
}

func InitGlobal(opts ...Option) error {
	client, err := New(opts...)
	if err != nil {
		return errs.Wrap(errs.ErrS3, "init s3", err)
	}
	Client = client
	return nil
}

func New(opts ...Option) (*S3Client, error) {
	c := &config{
		ssl:    true,
		region: "cn",
	}
	for _, opt := range opts {
		opt(c)
	}
	if len(c.bucket) == 0 {
		return nil, errs.New(errs.ErrParam, "nil bucket name")
	}

	credit := credentials.NewStaticCredentials(c.secretId, c.secretKey, "")
	sess, err := session.NewSession(&aws.Config{
		Credentials:      credit,
		Endpoint:         aws.String(c.endpoint),
		DisableSSL:       aws.Bool(!c.ssl),
		HTTPClient:       &http.Client{},
		Region:           aws.String(c.region),
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		return nil, errs.Wrap(errs.ErrS3, "init session fail", err)
	}
	client := s3.New(sess)
	return &S3Client{c: c, client: client, sess: sess}, nil
}

func (c *S3Client) DownloadByRange(ctx context.Context, fileid string, at int64) (io.ReadCloser, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(c.c.bucket),
		Key:    aws.String(fileid),
	}
	if at != 0 {
		input.Range = aws.String(fmt.Sprintf("%d-", at))
	}
	output, err := c.client.GetObject(input)
	if err != nil {
		return nil, errs.Wrap(errs.ErrS3, "get obj fail", err)
	}
	return output.Body, nil
}

func (c *S3Client) Download(ctx context.Context, fileid string) (io.ReadCloser, error) {
	return c.DownloadByRange(ctx, fileid, 0)
}

func (c *S3Client) Upload(ctx context.Context, fileid string, r io.ReadSeeker, sz int64, cks ...string) (string, error) {
	input := &s3.PutObjectInput{
		Body:   r,
		Bucket: aws.String(c.c.bucket),
		Key:    aws.String(fileid),
	}
	if len(cks) > 0 && len(cks[0]) > 0 {
		input.ContentMD5 = toBase64MD5CheckSum(cks[0])
	}
	rsp, err := c.client.PutObject(input)
	if err != nil {
		return "", errs.Wrap(errs.ErrS3, "write obj fail", err)
	}
	return c.unquote(*rsp.ETag), nil
}

func (c *S3Client) Remove(ctx context.Context, fileid string) error {
	_, err := c.client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(c.c.bucket),
		Key:    aws.String(fileid),
	})
	if err != nil {
		return errs.Wrap(errs.ErrS3, "delete fail", err)
	}
	return nil
}

func (c *S3Client) BeginUpload(ctx context.Context, fileid string) (string, error) {
	output, err := c.client.CreateMultipartUpload(&s3.CreateMultipartUploadInput{
		Bucket: aws.String(c.c.bucket),
		Key:    aws.String(fileid),
	})
	if err != nil {
		return "", errs.Wrap(errs.ErrS3, "create multi part upload fail", err)
	}
	return *output.UploadId, nil
}

func (c *S3Client) UploadPart(ctx context.Context, fileid string, uploadid string, partid int, file io.ReadSeeker, cks ...string) error {
	input := &s3.UploadPartInput{
		Body:       file,
		Bucket:     aws.String(c.c.bucket),
		Key:        aws.String(fileid),
		PartNumber: aws.Int64(int64(partid)),
		UploadId:   aws.String(uploadid),
	}
	if len(cks) > 0 {
		input.ContentMD5 = toBase64MD5CheckSum(cks[0])
	}
	_, err := c.client.UploadPart(input)
	if err != nil {
		return errs.Wrap(errs.ErrS3, "put part fail", err)
	}
	return nil
}

func (c *S3Client) listParts(ctx context.Context, fileid string, uploadid string) ([]*s3.Part, error) {
	output, err := c.client.ListParts(&s3.ListPartsInput{
		Bucket:              aws.String(c.c.bucket),
		ExpectedBucketOwner: new(string),
		Key:                 aws.String(fileid),
		UploadId:            aws.String(uploadid),
	})
	if err != nil {
		return nil, errs.Wrap(errs.ErrS3, "list part fail", err)
	}
	return output.Parts, nil
}

func (c *S3Client) parts2completeparts(src []*s3.Part) []*s3.CompletedPart {
	out := make([]*s3.CompletedPart, 0, len(src))
	for _, p := range src {
		out = append(out, &s3.CompletedPart{
			ChecksumCRC32:  p.ChecksumCRC32,
			ChecksumCRC32C: p.ChecksumCRC32C,
			ChecksumSHA1:   p.ChecksumSHA1,
			ChecksumSHA256: p.ChecksumSHA256,
			ETag:           p.ETag,
			PartNumber:     p.PartNumber,
		})
	}
	return out
}

func (c *S3Client) EndUpload(ctx context.Context, fileid string, uploadid string, partcount int) (string, error) {
	parts, err := c.listParts(ctx, fileid, uploadid)
	if err != nil {
		return "", err
	}
	if len(parts) != partcount {
		return "", errs.New(errs.ErrParam, "part count not match, need:%d, get:%d", partcount, len(parts))
	}
	rsp, err := c.client.CompleteMultipartUpload(&s3.CompleteMultipartUploadInput{
		Bucket: aws.String(c.c.bucket),
		Key:    aws.String(fileid),
		MultipartUpload: &s3.CompletedMultipartUpload{
			Parts: c.parts2completeparts(parts),
		},
		UploadId: aws.String(uploadid),
	})
	if err != nil {
		return "", errs.Wrap(errs.ErrS3, "finish upload fail", err)
	}
	return c.unquote(*rsp.ETag), nil
}

func (c *S3Client) DiscardMultiPartUpload(ctx context.Context, fileid string, uploadid string) error {
	_, err := c.client.AbortMultipartUpload(&s3.AbortMultipartUploadInput{
		Bucket:   aws.String(c.c.bucket),
		Key:      aws.String(fileid),
		UploadId: aws.String(uploadid),
	})
	if err != nil {
		return errs.Wrap(errs.ErrS3, "abort multipart upload fail", err)
	}
	return nil
}

type ObjectMetaInfo struct {
	ETag *string
}

func (c *S3Client) GetFileInfo(ctx context.Context, fileid string) (*ObjectMetaInfo, error) {
	out, err := c.client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(c.c.bucket),
		Key:    aws.String(fileid),
	})
	if err != nil {
		return nil, errs.Wrap(errs.ErrS3, "get obj info from s3 fail", err)
	}
	return &ObjectMetaInfo{
		ETag: aws.String(c.unquote(*out.ETag)),
	}, nil
}

func (c *S3Client) unquote(etag string) string {
	return strings.Trim(etag, "\"")
}
