package keypointsserver

import (
	"testing"

	skp "github.com/sirfrank96/go-server/sports-keypoints-proto"
)

func TestVerifyCreateUserRequest(t *testing.T) {
	// nil request
	err := verifyCreateUserRequest(nil)
	if err == nil {
		t.Errorf("verifyCreateUserRequest(nil) is supposed to have an error")
	}
	// empty request
	createUserRequest := &skp.CreateUserRequest{}
	err = verifyCreateUserRequest(createUserRequest)
	if err == nil {
		t.Errorf("verifyCreateUserRequest(%+v) is supposed to have an error", createUserRequest)
	}
	// only user name set
	createUserRequest.UserName = "user1"
	err = verifyCreateUserRequest(createUserRequest)
	if err == nil {
		t.Errorf("verifyCreateUserRequest(%+v) is supposed to have an error", createUserRequest)
	}
	// only user name and password set
	createUserRequest.Password = "password1"
	err = verifyCreateUserRequest(createUserRequest)
	if err == nil {
		t.Errorf("verifyCreateUserRequest(%+v) is supposed to have an error", createUserRequest)
	}
	// bad email
	createUserRequest.Email = "abc"
	err = verifyCreateUserRequest(createUserRequest)
	if err == nil {
		t.Errorf("verifyCreateUserRequest(%+v) is supposed to have an error", createUserRequest)
	}
	// bad email
	createUserRequest.Email = "abc.gmail.com"
	err = verifyCreateUserRequest(createUserRequest)
	if err == nil {
		t.Errorf("verifyCreateUserRequest(%+v) is supposed to have an error", createUserRequest)
	}
	// good request
	createUserRequest.Email = "abc@gmail.com"
	err = verifyCreateUserRequest(createUserRequest)
	if err != nil {
		t.Errorf("verifyCreateUserRequest(%+v) had an unexpected error: %s", createUserRequest, err.Error())
	}
}

func TestVerifyRegisterUserRequest(t *testing.T) {
	// nil request
	err := verifyRegisterUserRequest(nil)
	if err == nil {
		t.Errorf("(verifyRegisterUserRequest(nil) is supposed to have an error")
	}
	// empty request
	registerUserRequest := &skp.RegisterUserRequest{}
	err = verifyRegisterUserRequest(registerUserRequest)
	if err == nil {
		t.Errorf("(verifyRegisterUserRequest(%+v) is supposed to have an error", registerUserRequest)
	}
	// only user name set
	registerUserRequest.UserName = "user1"
	err = verifyRegisterUserRequest(registerUserRequest)
	if err == nil {
		t.Errorf("(verifyRegisterUserRequest(%+v) is supposed to have an error", registerUserRequest)
	}
	// good request
	registerUserRequest.Password = "password1"
	err = verifyRegisterUserRequest(registerUserRequest)
	if err != nil {
		t.Errorf("(verifyRegisterUserRequest(%+v) had an unexpected error: %s", registerUserRequest, err.Error())
	}
}

func TestVerifyReadUserRequest(t *testing.T) {
	// nil request
	err := verifyReadUserRequest(nil)
	if err == nil {
		t.Errorf("(verifyReadUserRequest(nil) is supposed to have an error")
	}
	// good request
	readUserRequest := &skp.ReadUserRequest{}
	err = verifyReadUserRequest(readUserRequest)
	if err != nil {
		t.Errorf("(verifyReadUserRequest(%+v) had an unexpected error: %s", readUserRequest, err.Error())
	}
}

func TestVerifyUpdateUserRequest(t *testing.T) {
	// nil request
	err := verifyUpdateUserRequest(nil)
	if err == nil {
		t.Errorf("(verifyUpdateUserRequest(nil) is supposed to have an error")
	}
	// empty request
	updateUserRequest := &skp.UpdateUserRequest{}
	err = verifyUpdateUserRequest(updateUserRequest)
	if err == nil {
		t.Errorf("(verifyUpdateUserRequest(%+v) is supposed to have an error", updateUserRequest)
	}
	// good request
	updateUserRequest.Password = "password1"
	err = verifyUpdateUserRequest(updateUserRequest)
	if err != nil {
		t.Errorf("verifyUpdateUserRequest(%+v) had an unexpected error: %s", updateUserRequest, err.Error())
	}
}

func TestVerifyDeleteUserRequest(t *testing.T) {
	// nil request
	err := verifyDeleteUserRequest(nil)
	if err == nil {
		t.Errorf("(verifyDeleteUserRequest(nil) is supposed to have an error")
	}
	// good request
	deleteUserRequest := &skp.DeleteUserRequest{}
	err = verifyDeleteUserRequest(deleteUserRequest)
	if err != nil {
		t.Errorf("(verifyDeleteUserRequest(%+v) had an unexpected error: %s", deleteUserRequest, err.Error())
	}
}

