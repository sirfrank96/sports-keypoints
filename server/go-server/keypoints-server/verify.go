package keypointsserver

import (
	"fmt"
	"net/mail"

	skp "github.com/sirfrank96/go-server/sports-keypoints-proto"
)

func verifyCreateUserRequest(request *skp.CreateUserRequest) error {
	if request == nil {
		return fmt.Errorf("request is empty")
	}
	if request.UserName == "" {
		return fmt.Errorf("please enter a non-empty username")
	}
	if request.Password == "" {
		return fmt.Errorf("please enter a non-empty password")
	}
	if request.Email == "" {
		return fmt.Errorf("please enter a non-empty email")
	}
	// make sure email is correct format
	emailParsed, err := mail.ParseAddress(request.Email)
	if err != nil {
		return fmt.Errorf("invalid email: %s", err.Error())
	}
	if emailParsed.Address != request.Email {
		return fmt.Errorf("invalid email, bad format")
	}
	return nil
}

func verifyRegisterUserRequest(request *skp.RegisterUserRequest) error {
	if request == nil {
		return fmt.Errorf("request is empty")
	}
	if request.UserName == "" {
		return fmt.Errorf("please enter a non-empty username")
	}
	if request.Password == "" {
		return fmt.Errorf("please enter a non-empty password")
	}
	return nil
}

func verifyReadUserRequest(request *skp.ReadUserRequest) error {
	if request == nil {
		return fmt.Errorf("request is empty")
	}
	return nil
}

func verifyUpdateUserRequest(request *skp.UpdateUserRequest) error {
	if request == nil {
		return fmt.Errorf("request is empty")
	}
	if request.UserName == "" && request.Password == "" && request.Email == "" {
		return fmt.Errorf("please add at least one field to be updated")
	}
	return nil
}

func verifyDeleteUserRequest(request *skp.DeleteUserRequest) error {
	if request == nil {
		return fmt.Errorf("request is empty")
	}
	return nil
}

func verifyUploadInputImageRequest(request *skp.UploadInputImageRequest) error {
	if request == nil {
		return fmt.Errorf("request is empty")
	}
	if request.ImageType == skp.ImageType_IMAGE_TYPE_UNSPECIFIED {
		return fmt.Errorf("please enter an image type")
	}
	if request.Image == nil || len(request.Image) == 0 {
		return fmt.Errorf("please upload an input image")
	}
	if request.Description == "" {
		return fmt.Errorf("please add a description")
	}
	if request.Timestamp == nil {
		return fmt.Errorf("please add a timestamp")
	}
	if err := request.Timestamp.CheckValid(); err != nil {
		return fmt.Errorf("invalid timestamp: %s", err.Error())
	}
	return nil
}

func verifyListInputImagesForUserRequest(request *skp.ListInputImagesForUserRequest) error {
	if request == nil {
		return fmt.Errorf("request is empty")
	}
	return nil
}

func verifyReadInputImageRequest(request *skp.ReadInputImageRequest) error {
	if request == nil {
		return fmt.Errorf("request is empty")
	}
	if request.InputImageId == "" {
		return fmt.Errorf("please enter an input image id")
	}
	return nil
}

func verifyDeleteInputImageRequest(request *skp.DeleteInputImageRequest) error {
	if request == nil {
		return fmt.Errorf("request is empty")
	}
	if request.InputImageId == "" {
		return fmt.Errorf("please enter an input image id")
	}
	return nil
}

func verifyCalibrateInputImageRequest(request *skp.CalibrateInputImageRequest) error {
	if request == nil {
		return fmt.Errorf("request is empty")
	}
	if request.InputImageId == "" {
		return fmt.Errorf("please enter an input image id")
	}
	return nil
}

func verifyCalculateGolfKeypointsRequest(request *skp.CalculateGolfKeypointsRequest) error {
	if request == nil {
		return fmt.Errorf("request is empty")
	}
	if request.InputImageId == "" {
		return fmt.Errorf("please enter an input image id")
	}
	return nil
}

func verifyReadGolfKeypointsRequest(request *skp.ReadGolfKeypointsRequest) error {
	if request == nil {
		return fmt.Errorf("request is empty")
	}
	if request.InputImageId == "" {
		return fmt.Errorf("please enter an input image id")
	}
	return nil
}

func verifyUpdateBodyKeypointsRequest(request *skp.UpdateBodyKeypointsRequest) error {
	if request == nil {
		return fmt.Errorf("request is empty")
	}
	if request.InputImageId == "" {
		return fmt.Errorf("please enter an input image id")
	}
	if request.UpdatedBodyKeypoints == nil {
		return fmt.Errorf("please add at least one body keypoint to update")
	}
	return nil
}

func verifyDeleteGolfKeypointsRequest(request *skp.DeleteGolfKeypointsRequest) error {
	if request == nil {
		return fmt.Errorf("request is empty")
	}
	if request.InputImageId == "" {
		return fmt.Errorf("please enter an input image id")
	}
	return nil
}
