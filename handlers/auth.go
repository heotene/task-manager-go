package handlers

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/mail"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"Myjob/database"
	"Myjob/middleware"
	"Myjob/models"
	"Myjob/utils"

	"golang.org/x/crypto/bcrypt"
)

// When form is submitted
func Register(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		name := strings.Join(
			strings.Fields(r.FormValue("fullname")),
			" ",
		)
		email := strings.ToLower(
			strings.TrimSpace(r.FormValue("email")),
		)
		phone := strings.TrimSpace(r.FormValue("phone"))

		phone = strings.ReplaceAll(phone, " ", "")
		phone = strings.ReplaceAll(phone, "-", "")
		phone = strings.ReplaceAll(phone, "(", "")
		phone = strings.ReplaceAll(phone, ")", "")
		if strings.HasPrefix(phone, "00") {
			phone = "+" + strings.TrimPrefix(phone, "00")
		}
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm_password")
		verificationMethod := r.FormValue("verification_method")

		// Check password length
		if len(password) < 8 {

			utils.Render(w, "register.html", models.PageData{
				Title: "Register",
				Error: "Password must be at least 8 characters long.",
			})

			return
		}

		hasNumber := false

		for _, char := range password {
			if unicode.IsDigit(char) {
				hasNumber = true
				break
			}
		}

		if !hasNumber {
			utils.Render(w, "register.html", models.PageData{
				Title: "Register",
				Error: "Password must contain at least one number.",
			})
			return
		}

		hasUppercase := false

		for _, char := range password {
			if unicode.IsUpper(char) {
				hasUppercase = true
				break
			}
		}

		if !hasUppercase {
			utils.Render(w, "register.html", models.PageData{
				Title: "Register",
				Error: "Password must contain at least one uppercase letter.",
			})
			return
		}

		hasSpecial := false

		for _, char := range password {
			if unicode.IsPunct(char) || unicode.IsSymbol(char) {
				hasSpecial = true
				break
			}
		}

		if !hasSpecial {
			utils.Render(w, "register.html", models.PageData{
				Title: "Register",
				Error: "Password must contain at least one special character.",
			})
			return
		}

		// Check passwords
		if password != confirmPassword {

			utils.Render(w, "register.html", models.PageData{
				Title: "Register",
				Error: "Passwords do not match.",
			})

			return
		}

		// Check verification method
		if verificationMethod != "email" &&
			verificationMethod != "phone" {

			utils.Render(w, "register.html", models.PageData{
				Title: "Register",
				Error: "Please select a verification method.",
			})

			return
		}

		// Email verification
		if verificationMethod == "email" && email == "" {

			utils.Render(w, "register.html", models.PageData{
				Title: "Register",
				Error: "Please enter your email address.",
			})

			return
		}
		// Validate email
		if verificationMethod == "email" {

			_, err := mail.ParseAddress(email)

			if err != nil {

				utils.Render(w, "register.html", models.PageData{
					Title: "Register",
					Error: "Please enter a valid email address.",
				})

				return
			}
		}
		if verificationMethod == "email" {

			_, err := database.GetUserByEmail(email)

			if err == nil {

				utils.Render(w, "register.html", models.PageData{
					Title: "Register",
					Error: "An account with this email already exists.",
				})

				return
			}
		}

		// Phone verification
		if verificationMethod == "phone" && phone == "" {

			utils.Render(w, "register.html", models.PageData{
				Title: "Register",
				Error: "Please enter your phone number.",
			})

			return
		}

		// Validate phone number
		if verificationMethod == "phone" {

			phonePattern := `^\+?[0-9]{10,15}$`

			matched, err := regexp.MatchString(phonePattern, phone)

			if err != nil {
				http.Error(
					w,
					"Unable to validate phone number.",
					http.StatusInternalServerError,
				)
				return
			}

			if !matched {

				utils.Render(w, "register.html", models.PageData{
					Title: "Register",
					Error: "Please enter a valid phone number.",
				})

				return
			}
		}
		if verificationMethod == "phone" {

			_, err := database.GetUserByPhone(phone)

			if err == nil {

				utils.Render(w, "register.html", models.PageData{
					Title: "Register",
					Error: "An account with this phone number already exists.",
				})

				return
			}
		}
		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword(
			[]byte(password),
			bcrypt.DefaultCost,
		)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Generate verification code
		verificationCode := fmt.Sprintf("%06d", rand.Intn(1000000))
		verificationExpires := time.Now().Add(10 * time.Minute)

		user := models.User{

			FullName: name,

			Email: email,

			Phone: phone,

			Password: string(hashedPassword),

			VerificationMethod: verificationMethod,

			EmailVerified: false,

			PhoneVerified: false,

			VerificationCode: verificationCode,

			VerificationExpires: verificationExpires.Format("2006-01-02 15:04:05"),
		}

		userID, err := database.CreateUser(user)

		if err != nil {
			utils.Render(w, "register.html", models.PageData{
				Title: "Register",
				Error: "Unable to create account. Email or phone may already be registered.",
			})

			return
		}

		user.ID = userID

		// Send verification code by email
		if verificationMethod == "email" {

			err = utils.SendVerificationCode(
				user.Email,
				verificationCode,
			)

			if err != nil {

				log.Println("VERIFICATION EMAIL ERROR:", err)

				utils.Render(w, "register.html", models.PageData{
					Title: "Register",
					Error: "Account was created, but we could not send the verification code. Please try again.",
				})

				return
			}
		}

		http.Redirect(
			w,
			r,
			"/verify?id="+strconv.Itoa(user.ID),
			http.StatusSeeOther,
		)

		return
	}

	utils.Render(w, "register.html", models.PageData{
		Title: "Register",
	})
}

