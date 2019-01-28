package main

import (
	"bytes"
	b64 "encoding/base64"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ascii "github.com/cantasaurus/goasciiart"
)

var (
	headersSuccess = map[string]string{
		"Access-Control-Allow-Origin": "*",
		"Content-Type":                "image/png",
	}
	headersError = map[string]string{
		"Access-Control-Allow-Origin": "*",
	}
	acceptedContentTypes = [2]string{"image/png", "image/jpeg"}
)

func acceptableImageType(req events.APIGatewayProxyRequest) bool {
	for _, elem := range acceptedContentTypes {
		if elem == req.Headers["Content-Type"] {
			return true
		}
	}
	return false
}

func base64DecodeImage(b64Image string) (image.Image, error) {
	imageData, err := b64.StdEncoding.DecodeString(b64Image)
	r := bytes.NewReader(imageData)
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func base64EncodeImage(imageData []byte) string {
	b64String := b64.StdEncoding.EncodeToString(imageData)
	return b64String
}

func clientError(status int, errorMessage string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       errorMessage,
		Headers:    headersError,
	}, nil
}

func PostHandler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if !acceptableImageType(req) {
		return clientError(http.StatusBadRequest, "Content-Type was not image/png or image/jpeg.")
	}

	img, err := base64DecodeImage(req.Body)
	if err != nil {
		return clientError(http.StatusUnprocessableEntity, "Could not base64 decode image.")
	}

	asciiBytes := ascii.Convert2Ascii(ascii.ScaleImage(img, 120))
	rgba, err := ascii.TextToImage(string(asciiBytes))
	if err != nil {
		return clientError(http.StatusUnprocessableEntity, "Unable to convert ascii text to image.")
	}

	imageData, err := ascii.ConvertToImage(rgba)
	b64EncodedImg := base64EncodeImage(imageData)
	return events.APIGatewayProxyResponse{
		StatusCode:      http.StatusOK,
		Body:            b64EncodedImg,
		Headers:         headersSuccess,
		IsBase64Encoded: true,
	}, nil
}

func main() {
	lambda.Start(PostHandler)
}
