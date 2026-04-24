package logic

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"dropshipbe/dropshipbe"
	"dropshipbe/internal/svc"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/zeromicro/go-zero/core/logx"
)

type UploadFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUploadFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadFileLogic {
	return &UploadFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// --- Files ---
func (l *UploadFileLogic) UploadFile(in *dropshipbe.UploadFileRequest) (*dropshipbe.UploadFileResponse, error) {
	// todo: add your logic here and delete this line
	numFiles := len(in.Files)
	if numFiles == 0 {
		return &dropshipbe.UploadFileResponse{}, nil
	}
	files := make([]*dropshipbe.UploadedFileInfo, numFiles)
	errChan := make(chan error, numFiles)

	var wg sync.WaitGroup

	expirationDuration := time.Duration(l.svcCtx.Config.R2.LinkExpiration) * time.Minute

	for i, file := range in.Files {

		wg.Add(1)
		go func(index int, f *dropshipbe.FileData) {
			defer wg.Done()

			contentType := http.DetectContentType(f.Content)
			fileID := fmt.Sprintf("%d_%s", time.Now().UnixNano(), f.Filename)
			putInput := &s3.PutObjectInput{
				Bucket:      aws.String(l.svcCtx.Config.R2.BucketName),
				Key:         aws.String(fileID),
				Body:        bytes.NewReader(f.Content),
				ContentType: aws.String(contentType),
			}

			uploadCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()
			_, err := l.svcCtx.S3Client.PutObject(uploadCtx, putInput)
			if err != nil {
				l.Logger.Errorf("Lỗi khi tải file %s lên R2: %v", f.Filename, err)
				errChan <- fmt.Errorf("không thể tải file %s", f.Filename)
				return
			}

			presignedReq, err := l.svcCtx.PresignClient.PresignGetObject(l.ctx, &s3.GetObjectInput{
				Bucket: aws.String(l.svcCtx.Config.R2.BucketName),
				Key:    aws.String(fileID),
			}, s3.WithPresignExpires(expirationDuration))

			if err != nil {
				l.Logger.Errorf("Error %s: %v", f.Filename, err)
				errChan <- fmt.Errorf("Error can not create link %s", f.Filename)
				return
			}
			files[index] = &dropshipbe.UploadedFileInfo{
				FileId: fileID,
				Url:    presignedReq.URL,
			}
		}(i, file)

	}

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return nil, <-errChan
	}

	return &dropshipbe.UploadFileResponse{
		Files: files,
	}, nil
}