func Verify(w http.ResponseWriter, r *http.Request) {

	userIDString := r.URL.Query().Get("id")

	userID, err := strconv.Atoi(userIDString)

	if err != nil || userID <= 0 {

		utils.Render(w, "verify.html", models.PageData{
			Title: "Verify Account",
			Error: "Invalid user ID.",
		})

		return
	}

	// Get user
	user, err := database.GetUserByID(userID)

	if err != nil {

		utils.Render(w, "verify.html", models.PageData{
			Title: "Verify Account",
			Error: "User account could not be found.",
		})

		return
	}

	// Show verification page
	if r.Method == "GET" {

		utils.Render(w, "verify.html", models.PageData{
			Title: "Verify Account",
			User:  user,
		})

		return
	}

	// Get submitted code
	code := strings.TrimSpace(r.FormValue("code"))

	if code == "" {

		utils.Render(w, "verify.html", models.PageData{
			Title: "Verify Account",
			User:  user,
			Error: "Please enter your verification code.",
		})

		return
	}

	// Check if verification code has expired
	if user.VerificationExpires == "" {

		utils.Render(w, "verify.html", models.PageData{
			Title: "Verify Account",
			User:  user,
			Error: "Your verification code is invalid or has expired.",
		})

		return
	}

	expires, err := time.Parse(
		"2006-01-02 15:04:05",
		user.VerificationExpires,
	)

	if err != nil || time.Now().After(expires) {

		utils.Render(w, "verify.html", models.PageData{
			Title: "Verify Account",
			User:  user,
			Error: "Your verification code has expired. Please request a new code.",
		})

		return
	}

	// Check verification code
	if user.VerificationCode != code {

		utils.Render(w, "verify.html", models.PageData{
			Title: "Verify Account",
			User:  user,
			Error: "Incorrect verification code.",
		})

		return
	}

	// Verify account
	err = database.VerifyUser(
		userID,
		user.VerificationMethod,
	)

	if err != nil {

		log.Println("VERIFY USER ERROR:", err)

		utils.Render(w, "verify.html", models.PageData{
			Title: "Verify Account",
			User:  user,
			Error: "Unable to verify your account.",
		})

		return
	}

	// Clear verification code after successful verification
	err = database.ClearVerification(userID)

	if err != nil {

		log.Println("CLEAR VERIFICATION ERROR:", err)

		utils.Render(w, "verify.html", models.PageData{
			Title: "Verify Account",
			User:  user,
			Error: "Account was verified, but the verification code could not be cleared.",
		})

		return
	}

	// Account verified successfully
	http.Redirect(
		w,
		r,
		"/login",
		http.StatusSeeOther,
	)
}

func Login(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		login := r.FormValue("login")
		password := r.FormValue("password")

		var user models.User
		var err error

		// Try email first
		user, err = database.GetUserByEmail(login)

		log.Println("LOGIN VALUE:", login)
		log.Println("EMAIL LOOKUP ERROR:", err)

		if err != nil {
			user, err = database.GetUserByPhone(login)

			log.Println("PHONE LOOKUP ERROR:", err)
		}

		if err == nil {
			log.Println("USER FOUND ID:", user.ID)
			log.Println("USER EMAIL:", user.Email)
			log.Println("USER PHONE:", user.Phone)
		}

		// User not found
		if err != nil {
			utils.Render(w, "login.html", models.PageData{
				Title: "Login",
				Error: "Incorrect email/phone or password.",
			})
			return
		}

		// Check password
		err = bcrypt.CompareHashAndPassword(
			[]byte(user.Password),
			[]byte(password),
		)

		if err != nil {
			utils.Render(w, "login.html", models.PageData{
				Title: "Login",
				Error: "Incorrect email/phone or password.",
			})
			return
		}

		// Check verification
		if user.VerificationMethod == "email" &&
			!user.EmailVerified {

			utils.Render(w, "login.html", models.PageData{
				Title: "Login",
				Error: "Please verify your email before logging in.",
			})
			return
		}

		if user.VerificationMethod == "phone" &&
			!user.PhoneVerified {

			utils.Render(w, "login.html", models.PageData{
				Title: "Login",
				Error: "Please verify your phone number before logging in.",
			})
			return
		}

		// Create login session
		middleware.CreateSession(
			w,
			r,
			user.ID,
		)

		// Go to dashboard
		http.Redirect(
			w,
			r,
			"/dashboard",
			http.StatusSeeOther,
		)

		return
	}

	// Show login page
	utils.Render(w, "login.html", models.PageData{
		Title: "Login",
	})
}
func Dashboard(w http.ResponseWriter, r *http.Request) {

	userID := middleware.GetUserID(r)

	if userID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Get user
	user, err := database.GetUserByID(userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get basic task statistics
	total, completed, pending, err :=
		database.GetTaskStats(userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get dashboard analytics
	highPriority, duetoday, completionRate, err :=
		database.GetDashboardAnalytics(userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	recentTasks, err := database.GetRecentTasks(userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Display dashboard
	utils.Render(w, "dashboard.html", models.PageData{

		Title: "Dashboard",

		Name: user.FullName,

		User: user,

		TotalTasks: total,

		CompletedTasks: completed,

		PendingTasks: pending,

		HighPriorityTasks: highPriority,

		DueTodayTasks: duetoday,

		CompletionRate: completionRate,

		Tasks: recentTasks,
	})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	middleware.DestroySession(w, r)

	http.Redirect(w, r, "/", http.StatusSeeOther)

}
