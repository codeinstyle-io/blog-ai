package cli

import (
	"fmt"
	"log"

	"codeinstyle.io/captain/db"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

func getValidInput(prompt string, validator func(string) error) string {
	for {
		var input string
		fmt.Print(prompt)
		fmt.Scanln(&input)

		if err := validator(input); err != nil {
			fmt.Printf("Error: %v. Please try again.\n", err)
			continue
		}
		return input
	}
}

func getValidPassword(prompt string) string {
	fmt.Println("\nPassword requirements:")
	fmt.Println("- At least 8 characters long")
	fmt.Println("- At least one uppercase letter")
	fmt.Println("- At least one lowercase letter")
	fmt.Println("- At least one number")
	fmt.Println("- At least one special character (!@#$%^&*(),.?\":{}|<>)")
	fmt.Println()

	for {
		fmt.Print(prompt)
		passwordBytes, _ := term.ReadPassword(0)
		password := string(passwordBytes)
		fmt.Println() // Add newline after password input

		if err := validatePassword(password); err != nil {
			fmt.Printf("Error: %v. Please try again.\n", err)
			continue
		}
		return password
	}
}

func CreateUser(cmd *cobra.Command, args []string) {
	firstName := getValidInput("First Name: ", validateFirstName)
	lastName := getValidInput("Last Name: ", validateLastName)
	email := getValidInput("Email: ", validateEmail)
	password := getValidPassword("Password: ")

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	database := db.InitDB()
	user := db.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  string(hashedPassword),
	}

	if err := database.Create(&user).Error; err != nil {
		log.Printf("Failed to create user: %v\n", err)
		return
	}

	fmt.Println("User created successfully")
}

func UpdateUserPassword(cmd *cobra.Command, args []string) {
	email := getValidInput("Email: ", validateEmail)

	database := db.InitDB()
	var user db.User
	if err := database.Where("email = ?", email).First(&user).Error; err != nil {
		fmt.Println("User not found")
		return
	}

	for {
		fmt.Print("Old Password: ")
		oldPasswordBytes, _ := term.ReadPassword(0)
		fmt.Println()

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), oldPasswordBytes); err != nil {
			fmt.Println("Incorrect password. Please try again.")
			continue
		}
		break
	}

	newPassword := getValidPassword("New Password: ")

	for {
		fmt.Print("Confirm Password: ")
		confirmBytes, _ := term.ReadPassword(0)
		fmt.Println()
		confirmPassword := string(confirmBytes)

		if newPassword != confirmPassword {
			fmt.Println("Passwords don't match. Please try again.")
			newPassword = getValidPassword("New Password: ")
			continue
		}
		break
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	if err := database.Save(&user).Error; err != nil {
		log.Printf("Failed to update password: %v\n", err)
		return
	}

	fmt.Println("Password updated successfully")
}
