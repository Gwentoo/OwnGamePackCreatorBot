package googleDrive

func DeleteFromGoogleDrive(fileID string) error {
	return DriveService.Files.Delete(fileID).Do()
}