func TestVerifyUploadInputImageRequest(t *testing.T) {
	// nil request
	err := verifyUploadInputImageRequest(nil)
	if err == nil {
		t.Errorf("(verifyUploadInputImageRequest(nil) is supposed to have an error")
	}
	// empty request
	uploadInputImageRequest := &skp.UploadInputImageRequest{}
	err = verifyUploadInputImageRequest(uploadInputImageRequest)
	if err == nil {
		t.Errorf("(verifyUploadInputImageRequest(%+v) is supposed to have an error", uploadInputImageRequest)
	}
	// only image type set
	uploadInputImageRequest.ImageType = skp.ImageType_DTL
	err = verifyUploadInputImageRequest(uploadInputImageRequest)
	if err == nil {
		t.Errorf("(verifyUploadInputImageRequest(%+v) is supposed to have an error", uploadInputImageRequest)
	}
	// zero length image
	uploadInputImageRequest.Image = []byte{}
	err = verifyUploadInputImageRequest(uploadInputImageRequest)
	if err == nil {
		t.Errorf("(verifyUploadInputImageRequest(%+v) is supposed to have an error", uploadInputImageRequest)
	}
	// good request
	uploadInputImageRequest.Image = []byte{1, 2, 3, 5, 8}
	err = verifyUploadInputImageRequest(uploadInputImageRequest)
	if err != nil {
		t.Errorf("verifyUploadInputImageRequest(%+v) had an unexpected error: %s", uploadInputImageRequest, err.Error())
	}
}

func TestVerifyListInputImagesForUserRequest(t *testing.T) {
	// nil request
	err := verifyListInputImagesForUserRequest(nil)
	if err == nil {
		t.Errorf("(verifyListInputImagesForUserRequest(nil) is supposed to have an error")
	}
	// good request
	listInputImagesForUserRequest := &skp.ListInputImagesForUserRequest{}
	err = verifyListInputImagesForUserRequest(listInputImagesForUserRequest)
	if err != nil {
		t.Errorf("(verifyListInputImagesForUserRequest(%+v) had an unexpected error: %s", listInputImagesForUserRequest, err.Error())
	}
}

func TestVerifyReadInputImageRequest(t *testing.T) {
	// nil request
	err := verifyReadInputImageRequest(nil)
	if err == nil {
		t.Errorf("(verifyReadInputImageRequest(nil) is supposed to have an error")
	}
	// empty request
	readInputImageRequest := &skp.ReadInputImageRequest{}
	err = verifyReadInputImageRequest(readInputImageRequest)
	if err == nil {
		t.Errorf("(verifyReadInputImageRequest(%+v) is supposed to have an error", readInputImageRequest)
	}
	// good request
	readInputImageRequest.InputImageId = "image1"
	err = verifyReadInputImageRequest(readInputImageRequest)
	if err != nil {
		t.Errorf("verifyReadInputImageRequest(%+v) had an unexpected error: %s", readInputImageRequest, err.Error())
	}
}

func TestVerifyDeleteInputImageRequest(t *testing.T) {
	// nil request
	err := verifyDeleteInputImageRequest(nil)
	if err == nil {
		t.Errorf("(verifyDeleteInputImageRequest(nil) is supposed to have an error")
	}
	// empty request
	deleteInputImageRequest := &skp.DeleteInputImageRequest{}
	err = verifyDeleteInputImageRequest(deleteInputImageRequest)
	if err == nil {
		t.Errorf("(verifyDeleteInputImageRequest(%+v) is supposed to have an error", deleteInputImageRequest)
	}
	// good request
	deleteInputImageRequest.InputImageId = "image1"
	err = verifyDeleteInputImageRequest(deleteInputImageRequest)
	if err != nil {
		t.Errorf("verifyDeleteInputImageRequest(%+v) had an unexpected error: %s", deleteInputImageRequest, err.Error())
	}
}

func TestVerifyCalibrateInputImageRequest(t *testing.T) {
	// nil request
	err := verifyCalibrateInputImageRequest(nil)
	if err == nil {
		t.Errorf("(verifyCalibrateInputImageRequest(nil) is supposed to have an error")
	}
	// empty request
	calibrateInputImageRequest := &skp.CalibrateInputImageRequest{}
	err = verifyCalibrateInputImageRequest(calibrateInputImageRequest)
	if err == nil {
		t.Errorf("(verifyCalibrateInputImageRequest(%+v) is supposed to have an error", calibrateInputImageRequest)
	}
	// good request
	calibrateInputImageRequest.InputImageId = "image1"
	err = verifyCalibrateInputImageRequest(calibrateInputImageRequest)
	if err != nil {
		t.Errorf("verifyCalibrateInputImageRequest(%+v) had an unexpected error: %s", calibrateInputImageRequest, err.Error())
	}
}

