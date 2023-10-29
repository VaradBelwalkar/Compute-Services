package external
import (
    "net/smtp"
    "path/filepath"
    "os"
)

func AttachFilesAndFolders(msg *Message, files []string, folders []string) error {
    for _, file := range files {
        attachment, err := os.Open(file)
        if err != nil {
            return err
        }
        defer attachment.Close()

        fileInfo, err := attachment.Stat()
        if err != nil {
            return err
        }

        header := make(map[string]string)
        header["Content-Type"] = "application/octet-stream"
        header["Content-Disposition"] = "attachment; filename=" + fileInfo.Name()
        header["Content-Transfer-Encoding"] = "base64"

        msg.Attach(header, attachment)
    }

    for _, folder := range folders {
        err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
            if err != nil {
                return err
            }

            if info.IsDir() {
                return nil
            }

            attachment, err := os.Open(path)
            if err != nil {
                return err
            }
            defer attachment.Close()

            header := make(map[string]string)
            header["Content-Type"] = "application/octet-stream"
            header["Content-Disposition"] = "attachment; filename=" + info.Name()
            header["Content-Transfer-Encoding"] = "base64"

            msg.Attach(header, attachment)

            return nil
        })

        if err != nil {
            return err
        }
    }

    return nil
}

func SendMail(username string, eMail string, files []string, folders []string) int {
    // ...
    if eMail == ""{   
        return http.StatusBadRequest
     }
         toEmailAddress := eMail //Should be dynamic generated format
         to := []string{toEmailAddress}
     
         host := "smtp.gmail.com"
         port := "587"
         address := host + ":" + port 
         msg := []byte("To: "+toEmailAddress+"\r\n" +
         "Subject: Your OTP for Two-Factor Authentication\n\n\r\n" +
         "\r\n" +
         "Dear User,\n\nYour OTP for two-factor authentication is: " + "\n\nPlease enter this otp in your app to complete the authentication process.\n\nBest regards,\nDYPLUG\r\n")

    err := AttachFilesAndFolders(msg, files, folders)
    if err != nil {
        return http.StatusInternalServerError
    }

    auth := smtp.PlainAuth("", twofa.Official_Email, twofa.Official_Email_Password, host)

    err = smtp.SendMail(address, auth, twofa.Official_Email, to, msg)
    if err != nil {
        return http.StatusInternalServerError
    }
return http.StatusOK

}

//files := []string{"/path/to/file1.txt", "/path/to/file2.pdf"}
//folders := []string{"/path/to/folder1", "/path/to/folder2"}

//SendMail("username", "email@example.com", files, folders)