func TestVerifyCalculateGolfKeypointsRequest(t *testing.T) {
	// nil request
	err := verifyCalculateGolfKeypointsRequest(nil)
	if err == nil {
		t.Errorf("(verifyCalculateGolfKeypointsRequest(nil) is supposed to have an error")
	}
	// empty request
	calculateGolfKeypointsRequest := &skp.CalculateGolfKeypointsRequest{}
	err = verifyCalculateGolfKeypointsRequest(calculateGolfKeypointsRequest)
	if err == nil {
		t.Errorf("(verifyCalculateGolfKeypointsRequest(%+v) is supposed to have an error", calculateGolfKeypointsRequest)
	}
	// good request
	calculateGolfKeypointsRequest.InputImageId = "image1"
	err = verifyCalculateGolfKeypointsRequest(calculateGolfKeypointsRequest)
	if err != nil {
		t.Errorf("verifyCalculateGolfKeypoints(%+v) had an unexpected error: %s", calculateGolfKeypointsRequest, err.Error())
	}
}

func TestVerifyReadGolfKeypointsRequest(t *testing.T) {
	// nil request
	err := verifyReadGolfKeypointsRequest(nil)
	if err == nil {
		t.Errorf("(verifyReadGolfKeypointsRequest(nil) is supposed to have an error")
	}
	// empty request
	readGolfKeypointsRequest := &skp.ReadGolfKeypointsRequest{}
	err = verifyReadGolfKeypointsRequest(readGolfKeypointsRequest)
	if err == nil {
		t.Errorf("(verifyReadGolfKeypointsRequest(%+v) is supposed to have an error", readGolfKeypointsRequest)
	}
	// good request
	readGolfKeypointsRequest.InputImageId = "image1"
	err = verifyReadGolfKeypointsRequest(readGolfKeypointsRequest)
	if err != nil {
		t.Errorf("verifyReadGolfKeypoints(%+v) had an unexpected error: %s", readGolfKeypointsRequest, err.Error())
	}
}

func TestVerifyUpdateBodyKeypointsRequest(t *testing.T) {
	// nil request
	err := verifyUpdateBodyKeypointsRequest(nil)
	if err == nil {
		t.Errorf("(verifyUpdateBodyKeypointsRequest(nil) is supposed to have an error")
	}
	// empty request
	updateBodyKeypointsRequest := &skp.UpdateBodyKeypointsRequest{}
	err = verifyUpdateBodyKeypointsRequest(updateBodyKeypointsRequest)
	if err == nil {
		t.Errorf("(verifyUpdateBodyKeypointsRequest(%+v) is supposed to have an error", updateBodyKeypointsRequest)
	}
	// only input image set
	updateBodyKeypointsRequest.InputImageId = "image1"
	err = verifyUpdateBodyKeypointsRequest(updateBodyKeypointsRequest)
	if err == nil {
		t.Errorf("(verifyUpdateBodyKeypointsRequest(%+v) is supposed to have an error", updateBodyKeypointsRequest)
	}
	// good request
	updateBodyKeypointsRequest.UpdatedBodyKeypoints = &skp.Body25PoseKeypoints{
		Nose: &skp.Keypoint{
			X:          1.23,
			Y:          4.56,
			Confidence: 1.0,
		},
	}
	err = verifyUpdateBodyKeypointsRequest(updateBodyKeypointsRequest)
	if err != nil {
		t.Errorf("verifyUpdateBodyKeypoints(%+v) had an unexpected error: %s", updateBodyKeypointsRequest, err.Error())
	}
}

func TestVerifyDeleteGolfKeypointsRequest(t *testing.T) {
	// nil request
	err := verifyDeleteGolfKeypointsRequest(nil)
	if err == nil {
		t.Errorf("(verifyDeleteGolfKeypointsRequest(nil) is supposed to have an error")
	}
	// empty request
	deleteGolfKeypointsRequest := &skp.DeleteGolfKeypointsRequest{}
	err = verifyDeleteGolfKeypointsRequest(deleteGolfKeypointsRequest)
	if err == nil {
		t.Errorf("(verifyDeleteGolfKeypointsRequest(%+v) is supposed to have an error", deleteGolfKeypointsRequest)
	}
	// good request
	deleteGolfKeypointsRequest.InputImageId = "image1"
	err = verifyDeleteGolfKeypointsRequest(deleteGolfKeypointsRequest)
	if err != nil {
		t.Errorf("verifyDeleteGolfKeypoints(%+v) had an unexpected error: %s", deleteGolfKeypointsRequest, err.Error())
	}
}
